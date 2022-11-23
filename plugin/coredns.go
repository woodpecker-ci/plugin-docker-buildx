package plugin

import (
	"io"
	"net"
	"os"
	"os/exec"
)

func (p Plugin) startCoredns() {
	cmd := exec.Command("coredns", "-conf", "/etc/coredns/Corefile")
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

func getContainerIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.IsLoopback() {
					continue
				}
				if v.IP.To4() != nil {
					return v.IP.String(), nil
				}
			}
		}
	}
	return "", nil
}
