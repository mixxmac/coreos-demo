package main

import (
	"log"
	"net"
	"net/http"
	"os"
)

var html = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>Hello</title>
  </head>
  <body>
    <h1>Hello!</h1>
  </body>
</html>
`

func helloServer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(html))
}

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", helloServer)
	log.Printf("Starting app on port %s ...", port)
	if err := http.ListenAndServe(net.JoinHostPort("0.0.0.0", port), nil); err != nil {
		log.Fatal(err)
	}
}
