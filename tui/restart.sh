#!/bin/bash
# This is useful to restart the applications with some file system watcher tool
# such as nodmeon.

pkill tui
go build -race -o tui . && ./tui
