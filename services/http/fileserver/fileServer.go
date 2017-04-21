package main

import (
	"flag"
	"fmt"
	"net/http"
	"path/filepath"
)

var dirFlag = flag.String("dir", "", "Directory which will server as document root")
var portFlag = flag.String("port", "8080", "Port on which server will listen")

func main() {
	flag.Parse()
	dir := filepath.Clean(*dirFlag)
	fmt.Printf("Serving files from dir %s\n", dir)
	http.ListenAndServe(":"+(*portFlag), http.FileServer(http.Dir(dir)))
}
