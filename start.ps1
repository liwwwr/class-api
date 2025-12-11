Write-Host "🚀 ЗАПУСК СЕРВЕРА API" -ForegroundColor Green
Write-Host "====================" -ForegroundColor Green
Write-Host ""

# Проверяем Go
Write-Host "Проверяем Go..." -ForegroundColor Yellow
go version

Write-Host ""
Write-Host "Запускаем сервер..." -ForegroundColor Yellow
Write-Host "Сервер будет доступен по адресу: http://localhost:8080" -ForegroundColor Cyan
Write-Host "Чтобы остановить сервер - нажмите Ctrl+C" -ForegroundColor Red
Write-Host ""

go run main.go