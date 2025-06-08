import type { RegistryRecommendation } from "../types";

export const mockRegistryRecommendations: RegistryRecommendation[] = [
  { id: "reg-1", imageName: "my-app/feature-x-service", tags: ["pr-123", "sha-a1b2c3d"], sizeMB: 850, lastUsed: "2024-03-15T10:00:00Z", severity: "CRITICAL" },
  { id: "reg-2", imageName: "legacy-data-processor", tags: ["1.2.0-deprecated"], sizeMB: 1200, lastUsed: "2023-11-20T14:30:00Z", severity: "CRITICAL" },
  { id: "reg-3", imageName: "base-images/ubuntu-custom", tags: ["20.04-build-tools"], sizeMB: 2500, lastUsed: "2024-04-01T12:00:00Z", severity: "WARNING" },
  { id: "reg-4", imageName: "marketing/promo-site", tags: ["campaign-2023-q4"], sizeMB: 450, lastUsed: "2024-01-10T18:45:00Z", severity: "CRITICAL" },
  { id: "reg-5", imageName: "dev/experimental-ml-model", tags: ["alpha-0.1"], sizeMB: 3200, lastUsed: "2024-04-22T09:20:00Z", severity: "WARNING" },
  { id: "reg-6", imageName: "internal-tools/ci-runner", tags: ["v2.5.1", "latest-stable"], sizeMB: 600, lastUsed: "2024-02-28T23:00:00Z", severity: "CRITICAL" },
  { id: "reg-7", imageName: "my-app/backend-service", tags: ["hotfix-old-db-conn"], sizeMB: 950, lastUsed: "2024-05-01T11:10:00Z", severity: "WARNING" },
  { id: "reg-8", imageName: "testing/load-test-injector", tags: ["jmeter-5.4.1"], sizeMB: 1100, lastUsed: "2024-03-30T16:00:00Z", severity: "WARNING" },
  { id: "reg-9", imageName: "archived/project-phoenix", tags: ["final-backup"], sizeMB: 5400, lastUsed: "2022-07-15T00:00:00Z", severity: "CRITICAL" },
  { id: "reg-10", imageName: "data-pipelines/etl-worker", tags: ["spark-3.2-hadoop-3.3"], sizeMB: 4200, lastUsed: "2024-02-10T05:00:00Z", severity: "CRITICAL" },
  { id: "reg-11", imageName: "my-app/frontend-ui", tags: ["release-candidate-2.1"], sizeMB: 300, lastUsed: "2024-04-18T21:00:00Z", severity: "WARNING" },
]; 
