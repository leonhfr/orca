package uci

// compile time check that optionInteger implements option.
var _ option = optionInteger{}

// compile time check that optionBoolean implements option.
var _ option = optionBoolean{}
