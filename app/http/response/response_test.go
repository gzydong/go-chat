package response

import (
	"errors"
	"reflect"
	"testing"

	"go-chat/app/entity"
)

func TestNewError(t *testing.T) {
	type args struct {
		code    int
		message []interface{}
	}
	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "测试01",
			args: args{code: 500, message: nil},
			want: &Response{Code: 500, Message: entity.CodeMessageMap[500]},
		},
		{
			name: "测试02",
			args: args{code: 501, message: []interface{}{"test01"}},
			want: &Response{Code: 501, Message: "test01"},
		},
		{
			name: "测试03",
			args: args{code: 400, message: nil},
			want: &Response{Code: 400, Message: "业务错误"},
		},
		{
			name: "测试04",
			args: args{code: 400, message: []interface{}{errors.New("my test")}},
			want: &Response{Code: 400, Message: "my test"},
		},
		{
			name: "测试05",
			args: args{code: 400, message: []interface{}{1}},
			want: &Response{Code: 400, Message: "1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewError(tt.args.code, tt.args.message...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewError() = %v, want %v", got, tt.want)
			}
		})
	}
}
