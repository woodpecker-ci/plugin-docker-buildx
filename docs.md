---
name: Docker Buildx
icon: https://woodpecker-ci.org/img/logo.svg
description: plugin to build multiarch Docker images with buildx
authors: Woodpecker Authors
tags: [docker, image, container, build]
containerImage: woodpeckerci/plugin-docker-buildx
containerImageUrl: https://hub.docker.com/r/woodpeckerci/plugin-docker-buildx
url: https://codeberg.org/woodpecker-plugins/docker-buildx
---

Woodpecker CI plugin to build multiarch Docker images with buildx. This plugin is a fork of [thegeeklab/drone-docker-buildx](https://github.com/thegeeklab/drone-docker-buildx/) which itself is a fork of [drone-plugins/drone-docker](https://github.com/drone-plugins/drone-docker).

## Features

- Build without push
- Use custom registries
- Build based on existing tags when needed
- Push to multiple registries/repos

It will automatically generate buildkit configuration to use custom CA certificate if following conditions are met:

- Setting `buildkit_config` is not set
- Custom `registry`/`logins` value is provided
- File exists `/etc/docker/certs.d/<registry-value>/ca.crt`

> NB! To mount custom CA you can use Woodpecker CI runner configuration environment `WOODPECKER_BACKEND_DOCKER_VOLUMES` with value `/etc/ssl/certs:/etc/ssl/certs:ro,/etc/docker/certs.d:/etc/docker/certs.d:ro`. And have created file `/etc/docker/certs.d/<registry-value>/ca.crt` with CA certificate on runner server host.

## Settings

| Settings Name             | Default           | Description
| --------------------------| ----------------- | --------------------------------------------
| `dry-run`                 | `false`           | disables docker push
| `repo`                    | *none*            | sets repository name for the image (can be a list)
| `username`                | *none*            | sets username to authenticates with
| `password`                | *none*            | sets password / token to authenticates with
| `email`                   | *none*            | sets email address to authenticates with
| `registry`                | `https://index.docker.io/v1/` | sets docker registry to authenticate with
| `dockerfile`              | `Dockerfile`      | sets dockerfile to use for the image build
| `tag`/`tags`              | *none*            | sets repository tags to use for the image
| `platforms`               | *none*            | sets target platform for build

## auto_tag

If set to true, it will use the `default_tag` ("latest") on tag event or default branch.
If it's a tag event it will also assume sem versioning and add tags accordingly (`x`, `x.x` and `x.x.x`).
If it's not a tag event, and no default branch, automated tags are skipped.

## Examples

```yaml
  publish-next-agent:
    image: woodpeckerci/plugin-docker-buildx
    secrets: [docker_username, docker_password]
    settings:
      repo: woodpeckerci/woodpecker-agent
      dockerfile: docker/Dockerfile.agent.multiarch
      platforms: windows/amd64,darwin/amd64,darwin/arm64,freebsd/amd64,linux/amd64,linux/arm64/v8
      tag: next
    when:
      branch: ${CI_REPO_DEFAULT_BRANCH}
      event: push
```

```yaml
  publish:
    image: woodpeckerci/plugin-docker-buildx
    settings:
      platforms: linux/386,linux/amd64,linux/arm/v6,linux/arm64/v8,linux/ppc64le,linux/riscv64,linux/s390x
      repo: codeberg.org/${CI_REPO_OWNER}/hello
      registry: codeberg.org
      tags: latest
      username: ${CI_REPO_OWNER}
      password:
        from_secret: cb_token
```

```yaml
  docker-build:
    image: woodpeckerci/plugin-docker-buildx
    settings:
      repo: codeberg.org/${CI_REPO_OWNER}/hello
      registry: codeberg.org
      dry_run: true
      output: type=oci,dest=${CI_REPO_OWNER}-hello.tar
```

## Advanced Settings

| Settings Name             | Default           | Description
| --------------------------| ----------------- | --------------------------------------------
| `mirror`                  | *none*            | sets a registry mirror to pull images
| `storage_driver`          | *none*            | sets the docker daemon storage driver
| `storage_path`            | `/var/lib/docker` | sets the docker daemon storage path
| `bip`                     | *none*            | allows the docker daemon to bride ip address
| `mtu`                     | *none*            | sets docker daemon custom mtu setting
| `custom_dns`              | *none*            | sets custom docker daemon dns server
| `custom_dns_search`       | *none*            | sets custom docker daemon dns search domain
| `insecure`                | `false`           | allows the docker daemon to use insecure registries
| `ipv6`                    | `false`           | enables docker daemon IPv6 support
| `experimental`            | `false`           | enables docker daemon experimental mode
| `debug`                   | `false`           | enables verbose debug mode for the docker daemon
| `daemon_off`              | `false`           | disables the startup of the docker daemon
| `buildkit_config`         | *none*            | sets content of the docker [buildkit TOML config](https://github.com/moby/buildkit/blob/master/docs/buildkitd.toml.md)
| `tags_file`               | *none*            | overwrites `tags` option with values find in specified file
| `context`                 | `.`               | sets the path of the build context to use
| `auto_tag`                | `false`           | generates tag names automatically based on git branch and git tag, tags supplied via `tags` are additionally added to the auto_tags without suffix
| `default_suffix"`/`auto_tag_suffix`| *none*   | generates tag names with the given suffix
| `default_tag`             | `latest`          | overrides the default tag name used when generating with `auto_tag` enabled
| `label`/`labels`          | *none*            | sets labels to use for the image in format `<name>=<value>`
| `default_labels`/`auto_labels` | `true`       | sets docker image labels based on git information
| `build_args`              | *none*            | sets custom build arguments for the build
| `build_args_from_env`     | *none*            | forwards environment variables as custom arguments to the build
| `quiet`                   | `false`           | enables suppression of the build output
| `target`                  | *none*            | sets the build target to use
| `cache_from`              | *none*            | sets images to consider as cache sources
| `pull_image`              | `true`            | enforces to pull base image at build time
| `compress`                | `false`           | enables compression of the build context using gzip
| `config`                  | *none*            | sets content of the docker daemon json config
| `purge`                   | `true`            | enables cleanup of the docker environment at the end of a build
| `no_cache`                | `false`           | disables the usage of cached intermediate containers
| `add_host`                | *none*            | sets additional host:ip mapping
| `output`                  | *none*            | sets build output in format `type=<type>[,<key>=<value>]`
| `logins`                  | *none*            | option to log into multiple registries

## Multi registry push example

Only supported with `woodpecker >= 1.0.0` (next-da997fa3).

```yaml
settings:
  repo: a6543/tmp,codeberg.org/6543/tmp
  tag: demo
  logins:
    - registry: https://index.docker.io/v1/
      username: a6543
      password:
        from_secret: docker_token
    - registry: https://codeberg.org
      username: "6543"
      password:
        from_secret: cb_token
```
