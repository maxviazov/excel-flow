#!/bin/bash
set -e

echo "🚀 Excel Flow - Full Deployment"
echo "================================"
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Step 1: Terraform Init
echo -e "${BLUE}Step 1: Initializing Terraform...${NC}"
cd terraform
terraform init
echo -e "${GREEN}✅ Terraform initialized${NC}"
echo ""

# Step 2: Terraform Plan
echo -e "${BLUE}Step 2: Planning infrastructure...${NC}"
terraform plan -out=tfplan
echo -e "${GREEN}✅ Plan created${NC}"
echo ""

# Step 3: Terraform Apply
echo -e "${BLUE}Step 3: Creating AWS resources...${NC}"
terraform apply tfplan
echo -e "${GREEN}✅ Infrastructure created${NC}"
echo ""

# Get outputs
ECR_URL=$(terraform output -raw ecr_repository_url)
API_CERT=$(terraform output -json api_certificate_validation)
FRONTEND_CERT=$(terraform output -json frontend_certificate_validation)
S3_BUCKET=$(terraform output -raw s3_bucket_name)
CLOUDFRONT_DOMAIN=$(terraform output -raw cloudfront_domain)

cd ..

echo ""
echo -e "${YELLOW}⚠️  IMPORTANT: DNS Configuration Required${NC}"
echo "=========================================="
echo ""
echo "Add these DNS records to viazov.dev:"
echo ""
echo "1. API Certificate Validation:"
echo "$API_CERT" | jq -r 'to_entries[] | "   \(.value.type) \(.value.name) -> \(.value.value)"'
echo ""
echo "2. Frontend Certificate Validation:"
echo "$FRONTEND_CERT" | jq -r 'to_entries[] | "   \(.value.type) \(.value.name) -> \(.value.value)"'
echo ""
echo "3. API Domain:"
echo "   CNAME api.viazov.dev -> $(cd terraform && terraform output -raw alb_dns_name)"
echo ""
echo "4. Frontend Domain:"
echo "   CNAME excel.viazov.dev -> $CLOUDFRONT_DOMAIN"
echo ""
echo "Press Enter after adding DNS records..."
read

# Step 4: Build and Push Docker Image
echo -e "${BLUE}Step 4: Building Docker image...${NC}"
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin $ECR_URL
docker build -t excel-flow .
docker tag excel-flow:latest $ECR_URL:latest
docker push $ECR_URL:latest
echo -e "${GREEN}✅ Docker image pushed${NC}"
echo ""

# Step 5: Update ECS Service
echo -e "${BLUE}Step 5: Updating ECS service...${NC}"
aws ecs update-service --cluster excel-flow-cluster --service excel-flow-service --force-new-deployment --region us-east-1
echo -e "${GREEN}✅ ECS service updated${NC}"
echo ""

# Step 6: Update Frontend Config
echo -e "${BLUE}Step 6: Updating frontend config...${NC}"
cat > frontend/public/config.js << EOF
// API Configuration
const API_BASE_URL = window.location.hostname === 'localhost' 
    ? 'http://localhost:8080'
    : 'https://api.viazov.dev';
EOF
echo -e "${GREEN}✅ Frontend config updated${NC}"
echo ""

# Step 7: Deploy Frontend to S3
echo -e "${BLUE}Step 7: Deploying frontend to S3...${NC}"
cd frontend
aws s3 sync . s3://$S3_BUCKET \
    --exclude ".git/*" \
    --exclude "*.sh" \
    --exclude "README.md" \
    --cache-control "public, max-age=31536000" \
    --delete
cd ..
echo -e "${GREEN}✅ Frontend deployed${NC}"
echo ""

# Step 8: Invalidate CloudFront Cache
echo -e "${BLUE}Step 8: Invalidating CloudFront cache...${NC}"
DISTRIBUTION_ID=$(aws cloudfront list-distributions --query "DistributionList.Items[?Aliases.Items[?contains(@, 'excel.viazov.dev')]].Id" --output text)
aws cloudfront create-invalidation --distribution-id $DISTRIBUTION_ID --paths "/*"
echo -e "${GREEN}✅ Cache invalidated${NC}"
echo ""

echo ""
echo -e "${GREEN}🎉 Deployment Complete!${NC}"
echo "======================="
echo ""
echo "🌐 API: https://api.viazov.dev"
echo "🌐 Frontend: https://excel.viazov.dev"
echo ""
echo "Wait 2-3 minutes for ECS service to start and DNS to propagate."
