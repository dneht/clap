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

package static

const DefaultConf = `
{
    env: 'dev',
    port: 8008,
    debug: false,
    namespace: 'clap-system',
    timezone: 'Local',
	message: {
		dingding: {
			enable: false,
			url: ''
		},
	},
	document: {
		java: {
			dev: 'http://do.dev.zbyan.net',
			dev: 'http://do.dev.zbyan.net',
			dev: 'http://do.dev.zbyan.net',
			dev: 'http://do.dev.zbyan.net',
		},
	},
    container: {
        cros: {
            allow_origins: '*',
            allow_methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
        },
        csrf: {},
        secure: {
            content_type_nosniff: 'nosniff',
            x_frame_options: 'SAMEORIGIN',
            hsts_preload_enabled: false,
            referrer_policy: 'origin',
        },
        gzip_level: 6,
        static_path: '/opt/ui',
    },
    database: {
        hosts: ['127.0.0.1'],
        ports: [3306],
        user: 'user',
        password: '',
        database: 'clap',
        max_idle_con: 50,
        max_open_con: 150,
    },
    password: {
        type: 'ssha',
    },
    application: {
        deployment: {
            build_job_image: '',
            image_pull_policy: 'Always',
            image_pull_secret: 'registry-secret',
            maven_skip_tests: false,
            backoff_limit: 2,
            clean_after_finished: 24 * 60 *  60,
        },
    }
}
`

