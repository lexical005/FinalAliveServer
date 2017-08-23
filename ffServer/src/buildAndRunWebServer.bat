@echo off
go install ffServer\ffWebServer || goto error1

move /y ..\bin\ffWebServer.exe ..\..\ffbin\ffServerBin\ffWebServer\ffWebServer.exe

cd ..\..\ffbin\ffServerBin\ffWebServer

ffWebServer.exe
goto:eof

:error1
echo go install ffServer\ffWebServer error
pause
goto:eof
