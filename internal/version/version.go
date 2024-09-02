package version

import (
	"fmt"
	"io"
)

// Build структура для вывода била
type Build struct {
	BuildVersion string
	BuildDate    string
	BuildCommit  string
}

const hello = `Build version: %s
Build date: %s
Build commit: %s
`

// Print вывод билда
func (b Build) Print(w io.Writer) {
	p := []any{
		val(b.BuildVersion),
		val(b.BuildDate),
		val(b.BuildCommit),
	}
	if _, err := fmt.Fprintf(w, hello, p...); err != nil {
		panic("Fail to print the build version")
	}
}

func val(v string) string {
	if v == "" {
		return "N/A"
	}
	return v
}
