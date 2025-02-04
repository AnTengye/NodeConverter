package core

import (
	"fmt"
	"regexp"

	"github.com/AnTengye/NodeConvertor/lib/yemoji"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const (
	ClashProxies      = "proxies"
	ClashProxiesGroup = "proxy-groups"
	ClashRules        = "rules"
)
const (
	ClashKernelClash     = "clash"
	ClashKernelClashMeta = "clashmeta"
)

type baseClash struct {
	ProxyGroups []ProxyGroup `yaml:"proxy-groups"`
	Rules       []string     `yaml:"rules"`
}
type ProxyGroup struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	// 其他字段
	OtherFields map[string]interface{} `yaml:",inline"`
	Proxies     []string               `yaml:"proxies,omitempty"`
}

type Clash struct {
	kernel  string
	data    map[string]any
	Proxies []Node
	base    *baseClash
}

func NewClash(kernel string) *Clash {
	return &Clash{
		kernel:  kernel,
		Proxies: make([]Node, 0),
		data:    make(map[string]any),
		base:    &baseClash{},
	}
}

func (c *Clash) WithTemplate(filePath string) error {
	var data map[string]any
	yamlData, err := yamlFromFile(filePath)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(yamlData, &data); err != nil {
		return err
	}
	var base baseClash
	if err = yaml.Unmarshal(yamlData, &base); err != nil {
		return err
	}
	c.data = data
	c.base = &base
	return nil
}

func (c *Clash) AddProxy(n ...Node) {
	c.Proxies = append(c.Proxies, n...)
}

func (c *Clash) SetACLSSR(acl *ClashACLSSR) {
	c.base.Rules = acl.RuleSet
	c.base.ProxyGroups = acl.ProxyGroups
}

func (c *Clash) ToYaml() (string, error) {
	c.data[ClashProxies] = c.Proxies
	c.withProxyGroup()
	c.data[ClashRules] = c.base.Rules
	d, err := yaml.Marshal(c.data)
	if err != nil {
		zap.S().Errorw("ToYaml error", "err", err)
		return "", err
	}
	unicodePoints, err := yemoji.ParseUnicodePoints(d)
	if err != nil {
		return "", fmt.Errorf("generate clash with unicode error: %v", err)
	}
	return string(unicodePoints), nil
}

func (c *Clash) withProxyGroup() {
	proxies := make([]string, len(c.Proxies))
	for i, v := range c.Proxies {
		proxies[i] = v.Name()
	}
	for i, v := range c.base.ProxyGroups {
		if len(v.Proxies) == 0 {
			if c.kernel == ClashKernelClash {
				if filter, ok := v.OtherFields["filter"]; ok {
					compile, err := regexp.Compile(filter.(string))
					if err != nil {
						zap.S().Errorw("regexp compile error", "err", err)
						continue
					}
					matchProxies := make([]string, 0, len(proxies))
					for _, proxy := range proxies {
						match := compile.MatchString(proxy)
						if match {
							matchProxies = append(matchProxies, proxy)
						}
					}
					c.base.ProxyGroups[i].Proxies = matchProxies
				}
				if len(c.base.ProxyGroups[i].Proxies) == 0 {
					c.base.ProxyGroups[i].Proxies = proxies
				}
				delete(v.OtherFields, "filter")
				delete(v.OtherFields, "include-all-proxies")
			}
		}
	}
	c.data[ClashProxiesGroup] = c.base.ProxyGroups
}
