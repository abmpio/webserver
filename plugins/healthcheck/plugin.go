//go:build !windows
// +build !windows

package main

import (
	_ "github.com/abmpio/webserver/starter/healthcheck"
)

type Bootstrap struct {
}

func newBootstrap() Bootstrap {
	b := Bootstrap{}
	return b
}

func (b Bootstrap) BootstrapPlugin() (err error) {
	return nil
}

var PluginBootstrap = newBootstrap()
