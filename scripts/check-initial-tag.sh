#!/bin/bash

set -e

test "$(git tag -l | grep -c 'v0.0.0')" -eq 1

exit 0