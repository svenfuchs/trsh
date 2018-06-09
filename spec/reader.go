package spec

import (
	"encoding/json"
	"github.com/svenfuchs/trsh/http"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type reader struct {
	http *http.HTTP
	path string
}

func (r reader) run() map[string]interface{} {
	var data []byte
	if _, err := os.Stat(r.path); err == nil {
		data = r.read()
	} else if r.http != nil {
		data = r.write(r.fetch())
	} else {
		panic("No http given, and spec file not found: " + string(r.path))
	}
	return r.parse(data).(map[string]interface{})
}

func (r reader) fetch() []byte {
	res, err := r.http.Get("/", nil)
	if err != nil {
		log.Fatal(err)
	}
	return res.Body
}

func (r reader) read() []byte {
	b, err := ioutil.ReadFile(r.path)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func (r reader) write(b []byte) []byte {
	err := os.MkdirAll(path.Dir(r.path), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(r.path, b, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return b
}

func (r reader) parse(b []byte) interface{} {
	var data interface{}
	err := json.Unmarshal(b, &data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
