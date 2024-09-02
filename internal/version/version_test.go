package version

import (
	"bytes"
	"testing"
)

func TestBuild_Print(t *testing.T) {
	type fields struct {
		BuildVersion string
		BuildDate    string
		BuildCommit  string
	}
	tests := []struct {
		name   string
		fields fields
		wantW  string
	}{
		{
			name: "clean",
			wantW: `Build version: N/A
Build date: N/A
Build commit: N/A
`,
		},
		{
			name:   "specified version",
			fields: fields{"12", "34", "56"},
			wantW: `Build version: 12
Build date: 34
Build commit: 56
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Build{
				BuildVersion: tt.fields.BuildVersion,
				BuildDate:    tt.fields.BuildDate,
				BuildCommit:  tt.fields.BuildCommit,
			}
			w := &bytes.Buffer{}
			b.Print(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Print() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
