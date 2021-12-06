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
	"cana.io/clap/pkg/log"
	"cana.io/clap/pkg/model"
	"cana.io/clap/pkg/refer"
	"flag"
	"go.uber.org/zap/zapcore"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var nowSeq uint64
var nowProp *Property

var seqFlag = flag.Uint64("seq", 1, "Sequence")
var envFlag = flag.String("env", "dev", "Environment")
var confFlag = flag.String("conf", "", "Config file path")
var namespaceFlag = flag.String("namespace", "clap-system", "The namespace where the app & job is located")
var timezoneFlag = flag.String("timezone", "", "App timezone, please use go format")

var k8sMap = make(map[uint64]*kubernetes.Clientset)
var k8sConfMap = make(map[uint64]*rest.Config)
var crdMap = make(map[uint64]dynamic.Interface)

var deployMap = make(map[uint64]*refer.DeployInfo)
var envIdMap = make(map[uint64]*model.Environment)
var envNameMap = make(map[string]uint64)
var resIdMap = make(map[uint64]*model.Resource)
var resNameMap = make(map[string]uint64)
var resInfoMap = make(map[uint64]*map[string]interface{})

func Init() {
	log.Init(zapcore.DebugLevel)
	flag.Parse()
	DbInit()
	WebInit()
	SeqInit()
	EnvInit()
	ResInit()
}

func Seq() uint64 {
	return nowSeq
}

func Now() *Property {
	return nowProp
}

func IsOffline() bool {
	return *envFlag == "dev"
}

func Reset() {
	lock.Lock()
	defer lock.Unlock()

	k8sMap = make(map[uint64]*kubernetes.Clientset)
	k8sConfMap = make(map[uint64]*rest.Config)
	crdMap = make(map[uint64]dynamic.Interface)

	deployMap = make(map[uint64]*refer.DeployInfo)
	envIdMap = make(map[uint64]*model.Environment)
	envNameMap = make(map[string]uint64)
	EnvInit()

	resIdMap = make(map[uint64]*model.Resource)
	resNameMap = make(map[string]uint64)
	resInfoMap = make(map[uint64]*map[string]interface{})
	ResInit()
}
