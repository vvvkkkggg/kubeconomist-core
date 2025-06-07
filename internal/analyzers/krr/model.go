package krr

type krrOutput struct {
	Object      parameter `json:"object"`
	Recommended parameter `json:"recommended"`
}

type parameter struct {
	Requests resources `json:"requests"`
	Limits   resources `json:"limits"`
}

type resources struct {
	CPU    *float64 `json:"cpu"`
	Memory *float64 `json:"memory"`
}
