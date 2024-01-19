package clash

import (
	"gopkg.in/yaml.v3"
	"myrouter/config"
	"myrouter/models"
	"os"
)

// 解析 clash 的配置文件
func parseConfig() (*models.ClashConfig, error) {
	yamlFile, err := os.ReadFile(config.Conf.Clash.ConfigPath)
	if err != nil {
		return nil, err
	}

	var clashConfig models.ClashConfig
	err = yaml.Unmarshal(yamlFile, &clashConfig)
	if err != nil {
		return nil, err
	}

	return &clashConfig, nil
}

// 获取配置中所有代理组的名字。用于在规则中指定分流
func getProxyGroups() ([]string, error) {
	conf, err := parseConfig()
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(conf.ProxyGroups)+2)
	// 先添加默认的两个规则
	names = append(names, "DIRECT", "REJECT")
	for _, proxyGroup := range conf.ProxyGroups {
		names = append(names, proxyGroup.Name)
	}

	return names, nil
}
