@echo off
go install ffServer\ffWebServer

move /y ..\bin\ffWebServer.exe ..\..\ffbin\ffServerBin\ffWebServer\ffWebServer.exe

cd ..\..\ffbin\ffServerBin\ffWebServer

ffWebServer.exe
