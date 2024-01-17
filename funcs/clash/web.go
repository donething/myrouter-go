package clash

import (
	"encoding/json"
	"fmt"
	"myrouter/comm/logger"
	"myrouter/models"
	"net/http"
)

// PostData 请求保存自定义规则时的 JSON 表单数据
type PostData[T any] struct {
	Data T `json:"data"`
}

// GetRules 获取自定义规则
//
// GET /api/clash/rules/get
func GetRules(w http.ResponseWriter, _ *http.Request) {
	// 读取本地规则
	var result models.Result
	bs, err := getRules()
	if err != nil {
		logger.Error.Printf("读取自定义规则出错：%s\n", err)
		result = models.Result{Code: 2000, Msg: fmt.Sprintf("读取自定义规则出错：%s", err)}
	} else {
		result = models.Result{Code: 0, Msg: "自定义规则", Data: bs}
	}

	// 响应结果
	result.Response(w)
}

// SaveRules 保存自定义规则
//
// 谨慎：会覆盖原有的规则。一定要在原有规则基础上，修改、添加后保存!
//
// POST /api/clash/rules/save
//
// JSON 表单数据类型为 PostData[string]
func SaveRules(w http.ResponseWriter, r *http.Request) {
	// 解析 JSON
	var rules PostData[string]
	err := json.NewDecoder(r.Body).Decode(&rules)
	if err != nil {
		logger.Error.Printf("解析请求中的数据出错：%s\n", err)
		result := models.Result{Code: 2100, Msg: fmt.Sprintf("解析请求中的数据出错：%s", err)}
		result.Response(w)
		return
	}

	// 保存规则
	err = saveRules(rules.Data)
	var result models.Result
	if err != nil {
		logger.Error.Printf("保存自定义规则出错：%s\n", err)
		result = models.Result{Code: 2110, Msg: fmt.Sprintf("保存自定义规则出错：%s", err)}
	} else {
		result = models.Result{Code: 0, Msg: "已保存自定义规则"}
	}

	// 响应结果
	result.Response(w)
}

// BackToLastRules 恢复自定义规则到上次保存的内容
//
// POST /api/clash/rules/backtolast
//
// JSON 表单数据类型为 PostData[bool]
func BackToLastRules(w http.ResponseWriter, r *http.Request) {
	// 解析 JSON
	var backto PostData[bool]
	err := json.NewDecoder(r.Body).Decode(&backto)
	if err != nil {
		logger.Error.Printf("解析请求中的数据出错：%s\n", err)
		result := models.Result{Code: 2200, Msg: fmt.Sprintf("解析请求中的数据出错：%s", err)}
		result.Response(w)
		return
	}

	// 仅在为 true 时恢复到上次的规则
	if !backto.Data {
		logger.Info.Printf("根据传递的参数，不恢复自定义规则：'%s'\n", backto.Data)
		result := models.Result{Code: 2200, Msg: fmt.Sprintf("根据传递的参数，不恢复自定义规则：'%v'", backto.Data)}
		result.Response(w)
		return
	}

	// 恢复规则
	err = backToLastRules()
	var result models.Result
	if err != nil {
		logger.Error.Printf("恢复到上次保存的自定义规则出错：%s\n", err)
		result = models.Result{Code: 2210, Msg: fmt.Sprintf("恢复到上次保存的自定义规则出错：%s", err)}
	} else {
		result = models.Result{Code: 0, Msg: "已恢复到上次保存的自定义规则"}
	}

	// 响应结果
	result.Response(w)
}
