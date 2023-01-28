## Snowflake

An untested version of twiiter's snowflake id generator:

The ID is a 64 bit integer, with 1 unused bit.

- 41 bits are from timestamp adjusted to January 1 2015
- 10 bits are nodeID that is machine dependent (so across multiple machine its unique)
- 12 bits are sequence. 0 to 4095, for uniqueness between millisecond interval

### Requirements:
- go-lang


### Usage

```go
import (
    "github.com/go-batteries/snowflake"
)

snowflake.NextID()

// or

generator := NewSequenceGenerator()
generator.NextID()
```

