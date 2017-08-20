@echo off
go install ffServer\ffAgentGameServer

move /y ..\bin\ffAgentGameServer.exe ..\..\..\FinalAlive\Server\ffAgentGameServer\ffAgentGameServer.exe

cd ..\..\..\FinalAlive\Server\ffAgentGameServer

start.bat
