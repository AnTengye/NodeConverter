package core

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

var _ Node = (*HysteriaNode)(nil)

//ports¶
//配置则启用端口跳跃，忽略port，格式参考端口范围
//
//password¶
//认证密码
//
//up/down¶
//brutal 速率控制，若不写单位，默认为 Mbps
//
//obfs¶
//QUIC 流量混淆器类型，仅可设为 salamander，如果为空则禁用
//
//obfs-password¶
//QUIC 流量混淆器密码

type HysteriaNode struct {
	Normal       `yaml:",inline"`
	TLSConfig    `yaml:",inline"`
	version      int    // hysteria版本
	Ports        string `json:"ports" yaml:"ports,omitempty"`                 // 配置则启用端口跳跃，忽略port，格式参考端口范围
	Password     string `json:"password" yaml:"password,omitempty"`           // 认证密码
	Up           int    `json:"up" yaml:"up,omitempty"`                       // brutal 速率控制，若不写单位，则默认为 Mbps
	Down         int    `json:"down" yaml:"down,omitempty"`                   // brutal 速率控制，若不写单位，则默认为 Mbps
	Obfs         string `json:"obfs" yaml:"obfs,omitempty"`                   // QUIC 流量混淆器类型，仅可设为 salamander，如果为空则禁用
	ObfsPassword string `json:"obfs-password" yaml:"obfs-password,omitempty"` // QUIC 流量混淆器密码

	// 以下为h1用
	AuthStr  string `json:"auth-str" yaml:"auth-str,omitempty"` // yourpassword
	Protocol string `json:"protocol" yaml:"protocol,omitempty"` // 支持 udp/wechat-video/faketcp
}

