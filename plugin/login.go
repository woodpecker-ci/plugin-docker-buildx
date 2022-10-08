package plugin

import (
	"fmt"
	"os/exec"
)

// login to the registrys
func (p *Plugin) Login() error {
	registrys := make(map[string]bool)
	for _, login := range p.settings.Logins {
		if !registrys[login.Registry] && !login.anonymous() {
			// only log into a registry once
			registrys[login.Registry] = true
			cmd := commandLogin(login)
			err := cmd.Run()
			if err != nil {
				return fmt.Errorf("error authenticating: %s", err)
			}
		}
	}
	return nil
}

// helper function to create the docker login command.
func commandLogin(login Login) *exec.Cmd {
	if login.Email != "" {
		return commandLoginEmail(login)
	}
	return exec.Command(
		dockerExe, "login",
		"-u", login.Username,
		"-p", login.Password,
		login.Registry,
	)
}
