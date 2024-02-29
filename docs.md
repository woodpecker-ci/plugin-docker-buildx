---
name: Docker Buildx
icon: https://codeberg.org/woodpecker-plugins/docker-buildx/raw/branch/main/docker.svg
description: plugin to build multiarch Docker images with buildx
author: Woodpecker Authors
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

| Settings Name           | Default                       | Description                                        |
| ----------------------- | ----------------------------- | -------------------------------------------------- |
| `dry-run`               | `false`                       | disables docker push                               |
| `repo`                  | _none_                        | sets repository name for the image (can be a list) |
| `username`              | _none_                        | sets username to authenticates with                |
| `password`              | _none_                        | sets password / token to authenticates with        |
| `aws_access_key_id`     | _none_                        | sets AWS_ACCESS_KEY_ID for AWS ECR auth            |
| `aws_secret_access_key` | _none_                        | sets AWS_SECRET_ACCESS_KEY for AWS ECR auth        |
| `aws_region`            | `us-east-1`                   | sets AWS_DEFAULT_REGION for AWS ECR auth           |
| `password`              | _none_                        | sets password / token to authenticates with        |
| `email`                 | _none_                        | sets email address to authenticates with           |
| `registry`              | `https://index.docker.io/v1/` | sets docker registry to authenticate with          |
| `dockerfile`            | `Dockerfile`                  | sets dockerfile to use for the image build         |
| `tag`/`tags`            | _none_                        | sets repository tags to use for the image          |
| `platforms`             | _none_                        | sets target platform for build                     |
| `provenance`            | _none_                        | sets provenance for build                          |

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
    dry-run: true
    output: type=oci,dest=${CI_REPO_OWNER}-hello.tar
```

## Advanced Settings

| Settings Name                       | Default           | Description                                                                                                                                       |
| ----------------------------------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| `mirror`                            | _none_            | sets a registry mirror to pull images                                                                                                             |
| `storage_driver`                    | _none_            | sets the docker daemon storage driver                                                                                                             |
| `storage_path`                      | `/var/lib/docker` | sets the docker daemon storage path                                                                                                               |
| `bip`                               | _none_            | allows the docker daemon to bride ip address                                                                                                      |
| `mtu`                               | _none_            | sets docker daemon custom mtu setting                                                                                                             |
| `custom_dns`                        | _none_            | sets custom docker daemon dns server                                                                                                              |
| `custom_dns_search`                 | _none_            | sets custom docker daemon dns search domain                                                                                                       |
| `insecure`                          | `false`           | allows the docker daemon to use insecure registries                                                                                               |
| `ipv6`                              | `false`           | enables docker daemon IPv6 support                                                                                                                |
| `experimental`                      | `false`           | enables docker daemon experimental mode                                                                                                           |
| `debug`                             | `false`           | enables verbose debug mode for the docker daemon                                                                                                  |
| `daemon_off`                        | `false`           | disables the startup of the docker daemon                                                                                                         |
| `buildkit_debug`                    | `false`           | enables debug output of buildkit                                                                                                                  |
| `buildkit_config`                   | _none_            | sets content of the docker[buildkit TOML config](https://github.com/moby/buildkit/blob/master/docs/buildkitd.toml.md)                             |
| `buildkit_driveropt`                | _none_            | adds one or multiple`--driver-opt` buildx arguments for the default buildkit builder instance                                                     |
| `tags_file`                         | _none_            | overrides the`tags` option with values in a file named `.tags`; multiple tags can be specified separated by a newline                             |
| `context`                           | `.`               | sets the path of the build context to use                                                                                                         |
| `auto_tag`                          | `false`           | generates tag names automatically based on git branch and git tag, tags supplied via`tags` are additionally added to the auto_tags without suffix |
| `default_suffix"`/`auto_tag_suffix` | _none_            | generates tag names with the given suffix                                                                                                         |
| `default_tag`                       | `latest`          | overrides the default tag name used when generating with`auto_tag` enabled                                                                        |
| `label`/`labels`                    | _none_            | sets labels to use for the image in format`<name>=<value>`                                                                                        |
| `default_labels`/`auto_labels`      | `true`            | sets docker image labels based on git information                                                                                                 |
| `build_args`                        | _none_            | sets custom build arguments for the build                                                                                                         |
| `build_args_from_env`               | _none_            | forwards environment variables as custom arguments to the build                                                                                   |
| `secrets`                           | _none_            | Sets the build secrets for the build                                                                                                              |
| `quiet`                             | `false`           | enables suppression of the build output                                                                                                           |
| `target`                            | _none_            | sets the build target to use                                                                                                                      |
| `cache_from`                        | _none_            | sets configuration for cache source                                                                                                               |
| `cache_to`                          | _none_            | sets configuration for cache export                                                                                                               |
| `cache_images`                      | _none_            | a list of images to use as cache.                                                                                                                 |
| `pull_image`                        | `true`            | enforces to pull base image at build time                                                                                                         |
| `compress`                          | `false`           | enables compression of the build context using gzip                                                                                               |
| `config`                            | _none_            | sets content of the docker daemon json config                                                                                                     |
| `purge`                             | `true`            | enables cleanup of the docker environment at the end of a build                                                                                   |
| `no_cache`                          | `false`           | disables the usage of cached intermediate containers                                                                                              |
| `add_host`                          | _none_            | sets additional host:ip mapping                                                                                                                   |
| `output`                            | _none_            | sets build output in format`type=<type>[,<key>=<value>]`                                                                                          |
| `logins`                            | _none_            | option to log into multiple registries                                                                                                            |
| `env_file`                          | _none_            | load env vars from specified file                                                                                                                 |
| `ecr_create_repository`             | `false`           | creates the ECR repository if it does not exist                                                                                                   |
| `ecr_lifecycle_policy`              | _none_            | AWS ECR lifecycle policy                                                                                                                          |
| `ecr_repository_policy`             | _none_            | AWS ECR repository policy                                                                                                                         |
| `ecr_scan_on_push`                  | _none_            | AWS: whether to enable image scanning on push                                                                                                     |

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
      mirrors:
        - "my-docker-mirror-host.local"
    - registry: https://codeberg.org
      username: "6543"
      password:
        from_secret: cb_token
    - registry: https://<account-id>.dkr.ecr.<region>.amazonaws.com
      aws_region: <region>
      aws_access_key_id:
        from_secret: aws_access_key_id
      aws_secret_access_key:
        from_secret: aws_secret_access_key
