package store

var (
	RS *ReceiptStore
	PS *PointsStore
)

func InitializeStores() {
	RS = NewReceiptStore()
	PS = NewPointsStore()
}
