package genshinartis

import (
	"math"
)

/**
"Simple" Optimizer, very WORK IN PROGRESS
Right now is just made to optimize Xiao plunges
**/

type element int
type attackTag int

const (
	Physical = iota
	Pyro
	Hydro
	Anemo
	Electro
	Dendro
	Cryo
	Geo
)

type attack struct {
	tag           attackTag
	element       element
	offensiveStat stat
	multiplier    float32
}

type weapon struct {
	baseAtk float32
	stats   map[stat]float32
	passive func(map[stat]float32) map[stat]float32
}

type character struct {
	level      int
	baseHP     float32
	baseAtk    float32
	baseDef    float32
	bonusStats map[stat]float32
	artifacts  map[artifactSlot]*Artifact
	weapon     weapon
}

func (c character) artifactStats() map[stat]float32 {
	s := map[stat]float32{}
	for _, art := range c.artifacts {
		s[art.MainStat] = s[art.MainStat] + art.MainStatValue
		for _, subStat := range art.SubStats {
			s[subStat.Stat] = s[subStat.Stat] + subStat.Value
		}
	}
	return s
}

func (c character) stats() map[stat]float32 {
	stats := map[stat]float32{}
	wepStats := c.weapon.stats
	artStats := c.artifactStats()

	// merge all the stats
	for stat, v := range c.bonusStats {
		stats[stat] = stats[stat] + v
	}
	for stat, v := range wepStats {
		stats[stat] = stats[stat] + v
	}
	for stat, v := range artStats {
		stats[stat] = stats[stat] + v
	}
	// set bonuses
	for stat, v := range artifactSetBonus(c.artifacts) {
		stats[stat] = stats[stat] + v
	}

	// now calculate stats with base values
	stats[HP] = c.baseHP*(1+stats[HPP]/100) + stats[HP]
	stats[ATK] = (c.baseAtk+c.weapon.baseAtk)*(1+stats[ATKP]/100) + stats[ATK]
	stats[DEF] = c.baseDef*(1+stats[DEFP]/100) + stats[DEF]
	stats[CritRate] = 5 + stats[CritRate]
	stats[CritDmg] = 50 + stats[CritDmg]
	stats[EnergyRecharge] = 100 + stats[EnergyRecharge]

	// apply weapon passive
	if c.weapon.passive != nil {
		for stat, v := range c.weapon.passive(stats) {
			stats[stat] = stats[stat] + v
		}
	}

	return stats
}

type optimizationConfig struct {
	character character
	target    attack
	artifacts []*Artifact
}

func (c optimizationConfig) findBest(artifactFilter func([]*Artifact) []*Artifact, buildFilter func(map[artifactSlot]*Artifact) bool) (map[artifactSlot]*Artifact, float32) {
	artifacts := c.artifacts
	if artifactFilter != nil {
		artifacts = artifactFilter(artifacts)
	}

	flowers := []*Artifact{}
	plumes := []*Artifact{}
	sandss := []*Artifact{}
	goblets := []*Artifact{}
	circlets := []*Artifact{}

	for _, art := range artifacts {
		switch art.Slot {
		case SlotFlower:
			flowers = append(flowers, art)
		case SlotPlume:
			plumes = append(plumes, art)
		case SlotSands:
			sandss = append(sandss, art)
		case SlotGoblet:
			goblets = append(goblets, art)
		case SlotCirclet:
			circlets = append(circlets, art)
		}
	}

	var best map[artifactSlot]*Artifact
	var bestTargetValue float32

	for _, flower := range flowers {
		for _, plume := range plumes {
			for _, sands := range sandss {
				for _, goblet := range goblets {
					for _, circlet := range circlets {
						build := map[artifactSlot]*Artifact{
							SlotFlower:  flower,
							SlotPlume:   plume,
							SlotSands:   sands,
							SlotGoblet:  goblet,
							SlotCirclet: circlet,
						}

						if !buildFilter(build) {
							continue
						}

						c.character.artifacts = build
						value := c.calculateTargetValue()
						if value > bestTargetValue {
							best = build
							bestTargetValue = value
						}
					}
				}
			}
		}
	}

	return best, bestTargetValue
}

func (c optimizationConfig) calculateTargetValue() float32 {
	t := c.target
	stats := c.character.stats()

	resMult := float32(1.1) // TEMP
	defMult := float32(0.5) // TEMP
	mvStatValue := stats[c.target.offensiveStat]
	critMult := critMultiplier(stats[CritRate], stats[CritDmg])
	dmgMult := 1 + stats[AnemoDMG]/100 + stats[GlobalDMGBonus]/100 // TEMP
	return (t.multiplier/100*mvStatValue + stats[BaseDMGIncrease]) * critMult * dmgMult * resMult * defMult
}

func critMultiplier(critRate, critDmg float32) float32 {
	finalCR := float32(math.Min(100, float64(critRate))) / 100
	return finalCR * (1 + critDmg/100)
}
