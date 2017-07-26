@echo off
go install ffServer\ffAgentServer

move /y ..\bin\ffAgentServer.exe ..\..\ffbin\ffServerBin\ffAgentGameServer\ffAgentGameServer.exe

cd ..\..\ffbin\ffServerBin\ffAgentGameServer

ffAgentGameServer.exe
