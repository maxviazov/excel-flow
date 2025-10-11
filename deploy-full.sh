#!/bin/bash
set -e

echo "🚀 Full Deployment: viazov.dev"
echo ""

# Variables
DOMAIN="viazov.dev"
AWS_REGION=${AWS_REGION:-us-east-1}
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
ECR_REPO="excel-flow"
IMAGE_TAG="latest"

# Step 1: Skip Route53 - using Cloudflare DNS
echo "1️⃣ Using Cloudflare DNS for $DOMAIN..."
echo "✅ DNS managed by Cloudflare"
echo ""

# Step 2: Build and push Docker image
echo "2️⃣ Building Docker image..."
docker build -t $ECR_REPO:$IMAGE_TAG .
echo ""

echo "3️⃣ Logging in to ECR..."
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com
echo ""

echo "4️⃣ Creating ECR repository..."
aws ecr describe-repositories --repository-names $ECR_REPO --region $AWS_REGION 2>/dev/null || \
  aws ecr create-repository --repository-name $ECR_REPO --region $AWS_REGION
echo ""

echo "5️⃣ Pushing image to ECR..."
docker tag $ECR_REPO:$IMAGE_TAG $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPO:$IMAGE_TAG
docker push $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPO:$IMAGE_TAG
echo ""

# Step 3: Create terraform.tfvars (without Route53)
echo "6️⃣ Creating Terraform configuration..."
cat > terraform/terraform.tfvars << EOF
aws_region       = "$AWS_REGION"
domain_name      = "$DOMAIN"
EOF
echo "✅ Created terraform/terraform.tfvars"
echo ""

# Step 4: Deploy with Terraform
echo "7️⃣ Deploying infrastructure with Terraform..."
cd terraform
terraform init -upgrade
terraform apply -auto-approve
cd ..
echo ""

# Step 5: Get outputs and show Cloudflare instructions
ALB_DNS=$(cd terraform && terraform output -raw alb_dns && cd ..)

echo "✅ Deployment complete!"
echo ""
echo "🔗 ALB URL: $ALB_DNS"
echo ""
echo "📋 Next steps in Cloudflare:"
echo "1. Go to DNS → Records → Add record"
echo "2. Type: CNAME"
echo "3. Name: @ (or viazov.dev)"
echo "4. Target: $ALB_DNS"
echo "5. Proxy status: Proxied (orange cloud)"
echo "6. Click Save"
echo ""
echo "⏱️  Wait 2-3 minutes, then open: https://$DOMAIN"
