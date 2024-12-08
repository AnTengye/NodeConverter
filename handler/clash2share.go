package handler

import (
	"fmt"
	"strings"

	"github.com/AnTengye/NodeConvertor/core"
	"github.com/go-resty/resty/v2"
	"github.com/kataras/iris/v12"
	"gopkg.in/yaml.v3"
)

type ClashUrl struct {
	Url string `url:"url,required"`
}

type clashNode struct {
	Proxies []Proxy `yaml:"proxies"`
}

type Proxy struct {
	Type string `yaml:"type"`
}

func ClashToShare(ctx iris.Context) {
	var nu ClashUrl
	logger := ctx.Application().Logger()
	if err := ctx.ReadQuery(&nu); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}
	client := resty.New()
	clashResp, err := client.R().Get(nu.Url)
	if err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}
	var data map[string]interface{}
	if err := yaml.Unmarshal(clashResp.Body(), &data); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	proxies, ok := data["proxies"].([]interface{})
	if !ok {
		ctx.StopWithError(iris.StatusBadRequest, fmt.Errorf("proxies field not found"))
		return
	}
	result := make([]string, 0, len(proxies))
	for _, proxy := range proxies {
		proxyMap, ok := proxy.(map[string]interface{})
		if !ok {
			logger.Warn("proxy format is invalid")
			continue
		}

		// Get the proxy type
		proxyType, ok := proxyMap["type"].(string)
		if !ok {
			logger.Warn("proxy type is not a string")
			continue
		}

		// Convert the proxy map to []byte
		proxyBytes, _ := yaml.Marshal(proxyMap)

		var node core.Node
		switch proxyType {
		case "vless":
			node = core.NewVLESSNode()
		default:
			logger.Warnf("%s is not supported", proxyType)
			continue
		}
		err = node.FromClash(proxyBytes)
		if err != nil {
			logger.Errorf("convert failed: %v", err)
			continue
		}
		result = append(result, node.ToShare())
	}
	ctx.WriteString(strings.Join(result, "\n"))
}
