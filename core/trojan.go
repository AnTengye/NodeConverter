package core

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"net/url"
	"strconv"
	"strings"
)

var _ Node = (*TrojanNode)(nil)

type TrojanNode struct {
	Normal        `yaml:",inline"`
	TLSConfig     `yaml:",inline"`
	NetworkConfig `yaml:",inline"`
	Password      string          `yaml:"password,omitempty"` // trojan 服务器密码
	SSOpts        TrojanSSOptions `yaml:"ss-opts,omitempty"`
}

func (node *TrojanNode) Name() string {
	return node.Normal.Name
}

func (node *TrojanNode) Type() NodeType {
	return node.Normal.Type
}

type TrojanSSOptions struct {
	Enabled  bool   `yaml:"enabled,omitempty"`  // 启用 trojan-go 的 ss AEAD 加密
	Method   string `yaml:"method,omitempty"`   // 加密方法，支持 aes-128-gcm/aes-256-gcm/chacha20-ietf-poly1305
	Password string `yaml:"password,omitempty"` // ss AEAD 加密密码
}

// trojan-password
// Trojan 的密码。 不可省略，不能为空字符串，不建议含有非 ASCII 可打印字符。 必须使用 encodeURIComponent 编码。
//
// trojan-host
// 节点 IP / 域名。 不可省略，不能为空字符串。 IPv6 地址必须扩方括号。 IDN 域名（如“百度.cn”）必须使用 xn--xxxxxx 格式。
//
// port
// 节点端口。 省略时默认为 443。 必须取 [1,65535] 中的整数。
//
// tls或allowInsecure
// 没有这个字段。 TLS 默认一直启用，除非有传输插件禁用它。 TLS 认证必须开启。无法使用根CA校验服务器身份的节点，不适合分享。
//
// sni
// 自定义 TLS 的 SNI。 省略时默认与 trojan-host 同值。不得为空字符串。
//
// 必须使用 encodeURIComponent 编码。
//
// type
//
//	传输类型。 省略时默认为 original，但不可为空字符串。 目前可选值只有 original 和 ws，未来可能会有 h2、h2+ws 等取值。
//
// 当取值为 original 时，使用原始 Trojan 传输方式，无法方便通过 CDN。 当取值为 ws 时，使用 Websocket over TLS 传输。
//
// host
// 自定义 HTTP Host 头。 可以省略，省略时值同 trojan-host。 可以为空字符串，但可能带来非预期情形。
//
// 警告：若你的端口非标准端口（不是 80 / 443），RFC 标准规定 Host 应在主机名后附上端口号，例如 example.com:44333。至于是否遵守，请自行斟酌。
//
// 必须使用 encodeURIComponent 编码。
//
// path
// 当传输类型 type 取 ws、h2、h2+ws 时，此项有效。 不可省略，不可为空。 必须以 / 开头。 可以使用 URL 中的 & # ? 等字符，但应当是合法的 URL 路径。
//
// 必须使用 encodeURIComponent 编码。
//
// mux
// 没有这个字段。 当前服务器默认一直支持 mux。 启用 mux 与否各有利弊，应由客户端决定自己是否启用。URL的作用，是定位服务器资源，而不是规定用户使用偏好。
//
// encryption
// 用于保证 Trojan 流量密码学安全的加密层。 可省略，默认为 none，即不使用加密。 不可以为空字符串。
//
// 必须使用 encodeURIComponent 编码。
//
// 使用 Shadowsocks 算法进行流量加密时，其格式为：
//
// ss;method:password
// 其中 ss 是固定内容，method 是加密方法，必须为下列之一：
//
// aes-128-gcm
// aes-256-gcm
// chacha20-ietf-poly1305
// 其中的 password 是 Shadowsocks 的密码，不得为空字符串。 password 中若包含分号，不需要进行转义。 password 应为英文可打印 ASCII 字符。
//
// 其他加密方案待定。
//
// plugin
// 额外的插件选项。本字段保留。 可省略，但不可以为空字符串。
//
// URL Fragment (# 后内容)
// 节点说明。 不建议省略，不建议为空字符串。
//
// 必须使用 encodeURIComponent 编码。
func (node *TrojanNode) ToShare() string {
	builder := strings.Builder{}
	builder.WriteString("trojan://")
	builder.WriteString(node.Password)
	builder.WriteString("@")
	builder.WriteString(node.Server)
	builder.WriteString(":")
	builder.WriteString(strconv.Itoa(node.Port))
	builder.WriteString("?encryption=none")
	if node.Network != "" {
		builder.WriteString("&type=")
		builder.WriteString(node.Network)
	}
	if node.TLSConfig.SkipCertVerify {
		builder.WriteString("&allowInsecure=1")
	}
	if node.TLSConfig.TLS {
		builder.WriteString("&tls=1")
	}
	if node.TLSConfig.SNI != "" {
		builder.WriteString("&sni=")
		builder.WriteString(node.TLSConfig.SNI)
	}
	builder.WriteString("#")
	builder.WriteString(node.Name())
	return builder.String()
}

func (node *TrojanNode) FromShare(s string) error {
	parse, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("parse trojan url err: %v", err)
	}
	setBase(parse, &node.Normal)
	values := parse.Query()
	setNetwork(values, &node.NetworkConfig)
	setTLS(values, &node.TLSConfig)

	if parse.User != nil {
		node.Password = parse.User.Username()
	}
	if err := node.extra(values); err != nil {
		return err
	}
	if err := node.check(); err != nil {
		return err
	}
	return nil
}

func (node *TrojanNode) extra(extra url.Values) error {
	return nil
}
func (node *TrojanNode) check() error {
	return nil
}

func (node *TrojanNode) ToClash() string {
	d, err := yaml.Marshal(&node)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return string(d)
}

func (node *TrojanNode) FromClash(s []byte) error {
	if err := yaml.Unmarshal(s, node); err != nil {
		return fmt.Errorf("unmarshal trojan node error: %v", err)
	}
	return nil
}

func NewTrojanNode() Node {
	return &TrojanNode{}
}
