package clash

import "testing"

func Test_addRule(t *testing.T) {
	type args struct {
		rule string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "成功添加新规则google",
			args:    args{rule: "- DOMAIN-SUFFIX,google.com,DIRECT"},
			wantErr: false,
		},
		{
			name:    "新规则不规范(域名前缺少逗号)",
			args:    args{rule: "- DOMAIN-SUFFIX steamcontent.com,DIRECT"},
			wantErr: true,
		}, {
			name:    "已存在新规则的域名",
			args:    args{rule: "- DOMAIN-SUFFIX,    gooGLe.com,DIRECT"},
			wantErr: true,
		},
		{
			name:    "成功添加新规则baidu",
			args:    args{rule: "- DOMAIN-SUFFIX,baidu.com,DIRECT"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := addRule(tt.args.rule); (err != nil) != tt.wantErr {
				t.Errorf("addRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_delRule(t *testing.T) {
	type args struct {
		rule string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "删除不存在的规则",
			args:    args{rule: "aabb.com"},
			wantErr: true,
		},
		{
			name:    "删除存在的规则",
			args:    args{rule: "- DOMAIN-SUFFIX,baidu.com,DIRECT"},
			wantErr: false,
		},
		{
			name:    "删除上条已删除的规则",
			args:    args{rule: "- DOMAIN-SUFFIX,baidu.com,DIRECT"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := delRule(tt.args.rule); (err != nil) != tt.wantErr {
				t.Errorf("delRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
