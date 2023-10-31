#!/bin/bash
go install github.com/swaggo/swag/cmd/swag@v1.7.8
#  create docs.go at  docs/docs.go
#  create swagger.json at  docs/swagger.json
#  create swagger.yaml at  docs/swagger.yaml
swag init --parseDependency --parseInternal -g cmd/app/main.go
#swag init -g internal/controller/http/v1/router.go

# make migrations package
go install github.com/jteeuwen/go-bindata/...@latest
cd migrations/
go-bindata -pkg migrations .
cd -

go build -o go-clean-template cmd/app/main.go
 ./go-clean-template -folder dev -name hello -listen ":8080"
