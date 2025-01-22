package core

import "testing"

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &TrojanNode{
				Normal:        tt.fields.Normal,
				TLSConfig:     tt.fields.TLSConfig,
				NetworkConfig: tt.fields.NetworkConfig,
				Password:      tt.fields.Password,
				SSOpts:        tt.fields.SSOpts,
			}
			if err := node.FromShare(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("FromShare() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
