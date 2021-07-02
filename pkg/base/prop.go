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
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

type Database struct {
	Addresses  []string `json:"addr"`
	User       string   `json:"user"`
	Password   string   `json:"password"`
	Database   string   `json:"database"`
	MaxIdleCon int      `json:"max_idle_con"`
	MaxOpenCon int      `json:"max_open_con"`
}

type Property struct {
	Env       string                  `json:"env"`
	Namespace string                  `json:"namespace"`
	Timezone  string                  `json:"timezone"`
	Service   ServiceProp             `json:"service"`
	Message   MessageProp             `json:"message"`
	Document  map[string]DocumentProp `json:"document"`
	Package   PackageProp             `json:"package"`
}

type ServiceProp struct {
	Port       string                  `json:"port"`
	Debug      bool                    `json:"debug"`
	//Cros       middleware.CORSConfig   `json:"cros"`
	//Csrf       middleware.CSRFConfig   `json:"csrf"`
	//Secure     middleware.SecureConfig `json:"secure"`
	GzipLevel  int                     `json:"gzip_level"`
	StaticPath string                  `json:"static_path"`
	Password   PasswordProp            `json:"password"`
}

type MessageProp struct {
	DingDing MessageDingDingProp `json:"dingding"`
}

type MessageDingDingProp struct {
	Enable bool   `json:"enable"`
	ApiUrl string `json:"api_url"`
}

type DocumentProp struct {
	Enable bool   `json:"enable"`
	Token  string `json:"token"`
	ApiUrl string `json:"api_url"`
}

type PasswordProp struct {
	Type string `json:"type"`
}

type PackageProp struct {
	BuildJobImage      string `json:"build_job_image"`
	ImagePullPolicy    string `json:"image_pull_policy"`
	ImagePullSecret    string `json:"image_pull_secret"`
	MavenSkipTests     bool   `json:"maven_skip_tests"`
	BackoffLimit       int32  `json:"backoff_limit"`
	CleanAfterFinished int32  `json:"clean_after_finished"`
}

func DatabaseConf() *Database {
	var database *Database
	confPath := *confFlag
	if "" == confPath {
		addr := parseList(os.Getenv("CLAP_DB_ADDR"))
		if len(addr) > 0 {
			database = &Database{
				Addresses:  addr,
				User:       parseString(os.Getenv("CLAP_DB_USER"), "root"),
				Password:   parseString(os.Getenv("CLAP_DB_PASSWORD"), ""),
				Database:   parseString(os.Getenv("CLAP_DB_DATABASE"), "clap"),
				MaxIdleCon: parseInt(os.Getenv("CLAP_DB_MAX_IDLE_CON"), 50),
				MaxOpenCon: parseInt(os.Getenv("CLAP_DB_MAX_OPEN_CON"), 100),
			}
		}
	} else {
		fileData, err := ioutil.ReadFile(confPath)
		if nil != err {
			panic(err)
		}
		allData := string(fileData)
		var lines []string
		if strings.Contains(allData, "\n") {
			lines = strings.Split(allData, "\n")
		} else if strings.Contains(allData, "\r") {
			lines = strings.Split(allData, "\r")
		}
		if nil != lines {
			kv := make(map[string]string, 6)
			for _, one := range lines {
				data := strings.TrimSpace(one)
				idx := strings.Index(data, "=")
				if idx > 0 {
					key := strings.TrimSpace(data[:idx])
					value := strings.TrimSpace(data[idx+1:])
					kv[key] = value
				}
			}
			database = &Database{
				Addresses:  parseList(kv["database.addr"]),
				User:       parseString(kv["database.user"], "root"),
				Password:   parseString(kv["database.password"], ""),
				Database:   parseString(kv["database.database"], "clap"),
				MaxIdleCon: parseInt(kv["database.max_idle_con"], 50),
				MaxOpenCon: parseInt(kv["database.max_open_con"], 100),
			}
		}
	}
	if nil == database {
		panic(errors.New("database not exist"))
	}
	return database
}

