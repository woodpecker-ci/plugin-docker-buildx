package plugin

import (
	"io"
	"os"
)

const (
	dockerExe      = "/usr/local/bin/docker"
	dockerdExe     = "/usr/local/bin/dockerd"
	dockerHome     = "/root/.docker/"
	buildkitConfig = "/tmp/buildkit.toml"
)

func (p Plugin) startDaemon() {
	cmd := commandDaemon(p.settings.Daemon)
	if p.settings.Daemon.Debug {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	}
	go func() {
		trace(cmd)
		cmd.Run()
	}()
}
