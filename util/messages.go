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

package util

import (
	"bytes"
	"cana.io/clap/pkg/base"
	"cana.io/clap/pkg/log"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func DingDingDeployMessage(env, space, app, desc string) {
	DingDingMessage(map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": fmt.Sprintf("正在[%s]环境的[%s]空间发布[%s|%s]项目", env, space, app, desc),
		},
	})
}

func DingDingMessage(ding map[string]interface{}) {
	dingProp := base.Now().Message.DingDing
	if !dingProp.Enable {
		return
	}

	dingJson, err := json.Marshal(ding)
	if nil != err {
		log.Warnf("dingding to json error: %v, %v\n", ding, err)
		return
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Post(dingProp.ApiUrl,
		"application/json", bytes.NewBuffer(dingJson))
	if nil != err {
		log.Warnf("dingding send error: %v, %v\n", resp, err)
		return
	}
	defer resp.Body.Close()
}
