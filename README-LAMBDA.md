# Деплой на AWS Lambda

## 🚀 Быстрый старт

### 1. Настройка AWS credentials
```bash
aws configure
# Введите AWS Access Key ID, Secret Access Key, region (us-east-1)
```

### 2. Деплой инфраструктуры
```bash
./deploy-lambda.sh
```

### 3. Настройка DNS
После деплоя Terraform выведет записи для валидации SSL сертификатов:
```
certificate_validation_api = {
  "fish.viazov.dev" = {
    name  = "_xxx.fish.viazov.dev"
    type  = "CNAME"
    value = "yyy.acm-validations.aws."
  }
}
```

Добавьте эти CNAME записи в DNS настройки домена.

### 4. Направление домена на CloudFront
После валидации сертификатов добавьте CNAME запись:
```
fish.viazov.dev -> d1234567890.cloudfront.net
```

### 5. Деплой фронтенда
```bash
cd frontend && ./deploy-s3-lambda.sh
```

## 💰 Стоимость

- **Lambda**: $0 за первый 1M запросов в месяц
- **API Gateway**: $0 за первый 1M запросов в месяц  
- **S3**: ~$0.02 за GB в месяц
- **CloudFront**: $0.085 за первые 10TB трафика

**Итого**: ~$1-2 в месяц при редком использовании

## 🔧 Локальная разработка

```bash
# Запуск локального сервера
go run cmd/server/main.go

# Тестирование Lambda функции локально
go run cmd/lambda/main.go
```

## 📁 Структура

- `cmd/lambda/` - Lambda функция
- `terraform/lambda.tf` - инфраструктура AWS
- `frontend/` - статические файлы для S3
- `deploy-lambda.sh` - скрипт деплоя

## 🌐 URL

После деплоя приложение будет доступно по адресу:
**https://fish.viazov.dev**