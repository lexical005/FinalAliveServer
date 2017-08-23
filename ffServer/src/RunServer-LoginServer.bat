@echo off
go install ffServer\ffLoginServer || goto error1

move /y ..\bin\ffLoginServer.exe ..\..\..\FinalAlive\Server\ffLoginServer\ffLoginServer.exe

cd ..\..\..\FinalAlive\Server\ffLoginServer

start.bat
goto:eof

:error1
echo go install ffServer\ffLoginServer error
pause
goto:eof
