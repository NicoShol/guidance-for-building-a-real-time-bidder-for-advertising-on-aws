package auction

type Campaign struct {
	ID string
}

type Response struct {
	Request  *ExtendedRequest
	Item     *ExtendedItem
	Price    int64
	Campaign *Campaign
}