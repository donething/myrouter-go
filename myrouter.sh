#!/bin/bash

# 在路由器中自启动脚本，参考 README.md 中自启动的部分

# 先判断该程序是否以运行，已运行就不重复运行
PROC_NAME=myrouter
ProcNumber=`ps -ef |grep -w $PROC_NAME|grep -v grep|wc -l`
if [ $ProcNumber -le 0 ];then
	echo "将在N秒后，等路由器连上网络，即开始执行脚本……"
	sleep 30
	cd /home/donet/bin/myrouter
	chmod +x myrouter
	nohup ./myrouter &
else
   echo "myrouter 已运行，不需重复运行"
fi
