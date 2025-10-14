# Исправление excel.viazov.dev

## Проблема
excel.viazov.dev возвращает старый config.js через Cloudflare кеш.

## Решение

### Вариант 1: Очистить кеш Cloudflare (быстро)

1. Откройте Cloudflare Dashboard
2. Выберите домен viazov.dev
3. Перейдите в Caching → Configuration
4. Нажмите "Purge Everything" или "Purge by URL"
5. Если по URL, укажите: `https://excel.viazov.dev/config.js`

### Вариант 2: Отключить прокси Cloudflare (рекомендуется)

1. Откройте Cloudflare Dashboard → DNS
2. Найдите запись `excel` CNAME
3. Кликните на оранжевое облако (Proxied)
4. Переключите на серое облако (DNS only)
5. Сохраните

Это отключит кеширование Cloudflare и запросы пойдут напрямую к CloudFront.

### Вариант 3: Использовать прямой URL (работает сейчас)

Используйте S3 напрямую:
```
http://excel-viazov-dev.s3-website-us-east-1.amazonaws.com
```

Или CloudFront (после обновления контейнера):
```
https://d18sq2gf3s7zhe.cloudfront.net
```

## Текущий статус

- ✅ Backend обновлен с правильными заголовками МЗ
- ✅ Docker образ задеплоен
- ⏳ ECS контейнер перезапускается (~2 минуты)
- ⚠️ excel.viazov.dev требует очистки кеша Cloudflare

## Проверка после деплоя

Подождите 2 минуты и проверьте:
```bash
# Проверить API
curl http://excel-flow-alb-1086104942.us-east-1.elb.amazonaws.com/health

# Загрузить тестовый файл
curl -X POST -F "file=@testdata/משקל.xlsx" \
  http://excel-flow-alb-1086104942.us-east-1.elb.amazonaws.com/api/upload
```

Затем откройте:
```
http://excel-viazov-dev.s3-website-us-east-1.amazonaws.com
```
