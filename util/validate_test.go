package util

import "testing"

func TestValidateProjectName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"", args{
			"-",
		}, true},
		{"", args{
			"a",
		}, false},
		{"", args{
			"_-",
		}, false},
		{"", args{
			"_--___--__",
		}, false},
		{"", args{
			"_name",
		}, false},
		{"", args{
			"pro",
		}, false},
		{"", args{
			" pro",
		}, true},
		{"", args{
			"p ro",
		}, true},
		{"", args{
			"pro ",
		}, true},
		{"", args{
			"pro_",
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateProjectName(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("ValidateProjectName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateUserName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"", args{
			"-",
		}, true},
		{"", args{
			"a",
		}, false},
		{"", args{
			"_-",
		}, false},
		{"", args{
			"_--___--__",
		}, false},
		{"", args{
			"_name",
		}, false},
		{"", args{
			"pro",
		}, false},
		{"", args{
			" pro",
		}, true},
		{"", args{
			"p ro",
		}, true},
		{"", args{
			"pro ",
		}, true},
		{"", args{
			"pro_",
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateUserName(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("ValidateUserName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
