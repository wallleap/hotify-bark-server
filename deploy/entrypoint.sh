#!/usr/bin/env bash

set -e

# The image runs as a non-root user, so /etc is not writable at runtime.
# A default timezone symlink is baked into the image at build time (see
# Dockerfile); drop runtime timezone mutation entirely.

exec "$@"
