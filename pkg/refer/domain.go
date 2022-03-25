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

package refer

type UpdateDomainSecret struct {
	EnvId      uint64 `json:"envId"`
	Domain     string `json:"domain"`
	MainDomain string `json:"mainDomain"`
	Namespace  string `json:"namespace"`
	SecretName string `json:"secretName"`
}

type ExistDomainSecret struct {
	EnvId      uint64            `json:"envId"`
	Namespace  string            `json:"namespace"`
	SecretName string            `json:"secretName"`
	SecretData map[string]string `json:"secretData"`
}

type UpdateDomainResult struct {
	Secrets *[]UpdateDomainSecret      `json:"secrets"`
	Results *map[string]AcmeTaskResult `json:"results"`
}
