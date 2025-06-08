import type { VpcRecommendation } from "../types";

export const mockVpcRecommendations: VpcRecommendation[] = [
  { id: "vpc-1", ipAddress: "10.0.1.23", resourceName: "eks-prod-worker-node-1", region: "us-east-1", unusedForDays: 35, severity: "WARNING" },
  { id: "vpc-2", ipAddress: "10.0.2.112", resourceName: "old-staging-db-replica", region: "us-east-1", unusedForDays: 120, severity: "CRITICAL" },
  { id: "vpc-3", ipAddress: "192.168.4.5", resourceName: "jenkins-build-agent-temp", region: "eu-west-2", unusedForDays: 45, severity: "WARNING" },
  { id: "vpc-4", ipAddress: "10.1.5.88", resourceName: "unattached-nat-gateway-ip", region: "us-west-2", unusedForDays: 210, severity: "CRITICAL" },
  { id: "vpc-5", ipAddress: "172.16.10.30", resourceName: "legacy-vpn-endpoint", region: "ap-southeast-1", unusedForDays: 95, severity: "CRITICAL" },
  { id: "vpc-6", ipAddress: "10.0.3.51", resourceName: "k8s-api-server-ha-ip", region: "us-east-1", unusedForDays: 50, severity: "WARNING" },
  { id: "vpc-7", ipAddress: "192.168.10.101", resourceName: "dev-feature-branch-service", region: "eu-west-2", unusedForDays: 62, severity: "WARNING" },
  { id: "vpc-8", ipAddress: "10.2.2.2", resourceName: "decommissioned-monitoring-server", region: "us-west-2", unusedForDays: 180, severity: "CRITICAL" },
  { id: "vpc-9", ipAddress: "172.17.0.15", resourceName: "docker-bridge-network-ip", region: "ap-southeast-1", unusedForDays: 40, severity: "WARNING" },
  { id: "vpc-10", ipAddress: "10.0.4.200", resourceName: "bastion-host-ip-temp", region: "us-east-1", unusedForDays: 88, severity: "CRITICAL" },
  { id: "vpc-11", ipAddress: "192.168.20.12", resourceName: "gitlab-runner-cache-ip", region: "eu-west-2", unusedForDays: 75, severity: "CRITICAL" },
  { id: "vpc-12", ipAddress: "10.3.3.3", resourceName: "temp-load-balancer-ip", region: "us-west-2", unusedForDays: 55, severity: "WARNING" },
];
