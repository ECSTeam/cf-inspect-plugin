#!/bin/bash

set -e

(cf uninstall-plugin "inspect" || true) && go build -o inspect-plugin cf-inspect.go && cf install-plugin inspect-plugin
