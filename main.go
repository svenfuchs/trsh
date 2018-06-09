package main

import (
	"bufio"
	"github.com/svenfuchs/trsh/api"
	"github.com/svenfuchs/trsh/cmd"
	"github.com/svenfuchs/trsh/conf"
	"github.com/svenfuchs/trsh/jq"
	"github.com/svenfuchs/trsh/ops"
	"github.com/svenfuchs/trsh/pipe"
	// "github.com/svenfuchs/trsh/http"
	"github.com/svenfuchs/trsh/prompt"
	// "github.com/svenfuchs/trsh/spec"
	"os"
	"os/user"
	"path"
	"strings"
)

func main() {
	api.Init()
	cmd.Init()
	jq.Init()
	pipe.Init()

	u, _ := user.Current()
	p := path.Join(u.HomeDir, ".trsh/config.json")
	c := conf.New(p)

	if len(os.Args) > 1 {
		run(c, strings.Join(args(), " "))
	} else if isStdin() {
		stdin(func(str string) { run(c, str) })
	} else {
		prompt.New(c, "travis ~ ").Run()
	}
}

func run(conf *conf.Conf, str string) {
	o := ops.New(conf)
	if o.Accept(str) {
		o.Run(os.Stdin, os.Stdout, os.Stderr)
	}
}

func args() []string {
	a := os.Args[1:]
	if a[0] == "--" {
		a = a[1:]
	}
	return a
}

func stdin(f func(string)) {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		f(s.Text())
	}
}

func isStdin() bool {
	s, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	return s.Size() > 0
}
