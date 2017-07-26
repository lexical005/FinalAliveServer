@echo off
go install ffTool\ffExportExcel

move /y ..\bin\ffExportExcel.exe ..\..\ffbin\ffGameConfig\ffExportExcel.exe

cd ..\..\ffbin\ffGameConfig

runByPro.bat
