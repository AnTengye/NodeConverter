package handler

import (
	"encoding/base64"
	"fmt"
	"github.com/AnTengye/NodeConvertor/core"
	"github.com/kataras/iris/v12"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"strings"
)

type NodeUrlType int

const (
	subUrl   NodeUrlType = 1
	shareUrl NodeUrlType = 2
)

// target	必要	surge&ver=4	指想要生成的配置类型，详见上方 支持类型 中的参数
// url	必要	https%3A%2F%2Fwww.xxx.com	指机场所提供的订阅链接或代理节点的分享链接，需要经过 URLEncode 处理
// config	可选	https%3A%2F%2Fwww.xxx.com	指 外部配置 的地址 (包含分组和规则部分)，需要经过 URLEncode 处理，详见 外部配置 ，当此参数不存在时使用 程序的主程序目录中的配置文件
type SubReq struct {
	Target string `url:"target,required"`
	Url    string `url:"url,required"`
	Config string `url:"config,omitempty"`
}

func Sub(ctx iris.Context) {
	var req SubReq
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}
	var (
		nodes []core.Node
	)
	switch checkUrlOrShare(req.Url) {
	case subUrl:
		response, err := restyCli.R().Get(req.Url)
		if err != nil {
			ctx.StopWithError(iris.StatusBadRequest, fmt.Errorf("download from sub error: %v", err))
			return
		}
		nodes, err = handlerSubResponse(response.Body())
		if err != nil {
			ctx.StopWithError(iris.StatusBadRequest, fmt.Errorf("get sub error: %v", err))
			return
		}
	case shareUrl:
		node, err := handlerSingleShareUrl(req.Url)
		if err != nil {
			ctx.StopWithError(iris.StatusBadRequest, fmt.Errorf("get share error: %v", err))
			return
		}
		nodes = append(nodes, node)
	default:
		ctx.StopWithError(iris.StatusBadRequest, fmt.Errorf("url is invalid"))
		return
	}
	switch req.Target {
	case "clash":
		clash, err := core.NewClash(viper.GetString("TemplateFilePath"))
		if err != nil {
			ctx.StopWithError(iris.StatusBadRequest, fmt.Errorf("new clash error: %v", err))
			return
		}
		clash.AddProxy(nodes...)
		y, err := clash.ToYaml()
		if err != nil {
			ctx.StopWithError(iris.StatusBadRequest, fmt.Errorf("generate clash error: %v", err))
			return
		}
		ctx.WriteString(y)
		return
	case string(core.NodeTypeShadowSocks), string(core.NodeTypeVMess), string(core.NodeTypeTrojan), string(core.NodeTypeVLESS), "auto":
		outputList := make([]string, 0, len(nodes))
		for _, node := range nodes {
			outputList = append(outputList, node.ToShare())
		}
		ctx.WriteString(base64.StdEncoding.EncodeToString([]byte(strings.Join(outputList, "\n"))))
		return
	default:
		ctx.StopWithError(iris.StatusBadRequest, fmt.Errorf("target is invalid"))
		return
	}
}

func checkUrlOrShare(urlOrShare string) NodeUrlType {
	if strings.HasPrefix(urlOrShare, "http") {
		return subUrl
	}
	return shareUrl
}

func handlerSubResponse(urlResponse []byte) ([]core.Node, error) {
	sDec, err := base64.StdEncoding.DecodeString(string(urlResponse))
	if err != nil {
		return handlerClashResponse(urlResponse)
	}
	// 按行读取
	lines := strings.Split(string(sDec), "\n")
	result := make([]core.Node, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		node, err := handlerSingleShareUrl(line)
		if err != nil {
			return nil, err
		}
		result = append(result, node)
	}
	return result, nil
}

func handlerSingleShareUrl(shareUrl string) (core.Node, error) {
	split := strings.SplitN(shareUrl, "://", 2)
	if len(split) != 2 {
		return nil, fmt.Errorf("share_url is invalid, it should start with xxx://, but got %s", shareUrl)
	}
	var (
		node core.Node
	)
	switch core.NodeType(split[0]) {
	case core.NodeTypeVLESS, core.NodeTypeVMess:
		node = core.NewVLESSNode()
	case core.NodeTypeShadowSocks:
		node = core.NewShadowsocksNode()
	case core.NodeTypeTrojan:
		node = core.NewTrojanNode()
	default:
		return nil, fmt.Errorf("not support protocol: %s", split[0])
	}
	if convertErr := node.FromShare(shareUrl); convertErr != nil {
		return nil, fmt.Errorf("share_url[%s] convert failed: %v", shareUrl, convertErr)
	}
	return node, nil
}

func handlerClashResponse(clashResponse []byte) ([]core.Node, error) {
	var data map[string]interface{}
	if err := yaml.Unmarshal(clashResponse, &data); err != nil {
		return nil, err
	}
	proxies, ok := data["proxies"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("proxies field not found")
	}
	result := make([]core.Node, 0, len(proxies))
	for _, proxy := range proxies {
		proxyMap, ok := proxy.(map[string]interface{})
		if !ok {
			zap.S().Warn("proxy format is invalid")
			continue
		}

		// Get the proxy type
		proxyType, ok := proxyMap["type"].(core.NodeType)
		if !ok {
			zap.S().Warn("proxy type is not a string")
			continue
		}

		// Convert the proxy map to []byte
		proxyBytes, _ := yaml.Marshal(proxyMap)

		var node core.Node
		switch proxyType {
		case core.NodeTypeVLESS, core.NodeTypeVMess:
			node = core.NewVLESSNode()
		case core.NodeTypeShadowSocks:
			node = core.NewShadowsocksNode()
		case core.NodeTypeTrojan:
			node = core.NewTrojanNode()
		default:
			zap.S().Warnf("%s is not supported", proxyType)
			continue
		}
		err := node.FromClash(proxyBytes)
		if err != nil {
			zap.S().Errorf("convert failed: %v", err)
			continue
		}
		result = append(result, node)
	}
	return result, nil
}
