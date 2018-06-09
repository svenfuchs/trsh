package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

var conf = []string{
	"endpoint",
}

func IsOpt(str string) bool {
	for _, opt := range conf {
		if opt == str {
			return true
		}
	}
	return false
}

func MatchesOpt(str string) bool {
	for _, opt := range conf {
		if strings.HasPrefix(opt, str) {
			return true
		}
	}
	return false
}

func New(path string) *Conf {
	return &Conf{path: path}
}

type Conf struct {
	path string
}

func (o Conf) Get(key string) (interface{}, bool) {
	data := o.read()
	if v, ok := data[key]; ok {
		return v, true
	}
	return "", false
}

func (o Conf) Set(key string, value interface{}) {
	data := o.read()
	data[key] = value
	o.write(data)
}

func (o Conf) Del(key string) {
	data := o.read()
	delete(data, key)
	o.write(data)
}

func (o Conf) read() map[string]interface{} {
	b, err := ioutil.ReadFile(o.path)
	if err != nil {
		b = []byte("{}")
	}
	return o.parse(b)
}

func (o Conf) write(data map[string]interface{}) []byte {
	err := os.MkdirAll(path.Dir(o.path), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	b := o.dump(data)
	err = ioutil.WriteFile(o.path, b, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return b
}

func (o Conf) dump(data map[string]interface{}) []byte {
	b, _ := json.MarshalIndent(data, "", "  ")
	return b
}

func (o Conf) parse(b []byte) map[string]interface{} {
	var data interface{}
	err := json.Unmarshal(b, &data)
	if err != nil {
		log.Fatal(err)
	}
	return data.(map[string]interface{})
}
