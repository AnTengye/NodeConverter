package core

import "testing"

func TestHysteriaNode_FromShare(t *testing.T) {
	type fields struct {
		Normal       Normal
		TLSConfig    TLSConfig
		version      int
		Ports        string
		Password     string
		Up           string
		Down         string
		Obfs         string
		ObfsPassword string
		AuthStr      string
		Protocol     string
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
				s: "hysteria://51.159.226.1:14241?alpn=h3&auth=dongtaiwang.com&auth_str=dongtaiwang.com&delay=1004&downmbps=100&protocol=udp&insecure=1&peer=apple.com&udp=true&upmbps=100#%F0%9F%87%AB%F0%9F%87%B7%E6%B3%95%E5%9B%BD1-%20%E2%AC%87%EF%B8%8F%209.4MB%2Fs",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewHYSTERIANode()
			if err := node.FromShare(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("FromShare() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(node.ToClash())
		})
	}
}
