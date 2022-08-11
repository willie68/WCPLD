@echo off
go build -ldflags="-s -w" -o wcupl.exe cmd/main.go
.\3rd_party\7z.exe a wcupl.zip wcupl.exe
copy wcupl.exe e:\WinCupl\shared