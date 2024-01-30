#!/bin/sh /etc/rc.common

USE_PROCD=1
START=95

start_service() {
    procd_open_instance
    procd_set_param command /data/myrouter/myrouter -c /data/myrouter/myrouter.json
    procd_set_param stdout 1
    procd_set_param stderr 1
    procd_set_param file /var/log/myrouter.log
    procd_close_instance
}
