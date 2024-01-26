package clash

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"myrouter/comm/logger"
	"myrouter/comm/myauth"
	"myrouter/models"
	"net/http"
)

// GetRules 获取所有的自定义规则
//
// GET /api/clash/rules/all
//
// 返回 Result[[]string]
func GetRules(c *gin.Context) {
	// 读取本地自定义的规则
	rules, err := getRules()
	if err != nil {
		logger.Error.Printf("[%s]读取自定义规则出错：%s\n", c.GetString(myauth.KeyUser), err)
		c.JSON(http.StatusOK, models.Result{Code: 2000, Msg: fmt.Sprintf("读取自定义规则出错：%s", err)})
		return
	}

	logger.Info.Printf("[%s]已读取自定义规则\n", c.GetString(myauth.KeyUser))
	c.JSON(http.StatusOK, models.Result{Code: 0, Msg: "已读取自定义规则", Data: rules})
}

// AddRule 添加一条自定义规则
//
// POST /api/clash/rule/add
//
// JSON 表单数据类型为 PostData[string]
func AddRule(c *gin.Context) {
	// 解析 JSON
	var rule models.PostData[string]
	err := c.BindJSON(&rule)
	if err != nil {
		logger.Error.Printf("[%s]解析请求中的数据出错：%s\n", c.GetString(myauth.KeyUser), err)
		c.JSON(http.StatusOK, models.Result{Code: 2100, Msg: fmt.Sprintf("解析请求中的数据出错：%s", err)})
		return
	}

	// 添加规则
	err = addRule(rule.Data)
	if err != nil {
		logger.Error.Printf("[%s]添加自定义规则'%s'出错：%s\n", c.GetString(myauth.KeyUser), rule.Data, err)
		c.JSON(http.StatusOK, models.Result{Code: 2110, Msg: fmt.Sprintf("添加自定义规则出错：%s", err)})
		return
	}

	go ExecClash(c, cmdStart, "重启 Clash ")

	logger.Info.Printf("[%s]已添加自定义规则'%s'\n", c.GetString(myauth.KeyUser), rule.Data)
	c.JSON(http.StatusOK, models.Result{Code: 0, Msg: "已添加自定义规则"})
}

// DelRule 删除一条自定义规则
//
// POST /api/clash/rule/del
//
// JSON 表单数据类型为 PostData[string]
//
// 传递的规则`data`必须与规则文件中的对应的一行完全一致（即使是空格等多余符号）
func DelRule(c *gin.Context) {
	// 解析 JSON
	var rule models.PostData[string]
	err := c.BindJSON(&rule)
	if err != nil {
		logger.Error.Printf("[%s]解析请求中的数据出错：%s\n", c.GetString(myauth.KeyUser), err)
		c.JSON(http.StatusOK, models.Result{Code: 2200, Msg: fmt.Sprintf("解析请求中的数据出错：%s", err)})
		return
	}

	// 删除规则
	err = delRule(rule.Data)
	if err != nil {
		logger.Error.Printf("[%s]删除自定义规则出错：%s\n", c.GetString(myauth.KeyUser), err)
		c.JSON(http.StatusOK, models.Result{Code: 2210, Msg: fmt.Sprintf("删除自定义规则出错：%s", err)})
		return
	}

	go ExecClash(c, cmdStart, "重启 Clash ")

	logger.Info.Printf("[%s]已删除自定义规则'%s'\n", c.GetString(myauth.KeyUser), rule.Data)
	c.JSON(http.StatusOK, models.Result{Code: 0, Msg: "已删除自定义规则"})
}

// OverrideRules 覆盖所有的自定义规则
//
// 谨慎：会覆盖原有的规则。一定要在原有规则基础上，修改、添加后保存!
//
// POST /api/clash/rules/override
//
// JSON 表单数据类型为 PostData[string]
func OverrideRules(c *gin.Context) {
	// 解析 JSON
	var rules models.PostData[string]
	err := c.BindJSON(&rules)
	if err != nil {
		logger.Error.Printf("[%s]解析请求中的数据出错：%s\n", c.GetString(myauth.KeyUser), err)
		c.JSON(http.StatusOK, models.Result{Code: 2300, Msg: fmt.Sprintf("解析请求中的数据出错：%s", err)})
		return
	}

	// 保存规则
	err = overrideRules(rules.Data)
	if err != nil {
		logger.Error.Printf("[%s]覆盖自定义规则出错：%s\n", c.GetString(myauth.KeyUser), err)
		c.JSON(http.StatusOK, models.Result{Code: 2310, Msg: fmt.Sprintf("覆盖自定义规则出错：%s", err)})
	}

	go ExecClash(c, cmdStart, "重启 Clash ")

	logger.Info.Printf("[%s]已覆盖自定义规则'%s'\n", c.GetString(myauth.KeyUser), rules.Data)
	c.JSON(http.StatusOK, models.Result{Code: 0, Msg: "已覆盖自定义规则"})
}

