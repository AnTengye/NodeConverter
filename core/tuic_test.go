package core

import "testing"

func TestTuicNode_FromShare(t *testing.T) {
	type fields struct {
		Normal                Normal
		TLSConfig             TLSConfig
		Version               string
		Token                 string
		Uuid                  string
		Password              string
		Ip                    string
		HeartbeatInterval     int
		DisableSni            bool
		ReduceRtt             bool
		RequestTimeout        int
		UdpRelayMode          string
		CongestionController  string
		MaxUdpRelayPacketSize int
		FastOpen              bool
		MaxOpenStreams        int
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
			name:   "testurl",
			fields: fields{},
			args: args{
				s: "tuic://f3ba3679-f2f4-465f-92a5-de320b775a88:f3ba3679-f2f4-465f-92a5-de320b775a88@185.22.155.105:43516?alpn=h3&delay=594&allow_insecure=1&sni=www.bing.com&udp_relay_mode=native&version=5#%F0%9F%87%B7%F0%9F%87%BA%E4%BF%84%E7%BD%97%E6%96%AF1-%20%E2%AC%87%EF%B8%8F%202.6MB%2Fs",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewTUICNode()
			if err := node.FromShare(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("FromShare() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(node.ToClash())
		})
	}
}
