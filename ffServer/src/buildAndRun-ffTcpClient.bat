@echo off
go install ffTest\ffTcpClient

cd ..\bin\

start ffTcpClient.exe
