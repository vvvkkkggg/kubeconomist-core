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
  ipAddress: string;
  resourceName: string;
  region: string;
  unusedForDays: number;
  severity: 'WARNING' | 'CRITICAL';
}

export type StorageRecommendation = {
  id: string;
  type: 'Object Storage';
  bucketName: string;
  region: string;
  storageClass: 'Standard' | 'Glacier' | 'Deep Archive';
  potentialSavings: number;
  severity: 'OK' | 'WARNING' | 'CRITICAL';
} | {
  id: string;
  type: 'Block Storage';
  volumeName: string;
  instanceId: string;
  currentSizeGB: number;
  recommendedSizeGB: number;
  severity: 'OK' | 'WARNING' | 'CRITICAL';
};

export interface RegistryRecommendation {
    id: string;
    imageName: string;
    tags: string[];
    sizeMB: number;
    lastUsed: string;
    severity: 'WARNING' | 'CRITICAL';
}

export interface SubnetRecommendation {
    id: string;
    subnetId: string;
    cidrBlock: string;
    region: string;
    totalIps: number;
    usedIps: number;
    utilization: number;
    severity: 'OK' | 'WARNING' | 'CRITICAL';
} 
