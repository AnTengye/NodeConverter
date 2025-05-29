package core

import "testing"

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewVLESSNode()
			if err := node.FromShare(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("FromShare() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(node.ToClash())
		})
	}
}
