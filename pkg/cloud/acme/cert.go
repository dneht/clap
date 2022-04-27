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

package acme

import (
	"cana.io/clap/pkg/log"
	"cana.io/clap/pkg/model"
	"cana.io/clap/pkg/refer"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/alidns"
	"github.com/go-acme/lego/v4/providers/dns/dnspod"
	"github.com/go-acme/lego/v4/registration"
	"github.com/gofiber/fiber/v2"
)

func ApplyCert(task *model.Timetable, domains map[string][]string) (map[string]refer.AcmeTaskResult, error) {
	if nil == task || "" == task.TaskInfo || nil == domains || len(domains) == 0 {
		return nil, nil
	}
	var info refer.AcmeTaskInfo
	err := json.Unmarshal([]byte(task.TaskInfo), &info)
	if nil != err {
		log.Warnf("decode task info failed: %v", err)
		return nil, err
	}
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if nil != err {
		log.Warnf("generate private key failed: %v", err)
		return nil, err
	}
	user := &refer.AcmeUser{
		Email:      info.Email,
		PrivateKey: privateKey,
	}
	config := lego.NewConfig(user)
	config.Certificate.KeyType = certcrypto.RSA2048
	dnsProvider, err := getDNSProvider(task, &info)
	if nil != err || nil == dnsProvider {
		return nil, err
	}

	client, err := lego.NewClient(config)
	if nil != err {
		log.Warnf("create lego client failed: %v", err)
		return nil, err
	}
	err = client.Challenge.SetDNS01Provider(dnsProvider)
	if nil != err {
		log.Warnf("lego set dns provider failed: %v", err)
		return nil, err
	}

	register, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if nil != err {
		log.Warnf("lego register user failed: %v", err)
		return nil, err
	}
	user.Registration = register

	result := make(map[string]refer.AcmeTaskResult, len(domains))
	for key, group := range domains {
		request := certificate.ObtainRequest{
			Domains: group,
			Bundle:  true,
		}
		cert, err := client.Certificate.Obtain(request)
		if nil != err {
			log.Warnf("lego obtain cert[%v] failed: %v", group, err)
			return nil, err
		}
		_, body, errs := fiber.Get(cert.CertURL).String()
		if nil != errs && len(errs) > 0 {
			log.Warnf("get stable cert from url[%v][%v] failed: %v", group, cert.CertURL, errs)
			return nil, errs[0]
		}
		result[key] = refer.AcmeTaskResult{
			FullCert:   body,
			PrivateKey: string(cert.PrivateKey),
		}
	}
	return result, nil
}

func getDNSProvider(task *model.Timetable, info *refer.AcmeTaskInfo) (challenge.Provider, error) {
	switch task.TaskType {
	case refer.AcmeAliyunTaskType:
		if nil == info.AliyunConfig {
			return nil, nil
		}
		dnsProvider, err := getAliyunDNSProvider(info)
		if nil != err {
			log.Warnf("get aliyun dns provider failed: %v", err)
			return nil, err
		}
		return dnsProvider, nil
	case refer.AcmeDnspodTaskType:
		if nil == info.DnspodConfig {
			return nil, nil
		}
		dnsProvider, err := getAliyunDNSProvider(info)
		if nil != err {
			log.Warnf("get aliyun dns provider failed: %v", err)
			return nil, err
		}
		return dnsProvider, nil
	}
	log.Warnf("task type not match: %v", task)
	return nil, nil
}

func getAliyunDNSProvider(info *refer.AcmeTaskInfo) (*alidns.DNSProvider, error) {
	dnsConfig := alidns.NewDefaultConfig()
	dnsConfig.APIKey = info.AliyunConfig.APIKey
	dnsConfig.SecretKey = info.AliyunConfig.SecretKey
	if "" != info.AliyunConfig.SecurityToken {
		dnsConfig.SecurityToken = info.AliyunConfig.SecurityToken
	}
	return alidns.NewDNSProviderConfig(dnsConfig)
}

func getDnspodDNSProvider(info *refer.AcmeTaskInfo) (*dnspod.DNSProvider, error) {
	dnsConfig := dnspod.NewDefaultConfig()
	dnsConfig.LoginToken = info.DnspodConfig.LoginToken
	return dnspod.NewDNSProviderConfig(dnsConfig)
}
