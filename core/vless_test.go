package core

import (
	"strings"
	"testing"
)

func TestVlessNode_FromClash(t *testing.T) {
	type fields struct {
		Normal         Normal
		TLSConfig      TLSConfig
		NetworkConfig  NetworkConfig
		Uuid           string
		Flow           string
		PacketEncoding string
	}
	type args struct {
		s []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "networkhttp-test",
			args: args{
				s: []byte(`
name: "\U0001F1FA\U0001F1F8美国2 | ⬇️ 5.1MB/s"
network: ws
port: 443
server: 160.123.255.23
servername: ll.bgm2024.dpdns.org
skip-cert-verify: true
tls: true
type: vless
udp: true
uuid: 55520747-311e-4015-83ce-be46e2060ce3
ws-opts:
  headers:
    host: ll.bgm2024.dpdns.org
  path: /?ed=2560
`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &VlessNode{}
			if err := node.FromClash(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("FromClash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVlessNode_FromShare(t *testing.T) {
	type fields struct {
		Normal         Normal
		TLSConfig      TLSConfig
		NetworkConfig  NetworkConfig
		Uuid           string
		Flow           string
		PacketEncoding string
	}
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "test1",
			fields: fields{},
			args: args{
				s: "vmess://eyJ2IjoiMiIsInBzIjoi8J+HuvCfh7jnvo7lm70xNy0g4qyH77iPIDEuMk1CL3MiLCJhZGQiOiIxNTAuMjMwLjM1LjExMyIsInBvcnQiOjgwODAsImlkIjoiOTQ2MjU4NGItOWVkNi00NjQ5LTk5MTUtMmFmZDhlNDM2ZDc3IiwidHlwZSI6IiIsImFpZCI6MCwibmV0Ijoid3MiLCJ0bHMiOiIiLCJwYXRoIjoiLzk0NjI1ODRiLTllZDYtNDY0OS05OTE1LTJhZmQ4ZTQzNmQ3Ny12bSIsImhvc3QiOiJ3d3cuYmluZy5jb20ifQ==",
			},
			wantErr: false,
		},
		{
			name:   "test-vless-succes",
			fields: fields{},
			args: args{
				s: "vless://55520747-311e-4015-83ce-be46e2060ce3@220.118.248.123:30012?encryption=none&security=tls&sni=ca.bgm2024.dpdns.org&allowInsecure=1&type=ws&host=ca.bgm2024.dpdns.org&path=%2F%3Fed%3D2560#%F0%9F%87%BA%F0%9F%87%B8%E7%BE%8E%E5%9B%BD5-%20%E2%AC%87%EF%B8%8F%203.4MB%2Fs",
			},
			wantErr: false,
		},
		{
			name:   "test-vless-ws-opts",
			fields: fields{},
			args: args{
				s: "vless://53fa8faf-ba4b-4322-9c69-a3e5b1555049@185.59.218.20:8880?security=none&type=ws&path=%2F%3Fed%3D2560&host=reedfree8mahsang2.redorg.ir&sni=reedfree8mahsang2.redorg.ir#%F0%9F%87%AD%F0%9F%87%B0%E9%A6%99%E6%B8%AF1-%20%E2%AC%87%EF%B8%8F%202.9MB%2Fs",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewVLESSNode()
			if err := node.FromShare(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("FromShare() error = %v, wantErr %v", err, tt.wantErr)
			}
			toShare := node.ToShare()
			node2 := NewVLESSNode()
			if err := node2.FromShare(toShare); (err != nil) != tt.wantErr {
				t.Errorf("ToShare() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(strings.Split(node.ToClash(), "\n")) != len(strings.Split(node2.ToClash(), "\n")) {
				t.Errorf("node1(%s)\n---------\n%s\nnode2(%s)\n---------\n%s", tt.args.s, node.ToClash(), toShare, node2.ToClash())
			}
		})
	}
}
