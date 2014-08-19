package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/coreos/go-etcd/etcd"
)

var (
	client *etcd.Client
	data   map[string]string
	mu     = &sync.Mutex{}
)

func init() {
	data = make(map[string]string)
	client = etcd.NewClient([]string{"http://127.0.0.1:4001"})
}

func watchConfig() {
	for {
		resp, err := client.Watch("/app", 0, true, nil, nil)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		log.Println("Got a new key: " + resp.Node.Key)
		mu.Lock()
		data[filepath.Base(resp.Node.Key)] = resp.Node.Value
		mu.Unlock()
	}
}

var html = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>App Database Settings</title>
  </head>
  <body>
    <h1>Database Settings</h1>
    <table>
      <tr>
        <td>host</td>
        <td>{{.host}}</td>
      </tr>
      <tr>
        <td>username</td>
        <td>{{.username}}</td>
      </tr>
      <tr>
        <td>password</td>
        <td>{{.password}}</td>
      </tr>
    </table>
  </body>
</html>
`

func appServer(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("html").Parse(html)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
}

func main() {
	go watchConfig()
	port := os.Getenv("PORT")
	http.HandleFunc("/", appServer)
	log.Printf("Starting app on port %s ...", port)
	if err := http.ListenAndServe(net.JoinHostPort("0.0.0.0", port), nil); err != nil {
		log.Fatal(err)
	}
}
