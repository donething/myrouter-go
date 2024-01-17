#!/bin/bash

### 在路由器中自启动脚本，适用于 RedmiAx6000
# 1. 将本脚本上传到路由器上，注意文件名要为`myrouter_redmi.sh`，和下面的完整执行路径一致
# 2. 需在路由器中打开`/data/auto_ssh/auto_ssh.sh`，
# 3. 在其最后一行添加`sh /data/myrouter/myrouter_redmi.sh`
# 4. 最后，加上执行权限`chmod +x /data/myrouter/myrouter_redmi.sh`
#
# 参考 https://blog.csdn.net/weixin_45945615/article/details/130319222

cd /data/myrouter
chmod +x ./myrouter
./myrouter >> /dev/null 2>&1 &
