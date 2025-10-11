# 🚀 Инструкция по деплою Excel Flow на AWS

## Что было сделано

1. ✅ **Исправлена проблема RTL** - Excel файлы теперь отображаются справа налево
2. ✅ **Создан веб-интерфейс** - простой UI на иврите с drag-and-drop
3. ✅ **HTTP API сервер** - для обработки файлов через веб
4. ✅ **Dockerfile** - для контейнеризации приложения
5. ✅ **Terraform** - для автоматического деплоя на AWS

## Быстрый старт (локально)

```bash
# Запустить сервер локально
go run cmd/server/main.go

# Открыть в браузере
open http://localhost:8080
```

## Деплой на AWS

### Предварительные требования

1. AWS CLI установлен и настроен (`aws configure`)
2. Docker установлен
3. Terraform установлен
4. Домен viazov.dev настроен в Route53

### Шаг 1: Получить Zone ID для Route53

```bash
aws route53 list-hosted-zones
```

Найди Zone ID для viazov.dev

### Шаг 2: Настроить переменные

```bash
cd terraform
cp terraform.tfvars.example terraform.tfvars
```

Отредактируй `terraform.tfvars`:
```hcl
aws_region       = "us-east-1"
domain_name      = "excel.viazov.dev"  # или просто viazov.dev
route53_zone_id  = "Z0123456789ABC"    # твой Zone ID
```

### Шаг 3: Деплой

```bash
# Из корня проекта
./deploy.sh
```

Скрипт автоматически:
- Соберет Docker образ
- Загрузит в ECR
- Развернет инфраструктуру через Terraform
- Настроит HTTPS сертификат
- Создаст Route53 запись

### Шаг 4: Подтвердить сертификат

После запуска Terraform, нужно подтвердить SSL сертификат:

1. Зайди в AWS Console → Certificate Manager
2. Найди сертификат для твоего домена
3. Скопируй CNAME записи для валидации
4. Добавь их в Route53 (или Terraform сделает это автоматически)

### Шаг 5: Проверить

```bash
# Получить URL балансировщика
cd terraform
terraform output alb_dns

# Проверить статус
curl https://excel.viazov.dev
```

## Архитектура

```
Internet
   ↓
Route53 (viazov.dev)
   ↓
Application Load Balancer (HTTPS)
   ↓
ECS Fargate (Go сервер)
   ↓
ECR (Docker образ)
```

## Стоимость (примерно)

- **Fargate**: ~$15/месяц (256 CPU, 512 MB RAM)
- **ALB**: ~$16/месяц
- **Route53**: $0.50/месяц
- **ECR**: ~$1/месяц

**Итого**: ~$32-35/месяц

## Обновление приложения

```bash
# Пересобрать и задеплоить
./deploy.sh
```

## Откат

```bash
cd terraform
terraform destroy
```

## Troubleshooting

### Проблема: Сертификат не валидируется

Проверь DNS записи:
```bash
dig excel.viazov.dev
```

### Проблема: ECS задача не запускается

Проверь логи:
```bash
aws logs tail /ecs/excel-flow --follow
```

### Проблема: 502 Bad Gateway

Проверь health check:
```bash
aws elbv2 describe-target-health --target-group-arn <ARN>
```

## Альтернатива: Быстрый деплой без домена

Если нужно показать шефу прямо сейчас без настройки домена:

```bash
# Только инфраструктура без Route53
cd terraform
terraform apply -target=aws_lb.main -target=aws_ecs_service.app

# Получить URL
terraform output alb_dns
# Используй этот URL: http://excel-flow-alb-123456789.us-east-1.elb.amazonaws.com
```

## Что показать шефу

1. Открыть https://excel.viazov.dev (или ALB URL)
2. Загрузить Excel файл (drag-and-drop)
3. Нажать "התחל עיבוד"
4. Показать лог процесса
5. Показать результат (количество строк до/после)
6. Скачать итоговый файл
7. Открыть в Excel - показать что RTL работает правильно!

Удачи с демо! 🚀
