#!/bin/bash
set -e

echo "🔄 Обновление Excel Flow на AWS..."

# Переменные
AWS_REGION=${AWS_REGION:-us-east-1}
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text 2>/dev/null || echo "")
ECR_REPO="excel-flow"
IMAGE_TAG="$(date +%Y%m%d-%H%M%S)"
ECR_URI="$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPO"

if [ -z "$AWS_ACCOUNT_ID" ]; then
    echo "❌ AWS credentials не настроены"
    echo "Запустите: aws configure"
    exit 1
fi

# Сборка Docker образа
echo "📦 Сборка Docker образа..."
docker build -t $ECR_REPO:$IMAGE_TAG .

# Логин в ECR
echo "🔐 Логин в ECR..."
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com

# Тегирование и пуш образа
echo "⬆️  Загрузка образа в ECR (тег: $IMAGE_TAG)..."
docker tag $ECR_REPO:$IMAGE_TAG $ECR_URI:$IMAGE_TAG
docker tag $ECR_REPO:$IMAGE_TAG $ECR_URI:latest
docker push $ECR_URI:$IMAGE_TAG
docker push $ECR_URI:latest

# Получаем текущую task definition
echo "📋 Получение текущей task definition..."
TASK_DEF=$(aws ecs describe-task-definition --task-definition excel-flow --region $AWS_REGION)
NEW_TASK_DEF=$(echo $TASK_DEF | jq --arg IMAGE "$ECR_URI:$IMAGE_TAG" '.taskDefinition | .containerDefinitions[0].image = $IMAGE | del(.taskDefinitionArn, .revision, .status, .requiresAttributes, .compatibilities, .registeredAt, .registeredBy)')

# Регистрируем новую task definition
echo "📝 Регистрация новой task definition..."
NEW_TASK_ARN=$(aws ecs register-task-definition --region $AWS_REGION --cli-input-json "$NEW_TASK_DEF" --query 'taskDefinition.taskDefinitionArn' --output text)

# Обновляем сервис с новой task definition
echo "🔄 Обновление ECS сервиса..."
aws ecs update-service \
    --cluster excel-flow-cluster \
    --service excel-flow-service \
    --task-definition $NEW_TASK_ARN \
    --force-new-deployment \
    --region $AWS_REGION

# Ждем завершения обновления
echo "⏳ Ожидание запуска новой задачи (30 секунд)..."
sleep 30

echo ""
echo "✅ Обновление запущено!"
echo "🏷️  Новый образ: $ECR_URI:$IMAGE_TAG"
echo "⏳ Новая версия будет развернута через 1-2 минуты"
echo ""
echo "Проверить статус:"
echo "  aws ecs describe-services --cluster excel-flow-cluster --services excel-flow-service --region $AWS_REGION --query 'services[0].{Running:runningCount,TaskDef:taskDefinition}'"
