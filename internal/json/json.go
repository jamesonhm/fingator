package json

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func decodeResp(r *http.Response, v *any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("error decoding json: %w", err)
	}
	return nil
}
