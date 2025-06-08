import type { SubnetRecommendation } from "../types";

export const mockSubnetRecommendations: SubnetRecommendation[] = [
    { id: "sub-1", subnetId: "subnet-0a1b2c3d4e5f6g7h8", cidrBlock: "10.0.1.0/24", region: "us-east-1", totalIps: 251, usedIps: 230, utilization: 91.6, severity: "CRITICAL" },
    { id: "sub-2", subnetId: "subnet-8h7g6f5e4d3c2b1a0", cidrBlock: "10.0.2.0/24", region: "us-east-1", totalIps: 251, usedIps: 201, utilization: 80.1, severity: "WARNING" },
    { id: "sub-3", subnetId: "subnet-9i8j7k6l5m4n3o2p1", cidrBlock: "192.168.0.0/22", region: "eu-west-2", totalIps: 1021, usedIps: 350, utilization: 34.3, severity: "OK" },
    { id: "sub-4", subnetId: "subnet-1p2o3n4m5l6k7j8i9", cidrBlock: "10.1.0.0/20", region: "us-west-2", totalIps: 4091, usedIps: 3800, utilization: 92.9, severity: "CRITICAL" },
    { id: "sub-5", subnetId: "subnet-za9y8x7w6v5u4t3s2", cidrBlock: "172.16.0.0/16", region: "ap-southeast-1", totalIps: 65531, usedIps: 12000, utilization: 18.3, severity: "OK" },
    { id: "sub-6", subnetId: "subnet-1a2b3c4d5e6f7g8h9", cidrBlock: "10.0.3.0/24", region: "us-east-1", totalIps: 251, usedIps: 190, utilization: 75.7, severity: "WARNING" },
    { id: "sub-7", subnetId: "subnet-q1w2e3r4t5y6u7i8o", cidrBlock: "192.168.4.0/22", region: "eu-west-2", totalIps: 1021, usedIps: 850, utilization: 83.2, severity: "WARNING" },
    { id: "sub-8", subnetId: "subnet-p0o9i8u7y6t5r4e3w", cidrBlock: "10.2.0.0/20", region: "us-west-2", totalIps: 4091, usedIps: 1500, utilization: 36.7, severity: "OK" },
    { id: "sub-9", subnetId: "subnet-mnbvcxzlkjhgfdsapoi", cidrBlock: "172.17.0.0/16", region: "ap-southeast-1", totalIps: 65531, usedIps: 60000, utilization: 91.5, severity: "CRITICAL" },
    { id: "sub-10", subnetId: "subnet-qazwsxedcrfvtgbyhnujm", cidrBlock: "10.0.4.0/24", region: "us-east-1", totalIps: 251, usedIps: 50, utilization: 19.9, severity: "OK" },
]; 
