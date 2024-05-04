server:
  app_name: <xpfo{ .ProjectName }xpfo>
  host: 0.0.0.0
  port: 17878
  is_debug: true

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

