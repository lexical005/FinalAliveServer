@echo off
go install ffTool\ffExportExcel

move /y ..\bin\ffExportExcel.exe ..\..\..\FinalAlive\Config\Game\ffExportExcel.exe

cd ..\..\..\FinalAlive\Config\Game

runByPro.bat
