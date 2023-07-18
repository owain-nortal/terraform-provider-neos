#!/bin/bash
go generate ./... 
gofmt -s -w .
golangci-lint run