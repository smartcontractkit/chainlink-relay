#!/bin/bash

set -e

cd "$(dirname "$0")"/..

go generate -v ./core/services/monitoring/generate.go
