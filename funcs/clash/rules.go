package clash

import (
	"fmt"
	"github.com/donething/utils-go/dofile"
	"io"
	. "myrouter/config"
	"os"
	"strings"
)

// 备份文件时，在末尾添加的标识
const bakTag = ".bak"

var (
	ErrExists = fmt.Errorf("已存在该域名的规则")
)

// 增加一条自定义规则
//
// 已存在该域名的规则时，将返回 ErrExists
func addRule(rule string) error {
	// 判断新规则是否规范
	newRuleParts := strings.Split(rule, ",")
	if len(newRuleParts) < 3 {
		return fmt.Errorf("新规则不规范'%s'", rule)
	}
	newRuleDomain := strings.ToLower(strings.TrimSpace(newRuleParts[1]))

	// 读取规则文件，判断域名部分是否重复
	rules, err := getRules()
	if err != nil {
		return fmt.Errorf("获取自定义规则出错：%w", err)
	}

	for _, r := range rules {
		if strings.TrimSpace(r) == "" {
			continue
		}
		// 先要判断原规则是否规范
		parts := strings.Split(r, ",")
		if len(parts) < 3 {
			return fmt.Errorf("原规则中有不规范的规则'%s'", r)
		}

		// 已存在该规则，返回已存在错误
		if strings.ToLower(strings.TrimSpace(parts[1])) == newRuleDomain {
			return ErrExists
		}
	}

	// 增加换行符
	rules = append(rules, rule)

	// 保存规则
	return overrideRules(strings.Join(rules, "\n"))
}

// 删除一条自定义规则
//
// 传递的规则必须与规则文件中的对应的一行完全一致（即使是空格等多余符号）。如"- DOMAIN-SUFFIX,baidu.com,DIRECT"
func delRule(rule string) error {
	// 读取原规则文件
	bs, err := os.ReadFile(Conf.Clash.RulesPath)
	if err != nil {
		return err
	}

	// 判断是否存在该规则
	rules := strings.Split(string(bs), "\n")
	index := findIndex(rules, rule)
	if index < 0 {
		return fmt.Errorf("自定义规则中没有该规则")
	}

	// 存在，则删除后，保存
	newRules := append(rules[:index], rules[index+1:]...)
	return overrideRules(strings.Join(newRules, "\n"))
}

// 获取所有的自定义规则（不包括'#'开头的注释）
func getRules() ([]string, error) {
	// 使用 OpenFile 打开文件，为了不存在时自动创建
	file, err := os.OpenFile(Conf.Clash.RulesPath, os.O_RDONLY|os.O_CREATE, 0644)
	// 还没有自定义规则，即文件不存在时，不作为错误只返回空字符串""
	if err != nil {
		return nil, err
	}

	bs, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	rulesStr := strings.TrimSpace(string(bs))
	// 当 rulesStr 为 ""时，strings.Split("", "\n") 也将返回 [""]，所以提前返回 []
	if rulesStr == "" {
		return []string{}, nil
	}

	// 排除'#'开头的注释
	all := strings.Split(rulesStr, "\n")
	rules := make([]string, 0, len(all))
	for _, r := range all {
		if strings.HasPrefix(strings.TrimSpace(r), "#") {
			continue
		}

		rules = append(rules, r)
	}

	return rules, nil
}

// 覆盖所有的自定义规则。将提前备份原文件为"*.bak"
//
// 谨慎：会覆盖原有的规则。一定要在原有规则基础上，修改、添加后保存!
func overrideRules(rules string) error {
	// 先备份原规则文件
	err := backupFile(Conf.Clash.RulesPath)
	if err != nil {
		return fmt.Errorf("备份原规则文件'%s'出错：%w", Conf.Clash.RulesPath, err)
	}

	// 写入新规则
	srcFile, err := os.OpenFile(Conf.Clash.RulesPath, os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	_, err = srcFile.WriteString(rules)
	return err
}

// 恢复自定义规则到上次保存的内容
func backToLastRules() error {
	// 判断备份文件是否存在
	var path = Conf.Clash.RulesPath
	var bakPath = fmt.Sprintf("%s%s", path, bakTag)
	exists, err := dofile.Exists(bakPath)
	if err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("备份文件不存在")
	}

	// 恢复
	err = os.Rename(bakPath, path)
	if err != nil {
		return err
	}

	// 完成
	return nil
}

// 备份文件
func backupFile(srcPath string) error {
	// 先备份原规则文件
	srcFile, err := os.OpenFile(srcPath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	dstFile, err := os.Create(fmt.Sprintf("%s%s", srcPath, bakTag))
	if err != nil {
		srcFile.Close()
		return err
	}
	defer dstFile.Close()

	// 开始备份
	_, err = io.Copy(dstFile, srcFile)
	srcFile.Close()

	return err
}

// 发现索引
func findIndex(array []string, searchStr string) int {
	for i, item := range array {
		if item == searchStr {
			return i
		}
	}

	return -1
}
