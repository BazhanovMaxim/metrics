package json

import (
	"bytes"
	"encoding/json"
	"io"
)

func MarshalToJSON(model any) ([]byte, error) {
	return json.Marshal(model)
}

func MarshalJSON(reader io.Reader, model any) error {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(reader); err != nil {
		return err
	}
	if err := json.Unmarshal(buf.Bytes(), model); err != nil {
		return err
	}
	return nil
}

func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func UnmarshalJSON(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
