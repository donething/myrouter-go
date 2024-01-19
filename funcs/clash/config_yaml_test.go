package clash

import (
	"testing"
)

func Test_getProxyGroups(t *testing.T) {
	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{
			name:    "获取所有代理组",
			want:    10,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getProxyGroups()
			if (err != nil) != tt.wantErr {
				t.Errorf("getProxyGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("getProxyGroups() got = %v, want %v", got, tt.want)
			}
		})
	}
}
