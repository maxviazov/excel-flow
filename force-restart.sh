#!/bin/bash
set -e

AWS_REGION=${AWS_REGION:-us-east-1}

echo "🔄 Принудительный перезапуск сервиса..."

# Останавливаем все задачи
echo "🛑 Останавливаем все задачи..."
TASKS=$(aws ecs list-tasks --cluster excel-flow-cluster --service-name excel-flow-service --region $AWS_REGION --query 'taskArns[]' --output text)
for TASK in $TASKS; do
    echo "  Останавливаем задачу: $TASK"
    aws ecs stop-task --cluster excel-flow-cluster --task $TASK --region $AWS_REGION --no-cli-pager > /dev/null 2>&1 || true
done

# Ждем 5 секунд
echo "⏳ Ожидание 5 секунд..."
sleep 5

# Принудительное обновление сервиса
echo "🔄 Принудительное обновление сервиса..."
aws ecs update-service \
    --cluster excel-flow-cluster \
    --service excel-flow-service \
    --force-new-deployment \
    --region $AWS_REGION \
    --no-cli-pager

echo ""
echo "✅ Перезапуск инициирован!"
echo "⏳ Подождите 2-3 минуты и проверьте статус"
