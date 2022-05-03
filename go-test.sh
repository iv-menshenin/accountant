#!/bin/bash

export TST_MONGO_HOST=mongo
export TST_MONGO_USERNAME=simachat
export TST_MONGO_PASSWORD=mongopassword1
export TST_MONGO_DBNAME=autotests
go test ./...
