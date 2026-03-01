package bidhandler

import (
	"apps_minimalist/bidder/code/auction"
	"apps_minimalist/bidder/code/ksuid"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fastjson"
)

// Used for pooling handler data structures
type pool struct {
	pool 	sync.Pool
}

type persistentData struct {
	parser 			*fastjson.Parser
	ksuidSequence 	*ksuid.Sequence

	auctionRequest 	auction.Request
	byteResponse 	[]byte
}

// Get returns persistentData from pool
// After use, the persistentData created must be Put to the pool
func (pp *pool) Get() *persistentData {
	v := pp.pool.Get()
	if v == nil {
		log.Trace().Msg("allocating bidhandler persistent data")
		return newPersistenData()
	}
}

func (pp *pool) Put(p *persistentData) {
	pp.pool.Put(p)
}

func newPersistenData() *persistentData {
	return &persistentData{
		parser: 	&fastjson.Parser{},
		ksuidSequence: ksuid.NewSequence(),
	}
}