[![Coverage Status](https://coveralls.io/repos/github/xplorfin/lndmock/badge.svg?branch=master)](https://coveralls.io/github/xplorfin/lndmock?branch=master)
[![Renovate enabled](https://img.shields.io/badge/renovate-enabled-brightgreen.svg)](https://app.renovatebot.com/dashboard#github/xplorfin/lndmock)
[![Build status](https://github.com/xplorfin/lndmock/workflows/test/badge.svg)](https://github.com/xplorfin/lndmock/actions?query=workflow%3Atest)
[![Build status](https://github.com/xplorfin/lndmock/workflows/goreleaser/badge.svg)](https://github.com/xplorfin/lndmock/actions?query=workflow%3Agoreleaser)
[![](https://godoc.org/github.com/xplorfin/lndmock?status.svg)](https://godoc.org/github.com/xplorfin/lndmock)
[![Go Report Card](https://goreportcard.com/badge/github.com/xplorfin/lndmock)](https://goreportcard.com/report/github.com/xplorfin/lndmock)

# What is this?

This is a helper library by [entropy](http://entropy.rocks/) that hopes to imitate some of the functionality of [polar](https://github.com/jamaljsr/polar) for continuous integration. Right now this supports creating btcd and lnd nodes and funding them. You can check out an example in `lnd_test`. This also supports mocking bolt11 invoices

## Note:

This library was open sourced as a dependency for another project. While this is functional, documentation may be lacking for a bit.

## Helpful tips:

This command will remove lndmock containers:

<!-- TODO: move this into Makefile -->
`docker stop $(docker ps -a -q --filter label="created-by=docker-utils") || true && docker rm -v $(docker ps -a -q --filter label="created-by=docker-utils") || true && docker network rm lightning-network || true && docker volume rm $( docker volume ls -q --filter label="created-by=docker-utils") || true`
<!-- see:  https://docs.docker.com/engine/reference/commandline/ps/ -->

# Docker Images

Docker images are currently built in [this](https://github.com/xplorfin/lnd-docker-images) repository. Eventually, we'd like to support bitcoind, eclair, c-lightning, etc and automate retroactive/proactive builds (e.g. a consistent lnd master image). Polar has some docs [here](https://github.com/jamaljsr/polar/blob/master/docs/custom-nodes.md) and we're probably best off following their lead