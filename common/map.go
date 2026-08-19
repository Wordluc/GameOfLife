package common

func ReverseMapToValue[key, value comparable](m map[key]value) map[value]key {
	n := make(map[value]key)
	for k, v := range m {
		n[v] = k
	}
	return n
}

func ReverseMapToSlice[key, value comparable](m map[key][]value) map[value][]key {
	n := make(map[value][]key)
	for key, values := range m {
		for _, v := range values {
			n[v] = append(n[v], key)
		}
	}
	return n
}
