function(app={})

assert "id" in app : "must set app id";
assert "uk" in app : "must set app uk";
assert "name" in app : "must set app name";
assert "type" in app : "must set app type";
assert "space" in app : "must set app space";
assert "image" in app : "must set app image";
assert "labels" in app : "must set app labels";
assert "selector" in app : "must set app selector";
assert "namespace" in app : "must set app namespace";
assert "component" in app : "must set app component";

assert std.length(app.containers) > 0 : "must set one container";

local namespace = app.namespace;
local component = app.component;
local firstcontainer = app.containers[0];
local getservicename(data) = if "name" in data then app.name + "-" + data.name else app.name;
local getconfigname(data) = if "name" in data then app.name + "-" + data.name else app.name + "-config";
local getsecretname(data) = if "name" in data then app.name + "-" + data.name else app.name + "-secret";
local mergeallenvs(data) = if "generalEnvs" in app then if "env" in data then std.uniq(std.sort(data.env + app.generalEnvs, function(one) one.name), function(one) one.name) else app.generalEnvs else data.env;
local appdomainsuffix = if "domain" in app then if std.startsWith(app.domain, "-") then app.domain else "." + app.domain else "";
local getfulldomain(data) = if "fullDomain" in data then data.fullDomain else if "preDomain" in data then data.preDomain + appdomainsuffix else app.name + appdomainsuffix;
local getallcontour = std.flatMap(function(one) (if "routers" in one then std.filterMap(function(router) ("tcpEnable" in router && router.tcpEnable) || ("httpPrefix" in router && std.length(router.httpPrefix) > 0) || ("httpHeader" in router && std.length(router.httpHeader) > 0) || ("wsPrefix" in router && std.length(router.wsPrefix) > 0) || ("wsHeader" in router && std.length(router.wsHeader) > 0), function(router) {name: getservicename(one), router: router}, one.routers) else []), (if "accessPortals" in app then std.filter(function(access) "type" in access && "Contour" == access.type, app.accessPortals) else []));
local getallsecret = if "volumeMounts" in app then std.filter(function(vol) "Secret" == vol.type && "data" in vol && std.length(vol.data) > 0, app.volumeMounts) else [];
local getallconfig = if "volumeMounts" in app then std.filter(function(vol) "Config" == vol.type && "data" in vol && std.length(vol.data) > 0, app.volumeMounts) else [];

