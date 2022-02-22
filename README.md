# OpenRC exporter

Prometheus exporter which exposes metrics on services managed by [OpenRC].

[OpenRC]: https://github.com/OpenRC/openrc/

# Dependencies

openrc-exporter uses [cgo] for binding to librc of openrc.

[cgo]: https://golang.org/cmd/cgo/

- openrc-dev

# Metrics

The following metrics are exposed in addition to the default golang metrics.

**openrc_runlevel_current**

The current runlevel

**openrc_service_daemons_crashed**

1 if the daemons started with `start-stop-daemon` are crashed, 0 otherwise

**openrc_service_status**

The states that a service is in. A service can be in multiple state at the same time. (1 means in state, 0 means not in state)

**openrc_service_uptime_seconds**

Number of seconds that the service has been running. Only available for services with `supervisor=supervise-daemon`.

**openrc_service_start_count**

Number of times a service has been restarted

**openrc_service_respawn_delay_seconds**

Number of seconds to wait before restarting a process that crashed. Only available for services with `supervisor=supervise-daemon`

**openrc_service_respawn_max**

Maximum number of restarts within `respawn_period` before giving up. Only available for services with `supervisor=supervise-daemon`

**openrc_service_respawn_period_seconds**

The period within which respawn counts towards respawn_max. Only available for services with `supervisor=supervise-daemon`


# Example

```
# HELP openrc_runlevel_current The current runlevel
# TYPE openrc_runlevel_current gauge
openrc_runlevel_current{runlevel="default"} 1
# HELP openrc_service_daemons_crashed Whether the daemons started with start-stop-daemon are crashed
# TYPE openrc_service_daemons_crashed gauge
openrc_service_daemons_crashed{service="alertmanager"} 0
openrc_service_daemons_crashed{service="znc"} 0
# HELP openrc_service_respawn_delay_seconds Number of seconds to wait before restarting a process that crashed
# TYPE openrc_service_respawn_delay_seconds gauge
openrc_service_respawn_delay_seconds{service="alertmanager"} 10
openrc_service_respawn_delay_seconds{service="znc"} 0
# HELP openrc_service_respawn_max Maximum number of restarts within respawn_period before giving up (supervise-daemon only)
# TYPE openrc_service_respawn_max gauge
openrc_service_respawn_max{service="alertmanager"} 3
openrc_service_respawn_max{service="znc"} 0
# HELP openrc_service_respawn_period_seconds The period within which respawn counts towards respawn_max (supervise-daemon only)
# TYPE openrc_service_respawn_period_seconds gauge
openrc_service_respawn_period_seconds{service="alertmanager"} 20
openrc_service_respawn_period_seconds{service="znc"} 0
# HELP openrc_service_start_count Number of times a service has been restarted
# TYPE openrc_service_start_count counter
openrc_service_start_count{service="alertmanager"} 0
openrc_service_start_count{service="znc"} 0
# HELP openrc_service_status States that the service is in (1 means the service is in the state, 0 otherwise)
# TYPE openrc_service_status gauge
openrc_service_status{service="alertmanager",state="crashed"} 0
openrc_service_status{service="alertmanager",state="failed"} 0
openrc_service_status{service="alertmanager",state="hotplugged"} 0
openrc_service_status{service="alertmanager",state="inactive"} 0
openrc_service_status{service="alertmanager",state="scheduled"} 0
openrc_service_status{service="alertmanager",state="started"} 1
openrc_service_status{service="alertmanager",state="starting"} 0
openrc_service_status{service="alertmanager",state="stopped"} 0
openrc_service_status{service="alertmanager",state="stopping"} 0
openrc_service_status{service="alertmanager",state="wasinactive"} 0
openrc_service_status{service="znc",state="crashed"} 0
openrc_service_status{service="znc",state="failed"} 0
openrc_service_status{service="znc",state="hotplugged"} 0
openrc_service_status{service="znc",state="inactive"} 0
openrc_service_status{service="znc",state="scheduled"} 0
openrc_service_status{service="znc",state="started"} 1
openrc_service_status{service="znc",state="starting"} 0
openrc_service_status{service="znc",state="stopped"} 0
openrc_service_status{service="znc",state="stopping"} 0
openrc_service_status{service="znc",state="wasinactive"} 0
# HELP openrc_service_uptime_seconds Number of seconds that the service is running (0 if not running)
# TYPE openrc_service_uptime_seconds gauge
openrc_service_uptime_seconds{service="alertmanager"} 2.0787426715027332e+07
openrc_service_uptime_seconds{service="znc"} 0
```

# Usage

```
Usage of openrc-exporter:
  -listen-address string
    	Listening address (default ":9816")
  -version
    	Print version and exit
```

# License

The license is [AGPL-3.0-only]. See [LICENSE](LICENSE).

[AGPL-3.0-only]: https://spdx.org/licenses/AGPL-3.0-only.html

# Contributing

You can send patches to
[~tomleb/public-inbox@lists.sr.ht](https://lists.sr.ht/~tomleb/public-inbox).
