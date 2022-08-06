package genshinartis

import "math/rand"

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

func (s artifactStat) goodKey() string {
	switch s {
	case HP:
		return "hp"
	case ATK:
		return "atk"
	case DEF:
		return "def"
	case HPP:
		return "hp_"
	case ATKP:
		return "atk_"
	case DEFP:
		return "def_"
	case EnergyRecharge:
		return "enerRech_"
	case ElementalMastery:
		return "eleMas"
	case CritRate:
		return "critRate_"
	case CritDmg:
		return "critDMG_"
	case PyroDMG:
		return "pyro_dmg_"
	case ElectroDMG:
		return "electro_dmg_"
	case CryoDMG:
		return "cryo_dmg_"
	case HydroDMG:
		return "hydro_dmg_"
	case AnemoDMG:
		return "anemo_dmg_"
	case GeoDMG:
		return "geo_dmg_"
	case PhysDMG:
		return "physical_dmg_"
	case HealingBonus:
		return "heal_"
	}
	return "unknown"
}

// Weights from https://genshin-impact.fandom.com/wiki/Artifacts/Distribution
// And https://genshin-impact.fandom.com/wiki/Artifacts/Stats

func (s artifactStat) RandomRollValue() float32 {
	// For all substats, the following is true:
	// Mid-high roll = highest roll * 0.9
	// Mid roll      = highest roll * 0.8
	// Low roll      = highest roll * 0.7
	var highRoll float32
	switch s {
	case HP:
		highRoll = 298.75
	case ATK:
		highRoll = 19.45
	case DEF:
		highRoll = 23.15
	case HPP:
		highRoll = 5.83
	case ATKP:
		highRoll = 5.83
	case DEFP:
		highRoll = 7.29
	case ElementalMastery:
		highRoll = 23.31
	case EnergyRecharge:
		highRoll = 6.48
	case CritRate:
		highRoll = 3.89
	case CritDmg:
		highRoll = 7.77
	}

	switch rand.Intn(4) {
	case 0:
		return highRoll * 0.7
	case 1:
		return highRoll * 0.8
	case 2:
		return highRoll * 0.9
	default:
		return highRoll
	}
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
