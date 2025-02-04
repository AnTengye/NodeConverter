package core

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/AnTengye/NodeConvertor/lib/network"
	"go.uber.org/zap"
)

type ClashACLSSR struct {
	RuleSet     []string
	ProxyGroups []ProxyGroup
}

func NewClashACLSSRFromBytes(data []byte) (*ClashACLSSR, error) {
	zap.S().Infow("loading ACLSSR config")
	var (
		rs         []string
		pg         []ProxyGroup
		ruleSet    []string
		ruleSetMap = make(map[string][]byte)
	)
	group := sync.WaitGroup{}
	mutex := sync.Mutex{}
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ruleset=") {
			split := strings.SplitN(line[8:], ",", 2)
			if len(split) < 2 {
				return nil, fmt.Errorf("ruleset[%s] is invalid", line)
			}
			groupName := split[0]
			ruleSet = append(ruleSet, line[8:])
			if strings.HasPrefix(split[1], "[]") {
				//ruleset=🎯 全球直连,[]GEOIP,CN
				//ruleset=🐟 漏网之鱼,[]FINAL
				continue
			}
			group.Add(1)
			go func(groupName string, url string) {
				defer group.Done()
				resp, err := network.CacheGET(url)
				if err != nil {
					zap.S().Warnw("fetch ruleset error", zap.String("url", url), zap.Error(err))
					return
				}
				mutex.Lock()
				ruleSetMap[url] = resp
				mutex.Unlock()
			}(groupName, split[1])
		} else if strings.HasPrefix(line, "custom_proxy_group=") {
			proxyGroup, err := handlerProxyGroup(line[19:])
			if err != nil {
				return nil, err
			}
			pg = append(pg, proxyGroup)
		}
	}
	zap.S().Infow("wait for ruleset download finish")
	group.Wait()
	zap.S().Infow("ruleset download finish")

	for _, v := range ruleSet {
		split := strings.SplitN(v, ",", 2)
		groupName := split[0]
		url := split[1]
		if strings.HasPrefix(url, "[]") {
			//ruleset=🎯 全球直连,[]GEOIP,CN
			//ruleset=🐟 漏网之鱼,[]FINAL
			if url == "[]FINAL" {
				rs = append(rs, "MATCH,"+groupName)
			} else {
				rs = append(rs, split[1][2:]+","+groupName)
			}
			continue
		}
		if resp, ok := ruleSetMap[url]; ok {
			items, err := handlerWithRuleItem(groupName, resp)
			if err != nil {
				zap.S().Warnw("parse ruleset error", zap.String("url", url), zap.Error(err))
				continue
			}
			rs = append(rs, items...)
		}
	}
	return &ClashACLSSR{RuleSet: rs, ProxyGroups: pg}, nil
}

func handlerWithRuleItem(key string, data []byte) ([]string, error) {
	result := make([]string, 0, 20)
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "URL-REGEX") || strings.HasPrefix(line, "USER-AGENT") {
			zap.S().Debugw("skip line", "line", line)
			continue
		}
		if strings.HasSuffix(line, ",no-resolve") {
			//IP-CIDR,0.0.0.0/8,no-resolve -> IP-CIDR,0.0.0.0/8,key,no-resolve
			result = append(result, strings.Replace(line, ",no-resolve", ","+key+",no-resolve", 1))
		} else {
			result = append(result, line+","+key)
		}
	}
	return result, nil
}

func handlerProxyGroup(data string) (ProxyGroup, error) {
	pg := ProxyGroup{
		OtherFields: make(map[string]interface{}),
	}
	split := strings.Split(data, "`")
	if len(split) < 2 {
		return pg, fmt.Errorf("proxy group is invalid")
	}
	pg.Name = split[0]
	pg.Type = split[1]
	switch pg.Type {
	case "url-test":
		if len(split) != 5 {
			return pg, fmt.Errorf("url-test is invalid")
		}
		params := strings.Split(split[4], ",")
		if len(params) != 3 {
			return pg, fmt.Errorf("url-test need 3 params")
		}
		pg.OtherFields["include-all-proxies"] = true
		if split[2] != ".*" {
			pg.OtherFields["filter"] = split[2]
		}
		interval, parseErr := strconv.ParseInt(params[0], 10, 64)
		if parseErr != nil {
			interval = 300
		}
		tolerance, parseErr := strconv.ParseInt(params[2], 10, 64)
		if parseErr != nil {
			tolerance = 50
		}
		pg.OtherFields["url"] = split[3]
		pg.OtherFields["interval"] = interval
		pg.OtherFields["tolerance"] = tolerance
	case "select":
		//奈飞节点`select`(NF|奈飞|解锁|Netflix|NETFLIX|Media)
		//节点选择`select`[]♻️ 自动选择`[]🇭🇰 香港节点`[]DIRECT
		//手动切换`select`.*
		for _, v := range split[2:] {
			if v[:2] == "[]" {
				pg.Proxies = append(pg.Proxies, v[2:])
			} else {
				pg.OtherFields["include-all-proxies"] = true
				pg.OtherFields["filter"] = v
				break
			}
		}
	}
	return pg, nil
}
