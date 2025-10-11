#!/bin/bash
echo "Checking Route53 hosted zones..."
aws route53 list-hosted-zones --query "HostedZones[?Name=='viazov.dev.'].{Name:Name,Id:Id}" --output table
