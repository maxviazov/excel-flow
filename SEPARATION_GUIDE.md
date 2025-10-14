# 🎯 Руководство по разделению Backend и Frontend

## ✅ Что сделано

### 1. Frontend (готов к деплою на S3)
- ✅ Создана структура `frontend/`
- ✅ Разделены public и admin приложения
- ✅ Добавлена конфигурация API (`config.js`)
- ✅ Создан скрипт деплоя на S3

### 2. Backend (API-only)
- ✅ Убрана раздача статики
- ✅ Оставлены только API endpoints
- ✅ Добавлен `/health` endpoint
- ✅ CORS настроен для S3

## 🚀 Шаги деплоя

### Шаг 1: Деплой Backend на ECS

```bash
# В корне проекта
./deploy.sh
```

После деплоя получите URL API (например: `https://api.excel-flow.com`)

### Шаг 2: Обновите конфигурацию Frontend

Откройте `frontend/public/config.js` и замените:

```javascript
const API_BASE_URL = 'https://YOUR_API_DOMAIN_HERE';
```

на ваш реальный URL:

```javascript
const API_BASE_URL = 'https://api.excel-flow.com';
```

### Шаг 3: Деплой Frontend на S3

```bash
cd frontend
./deploy-s3.sh
```

Скрипт:
1. Создаст S3 bucket
2. Настроит статический хостинг
3. Загрузит файлы
4. Настроит публичный доступ
5. Выведет URL сайта

### Шаг 4: (Опционально) Настройте CloudFront

Для HTTPS и CDN:

1. Создайте CloudFront distribution
2. Origin: ваш S3 bucket
3. SSL Certificate: создайте в ACM
4. Custom domain: `excel-flow.yourdomain.com`

## 🧪 Локальное тестирование

### Backend
```bash
go run cmd/server/main.go cmd/server/admin_handlers.go
```
API доступен на: http://localhost:8080

### Frontend
```bash
cd frontend
python3 -m http.server 3000
```
Откройте:
- http://localhost:3000/public - главное приложение
- http://localhost:3000/admin - админ-панель

## 📝 Важные файлы

### Backend
- `cmd/server/main.go` - API сервер
- `cmd/server/admin_handlers.go` - админ endpoints
- `deploy.sh` - деплой на ECS

### Frontend
- `frontend/public/config.js` - **ОБЯЗАТЕЛЬНО обновите API URL**
- `frontend/deploy-s3.sh` - деплой на S3
- `frontend/README.md` - подробная документация

## 💰 Стоимость

### До разделения
- ECS 24/7: ~$50-70/месяц

### После разделения
- ECS (API only): ~$15-30/месяц
- S3 + CloudFront: ~$1-3/месяц
- **Экономия**: 50-60%

## 🔧 Troubleshooting

### CORS ошибки
Убедитесь, что в `cmd/server/main.go` настроен CORS:
```go
w.Header().Set("Access-Control-Allow-Origin", "*")
```

### 404 на S3
Проверьте:
1. Bucket policy настроен
2. Static website hosting включен
3. Файлы загружены правильно

### API не отвечает
Проверьте:
1. ECS task запущен
2. Security group разрешает порт 8080
3. Health check `/health` работает

## 📚 Дополнительная информация

- `ARCHITECTURE.md` - подробная архитектура
- `frontend/README.md` - документация фронтенда
- `DEPLOY.md` - деплой бэкенда

## ✨ Следующие шаги

1. ✅ Задеплойте backend на ECS
2. ✅ Обновите `config.js` с URL API
3. ✅ Задеплойте frontend на S3
4. ⭐ (Опционально) Настройте CloudFront для HTTPS
5. ⭐ (Опционально) Настройте custom domain
6. ⭐ (Опционально) Добавьте мониторинг CloudWatch

Готово! 🎉
