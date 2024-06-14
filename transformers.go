package memo

import "encoding/json"

func toJson(key string) Transformer {
	return func(s string) (string, error) {
		data := make(map[string]string)
		data[key] = s
		bz, err := json.Marshal(data)
		return string(bz), err
	}
}

func fromJson(key string) Transformer {
	return func(s string) (string, error) {
		data := make(map[string]string)
		err := json.Unmarshal([]byte(s), &data)
		if err != nil {
			return "", err
		}
		return data[key], nil
	}
}
