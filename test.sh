#!/bin/bash
curl -X GET http://localhost:8080/v1/hello/say
echo ''
curl -X POST -d '{"name": "Alvin"}' http://localhost:8080/v1/hello/greeting
echo ''
