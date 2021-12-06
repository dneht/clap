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

import (
	"crypto"
	"github.com/go-acme/lego/v4/registration"
)

const (
	AcmeAliyunTaskType = "acme-aliyun"
	AcmeDnspodTaskType = "acme-dnspod"
)

type AcmeUser struct {
	Email        string
	Registration *registration.Resource
	PrivateKey   crypto.PrivateKey
}

func (u *AcmeUser) GetEmail() string {
	return u.Email
}
func (u AcmeUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *AcmeUser) GetPrivateKey() crypto.PrivateKey {
	return u.PrivateKey
}

type AcmeTaskInfo struct {
	Email        string            `json:"email"`
	AliyunConfig *AcmeAliyunConfig `json:"aliyun"`
	DnspodConfig *AcmeDnspodConfig `json:"dnspod"`
}

type AcmeAliyunConfig struct {
	APIKey        string `json:"apiKey"`
	SecretKey     string `json:"secretKey"`
	SecurityToken string `json:"securityToken"`
}

type AcmeDnspodConfig struct {
	LoginToken string `json:"loginToken"`
}

type AcmeTaskResult struct {
	FullCert   string `json:"fullCert"`
	PrivateKey string `json:"privateKey"`
}