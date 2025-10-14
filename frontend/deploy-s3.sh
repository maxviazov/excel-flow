#!/bin/bash

# Скрипт деплоя фронтенда на S3

BUCKET_NAME="excel-flow-frontend"
REGION="us-east-1"

echo "🚀 Deploying frontend to S3..."

# Проверка наличия bucket
if ! aws s3 ls "s3://$BUCKET_NAME" 2>&1 | grep -q 'NoSuchBucket'; then
    echo "✅ Bucket exists"
else
    echo "📦 Creating bucket..."
    aws s3 mb "s3://$BUCKET_NAME" --region "$REGION"
fi

# Настройка статического хостинга
echo "⚙️ Configuring static website hosting..."
aws s3 website "s3://$BUCKET_NAME" \
    --index-document public/index.html \
    --error-document public/index.html

# Загрузка файлов
echo "📤 Uploading files..."
aws s3 sync . "s3://$BUCKET_NAME" \
    --exclude ".git/*" \
    --exclude "*.sh" \
    --exclude "README.md" \
    --cache-control "public, max-age=31536000" \
    --delete

# Настройка публичного доступа
echo "🔓 Setting public access..."
aws s3api put-bucket-policy --bucket "$BUCKET_NAME" --policy "{
  \"Version\": \"2012-10-17\",
  \"Statement\": [{
    \"Sid\": \"PublicReadGetObject\",
    \"Effect\": \"Allow\",
    \"Principal\": \"*\",
    \"Action\": \"s3:GetObject\",
    \"Resource\": \"arn:aws:s3:::$BUCKET_NAME/*\"
  }]
}"

# Получение URL
WEBSITE_URL="http://$BUCKET_NAME.s3-website-$REGION.amazonaws.com"

echo ""
echo "✅ Deployment complete!"
echo "🌐 Website URL: $WEBSITE_URL"
echo ""
echo "📝 Don't forget to update config.js with your API URL!"
