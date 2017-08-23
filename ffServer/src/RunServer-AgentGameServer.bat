@echo off
go install ffServer\ffAgentGameServer || goto error1

move /y ..\bin\ffAgentGameServer.exe ..\..\..\FinalAlive\Server\ffAgentGameServer\ffAgentGameServer.exe

cd ..\..\..\FinalAlive\Server\ffAgentGameServer

start.bat
goto:eof

:error1
echo go install ffServer\ffAgentGameServer error
pause
goto:eof
