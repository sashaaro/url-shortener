// Package main наш главный испытуемый.
package main

import "os"

func main() {
	// Эти os.Exit() должны вызвать негодование анализатора.
	defer os.Exit(0) // want "Call os.Exit on function main of package main"
	go os.Exit(1)    // want `Call os.Exit on function main of package main`
	os.Exit(2)       // want "Call os.Exit on function main of package main"
}

func good() {
	// В этой функции разрешено использовать os.Exit()
	os.Exit(3)
}
