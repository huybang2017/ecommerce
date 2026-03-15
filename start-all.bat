@echo off
chcp 65001 >nul
echo Starting all services and frontends...

REM Create log directory if not exists
if not exist "logs" mkdir logs

REM Get current date and time for log files
set timestamp=%date:~-4%%date:~3,2%%date:~0,2%_%time:~0,2%%time:~3,2%%time:~6,2%
set timestamp=%timestamp: =0%

echo.
echo ======================================
echo Starting Go Services
echo ======================================
echo.

REM Start API Gateway
echo Starting api-gateway...
start /B "" cmd /c "cd api-gateway && go run cmd/main.go > ../logs/api-gateway_%timestamp%.log 2>&1"
timeout /t 2 /nobreak >nul

REM Start Identity Service
echo Starting identity-service...
start /B "" cmd /c "cd identity-service && go run cmd/main.go > ../logs/identity-service_%timestamp%.log 2>&1"
timeout /t 2 /nobreak >nul

REM Start Order Service
echo Starting order-service...
start /B "" cmd /c "cd order-service && go run cmd/main.go > ../logs/order-service_%timestamp%.log 2>&1"
timeout /t 2 /nobreak >nul

REM Start Product Service
echo Starting product-service...
start /B "" cmd /c "cd product-service && go run cmd/main.go > ../logs/product-service_%timestamp%.log 2>&1"
timeout /t 2 /nobreak >nul

REM Start Search Service
echo Starting search-service...
start /B "" cmd /c "cd search-service && go run cmd/main.go > ../logs/search-service_%timestamp%.log 2>&1"
timeout /t 2 /nobreak >nul

@REM echo.
@REM echo ======================================
@REM echo Starting Frontend Applications
@REM echo ======================================
@REM echo.

@REM REM Start Admin Frontend
@REM echo Starting admin frontend...
@REM start /B "" cmd /c "cd admin && npm start > ../logs/admin_%timestamp%.log 2>&1"
@REM timeout /t 2 /nobreak >nul

@REM REM Start Seller Frontend
@REM echo Starting seller frontend...
@REM start /B "" cmd /c "cd seller && npm start > ../logs/seller_%timestamp%.log 2>&1"
@REM timeout /t 2 /nobreak >nul

@REM REM Start Client Frontend
@REM echo Starting client frontend...
@REM start /B "" cmd /c "cd client && npm run dev > ../logs/client_%timestamp%.log 2>&1"
@REM timeout /t 2 /nobreak >nul

echo.
echo ======================================
echo All services started!
echo ======================================
echo.
echo Log files are being written to the 'logs' directory with timestamp: %timestamp%
echo.
echo Press any key to exit this window (services will continue running)...
pause >nul
