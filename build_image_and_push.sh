#!/bin/bash

PLATFORMS=linux/amd64,linux/arm64

docker buildx build \
	--platform=$PLATFORMS \
	-t luksamuk/minerva_tui:latest \
	--push \
	.

