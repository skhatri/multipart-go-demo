package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Server struct {
	saveDir string
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	multipartFile, header, err := r.FormFile("file")
	data := make(map[string]interface{}, 0)
	if err != nil {
		data["code"] = "bad_request"
		data["message"] = err.Error()
		s.send(w, data, 400)
	} else {
		d, err := ioutil.ReadAll(multipartFile)
		if err != nil {
			data["code"] = "read_error"
			data["message"] = err.Error()
			s.send(w, data, 400)
		} else {
			fileLoc := fmt.Sprintf("%s/%s", s.saveDir, strings.Replace(header.Filename, "/", "_", -1))
			newF, err := os.Create(fileLoc)
			if err != nil {
				data["code"] = "write_error"
				data["message"] = err.Error()
				s.send(w, data, 500)
			} else {
				newF.Write(d)
				data["filename"] = header.Filename
				data["length"] = header.Size
				data["stored"] = fileLoc
				s.send(w, data, 200)
			}
		}
	}
}

func (s *Server) send(w http.ResponseWriter, m map[string]interface{}, statusCode int) {
	w.WriteHeader(statusCode)
	buff := bytes.Buffer{}
	json.NewEncoder(&buff).Encode(m)
	w.Header().Add("Content-Type", "application/json")
	w.Write(buff.Bytes())
}

func main() {
	saveDir := "/tmp"
	if storePath := os.Getenv("STORE_PATH"); storePath != "" {
		saveDir = storePath
	}
	os.MkdirAll(saveDir, os.ModePerm)
	s := Server{
		saveDir: saveDir,
	}
	address := "0.0.0.0:8199"
	fmt.Println("listening multipart file upload server on", address,". parameter name is 'file'")
	fmt.Println(`
example upload using curl:

curl -X POST -H "Content-Type: multipart/form-data; boundary=123-UPLOAD-SEPARATOR" \
-d '--123-UPLOAD-SEPARATOR
Content-Disposition: form-data; name="file"; filename="test.txt"
Content-Type: text/plain

test
--123-UPLOAD-SEPARATOR--
' "http://localhost:8199"

`)
	http.ListenAndServe(address, &s)
}
