# myrouter

# 功能

## 重启路由器

POST `/api/reboot`

## 网络唤醒设置

POST `/api/wol`

需要在配置文件中设置 路由器的本地 IP、目标 MAC 地址

还要修改主板、Windows 设置

* 电脑主板
    * 网络唤醒 允许
    * PCIE设备唤醒 允许
* Windows 系统
    * 电源，快速启动 关闭（关闭后，即使关闭计算机，网线水晶头接口处的灯一直亮着）
    * 网络适配器，电源管理
        * 允许此设备唤醒计算机 允许
        * 只允许唤数据包唤醒计算机

# 自启动运行

自启动，参考 [007+CPE刷clnc小白教程](https://yaohuo.me/bbs/book_view.aspx?sitei=1000&classid=203&id=1097747&vpage=&lpage=)

1. 开一个终端（A）使用`adb`连接路由器，进入路由器的`shell`，备份`adbd-init`文件：`cp /etc/init.d/adbd-init /etc/init.d/adbd-init.bak`
2. 再开一个终端（B），工作路径在本地电脑下，拉取`adbd-init`文件到本地：`adb pull /etc/init.d/adbd-init ./`
3. 编辑`adbd-init`，在`start)`的`case`最后面，添加启动脚本：`/home/donet/bin/myrouter/myrouter.sh`
4. 切到终端A，在路由器`shell`中执行：`mount -o remount / /`
5. 切到终端B，将修改后的`adbd-init`推送到路由器中：`adb push ./adbd-init /etc/init.d/`
6. 切换到终端A，设置权限：`chmod 644 /etc/init.d/adbd-init`