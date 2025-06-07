package billing

import "testing"

func TestGenCPUNameForGrep(t *testing.T) {
	tests := []struct {
		name         string
		platform     string
		coreFraction string
		want         string
	}{
		{
			name:         "Standard V1 with full core",
			platform:     "standard-v1",
			coreFraction: "100",
			want:         "Intel Broadwell. 100%",
		},
		{
			name:         "Standard V2 with half core",
			platform:     "standard-v2",
			coreFraction: "50",
			want:         "Intel Cascade Lake. 50%",
		},
		{
			name:         "Standard V3 with quarter core",
			platform:     "standard-v3",
			coreFraction: "25",
			want:         "Intel Ice Lake. 25%",
		},
		{
			name:         "High-frequency V3",
			platform:     "highfreq-v3",
			coreFraction: "100",
			want:         "Intel Ice Lake (Compute-Optimized). 100%",
		},
		{
			name:         "AMD platform",
			platform:     "amd-v1",
			coreFraction: "75",
			want:         "AMD Zen 3. 75%",
		},
		{
			name:         "Unknown platform",
			platform:     "unknown-platform",
			coreFraction: "100",
			want:         ". 100%",
		},
		{
			name:         "Empty core fraction",
			platform:     "standard-v1",
			coreFraction: "",
			want:         "Intel Broadwell. %",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := genCPUNameForGrep(tt.platform, tt.coreFraction)
			if got != tt.want {
				t.Errorf("genCPUNameForGrep(%q, %q) = %q, want %q", tt.platform, tt.coreFraction, got, tt.want)
			}
		})
	}
}
