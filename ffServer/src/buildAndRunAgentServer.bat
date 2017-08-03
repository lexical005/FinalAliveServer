@echo off
go install ffServer\ffAgentServer

move /y ..\bin\ffAgentGameServer.exe ..\..\..\FinalAlive\Server\ffAgentGameServer\ffAgentGameServer.exe

cd ..\..\..\FinalAlive\Server\ffAgentGameServer

ffAgentGameServer.exe
