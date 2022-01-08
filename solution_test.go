package main

import "testing"

func TestParams(t *testing.T) {
	type args struct {
		uri    string
		params map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				uri: "/api/judge",
				params: map[string]interface{}{
					"key": "value",
				},
			},
			want: "https://oj.ismdeep.com/api/judge?key=value",
		},
		{
			name: "",
			args: args{
				uri: "/api/judge_api/get_pending",
				params: map[string]interface{}{
					"secure_code": "123456789",
					"sid":         1000,
				},
			},
			want: "https://oj.ismdeep.com/api/judge_api/get_pending?secure_code=123456789&sid=1000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateURL(tt.args.uri, tt.args.params); got != tt.want {
				t.Errorf("Params() = %v, want %v", got, tt.want)
			}
		})
	}
}
