// Package pkg призван проверить что тестируемый анализёр не реагирует на os.Exit() в не main пакетах.
package pkg

import "os"

func main() {
	// все os.Exit() из этого пакета должны игнорироваться анализатором
	defer os.Exit(0)
	go os.Exit(1)
	os.Exit(2)
}
