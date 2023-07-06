package gas

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"web3pro.com/blockchain/internal"
	"web3pro.com/blockchain/internal/errors"
	models "web3pro.com/blockchain/models/gas"
)

// Estimate
//
// Calculates the average priority fee per gas for the last n blocks. High, medium and low priority fee per gas are returned.
// /**

func Estimate(c echo.Context) error {
	client, err := ethclient.Dial(internal.Provider)
	if err != nil {
		errors.HandleError(c, err, "Provider connection failed", "Failed to connect to provider.", http.StatusInternalServerError, string(errors.Provider))
		return err
	}

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		errors.HandleError(c, err, "Provider query failed", "Failed to obtain latest block.", http.StatusInternalServerError, string(errors.Provider))
		return err
	}

	// Use GasEstimateBlockScan to determine how many blocks to scan for fee history
	blocks, err := strconv.ParseUint(internal.GasEstimateBlockScan, 10, 64)
	if err != nil {
		errors.HandleError(c, err, "Configuration error", "Failed to parse GasEstimateBlockScan.", http.StatusInternalServerError, string(errors.Configuration))
		return err
	}

	feeHistory, err := client.FeeHistory(context.Background(), blocks, header.Number, []float64{25, 50, 75})
	if err != nil {
		errors.HandleError(c, err, "Provider query failed", "Failed to obtain fee history.", http.StatusInternalServerError, string(errors.Provider))
		return err
	}

	feeData := FeeHistoryToDataPerBlock(feeHistory, blocks)
	slow := averagePriorityFeePerGas(feeData, 0)
	average := averagePriorityFeePerGas(feeData, 1)
	fast := averagePriorityFeePerGas(feeData, 2)

	res := struct {
		MaxPriorityFee models.GasEstimateResponse `json:"maxPriorityFee"`
	}{
		MaxPriorityFee: models.GasEstimateResponse{
			Slow:    slow.String(),
			Average: average.String(),
			Fast:    fast.String(),
		},
	}

	log.Println("Gas estimate response: ", res)
	return c.JSON(http.StatusOK, res)
}

func averagePriorityFeePerGas(feeData []*BlockFeeData, index uint) *big.Int {
	var total = big.NewInt(0)
	for _, blockFeeData := range feeData {
		total = total.Add(total, blockFeeData.priorityFeePerGas[index])
	}
	average := new(big.Int).Div(total, big.NewInt(int64(len(feeData))))
	return average
}

// BlockFeeData
//
// Stores fee data for a block.
// /**
type BlockFeeData struct {
	number            *big.Int
	baseFeePerGas     *big.Int
	gasUsedRatio      float64
	priorityFeePerGas []*big.Int
}

// FeeHistoryToDataPerBlock /**
/**

Converts FeeHistory to an array of BlockFeeData.
This simplifies the process of calculating the average priority fee per gas.
*/
func FeeHistoryToDataPerBlock(history *ethereum.FeeHistory, blocks uint64) []*BlockFeeData {
	blockNum := history.OldestBlock
	untilBlock := big.NewInt(0).Sub(blockNum, new(big.Int).SetUint64(blocks))
	index := 0
	var feeData []*BlockFeeData
	for {
		// If untilBlock >= blockNum, break out of loop
		if untilBlock.Cmp(blockNum) >= 0 {
			break
		}

		blockFeeData := &BlockFeeData{
			number:            blockNum,
			baseFeePerGas:     history.BaseFee[index],
			gasUsedRatio:      history.GasUsedRatio[index],
			priorityFeePerGas: history.Reward[index],
		}
		feeData = append(feeData, blockFeeData)
		blockNum = blockNum.Sub(blockNum, big.NewInt(1))
		index++
	}
	return feeData
}
