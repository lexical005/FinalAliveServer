@echo off
go install ffTool\ffExportExcel || goto error1

move /y ..\bin\ffExportExcel.exe ..\..\..\FinalAlive\Config\Game\ffExportExcel.exe

cd ..\..\..\FinalAlive\Config\Game

make_all.bat
goto:eof

:error1
echo go install ffServer\ffExportExcel error
pause
goto:eof
