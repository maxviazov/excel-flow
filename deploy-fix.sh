#!/bin/bash

echo "🚀 Deploying Excel Flow with city aliases fix..."

# Build and push Docker image
echo "📦 Building Docker image..."
docker build --no-cache -t excel-flow:latest .

echo "🏷️ Tagging image..."
docker tag excel-flow:latest 138008497687.dkr.ecr.us-east-1.amazonaws.com/excel-flow:latest

echo "🔐 Logging into ECR..."
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 138008497687.dkr.ecr.us-east-1.amazonaws.com

echo "⬆️ Pushing image..."
docker push 138008497687.dkr.ecr.us-east-1.amazonaws.com/excel-flow:latest

echo "🔄 Updating ECS service..."
aws ecs update-service --cluster excel-flow-cluster --service excel-flow-service --force-new-deployment --region us-east-1

echo "✅ Deployment initiated! Check AWS console for status."
echo "🌐 App will be available at: https://api.viazov.dev"