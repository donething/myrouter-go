#!/bin/sh /etc/rc.common
START=99
STOP=50

DESCRIPTION="MyRouter Start Service"
cmd="/data/myrouter/myrouter -c /data/myrouter/myrouter.json"
name="myrouter"

pid_file="/var/run/$name.pid"
stdout_log="/var/log/$name.log"
stderr_log="/var/log/$name.err"

get_pid() {
    cat "$pid_file"
}

is_running() {
    [ -f "$pid_file" ] && cat /proc/$(get_pid)/stat > /dev/null 2>&1
}

start() {
	if is_running; then
		echo "$name Already started"
	else
		echo "Starting $name"

		$cmd >> "$stdout_log" 2>> "$stderr_log" &
		echo $! > "$pid_file"
		if ! is_running; then
			echo "Unable to start, see $stdout_log and $stderr_log"
			exit 1
		fi
	fi
}

stop() {
	if is_running; then
		echo -n "Stopping $name.."
		kill $(get_pid)
		for i in $(seq 1 10)
		do
			if ! is_running; then
				break
			fi
			echo -n "."
			sleep 1
		done
		echo
		if is_running; then
			echo "$name Not stopped; may still be shutting down or shutdown may have failed"
			exit 1
		else
			echo "$name Stopped"
			if [ -f "$pid_file" ]; then
				rm "$pid_file"
			fi
		fi
	else
		echo "$name Not running"
	fi
}

restart() {
	stop
	if is_running; then
		echo "$name Unable to stop, will not attempt to start"
		exit 1
	fi
	start
}
