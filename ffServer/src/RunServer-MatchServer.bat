@echo off
go install ffServer\ffMatchServer || goto error1

move /y ..\bin\ffMatchServer.exe ..\..\..\FinalAlive\Server\ffMatchServer\ffMatchServer.exe

cd ..\..\..\FinalAlive\Server\ffMatchServer

start.bat
goto:eof

:error1
echo go install ffServer\ffMatchServer error
pause
goto:eof
