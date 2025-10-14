#!/bin/bash
set -e

echo "🚀 Deploying Excel Flow to AWS..."

# Check AWS credentials
if ! aws sts get-caller-identity &>/dev/null; then
    echo "❌ AWS credentials not configured"
    echo "Please run: aws configure"
    echo "Or set AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables"
    exit 1
fi

# Variables
AWS_REGION=${AWS_REGION:-us-east-1}
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
ECR_REPO="excel-flow"
IMAGE_TAG="latest"

echo "📦 Building Docker image..."
docker build -t $ECR_REPO:$IMAGE_TAG .

echo "🔐 Logging in to ECR..."
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com

echo "📦 Ensuring ECR repository exists..."
aws ecr describe-repositories --repository-names $ECR_REPO --region $AWS_REGION 2>/dev/null || \
  aws ecr create-repository --repository-name $ECR_REPO --region $AWS_REGION

echo "⬆️ Pushing image to ECR..."
docker tag $ECR_REPO:$IMAGE_TAG $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPO:$IMAGE_TAG
docker push $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPO:$IMAGE_TAG

echo "🔄 Updating ECS service..."
aws ecs update-service --cluster excel-flow-cluster --service excel-flow-service --force-new-deployment --region $AWS_REGION

echo "✅ Deployment complete!"
echo "🌐 Application: https://excel.viazov.dev"