// BackToLastRules 恢复自定义规则到上次保存的内容
//
// POST /api/clash/rules/backtolast
//
// JSON 表单数据类型为 PostData[bool]
func BackToLastRules(c *gin.Context) {
	if !checkNextStep(c) {
		return
	}

	// 恢复规则
	err := backToLastRules()
	if err != nil {
		logger.Error.Printf("[%s]恢复到上次保存的自定义规则出错：%s\n", c.GetString(myauth.KeyUser), err)
		c.JSON(http.StatusOK, models.Result{Code: 2420, Msg: fmt.Sprintf("恢复到上次保存的自定义规则出错：%s", err)})
		return
	}

	go ExecClash(c, cmdStart, "重启 Clash ")

	logger.Info.Printf("[%s]已恢复到上次保存的自定义规则\n", c.GetString(myauth.KeyUser))
	c.JSON(http.StatusOK, models.Result{Code: 0, Msg: "已恢复到上次保存的自定义规则"})
}

// GetProxyGroups 读取所有的代理组
//
// GET /api/clash/config/proxygroups
//
// 返回 Result[[]string]
func GetProxyGroups(c *gin.Context) {
	proxygroups, err := getProxyGroups()
	if err != nil {
		logger.Error.Printf("[%s]读取所有代理组出错：%s\n", c.GetString(myauth.KeyUser), err)
		c.JSON(http.StatusOK, models.Result{Code: 2500, Msg: fmt.Sprintf("读取所有代理组出错：%s", err)})
		return
	}

	logger.Info.Printf("[%s]已读取所有代理组\n", c.GetString(myauth.KeyUser))
	c.JSON(http.StatusOK, models.Result{Code: 0, Msg: "已读取所有代理组", Data: proxygroups})
}

// GetClashRenderData 获取客户端渲染Clash部分网页所需的数据
//
// GET /api/clash/data/render
//
// 返回 Result[RenderData]
func GetClashRenderData(c *gin.Context) {
	// 读取本地自定义的规则
	rules, err := getRules()
	if err != nil {
		logger.Error.Printf("[%s]读取自定义规则出错：%s\n", c.GetString(myauth.KeyUser), err)
		c.JSON(http.StatusOK, models.Result{Code: 2600, Msg: fmt.Sprintf("读取自定义规则出错：%s", err)})
		return
	}

	// 读取所有的代理组
	proxygroups, err := getProxyGroups()
	if err != nil {
		logger.Error.Printf("[%s]读取所有代理组出错：%s\n", c.GetString(myauth.KeyUser), err)
		c.JSON(http.StatusOK, models.Result{Code: 2610, Msg: fmt.Sprintf("读取所有代理组出错：%s", err)})
		return
	}

	data := models.RenderData{Rules: rules, ProxyGroups: proxygroups}
	logger.Info.Printf("[%s]已读取渲染Clash部分网页所需的数据\n", c.GetString(myauth.KeyUser))
	c.JSON(http.StatusOK, models.Result{Code: 0, Msg: "已读取渲染Clash部分网页所需的数据", Data: data})
}

// StartClash 开始/重启 Clash
//
// POST /api/clash/manager/start 最后的一个值'start'需要和后端的`clash.go`中的`Cmd`值一致，才能正确发起请求
//
// JSON 表单数据类型为 PostData[bool]
func StartClash(c *gin.Context) {
	if !checkNextStep(c) {
		return
	}

	go ExecClash(c, cmdStart, "重启 Clash ")

	logger.Info.Printf("[%s]正在重启 Clash...\n", c.GetString(myauth.KeyUser))
	c.JSON(http.StatusOK, models.Result{Code: 0, Msg: "正在重启 Clash..."})
}

// StopClash 停止 Clash
//
// POST /api/clash/manager/stop 最后的一个值'stop'需要和后端的`clash.go`中的`Cmd`值一致，才能正确发起请求
//
// JSON 表单数据类型为 PostData[bool]
func StopClash(c *gin.Context) {
	if !checkNextStep(c) {
		return
	}

	go ExecClash(c, cmdStop, "停止 Clash ")

	logger.Info.Printf("[%s]正在停止 Clash...\n", c.GetString(myauth.KeyUser))
	c.JSON(http.StatusOK, models.Result{Code: 0, Msg: "正在停止 Clash..."})
}

// 检查是否继续操作
func checkNextStep(c *gin.Context) bool {
	// 解析 JSON
	var data models.PostData[bool]
	err := c.BindJSON(&data)
	if err != nil {
		logger.Error.Printf("[%s]解析请求中的数据出错'%s'：%s\n", c.GetString(myauth.KeyUser), c.FullPath(), err)
		c.JSON(http.StatusOK, models.Result{Code: 10, Msg: fmt.Sprintf("解析请求中的数据出错：%s", err)})
		return false
	}

	// 仅在为 true 时恢复到上次的规则
	if !data.Data {
		logger.Error.Printf("[%s]根据参数'%t'，不继续执行\n", c.GetString(myauth.KeyUser), data.Data)
		c.JSON(http.StatusOK, models.Result{Code: 11, Msg: fmt.Sprintf("根据参数'%t'，不继续执行", data.Data)})
		return false
	}

	return true
}
