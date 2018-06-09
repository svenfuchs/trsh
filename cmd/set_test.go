package cmd_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/svenfuchs/trsh/cmd"
	"github.com/svenfuchs/trsh/conf"
	"github.com/svenfuchs/trsh/ops"
	"io/ioutil"
	"os"
	"testing"
)

// Accept

func setAccept(str string) (ops.Op, bool) {
	cmd.Init()
	o := ops.New(conf.New("/tmp/trsh.test/config.json"))
	c := cmd.New(o)
	ok, _ := c.Accept(str)
	return c, ok
}

func TestSetAcceptColon(t *testing.T) {
	s, ok := setAccept(":")
	assert.True(t, ok)
	assert.Equal(t, ":", s.Str())
}

func TestSetAcceptColonPartOfSet(t *testing.T) {
	s, ok := setAccept(":s")
	assert.True(t, ok)
	assert.Equal(t, ":", s.Str())
}

func TestSetAcceptColonSet(t *testing.T) {
	s, ok := setAccept(":set")
	assert.True(t, ok)
	assert.Equal(t, ":set", s.Str())
}

func TestSetAcceptColonSetAndSpace(t *testing.T) {
	s, ok := setAccept(":set ")
	assert.True(t, ok)
	assert.Equal(t, ":set", s.Str())
}

func TestSetAcceptColonSetAndSpaceAndUnknownOptName(t *testing.T) {
	s, ok := setAccept(":set wat")
	assert.False(t, ok)
	assert.Equal(t, ":set", s.Str())
}

func TestSetAcceptColonSetAndSpaceAndPartOfOptName(t *testing.T) {
	s, ok := setAccept(":set end")
	assert.True(t, ok)
	assert.Equal(t, ":set", s.Str())
}

func TestSetAcceptColonSetAndSpaceAndOptName(t *testing.T) {
	s, ok := setAccept(":set endpoint")
	assert.True(t, ok)
	assert.Equal(t, ":set endpoint=", s.Str())
}

func TestSetAcceptColonSetAndSpaceAndOptNameAndEqual(t *testing.T) {
	s, ok := setAccept(":set endpoint=")
	assert.True(t, ok)
	assert.Equal(t, ":set endpoint=", s.Str())
}

func TestSetAcceptColonSetAndSpaceAndOpt(t *testing.T) {
	s, ok := setAccept(":set endpoint=travis-ci.com")
	assert.True(t, ok)
	assert.Equal(t, ":set endpoint=travis-ci.com", s.Str())
}

// // Complete
//
// func TestSetCompleteEmpty(t *testing.T) {
// 	s := setSetup()
// 	s.setAccept("")
// 	assert.False(t, s.complete())
// }
//
// func TestSetCompleteUnknownOptName(t *testing.T) {
// 	s := setSetup()
// 	s.setAccept("wat")
// 	assert.False(t, s.complete())
// }
//
// func TestSetCompletePartOfOptName(t *testing.T) {
// 	s := setSetup()
// 	s.setAccept("end")
// 	assert.False(t, s.complete())
// }
//
// func TestSetCompleteOptName(t *testing.T) {
// 	s := setSetup()
// 	s.setAccept("endpoint")
// 	assert.False(t, s.complete())
// }
//
// func TestSetCompleteOptNameAndEqual(t *testing.T) {
// 	s := setSetup()
// 	s.setAccept("endpoint=")
// 	assert.False(t, s.complete())
// }
//
// func TestSetCompleteOpt(t *testing.T) {
// 	s := setSetup()
// 	s.setAccept("endpoint=travis-ci.com")
// 	assert.True(t, s.complete())
// }

// Run

var path = "/tmp/trsh.test/config.json"

func setup() {
	os.Remove(path)
}

func run(str string) ops.Op {
	s, _ := setAccept(str)
	s.Run(nil, nil, nil)
	return s
}

func read() string {
	b, _ := ioutil.ReadFile(path)
	return string(b)
}

func TestSetRunOptNameAndEqual(t *testing.T) {
	setup()
	run(":set endpoint=")
	assert.Equal(t, "{}", read())
}

func TestSetRunOpt(t *testing.T) {
	setup()
	run(":set endpoint=travis-ci.com")
	assert.Equal(t, "{\n  \"endpoint\": \"travis-ci.com\"\n}", read())
}
