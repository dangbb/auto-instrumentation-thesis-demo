package main

import (
	"fmt"
	"html"
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/sirupsen/logrus"

	config2 "microservice/config"
)

func logLogrus() {
	// Test ability to group actions in same goroutine
	logrus.SetLevel(logrus.DebugLevel)

	logrus.Trace("Something very low level.")
	logrus.Debug("Useful debugging information.")
	logrus.Info("Something noteworthy happened!")
	logrus.Warn("You should probably take a look at this.")
}

func main() {
	config := config2.Config{}
	kong.Parse(&config)

	http.HandleFunc("/customer", func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("Get request to %s", "/customer")
		fmt.Fprintf(w, "Hello, %q\n", html.EscapeString(r.URL.Path))

		go func() {
			logLogrus()

			go func() {
				logLogrus()
			}()
		}()
	})

	logrus.Infof("Start customer service at: 0.0.0.0:%d", config.HttpPort)
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.HttpPort), nil); err != nil {
		logrus.Fatalf("can listen to port %d\n", config.HttpPort)
	}
}
