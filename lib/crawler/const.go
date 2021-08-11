package crawler

const (
	MoonSwapBaseURL = "https://moonswap.fi/api/route/opt/swap"
	PathTokenPrice  = MoonSwapBaseURL + "/main/token-price"
	PathAllToken    = MoonSwapBaseURL + "/token/all-token"
)

const (
	CfxScanBaseURL    = "https://confluxscan.io"
	CfxScanApiVersion = "/v1"
	TokenInfoURL      = CfxScanBaseURL + CfxScanApiVersion + "/token/"
	ContractInfoURL   = CfxScanBaseURL + CfxScanApiVersion + "/contract/"
)

