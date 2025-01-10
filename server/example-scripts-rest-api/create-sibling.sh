#!/usr/bin/env bash
curl -X POST -H "Content-Type: application/json" -H "X-Auth-Key: key" -H "X-Auth-Secret: secret" -d '@./post-sibling.json' http://localhost:8088/v1/sibling
curl -X GET http://localhost:8088/siblings
