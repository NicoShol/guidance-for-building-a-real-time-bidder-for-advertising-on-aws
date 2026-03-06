package auction

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// defaultCampaign is a placeholder campaign used for testing.
// TODO: Replace with actual campaign selection logic.
var defaultCampaign = Campaign{ID: "constant-campaign"}

type Auction struct {
	pool *pool
}

// TDOD : understand if any "cache" needed & how to implement
func New() *Auction {
	return &Auction{
		pool: &pool{},
	}
}

func (a *Auction) Run(deadline time.Time, request *ExtendedRequest, response *Response) error {
	persistentData := a.pool.get()
	err := a.run(deadline, request, response, persistentData)
	a.pool.put(persistentData)

	if err != nil {
		return err
	}
	
	return nil
}

func (a *Auction) run(
	deadline time.Time,
	request *ExtendedRequest,
	response *Response,
	pd *persistentData,
) error {
	// Early deadline check
	if time.Now().After(deadline) {
		return ErrTimeout
	}

	if zerolog.GlobalLevel() == zerolog.TraceLevel {
		log.Trace().Msgf("received bidrequest with ID %s item %s, rng %d", 
			request.ID, request.Items[0].ID, pd.rng.Int())
	}
	
	// Future: When I add real auction logic with database lookups:
	// 
	// // Get device data
	// if err := a.deviceRepo.Get(deadline, deviceID, &device); err != nil {
	//     return err  // Already checks deadline internally
	// }
	//
	// // Check deadline before expensive campaign selection
	// if time.Now().After(deadline) {
	//     return ErrTimeout
	// }
	//
	// // Select campaigns
	// campaigns := a.selectCampaigns(device, pd)
	//
	// // Check deadline before bid calculation  
	// if time.Now().After(deadline) {
	//     return ErrTimeout
	// }
	//
	// // Calculate bid
	// bid := a.calculateBid(campaigns, pd)
	
	*response = Response{
		Request:  request,
		Item:     &request.Items[0],
		Price:    1,
		Campaign: &defaultCampaign,
	}
	
	return nil
}