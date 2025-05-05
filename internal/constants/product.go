package constants

// --- Type Definitions ---
type productTypes struct {
	FixedValueRecharge     string
	RangedValueRecharge    string
	FixedValuePinPurchase  string
	RangedValuePinPurchase string
	RangedValuePayment     string
}

type productServiceIDs struct {
	Mobile     int
	Utilities  int
	GiftCards  int
	Insurance  int
	Finance    int
	Government int
	Education  int
	ESIM       int
}

type productCountryISOCode struct {
	Global string
	Europe string
	India  string
}

// --- Constant Instances ---
var ProductTypes = productTypes{
	FixedValueRecharge:     "FIXED_VALUE_RECHARGE",
	RangedValueRecharge:    "RANGED_VALUE_RECHARGE",
	FixedValuePinPurchase:  "FIXED_VALUE_PIN_PURCHASE",
	RangedValuePinPurchase: "RANGED_VALUE_PIN_PURCHASE",
	RangedValuePayment:     "RANGED_VALUE_PAYMENT",
}

var ProductServiceIDs = productServiceIDs{
	Mobile:     1,
	Utilities:  3,
	GiftCards:  4,
	Insurance:  7,
	Finance:    10,
	Government: 11,
	Education:  12,
	ESIM:       13,
}

var ProductCountryISOCode = productCountryISOCode{
	Global: "GXX",
	Europe: "EXX",
	India:  "IND",
}

var ProductOperatorId = map[string]int{
	"Roblox":    5873,
	"Guatemala": 2060,
}
