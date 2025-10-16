# 🚀 Deploy Phase 2 to AWS

## ✅ Готово к деплою

**Ветка:** `main`  
**Коммит:** `e89ec54`  
**Изменения:** Phase 2 UI/UX improvements

---

## 📋 Что нужно задеплоить

### 1. Frontend (S3 + CloudFront)
- Новый UI с прогресс-баром
- Toast уведомления
- История обработок
- Улучшенный drag & drop

### 2. Backend (ECS)
- Обновленный CORS (добавлен localhost:3000)
- Очищен от заглушек admin handlers

---

## 🔧 Деплой

### Вариант 1: Полный деплой (Backend + Frontend)

```bash
# Настроить AWS credentials
export AWS_ACCESS_KEY_ID="your-key"
export AWS_SECRET_ACCESS_KEY="your-secret"
export AWS_REGION="us-east-1"

# Или через aws configure
aws configure

# Запустить деплой
./deploy-simple.sh
```

### Вариант 2: Только Frontend (быстрее)

```bash
# Деплой на S3
cd frontend
aws s3 sync public/ s3://excel-viazov-dev/ --delete

# Инвалидация CloudFront кеша
aws cloudfront create-invalidation \
  --distribution-id E1GG21S86PMA83 \
  --paths "/*"
```

### Вариант 3: Только Backend

```bash
# Build и push Docker image
docker build -t excel-flow:latest .

# Tag для ECR
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
docker tag excel-flow:latest $AWS_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/excel-flow:latest

# Login to ECR
aws ecr get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin \
  $AWS_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com

# Push
docker push $AWS_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/excel-flow:latest

# Update ECS service
aws ecs update-service \
  --cluster excel-flow-cluster \
  --service excel-flow-service \
  --force-new-deployment \
  --region us-east-1
```

---

## ✅ После деплоя

### Проверить:

1. **Frontend:** https://excel.viazov.dev
   - ✅ Drag & drop работает
   - ✅ Прогресс-бар показывается
   - ✅ Toast уведомления появляются
   - ✅ История сохраняется
   - ✅ Предпросмотр файла

2. **Backend:** https://api.viazov.dev
   - ✅ `/health` возвращает `{"status":"ok"}`
   - ✅ CORS работает
   - ✅ Загрузка файлов
   - ✅ Обработка файлов

3. **Интеграция:**
   - ✅ Загрузить файл
   - ✅ Обработать файл
   - ✅ Скачать результат
   - ✅ Проверить историю
   - ✅ Скачать из истории

---

## 📊 Что изменилось

### Frontend
- `index.html` - новые стили и разметка
- `app.js` - полностью переписан с новыми фичами
- Добавлены: прогресс-бар, toast, история, предпросмотр

### Backend
- `cmd/server/main.go` - очищен от admin handlers
- CORS обновлен (добавлен localhost:3000)

### Документация
- `ROADMAP.md` - новая дорожная карта
- `docs/PHASE2_FEATURES.md` - описание фич
- `PHASE2_SUMMARY.md` - итоги Phase 2

---

## 🐛 Troubleshooting

### Frontend не обновляется
```bash
# Очистить CloudFront кеш
aws cloudfront create-invalidation \
  --distribution-id E1GG21S86PMA83 \
  --paths "/*"

# Проверить в браузере (Ctrl+Shift+R для hard refresh)
```

### Backend не запускается
```bash
# Проверить логи ECS
aws logs tail /ecs/excel-flow --since 5m --format short

# Проверить статус сервиса
aws ecs describe-services \
  --cluster excel-flow-cluster \
  --services excel-flow-service
```

### CORS ошибки
- Проверить, что origin добавлен в `allowedOrigins`
- Проверить, что backend перезапущен после изменений

---

## 📝 Команды для быстрого доступа

```bash
# Логи backend
aws logs tail /ecs/excel-flow --follow

# Статус ECS
aws ecs describe-services \
  --cluster excel-flow-cluster \
  --services excel-flow-service \
  --query 'services[0].deployments'

# Список файлов в S3
aws s3 ls s3://excel-viazov-dev/

# Инвалидация CloudFront
aws cloudfront create-invalidation \
  --distribution-id E1GG21S86PMA83 \
  --paths "/*"
```

---

## 🎉 После успешного деплоя

1. Протестировать все новые фичи
2. Обновить STATUS.md с новой версией
3. Создать git tag: `git tag v2.1 && git push --tags`
4. Начать Phase 4 (Admin Panel)

---

## 📞 Контакты

- **Frontend:** https://excel.viazov.dev
- **Backend API:** https://api.viazov.dev
- **Health Check:** https://api.viazov.dev/health
- **CloudFront:** E1GG21S86PMA83
- **S3 Bucket:** excel-viazov-dev
- **ECS Cluster:** excel-flow-cluster
- **ECS Service:** excel-flow-service
