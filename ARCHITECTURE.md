# Excel Flow - Архитектура

## 🏗️ Разделение на Backend и Frontend

### Backend (Go API) - ECS Fargate
**Расположение**: `cmd/server/`, `internal/`  
**Деплой**: AWS ECS Fargate  
**Стоимость**: ~$15-30/месяц (работает только при обработке)

**Endpoints**:
- `POST /api/upload` - загрузка файла
- `POST /api/process` - обработка файла
- `GET /api/download/:filename` - скачивание результата
- `GET /api/admin/cities` - управление городами
- `GET /api/admin/drivers` - управление водителями
- `GET /health` - health check

**Особенности**:
- Только API, без статики
- CORS для доступа с S3
- Обработка Excel файлов
- Справочники в SQLite

### Frontend (Static) - S3 + CloudFront
**Расположение**: `frontend/`  
**Деплой**: AWS S3 + CloudFront  
**Стоимость**: ~$1-3/месяц

**Структура**:
```
frontend/
├── public/          # Главное приложение
│   ├── index.html
│   ├── app.js
│   └── config.js    # API_BASE_URL
└── admin/           # Админ-панель
    ├── index.html
    ├── app.js
    └── style.css
```

**Особенности**:
- Чистый HTML/CSS/JS
- Конфигурация API через config.js
- CDN через CloudFront
- HTTPS из коробки

## 🔄 Взаимодействие

```
┌─────────────┐         HTTPS          ┌──────────────┐
│   Browser   │ ──────────────────────> │  S3/CloudFront│
│             │                         │   (Frontend)  │
└─────────────┘                         └──────────────┘
       │                                        
       │ API Calls (CORS)                      
       │                                        
       v                                        
┌─────────────┐         HTTPS          ┌──────────────┐
│   Browser   │ ──────────────────────> │  ECS Fargate │
│             │                         │   (Backend)  │
└─────────────┘                         └──────────────┘
```

## 📦 Деплой

### Backend
```bash
cd /path/to/excel-flow
./deploy.sh
```

### Frontend
```bash
cd frontend
# Обновите config.js с URL бэкенда
./deploy-s3.sh
```

## 💰 Экономия

**До разделения**:
- ECS 24/7 с статикой: ~$50-70/месяц

**После разделения**:
- ECS (только API): ~$15-30/месяц
- S3 + CloudFront: ~$1-3/месяц
- **Итого**: ~$16-33/месяц (экономия 50-60%)

## 🔒 Безопасность

- Backend: CORS настроен только для фронтенд домена
- Frontend: HTTPS через CloudFront
- API: Rate limiting (опционально)
- S3: Публичный доступ только на чтение

## 🚀 Масштабирование

- **Backend**: ECS автоматически масштабируется при нагрузке
- **Frontend**: CloudFront CDN обслуживает из ближайшей точки
- **База данных**: SQLite в контейнере (для production рекомендуется RDS)

## 📊 Мониторинг

- Backend: CloudWatch Logs + Metrics
- Frontend: CloudFront Access Logs
- Health Check: `/health` endpoint
