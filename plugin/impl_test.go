package plugin

import (
	"testing"

	"codeberg.org/6543/go-yaml2json"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

var defaultSettings = Settings{
	Daemon: Daemon{
		StoragePath: "/var/lib/docker",
	},
	Build: Build{
		Context:         ".",
		Tags:            *cli.NewStringSlice("latest"),
		TagsDefaultName: "latest",
		LabelsAuto:      true,
		Pull:            true,
	},
	DefaultLogin: Login{
		Registry: "https://index.docker.io/v1/",
	},
	LoginsRaw: "[]",
	Cleanup:   true,
}

func TestDefaultLogin(t *testing.T) {
	s := defaultSettings
	assert.NoError(t, newSettingsOnly(&s).Validate())
	if assert.Len(t, s.Logins, 1) {
		assert.EqualValues(t, defaultSettings.DefaultLogin.Registry, s.Logins[0].Registry)
	}

	// only use login to auth to registrys
	loginsRaw, err := yaml2json.Convert([]byte(`
- registry: https://index.docker.io/v1/
  username: docker_username
  password: docker_password
- registry: https://codeberg.org
  username: cb_username
  password: cb_password`))
	assert.NoError(t, err)
	s.LoginsRaw = string(loginsRaw)
	assert.NoError(t, newSettingsOnly(&s).Validate())
	if assert.Len(t, s.Logins, 2) {
		assert.EqualValues(t, defaultSettings.DefaultLogin.Registry, s.Logins[0].Registry)
	}

	// mixed login settings ('logins' and 'username', 'password' are used)
	s = defaultSettings
	loginsRaw, err = yaml2json.Convert([]byte(`
- registry: https://codeberg.org
  username: cb_username
  password: cb_password`))
	assert.NoError(t, err)
	s.LoginsRaw = string(loginsRaw)
	s.DefaultLogin.Username = "docker_username"
	s.DefaultLogin.Password = "docker_password"
	assert.NoError(t, newSettingsOnly(&s).Validate())
	if assert.Len(t, s.Logins, 2) {
		assert.EqualValues(t, defaultSettings.DefaultLogin.Registry, s.Logins[0].Registry)
	}

	// ignore default registry
	s = defaultSettings
	loginsRaw, err = yaml2json.Convert([]byte(`
- registry: https://codeberg.org
  username: cb_username
  password: cb_password`))
	assert.NoError(t, err)
	s.LoginsRaw = string(loginsRaw)
	assert.NoError(t, newSettingsOnly(&s).Validate())
	if assert.Len(t, s.Logins, 1) {
		assert.EqualValues(t, "https://codeberg.org", s.Logins[0].Registry)
	}
}
