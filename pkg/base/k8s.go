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
	"cana.io/clap/pkg/refer"
	"encoding/base64"
	"k8s.io/client-go/dynamic"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func K8SCli(conf *refer.K8SConf) (*rest.Config, *kubernetes.Clientset, error) {
	config, err := buildConfig(conf)
	if err != nil {
		return nil, nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	return config, client, err
}

func K8SDynamic(conf *refer.K8SConf) (*rest.Config, dynamic.Interface, error) {
	config, err := buildConfig(conf)
	if err != nil {
		return nil, nil, err
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	return config, client, err
}

func buildConfig(conf *refer.K8SConf) (*rest.Config, error) {
	if conf.Inner {
		return rest.InClusterConfig()
	} else {
		K8sConfigGetter := func() (*clientcmdapi.Config, error) {
			apiConfig := new(clientcmdapi.Config)
			apiConfig.APIVersion = conf.Version
			apiConfig.Kind = "Config"
			apiConfig.Contexts = map[string]*clientcmdapi.Context{
				conf.ClusterName: {
					Cluster:  conf.ClusterName,
					AuthInfo: conf.ClusterUser,
				},
			}
			apiConfig.CurrentContext = conf.ClusterName
			apiConfig.Clusters = map[string]*clientcmdapi.Cluster{
				conf.ClusterName: {
					Server:                conf.Master,
					InsecureSkipTLSVerify: true,
				},
			}
			clientCert, clientCertErr := base64.StdEncoding.DecodeString(conf.ClientCert)
			if nil != clientCertErr {
				return nil, clientCertErr
			}
			clientKey, clientKeyErr := base64.StdEncoding.DecodeString(conf.ClientKey)
			if nil != clientKeyErr {
				return nil, clientKeyErr
			}
			apiConfig.AuthInfos = map[string]*clientcmdapi.AuthInfo{
				conf.ClusterUser: {
					ClientCertificateData: clientCert,
					ClientKeyData:         clientKey,
				},
			}
			return apiConfig, nil
		}
		return clientcmd.BuildConfigFromKubeconfigGetter("", K8sConfigGetter)
	}
}