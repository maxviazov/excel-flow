#!/bin/bash

set -e

BUCKET_NAME="excel-viazov-dev"

echo "🚀 Deploying frontend to S3..."

# Sync files to S3
aws s3 sync public/ s3://$BUCKET_NAME/ --delete

# Set proper content types
aws s3 cp s3://$BUCKET_NAME/index.html s3://$BUCKET_NAME/index.html --content-type "text/html; charset=utf-8" --metadata-directive REPLACE
aws s3 cp s3://$BUCKET_NAME/app.js s3://$BUCKET_NAME/app.js --content-type "application/javascript; charset=utf-8" --metadata-directive REPLACE
aws s3 cp s3://$BUCKET_NAME/config.js s3://$BUCKET_NAME/config.js --content-type "application/javascript; charset=utf-8" --metadata-directive REPLACE
aws s3 cp s3://$BUCKET_NAME/i18n.js s3://$BUCKET_NAME/i18n.js --content-type "application/javascript; charset=utf-8" --metadata-directive REPLACE

echo "✅ Frontend deployed to S3!"
echo "🌐 Website URL: https://excel.viazov.dev"