// 4.1 基本信息段
// 4.1.1 协议名称 protocol
// 所使用的协议名称。取值必须为 vmess 或 hysteria。
//
// 不可省略，不能为空字符串。
//
// 4.1.2 uuid
// UUID。对应配置文件该项出站中 settings.vnext[0].users[0].id 的值。
//
// 不可省略，不能为空字符串。
//
// 4.1.3 remote-host
// 服务器的域名或 IP 地址。
//
// 不可省略，不能为空字符串。
//
// IPv6 地址必须括上方括号。
//
// IDN 域名（如“百度.cn”）必须使用 xn--xxxxxx 格式。
//
// 4.1.4 remote-port
// 服务器的端口号。
//
// 不可省略，必须取 [1,65535] 中的整数。
//
// 4.1.5 descriptive-text
// 服务器的描述信息。
//
// 可省略，不推荐为空字符串。
//
// 必须使用 encodeURIComponent 转义。
//
// 4.2 协议相关段
// 4.2.1 传输方式 type （@RPRX 修改于 2024-03-05，2024-06-18，2024-11-11）
// 协议的传输方式。对应配置文件出站中 settings.vnext[0].streamSettings.network 的值。
//
// 当前的取值必须为 tcp、kcp、ws、http、grpc、httpupgrade、xhttp 其中之一，
// 分别对应 RAW、mKCP、WebSocket、HTTP/2/3、gRPC、HTTPUpgrade、XHTTP 传输方式。
//
// 4.2.2 (VMess/HYSTERIA) encryption
// 当协议为 VMess 时，对应配置文件出站中 settings.security，可选值为 auto / aes-128-gcm / chacha20-poly1305 / none。
//
// 省略时默认为 auto，但不可以为空字符串。除非指定为 none，否则建议省略。
//
// 当协议为 HYSTERIA 时，对应配置文件出站中 settings.encryption，当前可选值只有 none。
//
// 省略时默认为 none，但不可以为空字符串。
//
// 特殊说明：之所以不使用 security 而使用 encryption，是因为后面还有一个底层传输安全类型 security 与这个冲突。
// 由 @huyz 提议，将此字段重命名为 encryption，这样不仅能避免命名冲突，还与 HYSTERIA 保持了一致。
//
// 4.2.3 (VMess) alterId、aid 等
// 没有这些字段。旧的 VMess 因协议设计出现致命问题，不再适合使用或分享。
//
// 此分享标准仅针对 VMess AEAD 和 HYSTERIA。
//
// 4.3 传输层相关段
// 4.3.1 底层传输安全 security （@RPRX 修改于 2023-03-19）
// 设定底层传输所使用的 TLS 类型。当前可选值有 none，tls 和 reality。
//
// 省略时默认为 none，但不可以为空字符串。
//
// 4.3.2 (HTTP/2/3) path （@RPRX 修改于 2024-11-11）
// HTTP/2/3 的路径。省略时默认为 /，但不可以为空字符串。不推荐省略。
//
// 必须使用 encodeURIComponent 转义。
//
// 4.3.3 (HTTP/2/3) host （@RPRX 修改于 2024-11-11）
// 客户端进行 HTTP/2/3 通信时所发送的 Host 头部。
//
// 省略时复用 remote-host，但不可以为空字符串。
//
// 若有多个域名，可使用英文逗号隔开，但中间及前后不可有空格。
//
// 必须使用 encodeURIComponent 转义。
//
// 4.3.4 (WebSocket) path
// WebSocket 的路径。省略时默认为 /，但不可以为空字符串。不推荐省略。
//
// 必须使用 encodeURIComponent 转义。
//
// 4.3.5 (WebSocket) host
// WebSocket 请求时 Host 头的内容。不推荐省略，不推荐设为空字符串。
//
// 必须使用 encodeURIComponent 转义。
//
// 4.3.6 (mKCP) headerType
// mKCP 的伪装头部类型。当前可选值有 none / srtp / utp / wechat-video / dtls / wireguard。
//
// 省略时默认值为 none，即不使用伪装头部，但不可以为空字符串。
//
// 4.3.7 (mKCP) seed
// mKCP 种子。省略时不使用种子，但不可以为空字符串。建议 mKCP 用户使用 seed。
//
// 必须使用 encodeURIComponent 转义。
//
// 4.3.11 (gRPC) serviceName （@RPRX 修改于 2024-03-05）
// 对应 gRPC 的 ServiceName。建议仅使用英文字母数字和英文句号、下划线组成。
//
// 不建议省略，不可为空字符串。
//
// 修订：必须使用 encodeURIComponent 转义。#1815
//
// 4.3.12 (gRPC) mode
// 对应 gRPC 的传输模式，目前有以下三种：
//
// gun: 即原始的 gun 传输模式，将单个 []byte 封在 Protobuf 里通过 gRPC 发送（参考资料）；
// multi: 即 Xray-Core 的 multiMode，将多组 []byte 封在一条 Protobuf 里通过 gRPC 发送；
// guna: 即通过使用自定义 Codec 的方式，直接将数据包封在 gRPC 里发送。（参考资料）
// 省略时默认为 gun，不可以为空字符串。
//
// 4.3.13 (gRPC) authority （@RPRX 添加于 2024-03-05）
// 对应 gRPC 的 Authority。#3076
//
// 此项可能为空字符串。
//
// 必须使用 encodeURIComponent 转义。
//
// 4.3.14 (HTTPUpgrade) path （@RPRX 添加于 2024-03-05）
// HTTPUpgrade 的路径。省略时默认为 /，但不可以为空字符串。不推荐省略。
//
// 必须使用 encodeURIComponent 转义。
//
// 4.3.15 (HTTPUpgrade) host （@RPRX 添加于 2024-03-05）
// HTTPUpgrade 请求时 Host 头的内容。不推荐省略，不推荐设为空字符串。
//
// 必须使用 encodeURIComponent 转义。
//
// 4.3.16 (XHTTP) path （@RPRX 添加于 2024-06-18，修改于 2024-11-11）
// XHTTP 的路径。省略时默认为 /，但不可以为空字符串。不推荐省略。
//
// 必须使用 encodeURIComponent 转义。
//
// 4.3.17 (XHTTP) host （@RPRX 添加于 2024-06-18，修改于 2024-11-11）
// XHTTP 请求时 Host 头的内容。不推荐省略，不推荐设为空字符串。
//
// 必须使用 encodeURIComponent 转义。
//
// 4.3.18 (XHTTP) mode （@RPRX 添加于 2024-11-11）
// XHTTP 的 mode：#3994
//
// 4.3.19 (XHTTP) extra （@RPRX 添加于 2024-11-11）
// XHTTP 的 extra：#4000
//
// 必须使用 encodeURIComponent 转义。
//
// 4.4 TLS 相关段
// 4.4.0 fp （@RPRX 添加于 2023-02-01，修改于 2023-03-19）
// TLS Client Hello 指纹，对应配置文件中的 fingerprint 项目。
//
// 省略时默认为 chrome，不可以为空字符串。
//
// 若使用 REALITY，此项不可省略。
//
// Q: 为什么该项在分享链接中为 fp 而不是 fingerprint？
// A: 类似 sni、alpn，尽量缩短分享链接长度。
//
// Q: 为什么省略时默认为 chrome？
// A: Golang TLS Client Hello 指纹已被针对，而 Chrome 是目前市占率最高的浏览器。
//
// Q: 为什么是 chrome 而不是 random？
// A: 一是避免指纹总比例接近 Xray 预置比例而暴露统计学特征，二是避免 uTLS 对其中的某个指纹实现不当而暴露“一票否决”的特征。
//
// 4.4.1 sni
// TLS SNI，对应配置文件中的 serverName 项目。
//
// 省略时复用 remote-host，但不可以为空字符串。
//
// 4.4.2 alpn
// TLS ALPN，对应配置文件中的 alpn 项目。
//
// 多个 ALPN 之间用英文逗号隔开，中间无空格。
//
// 省略时由内核决定具体行为，但不可以为空字符串。
//
// 必须使用 encodeURIComponent 转义。
//
// 4.4.3 allowInsecure
// 没有这个字段。不安全的节点，不适合分享。
//
// 4.4.4 (XTLS) flow （@RPRX 修改于 2024-11-11）
// XTLS 的流控方式。可选值为 xtls-rprx-vision 等。
//
// 若使用 XTLS，此项不可省略，否则无此项。此项不可为空字符串。
//
// 4.4.5 (REALITY) pbk （@RPRX 添加于 2023-03-19）
// REALITY 的公钥，对应配置文件中的 publicKey 项目。
//
// 若使用 REALITY，此项不可省略，否则无此项。此项不可为空字符串。
//
// 4.4.6 (REALITY) sid （@RPRX 添加于 2023-03-19）
// REALITY 的 ID，对应配置文件中的 shortId 项目。
//
// 无需特殊处理。此项可能为空字符串。
//
// 4.4.7 (REALITY) spx （@RPRX 添加于 2023-03-19）
// REALITY 的爬虫，对应配置文件中的 spiderX 项目。
//
// 必须使用 encodeURIComponent 转义。此项可能为空字符串。
// 通用格式(HYSTERIA+reality+uTLS+Vision)
func (node *HysteriaNode) ToShare() string {
	return node.buildBaseShareURI(
		string(node.Type()),
		func(builder *strings.Builder) {
			builder.WriteString(node.Password)
			builder.WriteString("@")
		},
		func(builder *strings.Builder) {
			builder.WriteString("?encryption=none")
			tlsValues := node.TlSToValues()
			builder.WriteString("&")
			builder.WriteString(tlsValues.Encode())
			if node.AuthStr != "" {
				builder.WriteString("&auth_str=")
				builder.WriteString(node.AuthStr)
			}
			if node.Down != 0 {
				builder.WriteString("&downmbps=")
				builder.WriteString(strconv.Itoa(node.Down))
			}
			if node.Up != 0 {
				builder.WriteString("&upmbps=")
				builder.WriteString(strconv.Itoa(node.Up))
			}
			if node.Protocol != "" {
				builder.WriteString("&protocol=")
				builder.WriteString(node.Protocol)
			}
			if node.Obfs != "" {
				builder.WriteString("&obfs=")
				builder.WriteString(node.Obfs)
			}
			if node.ObfsPassword != "" {
				builder.WriteString("&obfs_password=")
				builder.WriteString(node.ObfsPassword)
			}
		},
	)
}

