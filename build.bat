@echo off&setlocal enabledelayedexpansion
title golang build
chcp 65001 >nul
:this is main name
set output=checkport
:this is version
set version=v1.0
@echo output:%output%
@echo build version:%version%

set outfilename=%output%_windows_386_%version%.exe
@echo build windows/386 …… %outfilename%
set GOOS=windows&&set GOARCH=386&& go build   -ldflags "-w -s" -buildmode=pie -o %outfilename%  main.go

set outfilename=%output%_windows_amd64_%version%.exe
@echo build windows/amd64 …… %outfilename%
set GOOS=windows&&set GOARCH=amd64&& go build   -ldflags "-w -s" -buildmode=pie -o %outfilename%  main.go

set outfilename=%output%_linux_386_%version%
@echo build linux/386 …… %outfilename%
set GOOS=linux&&set GOARCH=386&& go build   -ldflags "-w -s"  -o %outfilename%  main.go

set outfilename=%output%_linux_amd64_%version%
@echo build linux/amd64 …… %outfilename%
set GOOS=linux&&set GOARCH=amd64&& go build   -ldflags "-w -s" -buildmode=pie -o %outfilename%  main.go

set outfilename=%output%_darwin_amd64_%version%
@echo build darwin/amd64 …… %outfilename%
set GOOS=darwin&&set GOARCH=amd64&& go build   -ldflags "-w -s" -o %outfilename%   main.go
@echo build finished!
@pause