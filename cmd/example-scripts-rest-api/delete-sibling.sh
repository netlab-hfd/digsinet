#!/usr/bin/env bash
curl -X DELETE -H "Content-Type: application/json" \
  -H "X-Auth-Key: key" -H "X-Auth-Secret: secret" \
  -d '' \
  http://localhost:8088/v1/sibling/sec