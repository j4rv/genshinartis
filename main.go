package genshinartis

import (
	"fmt"
	"math/rand"
	"sort"
)

const MaxSubstats = 4
const DomainBase4Chance = 1.0 / 5.0
const StrongboxBase4Chance = 1.0 / 3.0
const BossBase4Chance = 1.0 / 3.0
const AverageDropsPerDomainRun = 1.065
const DomainExtraArtifactChance = 0.065

type ArtifactSubstat struct {
	Stat  stat
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

type Artifact struct {
	Set           artifactSet
	Slot          artifactSlot
	MainStat      stat
	MainStatValue float32
	SubStats      [MaxSubstats]*ArtifactSubstat
	IsFourLiner   bool
}

func (a Artifact) String() string {
	subsStr := ""
	for _, s := range a.SubStats {
		subsStr += s.String() + "\n"
	}
	return fmt.Sprintf("Set: %s, main stat: %s\n%s", a.Set, a.MainStat, subsStr)
}

func (a Artifact) subsQuality(wantedSubWeights map[stat]float32) float32 {
	var quality float32
	for _, sub := range a.SubStats {
		maxPossibleValue := substatValues[sub.Stat][3]
		quality += wantedSubWeights[sub.Stat] * float32(sub.Value) / maxPossibleValue
	}
	return quality
}

func (a Artifact) cv() float32 {
	var cv float32
	for _, sub := range a.SubStats {
		switch sub.Stat {
		case CritRate:
			cv += sub.Value * 2
		case CritDmg:
			cv += sub.Value
		}
	}
	return cv
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
	a.MainStatValue = mainStatValues[a.MainStat]
}

func (a *Artifact) randomizeSubstats(base4Chance float32) {
	numRolls := 3 + 5 // starts with 3 subs by default
	if rand.Float32() <= base4Chance {
		numRolls++ // starts with 4 subs
		a.IsFourLiner = true
	}

	a.SubStats = [MaxSubstats]*ArtifactSubstat{}
	possibleStats := weightedSubstats(a.MainStat)
	var subs [MaxSubstats]stat
	for i := 0; i < numRolls; i++ {
		// First 4 rolls
		if i < MaxSubstats {
			artiStat := weightedRand(possibleStats)
			subs[i] = artiStat
			a.SubStats[i] = &ArtifactSubstat{Stat: artiStat, Rolls: 1}
			delete(possibleStats, artiStat)
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
	artifact.randomizeSet(AllArtifactSets...)
	artifact.randomizeSlot()
	artifact.ranzomizeMainStat()
	artifact.randomizeSubstats(base4Chance)
	return &artifact
}

func RandomArtifactOfSlot(slot artifactSlot, base4Chance float32) *Artifact {
	var artifact Artifact
	artifact.randomizeSet(AllArtifactSets...)
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
	subValue map[stat]float32,
	n int) []*Artifact {
	type SetSlotStat struct {
		set      artifactSet
		slot     artifactSlot
		mainStat stat
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
