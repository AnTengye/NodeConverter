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
			name:   "test1",
			fields: fields{},
			args: args{
				s: "hysteria://51.159.226.1:14241?alpn=h3&auth=dongtaiwang.com&auth_str=dongtaiwang.com&delay=1004&downmbps=100&protocol=udp&insecure=1&peer=apple.com&udp=true&upmbps=100#%F0%9F%87%AB%F0%9F%87%B7%E6%B3%95%E5%9B%BD1-%20%E2%AC%87%EF%B8%8F%209.4MB%2Fs",
			},
			wantErr: false,
		},
		{
			name:   "test2",
			fields: fields{},
			args: args{
				s: "hysteria2://dongtaiwang.com@51.159.111.32:31180?insecure=1&sni=apple.com#%F0%9F%87%AB%F0%9F%87%B7%E6%B3%95%E5%9B%BD3-%20%E2%AC%87%EF%B8%8F%202.1MB%2Fs",
			},
			wantErr: false,
		},
		{
			name:   "test3",
			fields: fields{},
			args: args{
				s: "hysteria://195.154.200.40:15010?alpn=h3&auth=dongtaiwang.com&auth_str=dongtaiwang.com&delay=1086&downmbps=100&obfsParam=&insecure=1&peer=apple.com&upmbps=100#%F0%9F%87%AB%F0%9F%87%B7%E6%B3%95%E5%9B%BD4-%20%E2%AC%87%EF%B8%8F%205.9MB%2Fs\", \"error\": \"share_url[hysteria://195.154.200.40:15010?alpn=h3&auth=dongtaiwang.com&auth_str=dongtaiwang.com&delay=1086&downmbps=100&obfsParam=&insecure=1&peer=apple.com&upmbps=100#%F0%9F%87%AB%F0%9F%87%B7%E6%B3%95%E5%9B%BD4-%20%E2%AC%87%EF%B8%8F%205.9MB%2Fs",
			},
			wantErr: false,
		},
		{
			name:   "test-obs-password",
			fields: fields{},
			args: args{
				s: "hysteria2://dongtaiwang.com@108.181.5.130:57773?insecure=1&sni=apple.com#%F0%9F%87%BA%F0%9F%87%B8%E7%BE%8E%E5%9B%BD1-%20%E2%AC%87%EF%B8%8F%201.4MB%2Fs",
			},
			wantErr: false,
		},
		{
			name:   "test-updown",
			fields: fields{},
			args: args{
				s: "hysteria://108.181.24.77:11512?alpn=h3&auth_str=dongtaiwang.com&downmbps=1000 Mbps&insecure=1&peer=apple.com&udp=true&upmbps=1000 Mbps#%F0%9F%87%BA%F0%9F%87%B8%E7%BE%8E%E5%9B%BD3-%20%E2%AC%87%EF%B8%8F%203.8MB%2Fs",
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
