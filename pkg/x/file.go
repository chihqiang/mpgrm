package x

import (
	"github.com/samber/lo"
	"path/filepath"
	"strings"
)

func MatchedFiles(patterns []string) []string {
	var files []string
	exclude := make(map[string]struct{}) // 排除的文件列表（目前没有使用）

	for _, p := range patterns {
		p = strings.TrimSpace(p) // 去掉前后空格
		if p == "" {
			continue // 如果为空则跳过
		}
		// 处理通配符匹配
		matched, err := filepath.Glob(p)
		if err != nil {
			continue // 匹配出错则跳过
		}
		files = append(files, matched...) // 将匹配到的文件加入结果列表
	}
	// 转成绝对路径并过滤掉 exclude 的文件
	var result []string
	for _, f := range files {
		abs, _ := filepath.Abs(f) // 转换为绝对路径，忽略错误
		if _, ok := exclude[abs]; !ok {
			result = append(result, abs) // 不在排除列表的文件加入最终结果
		}
	}
	return lo.Uniq(result)
}
