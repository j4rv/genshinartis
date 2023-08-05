package genshinartis

func artifactSetBonus(artifactBuild map[artifactSlot]*Artifact) map[stat]float32 {
	bonus := map[stat]float32{}
	setCount := map[artifactSet]int{}
	for _, artifact := range artifactBuild {
		setCount[artifact.Set] = setCount[artifact.Set] + 1
	}
	for set, count := range setCount {
		if count >= 2 {
			set2pBonus := twoPieceBonus(set)
			for stat, value := range set2pBonus {
				bonus[stat] = bonus[stat] + value
			}
		}
		if count >= 4 {
			set4pBonus := fourPieceBonus(set)
			for stat, value := range set4pBonus {
				bonus[stat] = bonus[stat] + value
			}
		}
	}
	return bonus
}

// Very WIP, no stacks config, not all sets, etc

func twoPieceBonus(set artifactSet) map[stat]float32 {
	bonus := map[stat]float32{}
	switch set {
	case "VermillionHereafter":
		bonus[ATKP] = 18
	}
	return bonus
}

func fourPieceBonus(set artifactSet) map[stat]float32 {
	bonus := map[stat]float32{}
	switch set {
	case "VermillionHereafter":
		bonus[ATKP] = 8 + 10*4
	case "MarechausseeHunter":
		bonus[CritRate] = 12 * 3
	}
	return bonus
}
