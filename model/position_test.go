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
				EnterAvgPriceNet:    107,
				EnterFillValueNet:   1070,
				ExitFillFees:        map[string]float64{
					"ExchangeFee": 0,
					"SlippageFee": 0,
					"NetworkFee": 0,
					"TotalFees": 0,
				},
				ExitAvgPriceGross:   0,
				ExitAvgPriceNet:     0,
				ExitFillValueNet:    0,
				CurrentSymbolPrice:  100,
				CurrentMarketValue:  1070,
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

func TestPosition_Update(t *testing.T) {
	type fields struct {
		LastUpdateTraceId   uuid.UUID
		LastUpdateTimestamp time.Time
		Symbol              string
		Direction           string
		Quantity            float64
		EnterFillFees       map[string]float64
		EnterAvgPriceGross  float64
		EnterAvgPriceNet    float64
		EnterFillValueNet   float64
		ExitFillFees        map[string]float64
		ExitAvgPriceGross   float64
		ExitAvgPriceNet     float64
		ExitFillValueNet    float64
		CurrentSymbolPrice  float64
		CurrentMarketValue  float64
		UnrealProfitLoss    float64
		ResultProfitLoss    float64
	}
	type args struct {
		market MarketEvent
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Position{
				LastUpdateTraceId:   tt.fields.LastUpdateTraceId,
				LastUpdateTimestamp: tt.fields.LastUpdateTimestamp,
				Symbol:              tt.fields.Symbol,
				Direction:           tt.fields.Direction,
				Quantity:            tt.fields.Quantity,
				EnterFillFees:       tt.fields.EnterFillFees,
				EnterAvgPriceGross:  tt.fields.EnterAvgPriceGross,
				EnterAvgPriceNet:    tt.fields.EnterAvgPriceNet,
				EnterFillValueNet:   tt.fields.EnterFillValueNet,
				ExitFillFees:        tt.fields.ExitFillFees,
				ExitAvgPriceGross:   tt.fields.ExitAvgPriceGross,
				ExitAvgPriceNet:     tt.fields.ExitAvgPriceNet,
				ExitFillValueNet:    tt.fields.ExitFillValueNet,
				CurrentSymbolPrice:  tt.fields.CurrentSymbolPrice,
				CurrentMarketValue:  tt.fields.CurrentMarketValue,
				UnrealProfitLoss:    tt.fields.UnrealProfitLoss,
				ResultProfitLoss:    tt.fields.ResultProfitLoss,
			}
			if err := p.Update(tt.args.market); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPosition_Exit(t *testing.T) {
	type fields struct {
		LastUpdateTraceId   uuid.UUID
		LastUpdateTimestamp time.Time
		Symbol              string
		Direction           string
		Quantity            float64
		EnterFillFees       map[string]float64
		EnterAvgPriceGross  float64
		EnterAvgPriceNet    float64
		EnterFillValueNet   float64
		ExitFillFees        map[string]float64
		ExitAvgPriceGross   float64
		ExitAvgPriceNet     float64
		ExitFillValueNet    float64
		CurrentSymbolPrice  float64
		CurrentMarketValue  float64
		UnrealProfitLoss    float64
		ResultProfitLoss    float64
	}
	type args struct {
		fill FillEvent
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Position{
				LastUpdateTraceId:   tt.fields.LastUpdateTraceId,
				LastUpdateTimestamp: tt.fields.LastUpdateTimestamp,
				Symbol:              tt.fields.Symbol,
				Direction:           tt.fields.Direction,
				Quantity:            tt.fields.Quantity,
				EnterFillFees:       tt.fields.EnterFillFees,
				EnterAvgPriceGross:  tt.fields.EnterAvgPriceGross,
				EnterAvgPriceNet:    tt.fields.EnterAvgPriceNet,
				EnterFillValueNet:   tt.fields.EnterFillValueNet,
				ExitFillFees:        tt.fields.ExitFillFees,
				ExitAvgPriceGross:   tt.fields.ExitAvgPriceGross,
				ExitAvgPriceNet:     tt.fields.ExitAvgPriceNet,
				ExitFillValueNet:    tt.fields.ExitFillValueNet,
				CurrentSymbolPrice:  tt.fields.CurrentSymbolPrice,
				CurrentMarketValue:  tt.fields.CurrentMarketValue,
				UnrealProfitLoss:    tt.fields.UnrealProfitLoss,
				ResultProfitLoss:    tt.fields.ResultProfitLoss,
			}
			if err := p.Exit(tt.args.fill); (err != nil) != tt.wantErr {
				t.Errorf("Exit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_calculateProfitLoss(t *testing.T) {
	type args struct {
		direction      string
		quantity       float64
		exitFillValue  float64
		enterFillValue float64
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateProfitLoss(tt.args.direction, tt.args.quantity, tt.args.exitFillValue, tt.args.enterFillValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateProfitLoss() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calculateProfitLoss() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateProfitLossV2(t *testing.T) {
	type args struct {
		position Position
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateProfitLossV2(tt.args.position)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateProfitLossV2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calculateProfitLossV2() got = %v, want %v", got, tt.want)
			}
		})
	}
}