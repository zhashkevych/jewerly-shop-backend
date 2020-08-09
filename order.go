package jewerly

const (
	OrderStatusCreated = iota + 1
	OrderStatusPaymentSuccess
	OrderStatusPaymentFailed
	OrderStatusProcessed
)

type Order struct {

}
