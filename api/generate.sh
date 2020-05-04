#!/usr/bin/env bash
set -e

rootdir=`dirname $0`
pushd $rootdir

pushd ${rootdir}/../
  $(GO111MODULE=on go get -u "github.com/algorand/oapi-codegen/...@v1.3.5-algorand4")
popd


# Convert v2 to v3
curl -s -X POST "https://converter.swagger.io/api/convert" -H "accept: application/json" -H "Content-Type: application/json" -d @./indexer.oas2.json  -o 3.json

# Sort keys, format json and rename 3.json -> algod.oas3.yml
python3 -c "import json; import sys; json.dump(json.load(sys.stdin), sys.stdout, indent=2, sort_keys=True)" < 3.json > indexer.oas3.yml
rm 3.json

echo "generating code."
oapi-codegen -package generated -type-mappings integer=uint64 -generate types -o generated/types.go indexer.oas3.yml
oapi-codegen -package generated -type-mappings integer=uint64 -generate server,spec -o generated/routes.go indexer.oas3.yml

