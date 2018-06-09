package ops_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/svenfuchs/trsh/api"
	"github.com/svenfuchs/trsh/cmd"
	"github.com/svenfuchs/trsh/conf"
	"github.com/svenfuchs/trsh/jq"
	"github.com/svenfuchs/trsh/ops"
	"github.com/svenfuchs/trsh/pipe"
	s "github.com/svenfuchs/trsh/spec"
	p "path"
	"runtime"
	"testing"
)

var _, file, _, _ = runtime.Caller(0)
var path = p.Join(p.Dir(file), "../test/manifest.json")
var spec = s.Load(nil, path)

// Accept

func accept(str string) (*ops.Ops, bool) {
	api.Init(spec)
	jq.Init(spec)
	pipe.Init()
	cmd.Init()
	ops.Init(spec)
	o := ops.New(conf.New("/tmp/trsh.test/config.json"))
	return o, o.Accept(str)
}

func TestOpsAcceptEmpty(t *testing.T) {
	o, r := accept("")
	assert.False(t, r)
	assert.Equal(t, "", o.Str())
}

func TestOpsAcceptUnknownResource(t *testing.T) {
	o, r := accept("fux")
	assert.False(t, r)
	assert.Equal(t, "", o.Str())
}

func TestOpsAcceptPartOfResource(t *testing.T) {
	o, r := accept("us")
	assert.True(t, r)
	assert.Equal(t, "user", o.Str())
}

func TestOpsAcceptResource(t *testing.T) {
	o, r := accept("user")
	assert.True(t, r)
	assert.Equal(t, "user", o.Str())
}

func TestOpsAcceptResourceAndPipeToDot(t *testing.T) {
	o, r := accept("user | .")
	assert.False(t, r)
	assert.Equal(t, "", o.Str())
}

func TestOpsAcceptApiAndSpace(t *testing.T) {
	o, r := accept("user find id=1 ")
	assert.True(t, r)
	assert.Equal(t, "user find id=1", o.Str())
}

func TestOpsAcceptApiAndSpaceAndPipe(t *testing.T) {
	o, r := accept("user find id=1 |")
	assert.True(t, r)
	assert.Equal(t, "user find id=1 | ", o.Str())
}

func TestOpsAcceptApiAndSpaceAndPipeAndSpace(t *testing.T) {
	o, r := accept("user find id=1 | ")
	assert.True(t, r)
	assert.Equal(t, "user find id=1 | ", o.Str())
}

func TestOpsAcceptApiAndPipeToDotAttribute(t *testing.T) {
	o, r := accept("user find id=1 | .name")
	assert.True(t, r)
	assert.Equal(t, "user find id=1 | .name", o.Str())
}

func TestOpsAcceptApiAndPipeToBrackets(t *testing.T) {
	o, r := accept("builds find repository.id=1 | .builds [] .number")
	assert.True(t, r)
	assert.Equal(t, "builds find repository.id=1 | .builds [] .number", o.Str())
}

func TestOpsAccepCmd(t *testing.T) {
	o, r := accept(":set endpoint=api.travis-ci.com")
	assert.True(t, r)
	assert.Equal(t, ":set endpoint=api.travis-ci.com", o.Str())
}

// Hint

// func TestOpsHintEmpty(t *testing.T) {
// 	o, r := accept("")
// 	assert.False(t, r)
// 	assert.Equal(t, "", o.Str())
// }
//
