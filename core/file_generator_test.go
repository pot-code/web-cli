package core

import "testing"

func TestFileGenerator_Gen(t *testing.T) {
	type fields struct {
		file    string
		data    DataProvider
		cleaned bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"", fields{
			"../test/test", func() []byte { return []byte("test") }, false,
		}, false},
		{"generate same file", fields{
			"../test/test", func() []byte { return []byte("test") }, false,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gt := &FileGenerator{
				file:    tt.fields.file,
				data:    tt.fields.data,
				cleaned: tt.fields.cleaned,
			}
			if err := gt.Gen(); (err != nil) != tt.wantErr {
				t.Errorf("FileGenerator.Gen() error = %v, wantErr %v", err, tt.wantErr)
			}
			gt.Cleanup()
		})
	}
}