// FromShare
// hysteria://1.1.1.1:14241?alpn=h3&auth=dongtaiwang.com&auth_str=dongtaiwang.com&delay=1004&downmbps=100&protocol=udp&insecure=1&peer=apple.com&udp=true&upmbps=100#%F0%9F%87%AB%F0%9F%87%B7%E6%B3%95%E5%9B%BD1-%20%E2%AC%87%EF%B8%8F%209.4MB%2Fs
func (node *HysteriaNode) FromShare(s string) error {
	node.ClientFingerprint = "chrome"
	node.Down = 1000
	node.Up = 1000

	parse, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("parse hysteria url err: %v", err)
	}
	values := parse.Query()
	setBase(parse, &node.Normal)
	setTLS(values, &node.TLSConfig)
	if parse.User != nil {
		node.Password = parse.User.Username()
	}
	node.convertValues(values)
	if err := node.check(); err != nil {
		return err
	}
	return nil
}

func (node *HysteriaNode) convertValues(values url.Values) {
	for k, v := range values {
		switch k {
		case "auth_str":
			node.AuthStr = v[0]
		case "downmbps":
			node.Down = mbpsConvert(v[0])
		case "upmbps":
			node.Up = mbpsConvert(v[0])
		case "protocol":
			node.Protocol = v[0]
		case "obfs":
			node.Obfs = v[0]
		case "obfs-password":
			node.ObfsPassword = v[0]
		}
	}
}

