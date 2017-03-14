Redis-exporter - Redis stats to Sensision Metrics
=============

[![Build Status](https://travis-ci.org/miton18/redis-exporter.svg?branch=master)](https://travis-ci.org/miton18/redis-exporter)

Redis-Exporter scrapes Redis stats and expose them as a Sensision HTTP endpoint.

Redis-Exporter features:
 - **Simple**: Redis-Exporter fetch stats through TCP connexion.
 - **Highly scalable**: Redis-Exporter can export stats of thousands Redis.
 - **Pluggable**: Export your metrics via [Beamium](https://github.com/runabove/beamium).
 - **Versatile**: Redis-Exporter can flush metrics to files.

## Datapoints

Redis-Exporter will automatically collect all metrics exposed by Redis on its
`info` command. To help
classify these metrics, Redis-Exporter will name is ``redis_stats.[METRIC_NAME]``.
Additional labels may be configured in the documentation. See below.

*Sample metrics*
```
1489476345223746// redis_stats_uptime_in_seconds{host=localhost} 96
1489476345223746// redis_stats_used_memory{host=localhost} 822248
1489476345223746// redis_stats_instantaneous_ops_per_sec{host=localhost} 254
1489476345223746// redis_stats_total_net_output_bytes{host=localhost} 31577
1489476345223746// redis_stats_cluster_enabled{host=localhost} 1
```

For more informations, please see Redis Managment Guide: 
http://redis.io/commands/info

## Building

Redis-Exporter is pretty easy to build.
 - Clone the repository
 - Setup a minimal working config (see bellow)
 - Build and run `go run redis-exporter.go`

## Usage
```
redis-exporter [flags]

Flags:
      --config string   config file to use
      --listen string   listen address (default "127.0.0.1:9100")
  -v, --verbose         verbose output
```

## Configuration
Redis-Exporter come with a simple default [config file](config.sample.yaml).

Configuration is load and override in the following order:
 - /etc/redis-exporter/config.yaml
 - ~/redis-exporter/config.yaml
 - ./config.yaml
 - config filepath from command line

### Definitions
Config is composed of three main parts and some config fields:

#### Sources
Redis-Exporter can have one to many Redis stats sources. A *source* is defined as follow:
``` yaml
sources: # Sources definitions
  - uri: 127.0.0.1:6379 # address and port of Redis
    labels: # Labels are added to every metrics (Optional)
      label_name : label_value # Label definition
```

#### Metrics
Redis-Exporter can expose some or all Redis stats:
``` yaml
metrics: # Metrics to collect (Optional, all if unset)
    - uptime_in_seconds
    - connected_clients
    - client_longest_output_list
    - client_biggest_input_buf
    - blocked_clients
    - used_memory
    - used_memory_rss
    - used_memory_peak
    - total_system_memory
    - used_memory_lua
    - maxmemory
    - mem_fragmentation_ratio
    - loading
    - rdb_changes_since_last_save
    - rdb_bgsave_in_progress
    - rdb_last_save_time
    - rdb_last_bgsave_status
    - rdb_last_bgsave_time_sec
    - rdb_current_bgsave_time_sec
    - aof_enabled
    - aof_rewrite_in_progress
    - aof_rewrite_scheduled
    - aof_last_rewrite_time_sec
    - aof_current_rewrite_time_sec
    - aof_last_bgrewrite_status
    - aof_last_write_status
    - total_connections_received
    - total_commands_processed
    - instantaneous_ops_per_sec
    - total_net_input_bytes
    - total_net_output_bytes
    - instantaneous_input_kbps
    - instantaneous_output_kbps
    - rejected_connections
    - sync_full
    - sync_partial_ok
    - sync_partial_err
    - expired_keys
    - evicted_keys
    - keyspace_hits
    - keyspace_misses
    - pubsub_channels
    - pubsub_patterns
    - latest_fork_usec
    - migrate_cached_sockets
    - connected_slaves
    - master_repl_offset
    - repl_backlog_active
    - repl_backlog_size
    - repl_backlog_first_byte_offset
    - repl_backlog_histlen
    - used_cpu_sys
    - used_cpu_user
    - used_cpu_sys_children
    - used_cpu_user_children
    - cluster_enabled
```

#### Labels
Redis-Exporter can add static labels to collected metrics. A *label* is defined as follow:
``` yaml
labels: # Labels definitions (Optional)
  label_name: label_value # Label definition             (Required)
```

#### Parameters
Redis-Exporter can be customized through parameters. See available parameters bellow:
``` yaml
# parameters definitions (Optional)
scanDuration: 1000 # Duration within all the sources should be scraped (Optional, default: 1000)
maxConcurrent: 200 # Max concurrent scrape allowed (Optional, default: 50)
scrapeTimeout: 5000 # Stats fetch timeout (Optional, default: 5000)
flushPath: /opt/beamium/sinks/warp- # Path to flush metrics + filename header (Optional, default: no flush)
flushPeriod: 10000 # Flush period (Optional, 10000)
```
