package auction

import (
	"math/rand"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type pool struct {
	pool sync.Pool
}

type persistentData struct {
	rng *rand.Rand
}

// get returns persistentData from pool.
// The persistentData must be put to pool after use.
func (pp *pool) get() *persistentData {
	v := pp.pool.Get()
	if v == nil {
		log.Trace().Msg("allocating auction persistent data")
		return newPersistentData()
	}

	return v.(*persistentData)
}

// put returns p to pool.
func (pp *pool) put(p *persistentData) {
	pp.pool.Put(p)
}

func newPersistentData() *persistentData {
	return &persistentData{
		rng: rand.New(rand.NewSource(time.Now().Unix())),
	}
}