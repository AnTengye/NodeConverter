package core

import (
	"encoding/base64"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"net/url"
	"strconv"
	"strings"
)

var _ Node = (*ShadowsocksNode)(nil)

//   - name: "ss1"
//     type: ss
//     server: server
//     port: 443
//     cipher: aes-128-gcm
//     password: "password"
//     udp: true
//     udp-over-tcp: false
//     udp-over-tcp-version: 2
//     ip-version: ipv4
//     plugin: obfs
//     plugin-opts:
//     mode: tls
//     smux:
//     enabled: false

// password¶
// Shadowsocks 密码
//
// udp-over-tcp¶
// 启用 UDP over TCP，默认 false
//
// udp-over-tcp-version¶
// UDP over TCP 的协议版本，默认 1。可选值 1/2。
//
// 插件¶
// plugin¶
// 插件，支持 obfs/v2ray-plugin/shadow-tls/restls
//
// plugin-opts¶
// 插件设置
type ShadowsocksNode struct {
	Normal            `yaml:",inline"`
	TLSConfig         `yaml:",inline"`
	NetworkConfig     `yaml:",inline"`
	Cipher            string         `json:"cipher" yaml:"cipher,omitempty"`                             // 加密,如：aes-128-ctr
	Password          string         `json:"password" yaml:"password,omitempty"`                         // 密码
	UdpOverTcp        bool           `json:"udp-over-tcp" yaml:"udp-over-tcp,omitempty"`                 // 启用 UDP over TCP，默认 false
	UdpOverTcpVersion int            `json:"udp-over-tcp-version" yaml:"udp-over-tcp-version,omitempty"` // UDP over TCP 的协议版本，默认 1。可选值 1/2。
	Plugin            string         `json:"plugin" yaml:"plugin,omitempty"`                             // 插件，支持 obfs/v2ray-plugin/shadow-tls/restls
	PluginOpts        map[string]any `json:"plugin-opts" yaml:"plugin-opts,omitempty"`                   // 插件设置
}

// SIP002 purposed a new URI scheme, following RFC3986:
//
// SS-URI = "ss://" userinfo "@" hostname ":" port [ "/" ] [ "?" plugin ] [ "#" tag ]
// userinfo = websafe-base64-encode-utf8(method  ":" password)
//
//	method ":" password
//
// Note that encoding userinfo with Base64URL is recommended but optional for Stream and AEAD (SIP004).
// But for AEAD-2022 (SIP022), userinfo MUST NOT be encoded with Base64URL. When userinfo is not encoded, method and password MUST be percent encoded.
//
// The last / should be appended if plugin is present, but is optional if only tag is present.
// Example: ss://YmYtY2ZiOnRlc3Q@192.168.100.1:8888/?plugin=url-encoded-plugin-argument-value&unsupported-arguments=should-be-ignored#Dummy+profile+name.
// This kind of URIs can be parsed by standard libraries provided by most languages.
//
// For plugin argument, we use the similar format as TOR_PT_SERVER_TRANSPORT_OPTIONS,
// which have the format like simple-obfs;obfs=http;obfs-host=example.com where colons, semicolons, equal signs and backslashes MUST be escaped with a backslash.
//
// Examples:
//
// With user info encoded with Base64URL:
//
// ss://YWVzLTEyOC1nY206dGVzdA@192.168.100.1:8888#Example1
// ss://cmM0LW1kNTpwYXNzd2Q@192.168.100.1:8888/?plugin=obfs-local%3Bobfs%3Dhttp#Example2
// Plain user info:
//
// ss://2022-blake3-aes-256-gcm:YctPZ6U7xPPcU%2Bgp3u%2B0tx%2FtRizJN9K8y%2BuKlW2qjlI%3D@192.168.100.1:8888#Example3
// ss://2022-blake3-aes-256-gcm:YctPZ6U7xPPcU%2Bgp3u%2B0tx%2FtRizJN9K8y%2BuKlW2qjlI%3D@192.168.100.1:8888/?plugin=v2ray-plugin%3Bserver#Example3
func (node *ShadowsocksNode) ToShare() string {
	builder := strings.Builder{}
	builder.WriteString("ss://")
	if node.Cipher != "" {
		builder.WriteString(node.Cipher)
		builder.WriteString(":")
	}
	builder.WriteString(node.Password)
	builder.WriteString("@")
	builder.WriteString(node.Server)
	builder.WriteString(":")
	builder.WriteString(strconv.Itoa(node.Port))
	builder.WriteString("#")
	builder.WriteString(node.Name())
	return builder.String()
}

func (node *ShadowsocksNode) FromShare(s string) error {
	parse, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("parse ss url err: %v", err)
	}
	if !strings.Contains(parse.Host, ".") {
		// base64解码parse.host之后，再重新url.Parse
		decodeString, err := base64.URLEncoding.DecodeString(parse.Host)
		if err != nil {
			return fmt.Errorf("shadowsocks decode error: %v", err)
		}
		parse, err = url.Parse(parse.Scheme + "://" + string(decodeString) + "#" + parse.Fragment)
		if err != nil {
			return fmt.Errorf("reparse ss url err: %v", err)
		}
	}
	setBase(parse, &node.Normal)
	values, err := url.ParseQuery(parse.RawQuery)
	if err != nil {
		return fmt.Errorf("parse trojan url[%s] err: %v", parse.RawQuery, err)
	}
	setNetwork(values, &node.NetworkConfig)
	setTLS(values, &node.TLSConfig)
	if parse.User != nil {
		node.Password, _ = parse.User.Password()
		node.Cipher = parse.User.Username()
	}
	if err := node.extra(values); err != nil {
		return err
	}
	if err := node.check(); err != nil {
		return err
	}
	return nil
}

func (node *ShadowsocksNode) extra(extra url.Values) error {
	return nil
}
func (node *ShadowsocksNode) check() error {
	return nil
}

func (node *ShadowsocksNode) ToClash() string {
	d, err := yaml.Marshal(&node)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return string(d)
}

func (node *ShadowsocksNode) FromClash(s []byte) error {
	if err := yaml.Unmarshal(s, node); err != nil {
		return fmt.Errorf("unmarshal Shadowsocks node error: %v", err)
	}
	return nil
}

func (node *ShadowsocksNode) Name() string {
	return node.Normal.Name
}

func (node *ShadowsocksNode) Type() NodeType {
	return node.Normal.Type
}

func NewShadowsocksNode() Node {
	return &ShadowsocksNode{}
}
