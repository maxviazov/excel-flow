# 🔧 Быстрое исправление деплоя

## Проблема

SSL сертификаты требуют DNS валидации ДО того, как их можно использовать в ALB/CloudFront.

## Решение: Поэтапный деплой

### Шаг 1: Очистка текущего состояния

```bash
cd terraform
terraform destroy -auto-approve
cd ..
```

### Шаг 2: Временно отключить HTTPS

Закомментируйте в `terraform/main.tf`:
- `aws_lb_listener.https` (строки ~295-310)
- `aws_acm_certificate.api` (строки ~280-295)
- `aws_acm_certificate.frontend` (строки ~390-405)
- `aws_cloudfront_distribution.frontend` (строки ~410-470)

### Шаг 3: Деплой без SSL

```bash
cd terraform
terraform init
terraform apply -auto-approve
cd ..
```

### Шаг 4: Получить outputs

```bash
cd terraform
terraform output
cd ..
```

Вы получите:
- ALB DNS name
- ECR repository URL

### Шаг 5: Деплой backend

```bash
# Логин в ECR
ECR_URL=$(cd terraform && terraform output -raw ecr_repository_url)
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin $ECR_URL

# Собрать и загрузить
docker build -t excel-flow .
docker tag excel-flow:latest $ECR_URL:latest
docker push $ECR_URL:latest

# Обновить ECS
aws ecs update-service --cluster excel-flow-cluster --service excel-flow-service --force-new-deployment --region us-east-1
```

### Шаг 6: Деплой frontend

```bash
# Обновить config с HTTP (временно)
cat > frontend/public/config.js << 'EOF'
const API_BASE_URL = window.location.hostname === 'localhost' 
    ? 'http://localhost:8080'
    : 'http://ALB_DNS_NAME_HERE';  // Замените на ваш ALB DNS
EOF

# Загрузить в S3
cd frontend
aws s3 sync . s3://excel-viazov-dev --delete
cd ..
```

### Шаг 7: Проверка

```bash
# Backend
curl http://ALB_DNS_NAME/health

# Frontend
curl http://excel-viazov-dev.s3-website-us-east-1.amazonaws.com
```

## Альтернатива: Упрощенный деплой

Используйте существующую инфраструктуру без доменов:

1. Backend: ALB DNS напрямую
2. Frontend: S3 website endpoint напрямую

Это будет работать без SSL, но функционально полностью.

## Добавление SSL позже

После того как все работает:

1. Создайте ACM сертификаты
2. Добавьте DNS записи для валидации
3. Дождитесь валидации (5-10 минут)
4. Раскомментируйте HTTPS listener
5. Запустите `terraform apply`
6. Обновите DNS: api.viazov.dev → ALB, excel.viazov.dev → CloudFront

---

**Рекомендация**: Начните с HTTP версии, убедитесь что все работает, потом добавьте SSL.
