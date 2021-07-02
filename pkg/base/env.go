/*
Copyright 2020 Dasheng.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package base

import (
	"encoding/json"
	"errors"
	"k8s.io/client-go/dynamic"
	"strconv"
	"sync"

	"cana.io/clap/pkg/model"
	"cana.io/clap/pkg/refer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var lock sync.Mutex
var k8sMap = make(map[uint64]*kubernetes.Clientset)
var k8sConfMap = make(map[uint64]*rest.Config)
var crdMap = make(map[uint64]dynamic.Interface)

var deployMap = make(map[uint64]*refer.DeployInfo)
var envIdMap = make(map[uint64]*model.Environment)
var envNameMap = make(map[string]uint64)

func EnvInit() {
	list := dangListFullEnv()
	for _, one := range *list {
		val := one
		envIdMap[one.Id] = &val
		envNameMap[one.Env] = one.Id
	}
}

func Env(envId uint64) (*model.Environment, error) {
	envInfo, getOk := envIdMap[envId]
	if getOk {
		return envInfo, nil
	}
	return nil, errors.New("env " + strconv.FormatUint(envId, 10) + " not exist")
}

func Deploy(envId uint64) (*refer.DeployInfo, error) {
	deploy, ok := deployMap[envId]
	if ok {
		return deploy, nil
	}
	envInfo, envErr := Env(envId)
	if nil != envErr {
		return nil, envErr
	}
	deployErr := json.Unmarshal([]byte(envInfo.DeployInfo), &deploy)
	if nil != deployErr {
		return nil, deployErr
	}
	deployMap[envId] = deploy
	return deploy, nil
}

func K8S(envId uint64) (*kubernetes.Clientset, *rest.Config, error) {
	k8sCli, ok := k8sMap[envId]
	k8sConf := k8sConfMap[envId]
	if ok {
		return k8sCli, k8sConf, nil
	}
	lock.Lock()
	deploy, err := Deploy(envId)
	if nil != err {
		return nil, nil, err
	}
	k8sConf, k8sCli, k8sErr := K8SCli(&deploy.K8SInfo)
	if nil != k8sErr {
		return nil, nil, k8sErr
	}
	k8sMap[envId] = k8sCli
	k8sConfMap[envId] = k8sConf
	lock.Unlock()
	return k8sCli, k8sConf, nil
}

func K8D(envId uint64) (dynamic.Interface, *rest.Config, error) {
	crdCli, ok := crdMap[envId]
	k8sConf := k8sConfMap[envId]
	if ok {
		return crdCli, k8sConf, nil
	}
	lock.Lock()
	deploy, err := Deploy(envId)
	if nil != err {
		return nil, nil, err
	}
	k8sConf, crdCli, k8sErr := K8SDynamic(&deploy.K8SInfo)
	if nil != k8sErr {
		return nil, nil, k8sErr
	}
	crdMap[envId] = crdCli
	k8sConfMap[envId] = k8sConf
	lock.Unlock()
	return crdCli, k8sConf, nil
}

func Reset() {
	lock.Lock()
	k8sMap = make(map[uint64]*kubernetes.Clientset)
	k8sConfMap = make(map[uint64]*rest.Config)
	crdMap = make(map[uint64]dynamic.Interface)

	deployMap = make(map[uint64]*refer.DeployInfo)
	envIdMap = make(map[uint64]*model.Environment)
	EnvInit()
	lock.Unlock()
}

func dangListFullEnv() *[]model.Environment {
	var list []model.Environment
	err := Engine.Find(&list)
	if nil != err {
		panic(err)
	}
	return &list
}
