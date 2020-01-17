package okex

import "time"

/*
 OKEX api result definition
 @author Lingting Fu
 @date 2018-12-27
 @version 1.0.0
*/

type SwapPositionHolding struct {
	LiquidationPrice float64   `json:"liquidation_price,string"`
	Position         float64   `json:"position,string"`
	AvailPosition    float64   `json:"avail_position,string"`
	AvgCost          float64   `json:"avg_cost,string"`
	SettlementPrice  float64   `json:"settlement_price,string"`
	InstrumentId     string    `json:"instrument_id"`
	Leverage         float64   `json:"leverage,string"`
	RealizedPnl      float64   `json:"realized_pnl,string"`
	Side             string    `json:"side"`
	Timestamp        time.Time `json:"timestamp"`
	Margin           string    `json:"margin";default:""`
}

type SwapPosition struct {
	BizWarmTips
	MarginMode string                `json:"margin_mode"`
	Holding    []SwapPositionHolding `json:"holding"`
}

type SwapPositionList []SwapPosition

type SwapAccountInfo struct {
	Equity            float64   `json:"equity,string"`
	FixedBalance      float64   `json:"fixed_balance,string"`
	InstrumentId      string    `json:"instrument_id"`
	MaintMarginRatio  string    `json:"maint_margin_ratio"`
	Margin            float64   `json:"margin,string"`
	MarginFrozen      float64   `json:"margin_frozen,string"`
	MarginMode        string    `json:"margin_mode"`
	MarginRatio       float64   `json:"margin_ratio,string"`
	MaxWithdraw       float64   `json:"max_withdraw,string"`
	RealizedPnl       float64   `json:"realized_pnl,string"`
	Timestamp         time.Time `json:"timestamp,string"`
	TotalAvailBalance float64   `json:"total_avail_balance,string"`
	UnrealizedPnl     float64   `json:"unrealized_pnl,string"`
}

type SwapAccounts struct {
	BizWarmTips
	Info []SwapAccountInfo `json:"info"`
}

type SwapAccount struct {
	Info SwapAccountInfo `json:"info"`
}

type BaseSwapOrderResult struct {
	OrderId      string `json:"order_id"`
	ClientOid    string `json:"client_oid"`
	ErrorMessage string `json:"error_message"`
	ErrorCode    string `json:"error_code"`
	Result       string `json:"result"`
}

type SwapOrderResult struct {
	BaseSwapOrderResult
	BizWarmTips
}

type SwapOrdersResult struct {
	BizWarmTips
	OrderInfo []BaseSwapOrderResult `json:"order_info"`
}

type SwapCancelOrderResult struct {
	ErrorMessage string `json:"error_message"`
	ErrorCode    string `json:"error_code"`
	OrderId      string `json:"order_id"`
	Result       string `json:"result"`
}

type SwapBatchCancelOrderResult struct {
	BizWarmTips
	InstrumentId string   `json:"instrument_id"`
	Ids          []string `json:"ids"`
	Result       string   `json:"result"`
}

type BaseOrderInfo struct {
	InstrumentId string    `json:"instrument_id"`
	State        string    `json:"state"`
	OrderId      string    `json:"order_id"`
	Timestamp    time.Time `json:"timestamp,string"`
	Price        float64   `json:"price,string"`
	PriceAvg     float64   `json:"price_avg,string"`
	Size         float64   `json:"size,string"`
	Fee          float64   `json:"fee,string"`
	FilledQty    float64   `json:"filled_qty,string"`
	ContractVal  float64   `json:"contract_val,string"`
	Type         string    `json:"type"`
	OrderType    string    `json:"order_type"`
	ClientOid    string    `json:"client_oid"`
}

type SwapOrdersInfo struct {
	BizWarmTips
	OrderInfo []BaseOrderInfo `json:"order_info"`
}

type BaseFillInfo struct {
	InstrumentId string `json:"instrument_id"`
	OrderQty     string `json:"order_qty"`
	TradeId      string `json:"trade_id"`
	Fee          string `json:"fee"`
	OrderId      string `json:"order_id"`
	Timestamp    string `json:"timestamp"`
	Price        string `json:"price"`
	Side         string `json:"side"`
	ExecType     string `json:"exec_type"`
}

type SwapFillsInfo []BaseFillInfo

