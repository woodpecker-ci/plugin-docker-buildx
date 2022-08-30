---
name: Docker Buildx
icon: https://woodpecker-ci.org/img/logo.svg
description: plugin to build multiarch Docker images with buildx
---

Woodpecker CI plugin to build multiarch Docker images with buildx. This plugin is a fork of [thegeeklab/drone-docker-buildx](https://github.com/thegeeklab/drone-docker-buildx/) which itself is a fork of [drone-plugins/drone-docker](https://github.com/drone-plugins/drone-docker). You can find the full documentation at You can find the full documentation at [woodpecker-plugins.codeberg.page](https://woodpecker-plugins.codeberg.page/plugins/drone-docker-buildx).

## Features

- Build without push
- Use custom registries
- Build based on existing tags when needed.

## Settings

| Settings Name             | Default           | Description
| --------------------------| ----------------- | --------------------------------------------
| `dry-run`                 | `false`           | disables docker push
| `repo`                    | *none*            | sets repository name for the image
| `username`                | *none*            | sets username to authenticates with
| `password`                | *none*            | sets password / token to authenticates with
| `email`                   | *none*            | sets email address to authenticates with
| `registry`                | `https://index.docker.io/v1/` | sets docker registry to authenticate with
| `dockerfile`              | `Dockerfile`      | sets dockerfile to use for the image build
| `tag`/`tags`              | @".tags"          | sets repository tags to use for the image
| `platforms`               | *none*           | sets target platform for build

## Examples

```yml
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

```yml
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
| `buildkit_config`         | *none*            | sets content of the docker buildkit json config
| `context`                 | `.`               | sets the path of the build context to use
| `default_tags`/`auto_tag` | `false`           | generates tag names automatically based on git branch and git tag
| `default_suffix"`/`auto_tag_suffix`| *none*   | generates tag names with the given suffix
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
