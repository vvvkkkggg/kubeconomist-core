package krr

type Report struct {
	Scans []KrrOutput `json:"scans"`
}

type KrrOutput struct {
	Object      ObjectParameters      `json:"object"`
	Recommended RecommendedParameters `json:"recommended"`
	Severity    string                `json:"severity"`
}

type ObjectParameters struct {
	Cluster   string `json:"cluster"`
	Name      string `json:"name"`
	Container string `json:"container"`
	Pods      []struct {
		Name string `json:"name"`
	} `json:"pods"`
	Namespace   string             `json:"namespace"`
	Kind        string             `json:"kind"`
	Allocations ObjectResourceInfo `json:"allocations"`
	Warnings    []string           `json:"warnings"`
	Labels      map[string]string  `json:"labels"`
}

type ObjectResourceInfo struct {
	Requests Resources         `json:"requests"`
	Limits   Resources         `json:"limits"`
	Info     map[string]string `json:"info"`
}

type Resources struct {
	CPU    *float64 `json:"cpu"`
	Memory *float64 `json:"memory"`
}

type RecommendedParameters struct {
	Requests ExtendedResource  `json:"requests"`
	Limits   ExtendedResource  `json:"limits"`
	Info     map[string]string `json:"info"`
}

type ExtendedResource struct {
	CPU    ValueWithInfo `json:"cpu"`
	Memory ValueWithInfo `json:"memory"`
}

type ValueWithInfo struct {
	Value    *float64 `json:"value"`
	Severity string   `json:"severity"`
}
