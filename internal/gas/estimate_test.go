package gas

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"web3pro.com/blockchain/internal"
	models "web3pro.com/blockchain/models/gas"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestEstimate(t *testing.T) {
	// Create a simulated backend
	backend := backends.NewSimulatedBackend(core.GenesisAlloc{
		common.HexToAddress("0x0123456789abcdef0123456789abcdef01234567"): {Balance: big.NewInt(1000000000000000000)},
	}, 8000000)

	// Set up the gas estimate block scan value
	internal.GasEstimateBlockScan = "3"

	// Mock the fee history for testing
	mockFeeHistory := &ethereum.FeeHistory{
		OldestBlock:  big.NewInt(100),
		BaseFee:      []*big.Int{big.NewInt(100), big.NewInt(200), big.NewInt(300)},
		GasUsedRatio: []float64{0.25, 0.5, 0.75},
		Reward:       [][]*big.Int{{big.NewInt(500), big.NewInt(600), big.NewInt(700)}, {big.NewInt(500), big.NewInt(600), big.NewInt(700)}, {big.NewInt(500), big.NewInt(600), big.NewInt(700)}},
	}

	// Replace the actual Ethereum client with the simulated backend
	setClient(backend)

	// Perform a simulated Ethereum operation (fund an account, for example)
	account := common.HexToAddress("0x0123456789abcdef0123456789abcdef01234567")
	_ = backend.FundAddress(account, new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)))

	// Call the Estimate handler
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	err := Estimate(c)

	// Assertions

	// Check if there was no error
	assert.NoError(t, err)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rec.Code)

	// Check the response body
	expectedResponse := struct {
		MaxPriorityFee models.GasEstimateResponse `json:"maxPriorityFee"`
	}{
		MaxPriorityFee: models.GasEstimateResponse{
			Slow:    "1500000000",
			Average: "11219762040",
			Fast:    "19624999984",
		},
	}
	assert.JSONEq(t, `{"maxPriorityFee":{"slow":"1500000000","average":"11219762040","fast":"19624999984"}}`, rec.Body.String())
}
