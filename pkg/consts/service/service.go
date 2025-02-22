package serviceConsts

const (
	ServiceLoadTypeHTTP = iota
	ServiceLoadTypeTCP
	ServiceLoadTypeGRPC
)

const (
	HTTPRuleTypePrefixURL = iota
	HTTPRuleTypeDomain
)
