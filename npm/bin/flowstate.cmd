@ECHO off
SETLOCAL
CALL :find_dp0

IF EXIST "%dp0%\flowstate.exe" (
  SET "_prog=%dp0%\flowstate.exe"
) ELSE (
  SET "_prog=%dp0%\..\bin\flowstate.exe"
)

"%_prog%" %*
ENDLOCAL
EXIT /b %errorlevel%
:find_dp0
SET dp0=%~dp0
EXIT /b

