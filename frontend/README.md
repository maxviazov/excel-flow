# Excel Flow - Frontend

Статический фронтенд для деплоя на S3 + CloudFront.

## Структура

```
frontend/
├── public/          # Главное приложение
│   ├── index.html
│   ├── app.js
│   └── config.js    # Конфигурация API
└── admin/           # Админ-панель
    ├── index.html
    ├── app.js
    └── style.css
```

## Конфигурация

Перед деплоем обновите `public/config.js`:

```javascript
const API_BASE_URL = 'https://your-api-domain.com';
```

## Локальная разработка

```bash
# Запустите простой HTTP сервер
python3 -m http.server 3000

# Или используйте npx
npx serve .
```

Откройте:
- http://localhost:3000/public - главное приложение
- http://localhost:3000/admin - админ-панель

## Деплой на S3

### 1. Создайте S3 bucket

```bash
aws s3 mb s3://excel-flow-frontend
```

### 2. Настройте bucket для статического хостинга

```bash
aws s3 website s3://excel-flow-frontend \
  --index-document public/index.html \
  --error-document public/index.html
```

### 3. Загрузите файлы

```bash
aws s3 sync . s3://excel-flow-frontend \
  --exclude ".git/*" \
  --exclude "README.md" \
  --cache-control "public, max-age=31536000"
```

### 4. Настройте публичный доступ

```bash
aws s3api put-bucket-policy --bucket excel-flow-frontend --policy '{
  "Version": "2012-10-17",
  "Statement": [{
    "Sid": "PublicReadGetObject",
    "Effect": "Allow",
    "Principal": "*",
    "Action": "s3:GetObject",
    "Resource": "arn:aws:s3:::excel-flow-frontend/*"
  }]
}'
```

## CloudFront (опционально)

Для HTTPS и CDN создайте CloudFront distribution:

1. Origin: S3 bucket
2. Default root object: `public/index.html`
3. SSL Certificate: ACM certificate
4. Custom domain: `excel-flow.yourdomain.com`

## Стоимость

- S3: ~$0.023/GB хранение + $0.09/GB трафик
- CloudFront: ~$0.085/GB трафик (первые 10TB)
- **Итого**: ~$1-3/месяц для небольшого проекта
