package status

import (
	"myrouter/models"
	"strings"
	"testing"
)

func TestGetLocalIPAddr(t *testing.T) {
	tests := []struct {
		name    string
		want    *models.IPAddr
		wantErr bool
	}{
		{
			name: "本地IP地址",
			want: &models.IPAddr{IPv6: "2408:"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLocalIPAddr()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLocalIPAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.HasPrefix(got.IPv6, tt.want.IPv6) {
				t.Errorf("GetLocalIPAddr() got = %v, want %v", got, tt.want)
			}
		})
	}
}
