📌 API Доски объявлений класса
Простое REST API для управления школьными объявлениями.

🚀 Запуск
bash
go run main.go
Сервер: http://localhost:8080

📡 Методы API
GET /api/announcements — все объявления

GET /api/announcements/{id} — одно объявление

POST /api/announcements — создать

PUT /api/announcements/{id} — обновить

DELETE /api/announcements/{id} — удалить

GET /health — проверка сервера

💡 Примеры
bash
# Получить все
curl http://localhost:8080/api/announcements

# Создать
curl -X POST http://localhost:8080/api/announcements \
  -H "Content-Type: application/json" \
  -d '{"title":"Новое","content":"Текст","author":"Имя"}'
📁 Файлы
main.go — код сервера

bulletin.json — хранилище данных

go.mod — зависимости

