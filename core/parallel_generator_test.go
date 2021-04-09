package core

import "testing"

func TestParallelGenerator_Gen(t *testing.T) {
	type fields struct {
		subtasks []Generator
		cleaned  bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"generate same files", fields{
			[]Generator{
				NewFileGenerator("../test/test", func() []byte { return []byte("test") }),
				NewFileGenerator("../test/test", func() []byte { return []byte("test") }),
				NewFileGenerator("../test/test", func() []byte { return []byte("test") }),
				NewFileGenerator("../test/test", func() []byte { return []byte("test") }),
				NewFileGenerator("../test/test", func() []byte { return []byte("test") }),
			}, false,
		}, false},
		{"generate some files", fields{
			[]Generator{
				NewFileGenerator("../test/test1", func() []byte { return []byte("test") }),
				NewFileGenerator("../test/test2", func() []byte { return []byte("test") }),
				NewFileGenerator("../test/test3", func() []byte { return []byte("test") }),
				NewFileGenerator("../test/test4", func() []byte { return []byte("test") }),
				NewFileGenerator("../test/a/b/c/d/e/f/g/test", func() []byte { return []byte("test") }),
			}, false,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := ParallelGenerator{
				subtasks: tt.fields.subtasks,
				cleaned:  tt.fields.cleaned,
			}
			if err := pg.Gen(); (err != nil) != tt.wantErr {
				t.Errorf("ParallelGenerator.Gen() error = %v, wantErr %v", err, tt.wantErr)
			}
			pg.Cleanup()
		})
	}
}
