// Package utils
package utils

// Must кидает panic если err != null
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// Filter - фильтрация в слайсе
func Filter(slice []string, f func(string) bool) []string {
	var r []string
	for _, s := range slice {
		if f(s) {
			r = append(r, s)
		}
	}
	return r
}
