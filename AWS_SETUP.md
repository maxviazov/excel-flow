# 🔐 Настройка AWS для деплоя

## Шаг 1: Получение AWS Credentials

1. Войдите в AWS Console: https://console.aws.amazon.com
2. Перейдите в IAM → Users → Ваш пользователь
3. Security credentials → Create access key
4. Выберите "Command Line Interface (CLI)"
5. Скопируйте Access Key ID и Secret Access Key

## Шаг 2: Настройка AWS CLI

```bash
aws configure
```

Введите:
- **AWS Access Key ID**: ваш Access Key
- **AWS Secret Access Key**: ваш Secret Key
- **Default region name**: `us-east-1`
- **Default output format**: `json`

## Шаг 3: Проверка

```bash
aws sts get-caller-identity
```

Должно вывести:
```json
{
    "UserId": "...",
    "Account": "...",
    "Arn": "arn:aws:iam::..."
}
```

## Шаг 4: Запуск деплоя

После настройки credentials запустите:

```bash
./deploy-full.sh
```

## Альтернатива: Использование существующего профиля

Если у вас уже есть AWS профиль:

```bash
export AWS_PROFILE=your-profile-name
./deploy-full.sh
```

## Необходимые права (IAM Policy)

Ваш IAM пользователь должен иметь права на:
- EC2 (VPC, Subnets, Security Groups)
- ECS (Cluster, Service, Task Definition)
- ECR (Repository)
- ELB (Application Load Balancer)
- S3 (Bucket operations)
- CloudFront (Distribution)
- ACM (Certificate Manager)
- IAM (Role creation)
- CloudWatch (Logs)

Рекомендуемая политика: `AdministratorAccess` (для первого деплоя)

## Безопасность

⚠️ **Важно**:
- Не коммитьте credentials в Git
- Используйте IAM роли где возможно
- Регулярно ротируйте access keys
- Используйте MFA для AWS Console

## Следующие шаги

После настройки AWS CLI:

1. ✅ Запустите `./deploy-full.sh`
2. ✅ Добавьте DNS записи в viazov.dev
3. ✅ Дождитесь валидации SSL сертификатов
4. ✅ Проверьте работу приложения

---

**Готовы продолжить?** Настройте AWS credentials и запустите деплой! 🚀
