package clash

import (
	"errors"
	"fmt"
	"github.com/donething/utils-go/dofile"
	"io"
	. "myrouter/config"
	"os"
)

// 备份文件时，在末尾添加的标识
const bakTag = ".bak"

// 获取自定义规则
func getRules() (string, error) {
	bs, err := os.ReadFile(Conf.Clash.RulesPath)
	// 还没有自定义规则，即文件不存在时，不作为错误只返回空字符串""
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}

		return "", err
	}

	return string(bs), nil
}

// 保存自定义规则。将提前备份原文件为"*.bak"
//
// 谨慎：会覆盖原有的规则。一定要在原有规则基础上，修改、添加后保存!
func saveRules(rules string) error {
	// 先备份原规则文件
	srcFile, err := os.OpenFile(Conf.Clash.RulesPath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	dstFile, err := os.Create(fmt.Sprintf("%s%s", Conf.Clash.RulesPath, bakTag))
	if err != nil {
		srcFile.Close()
		return err
	}
	defer dstFile.Close()

	// 开始备份
	_, err = io.Copy(dstFile, srcFile)
	srcFile.Close()
	if err != nil {
		return err
	}

	// 写入新规则
	srcFile, err = os.OpenFile(Conf.Clash.RulesPath, os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	_, err = srcFile.WriteString(rules)
	if err != nil {
		return err
	}

	return nil
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
