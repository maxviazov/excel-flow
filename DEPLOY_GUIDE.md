# 🚀 Руководство по деплою на viazov.dev

## Предварительные требования

1. ✅ AWS CLI установлен и настроен
2. ✅ Terraform установлен
3. ✅ Docker установлен и запущен
4. ✅ Доступ к DNS настройкам viazov.dev

## Архитектура

```
┌─────────────────────────────────────────────────────────┐
│                     viazov.dev                          │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  excel.viazov.dev          api.viazov.dev              │
│  (Frontend)                (Backend API)                │
│       │                          │                      │
│       ▼                          ▼                      │
│  CloudFront                    ALB                      │
│       │                          │                      │
│       ▼                          ▼                      │
│   S3 Bucket                ECS Fargate                  │
│   (Static)                 (Docker)                     │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

## Быстрый деплой

```bash
./deploy-full.sh
```

Скрипт автоматически:
1. Создаст всю инфраструктуру в AWS
2. Выведет DNS записи для добавления
3. Соберет и задеплоит backend
4. Задеплоит frontend
5. Настроит SSL сертификаты

## Пошаговая инструкция

### Шаг 1: Создание инфраструктуры

```bash
cd terraform
terraform init
terraform plan
terraform apply
```

Terraform создаст:
- VPC с публичными подсетями
- ECS Cluster + Fargate
- Application Load Balancer
- ECR Repository
- S3 Bucket для фронтенда
- CloudFront Distribution
- ACM Certificates для SSL

### Шаг 2: Настройка DNS

После `terraform apply` вы получите DNS записи. Добавьте их в viazov.dev:

#### Для api.viazov.dev:
```
CNAME _xxx.api.viazov.dev -> _xxx.acm-validations.aws
CNAME api.viazov.dev -> excel-flow-alb-xxx.us-east-1.elb.amazonaws.com
```

#### Для excel.viazov.dev:
```
CNAME _xxx.excel.viazov.dev -> _xxx.acm-validations.aws
CNAME excel.viazov.dev -> dxxx.cloudfront.net
```

**Важно**: Дождитесь валидации сертификатов (5-10 минут)

### Шаг 3: Деплой Backend

```bash
# Получить URL ECR
ECR_URL=$(cd terraform && terraform output -raw ecr_repository_url)

# Логин в ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin $ECR_URL

# Собрать и загрузить образ
docker build -t excel-flow .
docker tag excel-flow:latest $ECR_URL:latest
docker push $ECR_URL:latest

# Обновить ECS service
aws ecs update-service --cluster excel-flow-cluster --service excel-flow-service --force-new-deployment --region us-east-1
```

### Шаг 4: Деплой Frontend

```bash
# Обновить config.js
cat > frontend/public/config.js << 'EOF'
const API_BASE_URL = window.location.hostname === 'localhost' 
    ? 'http://localhost:8080'
    : 'https://api.viazov.dev';
EOF

# Загрузить в S3
S3_BUCKET=$(cd terraform && terraform output -raw s3_bucket_name)
cd frontend
aws s3 sync . s3://$S3_BUCKET --delete
cd ..

# Инвалидировать CloudFront кэш
DISTRIBUTION_ID=$(aws cloudfront list-distributions --query "DistributionList.Items[?Aliases.Items[?contains(@, 'excel.viazov.dev')]].Id" --output text)
aws cloudfront create-invalidation --distribution-id $DISTRIBUTION_ID --paths "/*"
```

## Проверка деплоя

### Backend API
```bash
curl https://api.viazov.dev/health
# Ожидается: {"status":"ok"}
```

### Frontend
Откройте в браузере:
- https://excel.viazov.dev/public - главное приложение
- https://excel.viazov.dev/admin - админ-панель

## Обновление

### Backend
```bash
# Пересобрать и загрузить образ
docker build -t excel-flow .
docker tag excel-flow:latest $ECR_URL:latest
docker push $ECR_URL:latest

# Перезапустить service
aws ecs update-service --cluster excel-flow-cluster --service excel-flow-service --force-new-deployment --region us-east-1
```

### Frontend
```bash
cd frontend
aws s3 sync . s3://$S3_BUCKET --delete
aws cloudfront create-invalidation --distribution-id $DISTRIBUTION_ID --paths "/*"
```

## Мониторинг

### Логи Backend
```bash
aws logs tail /ecs/excel-flow --follow --region us-east-1
```

### Метрики ECS
```bash
aws ecs describe-services --cluster excel-flow-cluster --services excel-flow-service --region us-east-1
```

### CloudFront статистика
AWS Console → CloudFront → Distributions → Monitoring

## Стоимость

- **ECS Fargate**: ~$15-30/месяц (0.25 vCPU, 0.5GB RAM)
- **ALB**: ~$16/месяц
- **S3**: ~$0.50/месяц (хранение + запросы)
- **CloudFront**: ~$1-2/месяц (трафик)
- **Route53**: $0.50/месяц (hosted zone)
- **Итого**: ~$33-50/месяц

## Удаление инфраструктуры

```bash
# Удалить образы из ECR
aws ecr batch-delete-image --repository-name excel-flow --image-ids imageTag=latest --region us-east-1

# Очистить S3
aws s3 rm s3://$S3_BUCKET --recursive

# Удалить через Terraform
cd terraform
terraform destroy
```

## Troubleshooting

### ECS task не запускается
```bash
# Проверить логи
aws logs tail /ecs/excel-flow --follow

# Проверить task definition
aws ecs describe-task-definition --task-definition excel-flow
```

### CloudFront возвращает 403
- Проверьте S3 bucket policy
- Проверьте что файлы загружены
- Инвалидируйте кэш

### SSL сертификат не валидируется
- Проверьте DNS записи
- Подождите 5-10 минут
- Проверьте в ACM Console

## Поддержка

Документация:
- `ARCHITECTURE.md` - архитектура системы
- `TEST_REPORT.md` - результаты тестирования
- `SEPARATION_GUIDE.md` - разделение на backend/frontend
