package cmd_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/svenfuchs/trsh/cmd"
	"github.com/svenfuchs/trsh/conf"
	"github.com/svenfuchs/trsh/ops"
	"testing"
)

// Accept

func cmdAccept(str string) (ops.Op, bool) {
	cmd.Init()
	o := ops.New(conf.New("/tmp/trsh.test/config.json"))
	c := cmd.New(o)
	ok, _ := c.Accept(str)
	return c, ok
}

func TestCmdAcceptColon(t *testing.T) {
	c, ok := cmdAccept(":")
	assert.True(t, ok)
	assert.Equal(t, ":", c.Str())
}

func TestCmdAcceptColonUnknownCmdName(t *testing.T) {
	c, ok := cmdAccept(":fux")
	assert.False(t, ok)
	assert.Equal(t, ":", c.Str())
}

func TestCmdAcceptColonPartOfCmdName(t *testing.T) {
	c, ok := cmdAccept(":s")
	assert.True(t, ok)
	assert.Equal(t, ":", c.Str())
}

func TestCmdAcceptColonCmdName(t *testing.T) {
	c, ok := cmdAccept(":set")
	assert.True(t, ok)
	assert.Equal(t, ":set", c.Str())
}

// Hint

func cmdHint(str string) string {
	c, _ := cmdAccept(str)
	return c.Hint(str)
}

// func TestCmdHintEmpty(t *testing.T) {
// 	h := cmdHint("")
// 	assert.Equal(t, "active for_owner organization.login=", h)
// }

// func TestCmdHintUnknownOptName(t *testing.T) {
// 	c, ok := cmdAccept("wat")
// 	assert.False(t, c.Complete())
// }
//
// func TestCmdHintPartOfOptName(t *testing.T) {
// 	c, ok := cmdAccept("end")
// 	assert.False(t, c.Complete())
// }
//
// func TestCmdHintOptName(t *testing.T) {
// 	c, ok := cmdAccept("endpoint")
// 	assert.False(t, c.Complete())
// }
//
// func TestCmdHintOptNameAndEqual(t *testing.T) {
// 	c, ok := cmdAccept("endpoint=")
// 	assert.False(t, c.Complete())
// }
//
// func TestCmdHintOpt(t *testing.T) {
// 	c, ok := cmdAccept("endpoint=travis-ci.com")
// 	assert.True(t, c.Complete())
// }

// // Complete
//
// func TestCmdCompleteEmpty(t *testing.T) {
// 	s := setCmdup()
// 	s.accept("")
// 	assert.False(t, s.complete())
// }
//
// func TestCmdCompleteUnknownOptName(t *testing.T) {
// 	s := setCmdup()
// 	s.accept("wat")
// 	assert.False(t, s.complete())
// }
//
// func TestCmdCompletePartOfOptName(t *testing.T) {
// 	s := setCmdup()
// 	s.accept("end")
// 	assert.False(t, s.complete())
// }
//
// func TestCmdCompleteOptName(t *testing.T) {
// 	s := setCmdup()
// 	s.accept("endpoint")
// 	assert.False(t, s.complete())
// }
//
// func TestCmdCompleteOptNameAndEqual(t *testing.T) {
// 	s := setCmdup()
// 	s.accept("endpoint=")
// 	assert.False(t, s.complete())
// }
//
// func TestCmdCompleteOpt(t *testing.T) {
// 	s := setCmdup()
// 	s.accept("endpoint=travis-ci.com")
// 	assert.True(t, s.complete())
// }
