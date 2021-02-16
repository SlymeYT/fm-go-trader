package model

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestPosition_Enter(t *testing.T) {
	testUUID := uuid.New()
	testTimestamp := time.Now()

	testCases := []struct {
		name string
		input Position
		arg FillEvent
		expected Position
	}{
		{
			name: "TestPosition_Enter_longDecision",
			input: Position{},
			arg: FillEvent{
				TraceId: testUUID,
				Timestamp: testTimestamp,
				Symbol: "ETH-USD",
				Exchange: "Binance",
				Quantity: 10,
				Decision: DecisionLong,
				FillValueGross: 1000, // abs(Quantity) * £100
				ExchangeFee: 10,
				SlippageFee: 50,
				NetworkFee: 10,
			},
			expected: Position{
				LastUpdateTraceId:   testUUID,
				LastUpdateTimestamp: testTimestamp,
				Symbol:              "ETH-USD",
				Direction:           DecisionLong,
				Quantity:            10,
				EnterFillFees:       map[string]float64{
					"ExchangeFee": 10,
					"SlippageFee": 50,
					"NetworkFee": 10,
					"TotalFees": 70,
				},
				EnterAvgPriceGross:  100,
				EnterFillValueGross:   1000,
				ExitFillFees:        map[string]float64{
					"ExchangeFee": 0,
					"SlippageFee": 0,
					"NetworkFee": 0,
					"TotalFees": 0,
				},
				ExitAvgPriceGross:   0,
				ExitFillValueGross:    0,
				CurrentSymbolPrice:  100,
				CurrentMarketValue:  1000,
				UnrealProfitLoss:    0,
				ResultProfitLoss:    0,
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.input.Enter(testCase.arg)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(testCase.input, testCase.expected); diff != "" {
				t.Fatalf("(-want +got):\n%s", diff)
			}
		})
	}
}