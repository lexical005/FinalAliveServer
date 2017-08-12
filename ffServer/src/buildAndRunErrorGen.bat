@echo off
go install ffTool\ffError

move /y ..\bin\ffError.exe ..\..\..\FinalAlive\Config\Error\ffError.exe

cd ..\..\..\FinalAlive\Config\Error

runByPro.bat
