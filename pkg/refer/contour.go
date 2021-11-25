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
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const ContourApiVersion = "projectcontour.io/v1"

const ContourPlural = "httpproxies"

var ContourGvr = schema.GroupVersionResource{
	Group:    "projectcontour.io",
	Version:  "v1",
	Resource: ContourPlural,
}

type ContourDataInfo struct {
	Fqdn       string                `json:"fqdn,omitempty"`
	Tls        ContourTlsInfo        `json:"tls,omitempty"`
	CorsPolicy ContourCorsInfo       `json:"corsPolicy,omitempty"`
	Includes   []ContourIncludeInfo  `json:"includes,omitempty"`
	Routers    []ContourRouterInfo   `json:"routes,omitempty"`
	Tcpproxy   []ContourTcpproxyInfo `json:"tcpproxy,omitempty"`
}

type ContourTlsInfo struct {
	SecretName  string `json:"secretName,omitempty"`
	Passthrough bool   `json:"passthrough,omitempty"`
}

type ContourIncludeInfo struct {
	Name       string               `json:"name,omitempty"`
	Namespace  string               `json:"namespace,omitempty"`
	Conditions ContourConditionInfo `json:"conditions,omitempty"`
}

/**
LoadBalanceStrategy: RoundRobin, WeightedLeastRequest, Cookie
*/
type ContourRouterInfo struct {
	Conditions          ContourConditionInfo   `json:"conditions,omitempty"`
	Services            []ContourServiceInfo   `json:"services,omitempty"`
	InnerService        ContourServiceInfo     `json:"innerService,omitempty"`
	RespTimeout         string                 `json:"respTimeout,omitempty"`
	IdleTimeout         string                 `json:"idleTimeout,omitempty"`
	RetryCount          uint64                 `json:"retryCount,omitempty"`
	RetryPerTryTimeout  string                 `json:"retryPerTryTimeout,omitempty"`
	EnableWebsockets    bool                   `json:"enableWebsockets,omitempty"`
	LoadBalanceStrategy string                 `json:"loadBalanceStrategy,omitempty"`
	HealthCheckPolicy   ContourHealthCheckInfo `json:"healthCheckPolicy,omitempty"`
}

type ContourTcpproxyInfo struct {
	Services []ContourServiceInfo `json:"services,omitempty"`
}

type ContourConditionInfo struct {
	Prefix []string            `json:"prefix,omitempty"`
	Header []ContourHeaderInfo `json:"header,omitempty"`
}

type ContourHeaderInfo struct {
	Name     string `json:"name,omitempty"`
	Contains string `json:"contains,omitempty"`
}

type ContourServiceInfo struct {
	Name   string `json:"name,omitempty"`
	Port   int32  `json:"port,omitempty"`
	Weight uint32 `json:"weight,omitempty"`
}

type ContourCorsInfo struct {
	AllowCredentials bool     `json:"allowCredentials,omitempty"`
	AllowOrigin      []string `json:"allowOrigin,omitempty"`
	AllowMethods     []string `json:"allowMethods,omitempty"`
	AllowHeaders     []string `json:"allowHeaders,omitempty"`
	ExposeHeaders    []string `json:"exposeHeaders,omitempty"`
	MaxAge           string   `json:"maxAge,omitempty"`
}

type ContourHealthCheckInfo struct {
	Path                    string `json:"path,omitempty"`
	IntervalSeconds         uint64 `json:"intervalSeconds,omitempty"`
	TimeoutSeconds          uint64 `json:"timeoutSeconds,omitempty"`
	UnhealthyThresholdCount uint64 `json:"unhealthyThresholdCount,omitempty"`
	HealthyThresholdCount   uint64 `json:"healthyThresholdCount,omitempty"`
}
