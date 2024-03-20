# plugin-docker-buildx

<br/>
<p align="center">
<a href="https://ci.codeberg.org/repos/3265" target="_blank">
  <img src="https://ci.codeberg.org/api/badges/3265/status.svg" alt="status-badge" />
</a>
<a href="https://codeberg.org/woodpecker-plugins/docker-buildx/releases" title="Latest release">
  <img src="https://img.shields.io/gitea/v/release/woodpecker-plugins/docker-buildx?gitea_url=https%3A%2F%2Fcodeberg.org
" alt="Latest release">
</a>
  <a href="https://matrix.to/#/#woodpecker:matrix.org" title="Join the Matrix space at https://matrix.to/#/#woodpecker:matrix.org">
    <img src="https://img.shields.io/matrix/woodpecker:matrix.org?label=matrix" alt="Matrix space">
  </a>
  <a href="https://hub.docker.com/r/woodpeckerci/plugin-docker-buildx" title="Docker pulls">
    <img src="https://img.shields.io/docker/pulls/woodpeckerci/plugin-docker-buildx" alt="Docker pulls">
  </a>
  <a href="https://opensource.org/licenses/Apache-2.0" title="License: Apache-2.0">
    <img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License: Apache-2.0">
  </a>
</p>
<br/>

Woodpecker CI plugin to build multiarch Docker images with [buildx](https://duckduckgo.com/?q=docker+buildx&ia=web).
This plugin was initially a fork of [thegeeklab/drone-docker-buildx](https://github.com/thegeeklab/drone-docker-buildx/) (now archived in favor of this plugin) which itself was a fork of [drone-plugins/drone-docker](https://github.com/drone-plugins/drone-docker).
I also contains the ability to publish to AWS ECR which was previously provided by [drone-plugins/drone-ecr](https://github.com/drone-plugins/drone-docker/tree/master/cmd/drone-ecr).
You can find the full documentation at [woodpecker-ci.org](https://woodpecker-ci.org/plugins/Docker%20Buildx) ([docs.md](./docs.md)).

## Images

Images are available on [Dockerhub](https://hub.docker.com/r/woodpeckerci/plugin-docker-buildx) and in the [Codeberg registry](https://codeberg.org/woodpecker-plugins/-/packages/container/docker-buildx/latest).

## License

This project is licensed under the Apache-2.0 License - see the [LICENSE](https://codeberg.org/woodpecker-plugins/plugin-docker-buildx/src/branch/main/LICENSE) file for details.

## Maintainers

This plugin is maintained by @6543 and @pat-s.
