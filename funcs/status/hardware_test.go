package status

import (
	"myrouter/models"
	"reflect"
	"testing"
)

func TestGetRouteStatus(t *testing.T) {
	tests := []struct {
		name string
		want *models.RouterStatus
	}{
		{
			name: "获取 CPU、内存的状态",
			want: &models.RouterStatus{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRouterStatus(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRouterStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
