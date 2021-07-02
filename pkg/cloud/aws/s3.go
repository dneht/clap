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

package aws

import (
	"cana.io/clap/pkg/base"
	"cana.io/clap/pkg/refer"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"sync"
)

var lock sync.Mutex
var s3Map = make(map[uint64]*session.Session)

func GetS3(envId uint64) (*session.Session, error) {
	s3Cli, ok := s3Map[envId]
	if ok {
		return s3Cli, nil
	}
	lock.Lock()
	deploy, err := base.Deploy(envId)
	if nil != err {
		return nil, err
	}
	s3Cli, s3Err := initS3Cli(&deploy.S3Info)
	if nil != s3Err {
		return nil, s3Err
	}
	s3Map[envId] = s3Cli
	lock.Unlock()
	return s3Cli, nil
}

func ResetS3() {
	lock.Lock()
	s3Map = make(map[uint64]*session.Session)

	lock.Unlock()
}

func initS3Cli(conf *refer.S3Conf) (*session.Session, error) {
	creds := credentials.NewStaticCredentials(conf.AccessKey, conf.SecretKey, "")
	config := &aws.Config{
		Region:           aws.String(conf.Region),
		Endpoint:         aws.String(conf.Endpoint),
		S3ForcePathStyle: aws.Bool(conf.ForcePath),
		DisableSSL:       aws.Bool(conf.DisableSSL),
		Credentials:      creds,
	}
	client, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}
