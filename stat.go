package genshinartis

import (
	"math/rand"
)

type artifactStat int

const (
	HP artifactStat = iota
	ATK
	DEF
	HPP
	ATKP
	DEFP
	EnergyRecharge
	ElementalMastery
	CritRate
	CritDmg
	PyroDMG
	ElectroDMG
	CryoDMG
	HydroDMG
	AnemoDMG
	GeoDMG
	PhysDMG
	HealingBonus
)

// Weights from https://genshin-impact.fandom.com/wiki/Artifacts/Distribution
// And https://genshin-impact.fandom.com/wiki/Artifacts/Stats
var substatValues map[artifactStat][4]float32 = map[artifactStat][4]float32{
	HP:               {209.13, 239.00, 268.88, 298.75},
	ATK:              {13.62, 15.56, 17.51, 19.45},
	DEF:              {16.20, 18.52, 20.83, 23.15},
	HPP:              {4.08, 4.66, 5.25, 5.83},
	ATKP:             {4.08, 4.66, 5.25, 5.83},
	DEFP:             {5.10, 5.83, 6.56, 7.29},
	ElementalMastery: {16.32, 18.65, 20.98, 23.31},
	EnergyRecharge:   {4.53, 5.18, 5.83, 6.48},
	CritRate:         {2.72, 3.11, 3.50, 3.89},
	CritDmg:          {5.44, 6.22, 6.99, 7.77},
}

func (s artifactStat) String() string {
	switch s {
	case HP:
		return "HP"
	case ATK:
		return "ATK"
	case DEF:
		return "DEF"
	case HPP:
		return "HP%"
	case ATKP:
		return "ATK%"
	case DEFP:
		return "DEF%"
	case EnergyRecharge:
		return "Energy Recharge%"
	case ElementalMastery:
		return "Elemental Mastery"
	case CritRate:
		return "CRIT Rate%"
	case CritDmg:
		return "CRIT DMG%"
	case PyroDMG:
		return "Pyro DMG%"
	case ElectroDMG:
		return "Electro DMG%"
	case CryoDMG:
		return "Cryo DMG%"
	case HydroDMG:
		return "Hydro DMG%"
	case AnemoDMG:
		return "Anemo DMG%"
	case GeoDMG:
		return "Geo DMG%"
	case PhysDMG:
		return "Physical DMG%"
	case HealingBonus:
		return "Healing Bonus%"
	}
	return "Unknown"
}

func (s artifactStat) RandomRollValue() float32 {
	return substatValues[s][rand.Intn(4)]
}

var sandsWeightedStats = map[artifactStat]int{
	HPP:              2668,
	ATKP:             2666,
	DEFP:             2666,
	EnergyRecharge:   1000,
	ElementalMastery: 1000,
}

var gobletWeightedStats = map[artifactStat]int{
	HPP:              2125,
	ATKP:             2125,
	DEFP:             2000,
	PyroDMG:          500,
	ElectroDMG:       500,
	CryoDMG:          500,
	HydroDMG:         500,
	AnemoDMG:         500,
	GeoDMG:           500,
	PhysDMG:          500,
	ElementalMastery: 250,
}

var circletWeightedStats = map[artifactStat]int{
	HPP:              2200,
	ATKP:             2200,
	DEFP:             2200,
	CritRate:         1000,
	CritDmg:          1000,
	HealingBonus:     1000,
	ElementalMastery: 400,
}

const (
	flatSubstatWeight   = 150
	commonSubstatWeight = 100
	critSubstatWeight   = 75
)

func weightedSubstats(mainStat artifactStat) map[artifactStat]int {
	weightedSubs := map[artifactStat]int{
		HP:               flatSubstatWeight,
		ATK:              flatSubstatWeight,
		DEF:              flatSubstatWeight,
		HPP:              commonSubstatWeight,
		ATKP:             commonSubstatWeight,
		DEFP:             commonSubstatWeight,
		EnergyRecharge:   commonSubstatWeight,
		ElementalMastery: commonSubstatWeight,
		CritRate:         critSubstatWeight,
		CritDmg:          critSubstatWeight,
	}
	delete(weightedSubs, mainStat)
	return weightedSubs
}
