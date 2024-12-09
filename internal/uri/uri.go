package uri

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/jamesonhm/fingator/internal/polygon/models"
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
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (b *URIBuilder) EncodeParams(path string, params any) string {
	epath := encodePath(path, params)
	return b.baseURL + epath
}

func encodePath(path string, params interface{}) string {
	pv := reflect.ValueOf(params)
	i := reflect.Indirect(pv)
	pt := i.Type()
	for i := 0; i < pt.NumField(); i++ {
		field := pt.Field(i)
		tag := field.Tag.Get(pathTag)
		if tag != "" {
			ft := pv.Field(i).Type().String()
			fmt.Printf("PATH ENCODER: field Type = %s\n", ft)
			fv := pv.Field(i).Interface()
			sfv, _ := formatFieldValue(ft, fv)
			// insert tag and value (fv) to path
			path = strings.ReplaceAll(path, fmt.Sprintf("{%s}", tag), url.PathEscape(sfv))
		}
	}
	return path
}

func formatFieldValue(fieldType string, fieldValue interface{}) (string, error) {
	var sfv string
	switch fieldType {
	case "float64":
		sfv = fmt.Sprintf("%g", fieldValue)
	case "models.Date":
		sfv = fieldValue.(models.Date).Format()
	default:
		sfv = fmt.Sprintf("%v", fieldValue)
	}
	return sfv, nil
}
