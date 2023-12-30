package plugin

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"codeberg.org/woodpecker-plugins/plugin-docker-buildx/utils"
)

// Daemon defines Docker daemon parameters.
type Daemon struct {
	Registry          string          // Docker registry
	Mirror            string          // Docker registry mirror
	Insecure          bool            // Docker daemon enable insecure registries
	StorageDriver     string          // Docker daemon storage driver
	StoragePath       string          // Docker daemon storage path
	Disabled          bool            // DOcker daemon is disabled (already running)
	Debug             bool            // Docker daemon started in debug mode
	Bip               string          // Docker daemon network bridge IP address
	DNS               cli.StringSlice // Docker daemon dns server
	DNSSearch         cli.StringSlice // Docker daemon dns search domain
	MTU               string          // Docker daemon mtu setting
	IPv6              bool            // Docker daemon IPv6 networking
	Experimental      bool            // Docker daemon enable experimental mode
	BuildkitConfig    string          // Docker buildkit config
	BuildkitDriverOpt cli.StringSlice // Docker buildkit driveropt args
}

// Login defines Docker login parameters.
type Login struct {
	// Generic
	Registry string // Docker registry address
	Username string // Docker registry username
	Password string // Docker registry password
	Email    string // Docker registry email
	Config   string // Docker Auth Config

	// ECR
	Aws_access_key_id     string `json:"aws_access_key_id"`     // AWS access key id
	Aws_secret_access_key string `json:"aws_secret_access_key"` // AWS secret access key
	Aws_region            string `json:"aws_region"`            // AWS region
}

// Build defines Docker build parameters.
type Build struct {
	Remote          string          // Git remote URL
	Ref             string          // Git commit ref
	Branch          string          // Git repository branch
	Dockerfile      string          // Docker build Dockerfile
	Context         string          // Docker build context
	TagsAuto        bool            // Docker build auto tag
	TagsDefaultName string          // Docker build auto tag name override
	TagsSuffix      string          // Docker build tags with suffix
	Tags            cli.StringSlice // Docker build tags
	TagsFile        string          // Docker build tags read from an file
	LabelsAuto      bool            // Docker build auto labels
	Labels          cli.StringSlice // Docker build labels
	Platforms       cli.StringSlice // Docker build target platforms
	Args            cli.StringSlice // Docker build args
	ArgsEnv         cli.StringSlice // Docker build args from env
	Target          string          // Docker build target
	Output          string          // Docker build output
	Pull            bool            // Docker build pull
	CacheFrom       cli.StringSlice // Docker build cache-from
	CacheTo         cli.StringSlice // Docker build cache-to
	Compress        bool            // Docker build compress
	Repo            cli.StringSlice // Docker build repository
	NoCache         bool            // Docker build no-cache
	AddHost         cli.StringSlice // Docker build add-host
	Quiet           bool            // Docker build quiet
}

// Settings for the Plugin.
type Settings struct {
	// ECR
	AwsRegion           string `json:"aws_region"`            // AWS region
	EcrScanOnPush       bool   `json:"ecr_scan_on_push"`      // ECR scan on push
	EcrRepositoryPolicy string `json:"ecr_repository_policy"` // ECR repository policy
	EcrLifecyclePolicy  string `json:"ecr_lifecycle_policy"`  // ECR lifecycle policy
	EcrCreateRepository bool   `json:"ecr_create_repository"` // ECR create repository
	AwsAccessKeyId      string `json:"aws_access_key_id"`     // AWS access key id
	AwsSecretAccessKey  string `json:"aws_secret_access_key"` // AWS secret access key

	// Generic
	Daemon       Daemon
	Logins       []Login
	LoginsRaw    string
	DefaultLogin Login
	Build        Build
	Dryrun       bool
	Cleanup      bool
}

func (l Login) anonymous() bool {
	return l.Username == "" || l.Password == ""
}

