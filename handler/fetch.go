package handler

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/AnTengye/NodeConverter/core"
	"github.com/AnTengye/NodeConverter/lib/network"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const maxLoggedShareURLLength = 160

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
		var (
			data []byte
			err  error
		)
		if !strings.HasPrefix(u, "http://api:15831") {
			// 从缓存中获取非内部api
			response, networkErr := network.CacheGET(u)
			if networkErr != nil {
				zap.S().Errorw("fetch subscription failed", "url", u, "stage", "download", zap.Error(networkErr))
				return nil, fmt.Errorf("download from sub error: %v", networkErr)
			}
			data = response
		} else {
			resp, networkErr := network.RestyCli.R().Get(u)
			if networkErr != nil {
				zap.S().Errorw("fetch subscription failed", "url", u, "stage", "download", zap.Error(networkErr))
				return nil, fmt.Errorf("download from sub error: %v", networkErr)
			}
			data = resp.Bytes()
		}
		nodes, err = handlerSubResponse(data)
		if err != nil {
			zap.S().Errorw("fetch subscription failed", "url", u, "stage", "parse", zap.Error(err))
			network.DeleteCache(u)
			return nil, fmt.Errorf("get sub error: %v", err)
		}
	case shareUrl:
		node, err := handlerSingleShareUrl(u)
		if err != nil {
			zap.S().Errorw("fetch share url failed", "share_url", limitForLog(u), zap.Error(err))
			return nil, fmt.Errorf("get share error: %v", err)
		}
		nodes = append(nodes, node)
	default:
		zap.S().Errorw("fetch nodes failed", "input", limitForLog(u), "stage", "validate", "error", "url is invalid")
		return nil, fmt.Errorf("url is invalid")
	}
	return nodes, nil
}

func handlerSubResponse(urlResponse []byte) ([]core.Node, error) {
	sDec, err := base64.StdEncoding.DecodeString(string(urlResponse))
	if err != nil {
		zap.S().Infow("subscription is not base64, fallback to clash decode", zap.Error(err))
		nodes, clashErr := handlerClashResponse(urlResponse)
		if clashErr != nil {
			zap.S().Errorw("subscription clash decode failed", zap.Error(clashErr))
			return nil, clashErr
		}
		return nodes, nil
	}
	// 按行读取
	lines := strings.Split(string(sDec), "\n")
	result := make([]core.Node, 0, len(lines))
	failed := 0
	for idx, line := range lines {
		if line == "" {
			continue
		}
		node, err := handlerSingleShareUrl(line)
		if err != nil {
			failed++
			zap.S().Errorw("subscription node conversion failed",
				"line", idx+1,
				"share_url", limitForLog(line),
				zap.Error(err),
			)
			continue
		}
		result = append(result, node)
	}
	if failed > 0 {
		zap.S().Warnw("subscription parsed with invalid nodes",
			"total_lines", len(lines),
			"success_nodes", len(result),
			"failed_nodes", failed,
		)
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
	case core.NodeTypeTUIC:
		node = core.NewTUICNode()
	case core.NodeTypeHysteria, core.NodeTypeHysteria2:
		node = core.NewHYSTERIANode()
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
			zap.S().Warnw("clash proxy format is invalid", "proxy", proxy)
			continue
		}

		// Get the proxy type
		proxyType, ok := proxyMap["type"].(string)
		if !ok {
			zap.S().Warnw("clash proxy type is invalid", "proxy_name", proxyMap["name"], zap.Any("proxy_type", proxyMap["type"]))
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
		case core.NodeTypeTUIC:
			node = core.NewTUICNode()
		case core.NodeTypeHysteria, core.NodeTypeHysteria2:
			node = core.NewHYSTERIANode()
		default:
			zap.S().Warnw("clash proxy type is not supported", "proxy_name", proxyMap["name"], "proxy_type", proxyType)
			continue
		}
		err := node.FromClash(proxyBytes)
		if err != nil {
			zap.S().Errorw("clash proxy conversion failed",
				"proxy_name", proxyMap["name"],
				"proxy_type", proxyType,
				zap.Error(err),
			)
			continue
		}
		result = append(result, node)
	}
	return result, nil
}

func limitForLog(value string) string {
	if len(value) <= maxLoggedShareURLLength {
		return value
	}
	return value[:maxLoggedShareURLLength] + "...(truncated)"
}
