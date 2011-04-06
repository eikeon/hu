package main

import (
	"os"
	"path"
	"flag"
	"http"
	"io"
	"log"
	"strings"
	"bufio"
)


func StaticHandler(w http.ResponseWriter, req *http.Request) {
	w.SetHeader("Content-Type", "text/css")
	if strings.Contains(req.URL.Path, "^") {
		w.SetHeader("Cache-Control", "max-age=3153600")
	} else {
		w.SetHeader("Cache-Control", "max-age=1, must-revalidate")
	}

	var filename = path.Join("static", req.URL.Path)

	f, err := os.Open(filename, os.O_RDONLY, 0)
	if err != nil {
		log.Print(err)
		w.SetHeader("Cache-Control", "max-age=10, must-revalidate")
		w.WriteHeader(http.StatusNotFound)
 		return
	}

	var input = bufio.NewReader(f)
	io.Copy(w, input)
}

var addr = flag.String("addr", ":9999", "http service address")

func main() {
	flag.Parse()
	http.Handle("h.eikeon.com/", http.HandlerFunc(StaticHandler))
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Print("ListenAndServe:", err)
	}
}

// TODO: gzip compression

	//"compress/gzip"

	//w.SetHeader("Vary", "Accept-Encoding")

	// w.SetHeader("Content-Encoding", "gzip")
	// ww, err := gzip.NewWriter(w)
	// if err != nil {
	// 	log.Print("gzip", err)
	// }
