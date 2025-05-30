package core

import (
	"net/url"
	"strconv"
	"time"
)

func setBase(u *url.URL, normal *Normal) {
	normal.Type = NodeType(u.Scheme)
	portStr := u.Port()
	if portStr == "" {
		portStr = "443"
	}
	port, _ := strconv.Atoi(portStr)
	normal.Server = u.Hostname()
	normal.Port = port
	if u.Fragment != "" {
		normal.Name = u.Fragment
	} else {
		normal.Name = u.Scheme + "-" + time.Now().Format("15-04-05")
	}
}

func setNetwork(values url.Values, n *NetworkConfig) {
	n.Network = values.Get("type")
	if n.Network == "" {
		n.Network = values.Get("net")
	}
	switch n.Network {
	case "ws":
		n.WSOpts = &WSNetworkConfig{}
		n.WSOpts.Path = values.Get("path")
		if values.Has("host") {
			n.WSOpts.Headers = make(map[string]any)
			n.WSOpts.Headers["host"] = values.Get("host")
		}
	case "grpc":
		n.GRPCOpts = &GRPCNetworkConfig{}
		n.GRPCOpts.GRPCServiceName = values.Get("serviceName")
	case "h2":
		n.H2Opts = &H2NetworkConfig{}
		n.H2Opts.Host = []string{values.Get("host")}
		n.H2Opts.Path = values.Get("path")
	case "http":
		n.HTTPOpts = &HTTPNetworkConfig{}
		n.HTTPOpts.Method = values.Get("method")
		n.HTTPOpts.Path = []string{values.Get("path")}
	default:
		return
	}
}

func setTLS(values url.Values, t *TLSConfig) {
	t.SNI = values.Get("sni")
	t.ServerName = values.Get("sni")
	security := values.Get("security")
	if security == "reality" {
		t.RealityOpts = &RealityTlsConfig{}
		t.RealityOpts.PublicKey = values.Get("pbk")
		t.RealityOpts.ShortID = values.Get("sid")
	}
	t.TLS = security == "tls" || security == "reality"
	t.Fingerprint = values.Get("fp")
	t.ClientFingerprint = values.Get("fp")
	if values.Has("alpn") {
		t.ALPN = []string{values.Get("alpn")}
	}
	if values.Get("allowInsecure") == "1" || values.Get("insecure") == "1" || values.Get("allow_insecure") == "1" {
		t.SkipCertVerify = true
	}
}
