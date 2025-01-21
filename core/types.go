package core

type NodeType string

const (
	NodeTypeShadowSocks NodeType = "ss"
	NodeTypeVMess       NodeType = "vmess"
	NodeTypeTrojan      NodeType = "trojan"
	NodeTypeVLESS       NodeType = "vless"
)

// 参考文档：
// https://wiki.metacubex.one/

type Normal struct {
	Name string   `json:"name" yaml:"name"` // 必须，代理名称，不可重复
	Type NodeType `json:"type" yaml:"type"` // 必须，代理节点类型

	Server string `json:"server" yaml:"server"`
	Port   int    `json:"port" yaml:"port"`

	IPVersion string `json:"ip-version" yaml:"ip-version,omitempty"`
	UDP       bool   `json:"udp" yaml:"udp,omitempty"`

	InterfaceName string `json:"interface-name" yaml:"interface-name,omitempty"`
	RoutingMark   int    `json:"routing-mark" yaml:"routing-mark,omitempty"`
	TFO           bool   `json:"tfo" yaml:"tfo,omitempty"`
	MPTCP         bool   `json:"mptcp" yaml:"mptcp,omitempty"`

	DialerProxy string `json:"dialer-proxy" yaml:"dialer-proxy,omitempty"`

	Smux Smux `json:"smux" yaml:"smux,omitempty"`
}

type Smux struct {
	Enabled        bool   `json:"enabled" yaml:"enabled,omitempty"`
	Protocol       string `json:"protocol" yaml:"protocol,omitempty"`
	MaxConnections int    `json:"max-connections" yaml:"max-connections,omitempty"`
	MinStreams     int    `json:"min-streams" yaml:"min-streams,omitempty"`
	MaxStreams     int    `json:"max-streams" yaml:"max-streams,omitempty"`
	Statistic      bool   `json:"statistic" yaml:"statistic,omitempty"`
	OnlyTCP        bool   `json:"only-tcp" yaml:"only-tcp,omitempty"`
	Padding        bool   `json:"padding" yaml:"padding,omitempty"`
	BrutalOpts     struct {
		Enabled bool `json:"enabled" yaml:"enabled,omitempty"`
		Up      int  `json:"up" yaml:"up,omitempty"`
		Down    int  `json:"down" yaml:"down,omitempty"`
	} `json:"brutal-opts" yaml:"brutal-opts,omitempty"`
}

type TLSConfig struct {
	TLS               bool              `json:"tls" yaml:"tls,omitempty"`
	SNI               string            `json:"sni" yaml:"sni,omitempty"`
	ServerName        string            `json:"servername" yaml:"servername,omitempty"`
	Fingerprint       string            `json:"fingerprint" yaml:"fingerprint,omitempty"`
	ALPN              []string          `json:"alpn" yaml:"alpn,omitempty"`
	SkipCertVerify    bool              `json:"skip-cert-verify" yaml:"skip-cert-verify,omitempty"`
	ClientFingerprint string            `json:"client-fingerprint" yaml:"client-fingerprint,omitempty"`
	RealityOpts       *RealityTlsConfig `json:"reality-opts" yaml:"reality-opts,omitempty"`
}

type RealityTlsConfig struct {
	PublicKey string `json:"public-key" yaml:"public-key,omitempty"`
	ShortID   string `json:"short-id" yaml:"short-id,omitempty"`
}

// network: http
// http-opts:
//
//	method: "GET"
//	path:
//	- '/'
//	- '/video'
//	headers:
//	  Connection:
//	  - keep-alive
type NetworkConfig struct {
	Network string `json:"network" yaml:"network,omitempty"`

	HTTPOpts *HTTPNetworkConfig `json:"http-opts" yaml:"http-opts,omitempty"`
	H2Opts   *H2NetworkConfig   `json:"h2-opts" yaml:"h2-opts,omitempty"`
	GRPCOpts *GRPCNetworkConfig `json:"grpc-opts" yaml:"grpc-opts,omitempty"`
	WSOpts   *WSNetworkConfig   `json:"ws-opts" yaml:"ws-opts,omitempty"`
}

type HTTPNetworkConfig struct {
	Method  string              `json:"method" yaml:"method,omitempty"`
	Path    []string            `json:"path" yaml:"path,omitempty"`
	Headers map[string][]string `json:"headers" yaml:"headers,omitempty"`
}

// h2-opts:
//
//	host:
//	- example.com
//	path: /
type H2NetworkConfig struct {
	Host []string `json:"host" yaml:"host,omitempty"`
	Path string   `json:"path" yaml:"path,omitempty"`
}

type GRPCNetworkConfig struct {
	GRPCServiceName string `json:"grpc-service-name" yaml:"grpc-service-name,omitempty"`
}

// ws-opts:
//
//	path: /path
//	headers:
//	  Host: example.com
//	max-early-data:
//	early-data-header-name:
//	v2ray-http-upgrade: false
//	v2ray-http-upgrade-fast-open: false
type WSNetworkConfig struct {
	Path                     string              `json:"path" yaml:"path,omitempty"`
	Headers                  map[string][]string `json:"headers" yaml:"headers,omitempty"`
	MaxEarlyData             int                 `json:"max-early-data" yaml:"max-early-data,omitempty"`
	EarlyDataHeaderName      string              `json:"early-data-header-name" yaml:"early-data-header-name,omitempty"`
	V2RayHTTPUpgrade         bool                `json:"v2ray-http-upgrade" yaml:"v2ray-http-upgrade,omitempty"`
	V2RayHTTPUpgradeFastOpen bool                `json:"v2ray-http-upgrade-fast-open" yaml:"v2ray-http-upgrade-fast-open,omitempty"`
}
