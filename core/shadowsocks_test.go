package core

import "testing"

func TestShadowsocksNode_FromShare(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "SIP002",
			args: args{
				s: "ss://Y2hhY2hhMjAtaWV0Zi1wb2x5MTMwNTpmY2I4YTlhYy03YzQyLTQ5YWMtYjc2MS1lM2QwNmYwZWRiN2NAdGVzdC5jb206MTIzNAo=#%E6%9B%B4%EF%BC%9A01%2F22%7C%E7%BB%88%EF%BC%9A2025%2F02%2F28",
			},
			wantErr: false,
		},
		{
			name: "ss",
			args: args{
				s: "ss://chacha20-ietf-poly1305:fcb8a9ac-7c42-49ac-b761-e3d06f0edb7c@test.com:1234#%E6%9B%B4%EF%BC%9A01%2F22%7C%E7%BB%88%EF%BC%9A2025%2F02%2F28",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &ShadowsocksNode{}
			if err := node.FromShare(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("FromShare() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