// Init initialise plugin settings
func (p *Plugin) InitSettings() error {
	if p.settings.LoginsRaw != "" {
		if err := json.Unmarshal([]byte(p.settings.LoginsRaw), &p.settings.Logins); err != nil {
			return fmt.Errorf("could not unmarshal logins: %v", err)
		}
	}

	p.settings.Build.Branch = p.pipeline.Repo.Branch
	p.settings.Build.Ref = p.pipeline.Commit.Ref

	// check if any Login struct contains AWS credentials
	for _, login := range p.settings.Logins {
		if strings.Contains(login.Registry, "amazonaws.com") {
			p.EcrInit()
		}
	}

	if p.settings.AwsAccessKeyId != "" && p.settings.AwsSecretAccessKey != "" {
		p.EcrInit()
	}

	if len(p.settings.Logins) == 0 {
		p.settings.Logins = []Login{p.settings.DefaultLogin}
	} else if !p.settings.DefaultLogin.anonymous() {
		p.settings.Logins = prepend(p.settings.Logins, p.settings.DefaultLogin)
	}

	p.settings.Daemon.Registry = p.settings.Logins[0].Registry

	return nil
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	if err := p.InitSettings(); err != nil {
		return err
	}

	if !isSingleTag(p.settings.Build.TagsDefaultName) {
		return fmt.Errorf("'%s' is not a valid, single tag", p.settings.Build.TagsDefaultName)
	}

	// beside the default login all other logins need to set a username and password
	for _, l := range p.settings.Logins[1:] {
		if l.anonymous() {
			return fmt.Errorf("beside the default login all other logins need to set a username and password")
		}
	}

	// overload tags flag with tags.file if set
	if p.settings.Build.TagsFile != "" {
		tagsFile, err := os.ReadFile(p.settings.Build.TagsFile)
		if err != nil {
			return fmt.Errorf("could not read tags file: %w", err)
		}

		// split file content into slice of tags
		tagsFileList := strings.Split(strings.TrimSpace(string(tagsFile)), "\n")
		// trim space of each tag
		tagsFileList = utils.Map(tagsFileList, func(s string) string { return strings.TrimSpace(s) })

		// finally overwrite
		p.settings.Build.Tags = *cli.NewStringSlice(tagsFileList...)
	}

	if p.settings.Build.TagsAuto {
		// we only generate tags on default branch or an tag event
		if UseDefaultTag(
			p.settings.Build.Ref,
			p.settings.Build.Branch,
		) {
			tag, err := DefaultTagSuffix(
				p.settings.Build.Ref,
				p.settings.Build.TagsDefaultName,
				p.settings.Build.TagsSuffix,
			)
			if err != nil {
				logrus.Printf("cannot build docker image for %s, invalid semantic version", p.settings.Build.Ref)
				return err
			}

			// include user supplied tags
			tag = append(tag, p.sanitizedUserTags()...)

			p.settings.Build.Tags = *cli.NewStringSlice(tag...)
		} else {
			logrus.Printf("skipping automated docker build for %s", p.settings.Build.Ref)
			return nil
		}
	} else {
		p.settings.Build.Tags = *cli.NewStringSlice(p.sanitizedUserTags()...)
	}

	if p.settings.Build.LabelsAuto {
		p.settings.Build.Labels = *cli.NewStringSlice(p.Labels()...)
	}

	return nil
}

func (p *Plugin) sanitizedUserTags() []string {
	// ignore empty tags
	var tags []string
	for _, t := range p.settings.Build.Tags.Value() {
		t = strings.TrimSpace(t)
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}

func (p *Plugin) writeBuildkitConfig() error {
	// no buildkit config, automatically generate buildkit configuration to use a custom CA certificate for each registry
	if p.settings.Daemon.BuildkitConfig == "" && p.settings.Daemon.Registry != "" {
		for _, login := range p.settings.Logins {
			if registry := login.Registry; registry != "" {
				u, err := url.Parse(registry)
				if err != nil {
					return fmt.Errorf("could not parse registry address: %s: %v", registry, err)
				}
				if u.Host != "" {
					registry = u.Host
				}

				caPath := fmt.Sprintf("/etc/docker/certs.d/%s/ca.crt", registry)
				ca, err := os.Open(caPath)
				if err != nil && !os.IsNotExist(err) {
					logrus.Warnf("error reading %s: %v", caPath, err)
				} else if err == nil {
					ca.Close()
					p.settings.Daemon.BuildkitConfig += fmt.Sprintf(buildkitConfigTemplate, registry, caPath)
				}
			}
		}
	}

	// save buildkit config as described
	if p.settings.Daemon.BuildkitConfig != "" {
		err := os.WriteFile(buildkitConfig, []byte(p.settings.Daemon.BuildkitConfig), 0o600)
		if err != nil {
			return fmt.Errorf("error writing buildkit.toml: %s", err)
		}
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	// start the Docker daemon server
	if !p.settings.Daemon.Disabled {
		// If no custom DNS value set start internal DNS server
		if len(p.settings.Daemon.DNS.Value()) == 0 {
			ip, err := getContainerIP()
			if err != nil {
				logrus.Warnf("error detecting IP address: %v", err)
			} else if ip != "" {
				p.startCoredns()
				p.settings.Daemon.DNS.Set(ip)
			}
		}

		p.startDaemon()
	}

	// poll the docker daemon until it is started. This ensures the daemon is
	// ready to accept connections before we proceed.
	for i := 0; i < 15; i++ {
		cmd := commandInfo()
		err := cmd.Run()
		if err == nil {
			break
		}
		time.Sleep(time.Second * 1)
	}

	// Create Auth Config File
	if p.settings.Logins[0].Config != "" {
		os.MkdirAll(dockerHome, 0o600)

		path := filepath.Join(dockerHome, "config.json")
		err := os.WriteFile(path, []byte(p.settings.Logins[0].Config), 0o600)
		if err != nil {
			return fmt.Errorf("error writing config.json: %s", err)
		}
	}

	// login to the Docker registry
	if err := p.Login(); err != nil {
		return err
	}

	if err := p.writeBuildkitConfig(); err != nil {
		return err
	}

	switch {
	case p.settings.Logins[0].Password != "":
		fmt.Println("Detected registry credentials")
	case p.settings.Logins[0].Config != "":
		fmt.Println("Detected registry credentials file")
	default:
		fmt.Println("Registry credentials or Docker config not provided. Guest mode enabled.")
	}

	// add proxy build args
	addProxyBuildArgs(&p.settings.Build)

	var cmds []*exec.Cmd
	cmds = append(cmds, commandVersion()) // docker version
	cmds = append(cmds, commandInfo())    // docker info
	cmds = append(cmds, commandBuilder(p.settings.Daemon))
	cmds = append(cmds, commandBuildx())
	cmds = append(cmds, commandBuild(p.settings.Build, p.settings.Dryrun)) // docker build

	// execute all commands in batch mode.
	for _, cmd := range cmds {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		trace(cmd)

		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func prepend[Type any](slice []Type, elems ...Type) []Type {
	return append(elems, slice...)
}
