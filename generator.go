package snowflake

import (
	"sync"
	"time"
)

/*
The ID is 64 bit. One bit is unused
Next 41 bits from timestamp adjusted to January 1 2015
Next 10 bits are nodeID that is machine dependent (so across multiple machine its unique)
Next 12 bits are sequence. 0 to 4095. So as to create unique ids within same millisecond
  In case, that's exhausted we wait till the next millisecond
*/

const (
	UNUSED_BITS   = 1
	EPOCH_BITS    = 41
	NODE_ID_BITS  = 10
	SEQUENCE_BITS = 12
)

const (
	maxNodeId   = 1<<NODE_ID_BITS - 1
	maxSequence = 1<<SEQUENCE_BITS - 1
	customEpoch = 1420070400000
)

type SequenceGenerator struct {
	mu            sync.Mutex
	lastTimeStamp int64
	nodeID        int64
	sequence      int64
}

var (
	do        sync.Once
	snowflake *SequenceGenerator
)

func init() {
	do.Do(func() {
		snowflake = NewSequenceGenerator()
	})
}

func NextID() int64 {
	return snowflake.NextID()
}

func NewSequenceGenerator() *SequenceGenerator {
	nodeID := createNodeId()
	return &SequenceGenerator{nodeID: nodeID}
}

func (s *SequenceGenerator) NextID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	var currentTimestamp = timestamp()

	if currentTimestamp < s.lastTimeStamp {
		panic("system clock is destroyed")
	}

	if currentTimestamp == s.lastTimeStamp {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 { //
			currentTimestamp = s.waitNextMillisecond(currentTimestamp)
		}
	} else {
		s.sequence = 0
	}

	s.lastTimeStamp = currentTimestamp

	id := currentTimestamp << (NODE_ID_BITS + SEQUENCE_BITS)
	id |= (s.nodeID << SEQUENCE_BITS)
	id |= s.sequence

	return id
}

func (s *SequenceGenerator) waitNextMillisecond(currentTimestamp int64) int64 {
	for s.lastTimeStamp == currentTimestamp {
		currentTimestamp = timestamp()
	}

	return currentTimestamp
}

func timestamp() int64 {
	return time.Now().UTC().UnixNano()/int64(time.Millisecond) - customEpoch
}
