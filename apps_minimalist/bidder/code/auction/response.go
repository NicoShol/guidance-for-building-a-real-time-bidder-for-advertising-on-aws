package auction

type Response struct {
	Request  *Request
	Item     *Item
	Price    int64
}