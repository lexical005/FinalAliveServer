@echo off
go install ffServer\ffMatchServer

move /y ..\bin\ffMatchServer.exe ..\..\..\FinalAlive\Server\ffMatchServer\ffMatchServer.exe

cd ..\..\..\FinalAlive\Server\ffMatchServer

start.bat
