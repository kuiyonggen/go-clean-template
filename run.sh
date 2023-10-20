#!/bin/bash
go build -o go-clean-template cmd/app/main.go
 ./go-clean-template -folder dev -name hello -listen ":8080"
