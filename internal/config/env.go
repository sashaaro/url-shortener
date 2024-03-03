package config

import "os"

var HTTPPort = os.Getenv("HTTP_PORT")

func init() {
	if HTTPPort == "" {
		HTTPPort = "8080"
	}
}
