package main

import (
	"github.com/zinvapel/warden-core"
	"github.com/zinvapel/warden-core/pkg/environment"
)

var env *environment.Environment

func init() {
	defaultEnv, err := warden.DefaultEnvironment()

	if err != nil {
		panic(err)
	}

	env = defaultEnv
}

func main() {
	warden.NewWarden().WithEnv(env).Run()
}
