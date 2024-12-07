package main

type Node interface {
	ToShare() string
	ToClash() string
	FromShare(string) error
	FromClash(string) error
}

type Normal struct {
	Name string `json:"name" yaml:"name"` // 必须，代理名称，不可重复
	Type string `json:"type" yaml:"type"` // 必须，代理节点类型

	Server string `json:"server" yaml:"server"`
	Port   int    `json:"port" yaml:"port"`

	IPVersion string `json:"ip-version" yaml:"ip-version"`
	UDP       bool   `json:"udp" yaml:"udp"`

	InterfaceName string `json:"interface-name" yaml:"interface-name"`
	RoutingMark   int    `json:"routing-mark" yaml:"routing-mark"`
	TFO           bool   `json:"tfo" yaml:"tfo"`
	MPTCP         bool   `json:"mptcp" yaml:"mptcp"`

	DialerProxy string `json:"dialer-proxy" yaml:"dialer-proxy"`

	Smux Smux `json:"smux" yaml:"smux"`
}

type Smux struct {
	Enabled        bool   `json:"enabled" yaml:"enabled"`
	Protocol       string `json:"protocol" yaml:"protocol"`
	MaxConnections int    `json:"max-connections" yaml:"max-connections"`
	MinStreams     int    `json:"min-streams" yaml:"min-streams"`
	MaxStreams     int    `json:"max-streams" yaml:"max-streams"`
	Statistic      bool   `json:"statistic" yaml:"statistic"`
	OnlyTCP        bool   `json:"only-tcp" yaml:"only-tcp"`
	Padding        bool   `json:"padding" yaml:"padding"`
	BrutalOpts     struct {
		Enabled bool `json:"enabled" yaml:"enabled"`
		Up      int  `json:"up" yaml:"up"`
		Down    int  `json:"down" yaml:"down"`
	} `json:"brutal-opts" yaml:"brutal-opts"`
}

type TLSConfig struct {
	TLS               bool              `json:"tls" yaml:"tls"`
	SNI               string            `json:"sni" yaml:"sni"`
	ServerName        string            `json:"servername" yaml:"servername"`
	Fingerprint       string            `json:"fingerprint" yaml:"fingerprint"`
	ALPN              []string          `json:"alpn" yaml:"alpn"`
	SkipCertVerify    bool              `json:"skip-cert-verify" yaml:"skip-cert-verify"`
	ClientFingerprint string            `json:"client-fingerprint" yaml:"client-fingerprint"`
	RealityOpts       *RealityTlsConfig `json:"reality-opts" yaml:"reality-opts"`
}

type RealityTlsConfig struct {
	PublicKey string `json:"public-key" yaml:"public-key"`
	ShortID   string `json:"short-id" yaml:"short-id"`
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
	Network string `json:"network" yaml:"network"`

	HTTPOpts *HTTPNetworkConfig `json:"http-opts" yaml:"http-opts"`
	H2Opts   *H2NetworkConfig   `json:"h2-opts" yaml:"h2-opts"`
	GRPCOpts *GRPCNetworkConfig `json:"grpc-opts" yaml:"grpc-opts"`
	WSOpts   *WSNetworkConfig   `json:"ws-opts" yaml:"ws-opts"`
}

type HTTPNetworkConfig struct {
	Method  string              `json:"method" yaml:"method"`
	Path    []string            `json:"path" yaml:"path"`
	Headers map[string][]string `json:"headers" yaml:"headers"`
}

// h2-opts:
//
//	host:
//	- example.com
//	path: /
type H2NetworkConfig struct {
	Host []string `json:"host" yaml:"host"`
	Path string   `json:"path" yaml:"path"`
}

type GRPCNetworkConfig struct {
	GRPCServiceName string `json:"grpc-service-name" yaml:"grpc-service-name"`
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
	Path                     string              `json:"path" yaml:"path"`
	Headers                  map[string][]string `json:"headers" yaml:"headers"`
	MaxEarlyData             int                 `json:"max-early-data" yaml:"max-early-data"`
	EarlyDataHeaderName      string              `json:"early-data-header-name" yaml:"early-data-header-name"`
	V2RayHTTPUpgrade         bool                `json:"v2ray-http-upgrade" yaml:"v2ray-http-upgrade"`
	V2RayHTTPUpgradeFastOpen bool                `json:"v2ray-http-upgrade-fast-open" yaml:"v2ray-http-upgrade-fast-open"`
}
