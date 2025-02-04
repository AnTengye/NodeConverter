package handler

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/AnTengye/NodeConvertor/core"
	"github.com/AnTengye/NodeConvertor/lib/network"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func checkUrlOrShare(urlOrShare string) NodeUrlType {
	if strings.HasPrefix(urlOrShare, "http") {
		return subUrl
	}
	return shareUrl
}

func fetchNodes(u string) ([]core.Node, error) {
	var nodes []core.Node
	switch checkUrlOrShare(u) {
	case subUrl:
		response, err := network.CacheGET(u)
		if err != nil {
			return nil, fmt.Errorf("download from sub error: %v", err)
		}
		nodes, err = handlerSubResponse(response)
		if err != nil {
			network.DeleteCache(u)
			return nil, fmt.Errorf("get sub error: %v", err)
		}
	case shareUrl:
		node, err := handlerSingleShareUrl(u)
		if err != nil {
			return nil, fmt.Errorf("get share error: %v", err)
		}
		nodes = append(nodes, node)
	default:
		return nil, fmt.Errorf("url is invalid")
	}
	return nodes, nil
}

func handlerSubResponse(urlResponse []byte) ([]core.Node, error) {
	sDec, err := base64.StdEncoding.DecodeString(string(urlResponse))
	if err != nil {
		zap.S().Infow("decode base64 error, change to clash decode", "error", err)
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
		proxyType, ok := proxyMap["type"].(string)
		if !ok {
			zap.S().Warnw("proxy type is not a string", zap.Any("proxy", proxyMap["type"]))
			continue
		}

		// Convert the proxy map to []byte
		proxyBytes, _ := yaml.Marshal(proxyMap)

		var node core.Node
		switch core.NodeType(proxyType) {
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
