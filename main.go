package genshinartis

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

const MaxSubstats = 4
const DomainBase4Chance = 0.2
const StrongboxBase4Chance = 0.33333333333
const AverageDropsPerDomainRun = 1.065

type ArtifactSubstat struct {
	Stat  artifactStat
	Rolls int
	Value float32
}

func (s *ArtifactSubstat) randomizeValue() {
	s.Value = 0
	for i := 0; i < s.Rolls; i++ {
		s.Value = s.Value + s.Stat.RandomRollValue()
	}
}

func (s *ArtifactSubstat) String() string {
	return fmt.Sprintf("%s: %.1f", s.Stat, s.Value)
}

func (s *ArtifactSubstat) RoundedValue() float32 {
	switch s.Stat {
	case ATK, DEF, HP, ElementalMastery:
		return float32(math.Round(float64(s.Value)))
	default:
		return float32(math.Round(float64(s.Value)*10) / 10)
	}
}

type Artifact struct {
	Set      artifactSet
	Slot     artifactSlot
	MainStat artifactStat
	SubStats [MaxSubstats]*ArtifactSubstat
}

func (a Artifact) subsQuality(subValue map[artifactStat]float32) float32 {
	var quality float32
	for _, sub := range a.SubStats {
		quality += float32(sub.Rolls) * subValue[sub.Stat]
	}
	return quality
}

func (a *Artifact) randomizeSet(options ...artifactSet) {
	a.Set = options[rand.Intn(len(options))]
}

func (a *Artifact) randomizeSlot() {
	a.Slot = artifactSlot(rand.Intn(5))
}

func (a *Artifact) ranzomizeMainStat() {
	switch a.Slot {
	case SlotFlower:
		a.MainStat = HP
	case SlotPlume:
		a.MainStat = ATK
	case SlotSands:
		a.MainStat = weightedRand(sandsWeightedStats)
	case SlotGoblet:
		a.MainStat = weightedRand(gobletWeightedStats)
	case SlotCirclet:
		a.MainStat = weightedRand(circletWeightedStats)
	}
}

func (a *Artifact) randomizeSubstats(base4Chance float32) {
	numRolls := 3 + 5 // starts with 3 subs by default
	if rand.Float32() <= base4Chance {
		numRolls++ // starts with 4 subs
	}

	a.SubStats = [MaxSubstats]*ArtifactSubstat{}
	possibleStats := weightedSubstats(a.MainStat)
	var subs [MaxSubstats]artifactStat
	for i := 0; i < numRolls; i++ {
		// First 4 rolls
		if i < MaxSubstats {
			stat := weightedRand(possibleStats)
			subs[i] = stat
			a.SubStats[i] = &ArtifactSubstat{Stat: stat, Rolls: 1}
			delete(possibleStats, stat)
		} else {
			// Rest of rolls
			index := rand.Intn(MaxSubstats)
			a.SubStats[index].Rolls += 1
		}
	}

	for _, substat := range a.SubStats {
		substat.randomizeValue()
	}
}

func RandomArtifact(base4Chance float32) *Artifact {
	var artifact Artifact
	artifact.randomizeSet(allArtifactSets...)
	artifact.randomizeSlot()
	artifact.ranzomizeMainStat()
	artifact.randomizeSubstats(base4Chance)
	return &artifact
}

func RandomArtifactOfSlot(slot artifactSlot, base4Chance float32) *Artifact {
	var artifact Artifact
	artifact.randomizeSet(allArtifactSets...)
	artifact.Slot = slot
	artifact.ranzomizeMainStat()
	artifact.randomizeSubstats(base4Chance)
	return &artifact
}

func RandomArtifactOfSet(set string, base4Chance float32) *Artifact {
	var artifact Artifact
	artifact.Set = artifactSet(set)
	artifact.randomizeSlot()
	artifact.ranzomizeMainStat()
	artifact.randomizeSubstats(base4Chance)
	return &artifact
}

func RandomArtifactFromDomain(setA, setB string) *Artifact {
	var artifact Artifact
	artifact.randomizeSet(artifactSet(setA), artifactSet(setB))
	artifact.randomizeSlot()
	artifact.ranzomizeMainStat()
	artifact.randomizeSubstats(DomainBase4Chance)
	return &artifact
}

// RemoveTrashArtifacts processes a slice of artifacts and keeps the best ones that have the correct mainstat
// subValue: To know which artifacts are more desirable
// n: Amount of artifacts to keep for every set, slot and main stat (example: n = 10, it will keep at most 10 gladiator atk sands)
func RemoveTrashArtifacts(arts []*Artifact,
	subValue map[artifactStat]float32,
	n int) []*Artifact {
	type SetSlotStat struct {
		set  artifactSet
		slot artifactSlot
		stat artifactStat
	}
	processed := map[SetSlotStat][]*Artifact{}
	for _, art := range arts {
		sss := SetSlotStat{art.Set, art.Slot, art.MainStat}
		processed[sss] = append(processed[sss], art)
	}

	result := []*Artifact{}
	for _, aa := range processed {
		// Ordering the artifacts in processed by sub quality
		sort.Slice(aa, func(i, j int) bool {
			return aa[i].subsQuality(subValue) > aa[j].subsQuality(subValue)
		})
		// Keeping the n best
		if len(aa) > n {
			aa = aa[0:n]
		}
		result = append(result, aa...)
	}
	return result
}
