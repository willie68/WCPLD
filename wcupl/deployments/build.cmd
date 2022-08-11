@echo off
go build -ldflags="-s -w" -o wcupl.exe cmd/main.go
copy wcupl.exe e:\WinCupl\shared