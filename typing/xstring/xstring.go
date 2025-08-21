package xstring

func Pick(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
