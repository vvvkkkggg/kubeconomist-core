export interface Allocation {
  requests: {
    cpu: number | null;
    memory: number | null;
  };
  limits: {
    cpu: number | null;
    memory: number | null;
  };
}

export interface Recommended {
  value: number | string | null;
  severity: string;
}

export interface Scan {
  object: {
    cluster: string;
    clusterId: string;
    name: string;
    container: string;
    namespace: string;
    kind: string;
    allocations: Allocation;
  };
  recommended: {
    requests: {
      cpu: Recommended;
      memory: Recommended;
    };
    limits: {
      cpu: Recommended;
      memory: Recommended;
    };
  };
  severity: string;
}

export type Severity = "CRITICAL" | "WARNING" | "OK" | "GOOD" | "UNKNOWN";

export interface KrrReport {
  scans: Scan[];
  score: number;
  description: string;
  date: string;
  strategy: object;
  errors: object[];
  clusterSummary: object;
  config: object;
}

export interface VpcRecommendation {
  id: string;
  cloudId: string;
  folderId: string;
  ipAddress: string;
  isUsed: boolean;
  isReserved: boolean;
}

export type StorageRecommendation = {
  id: string;
  type: 'Object Storage';
  bucketName: string;
  region: string;
  storageClass: 'Standard' | 'Glacier' | 'Deep Archive';
  potentialSavings: number;
  severity: Severity;
} | {
  id: string;
  type: 'Block Storage';
  volumeName: string;
  instanceId: string;
  currentSizeGB: number;
  recommendedSizeGB: number;
  severity: Severity;
};

export interface RegistryRecommendation {
    id: string;
    imageName: string;
    tags: string[];
    sizeMB: number;
    lastUsed: string;
    severity: Severity;
}

export interface DnsRecommendation {
    id: string; // zone_id
    cloudId: string;
    folderId: string;
    zoneId: string;
    isUsed: boolean;
}

export interface NodeOptimizerRecommendation {
  id: string; // instance_id
  cloudId: string;
  folderId: string;
  instanceId: string;
  currentCores: number;
  desiredCores: number;
  currentMemoryGB: number;
  desiredMemoryGB: number;
  currentPrice: number;
  desiredPrice: number;
  savings: number;
}

export interface PlatformOptimizerRecommendation {
  id: string; // node_group_id
  cloudId: string;
  folderId: string;
  nodeGroupId: string;
  currentPlatform: string;
  desiredPlatform: string;
  currentPrice: number;
  desiredPrice: number;
  savings: number;
}
