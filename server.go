package main

import (
	"log"
	"os"

	"net/http"

	"github.com/grigagod/chat-example/server"
	"golang.org/x/net/websocket"
)

var envVars = map[string]string{
	"PORT":    "",
	"PQL_DSN": "",
}

func checkEnvVars() {
	for k := range envVars {
		val := os.Getenv(k)
		if val == "" {
			log.Fatalf("Error: environment variable %s is not set", k)
		}
		envVars[k] = val
	}
}

func main() {
	os.Setenv("PORT", "8001")
	os.Setenv("PQL_DSN", "host=localhost user=qiryl password=zxcv dbname=chat port=5432 TimeZone=Europe/Minsk")

	checkEnvVars()

	serverConfig := server.Config{
		DSN:       envVars["PQL_DSN"],
		Keepalive: 15,
	}

	server := server.CreateServer(serverConfig)

	log.Printf("Listening on port: %s\n", envVars["PORT"])

	http.Handle("/", websocket.Handler(server.WebsockHandler))

	err := http.ListenAndServe(":"+envVars["PORT"], nil)
	log.Printf("Error occurred in http listener: %s\n", err)
}
