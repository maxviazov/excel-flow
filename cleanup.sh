#!/bin/bash
set -e

echo "🧹 Cleaning up existing infrastructure..."
cd terraform

# Create dummy tfvars if not exists
if [ ! -f terraform.tfvars ]; then
    cat > terraform.tfvars << EOF
aws_region       = "us-east-1"
domain_name      = "viazov.dev"
route53_zone_id  = "dummy"
EOF
fi

terraform destroy -auto-approve || true
rm -f terraform.tfvars
cd ..

echo "✅ Cleanup complete!"
echo "Now you can run: ./deploy-full.sh"
