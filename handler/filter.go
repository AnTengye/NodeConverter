package handler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/AnTengye/NodeConverter/core"
)

type FilterCondition struct {
	IncludeReg string
	ExcludeReg string
	RenameReg  string
}

func filterNodes(nodes []core.Node, filterCondition FilterCondition) ([]core.Node, error) {
	var (
		includeReg, excludeReg, renameReg *regexp.Regexp
		renameReplacement                 string
		err                               error
	)
	if filterCondition.IncludeReg != "" {
		includeReg, err = regexp.Compile(filterCondition.IncludeReg)
		if err != nil {
			return nil, fmt.Errorf("include reg error: %v", err)
		}
	}
	if filterCondition.ExcludeReg != "" {
		excludeReg, err = regexp.Compile(filterCondition.ExcludeReg)
		if err != nil {
			return nil, fmt.Errorf("exclude reg error: %v", err)
		}
	}
	if filterCondition.RenameReg != "" {
		// 分割正则表达式和替换字符串
		parts := strings.Split(filterCondition.RenameReg, "@")
		if len(parts) != 2 {
			return nil, fmt.Errorf("rename reg error: need @")
		}

		pattern := parts[0]
		renameReplacement = parts[1]
		renameReg, err = regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("rename reg error: %v", err)
		}
	}
	filteredNodes := make([]core.Node, 0, len(nodes))
	for _, node := range nodes {
		if excludeReg != nil && excludeReg.MatchString(node.Name()) {
			continue
		}
		if includeReg != nil && !includeReg.MatchString(node.Name()) {
			continue
		}
		if renameReg != nil {
			node.SetName(renameReg.ReplaceAllString(node.Name(), renameReplacement))
		}
		filteredNodes = append(filteredNodes, node)
	}
	return filteredNodes, nil
}
