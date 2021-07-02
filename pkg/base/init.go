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
	"flag"
)

var nowSeq uint64
var nowProp *Property

var seqFlag = flag.Uint64("seq", 1, "Sequence")
var envFlag = flag.String("env", "dev", "Environment")
var confFlag = flag.String("conf", "", "Config file path")
var namespaceFlag = flag.String("namespace", "clap-system", "The namespace where the app & job is located")
var timezoneFlag = flag.String("timezone", "Local", "App timezone, please use go format")

func Init() {
	flag.Parse()
	DbInit()
	WebInit()
	SeqInit()
	EnvInit()
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
