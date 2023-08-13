package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var queue = make(chan string, 100)

func main() {
	cmd := os.Args[1]

	switch cmd {
	case "run":
		run()

	case "serve":
		go worker()
		serve()
	}
}

func run() {
	cmdStr := os.Args[2]
	http.Post("http://127.0.0.1:9000/api/run", "text/plain", strings.NewReader(cmdStr))
}

func serveHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/run":
		cmdStr, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		queue <- string(cmdStr)
	}
}

func worker() {
	for cmdStr := range queue {
		cmd := exec.Command("sh", "-c", cmdStr)
		out, _ := cmd.CombinedOutput()
		fmt.Print(string(out))
	}
}

func serve() {
	http.ListenAndServe(":9000", http.HandlerFunc(serveHandler))
}
