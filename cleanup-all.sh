#!/bin/bash
set -e

echo "🗑️ Complete cleanup: AWS + Docker + Terraform state"

# Unset AWS config
unset AWS_CONFIG_FILE AWS_SHARED_CREDENTIALS_FILE

# 1. Delete ECS service
echo "1. Deleting ECS service..."
aws ecs delete-service --cluster excel-flow-cluster --service excel-flow-service --region us-east-1 --force 2>/dev/null || echo "  ✓ Service already deleted"

# 2. Wait and delete cluster
echo "2. Waiting 30s for service deletion..."
sleep 30

echo "3. Deleting ECS cluster..."
aws ecs delete-cluster --cluster excel-flow-cluster --region us-east-1 2>/dev/null || echo "  ✓ Cluster already deleted"

# 4. Delete task definitions
echo "4. Deregistering task definitions..."
TASK_ARNS=$(aws ecs list-task-definitions --family-prefix excel-flow --region us-east-1 --query 'taskDefinitionArns' --output text 2>/dev/null || echo "")
for arn in $TASK_ARNS; do
  aws ecs deregister-task-definition --task-definition $arn --region us-east-1 2>/dev/null || true
done

# 5. Delete ALB
echo "5. Deleting ALB..."
ALB_ARN=$(aws elbv2 describe-load-balancers --region us-east-1 --query "LoadBalancers[?LoadBalancerName=='excel-flow-alb'].LoadBalancerArn" --output text 2>/dev/null || echo "")
if [ ! -z "$ALB_ARN" ]; then
  aws elbv2 delete-load-balancer --load-balancer-arn $ALB_ARN --region us-east-1
  echo "  Waiting for ALB deletion..."
  sleep 30
else
  echo "  ✓ ALB already deleted"
fi

# 6. Delete target group
echo "6. Deleting target group..."
TG_ARN=$(aws elbv2 describe-target-groups --region us-east-1 --query "TargetGroups[?TargetGroupName=='excel-flow-tg'].TargetGroupArn" --output text 2>/dev/null || echo "")
if [ ! -z "$TG_ARN" ]; then
  aws elbv2 delete-target-group --target-group-arn $TG_ARN --region us-east-1 2>/dev/null || echo "  ✓ Target group already deleted"
else
  echo "  ✓ Target group already deleted"
fi

# 7. Delete security groups
echo "7. Deleting security groups..."
for sg in $(aws ec2 describe-security-groups --region us-east-1 --query "SecurityGroups[?starts_with(GroupName, 'excel-flow')].GroupId" --output text 2>/dev/null || echo ""); do
  if [ ! -z "$sg" ]; then
    aws ec2 delete-security-group --group-id $sg --region us-east-1 2>/dev/null || echo "  ✓ SG $sg already deleted"
  fi
done

# 8. Delete VPC resources
echo "8. Deleting VPC resources..."
for vpc in $(aws ec2 describe-vpcs --region us-east-1 --query "Vpcs[?Tags[?Key=='Name' && Value=='excel-flow-vpc']].VpcId" --output text 2>/dev/null || echo ""); do
  if [ ! -z "$vpc" ]; then
    # Delete subnets
    for subnet in $(aws ec2 describe-subnets --region us-east-1 --filters "Name=vpc-id,Values=$vpc" --query "Subnets[].SubnetId" --output text 2>/dev/null || echo ""); do
      aws ec2 delete-subnet --subnet-id $subnet --region us-east-1 2>/dev/null || true
    done
    
    # Delete route tables
    for rt in $(aws ec2 describe-route-tables --region us-east-1 --filters "Name=vpc-id,Values=$vpc" --query "RouteTables[?Associations[0].Main==\`false\`].RouteTableId" --output text 2>/dev/null || echo ""); do
      aws ec2 delete-route-table --route-table-id $rt --region us-east-1 2>/dev/null || true
    done
    
    # Detach and delete IGW
    for igw in $(aws ec2 describe-internet-gateways --region us-east-1 --filters "Name=attachment.vpc-id,Values=$vpc" --query "InternetGateways[].InternetGatewayId" --output text 2>/dev/null || echo ""); do
      aws ec2 detach-internet-gateway --internet-gateway-id $igw --vpc-id $vpc --region us-east-1 2>/dev/null || true
      aws ec2 delete-internet-gateway --internet-gateway-id $igw --region us-east-1 2>/dev/null || true
    done
    
    # Delete VPC
    aws ec2 delete-vpc --vpc-id $vpc --region us-east-1 2>/dev/null || echo "  ✓ VPC already deleted"
  fi
done

# 9. Delete IAM role
echo "9. Deleting IAM role..."
aws iam detach-role-policy --role-name excel-flow-ecs-execution --policy-arn arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy --region us-east-1 2>/dev/null || true
aws iam delete-role --role-name excel-flow-ecs-execution --region us-east-1 2>/dev/null || echo "  ✓ IAM role already deleted"

# 10. Delete CloudWatch log group
echo "10. Deleting CloudWatch log group..."
aws logs delete-log-group --log-group-name /ecs/excel-flow --region us-east-1 2>/dev/null || echo "  ✓ Log group already deleted"

# 11. Delete ACM certificate
echo "11. Deleting ACM certificate..."
CERT_ARN=$(aws acm list-certificates --region us-east-1 --query "CertificateSummaryList[?DomainName=='viazov.dev'].CertificateArn" --output text 2>/dev/null || echo "")
if [ ! -z "$CERT_ARN" ]; then
  aws acm delete-certificate --certificate-arn $CERT_ARN --region us-east-1 2>/dev/null || echo "  ✓ Certificate already deleted"
else
  echo "  ✓ Certificate already deleted"
fi

# 12. Clean Terraform state
echo "12. Cleaning Terraform state..."
rm -f terraform/terraform.tfstate* terraform/.terraform.lock.hcl

# 13. Clean Docker images
echo "13. Cleaning Docker images..."
docker rmi excel-flow:latest 2>/dev/null || echo "  ✓ Local image already deleted"
docker rmi 138008497687.dkr.ecr.us-east-1.amazonaws.com/excel-flow:latest 2>/dev/null || echo "  ✓ ECR image already deleted"

echo ""
echo "✅ Complete cleanup finished!"
echo ""
echo "Next steps:"
echo "  1. Run: bash deploy-full.sh"
echo "  2. Wait 5-10 minutes for deployment"
echo "  3. Update Cloudflare DNS"
