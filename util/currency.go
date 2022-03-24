package util

const (
	IDR = "IDR"
	EUR = "EUR"
	USD = "USD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case IDR, EUR, USD:
		return true
	}
	return false
}
