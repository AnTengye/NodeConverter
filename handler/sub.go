package handler

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/AnTengye/NodeConverter/core"
	"github.com/AnTengye/NodeConverter/lib/network"
	"github.com/kataras/iris/v12"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type NodeUrlType int

type FilterNodeFunc func(node core.Node) bool

const (
	subUrl   NodeUrlType = 1
	shareUrl NodeUrlType = 2
)

// target 必要	surge&ver=4	指想要生成的配置类型，详见上方 支持类型 中的参数
// url	必要	https%3A%2F%2Fwww.xxx.com	指机场所提供的订阅链接或代理节点的分享链接，需要经过 URLEncode 处理
// config	可选	https%3A%2F%2Fwww.xxx.com	指 外部配置 的地址 (包含分组和规则部分)，需要经过 URLEncode 处理，详见 外部配置 ，当此参数不存在时使用 程序的主程序目录中的配置文件
type SubReq struct {
	Target  string `url:"target,required"`
	Url     string `url:"url,required"`
	Config  string `url:"config,omitempty"`
	Exclude string `url:"exclude,omitempty"`
	Include string `url:"include,omitempty"`
	Rename  string `url:"rename,omitempty"`
}

func Sub(ctx iris.Context) {
	var req SubReq
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}
	reqUrls := strings.Split(req.Url, "|")
	var (
		nodes []core.Node
		err   error
	)
	for _, url := range reqUrls {
		// 获取
		zap.S().Debugw("fetch url", "url", url)
		n, fetchErr := fetchNodes(url)
		if fetchErr != nil {
			ctx.StopWithError(iris.StatusBadRequest, fmt.Errorf("fetch nodes[%s] error: %v", url, fetchErr))
			return
		}
		nodes = append(nodes, n...)
	}

	// 过滤
	condition := FilterCondition{
		IncludeReg: viper.GetString("Common.Include"),
		ExcludeReg: viper.GetString("Common.Exclude"),
		RenameReg:  viper.GetString("Common.Rename"),
	}
	if req.Include != "" {
		condition.IncludeReg = req.Include
	}
	if req.Exclude != "" {
		condition.ExcludeReg = req.Exclude
	}
	if req.Rename != "" {
		condition.RenameReg = req.Rename
	}
	zap.S().Debugw("filter nodes start", zap.Int("nodes", len(nodes)))
	nodes, err = filterNodes(nodes, condition)
	if err != nil {
		ctx.StopWithError(iris.StatusBadRequest, fmt.Errorf("filter nodes error: %v", err))
		return
	}
	zap.S().Debugw("filter nodes end, start convert", zap.Int("nodes", len(nodes)))
	// 转换
	result, err := convertNodes(nodes, req.Target, req.Config)
	if err != nil {
		ctx.StopWithError(iris.StatusBadRequest, fmt.Errorf("convert nodes error: %v", err))
		return
	}
	ctx.WriteString(result)
	zap.S().Debugw("done!", zap.String("url", req.Url))

}

func convertNodes(nodes []core.Node, target string, config string) (string, error) {
	switch target {
	case "clash", "clashmeta":
		clash := core.NewClash(target)
		if config != "" {
			configBytes, err := network.CacheGET(config)
			if err != nil {
				return "", fmt.Errorf("download from config error: %v", err)
			}
			clashACLSSR, err := core.NewClashACLSSRFromBytes(configBytes)
			if err != nil {
				return "", err
			}
			clash.SetACLSSR(clashACLSSR)
		} else {
			err := clash.WithTemplate(viper.GetString("Advanced.TemplateFilePath"))
			if err != nil {
				return "", fmt.Errorf("new clash with template error: %v", err)
			}
		}
		if target == core.ClashKernelClash {
			// clash 不支持vless 过滤vless类型
			filterVless := make([]core.Node, 0, len(nodes))
			for _, v := range nodes {
				if v.Type() == core.NodeTypeVLESS {
					continue
				}
				filterVless = append(filterVless, v)
			}
			nodes = filterVless
		}
		clash.AddProxy(nodes...)
		y, err := clash.ToYaml()
		if err != nil {
			return "", fmt.Errorf("generate clash error: %v", err)
		}
		return y, nil
	case string(core.NodeTypeShadowSocks), string(core.NodeTypeVMess), string(core.NodeTypeTrojan), string(core.NodeTypeVLESS), "auto":
		outputList := make([]string, 0, len(nodes))
		for _, node := range nodes {
			outputList = append(outputList, node.ToShare())
		}
		return base64.StdEncoding.EncodeToString([]byte(strings.Join(outputList, "\n"))), nil
	default:
		return "", fmt.Errorf("target is invalid")
	}
}
