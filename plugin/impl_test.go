package plugin

import (
	"fmt"
	"os"
	"testing"

	"codeberg.org/6543/go-yaml2json"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

var defaultTestSettings = Settings{
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
	LoginsRaw:       "[]",
	Cleanup:         true,
	CustomCertStore: "/etc/docker/certs.d/",
}

func TestDefaultLogin(t *testing.T) {
	s := defaultTestSettings
	assert.NoError(t, newSettingsOnly(&s).Validate())
	if assert.Len(t, s.Logins, 1) {
		assert.EqualValues(t, defaultTestSettings.DefaultLogin.Registry, s.Logins[0].Registry)
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
		assert.EqualValues(t, defaultTestSettings.DefaultLogin.Registry, s.Logins[0].Registry)
	}

	// mixed login settings ('logins' and 'username', 'password' are used)
	s = defaultTestSettings
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
		assert.EqualValues(t, defaultTestSettings.DefaultLogin.Registry, s.Logins[0].Registry)
	}

	// ignore default registry
	s = defaultTestSettings
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

func TestWriteBuildkitConfig(t *testing.T) {
	settings := defaultTestSettings
	assert.NoError(t, newSettingsOnly(&settings).Validate())
	assert.EqualValues(t, "", settings.Daemon.BuildkitConfig)

	settings = defaultTestSettings
	settings.Daemon.BuildkitDebug = true
	assert.NoError(t, newSettingsOnly(&settings).Validate())
	assert.EqualValues(t, "debug = true\n", settings.Daemon.BuildkitConfig)

	settings = defaultTestSettings
	settings.Daemon.Mirror = "mirror.example.com"
	assert.NoError(t, newSettingsOnly(&settings).Validate())
	assert.EqualValues(t, "[registry]\n[registry.'docker.io']\nmirrors = ['mirror.example.com']\n", settings.Daemon.BuildkitConfig)

	settings = defaultTestSettings
	settings.DefaultLogin.Registry = "codeberg.org"
	tmpDir, err := os.MkdirTemp("", "go-test-*")
	assert.NoError(t, err)
	settings.CustomCertStore = tmpDir
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, os.Mkdir(tmpDir+"/codeberg.org", os.ModePerm))
	caFile, err := os.Create(tmpDir + "/codeberg.org/" + "ca.crt")
	assert.NoError(t, err)
	assert.NoError(t, caFile.Close())

	assert.NoError(t, newSettingsOnly(&settings).Validate())
	assert.EqualValues(t, fmt.Sprintf("[registry]\n[registry.'codeberg.org']\nca = ['%s/codeberg.org/ca.crt']\n", tmpDir), settings.Daemon.BuildkitConfig)
}
