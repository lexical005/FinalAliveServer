@echo off
go install ffTool\ffGenTomlReader

move /y ..\bin\ffGenTomlReader.exe ..\..\..\FinalAlive\Config\Game\ffGenTomlReader.exe

cd ..\..\..\FinalAlive\Config\Game

ffGenTomlReader.exe -gocodedir "../../../FinalAliveServer/ffServer/src/ffAutoGen/ffClientToml" -readername ffClientTomlTranslator -proto proto -csharp csharp

pause