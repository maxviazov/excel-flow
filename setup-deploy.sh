#!/bin/bash
set -e

echo "🚀 Excel Flow Deployment Setup"
echo ""

# Check AWS credentials
echo "1️⃣ Checking AWS credentials..."
if ! aws sts get-caller-identity &>/dev/null; then
    echo "❌ AWS credentials not configured"
    echo "Run: aws configure"
    exit 1
fi
echo "✅ AWS credentials OK"
echo ""

# Get Zone ID
echo "2️⃣ Getting Route53 Zone ID for viazov.dev..."
ZONE_ID=$(aws route53 list-hosted-zones --query "HostedZones[?Name=='viazov.dev.'].Id" --output text | cut -d'/' -f3)

if [ -z "$ZONE_ID" ]; then
    echo "❌ Zone viazov.dev not found in Route53"
    echo "Create hosted zone first or use different domain"
    exit 1
fi
echo "✅ Found Zone ID: $ZONE_ID"
echo ""

# Create terraform.tfvars
echo "3️⃣ Creating terraform/terraform.tfvars..."
cat > terraform/terraform.tfvars << EOF
aws_region       = "us-east-1"
domain_name      = "excel.viazov.dev"
route53_zone_id  = "$ZONE_ID"
EOF
echo "✅ Created terraform.tfvars"
echo ""

# Show next steps
echo "📋 Next steps:"
echo "1. Review terraform/terraform.tfvars"
echo "2. Run: ./deploy.sh"
echo ""
echo "🎯 Your app will be available at: https://excel.viazov.dev"
