package core

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"gopkg.in/yaml.v3"
)

var _ Node = (*TuicNode)(nil)

type TuicNode struct {
	Normal                `yaml:",inline"`
	TLSConfig             `yaml:",inline"`
	Version               string `json:"version" yaml:"version,omitempty"`                                     // 用于指定 TUIC 的版本
	Token                 string `json:"token" yaml:"token,omitempty"`                                         // 必须，用于 TUIC V4 的用户标识，使用 TUIC V5 时不可书写
	Uuid                  string `json:"uuid" yaml:"uuid,omitempty"`                                           // 必须，用于 TUICV5 的用户唯一识别码，使用 TUIC V4 时不可书写
	Password              string `json:"password" yaml:"password,omitempty"`                                   // 必须，用于 TUICV5 的用户密码，使用 TUIC V4 时不可书写
	Ip                    string `json:"ip" yaml:"ip,omitempty"`                                               //用于覆盖“server”选项中设置的服务器地址的 DNS 查找结果
	HeartbeatInterval     int    `json:"heartbeat-interval" yaml:"heartbeat-interval,omitempty"`               // 发送保持连接活动的心跳包的间隔时间，单位为毫秒
	DisableSni            bool   `json:"disable-sni" yaml:"disable-sni,omitempty"`                             //  设置是否在 TLS 握手中禁用 SNI（服务器名称指示）
	ReduceRtt             bool   `json:"reduce-rtt" yaml:"reduce-rtt,omitempty"`                               // 设置是否在客户端启用 QUIC 的 0-RTT 握手
	RequestTimeout        int    `json:"request-timeout" yaml:"request-timeout,omitempty"`                     // 设置建立到 TUIC 代理服务器的连接的超时时间，单位为毫秒
	UdpRelayMode          string `json:"udp-relay-mode" yaml:"udp-relay-mode,omitempty"`                       // 设置 UDP 数据包中继模式，可以是 native/quic
	CongestionController  string `json:"congestion-controller" yaml:"congestion-controller,omitempty"`         // 设置拥塞控制算法，可选项为 cubic/new_reno/bbr
	MaxUdpRelayPacketSize int    `json:"max-udp-relay-packet-size" yaml:"max-udp-relay-packet-size,omitempty"` // 设置最大的 UDP 数据包中继大小，单位为字节
	FastOpen              bool   `json:"fast-open" yaml:"fast-open,omitempty"`                                 // 设置是否启用 Fast Open，这可以减少连接建立时间
	MaxOpenStreams        int    `json:"max-open-streams" yaml:"max-open-streams,omitempty"`                   // 设置最大打开流的数量
}

type TuicCustomParams struct {
}

// 4.1 基本信息段
// 4.1.1 协议名称 protocol
// 所使用的协议名称。取值必须为 vmess 或 tuic。
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
// 4.2.2 (VMess/TUIC) encryption
// 当协议为 VMess 时，对应配置文件出站中 settings.security，可选值为 auto / aes-128-gcm / chacha20-poly1305 / none。
//
// 省略时默认为 auto，但不可以为空字符串。除非指定为 none，否则建议省略。
//
// 当协议为 TUIC 时，对应配置文件出站中 settings.encryption，当前可选值只有 none。
//
// 省略时默认为 none，但不可以为空字符串。
//
// 特殊说明：之所以不使用 security 而使用 encryption，是因为后面还有一个底层传输安全类型 security 与这个冲突。
// 由 @huyz 提议，将此字段重命名为 encryption，这样不仅能避免命名冲突，还与 TUIC 保持了一致。
//
// 4.2.3 (VMess) alterId、aid 等
// 没有这些字段。旧的 VMess 因协议设计出现致命问题，不再适合使用或分享。
//
// 此分享标准仅针对 VMess AEAD 和 TUIC。
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
// 通用格式(TUIC+reality+uTLS+Vision)
func (node *TuicNode) ToShare() string {
	return node.buildBaseShareURI(
		string(node.Type()),
		func(builder *strings.Builder) {
			if node.Uuid != "" {
				fmt.Fprintf(builder, "%s:%s@", node.Uuid, node.Password)
			}
		},
		func(builder *strings.Builder) {
			builder.WriteString("?encryption=none")
			builder.WriteString("&")
			builder.WriteString(node.TlSToValues().Encode())
			if node.UdpRelayMode != "" {
				builder.WriteString("&udp_relay_mode=")
				builder.WriteString(node.UdpRelayMode)
			}
			if node.Version != "" {
				builder.WriteString("&version=")
				builder.WriteString(node.Version)
			}
		},
	)
}

func (node *TuicNode) FromShare(s string) error {
	node.ClientFingerprint = "chrome"
	node.DisableSni = true
	node.ReduceRtt = true
	node.RequestTimeout = 8000
	node.UdpRelayMode = "native"

	parse, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("parse tuic url err: %v", err)
	}
	values := parse.Query()
	setBase(parse, &node.Normal)
	setTLS(values, &node.TLSConfig)
	if parse.User != nil {
		node.Uuid = parse.User.Username()
		node.Password, _ = parse.User.Password()
	}
	node.convertValues(values)
	if err := node.check(); err != nil {
		return err
	}
	return nil
}

func (node *TuicNode) convertValues(values url.Values) {
	for k, v := range values {
		switch k {
		case "disable_sni":
			if v[0] == "true" {
				node.DisableSni = true
			} else {
				node.DisableSni = false
			}
		case "reduce_rtt":
			if v[0] == "true" {
				node.ReduceRtt = true
			} else {
				node.ReduceRtt = false
			}

		}
	}
}

func (node *TuicNode) check() error {
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
		node.ServerName = node.Server
	}
	return nil
}

func (node *TuicNode) ToClash() string {
	d, err := yaml.Marshal(&node)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return string(d)
}

func (node *TuicNode) FromClash(s []byte) error {
	if err := yaml.Unmarshal(s, node); err != nil {
		return fmt.Errorf("unmarshal tuic node error: %v", err)
	}
	return nil
}

func (node *TuicNode) Name() string {
	return node.Normal.Name
}

func (node *TuicNode) Type() NodeType {
	return node.Normal.Type
}

func NewTUICNode() Node {
	return &TuicNode{}
}
