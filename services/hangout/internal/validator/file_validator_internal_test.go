package validator

import "testing"

func TestIsValidMimeType_TableDriven(t *testing.T) {
	fv := NewFileValidator()

	tests := []struct {
		name string
		mime string
		ext  string
		want bool
	}{
		{name: "unknown extension returns false", mime: "image/svg+xml", ext: ".svg", want: false},
		{name: "jpeg with params", mime: "image/jpeg; charset=utf-8", ext: ".jpg", want: true},
		{name: "uppercase mime prefix", mime: "Image/JPEG", ext: ".jpg", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fv.isValidMimeType(tt.mime, tt.ext)
			if got != tt.want {
				t.Fatalf("isValidMimeType(%q, %q) = %v, want %v", tt.mime, tt.ext, got, tt.want)
			}
		})
	}
}
