package uri

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

const (
	pathTag  = "path"
	queryTag = "query"
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

func (b *URIBuilder) EncodeParams(path string, params interface{}) string {
	epath := encodePath(path, params)
	return epath
}

func encodePath(path string, params interface{}) string {
	pt := reflect.TypeOf(params)
	pv := reflect.ValueOf(params)
	for i := 0; i < pt.NumField(); i++ {
		field := pt.Field(i)
		tag := field.Tag.Get(pathTag)
		if tag != "" {
			//ft := pv.Field(i).Type().String()
			fv := pv.Field(i).String()
			// insert tag and value (fv) to path
			path = strings.ReplaceAll(path, fmt.Sprintf("{%s}", tag), url.PathEscape(fv))
		}
	}
	return path
}

//func (b *URIBuilder) AddPathParam(segment string) *URIBuilder {
//	b.pathSegments = append(b.pathSegments, url.PathEscape(segment))
//	return b
//}
//
//func (b *URIBuilder) AddQueryParam(key string, value any) *URIBuilder {
//	var strVal string
//	switch value.(type) {
//	case string:
//		strVal = value.(string)
//	case int:
//		strVal = fmt.Sprintf("%d", value)
//	case fmt.Stringer:
//		strVal = value.(fmt.Stringer).String()
//	default:
//		strVal = fmt.Sprintf("%v", value)
//	}
//
//	if strVal != "" {
//		b.queryParams.Add(key, strVal)
//	}
//	return b
//}
//
//func (b *URIBuilder) Build() string {
//	path := strings.Join(b.pathSegments, "/")
//	fullURL := fmt.Sprintf("%s/%s", b.baseURL, path)
//
//	if len(b.queryParams) > 0 {
//		fullURL += "?" + b.queryParams.Encode()
//	}
//	return fullURL
//}
