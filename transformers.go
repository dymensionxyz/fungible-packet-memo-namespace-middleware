package memo

func toJson(key string) Transformer {
	return func(s string) string {
		data := make(map[string]string)
		data[key] = s
	}
}

func fromJson(key string) Transformer {
}
