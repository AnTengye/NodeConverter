package core

import (
	"strings"
	"testing"
)

func TestTrojanNode_FromShare(t *testing.T) {
	type fields struct {
		Normal        Normal
		TLSConfig     TLSConfig
		NetworkConfig NetworkConfig
		Password      string
		SSOpts        TrojanSSOptions
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
			name: "trojan-cn",
			args: args{
				s: "trojan://0052c3c2-28bf-467d-b586-3b91d1c5bdc0@test.com:443?encryption=none&security=tls&type=h2&host=test.com&path=/0052c3c2-28bf-467d-b586-3b91d1c5bdc0#中文测试",
			},
			wantErr: false,
		},
		{
			name: "trojan-success",
			args: args{
				s: "trojan://BIRCauz5R2S8O3pZ56eg8yC43yla4T3D7znxCDFxSCODO9jA3NxAcaE3YYK3Sq8u9DxZg@closet.homeofbrave.net:18332?sni=closet.homeofbrave.net&fp=chrome#%F0%9F%87%BA%F0%9F%87%B8%E7%BE%8E%E5%9B%BD20-%20%E2%AC%87%EF%B8%8F%202.2MB%2Fs",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewTrojanNode()
			if err := node.FromShare(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("FromShare() error = %v, wantErr %v", err, tt.wantErr)
			}
			toShare := node.ToShare()
			node2 := NewTrojanNode()
			if err := node2.FromShare(toShare); (err != nil) != tt.wantErr {
				t.Errorf("ToShare() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(strings.Split(node.ToClash(), "\n")) != len(strings.Split(node2.ToClash(), "\n")) {
				t.Errorf("node1(%s)\n---------\n%s\nnode2(%s)\n---------\n%s", tt.args.s, node.ToClash(), toShare, node2.ToClash())
			}
		})
	}
}
