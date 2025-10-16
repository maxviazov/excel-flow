#!/bin/bash
set -e

echo "🚀 Deploying Excel Flow to AWS..."

# Variables
AWS_REGION=${AWS_REGION:-us-east-1}
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
ECR_REPO="excel-flow"
IMAGE_TAG="latest"

# Build Docker image
echo "📦 Building Docker image..."
docker buildx build --platform linux/amd64 -t $ECR_REPO:$IMAGE_TAG --load .

# Login to ECR
echo "🔐 Logging in to ECR..."
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com

# Create ECR repository if it doesn't exist
echo "📦 Creating ECR repository..."
aws ecr describe-repositories --repository-names $ECR_REPO --region $AWS_REGION 2>/dev/null || \
  aws ecr create-repository --repository-name $ECR_REPO --region $AWS_REGION

# Tag and push image
echo "⬆️ Pushing image to ECR..."
docker tag $ECR_REPO:$IMAGE_TAG $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPO:$IMAGE_TAG
docker push $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPO:$IMAGE_TAG

# Deploy with Terraform
echo "🏗️ Deploying infrastructure with Terraform..."
cd terraform
terraform init
terraform apply -auto-approve

echo "✅ Deployment complete!"
echo "🌐 Your application will be available at your domain once DNS propagates"
