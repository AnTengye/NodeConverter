package handler

import (
	"github.com/AnTengye/NodeConvertor/core"
	"testing"
)

// 模拟 core.Node 结构
type Node struct {
	name string
}

func (n *Node) Type() core.NodeType {
	panic("implement me")
}

func (n *Node) ToShare() string {
	panic("implement me")
}

func (n *Node) ToClash() string {
	panic("implement me")
}

func (n *Node) FromShare(s string) error {
	panic("implement me")
}

func (n *Node) FromClash(bytes []byte) error {
	panic("implement me")
}

func (n *Node) Name() string {
	return n.name
}

func (n *Node) SetName(name string) {
	n.name = name
}

// 单元测试
func TestFilterNodes(t *testing.T) {
	// 初始化测试数据
	nodes := []core.Node{
		&Node{name: "剩余流量：698.03 GB"},
		&Node{name: "套餐到期：2027-02-04"},
		&Node{name: "🇭🇰 香港 IEPL 01 | x4"},
		&Node{name: "🇭🇰 香港 IEPL 02 | x4"},
		&Node{name: "🇸🇬 新加坡 IEPL 01 | x4"},
		&Node{name: "🇸🇬 新加坡 IEPL 02 | x4"},
		&Node{name: "🇯🇵 日本 IEPL 01 | x4"},
		&Node{name: "🇯🇵 日本 IEPL 02 | x4"},
		&Node{name: "🇨🇳 台湾 IEPL 01 | x4"},
		&Node{name: "🇨🇳 台湾 IEPL 02 | x4"},
		&Node{name: "🇭🇰 香港 01"},
		&Node{name: "🇭🇰 香港 02"},
	}

	// 初始化过滤条件
	filterCondition := FilterCondition{
		ExcludeReg: "(到期|剩余流量|时间|官网|产品|平台)",
		IncludeReg: "",
		RenameReg:  "(.*)@[xf]$1",
	}

	// 调用过滤函数
	filteredNodes, err := filterNodes(nodes, filterCondition)
	if err != nil {
		t.Fatalf("filterNodes failed: %v", err)
	}

	// 验证过滤结果
	expectedNodes := []string{
		"[xf]🇭🇰 香港 IEPL 01 | x4",
		"[xf]🇭🇰 香港 IEPL 02 | x4",
		"[xf]🇸🇬 新加坡 IEPL 01 | x4",
		"[xf]🇸🇬 新加坡 IEPL 02 | x4",
		"[xf]🇯🇵 日本 IEPL 01 | x4",
		"[xf]🇯🇵 日本 IEPL 02 | x4",
		"[xf]🇨🇳 台湾 IEPL 01 | x4",
		"[xf]🇨🇳 台湾 IEPL 02 | x4",
		"[xf]🇭🇰 香港 01",
		"[xf]🇭🇰 香港 02",
	}

	if len(filteredNodes) != len(expectedNodes) {
		t.Fatalf("expected %d nodes, got %d", len(expectedNodes), len(filteredNodes))
	}

	for i, node := range filteredNodes {
		if node.Name() != expectedNodes[i] {
			t.Errorf("expected node %d to be %s, got %s", i, expectedNodes[i], node.Name())
		}
	}
}