func BuildProperty() {
	propList := dangListFullProp()
	propMap := make(map[string]string, len(*propList))
	for _, prop := range *propList {
		propMap[strings.TrimSpace(prop.Prop)] = strings.TrimSpace(prop.Value)
	}

	nowProp = &Property{
		Env:       *envFlag,
		Namespace: *namespaceFlag,
		Timezone:  *timezoneFlag,
		Service: ServiceProp{
			Port:  parseString(propMap["service.port"], "8008"),
			Debug: parseBool(propMap["service.debug"], false),
			//Cros: middleware.CORSConfig{
			//	AllowOrigins:     parseListElse(propMap["service.cros.allow_origins"], []string{"*"}),
			//	AllowMethods:     parseListElse(propMap["service.cros.allow_methods"], []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			//	AllowHeaders:     parseList(propMap["service.cros.allow_headers"]),
			//	AllowCredentials: parseBool(propMap["service.cros.allow_credentials"], false),
			//	ExposeHeaders:    parseList(propMap["service.cros.expose_headers"]),
			//	MaxAge:           parseInt(propMap["service.cros.max_age"], 0),
			//},
			//Csrf: middleware.CSRFConfig{
			//	TokenLength:    uint8(parseInt(propMap["service.csrf.token_length"], 32)),
			//	TokenLookup:    parseString(propMap["service.csrf.token_lookup"], "header:X-CSRF-Token"),
			//	ContextKey:     parseString(propMap["service.csrf.context_key"], "csrf"),
			//	CookieName:     parseString(propMap["service.csrf.cookie_name"], "csrf"),
			//	CookieDomain:   parseString(propMap["service.csrf.cookie_domain"], ""),
			//	CookiePath:     parseString(propMap["service.csrf.cookie_path"], ""),
			//	CookieMaxAge:   parseInt(propMap["service.csrf.cookie_max_age"], 86400),
			//	CookieSecure:   parseBool(propMap["service.csrf.cookie_secure"], false),
			//	CookieHTTPOnly: parseBool(propMap["service.csrf.cookie_http_only"], false),
			//},
			//Secure: middleware.SecureConfig{
			//	XSSProtection:         parseString(propMap["service.secure.xss_protection"], "1; mode=block"),
			//	ContentTypeNosniff:    parseString(propMap["service.secure.content_type_nosniff"], "nosniff"),
			//	XFrameOptions:         parseString(propMap["service.secure.x_frame_options"], "SAMEORIGIN"),
			//	ContentSecurityPolicy: parseString(propMap["service.secure.content_security_policy"], ""),
			//	ReferrerPolicy:        parseString(propMap["service.secure.referrer_policy"], "origin"),
			//	CSPReportOnly:         parseBool(propMap["service.secure.csp_report_only"], false),
			//	HSTSMaxAge:            parseInt(propMap["service.secure.hsts_max_age"], 0),
			//	HSTSPreloadEnabled:    parseBool(propMap["service.secure.hsts_preload_enabled"], false),
			//	HSTSExcludeSubdomains: parseBool(propMap["service.secure.hsts_exclude_subdomains"], false),
			//},
			GzipLevel:  parseInt(propMap["service.gzip_level"], 6),
			StaticPath: parseString(propMap["service.static_path"], "/opt/ui"),
			Password: PasswordProp{
				Type: parseString(propMap["service.password.type"], "ssha"),
			},
		},
		Message: MessageProp{
			DingDing: MessageDingDingProp{
				Enable: parseBool(propMap["message.dingding.enable"], false),
				ApiUrl: parseString(propMap["message.dingding.api_url"], ""),
			},
		},
		Document: parseDocument(propMap),
		Package: PackageProp{
			BuildJobImage:      parseString(propMap["package.build_job_image"], "dneht/clap-build:1.0.0"),
			ImagePullPolicy:    parseString(propMap["package.image_pull_policy"], "Always"),
			ImagePullSecret:    parseString(propMap["package.image_pull_secret"], ""),
			MavenSkipTests:     parseBool(propMap["package.maven_skip_tests"], true),
			BackoffLimit:       int32(parseInt(propMap["package.backoff_limit"], 2)),
			CleanAfterFinished: int32(parseInt(propMap["package.clean_after_finished"], 3*24*60*60)),
		},
	}
}
