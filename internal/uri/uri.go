package uri

import (
	"fmt"
	"net/url"
	"strings"
)

type URIBuilder struct {
	baseURL      string
	pathSegments []string
	queryParams  url.Values
}

func New(baseURL string) *URIBuilder {
	return &URIBuilder{
		baseURL:      strings.TrimRight(baseURL, "/"),
		pathSegments: []string{},
		queryParams:  url.Values{},
	}
}

func (b *URIBuilder) AddPathParam(segment string) *URIBuilder {
	b.pathSegments = append(b.pathSegments, url.PathEscape(segment))
	return b
}

func (b *URIBuilder) AddQueryParam(key string, value any) *URIBuilder {
	var strVal string
	switch value.(type) {
	case string:
		strVal = value.(string)
	case int:
		strVal = fmt.Sprintf("%d", value)
	case fmt.Stringer:
		strVal = value.(fmt.Stringer).String()
	default:
		strVal = fmt.Sprintf("%v", value)
	}

	if strVal != "" {
		b.queryParams.Add(key, strVal)
	}
	return b
}

func (b *URIBuilder) Build() string {
	path := strings.Join(b.pathSegments, "/")
	fullURL := fmt.Sprintf("%s/%s", b.baseURL, path)

	if len(b.queryParams) > 0 {
		fullURL += "?" + b.queryParams.Encode()
	}
	return fullURL
}
