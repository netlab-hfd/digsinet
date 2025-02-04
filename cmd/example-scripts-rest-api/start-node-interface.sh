#!/usr/bin/env bash
curl -X POST -H "Content-Type: application/json" \
  -H "X-Auth-Key: key" -H "X-Auth-Secret: secret" \
  -d '{ "path": "system/state/hostname" }' \
  http://localhost:8088/v1/sibling/sec/arista_ceos/start-node-iface