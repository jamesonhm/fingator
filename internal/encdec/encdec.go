package encdec

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html"

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

func DecodeHTMLResp(r *http.Response, v any) error {
	//bytes, _ := io.ReadAll(r.Body)
	//fmt.Println(string(bytes))

	node, err := html.Parse(r.Body)
	if err != nil {
		return fmt.Errorf("error parsing html")
	}
	//fmt.Printf("Node: %+v\n", node)
	if val, ok := v.(*html.Node); ok {
		*val = *node
	}
	return nil
}

func DecodeTxtResponse(r *http.Response, v any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	// Type assert v to a pointer to io.Reader
	rdrPtr, ok := v.(*io.Reader)
	if !ok {
		return fmt.Errorf("v must be a pointer to an io.Reader")
	}

	*rdrPtr = bytes.NewReader(body)
	return nil
}
