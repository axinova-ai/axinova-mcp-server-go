#!/bin/bash
# Script to open port 9001 for Portainer agents
# Region: cn-hangzhou

echo "Adding port 9001 rules to Aliyun security groups..."
echo ""

# ax-dev-service-01 (ax-dev-app: 120.26.30.40, internal: 172.18.80.46)
echo "1. ax-dev-app (sg-bp14ou72mavs7omlhpgv)"
aliyun ecs AuthorizeSecurityGroup \
  --RegionId cn-hangzhou \
  --SecurityGroupId sg-bp14ou72mavs7omlhpgv \
  --IpProtocol tcp \
  --PortRange 9001/9001 \
  --SourceCidrIp 172.18.80.50/32 \
  --Priority 100 \
  --Description "Portainer agent access from ax-tools"

# ax-dev-db-01 (172.18.80.47)
echo "2. ax-dev-db (sg-bp1fsz1rxpxkori6nmbc)"
aliyun ecs AuthorizeSecurityGroup \
  --RegionId cn-hangzhou \
  --SecurityGroupId sg-bp1fsz1rxpxkori6nmbc \
  --IpProtocol tcp \
  --PortRange 9001/9001 \
  --SourceCidrIp 172.18.80.50/32 \
  --Priority 100 \
  --Description "Portainer agent access from ax-tools"

# ax-prod-service-01 (ax-prod-app: 114.55.132.190, internal: 172.18.80.48)
echo "3. ax-prod-app (sg-bp1j3sza2om9vz27w45a)"
aliyun ecs AuthorizeSecurityGroup \
  --RegionId cn-hangzhou \
  --SecurityGroupId sg-bp1j3sza2om9vz27w45a \
  --IpProtocol tcp \
  --PortRange 9001/9001 \
  --SourceCidrIp 172.18.80.50/32 \
  --Priority 100 \
  --Description "Portainer agent access from ax-tools"

# ax-prod-db-01 (172.18.80.49)
echo "4. ax-prod-db (sg-bp11gctkdgbjoihn01n2)"
aliyun ecs AuthorizeSecurityGroup \
  --RegionId cn-hangzhou \
  --SecurityGroupId sg-bp11gctkdgbjoihn01n2 \
  --IpProtocol tcp \
  --PortRange 9001/9001 \
  --SourceCidrIp 172.18.80.50/32 \
  --Priority 100 \
  --Description "Portainer agent access from ax-tools"

echo ""
echo "✓ Port 9001 rules added to all 4 VPC machines"
echo ""
echo "Note: ax-sas-tools (121.40.188.25) is on a different network."
echo "If you want to add it, you'll need to configure its security group separately."
