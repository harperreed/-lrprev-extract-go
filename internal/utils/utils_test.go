package utils

import (
	"testing"
)

func TestExtractUUIDFromFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
		wantErr  bool
	}{
		{
			name:     "Valid filename with UUID",
			filename: "preview-12345678-1234-5678-1234-567812345678.lrprev",
			want:     "12345678-1234-5678-1234-567812345678",
			wantErr:  false,
		},
		{
			name:     "Filename without UUID",
			filename: "preview-nouuid.lrprev",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "Filename with invalid UUID format",
			filename: "preview-12345678-1234-5678-1234-56781234567G.lrprev",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "Empty filename",
			filename: "",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "Filename with multiple UUIDs",
			filename: "preview-12345678-1234-5678-1234-567812345678-87654321-4321-8765-4321-876543210987.lrprev",
			want:     "12345678-1234-5678-1234-567812345678",
			wantErr:  false,
		},
		{
			name:     "Filename with UUID not in base name",
			filename: "/path/to/12345678-1234-5678-1234-567812345678/preview.lrprev",
			want:     "12345678-1234-5678-1234-567812345678",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractUUIDFromFilename(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractUUIDFromFilename() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractUUIDFromFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}
