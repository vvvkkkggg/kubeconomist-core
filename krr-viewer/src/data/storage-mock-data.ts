import type { StorageRecommendation } from "../types";

export const mockStorageRecommendations: StorageRecommendation[] = [
  { id: "store-1", type: "Block Storage", volumeName: "prod-db-data-vol", instanceId: "i-0123456789abcdef0", currentSizeGB: 1024, recommendedSizeGB: 500, severity: "CRITICAL" },
  { id: "store-2", type: "Object Storage", bucketName: "company-assets-archive", region: "us-east-1", storageClass: "Standard", potentialSavings: 1250.50, severity: "CRITICAL" },
  { id: "store-3", type: "Block Storage", volumeName: "staging-app-logs", instanceId: "i-fdecba9876543210", currentSizeGB: 250, recommendedSizeGB: 100, severity: "WARNING" },
  { id: "store-4", type: "Object Storage", bucketName: "dev-team-test-data", region: "eu-west-2", storageClass: "Standard", potentialSavings: 350.00, severity: "WARNING" },
  { id: "store-5", type: "Block Storage", volumeName: "jenkins-workspace-vol", instanceId: "i-abcdef1234567890", currentSizeGB: 500, recommendedSizeGB: 200, severity: "WARNING" },
  { id: "store-6", type: "Object Storage", bucketName: "analytics-event-stream-cold", region: "us-west-2", storageClass: "Glacier", potentialSavings: 800.75, severity: "OK" },
  { id: "store-7", type: "Block Storage", volumeName: "shared-nfs-volume", instanceId: "i-0987654321fedcba", currentSizeGB: 2048, recommendedSizeGB: 1024, severity: "CRITICAL" },
  { id: "store-8", type: "Object Storage", bucketName: "marketing-campaign-assets", region: "ap-southeast-1", storageClass: "Standard", potentialSavings: 50.25, severity: "OK" },
  { id: "store-9", type: "Block Storage", volumeName: "monitoring-tsdb-data", instanceId: "i-fedcba0987654321", currentSizeGB: 750, recommendedSizeGB: 400, severity: "WARNING" },
  { id: "store-10", type: "Object Storage", bucketName: "legal-document-archive-2022", region: "eu-central-1", storageClass: "Standard", potentialSavings: 2200.00, severity: "CRITICAL" },
  { id: "store-11", type: "Block Storage", volumeName: "dev-user-home-dir-vol", instanceId: "i-1234567890abcdef", currentSizeGB: 100, recommendedSizeGB: 50, severity: "WARNING" },
]; 
