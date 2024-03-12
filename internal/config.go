package internal

import (
	"flag"
)

var HTTPAddr = flag.String("a", ":8080", "listen address")
var BaseURL = flag.String("b", "http://localhost:8080", "base url")
