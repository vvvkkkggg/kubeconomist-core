package krr

type KrrOutput struct {
	Object      Parameters `json:"object"`
	Recommended Parameters `json:"recommended"`
}

type Parameters struct {
	Cluster   string `json:"cluster"`
	Namespace string `json:"namespace"`
	Pods      []struct {
		Name string `json:"name"`
	} `json:"pods"`
	Container string    `json:"container"`
	Requests  Resources `json:"requests"`
	Limits    Resources `json:"limits"`
}

type Resources struct {
	CPU    *float64 `json:"cpu"`
	Memory *float64 `json:"memory"`
}
