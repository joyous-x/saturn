package utils

func TernaryString(condition bool, s string, f string) string {
	if condition {
		return s
	}
	return f
}
