# Excel Flow - Статус проекта

## ✅ Текущее состояние

Проект разделен на Backend (API) и Frontend (Static). Готов к деплою.

## 📁 Структура проекта

### Backend (Go API)
```
cmd/
  ├── server/    # API сервер (порт 8080)
  └── admin/     # Админ-панель API (порт 8081)

internal/
  ├── app/       # Бизнес-логика
  ├── drivers/   # Справочник водителей
  ├── textutil/  # Очистка текста
  ├── writer/    # Генерация Excel
  └── ...

configs/
  ├── dictionaries/  # Справочники (города, водители)
  └── pipeline.yaml  # Конфигурация
```

### Frontend (Static)
```
frontend/
  ├── public/          # Главное приложение
  │   ├── index.html
  │   ├── app.js
  │   └── config.js    # API URL
  ├── admin/           # Админ-панель
  │   ├── index.html
  │   ├── app.js
  │   └── style.css
  └── deploy-s3.sh     # Деплой на S3
```

## 🚀 Запуск

### Backend (локально)
```bash
go run cmd/server/main.go cmd/server/admin_handlers.go
```
API: http://localhost:8080

### Frontend (локально)
```bash
cd frontend
python3 -m http.server 3000
```
- http://localhost:3000/public
- http://localhost:3000/admin

## 📦 Деплой

### Backend → ECS Fargate
```bash
./deploy.sh
```

### Frontend → S3
```bash
cd frontend
# Обновите config.js с URL API
./deploy-s3.sh
```

## 🎯 Основные функции

1. ✅ Загрузка SAP Excel файлов
2. ✅ Обработка и группировка данных
3. ✅ Очистка HTML-entities
4. ✅ Транслитерация английских названий в иврит
5. ✅ Автоназначение водителей по коду города
6. ✅ Генерация отчетов для МОЗ
7. ✅ API для фронтенда
8. ✅ Управление справочниками

## 📊 Справочники

### Водители (8 записей)
- סיומה אורמן - 659-86-601
- ארתור מסרסקי - 695-70-102
- פבל ולר - 69-570-202
- מיכאל טולשין - 790-93-702
- רומן רצ'ימוב - 402-20-603
- אלכס רוגוזין - 402-39-703
- סרגיי נקיטין - 292-58-003
- אלכסיי יפרמוב - 368-91-801

### Города
85 уникальных кодов городов

## 💰 Стоимость

- **Backend (ECS)**: ~$15-30/месяц
- **Frontend (S3)**: ~$1-3/месяц
- **Итого**: ~$16-33/месяц

## 🧪 Тестирование

```bash
go test ./...
```

Все тесты проходят ✅

## 📝 Документация

- **SEPARATION_GUIDE.md** - руководство по разделению
- **ARCHITECTURE.md** - архитектура системы
- **frontend/README.md** - документация фронтенда
- **DEPLOY.md** - деплой бэкенда
- **docs/DRIVERS.md** - модуль водителей
- **docs/TEXT_SANITIZATION.md** - очистка текста

## 🔄 API Endpoints

### Main API
- `POST /api/upload` - загрузка файла
- `POST /api/process` - обработка файла
- `GET /api/download/:filename` - скачивание
- `GET /health` - health check

### Admin API
- `GET /api/admin/cities` - список городов
- `POST /api/admin/cities` - добавить город
- `DELETE /api/admin/cities` - удалить город
- `GET /api/admin/drivers` - список водителей
- `POST /api/admin/drivers` - добавить водителя
- `DELETE /api/admin/drivers` - удалить водителя

## ✨ Готово к продакшену!

Проект полностью разделен и готов к деплою на AWS 🚀
