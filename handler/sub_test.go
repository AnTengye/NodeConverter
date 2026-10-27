package handler

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/AnTengye/NodeConverter/lib/network"
	"github.com/kataras/iris/v12"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func TestSubPreservesRealityOptionsAndCertificatePin(t *testing.T) {
	template := filepath.Join(t.TempDir(), "clash-tpl.yaml")
	if err := os.WriteFile(template, []byte("proxy-groups: []\nrules: []\n"), 0600); err != nil {
		t.Fatal(err)
	}
	previousTemplate := viper.Get("Advanced.TemplateFilePath")
	viper.Set("Advanced.TemplateFilePath", template)
	t.Cleanup(func() { viper.Set("Advanced.TemplateFilePath", previousTemplate) })

	app := iris.New()
	app.Get("/sub", Sub)
	if err := app.Build(); err != nil {
		t.Fatal(err)
	}

	const pin = "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789"
	const reality = "vless://00000000-0000-4000-8000-000000000001@[2001:db8::1]:443?security=reality&type=tcp&sni=example.com&fp=chrome&pbk=test-public-key&sid=0123456789abcdef&ip-version=ipv6&support-x25519mlkem768=true#ipv6-reality"
	const hy2 = "hysteria2://test-password@192.0.2.20:443?sni=example.com&pinSHA256=" + pin + "&insecure=0#hy2"
	clash := "proxies:\n- name: ipv6-reality\n  type: vless\n  server: 2001:db8::1\n  port: 443\n  uuid: 00000000-0000-4000-8000-000000000001\n  tls: true\n  ip-version: ipv6\n  reality-opts:\n    public-key: test-public-key\n    short-id: 0123456789abcdef\n    support-x25519mlkem768: true\n- name: hy2\n  type: hysteria2\n  server: 192.0.2.20\n  port: 443\n  password: test-password\n  fingerprint: " + pin + "\n  skip-cert-verify: false\n"
	previousCacheGet := network.CacheGET
	t.Cleanup(func() { network.CacheGET = previousCacheGet })

	for _, tc := range []struct {
		name    string
		source  string
		content []byte
	}{
		{name: "share-links", source: reality + "|" + hy2},
		{name: "base64-subscription", source: "https://example.com/base64", content: []byte(base64.StdEncoding.EncodeToString([]byte(reality + "\n" + hy2)))},
		{name: "clash-subscription", source: "https://example.com/clash", content: []byte(clash)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			network.CacheGET = func(string) ([]byte, error) { return tc.content, nil }
			query := url.Values{"target": {"clashmeta"}, "url": {tc.source}}
			response := httptest.NewRecorder()
			app.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/sub?"+query.Encode(), nil))
			if response.Code != http.StatusOK {
				t.Fatalf("HTTP %d: %s", response.Code, response.Body.String())
			}
			var document struct {
				Proxies []map[string]any `yaml:"proxies"`
			}
			if err := yaml.Unmarshal(response.Body.Bytes(), &document); err != nil {
				t.Fatal(err)
			}
			if len(document.Proxies) != 2 {
				t.Fatalf("proxy count = %d, want 2", len(document.Proxies))
			}
			vless := document.Proxies[0]
			realityOpts, _ := vless["reality-opts"].(map[string]any)
			if vless["server"] != "2001:db8::1" || vless["ip-version"] != "ipv6" || realityOpts["support-x25519mlkem768"] != true {
				t.Errorf("IPv6 REALITY options lost in /sub: %+v", vless)
			}
			hysteria := document.Proxies[1]
			if hysteria["fingerprint"] != pin || hysteria["skip-cert-verify"] == true {
				t.Errorf("certificate settings lost in /sub: %+v", hysteria)
			}
		})
	}
}
