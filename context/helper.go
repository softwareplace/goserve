package context

import (
	"encoding/json"
	"io"
)

var (
	encoder = jsonEncoder
)

func jsonEncoder(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}
