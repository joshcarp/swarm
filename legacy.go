// This file is kept to ensure backward compatibility.

package swarm

import (
	"log"
	"math"
	"time"
)

func createRateLimiter(maxRPS int64, requestIncreaseRate string) (rateLimiter RateLimiter, err error) {
	if requestIncreaseRate != "-1" {
		if maxRPS > 0 {
			log.Println("The max RPS that boomer may generate is limited to", maxRPS, "with a increase rate", requestIncreaseRate)
			rateLimiter, err = NewRampUpRateLimiter(maxRPS, requestIncreaseRate, time.Second)
		} else {
			log.Println("The max RPS that boomer may generate is limited by a increase rate", requestIncreaseRate)
			rateLimiter, err = NewRampUpRateLimiter(math.MaxInt64, requestIncreaseRate, time.Second)
		}
	} else {
		if maxRPS > 0 {
			log.Println("The max RPS that boomer may generate is limited to", maxRPS)
			rateLimiter = NewStableRateLimiter(maxRPS, time.Second)
		}
	}
	return rateLimiter, err
}
