# `heroku-logs-exporter`

Prometheus exporter for Heroku applications. Written in Go. It parses incoming logs from [Heroku Log Drain](https://devcenter.heroku.com/articles/log-drains).

**This exporter is just a proof of concept and still being developed.**

## Installation and usage

`heroku-logs-exporter` listens on HTTP port 9841 by default.

By default it it exports metrics on `/metrics` and listens for incoming logs from Heroku Log Drain on `/logs`.

### Deploying `heroku-logs-exporter`

You have to deploy `heroku-logs-exporter` somewhere where Heroku Log Drain can send its logs to.

Run `heroku-logs-exporter -h` to see all options.

```sh
$ heroku-logs-exporter -web.logs-token-param-value "secret-token"
```

### Setting up Heroku Log Drain

When adding Heroku Log Drain you have to set application name using `app_name` query parameter. You can also set `token` parameter to authorize with `heroku-logs-exporter`.

```sh
$ heroku drains:add "http://example.com:9841/logs?app_name=your-app&token=secret-token" -a your-app
```

## Metrics

### Heroku Router

Metrics sent by Heroku Router to Heroku log. `connect` and `service` duration are collected as histogram and summary metrics.

Histogram buckets for `connect` metric are `.001, .002, .003, .004, .005, .01, 0.025, .05, .1, .25, .5, 1.0, 2.5, 5.0, 10.0, 20.0` in seconds.

Histogram buckets for `service` metric are `.005, .01, .02, 0.04, .06, .08, 0.1, .125, 0.15, 0.175, 0.2, 0.3, 0.4, .5, 1, 2.5, 5, 10, 15, 20` in seconds.

Summary quantiles with their absolute errors for both `connect` and `service` metrics are `0.01: 0.001, 0.1: 0.01, 0.5: 0.05, 0.9: 0.01, 0.95: 0.001, 0.99: 0.001`.

<details>
  <summary>Sample metrics</summary>

```
# HELP heroku_router_connect_duration_histogram_seconds Request connect duration reported by Heroku Router as histogram.
# TYPE heroku_router_connect_duration_histogram_seconds histogram
heroku_router_connect_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.001"} 56
heroku_router_connect_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.002"} 57
heroku_router_connect_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.003"} 57
heroku_router_connect_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.004"} 57
heroku_router_connect_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.005"} 57
heroku_router_connect_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.01"} 57
heroku_router_connect_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.025"} 57
heroku_router_connect_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.05"} 57
heroku_router_connect_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.1"} 57
heroku_router_connect_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.25"} 57
heroku_router_connect_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.5"} 57
heroku_router_connect_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="1"} 57
heroku_router_connect_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="2.5"} 57
heroku_router_connect_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="5"} 57
heroku_router_connect_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="10"} 57
heroku_router_connect_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="20"} 57
heroku_router_connect_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="+Inf"} 57
heroku_router_connect_duration_histogram_seconds_sum{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200"} 0.014000000000000005
# HELP heroku_router_connect_duration_seconds Request connect duration reported by Heroku Router as summary.
# TYPE heroku_router_connect_duration_seconds summary
heroku_router_connect_duration_seconds{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",quantile="0.01"} 0
heroku_router_connect_duration_seconds{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",quantile="0.1"} 0
heroku_router_connect_duration_seconds{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",quantile="0.5"} 0
heroku_router_connect_duration_seconds{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",quantile="0.9"} 0.001
heroku_router_connect_duration_seconds{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",quantile="0.95"} 0.001
heroku_router_connect_duration_seconds{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",quantile="0.99"} 0.002
heroku_router_connect_duration_seconds_sum{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200"} 0.014000000000000005
# HELP heroku_router_service_duration_histogram_seconds Request service duration reported by Heroku Router as histogram.
# TYPE heroku_router_service_duration_histogram_seconds histogram
heroku_router_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.005"} 0
heroku_router_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.01"} 0
heroku_router_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.02"} 0
heroku_router_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.04"} 12
heroku_router_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.06"} 38
heroku_router_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.08"} 49
heroku_router_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.1"} 54
heroku_router_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.125"} 55
heroku_router_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.15"} 56
heroku_router_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.175"} 56
heroku_router_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.2"} 57
heroku_router_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.3"} 57
heroku_router_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.4"} 57
heroku_router_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="0.5"} 57
heroku_router_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="1"} 57
heroku_router_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="2.5"} 57
heroku_router_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="5"} 57
heroku_router_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="10"} 57
heroku_router_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="15"} 57
heroku_router_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="20"} 57
heroku_router_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200",le="+Inf"} 57
heroku_router_service_duration_histogram_seconds_sum{app_name="slideslive",dyno="web.1",host="slideslive.com",method="GET",protocol="https",status="200"} 3.287
# HELP heroku_router_service_duration_seconds Request service duration reported by Heroku Router as summary.
# TYPE heroku_router_service_duration_seconds summary
heroku_router_service_duration_seconds{app_name="slideslive",dyno="web.1",host="slideslive.at",method="GET",protocol="https",status="200",quantile="0.01"} 0.034
heroku_router_service_duration_seconds{app_name="slideslive",dyno="web.1",host="slideslive.at",method="GET",protocol="https",status="200",quantile="0.1"} 0.034
heroku_router_service_duration_seconds{app_name="slideslive",dyno="web.1",host="slideslive.at",method="GET",protocol="https",status="200",quantile="0.5"} 0.034
heroku_router_service_duration_seconds{app_name="slideslive",dyno="web.1",host="slideslive.at",method="GET",protocol="https",status="200",quantile="0.9"} 0.034
heroku_router_service_duration_seconds{app_name="slideslive",dyno="web.1",host="slideslive.at",method="GET",protocol="https",status="200",quantile="0.95"} 0.034
heroku_router_service_duration_seconds{app_name="slideslive",dyno="web.1",host="slideslive.at",method="GET",protocol="https",status="200",quantile="0.99"} 0.034
heroku_router_service_duration_seconds_sum{app_name="slideslive",dyno="web.1",host="slideslive.at",method="GET",protocol="https",status="200"} 0.034
```

</details>

### Heroku Runtime

Enable [Heroku Labs: log-runtime-metrics](https://devcenter.heroku.com/articles/log-runtime-metrics) for these metrics to be collected. They are all collected as gauges.

```sh
$ heroku labs:enable log-runtime-metrics -a your-app
$ heroku restart -a your-app
```

<details>
  <summary>Sample metrics</summary>

```
# HELP heroku_runtime_metrics_load_avg_15m The load average for the dyno in the last 15 minutes. This reflects the number of CPU tasks that are in the ready queue (i.e. waiting to be processed).
# TYPE heroku_runtime_metrics_load_avg_15m gauge
heroku_runtime_metrics_load_avg_15m{app_name="slideslive",dyno="web.1",dyno_id="web.1"} 0.41
# HELP heroku_runtime_metrics_load_avg_1m The load average for the dyno in the last 1 minute. This reflects the number of CPU tasks that are in the ready queue (i.e. waiting to be processed).
# TYPE heroku_runtime_metrics_load_avg_1m gauge
heroku_runtime_metrics_load_avg_1m{app_name="slideslive",dyno="web.1",dyno_id="web.1"} 1.07
# HELP heroku_runtime_metrics_load_avg_5m The load average for the dyno in the last 5 minutes. This reflects the number of CPU tasks that are in the ready queue (i.e. waiting to be processed).
# TYPE heroku_runtime_metrics_load_avg_5m gauge
heroku_runtime_metrics_load_avg_5m{app_name="slideslive",dyno="web.1",dyno_id="web.1"} 0.65
# HELP heroku_runtime_metrics_memory_cache_bytes The portion of the dyno’s memory used as disk cache.
# TYPE heroku_runtime_metrics_memory_cache_bytes gauge
heroku_runtime_metrics_memory_cache_bytes{app_name="slideslive",dyno="web.1",dyno_id="web.1"} 2.609905664e+07
# HELP heroku_runtime_metrics_memory_pgpgin_pages The cumulative total of the pages written to disk. Sudden high variations on this number can indicate short duration spikes in swap usage. The other memory related metrics are point in time snapshots and can miss short spikes.
# TYPE heroku_runtime_metrics_memory_pgpgin_pages gauge
heroku_runtime_metrics_memory_pgpgin_pages{app_name="slideslive",dyno="web.1",dyno_id="web.1"} 1.8879057e+07
# HELP heroku_runtime_metrics_memory_pgpgout_pages The cumulative total of the pages read from disk. Sudden high variations on this number can indicate short duration spikes in swap usage. The other memory related metrics are point in time snapshots and can miss short spikes.
# TYPE heroku_runtime_metrics_memory_pgpgout_pages gauge
heroku_runtime_metrics_memory_pgpgout_pages{app_name="slideslive",dyno="web.1",dyno_id="web.1"} 1.8631374e+07
# HELP heroku_runtime_metrics_memory_quota_bytes The resident memory (memory_rss) value at which an R14 is triggered.
# TYPE heroku_runtime_metrics_memory_quota_bytes gauge
heroku_runtime_metrics_memory_quota_bytes{app_name="slideslive",dyno="web.1",dyno_id="web.1"} 1.073741824e+09
# HELP heroku_runtime_metrics_memory_rss_bytes The portion of the dyno’s memory held in RAM.
# TYPE heroku_runtime_metrics_memory_rss_bytes gauge
heroku_runtime_metrics_memory_rss_bytes{app_name="slideslive",dyno="web.1",dyno_id="web.1"} 9.9050586112e+08
# HELP heroku_runtime_metrics_memory_swap_bytes The portion of a dyno’s memory stored on disk.
# TYPE heroku_runtime_metrics_memory_swap_bytes gauge
heroku_runtime_metrics_memory_swap_bytes{app_name="slideslive",dyno="web.1",dyno_id="web.1"} 251658.24
# HELP heroku_runtime_metrics_memory_total_bytes The total memory being used by the dyno, equal to the sum of resident, cache, and swap memory.
# TYPE heroku_runtime_metrics_memory_total_bytes gauge
heroku_runtime_metrics_memory_total_bytes{app_name="slideslive",dyno="web.1",dyno_id="web.1"} 1.016856576e+09
```

</details>

### Heroku Postgres

These metrics are collected when you have Heroku Postgres addon. They are described in [Heroku Postgres Metrics Logs](https://devcenter.heroku.com/articles/heroku-postgres-metrics-logs).

<details>
  <summary>Sample metrics</summary>

```
# HELP heroku_postgres_metrics_active_connection_count The number of connections established on the database.
# TYPE heroku_postgres_metrics_active_connection_count gauge
heroku_postgres_metrics_active_connection_count{addon="postgresql-objective-42956",app_name="slideslive",source="HEROKU_POSTGRESQL_CHARCOAL"} 76
# HELP heroku_postgres_metrics_db_size_bytes The number of bytes contained in the database. This includes all table and index data on disk, including database bloat.
# TYPE heroku_postgres_metrics_db_size_bytes gauge
heroku_postgres_metrics_db_size_bytes{addon="postgresql-objective-42956",app_name="slideslive",source="HEROKU_POSTGRESQL_CHARCOAL"} 2.7259822959e+10
# HELP heroku_postgres_metrics_index_cache_hit_rate Ratio of index lookups served from shared buffer cache, rounded to five decimal points.
# TYPE heroku_postgres_metrics_index_cache_hit_rate gauge
heroku_postgres_metrics_index_cache_hit_rate{addon="postgresql-objective-42956",app_name="slideslive",source="HEROKU_POSTGRESQL_CHARCOAL"} 0.99737
# HELP heroku_postgres_metrics_load_avg_15m The average system load over a period of 15 minutes divided by the number of available CPUs.
# TYPE heroku_postgres_metrics_load_avg_15m gauge
heroku_postgres_metrics_load_avg_15m{addon="postgresql-objective-42956",app_name="slideslive",source="HEROKU_POSTGRESQL_CHARCOAL"} 0.05
# HELP heroku_postgres_metrics_load_avg_1m The average system load over a period of 1 minute divided by the number of available CPUs.
# TYPE heroku_postgres_metrics_load_avg_1m gauge
heroku_postgres_metrics_load_avg_1m{addon="postgresql-objective-42956",app_name="slideslive",source="HEROKU_POSTGRESQL_CHARCOAL"} 0.025
# HELP heroku_postgres_metrics_load_avg_5m The average system load over a period of 5 minutes divided by the number of available CPUs.
# TYPE heroku_postgres_metrics_load_avg_5m gauge
heroku_postgres_metrics_load_avg_5m{addon="postgresql-objective-42956",app_name="slideslive",source="HEROKU_POSTGRESQL_CHARCOAL"} 0.045
# HELP heroku_postgres_metrics_memory_cached_bytes Amount of memory being used by the OS for page cache.
# TYPE heroku_postgres_metrics_memory_cached_bytes gauge
heroku_postgres_metrics_memory_cached_bytes{addon="postgresql-objective-42956",app_name="slideslive",source="HEROKU_POSTGRESQL_CHARCOAL"} 3.060965376e+09
# HELP heroku_postgres_metrics_memory_free_bytes Amount of free memory available.
# TYPE heroku_postgres_metrics_memory_free_bytes gauge
heroku_postgres_metrics_memory_free_bytes{addon="postgresql-objective-42956",app_name="slideslive",source="HEROKU_POSTGRESQL_CHARCOAL"} 1.60415744e+08
# HELP heroku_postgres_metrics_memory_postgres_bytes Approximate amount of memory used by your database’s Postgres processes. This includes shared buffer cache as well as memory for each connection.
# TYPE heroku_postgres_metrics_memory_postgres_bytes gauge
heroku_postgres_metrics_memory_postgres_bytes{addon="postgresql-objective-42956",app_name="slideslive",source="HEROKU_POSTGRESQL_CHARCOAL"} 5.48995072e+08
# HELP heroku_postgres_metrics_memory_total_bytes Total amount of server memory available.
# TYPE heroku_postgres_metrics_memory_total_bytes gauge
heroku_postgres_metrics_memory_total_bytes{addon="postgresql-objective-42956",app_name="slideslive",source="HEROKU_POSTGRESQL_CHARCOAL"} 4.142166016e+09
# HELP heroku_postgres_metrics_read_iops Number of read operations in I/O sizes of 16KB blocks.
# TYPE heroku_postgres_metrics_read_iops gauge
heroku_postgres_metrics_read_iops{addon="postgresql-objective-42956",app_name="slideslive",source="HEROKU_POSTGRESQL_CHARCOAL"} 2.3317
# HELP heroku_postgres_metrics_table_cache_hit_rate Ratio of table lookups served from shared buffer cache, rounded to five decimal points.
# TYPE heroku_postgres_metrics_table_cache_hit_rate gauge
heroku_postgres_metrics_table_cache_hit_rate{addon="postgresql-objective-42956",app_name="slideslive",source="HEROKU_POSTGRESQL_CHARCOAL"} 0.99782
# HELP heroku_postgres_metrics_table_count The number of tables in the database.
# TYPE heroku_postgres_metrics_table_count gauge
heroku_postgres_metrics_table_count{addon="postgresql-objective-42956",app_name="slideslive",source="HEROKU_POSTGRESQL_CHARCOAL"} 136
# HELP heroku_postgres_metrics_tmp_disk_available_bytes Amount of bytes available on tmp mount.
# TYPE heroku_postgres_metrics_tmp_disk_available_bytes gauge
heroku_postgres_metrics_tmp_disk_available_bytes{addon="postgresql-objective-42956",app_name="slideslive",source="HEROKU_POSTGRESQL_CHARCOAL"} 7.2944943104e+10
# HELP heroku_postgres_metrics_tmp_disk_used_bytes Amount of bytes used on tmp mount.
# TYPE heroku_postgres_metrics_tmp_disk_used_bytes gauge
heroku_postgres_metrics_tmp_disk_used_bytes{addon="postgresql-objective-42956",app_name="slideslive",source="HEROKU_POSTGRESQL_CHARCOAL"} 3.3849344e+07
# HELP heroku_postgres_metrics_transactions_total The current transaction ID, which can be used to track writes over time.
# TYPE heroku_postgres_metrics_transactions_total gauge
heroku_postgres_metrics_transactions_total{addon="postgresql-objective-42956",app_name="slideslive",source="HEROKU_POSTGRESQL_CHARCOAL"} 4.2836737e+08
# HELP heroku_postgres_metrics_waiting_connection_count Number of connections waiting on a lock to be acquired. If many connections are waiting, this can be a sign of mishandled database concurrency.
# TYPE heroku_postgres_metrics_waiting_connection_count gauge
heroku_postgres_metrics_waiting_connection_count{addon="postgresql-objective-42956",app_name="slideslive",source="HEROKU_POSTGRESQL_CHARCOAL"} 0
# HELP heroku_postgres_metrics_wal_percentage_used Percentage of the WAL disk that has been used, between 0.0 and 1.0.
# TYPE heroku_postgres_metrics_wal_percentage_used gauge
heroku_postgres_metrics_wal_percentage_used{addon="postgresql-objective-42956",app_name="slideslive",source="HEROKU_POSTGRESQL_CHARCOAL"} 0.06710167833628045
# HELP heroku_postgres_metrics_write_iops Number of write operations in I/O sizes of 16KB blocks.
# TYPE heroku_postgres_metrics_write_iops gauge
heroku_postgres_metrics_write_iops{addon="postgresql-objective-42956",app_name="slideslive",source="HEROKU_POSTGRESQL_CHARCOAL"} 2.375
```

</details>

### Heroku Pgbouncer

These metrics are collected when you have Heroku PgBouncer pooler attachment. They are described in [Heroku Postgres Metrics Logs](https://devcenter.heroku.com/articles/heroku-postgres-metrics-logs#pgbouncer-metrics).

<details>
  <summary>Sample metrics</summary>

```
# HELP heroku_pgbouncer_metrics_server_idle_count The number of server connections that are not currently assigned to a client connection.
# TYPE heroku_pgbouncer_metrics_server_idle_count gauge
heroku_pgbouncer_metrics_server_idle_count{addon="postgresql-objective-42956",app_name="slideslive",source="pgbouncer"} 4
# HELP heroku_pgbouncer_metrics_avg_query_seconds The average query time of all queries executed through through poolec connections.
# TYPE heroku_pgbouncer_metrics_avg_query_seconds gauge
heroku_pgbouncer_metrics_avg_query_seconds{addon="postgresql-objective-42956",app_name="slideslive",source="pgbouncer"} 0
# HELP heroku_pgbouncer_metrics_avg_recv_bytes The average amount of bytes received from clients per second.
# TYPE heroku_pgbouncer_metrics_avg_recv_bytes gauge
heroku_pgbouncer_metrics_avg_recv_bytes{addon="postgresql-objective-42956",app_name="slideslive",source="pgbouncer"} 20400
# HELP heroku_pgbouncer_metrics_avg_sent_bytes The average amount of bytes sent to clients per second.
# TYPE heroku_pgbouncer_metrics_avg_sent_bytes gauge
heroku_pgbouncer_metrics_avg_sent_bytes{addon="postgresql-objective-42956",app_name="slideslive",source="pgbouncer"} 606498
# HELP heroku_pgbouncer_metrics_client_active_count The number of client connections to the pooler that have an active server connection assignment.
# TYPE heroku_pgbouncer_metrics_client_active_count gauge
heroku_pgbouncer_metrics_client_active_count{addon="postgresql-objective-42956",app_name="slideslive",source="pgbouncer"} 287
# HELP heroku_pgbouncer_metrics_client_waiting_count The number of client connections to the pooler that are waiting for a server connection assignment.
# TYPE heroku_pgbouncer_metrics_client_waiting_count gauge
heroku_pgbouncer_metrics_client_waiting_count{addon="postgresql-objective-42956",app_name="slideslive",source="pgbouncer"} 0
# HELP heroku_pgbouncer_metrics_max_wait_seconds The longest wait time of any client currently waiting for a server connection assignment.
# TYPE heroku_pgbouncer_metrics_max_wait_seconds gauge
heroku_pgbouncer_metrics_max_wait_seconds{addon="postgresql-objective-42956",app_name="slideslive",source="pgbouncer"} 0
# HELP heroku_pgbouncer_metrics_server_active_count The number of server connections that are currently assigned to a client connection.
# TYPE heroku_pgbouncer_metrics_server_active_count gauge
heroku_pgbouncer_metrics_server_active_count{addon="postgresql-objective-42956",app_name="slideslive",source="pgbouncer"} 0
```

</details>

### rack-timeout

These metrics are collected for Ruby application with `rack-timeout` gem installed. `wait` and `service` durations are collected as summary and histogram metrics.

Histogram buckets for both `wait` and `service` metrics are `.005, .01, .02, 0.04, .06, .08, 0.1, .125, 0.15, 0.175, 0.2, 0.3, 0.4, .5, 1, 2.5, 5, 10, 15, 20` in seconds.

Summary quantiles with their absolute errors for both `wait` and `service` metrics are `0.01: 0.001, 0.1: 0.01, 0.5: 0.05, 0.9: 0.01, 0.95: 0.001, 0.99: 0.001`.

<details>
  <summary>Sample metrics</summary>

```
# HELP heroku_rack_timeout_service_duration_histogram_seconds Request service duration reported by rack-timeout as histogram.
# TYPE heroku_rack_timeout_service_duration_histogram_seconds histogram
heroku_rack_timeout_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.005"} 0
heroku_rack_timeout_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.01"} 11
heroku_rack_timeout_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.02"} 38
heroku_rack_timeout_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.04"} 72
heroku_rack_timeout_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.06"} 85
heroku_rack_timeout_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.08"} 91
heroku_rack_timeout_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.1"} 91
heroku_rack_timeout_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.125"} 93
heroku_rack_timeout_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.15"} 94
heroku_rack_timeout_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.175"} 94
heroku_rack_timeout_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.2"} 94
heroku_rack_timeout_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.3"} 94
heroku_rack_timeout_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.4"} 94
heroku_rack_timeout_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.5"} 94
heroku_rack_timeout_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="1"} 94
heroku_rack_timeout_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="2.5"} 94
heroku_rack_timeout_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="5"} 94
heroku_rack_timeout_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="10"} 94
heroku_rack_timeout_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="15"} 94
heroku_rack_timeout_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="20"} 94
heroku_rack_timeout_service_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="+Inf"} 94
heroku_rack_timeout_service_duration_histogram_seconds_sum{app_name="slideslive",dyno="web.1"} 2.9359999999999995
# HELP heroku_rack_timeout_service_duration_seconds Request service duration reported rack-timeout Ruby gem.
# TYPE heroku_rack_timeout_service_duration_seconds summary
heroku_rack_timeout_service_duration_seconds{app_name="slideslive",dyno="web.1",quantile="0.01"} 0.009
heroku_rack_timeout_service_duration_seconds{app_name="slideslive",dyno="web.1",quantile="0.1"} 0.013
heroku_rack_timeout_service_duration_seconds{app_name="slideslive",dyno="web.1",quantile="0.5"} 0.039
heroku_rack_timeout_service_duration_seconds{app_name="slideslive",dyno="web.1",quantile="0.9"} 0.137
heroku_rack_timeout_service_duration_seconds{app_name="slideslive",dyno="web.1",quantile="0.99"} 0.385
heroku_rack_timeout_service_duration_seconds_sum{app_name="slideslive",dyno="web.1"} 17783.774000000667
# HELP heroku_rack_timeout_wait_duration_histogram_seconds Request wait duration reported by rack-timeout as histogram.
# TYPE heroku_rack_timeout_wait_duration_histogram_seconds histogram
heroku_rack_timeout_wait_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.005"} 16
heroku_rack_timeout_wait_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.01"} 67
heroku_rack_timeout_wait_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.02"} 89
heroku_rack_timeout_wait_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.04"} 93
heroku_rack_timeout_wait_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.06"} 93
heroku_rack_timeout_wait_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.08"} 94
heroku_rack_timeout_wait_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.1"} 94
heroku_rack_timeout_wait_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.125"} 94
heroku_rack_timeout_wait_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.15"} 94
heroku_rack_timeout_wait_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.175"} 94
heroku_rack_timeout_wait_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.2"} 94
heroku_rack_timeout_wait_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.3"} 94
heroku_rack_timeout_wait_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.4"} 94
heroku_rack_timeout_wait_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="0.5"} 94
heroku_rack_timeout_wait_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="1"} 94
heroku_rack_timeout_wait_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="2.5"} 94
heroku_rack_timeout_wait_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="5"} 94
heroku_rack_timeout_wait_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="10"} 94
heroku_rack_timeout_wait_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="15"} 94
heroku_rack_timeout_wait_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="20"} 94
heroku_rack_timeout_wait_duration_histogram_seconds_bucket{app_name="slideslive",dyno="web.1",le="+Inf"} 94
heroku_rack_timeout_wait_duration_histogram_seconds_sum{app_name="slideslive",dyno="web.1"} 0.9190000000000006
# HELP heroku_rack_timeout_wait_duration_seconds Request wait duration reported rack-timeout Ruby gem.
# TYPE heroku_rack_timeout_wait_duration_seconds summary
heroku_rack_timeout_wait_duration_seconds{app_name="slideslive",dyno="web.1",quantile="0.01"} 0.003
heroku_rack_timeout_wait_duration_seconds{app_name="slideslive",dyno="web.1",quantile="0.1"} 0.006
heroku_rack_timeout_wait_duration_seconds{app_name="slideslive",dyno="web.1",quantile="0.5"} 0.012
heroku_rack_timeout_wait_duration_seconds{app_name="slideslive",dyno="web.1",quantile="0.9"} 0.035
heroku_rack_timeout_wait_duration_seconds{app_name="slideslive",dyno="web.1",quantile="0.99"} 0.161
heroku_rack_timeout_wait_duration_seconds_sum{app_name="slideslive",dyno="web.1"} 6062.761000000734
```

</details>
