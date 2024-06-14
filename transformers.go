package memo

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"io"
	"strings"
)

func toJson(key string) Transformer {
	return func(value string) (string, error) {
		data := make(map[string]string)

		var b bytes.Buffer
		w := zlib.NewWriter(&b)
		w.Write([]byte(value))
		w.Close()
		compressed := b.Bytes()
		data[key] = string(compressed)

		bz, err := json.Marshal(data)
		return string(bz), err
	}
}

func fromJson(key string) Transformer {
	return func(memo string) (string, error) {
		data := make(map[string]string)
		err := json.Unmarshal([]byte(memo), &data)
		if err != nil {
			return "", err
		}
		compressed, ok := data[key]
		if !ok {
			return "", nil
		}
		r, err := zlib.NewReader(strings.NewReader(compressed))
		if err != nil {
			return "", err
		}
		buf := new(strings.Builder)
		_, err = io.Copy(buf, r)
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	}
}
