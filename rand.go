package genshinartis

import (
	"log"
	"math/rand"
)

func weightedRand(weightedVals map[stat]int) stat {
	sum := 0
	for _, weight := range weightedVals {
		sum += weight
	}

	i := rand.Intn(sum)
	for value, weight := range weightedVals {
		i -= weight
		if i < 0 {
			return value
		}
	}

	log.Println("fatal error in WeightedRand: should never reach this log")
	return 0
}
