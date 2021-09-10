package core

import "testing"

func TestCmdExecutor_Run(t *testing.T) {
	type fields struct {
		cmd *ShellCommand
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"", fields{
			cmd: &ShellCommand{
				Bin:  "ls",
				Args: []string{"-ahl"},
			},
		}, false},
		{"", fields{
			cmd: &ShellCommand{
				Bin:  "lab",
				Args: []string{"-ahl"},
			},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ce := ShellCmdExecutor{
				cmd: tt.fields.cmd,
			}
			if err := ce.Run(); (err != nil) != tt.wantErr {
				t.Errorf("CmdExecutor.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
