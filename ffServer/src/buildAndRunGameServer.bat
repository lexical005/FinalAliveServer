@echo off
go install ffServer\ffGameServer

move /y ..\bin\ffGameServer.exe ..\..\ffbin\ffServerBin\ffGameServer\ffGameServer.exe

cd ..\..\ffbin\ffServerBin\ffGameServer

ffGameServer.exe
