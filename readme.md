# Install [Go](https://github.com/golang/go) first
## You need to create the log files first
1. Log files configuration refer to `main.go`:
    >```logs.SetLogger(logs.AdapterFile, `{"filename":"logs/dev.log","level":7,"daily":true,"maxdays":10}`)```
1. Create folder
    >`mkdir logs`
1. Touch files
    >`touch logs/dev.log`
    >`touch logs/project.log`
## How to start local dev server:
1. Database configuration refer to `main.go`:
    >`orm.RegisterDataBase("default", "mysql", "root:Start123@/group_buying?charset=utf8")`
1. Start a terminal in the root direcrtory
1. Run 
    >`bee run` 
1. Visit pr
    >`localhost:6666/dev`

## How to pack the project
1. Start a terminal in the root direcrtory
1. Clean all exist file use `go clean`
1. Install all the related packages by: `go get`
1. Run
    >`bee pack -be GOOS=$targetenv`

`$targetenv` is the server system environments.

For example: `GOOS=linux` 