// 读取字符串，将1000 Mbps 转成 1000
func mbpsConvert(s string) int {
	atoi, err := strconv.Atoi(s)
	if err != nil {
		split := strings.Split(s, " ")
		if len(split) == 2 {
			v := 1
			switch split[1] {
			case "Mbps":
				v = 1 // 1000 Mbps
			case "Kbps":
				v = 1000
			}
			atoi, _ = strconv.Atoi(split[0])
			return atoi * v
		}
	}
	return atoi
}

func (node *HysteriaNode) check() error {
	node.Fingerprint = ""
	if node.RealityOpts != nil {
		if node.ClientFingerprint == "" {
			return fmt.Errorf("reality-opts need client-fingerprint")
		}
		if node.RealityOpts.PublicKey == "" {
			return fmt.Errorf("reality-opts need public-key")
		}
	}
	if node.ServerName == "" {
		node.SNI = node.Server
		node.ServerName = node.Server
	}
	if node.Up == 0 || node.Down == 0 {
		return fmt.Errorf("need up or down")
	}
	if node.Type() == NodeTypeHysteria {
		if node.AuthStr == "" {
			return fmt.Errorf("need auth_str")
		}
	} else {
		if node.Password == "" {
			return fmt.Errorf("need password")
		}
		if node.ObfsPassword == "" && node.Obfs != "" {
			return fmt.Errorf("need obfs-password")
		}
	}
	return nil
}

func (node *HysteriaNode) ToClash() string {
	d, err := yaml.Marshal(&node)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return string(d)
}

func (node *HysteriaNode) FromClash(s []byte) error {
	if err := yaml.Unmarshal(s, node); err != nil {
		return fmt.Errorf("unmarshal hysteria node error: %v", err)
	}
	return nil
}

func (node *HysteriaNode) Name() string {
	return node.Normal.Name
}

func (node *HysteriaNode) Type() NodeType {
	return node.Normal.Type
}

func NewHYSTERIANode() Node {
	return &HysteriaNode{}
}
