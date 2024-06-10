package mrp

const mrp_price = 3692

func GetMrpPrice(mrp *int) *int {
	if mrp == nil {
		return nil
	}
	price := (*mrp) * mrp_price
	return &price
}
