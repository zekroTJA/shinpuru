package util

func EnsureNotEmpty(str, def string) string {
	if str == "" {
		return def
	}
	return str
}
