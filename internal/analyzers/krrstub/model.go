package krrstub

type KrrOutput struct {
	Object      Parameters `json:"object"`
	Recommended Parameters `json:"recommended"`
}

type Parameters struct {
	Cluster   string `json:"cluster"`
	Namespace string `json:"namespace"`
	Pods      []Pod
	Container string    `json:"container"`
	Requests  Resources `json:"requests"`
	Limits    Resources `json:"limits"`
}

type Resources struct {
	CPU    *float64 `json:"cpu"`
	Memory *float64 `json:"memory"`
}

type Pod struct {
	Name string `json:"name"`
}
