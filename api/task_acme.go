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
	"cana.io/clap/pkg/cloud/acme"
	"cana.io/clap/pkg/log"
	"cana.io/clap/pkg/model"
	"cana.io/clap/pkg/refer"
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"strings"
	"time"
)

func execTaskAcmeDomains(task *model.Timetable) (interface{}, error) {
	secrets, domains := allNeedAcmeDomains()
	log.Infof("this time secrets is: %v, domains is: %v", secrets, domains)
	if nil == secrets || nil == domains || len(*secrets) == 0 || len(*domains) == 0 {
		log.Infof("no domains found to execute: %v, %v", secrets, domains)
		return nil, nil
	}

	var domainResult refer.UpdateDomainResult
	needExec := false
	latest, err := getLatestTimetableResult(task.Id)
	if nil == latest || nil != err {
		if nil == latest {
			needExec = true
		}
		log.Warnf("get task latest result empty or failed: %v, continue is %v", err, needExec)
	} else {
		log.Infof("get task latest result status: %v, time: %v", latest.LastStatus, latest.CreatedAt)
		if latest.LastStatus == 0 && "" != latest.LastResult {
			err = json.Unmarshal([]byte(latest.LastResult), &domainResult)
			if nil != err {
				log.Warnf("decode task latest result failed: %v, continue execute", err)
			}
			if nil != domainResult.Secrets && nil != domainResult.Results && len(*domainResult.Secrets) > 0 && len(*domainResult.Results) > 0 {
				secretMap := make(map[string]bool, len(*domainResult.Secrets))
				domainMap := make(map[string]bool, len(*domainResult.Results))
				for _, secret := range *domainResult.Secrets {
					secretMap[strconv.FormatUint(secret.EnvId, 10)+secret.Namespace+secret.SecretName] = true
				}
				for _, secret := range *secrets {
					_, ok := secretMap[strconv.FormatUint(secret.EnvId, 10)+secret.Namespace+secret.SecretName]
					if !ok {
						log.Infof("secret is change: %v -> %v", domainResult.Secrets, secrets)
						needExec = true
						break
					}
				}
				if !needExec {
					for main := range *domainResult.Results {
						domainMap[main] = true
					}
					for main := range *domains {
						_, ok := domainMap[main]
						if !ok {
							log.Infof("domain is change: %v -> %v", domainResult.Results, domains)
							needExec = true
							break
						}
					}
				}
			} else {
				needExec = true
			}
			if !needExec {
				if latest.CreatedAt.Before(time.Now().Add(-60 * 24 * time.Hour)) {
					log.Infof("latest is too old: %v", latest)
					needExec = true
				}
			}
		} else {
			needExec = true
		}
	}
	if needExec {
		results, err := acme.ApplyCert(task, domains)
		if nil != err {
			log.Warnf("generate acme cert failed: %v, skip execute", err)
			return nil, err
		}
		for _, secret := range *secrets {
			if "" != secret.SecretName && "" != secret.Domain && "" != secret.MainDomain {
				result, ok := (*results)[secret.MainDomain]
				if !ok {
					log.Warnf("can not get cert result: %v, %v", secret, results)
					continue
				}
				cli, _, err := base.K8S(secret.EnvId)
				if nil != err {
					log.Warnf("can not client cluster: %v, %v", secret, err)
					continue
				}
				err = createOrUpdateSecret(cli, secret.Namespace, secret.SecretName, &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name: secret.SecretName,
					},
					Data: map[string][]byte{
						corev1.TLSCertKey:       []byte(result.FullCert),
						corev1.TLSPrivateKeyKey: []byte(result.PrivateKey),
					},
					Type: corev1.SecretTypeTLS,
				})
				if nil != err {
					log.Warnf("create or update secret failed: %v, %v", secret, err)
				}
			}
		}
		domainResult = refer.UpdateDomainResult{
			Secrets: secrets,
			Results: results,
		}
		return &domainResult, nil
	} else {
		return nil, nil
	}
}

func allNeedAcmeDomains() (*[]refer.UpdateDomainSecret, *map[string][]string) {
	spaces, err := findAllSpaceForTask()
	if nil != err {
		log.Warnf("get space for task failed: %v", err)
		return nil, nil
	}

	updates := make([]refer.UpdateDomainSecret, 0, 32)
	domains := make(map[string][]string)
	updateDup := make(map[string]bool)
	domainDup := make(map[string]bool)
	wildcardDup := make(map[string]bool)
	for _, space := range *spaces {
		var info refer.SpaceRealInfo
		err = json.Unmarshal([]byte(space.SpaceInfo), &info)
		if nil != err {
			log.Warnf("decode space info for task failed: %v", err)
			continue
		}
		param := info.Param
		if "" != param.Domain && "" != param.TLS.SecretName {
			domain := strings.ToLower(strings.TrimSpace(param.Domain))
			split := strings.Split(domain, ".")
			if len(split) <= 1 {
				continue
			}

			var main string
			if len(split) == 2 {
				main = domain
			} else {
				main = strings.Join(split[len(split)-2:], ".")
				if strings.HasPrefix(domain, "-") {
					domain = strings.Join(split[1:], ".")
				}
			}

			var ok bool
			updateKey := strconv.FormatUint(space.EnvId, 10) + space.SpaceKeep + param.TLS.SecretName
			_, ok = updateDup[updateKey]
			if !ok {
				updateDup[updateKey] = true
				updates = append(updates, refer.UpdateDomainSecret{
					EnvId:      space.EnvId,
					Domain:     domain,
					MainDomain: main,
					Namespace:  space.SpaceKeep,
					SecretName: param.TLS.SecretName,
				})
			}
			var list []string
			_, ok = domainDup[domain]
			if !ok {
				nextDomain := "*." + domain
				domainDup[domain] = true
				wildcardDup[nextDomain] = true
				list, ok = domains[main]
				if !ok {
					list = make([]string, 0, 4)
				}
				list = append(list, domain, nextDomain)
				domains[main] = list
			}
		}
	}
	//wildcard overlap
	if len(wildcardDup) > 0 {
		for main, list := range domains {
			nList := make([]string, 0, len(list))
			for _, domain := range list {
				split := strings.Split(domain, ".")
				if len(split) < 2 {
					continue
				}
				if split[0] == "*" {
					nList = append(nList, domain)
				} else {
					wildDomain := "*." + strings.Join(split[1:], ".")
					_, ok := wildcardDup[wildDomain]
					if !ok {
						nList = append(nList, domain)
					}
				}
			}
			domains[main] = nList
		}
	}
	return &updates, &domains
}
