package publicCode

type FingerprintScanResult struct {
	URL           string   `json:"url"`
	Title         string   `json:"title"`
	Fingerprint   []string `json:"fingerprint"`
	Vulnerability []string `json:"vulnerability"`
}

type FingerprintScanRespond struct {
	StatusCode  string            `json:"status_code"`
	Headers     map[string]string `json:"headers"`
	Body        string            `json:"keyword"`
	FaviconHash string            `json:"favicon_hash"`
}

type PocScanResult struct {
	URL     string `json:"url"`
	PocName string `json:"pocName"`
	Print   string `json:"print"`
}
