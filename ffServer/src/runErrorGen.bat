@echo off
go install ffTool\ffError

move /y ..\bin\ffError.exe ..\..\ffbin\ffError\ffError.exe

cd ..\..\ffbin\ffError

runByPro.bat
