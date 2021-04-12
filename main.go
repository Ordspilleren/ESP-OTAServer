package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var dataDir = os.Getenv("OTA_DATA_DIR")

func serveOTA(w http.ResponseWriter, req *http.Request) {
	//version := req.Header.Get("HTTP_X_ESP8266_VERSION")
	checksum := req.Header.Get("x-ESP8266-sketch-md5")

	urlPath := filepath.FromSlash(req.URL.Path)
	firmwareFile := filepath.Join(dataDir, urlPath, "firmware.bin")

	if _, err := os.Stat(firmwareFile); os.IsNotExist(err) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	file, err := os.Open(firmwareFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}

	fileMD5 := hex.EncodeToString(h.Sum(nil))
	log.Printf("Received OTA request!\nDevice checksum: %s\nServer checksum: %s", checksum, fileMD5)

	if fileMD5 != "" && fileMD5 != checksum {
		w.Header().Set("x-MD5", fileMD5)
		http.ServeFile(w, req, firmwareFile)
	} else {
		http.Error(w, http.StatusText(http.StatusNotModified), http.StatusNotModified)
	}
}

func main() {
	http.HandleFunc("/", serveOTA)

	http.ListenAndServe(":8080", nil)
}
