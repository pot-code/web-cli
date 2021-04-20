package core

import "testing"

func TestCmdExecutor_Run(t *testing.T) {
	type fields struct {
		name string
		args []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"", fields{
			"ls", []string{"-ahl"},
		}, false},
		{"", fields{
			"lab", []string{"-ahl"},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ce := CmdExecutor{
				name: tt.fields.name,
				args: tt.fields.args,
			}
			if err := ce.Run(); (err != nil) != tt.wantErr {
				t.Errorf("CmdExecutor.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
