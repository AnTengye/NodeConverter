package core

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestVlessNodePreservesIPv6RealityOptions(t *testing.T) {
	for _, option := range []string{"true", "false", ""} {
		for _, format := range []string{"share", "clash"} {
			t.Run(format+"/support="+option, func(t *testing.T) {
				node := NewVLESSNode()
				if format == "share" {
					link := "vless://00000000-0000-4000-8000-000000000001@[2001:db8::1]:443?security=reality&type=tcp&sni=example.com&fp=chrome&pbk=test-public-key&sid=0123456789abcdef&ip-version=ipv6#ipv6-reality"
					if option != "" {
						link = strings.Replace(link, "#", "&support-x25519mlkem768="+option+"#", 1)
					}
					if err := node.FromShare(link); err != nil {
						t.Fatal(err)
					}
				} else {
					input := "name: ipv6-reality\ntype: vless\nserver: 2001:db8::1\nport: 443\nuuid: 00000000-0000-4000-8000-000000000001\ntls: true\nservername: example.com\nclient-fingerprint: chrome\nip-version: ipv6\nreality-opts:\n  public-key: test-public-key\n  short-id: 0123456789abcdef\n"
					if option != "" {
						input += "  support-x25519mlkem768: " + option + "\n"
					}
					if err := node.FromClash([]byte(input)); err != nil {
						t.Fatal(err)
					}
				}

				assertReality := func(node Node) {
					t.Helper()
					var proxy map[string]any
					if err := yaml.Unmarshal([]byte(node.ToClash()), &proxy); err != nil {
						t.Fatal(err)
					}
					if proxy["server"] != "2001:db8::1" || proxy["ip-version"] != "ipv6" {
						t.Errorf("IPv6 server or preference lost: %+v", proxy)
					}
					reality, ok := proxy["reality-opts"].(map[string]any)
					if !ok {
						t.Fatalf("missing REALITY options: %+v", proxy)
					}
					value, exists := reality["support-x25519mlkem768"]
					if option == "" {
						if exists {
							t.Errorf("unspecified REALITY compatibility option was added: %v", value)
						}
					} else if !exists || value != (option == "true") {
						t.Errorf("REALITY compatibility option = %v (exists %v), want %s", value, exists, option)
					}
				}
				assertReality(node)

				share := node.ToShare()
				parsed, err := url.Parse(share)
				if err != nil {
					t.Fatalf("invalid exported IPv6 link: %v", err)
				}
				if parsed.Host != "[2001:db8::1]:443" || parsed.Query().Get("ip-version") != "ipv6" || parsed.Query().Get("support-x25519mlkem768") != option {
					t.Errorf("exported link lost IPv6 or REALITY settings: %s", share)
				}
				restored := NewVLESSNode()
				if err := restored.FromShare(share); err != nil {
					t.Fatal(err)
				}
				assertReality(restored)
			})
		}
	}
}

func TestHysteriaNodePreservesCertificatePin(t *testing.T) {
	const pin = "AB:CD:EF:01:23:45:67:89:AB:CD:EF:01:23:45:67:89:AB:CD:EF:01:23:45:67:89:AB:CD:EF:01:23:45:67:89"
	for _, format := range []string{"share", "clash"} {
		for _, insecure := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/insecure=%t", format, insecure), func(t *testing.T) {
				node := NewHYSTERIANode()
				if format == "share" {
					insecureValue := "0"
					if insecure {
						insecureValue = "1"
					}
					link := "hysteria2://test-password@192.0.2.20:443?sni=example.com&fp=chrome&pinSHA256=" + url.QueryEscape(pin) + "&insecure=" + insecureValue + "#hy2"
					if err := node.FromShare(link); err != nil {
						t.Fatal(err)
					}
				} else {
					input := fmt.Sprintf("name: hy2\ntype: hysteria2\nserver: 192.0.2.20\nport: 443\npassword: test-password\nsni: example.com\nfingerprint: %s\nclient-fingerprint: chrome\nskip-cert-verify: %t\n", pin, insecure)
					if err := node.FromClash([]byte(input)); err != nil {
						t.Fatal(err)
					}
				}

				assertPin := func(node Node) {
					t.Helper()
					var proxy map[string]any
					if err := yaml.Unmarshal([]byte(node.ToClash()), &proxy); err != nil {
						t.Fatal(err)
					}
					if proxy["fingerprint"] != pin || proxy["client-fingerprint"] != "chrome" {
						t.Errorf("certificate and client fingerprints were lost or mixed: %+v", proxy)
					}
					skip, _ := proxy["skip-cert-verify"].(bool)
					if skip != insecure {
						t.Errorf("skip-cert-verify = %t, want %t", skip, insecure)
					}
				}
				assertPin(node)

				share := node.ToShare()
				parsed, err := url.Parse(share)
				if err != nil {
					t.Fatal(err)
				}
				if parsed.Query().Get("pinSHA256") != pin || parsed.Query().Get("fp") != "chrome" {
					t.Errorf("exported link lost or mixed fingerprints: %s", share)
				}
				restored := NewHYSTERIANode()
				if err := restored.FromShare(share); err != nil {
					t.Fatal(err)
				}
				assertPin(restored)
			})
		}
	}
}

func TestHysteriaNodeClientFingerprintDoesNotBecomeCertificatePin(t *testing.T) {
	node := NewHYSTERIANode()
	if err := node.FromShare("hysteria2://test-password@192.0.2.20:443?sni=example.com&fp=chrome#hy2"); err != nil {
		t.Fatal(err)
	}
	var proxy map[string]any
	if err := yaml.Unmarshal([]byte(node.ToClash()), &proxy); err != nil {
		t.Fatal(err)
	}
	if _, exists := proxy["fingerprint"]; exists {
		t.Errorf("client fingerprint became a certificate pin: %+v", proxy)
	}
}
