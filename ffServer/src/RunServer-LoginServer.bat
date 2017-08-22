@echo off
go install ffServer\ffLoginServer

move /y ..\bin\ffLoginServer.exe ..\..\..\FinalAlive\Server\ffLoginServer\ffLoginServer.exe

cd ..\..\..\FinalAlive\Server\ffLoginServer

start.bat
