#!/bin/sh

# ADD: to support gomod
# export GOPATH=$GOPATH:/Users/bogao/develop/gitme/gotapestry
export GO111MODULE=on

echo "Building Project..."
# echo "-[BUILD:1/6]- remove some folder first:"
# rm -rf ~/develop/go/pkg/darwin_amd64/github.com/elivoa/got

# echo "-[BUILD:2/6]- go install github.com/elivoa/gxl..."
# go install github.com/elivoa/gxl

# echo "-[BUILD:3/6]- go install github.com/elivoa/got..."
# go install github.com/elivoa/got

# echo "-250- go install got..."
# go install got

# echo "-[BUILD:4/6]- go install syd... ?? need to remove and build again?"
# go install syd
# go install syd/service
# go install syd/model
# go install syd/dal/userdao

# How to generate all components and pages package to build.
# echo "-[BUILD:5/6]- go install generated; to build all pages."
# go install syd/generated

echo "-[BUILD:6/6]- go run main.go..."

echo "---------------------------------------------------------------------------------------------------"
go run src/syd/generated/main.go
