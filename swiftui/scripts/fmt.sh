#!/bin/bash

cmd='swift-format --configuration ./swift-format.json'

if [ $# -ne 1 ]; then
  # format all
  $cmd --in-place "$@"
  exit $?
fi

# first pass to check if file need change
formatted=$($cmd "$1")
needInPlace=$($cmd "$1" | diff "$1" -)
if [ ! -z "$needInPlace" ]; then
  # fuck Xcode, we have to give it some time to complete whatever it does after
  # saving file before execute the formating so that it won't complain about
  # external change.
  sleep 0.3
  $cmd --in-place "$1"
fi