```

## Using `plugin-docker-buildx` behind a proxy

When performing a docker build behind a corporate proxy one needs to pass through the proxy settings to the plugin.

```yaml
variables:
  # proxy config
  - proxy_conf: &proxy_conf
      - http_proxy: "http://X.Y.Z.Z:3128"
      - https_proxy: "http://X.Y.Z.Z:3128"
      - no_proxy: ".my-subdomain.com"
  # deployment targets
  - &publish_repos "codeberg.org/test"
  # logins for deployment targets
  - publish_logins: &publish_logins
      - registry: https://codeberg.org
        username:
          from_secret: CODEBERG_USER
        password:
          from_secret: CODEBERG_TOKEN

steps:
  test:
    image: woodpeckerci/plugin-docker-buildx:2
    environment:
      # adding proxy in env for the plugin runtime itself.
      - <<: *proxy_conf
    privileged: true
    settings:
      dry-run: true
      repo: *publish_repos
      dockerfile: Dockerfile.multi
      platforms: linux/amd64
      auto_tag: true
      logins: *publish_logins
      # Adding custom dns server to lookup internal Docker Hub mirror.
      # custom_dns:
      #   - 192.168.55.31
      #   - 192.168.55.32
      # Adding an optional Docker Hub mirror for the nested dockerd.
      # mirror: https://my-mirror.example.com
      build_args:
        # passthrough proxy config to the build process and Dockerfile CMDs itself.
        - <<: *proxy_conf
      # add driver-opt http config to tell buildkit + buildx to resolve external checksums through a proxy.
      buildkit_driveropt:
        - "env.http_proxy=http://X.Y.Z.Z:3128"
        - "env.https_proxy=http://X.Y.Z.Z:3128"
        - "env.no_proxy=.my-subdomain.com"
```

## Using cache images

You can provide a list of images to use for cache.
These cache images are built with mode=max, image-manifest=true, and oci-mediatypes=true.
This is to provide better usage of cache and better compatibility with image stores like Harbor.

```yaml
steps:
  build:
    image: woodpeckerci/plugin-docker-buildx
    settings:
      repo: hari/radiant
      cache_images:
        - hari/radiant:cache
        - harbor.example.com/hari/radiant:cache
      logins:
        - registry: https://index.docker.io/v1/
          username: hari
          password:
            from_secret: docker_password
        - registry: https://harbor.example.com
          username: hari
          password:
            from_secret: harbor_password
```

## Using other cache types

You can specify cache_to and cache_from to use specific settings.
For example you can configure an s3 object as cache.

More details can be found [in the docker docs](https://docs.docker.com/build/cache/backends/).

```yaml
steps:
  build:
    image: woodpeckerci/plugin-docker-buildx
    settings:
      repo: hari/radiant
      cache_to: type=s3,region=east,bucket=mystuff,name=radiant-cache
      cache_from: type=s3,region=east,bucket=mystuff,name=radiant-cache
```
