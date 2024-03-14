package tse

import (
	"context"
	"testing"
)

func TestTseMiddleware_ExecCmdTask(t *testing.T) {
	type args struct {
		host   string
		cmd    string
		usrKey string
		sync   bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				host:   "11.71.58.246",
				cmd:    "ps ef",
				usrKey: "",
				sync:   false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tseMiddleware := TseMiddleware{tseUrl: "https://shell-api.alipay.com"}
			res := tseMiddleware.ExecCmdTask(context.Background(), tt.args.host, tt.args.cmd)
			if (res.ErrorCode != "") != tt.wantErr {
				return
			}
		})
	}
}
