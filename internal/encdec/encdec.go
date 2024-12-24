package encdec

import (
	"encoding/json"
	"encoding/xml"
	//"fmt"
	//"io"
	"net/http"

	"golang.org/x/net/html/charset"
	//"golang.org/x/text/encoding/charmap"
)

func DecodeJsonResp(r *http.Response, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func DecodeXmlResp(r *http.Response, v any) error {
	decoder := xml.NewDecoder(r.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	//decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
	//	if charset == "ISO-8859-1" {
	//		return charmap.ISO8859_1.NewDecoder().Reader(input), nil
	//	}
	//	return nil, fmt.Errorf("unsupported charset: %s", charset)
	//}
	return decoder.Decode(v)
}
