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

package api

import (
	"cana.io/clap/pkg/base"
	"cana.io/clap/pkg/refer"
	"context"
	"encoding/json"
	"errors"
	corev1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func handleDeployment(envId uint64, namespace, jsonStr, appDeployName string) (refer.DeployStatus, error) {
	k8s, _, err := base.K8S(envId)
	if nil != err {
		return refer.ConnectErrorDeployStatus, err
	}

	var median = new(refer.Deployment)
	err = json.Unmarshal([]byte(jsonStr), median)
	if nil != err {
		return refer.FailedDeployStatus, err
	}
	get, err := k8s.AppsV1().Deployments(namespace).Get(context.TODO(), appDeployName, metav1.GetOptions{})
	if nil != err {
		if k8serror.IsNotFound(err) {
			get = nil
			err = nil
		} else {
			return refer.FailedDeployStatus, err
		}
	}
	if nil == get {
		get, err = k8s.AppsV1().Deployments(namespace).Create(context.TODO(), median.Main, metav1.CreateOptions{})
	} else {
		get, err = k8s.AppsV1().Deployments(namespace).Update(context.TODO(), median.Main, metav1.UpdateOptions{})
	}
	if nil != err {
		return refer.FailedDeployStatus, err
	}

	status := refer.DefaultDeployStatus
	status, err = handleService(envId, &median.CommonExtend, namespace, status)
	status, err = handleContour(envId, &median.CommonExtend, namespace, status)
	status, err = handleSecret(envId, &median.CommonExtend, namespace, status)
	status, err = handleConfig(envId, &median.CommonExtend, namespace, status)
	return status, err
}

func handleStatefulSet(envId uint64, namespace, jsonStr, appDeployName string) (refer.DeployStatus, error) {
	k8s, _, err := base.K8S(envId)
	if nil != err {
		return refer.ConnectErrorDeployStatus, err
	}

	var median = new(refer.StatefulSet)
	err = json.Unmarshal([]byte(jsonStr), median)
	if nil != err {
		return refer.FailedDeployStatus, err
	}
	get, err := k8s.AppsV1().StatefulSets(namespace).Get(context.TODO(), appDeployName, metav1.GetOptions{})
	if nil != err {
		if k8serror.IsNotFound(err) {
			get = nil
			err = nil
		} else {
			return refer.FailedDeployStatus, err
		}
	}
	if nil == get {
		get, err = k8s.AppsV1().StatefulSets(namespace).Create(context.TODO(), median.Main, metav1.CreateOptions{})
	} else {
		get, err = k8s.AppsV1().StatefulSets(namespace).Update(context.TODO(), median.Main, metav1.UpdateOptions{})
	}
	if nil != err {
		return refer.FailedDeployStatus, err
	}

	status := refer.DefaultDeployStatus
	status, err = handleService(envId, &median.CommonExtend, namespace, status)
	status, err = handleContour(envId, &median.CommonExtend, namespace, status)
	status, err = handleSecret(envId, &median.CommonExtend, namespace, status)
	status, err = handleConfig(envId, &median.CommonExtend, namespace, status)
	return status, err
}

func handleService(envId uint64, extend *refer.CommonExtend, namespace string, status refer.DeployStatus) (refer.DeployStatus, error) {
	if nil == extend.Services {
		return status, nil
	}
	serviceList := *extend.Services
	if len(serviceList) <= 0 {
		return status, nil
	}

	var err error
	var k8s *kubernetes.Clientset
	k8s, _, err = base.K8S(envId)
	if nil != err {
		return refer.ConnectErrorDeployStatus, err
	}

	for serviceName, serviceData := range serviceList {
		if nil == serviceData {
			continue
		}
		err = createOrUpdateService(k8s, namespace, serviceName, serviceData)
		if nil != err {
			status = refer.UnknownErrorDeployStatus
			return status, err
		}
	}

	return status, err
}

func createOrUpdateService(k8s *kubernetes.Clientset, namespace string, serviceName string, serviceData *corev1.Service) error {
	if nil == k8s {
		return errors.New("could not get cluster info")
	}
	get, err := k8s.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if nil != err {
		if k8serror.IsNotFound(err) {
			get = nil
			err = nil
		} else {
			return err
		}
	}
	if nil == get {
		get, err = k8s.CoreV1().Services(namespace).Create(context.TODO(), serviceData, metav1.CreateOptions{})
	} else {
		get, err = k8s.CoreV1().Services(namespace).Update(context.TODO(), serviceData, metav1.UpdateOptions{})
	}
	return err
}

func handleContour(envId uint64, extend *refer.CommonExtend, namespace string, status refer.DeployStatus) (refer.DeployStatus, error) {
	if nil == extend.Contours {
		return status, nil
	}
	contourList := *extend.Contours
	if len(contourList) <= 0 {
		return status, nil
	}

	crd, _, err := base.K8D(envId)
	if nil != err {
		return refer.ConnectErrorDeployStatus, err
	}

	for contourName, contourData := range contourList {
		if nil == contourData {
			continue
		}
		err = createOrUpdateContour(crd, namespace, contourName, contourData)
		if nil != err {
			status = refer.UnknownErrorDeployStatus
			return status, err
		}
	}

	return status, err
}

func createOrUpdateContour(crd dynamic.Interface, namespace string, contourName string, contourData *unstructured.Unstructured) error {
	if nil == crd {
		return errors.New("could not get cluster info")
	}
	get, err := crd.Resource(refer.ContourGvr).Namespace(namespace).Get(context.TODO(), contourName, metav1.GetOptions{})
	if nil != err {
		if k8serror.IsNotFound(err) {
			get = nil
			err = nil
		} else {
			return err
		}
	}

	if nil == get {
		get, err = crd.Resource(refer.ContourGvr).Namespace(namespace).Create(context.TODO(), contourData, metav1.CreateOptions{})
	} else {
		contourData.SetResourceVersion(get.GetResourceVersion())
		get, err = crd.Resource(refer.ContourGvr).Namespace(namespace).Update(context.TODO(), contourData, metav1.UpdateOptions{})
	}
	return err
}

func handleSecret(envId uint64, extend *refer.CommonExtend, namespace string, status refer.DeployStatus) (refer.DeployStatus, error) {
	if nil == extend.Secrets {
		return status, nil
	}
	secretList := *extend.Secrets
	if len(secretList) <= 0 {
		return status, nil
	}

	var err error
	var k8s *kubernetes.Clientset
	k8s, _, err = base.K8S(envId)
	if nil != err {
		return refer.ConnectErrorDeployStatus, err
	}

	for secretName, secretData := range secretList {
		if nil == secretData {
			continue
		}
		err = createOrUpdateSecret(k8s, namespace, secretName, secretData)
		if nil != err {
			status = refer.UnknownErrorDeployStatus
			return status, err
		}
	}

	return status, err
}

func createOrUpdateSecret(k8s *kubernetes.Clientset, namespace string, secretName string, secretData *corev1.Secret) error {
	if nil == k8s {
		return errors.New("could not get cluster info")
	}
	get, err := k8s.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if nil != err {
		if k8serror.IsNotFound(err) {
			get = nil
			err = nil
		} else {
			return err
		}
	}
	if nil == get {
		get, err = k8s.CoreV1().Secrets(namespace).Create(context.TODO(), secretData, metav1.CreateOptions{})
	} else {
		get, err = k8s.CoreV1().Secrets(namespace).Update(context.TODO(), secretData, metav1.UpdateOptions{})
	}
	return err
}

func handleConfig(envId uint64, extend *refer.CommonExtend, namespace string, status refer.DeployStatus) (refer.DeployStatus, error) {
	if nil == extend.Configs {
		return status, nil
	}
	configList := *extend.Configs
	if len(configList) <= 0 {
		return status, nil
	}

	var err error
	var k8s *kubernetes.Clientset
	k8s, _, err = base.K8S(envId)
	if nil != err {
		return refer.ConnectErrorDeployStatus, err
	}

	for configName, configData := range configList {
		if nil == configData {
			continue
		}
		err = createOrUpdateConfig(k8s, namespace, configName, configData)
		if nil != err {
			status = refer.UnknownErrorDeployStatus
			return status, err
		}
	}

	return status, err
}

func createOrUpdateConfig(k8s *kubernetes.Clientset, namespace string, configName string, configData *corev1.ConfigMap) error {
	if nil == k8s {
		return errors.New("could not get cluster info")
	}
	get, err := k8s.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configName, metav1.GetOptions{})
	if nil != err {
		if k8serror.IsNotFound(err) {
			get = nil
			err = nil
		} else {
			return err
		}
	}
	if nil == get {
		get, err = k8s.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configData, metav1.CreateOptions{})
	} else {
		get, err = k8s.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configData, metav1.UpdateOptions{})
	}
	return err
}

func onlyUpdateConfig(k8s *kubernetes.Clientset, namespace string, configName string, configData *corev1.ConfigMap) error {
	if nil == k8s {
		return errors.New("could not get cluster info")
	}
	get, err := k8s.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configName, metav1.GetOptions{})
	if nil != err {
		return err
	}
	if nil == get {
		return errors.New("config not exist")
	} else {
		get, err = k8s.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configData, metav1.UpdateOptions{})
	}
	return err
}

//TODO handleBudget
func handleBudget() {

}

//TODO handlePolicy
func handlePolicy() {

}