type SwapAccountsSetting struct {
	BizWarmTips
	InstrumentId  string `json:"instrument_id"`
	LongLeverage  string `json:"long_leverage"`
	ShortLeverage string `json:"short_leverage"`
	MarginMode    string `json:"margin_mode"`
}

type BaseLedgerInfo struct {
	InstrumentId string `json:"instrument_id"`
	Fee          string `json:"fee"`
	Timestamp    string `json:"timestamp"`
	Amount       string `json:"amount"`
	LedgerId     string `json:"ledger_id"`
	Type         string `json:"type"`
}

type SwapAccountsLedgerList []BaseLedgerInfo

type BaseInstrumentInfo struct {
	InstrumentId    string `json:"instrument_id"`
	QuoteCurrency   string `json:"quote_currency"`
	TickSize        string `json:"tick_size"`
	ContractVal     string `json:"contract_val"`
	Listing         string `json:"listing"`
	UnderlyingIndex string `json:"underlying_index"`
	Delivery        string `json:"delivery"`
	Coin            string `json:"coin"`
	SizeIncrement   string `json:"size_increment"`
}

type SwapInstrumentList []BaseInstrumentInfo

type BaesDepthInfo []interface{}
type SwapInstrumentDepth struct {
	BizWarmTips
	Timestamp string          `json:"timestamp"`
	Time      string          `json:"time"`
	Bids      []BaesDepthInfo `json:"bids"`
	Asks      []BaesDepthInfo `json:"asks"`
}

type BaseTickerInfo struct {
	InstrumentId string    `json:"instrument_id"`
	Last         float64   `json:"last,string"`
	Timestamp    time.Time `json:"timestamp"`
	High24h      float64   `json:"high_24h,string"`
	Volume24h    float64   `json:"volume_24h,string"`
	Low24h       float64   `json:"low_24h,string"`
}

type SwapTickerList []BaseTickerInfo

type BaseTradeInfo struct {
	Timestamp string `json:"timestamp"`
	TradeId   string `json:"trade_id"`
	Side      string `json:"side"`
	Price     string `json:"price"`
	Size      string `json:"size"`
}

type SwapTradeList []BaseTradeInfo

type BaseCandleInfo []interface{}
type SwapCandleList []BaseCandleInfo

type SwapIndexInfo struct {
	BizWarmTips
	InstrumentId string `json:"instrument_id"`
	Index        string `json:"index"`
	Timestamp    string `json:"timestamp"`
}

type SwapRate struct {
	InstrumentId string `json:"instrument_id"`
	Timestamp    string `json:"timestamp"`
	Rate         string `json:"rate"`
}

type BaseInstrumentAmount struct {
	BizWarmTips
	InstrumentId string `json:"instrument_id"`
	Timestamp    string `json:"timestamp"`
	Amount       string `json:"amount"`
}

type SwapOpenInterest BaseInstrumentAmount

type SwapPriceLimit struct {
	BizWarmTips
	InstrumentId string `json:"instrument_id"`
	Lowest       string `json:"lowest"`
	Highest      string `json:"highest"`
	Timestamp    string `json:"timestamp"`
}

type BaseLiquidationInfo struct {
	InstrumentId string `json:"instrument_id"`
	Loss         string `json:"loss"`
	CreatedAt    string `json:"created_at"`
	Type         string `json:"type"`
	Price        string `json:"price"`
	Size         string `json:"size"`
}

type SwapLiquidationList []BaseLiquidationInfo

type SwapAccountHolds BaseInstrumentAmount

type SwapFundingTime struct {
	BizWarmTips
	InstrumentId string `json:"instrument_id"`
	FundingTime  string `json:"funding_time"`
}

type SwapMarkPrice struct {
	BizWarmTips
	InstrumentId string `json:"instrument_id"`
	MarkPrice    string `json:"mark_price"`
	Timestamp    string `json:"timestamp"`
}

type BaseHistoricalFundingRate struct {
	InstrumentId string `json:"instrument_id"`
	InterestRate string `json:"interest_rate"`
	FundingRate  string `json:"funding_rate"`
	FundingTime  string `json:"funding_time"`
	RealizedRate string `json:"realized_rate"`
}

type SwapHistoricalFundingRateList []BaseHistoricalFundingRate
