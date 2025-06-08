import type { KrrReport } from '../types';

export const mockReports: KrrReport[] = [
    {
        date: "2024-05-20T11:55:24Z",
        scans: [
            {
                object: {
                    cluster: "dev",
                    name: "my-app-1",
                    container: "my-app-1",
                    namespace: "my-namespace-name",
                    kind: "Deployment",
                    allocations: { requests: { cpu: null, memory: null }, limits: { cpu: null, memory: null } }
                },
                recommended: {
                    requests: { cpu: { value: 0.01, severity: "WARNING" }, memory: { value: 104857600.0, severity: "WARNING" } },
                    limits: { cpu: { value: null, severity: "GOOD" }, memory: { value: 104857600.0, severity: "WARNING" } }
                },
                severity: "WARNING"
            },
            {
                object: {
                    cluster: "dev",
                    name: "my-app-3",
                    container: "my-app-3",
                    namespace: "my-namespace-name",
                    kind: "StatefulSet",
                    allocations: { requests: { cpu: 1.0, memory: 4294967296.0 }, limits: { cpu: 1.0, memory: 4294967296.0 } }
                },
                recommended: {
                    requests: { cpu: { value: 0.085, severity: "CRITICAL" }, memory: { value: 2587885568.0, severity: "CRITICAL" } },
                    limits: { cpu: { value: null, severity: "WARNING" }, memory: { value: 2587885568.0, severity: "CRITICAL" } }
                },
                severity: "CRITICAL"
            }
        ],
        score: 52,
        description: "...",
        strategy: {},
        errors: [],
        clusterSummary: {},
        config: {}
    },
    {
        date: "2024-05-21T14:30:00Z",
        scans: [
            {
                object: {
                    cluster: "dev",
                    name: "my-app-2",
                    container: "my-app-2",
                    namespace: "my-namespace-name",
                    kind: "Deployment",
                    allocations: { requests: { cpu: 0.1, memory: 104857600.0 }, limits: { cpu: 0.5, memory: 536870912.0 } }
                },
                recommended: {
                    requests: { cpu: { value: "?", severity: "UNKNOWN" }, memory: { value: "?", severity: "UNKNOWN" } },
                    limits: { cpu: { value: "?", severity: "UNKNOWN" }, memory: { value: "?", severity: "UNKNOWN" } }
                },
                severity: "UNKNOWN"
            },
            {
                object: {
                    cluster: "dev",
                    name: "my-app-4",
                    container: "somecontainer",
                    namespace: "my-namespace-name",
                    kind: "StatefulSet",
                    allocations: { requests: { cpu: 0.1, memory: 268435456.0 }, limits: { cpu: null, memory: null } }
                },
                recommended: {
                    requests: { cpu: { value: 0.443, severity: "WARNING" }, memory: { value: 361758720.0, severity: "GOOD" } },
                    limits: { cpu: { value: null, severity: "GOOD" }, memory: { value: 361758720.0, severity: "WARNING" } }
                },
                severity: "WARNING"
            },
            {
                object: {
                    cluster: "dev",
                    name: "another-service",
                    container: "main",
                    namespace: "prod",
                    kind: "Deployment",
                    allocations: { requests: { cpu: 2, memory: 2147483648 }, limits: { cpu: 4, memory: 4294967296 } },
                },
                recommended: {
                    requests: { cpu: { value: 0.5, severity: "CRITICAL" }, memory: { value: 1073741824, severity: "CRITICAL" } },
                    limits: { cpu: { value: 1, severity: "WARNING" }, memory: { value: 2147483648, severity: "WARNING" } }
                },
                severity: "CRITICAL"
            }
        ],
        score: 78,
        description: "...",
        strategy: {},
        errors: [],
        clusterSummary: {},
        config: {}
    }
]; 
