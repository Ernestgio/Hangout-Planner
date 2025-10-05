package sanitizer

import (
	"bytes"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var (
	strictPolicy = bluemonday.StrictPolicy()
	ugcPolicy    = bluemonday.UGCPolicy()
)

func SanitizeString(s string) string {
	return strictPolicy.Sanitize(s)
}

func SanitizeMarkdown(md string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.New(goldmark.WithExtensions(extension.GFM)).Convert([]byte(md), &buf); err != nil {
		return "", err
	}

	safeHTML := ugcPolicy.Sanitize(buf.String())
	return safeHTML, nil
}
