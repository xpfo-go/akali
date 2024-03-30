server:
  app_name: {{ .ProjectName }}
  host: 0.0.0.0
  port: 17878
  is_debug: true

log:
  level: debug

database:
    host: 127.0.0.1
    port: 3306
    user: root
    password: root
    name: db_name
    max_open_conn: 100
    max_idle_conn: 25
    conn_max_lifetime_second: 600

