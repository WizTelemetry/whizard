package util

import (
	"reflect"
	"testing"
)

func TestCreateKeepTenantsRelabelConfig(t *testing.T) {
	testCases := []struct {
		name    string
		label   string
		regex   []string
		want    string
		wantErr bool
	}{
		{
			name:    "null case",
			label:   "",
			regex:   []string{},
			want:    "- action: keep\n  source_labels:\n    - \"\"\n  regex: \"\"\n",
			wantErr: false,
		},
		{
			name:    "good case",
			label:   "tenant_id",
			regex:   []string{"host", "member-1"},
			want:    "- action: keep\n  source_labels:\n    - tenant_id\n  regex: ^host$|^member-1$\n",
			wantErr: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			got, err := CreateKeepTenantsRelabelConfig(tt.label, tt.regex)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateKeepTenantsRelabelConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateKeepTenantsRelabelConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
