package conf_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/svenfuchs/trsh/conf"
	"io/ioutil"
	"os"
	"testing"
)

var path = "/tmp/trsh.test/config.json"

func setup() {
	os.Remove(path)
}

func set(name string, value interface{}) {
	c := conf.New(path)
	c.Set(name, value)
}

func get(name string) (interface{}, bool) {
	c := conf.New(path)
	return c.Get(name)
}

func del(name string) {
	c := conf.New(path)
	c.Del(name)
}

func read() string {
	b, _ := ioutil.ReadFile(path)
	return string(b)
}

func TestOptsSetString(t *testing.T) {
	setup()
	set("foo", "bar")
	assert.Equal(t, "{\n  \"foo\": \"bar\"\n}", read())
}

func TestOptsSetInt(t *testing.T) {
	setup()
	set("foo", 1)
	assert.Equal(t, "{\n  \"foo\": 1\n}", read())
}

func TestOptsSetBool(t *testing.T) {
	setup()
	set("foo", true)
	assert.Equal(t, "{\n  \"foo\": true\n}", read())
}

func TestOptsGetString(t *testing.T) {
	setup()
	set("foo", "bar")
	value, ok := get("foo")
	assert.True(t, ok)
	assert.Equal(t, "bar", value)
}

func TestOptsGetInt(t *testing.T) {
	setup()
	set("foo", 1)
	value, ok := get("foo")
	assert.True(t, ok)
	assert.Equal(t, 1.0, value) // wat?
}

func TestOptsGetBool(t *testing.T) {
	setup()
	set("foo", true)
	value, ok := get("foo")
	assert.True(t, ok)
	assert.Equal(t, true, value)
}

func TestOptsDel(t *testing.T) {
	setup()
	set("foo", "foo")
	set("bar", "bar")
	del("foo")
	value, ok := get("foo")
	assert.False(t, ok)
	assert.Equal(t, "", value)
	assert.Equal(t, "{\n  \"bar\": \"bar\"\n}", read())
}