{
    main: {
        kind: "Deployment",
        metadata: {
            name: app.name,
            namespace: namespace,
            labels: app.labels,
        },
        spec: {
            replicas: if "replicas" in app then app.replicas else 1,
            [if "minReadySeconds" in app then "minReadySeconds"]: app.minReadySeconds,
            [if "revisionHistoryLimit" in app then "revisionHistoryLimit"]: app.revisionHistoryLimit,
            selector: {
                matchLabels: app.selector
            },
            strategy:
            if "recreate" in app && app.recreate then {
                type: "Recreate"
            } else {
                rollingUpdate: {
                    maxSurge: if "maxSurge" in app then app.maxSurge else "1",
                    maxUnavailable: if "maxUnavailable" in app then app.maxUnavailable else "1",
                },
                type: "RollingUpdate",
            },
            template: {
                metadata: {
                    labels: app.labels
                },
                spec: {
                    [if "hostNetwork" in app then "hostNetwork"]: app.hostNetwork,
                    [if "hostAliases" in app then "hostAliases"]: app.hostAliases,
                    dnsPolicy: if "dnsPolicy" in app then app.dnsPolicy else "ClusterFirstWithHostNet",
                    restartPolicy: if "restartPolicy" in app then app.restartPolicy else "Always",
                    [if "terminationGracePeriodSeconds" in app then "terminationGracePeriodSeconds"]: app.terminationGracePeriodSeconds,
                    [if "securityContext" in app then "securityContext"]: app.securityContext,
                    [if "initContainers" in app && std.length(app.initContainers) > 0 then "initContainers"]: [
                        {
                            name: if "name" in one then one.name else app.name,
                            image: if "image" in one then one.image else app.image,
                            [if "env" in one || "generalEnvs" in app then "env"]: mergeallenvs(one),
                            [if "envFrom" in one then "envFrom"]: one.envFrom,
                            [if "ports" in one then "ports"]: one.ports,
                            [if "command" in one then "command"]: one.command,
                            [if "args" in one then "args"]: one.args,
                            [if "imagePullPolicy" in app then "imagePullPolicy"]: app.imagePullPolicy,
                            [if "securityContext" in one then "securityContext"]: one.securityContext,
                            volumeMounts: [
                                {
                                    name: "timezone",
                                    mountPath: "/etc/localtime",
                                    readOnly: true,
                                }
                            ] +
                            (if "volumeNames" in one && std.length(one.volumeNames) > 0 then [
                                {
                                    name: volume.name,
                                    mountPath: volume.mountPath,
                                    [if "readOnly" in volume && volume.readOnly then "readOnly"]: volume.readOnly,
                                } for volume in std.filter(function(vol) "mountPath" in vol && std.length(vol.mountPath) > 0 && std.count(one.volumeNames, vol.name) > 0, app.volumeMounts)
                            ] else [])
                        } for one in app.initContainers
                    ],
                    containers: [
                        {
                            name: if "name" in one then one.name else app.name,
                            image: if "image" in one then one.image else app.image,
                            [if "env" in one || "generalEnvs" in app then "env"]: mergeallenvs(one),
                            [if "envFrom" in one then "envFrom"]: one.envFrom,
                            [if "stdin" in one then "stdin"]: one.stdin,
                            [if "tty" in one then "tty"]: one.tty,
                            [if "ports" in one then "ports"]: one.ports,
                            [if "format" in one && "specs" in app && std.objectHas(app.specs, one.format) then "resources"]: {
                               requests: {
                                   cpu: if "requestCpu" in app.specs[one.format] then app.specs[one.format].requestCpu else "1",
                                   memory: if "requestMemory" in app.specs[one.format] then app.specs[one.format].requestMemory else "1Gi",
                               },
                               limits: {
                                   cpu: if "limitCpu" in app.specs[one.format] then app.specs[one.format].limitCpu else "2",
                                   memory: if "limitMemory" in app.specs[one.format] then app.specs[one.format].limitMemory else "8Gi",
                                   [if "limitNvidia" in app.specs[one.format] then "nvidia.com/gpu"]: app.specs[one.format].limitNvidia,
                               },
                            },
                            [if "command" in one then "command"]: one.command,
                            [if "args" in one then "args"]: one.args,
                            [if "lifecyclePostStart" in one || "lifecyclePreStop" in one then "lifecycle"]: {
                                [if "lifecyclePostStart" in one then "postStart"]: one.lifecyclePostStart,
                                [if "lifecyclePreStop" in one then "preStop"]: one.lifecyclePreStop,
                            },
                            [if "startupProbe" in one then "startupProbe"]: one.startupProbe,
                            [if "readinessProbe" in one then "readinessProbe"]: one.readinessProbe,
                            [if "livenessProbe" in one then "livenessProbe"]: one.livenessProbe,
                            [if "imagePullPolicy" in app then "imagePullPolicy"]: app.imagePullPolicy,
                            [if "securityContext" in one then "securityContext"]: one.securityContext,
                            volumeMounts: [
                                {
                                    name: "timezone",
                                    mountPath: "/etc/localtime",
                                    readOnly: true,
                                }
                            ] +
                            (if firstcontainer == one && "volumeMounts" in app && std.length(app.volumeMounts) > 0 && (!std.objectHas(one, "volumeNames") || std.length(one.volumeNames) == 0) then [
                                {
                                    name: volume.name,
                                    mountPath: volume.mountPath,
                                    [if "readOnly" in volume && volume.readOnly then "readOnly"]: volume.readOnly,
                                } for volume in std.filter(function(vol) "mountPath" in vol && std.length(vol.mountPath) > 0, app.volumeMounts)
                            ] else []) +
                            (if "volumeNames" in one && std.length(one.volumeNames) > 0 then [
                                {
                                    name: volume.name,
                                    mountPath: volume.mountPath,
                                    [if "readOnly" in volume && volume.readOnly then "readOnly"]: volume.readOnly,
                                } for volume in std.filter(function(vol) "mountPath" in vol && std.length(vol.mountPath) > 0 && std.count(one.volumeNames, vol.name) > 0, app.volumeMounts)
                            ] else [])
                        } for one in app.containers
                    ],
                    [if "imagePullSecrets" in app then "imagePullSecrets"]: app.imagePullSecrets,
                    [if "nodeSelector" in app then "nodeSelector"]: app.nodeSelector,
                    [if "affinity" in app then "affinity"]: {
                        [if "nodeExclusive" == app.affinity then "podAntiAffinity"]: {
                            "requiredDuringSchedulingIgnoredDuringExecution": [
                                {
                                    "topologyKey": "kubernetes.io/hostname",
                                    "labelSelector": {
                                        "matchExpressions": [
                                            {
                                                "key": app.constant.name,
                                                "operator": "In",
                                                "values": [
                                                    app.uk
                                                ]
                                            },
                                            {
                                                "key": app.constant.space,
                                                "operator": "In",
                                                "values": [
                                                    app.space
                                                ]
                                            },
                                            {
                                                "key": app.constant.component,
                                                "operator": "In",
                                                "values": [
                                                    component
                                                ]
                                            },
                                            {
                                                "key": app.constant.managed,
                                                "operator": "In",
                                                "values": [
                                                    "Clap"
                                                ]
                                            },
                                        ]
                                    },
                                },
                            ],
                        },
                    },
                    tolerations: [
                        {
                            effect: "NoSchedule",
                            key: "node-role.kubernetes.io/master",
                            operator: "Exists"
                        }
                    ] + if "tolerations" in app then app.tolerations else [],
                    volumes: [
                        {
                            hostPath: {
                                path: if "timezonePath" in app then app.timezonePath else "/usr/share/zoneinfo/Asia/Shanghai",
                            },
                            name: "timezone",
                        }
                    ] +
                    (if "volumeMounts" in app then std.flatMap(function(one) [one],
                    [
                        {
                            name: vol.name,
                            secret: {
                                secretName: getsecretname(vol),
                                [if "items" in vol then "items"]: vol.items,
                                [if "defaultMode" in vol then "defaultMode"]: vol.defaultMode,
                            },
                        } for vol in std.filter(function(vol) "Secret" == vol.type, app.volumeMounts)
                    ] + [
                        {
                            name: vol.name,
                            configMap: {
                                name: getconfigname(vol),
                                [if "items" in vol then "items"]: vol.items,
                                [if "defaultMode" in vol then "defaultMode"]: vol.defaultMode,
                            },
                        } for vol in std.filter(function(vol) "Config" == vol.type, app.volumeMounts)
                    ] + [
                        {
                            name: vol.name,
                            hostPath: {
                                path: vol.hostPath
                            }
                        } for vol in std.filter(function(vol) "HostPath" == vol.type && "hostPath" in vol, app.volumeMounts)
                    ] + [
                        {
                            name: vol.name,
                            persistentVolumeClaim: {
                                claimName: vol.claimName
                            }
                        } for vol in std.filter(function(vol) "VolumeClaim" == vol.type && "claimName" in vol, app.volumeMounts)
                    ]) else []),
                }
            },
        },
    },
    [if "accessPortals" in app && std.length(app.accessPortals) > 0 then "services"]: {
        [getservicename(svc)]: {
            kind: "Service",
            metadata: {
                name: getservicename(svc),
                namespace: namespace,
                labels: app.labels,
            },
            spec: {
                selector: app.selector,
                [if "type" in svc && (svc.type == "NodePort" || svc.type == "LoadBalancer") then "type"]: svc.type,
                [if "clusterIP" in svc && std.length(svc.clusterIP) > 0 then "clusterIP"]: svc.clusterIP,
                [if "sessionAffinity" in svc then "sessionAffinity"]: svc.sessionAffinity,
                [if "externalTrafficPolicy" in svc && (svc.type == "NodePort" || svc.type == "LoadBalancer") then "externalTrafficPolicy"]: svc.externalTrafficPolicy,
                ports: [
                    {
                        [if "name" in port && std.length(port.name) > 0 then "name"]: port.name,
                        [if "protocol" in port then "protocol"]: port.protocol,
                        port: port.port,
                        [if "targetPort" in port then "targetPort"]: port.targetPort,
                        [if "nodePort" in port && "NodePort" == svc.type then "nodePort"]: port.nodePort,
                    } for port in svc.routers
                ]
            }
        } for svc in app.accessPortals
    },
    [if std.length(getallsecret) > 0 then "secrets"]: {
        [getsecretname(secret)]: {
            kind: "Secret",
            metadata: {
                name: getsecretname(secret),
                namespace: namespace,
                labels: app.labels,
            },
            data: secret.data
        } for secret in std.filter(function(vol) "Secret" == vol.type && "data" in vol && std.length(vol.data) > 0, app.volumeMounts)
    },
    [if std.length(getallconfig) > 0 then "configs"]: {
        [getconfigname(config)]: {
            kind: "ConfigMap",
            metadata: {
                name: getconfigname(config),
                namespace: namespace,
                labels: app.labels,
            },
            data: config.data
        } for config in std.filter(function(vol) "Config" == vol.type && "data" in vol && std.length(vol.data) > 0, app.volumeMounts)
    },
    [if std.length(getallcontour) > 0 then "contours"]: {
        [getfulldomain(svc.router)]: {
            apiVersion: app.contour.apiVersion,
            kind: "HTTPProxy",
            metadata: {
                name: getfulldomain(svc.router),
                namespace: namespace,
                labels: app.labels,
            },
            spec: {
                virtualhost: {
                    fqdn: getfulldomain(svc.router),
                    [if ("tls" in app && "secretName" in app.tls) || ("tls" in svc.router && "secretName" in svc.router.tls && std.length(svc.router.tls.secretName) > 0) || ("tcpEnable" in svc.router && svc.router.tcpEnable) then "tls"]: {
                        passthrough: if ("tls" in app && "passthrough" in app.tls) then app.tls.passthrough else if ("tls" in app && "secretName" in app.tls) then false else !("tls" in svc.router && "secretName" in svc.router.tls && std.length(svc.router.tls.secretName) > 0),
                        [if ("tls" in app && "secretName" in app.tls) || ("tls" in svc.router && "secretName" in svc.router.tls && std.length(svc.router.tls.secretName) > 0) then "secretName"]: if ("tls" in app && "secretName" in app.tls) then app.tls.secretName else svc.router.tls.secretName,
                    },
                    [if "corsPolicy" in svc.router && std.length(svc.router.corsPolicy) > 0 then "corsPolicy"]: {
                        [if "allowCredentials" in svc.router.corsPolicy.allowCredentials && svc.router.corsPolicy.allowCredentials then "allowCredentials"]: svc.router.corsPolicy.allowCredentials,
                        [if "allowOrigin" in svc.router.corsPolicy.allowOrigin && std.length(svc.router.corsPolicy.allowOrigin) > 0 then "allowOrigin"]: svc.router.corsPolicy.allowOrigin,
                        [if "allowMethods" in svc.router.corsPolicy.allowMethods && std.length(svc.router.corsPolicy.allowMethods) > 0 then "allowMethods"]: svc.router.corsPolicy.allowMethods,
                        [if "allowHeaders" in svc.router.corsPolicy.allowHeaders && std.length(svc.router.corsPolicy.allowHeaders) > 0 then "allowHeaders"]: svc.router.corsPolicy.allowHeaders,
                        [if "exposeHeaders" in svc.router.corsPolicy.exposeHeaders && std.length(svc.router.corsPolicy.exposeHeaders) > 0 then "exposeHeaders"]: svc.router.corsPolicy.exposeHeaders,
                        [if "maxAge" in svc.router.corsPolicy.maxAge && std.length(svc.router.corsPolicy.maxAge) > 0 then "maxAge"]: svc.router.corsPolicy.maxAge,
                    },
                },
                [if "httpPrefix" in svc.router || "httpHeader" in svc.router || "wsPrefix" in svc.router || "wsHeader" in svc.router then "routes"]:
                (if "httpPrefix" in svc.router || "httpHeader" in svc.router then [
                    {
                        services: [
                            {
                                "name": svc.name,
                                "port": svc.router.port,
                                [if "protocol" in svc.router then "protocol"]: svc.router.protocol,
                            }
                        ],
                        conditions: std.flatMap(function(one) [one],
                        [
                            {
                                "prefix": prefix
                            } for prefix in (if "httpPrefix" in svc.router && std.length(svc.router.httpPrefix) > 0 then svc.router.httpPrefix else [])
                        ] + [
                            {
                                "header": header
                            } for header in (if "httpHeader" in svc.router && std.length(svc.router.httpHeader) > 0 then svc.router.httpHeader else [])
                        ]),
                        [if "respTimeout" in svc.router || "idleTimeout" in svc.router then "timeoutPolicy"]: {
                            [if "respTimeout" in svc.router then "response"]: svc.router.respTimeout,
                            [if "idleTimeout" in svc.router then "idle"]: svc.router.idleTimeout,
                        },
                        [if "retryCount" in svc.router then "retryPolicy"]: {
                            "retryCount": svc.router.retryCount,
                            [if "retryPerTryTimeout" in svc.router then "perTryTimeout"]: svc.router.retryPerTryTimeout,
                        },
                        [if "loadBalanceStrategy" in svc.router then "loadBalancerPolicy"]: {
                            "strategy": svc.router.loadBalanceStrategy,
                        },
                        [if "healthCheckPolicy" in svc.router && std.length(svc.router.healthCheckPolicy) > 0 then "healthCheckPolicy"]: svc.router.healthCheckPolicy,
                    }
                ] else []) + (if "wsPrefix" in svc.router || "wsHeader" in svc.router then [
                    {
                        services: [
                            {
                                "name": svc.name,
                                "port": svc.router.port,
                            }
                        ],
                        enableWebsockets: true,
                        conditions: std.flatMap(function(one) [one],
                        [
                            {
                                "prefix": prefix
                            } for prefix in (if "wsPrefix" in svc.router && std.length(svc.router.wsPrefix) > 0 then svc.router.wsPrefix else [])
                        ] + [
                            {
                                "header": header
                            } for header in (if "wsHeader" in svc.router && std.length(svc.router.wsHeader) > 0 then svc.router.wsHeader else [])
                        ]),
                        [if "loadBalanceStrategy" in svc.router then "loadBalancerPolicy"]: {
                            "strategy": svc.router.loadBalanceStrategy,
                        },
                        [if "healthCheckPolicy" in svc.router && std.length(svc.router.healthCheckPolicy) > 0 then "healthCheckPolicy"]: svc.router.healthCheckPolicy,
                    }
                ] else []),
                [if "tcpEnable" in svc.router && svc.router.tcpEnable then "tcpproxy"]: {
                    services: [
                        {
                            "name": svc.name,
                            "port": svc.router.port,
                        }
                    ],
                    [if "healthCheckPolicy" in svc.router && std.length(svc.router.healthCheckPolicy) > 0 then "healthCheckPolicy"]: svc.router.healthCheckPolicy,
                },
            }
        } for svc in getallcontour
    },
    [if "budget" in app then "budget"]: {

    },
    [if "policy" in app then "policy"]: {

    },
    [if "scaler" in app then "scaler"]: {

    },
}
