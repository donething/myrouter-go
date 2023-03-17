#!/bin/bash

# 在路由器中自启动脚本，参考 README.md 中自启动的部分

echo "将在N秒后，等路由器连上网络，即开始执行脚本……"
sleep 120
cd /root/do/bin/myrouter || exit
chmod +x ./myrouter
nohup ./myrouter >/dev/null 2>errs.log &
