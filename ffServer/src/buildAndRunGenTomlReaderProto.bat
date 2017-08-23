@echo off
go install ffTool\ffGenTomlReader || goto error1

move /y ..\bin\ffGenTomlReader.exe ..\..\..\FinalAlive\Config\Game\ffGenTomlReader.exe

cd ..\..\..\FinalAlive\Config\Game

ffGenTomlReader.exe -gocodedir "../../../FinalAliveServer/ffServer/src/ffAutoGen/ffClientToml" -readername ffClientTomlTranslator -proto proto -csharp csharp
pause
goto:eof

:error1
echo go install ffServer\ffGenTomlReader error
pause
goto:eof
