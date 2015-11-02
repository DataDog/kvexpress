#!/bin/bash

set -e

export PREDICTED_CHECKSUM="0ab71a1c8fef24ade8d650e2cc248aac1e499a45a0e9456ba0b47901f99176d8"

make deps
make build

echo "Launching Consul."
consul agent -data-dir `mktemp -d` -bootstrap -server -bind=127.0.0.1 1>/dev/null &
sleep 3
curl -s https://gist.githubusercontent.com/darron/94447bfab90617f16962/raw/d4cb39471724800ba9e731f99e5844167e93c5df/sorting.txt > sorting
echo "Putting 'sorting' into 'testing' key."
bin/kvexpress in -k testing -f sorting --sorted true
echo "Pulling 'testing' key out and saving it to 'output'."
bin/kvexpress out -k testing -f output

export CHECKSUM=$(shasum -a 256 output | cut -d ' ' -f 1)

echo "Testing clean command."
bin/kvexpress clean -f sorting
bin/kvexpress clean -f output

echo "Checksum : $CHECKSUM"
echo "Predicted: $PREDICTED_CHECKSUM"

if [[ "$CHECKSUM" == "$PREDICTED_CHECKSUM" ]]; then
  echo "Looks good."
  exit 0
else
  echo "Looks bad - checksums don't match."
  exit 1
fi
