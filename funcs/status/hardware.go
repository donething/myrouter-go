package status

import (
	"fmt"
	"github.com/donething/utils-go/dotext"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"myrouter/comm/logger"
	"myrouter/config"
	"myrouter/models"
	"time"
)

// GetRouterStatus 获取 CPU、内存的状态
func GetRouterStatus() *models.RouterStatus {
	// IP
	ip, err := GetLocalIPAddr()
	if err != nil {
		logger.Error.Printf("获取本地 IP 地址出错：%s\n", err)
	}

	// 内存
	v, err := mem.VirtualMemory()
	if err != nil {
		logger.Error.Printf("获取内存使用情况出错：%s\n", err)
	}
	m := models.Mem{
		Total:       dotext.BytesHumanReadable(v.Total),
		Available:   dotext.BytesHumanReadable(v.Available),
		Used:        dotext.BytesHumanReadable(v.Used),
		UsedPercent: v.UsedPercent,
	}

	// CPU
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		logger.Error.Printf("获取 CPU 使用情况出错：%s\n", err)
	}
	c := models.CPU{UsedPercent: fmt.Sprintf("%.2f", percent[0])}

	// 所有需要的状态
	return &models.RouterStatus{
		IPAddr: *ip,
		Mem:    m,
		CPU:    c,

		Logo: config.Conf.Router.Logo,
	}
}
