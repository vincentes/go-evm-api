package internal

import "os"

var Provider = ""
var GasEstimateBlockScan = "10"

func LoadEnvironmentVariables() {
	Provider = os.Getenv("MUMBAI_HTTP_PROVIDER")
	GasEstimateBlockScan = os.Getenv("ESTIMATE_BLOCK_SCAN")
}
