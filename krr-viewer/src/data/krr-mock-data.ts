import type { KrrReport } from '../types';

export const krrReport: KrrReport = {
  "scans": [
    {
      "object": {
        "cluster": "dev",
        "clusterId": "catjf4lqjrdc4u4cbepj",
        "name": "my-app-1",
        "container": "my-app-1",
        "namespace": "my-namespace-name",
        "kind": "Deployment",
        "allocations": {
          "requests": {
            "cpu": null,
            "memory": null
          },
          "limits": {
            "cpu": null,
            "memory": null
          }
        }
      },
      "recommended": {
        "requests": {
          "cpu": {
            "value": 0.01,
            "severity": "WARNING"
          },
          "memory": {
            "value": 104857600.0,
            "severity": "WARNING"
          }
        },
        "limits": {
          "cpu": {
            "value": null,
            "severity": "GOOD"
          },
          "memory": {
            "value": 104857600.0,
            "severity": "WARNING"
          }
        }
      },
      "severity": "WARNING"
    },
    {
      "object": {
        "cluster": "dev",
        "clusterId": "catjf4lqjrdc4u4cbepj",
        "name": "my-app-2",
        "container": "my-app-2",
        "namespace": "my-namespace-name",
        "kind": "Deployment",
        "allocations": {
          "requests": {
            "cpu": 0.1,
            "memory": 104857600.0
          },
          "limits": {
            "cpu": 0.5,
            "memory": 536870912.0
          }
        }
      },
      "recommended": {
        "requests": {
          "cpu": {
            "value": "?",
            "severity": "UNKNOWN"
          },
          "memory": {
            "value": "?",
            "severity": "UNKNOWN"
          }
        },
        "limits": {
          "cpu": {
            "value": "?",
            "severity": "UNKNOWN"
          },
          "memory": {
            "value": "?",
            "severity": "UNKNOWN"
          }
        }
      },
      "severity": "UNKNOWN"
    },
    {
      "object": {
        "cluster": "dev",
        "clusterId": "catjf4lqjrdc4u4cbepj",
        "name": "my-app-3",
        "container": "my-app-3",
        "namespace": "my-namespace-name",
        "kind": "StatefulSet",
        "allocations": {
          "requests": {
            "cpu": 1.0,
            "memory": 4294967296.0
          },
          "limits": {
            "cpu": 1.0,
            "memory": 4294967296.0
          }
        }
      },
      "recommended": {
        "requests": {
          "cpu": {
            "value": 0.085,
            "severity": "CRITICAL"
          },
          "memory": {
            "value": 2587885568.0,
            "severity": "CRITICAL"
          }
        },
        "limits": {
          "cpu": {
            "value": null,
            "severity": "WARNING"
          },
          "memory": {
            "value": 2587885568.0,
            "severity": "CRITICAL"
          }
        }
      },
      "severity": "CRITICAL"
    },
    {
      "object": {
        "cluster": "dev",
        "clusterId": "catjf4lqjrdc4u4cbepj",
        "name": "my-app-4",
        "container": "somecontainer",
        "namespace": "my-namespace-name",
        "kind": "StatefulSet",
        "allocations": {
          "requests": {
            "cpu": 0.1,
            "memory": 268435456.0
          },
          "limits": {
            "cpu": null,
            "memory": null
          }
        }
      },
      "recommended": {
        "requests": {
          "cpu": {
            "value": 0.443,
            "severity": "WARNING"
          },
          "memory": {
            "value": 361758720.0,
            "severity": "GOOD"
          }
        },
        "limits": {
          "cpu": {
            "value": null,
            "severity": "GOOD"
          },
          "memory": {
            "value": 361758720.0,
            "severity": "WARNING"
          }
        }
      },
      "severity": "WARNING"
    },
    {
      "object": {
        "cluster": "dev",
        "clusterId": "catjf4lqjrdc4u4cbepj",
        "name": "my-app-4",
        "container": "exporter",
        "namespace": "my-namespace-name",
        "kind": "StatefulSet",
        "allocations": {
          "requests": {
            "cpu": 0.02,
            "memory": 33554432.0
          },
          "limits": {
            "cpu": null,
            "memory": 268435456.0
          }
        }
      },
      "recommended": {
        "requests": {
          "cpu": {
            "value": 0.03,
            "severity": "GOOD"
          },
          "memory": {
            "value": 104857600.0,
            "severity": "GOOD"
          }
        },
        "limits": {
          "cpu": {
            "value": null,
            "severity": "GOOD"
          },
          "memory": {
            "value": 104857600.0,
            "severity": "OK"
          }
        }
      },
      "severity": "OK"
    },
    {
      "object": {
        "cluster": "dev",
        "clusterId": "catjf4lqjrdc4u4cbepj",
        "name": "alertmanager-kube-prometheus-stack-alertmanager-0",
        "container": "alertmanager",
        "namespace": "prometheus-operator-space",
        "kind": "StatefulSet",
        "allocations": { "requests": { "cpu": 0.1, "memory": 268435456 }, "limits": { "cpu": 0.1, "memory": 268435456 } }
      },
      "recommended": { "requests": { "cpu": { "value": 0.05, "severity": "CRITICAL" }, "memory": { "value": 134217728, "severity": "CRITICAL" } }, "limits": { "cpu": { "value": null, "severity": "GOOD" }, "memory": { "value": 134217728, "severity": "CRITICAL" } } },
      "severity": "CRITICAL"
    },
    {
      "object": {
        "cluster": "dev",
        "clusterId": "catjf4lqjrdc4u4cbepj",
        "name": "kube-prometheus-stack-grafana-6846988bcd-9jxqm",
        "container": "grafana",
        "namespace": "prometheus-operator-space",
        "kind": "Deployment",
        "allocations": { "requests": { "cpu": 0.2, "memory": 536870912 }, "limits": { "cpu": 0.5, "memory": 1073741824 } }
      },
      "recommended": { "requests": { "cpu": { "value": 0.1, "severity": "CRITICAL" }, "memory": { "value": 268435456, "severity": "CRITICAL" } }, "limits": { "cpu": { "value": null, "severity": "GOOD" }, "memory": { "value": 268435456, "severity": "CRITICAL" } } },
      "severity": "CRITICAL"
    },
    {
      "object": {
        "cluster": "dev",
        "clusterId": "catjf4lqjrdc4u4cbepj",
        "name": "kube-prometheus-stack-kube-state-metrics-7dc7f8f984-2xs56",
        "container": "kube-state-metrics",
        "namespace": "prometheus-operator-space",
        "kind": "Deployment",
        "allocations": { "requests": { "cpu": 0.1, "memory": 268435456 }, "limits": { "cpu": null, "memory": null } }
      },
      "recommended": { "requests": { "cpu": { "value": 0.05, "severity": "CRITICAL" }, "memory": { "value": 134217728, "severity": "CRITICAL" } }, "limits": { "cpu": { "value": null, "severity": "GOOD" }, "memory": { "value": null, "severity": "GOOD" } } },
      "severity": "CRITICAL"
    },
    {
      "object": {
        "cluster": "dev",
        "clusterId": "catjf4lqjrdc4u4cbepj",
        "name": "kube-prometheus-stack-operator-7dcbd76b88-cfngc",
        "container": "kube-prometheus-stack",
        "namespace": "prometheus-operator-space",
        "kind": "Deployment",
        "allocations": { "requests": { "cpu": 0.1, "memory": 268435456 }, "limits": { "cpu": null, "memory": null } }
      },
      "recommended": { "requests": { "cpu": { "value": 0.15, "severity": "WARNING" }, "memory": { "value": 200000000, "severity": "WARNING" } }, "limits": { "cpu": { "value": null, "severity": "GOOD" }, "memory": { "value": null, "severity": "GOOD" } } },
      "severity": "WARNING"
    }
  ],
  "score": 52,
  "description": "Simple Strategy",
  "date": "2024-07-29T12:00:00Z",
  "strategy": {
    "name": "simple"
  },
  "errors": [],
  "clusterSummary": {},
  "config": {}
}; 
