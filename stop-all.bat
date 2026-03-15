@echo off
echo Stopping services...

set ports=8000 8081 8080 8083 8082

for %%p in (%ports%) do (
    echo Killing port %%p ...
    for /f "tokens=5" %%a in ('netstat -ano ^| find ":%%p"') do (
        taskkill /F /PID %%a >nul 2>&1
    )
)

timeout /t 2 >nul

echo Cleaning logs...

if exist logs (
    rmdir /S /Q logs
)

echo Done.
pause
