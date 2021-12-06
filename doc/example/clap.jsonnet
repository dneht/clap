{
    env: 'dev',
    port: 8008,
    debug: false,
    namespace: 'clap-system',
    timezone: '',
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