server:
  app_name: <xpfo{ .ProjectName }xpfo>
  host: 0.0.0.0
  port: 17878
  is_debug: true
  pprof_username: admin
  pprof_password: admin

log:
  file_name: log
  level: debug
  max_age: 21

mysql:
    host: 127.0.0.1
    port: 3306
    user: root
    password: root
    name: db_name
    max_open_conn: 100
    max_idle_conn: 25
    conn_max_lifetime_second: 600
<xpfo{ if .EnableRedis }xpfo>

redis:
    host: 127.0.0.1
    port: 6379
    password: ""
    db: 0
    pool_size: 20
<xpfo{ end }xpfo>

auth:
    enabled: <xpfo{ .EnableAuth }xpfo>
    jwt_secret: "please-change-me"

rate_limit:
    enabled: <xpfo{ .EnableRate }xpfo>
    rps: 50
    burst: 100
