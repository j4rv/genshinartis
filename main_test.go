package genshinartis

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"testing"
	"time"
)

func TestOptimizer(t *testing.T) {
	piecesToGenerate := 862
	repetitions := 50
	set := "MarechausseeHunter" // MarechausseeHunter / VermillionHereafter

	c := character{
		level:   90,
		baseAtk: 349,
		weapon:  weaponPJWSFullStacks,
		bonusStats: map[stat]float32{
			ATK:             1203,             // Benny
			ATKP:            20 + 15 + 25,     // Tenacity, Noblesse, Pyro resonance, TTDS, etc
			BaseDMGIncrease: 208.27,           // Faru A4
			AnemoDMG:        22.5 + 95.2 + 15, // Faruzan, Xiao burst and Xiao A1
			CritRate:        24.2 - 5,         // Xiao main stat
			CritDmg:         40,               // Faruzan c6
		},
	}

	var bestTargetValueSum float32
	for i := 0; i < repetitions; i++ {
		var artis []*Artifact
		// Random pieces of an offset
		for i := 0; i < piecesToGenerate*2; i++ {
			arti := RandomArtifactOfSet("GladiatorsFinale", DomainBase4Chance)
			artis = append(artis, arti)
		}
		// Random pieces of a specific set
		for i := 0; i < piecesToGenerate; i++ {
			artis = append(artis, RandomArtifactOfSet(set, DomainBase4Chance))
			//artis = append(artis, RandomArtifactOfSet("VermillionHereafter", StrongboxBase4Chance))
		}
		subs := map[stat]float32{
			ATK:      0.2,
			ATKP:     0.8,
			CritRate: 1,
			CritDmg:  1,
		}
		artis = RemoveTrashArtifacts(artis, subs, 5)

		config := optimizationConfig{
			character: c,
			target: attack{
				element:       Anemo,
				offensiveStat: ATK,
				multiplier:    404,
			},
			artifacts: artis,
		}

		_, bestTargetValue := config.findBest(func(unfiltered []*Artifact) []*Artifact {
			filtered := []*Artifact{}
			for _, art := range unfiltered {
				if art.Slot == SlotSands {
					if art.MainStat != ATKP {
						continue
					}
				}
				if art.Slot == SlotGoblet {
					if !(art.MainStat == AnemoDMG) {
						continue
					}
				}
				if art.Slot == SlotCirclet {
					if !(art.MainStat == CritRate || art.MainStat == CritDmg) {
						continue
					}
				}
				filtered = append(filtered, art)
			}
			return filtered
		})
		log.Printf("Best value: %v", bestTargetValue)
		bestTargetValueSum += bestTargetValue
	}

	log.Printf("Best value AVG: %v", bestTargetValueSum/float32(repetitions))
	t.Error("Just to make VSCode show the logs ¯\\_(ツ)_/¯")
}

func TestRandomArtifactFromDomain(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	var set1Count, set2Count int
	set1, set2 := "Emblem", "Shimenawa"

	// Generate 1000 artifacts from two sets
	for i := 0; i < 1000; i++ {
		art := RandomArtifactFromDomain(set1, set2)
		if art.Set == artifactSet(set1) {
			set1Count++
		} else if art.Set == artifactSet(set2) {
			set2Count++
		} else {
			t.Error("Unexpected artifact set: " + art.Set)
		}
	}

	// Then check that the chances of getting an artifact from either set is ~50%
	if set1Count < 450 || set1Count > 550 {
		t.Error("Too many or too few artifacts from set 1: " + strconv.Itoa(set1Count))
	}
	if set2Count < 450 || set2Count > 550 {
		t.Error("Too many or too few artifacts from set 1: " + strconv.Itoa(set2Count))
	}
}

func TestRemoveTrashArtifacts(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	var artis []*Artifact
	set1, set2 := "Emblem", "Shimenawa"

	// Generate 1000 artifacts from two sets
	for i := 0; i < 10000; i++ {
		artis = append(artis, RandomArtifactFromDomain(set1, set2))
	}

	subs := map[stat]float32{
		ATKP:           1,
		CritRate:       1,
		CritDmg:        1,
		EnergyRecharge: 0.5,
		ATK:            0.25,
	}
	filtered := RemoveTrashArtifacts(artis, subs, 5)
	for _, a := range filtered {
		t.Log(*a)
	}
}

func TestExportToGOOD(t *testing.T) {
	var artis []*Artifact

	// Random pieces of any set
	for i := 0; i < 2000; i++ {
		arti := RandomArtifact(StrongboxBase4Chance)
		if arti.Set == "MarechausseeHunter" || arti.Set == "VermillionHereafter" {
			i--
			continue
		}
		artis = append(artis, arti)
	}

	// Random pieces of a specific set
	for i := 0; i < 1782; i++ {
		artis = append(artis, RandomArtifactOfSet("MarechausseeHunter", StrongboxBase4Chance))
		artis = append(artis, RandomArtifactOfSet("VermillionHereafter", StrongboxBase4Chance))
	}

	subs := map[stat]float32{
		ATK:              0.2,
		ATKP:             0.8,
		HP:               0.2,
		HPP:              0.8,
		DEF:              0.2,
		DEFP:             0.8,
		ElementalMastery: 1,
		EnergyRecharge:   1,
		CritRate:         1,
		CritDmg:          1,
	}
	artis = RemoveTrashArtifacts(artis, subs, 10)
	export := ExportToGOOD(artis)
	b, err := json.Marshal(export)
	if err != nil {
		t.Error(err)
	}
	os.WriteFile("goodExport_"+time.Now().Format("2006-01-02_15.04.05")+".json", b, 0755)
}

func TestExportToGOODEmblemHell(t *testing.T) {
	// 6 months of Emblem -> 1620 runs
	// 6 months of Emblem -> 4320 runs if max refreshing
	// 10000 resin -> 500 runs ~> 532 artifacts
	var artis []*Artifact
	for i := 0; i < 1620; i++ {
		artis = append(artis, RandomArtifactFromDomain("EmblemOfSeveredFate", "ShimenawasReminiscence"))
	}
	subs := map[stat]float32{
		CritRate:       1,
		CritDmg:        1,
		ATKP:           0.8,
		EnergyRecharge: 0.8,
		ATK:            0.2,
	}
	artis = RemoveTrashArtifacts(artis, subs, 12)
	export := ExportToGOOD(artis)
	b, err := json.Marshal(export)
	if err != nil {
		t.Error(err)
	}
	os.WriteFile("goodExportEmblemTest.json", b, 0755)
}

func TestChancesToUpgrade(t *testing.T) {
	var targetQuality float32 = 6.46
	var domainRuns float64 = 12000 * 1.065
	var repetitions float64 = 1000
	subWeights := map[stat]float32{
		CritRate:       1,
		CritDmg:        1,
		ATKP:           0.8,
		EnergyRecharge: 0.8,
		ATK:            0.2,
	}
	upgradeCount := 0.0
	for i := 0.0; i < repetitions; i++ {
		for j := 0.0; j < domainRuns; j++ {
			art := RandomArtifactFromDomain("Emblem", "Shime")
			if art.Set != "Emblem" {
				continue
			}
			if art.MainStat != EnergyRecharge {
				continue
			}
			if art.subsQuality(subWeights) >= targetQuality {
				upgradeCount++
				break
			}
		}
	}
	t.Log(upgradeCount / repetitions)
}

func TestAvgOfDendroGoblets(t *testing.T) {
	dendroGobletCount := 0.0
	for n := 0; n < 1000; n++ {
		for i := 0; i < 308; i++ {
			art := RandomArtifactOfSet("CrimsonWitchOfFlames", StrongboxBase4Chance)
			if art.MainStat == DendroDMG {
				dendroGobletCount++
			}
		}
	}
	t.Log(dendroGobletCount / 1000.0)
}

func TestAvgOfHighCVPyroGoblets(t *testing.T) {
	godPieces := 0.0
	for n := 0; n < 1000; n++ {
		for i := 0; i < 2000; i++ {
			art := RandomArtifactOfSet("CrimsonWitchOfFlames", StrongboxBase4Chance)
			if art.MainStat == DendroDMG && art.cv() >= 30 {
				godPieces++
			}
		}
	}
	t.Log(godPieces / 1000.0)
}

func TestAvgOf40CVCritCirclets(t *testing.T) {
	godPieces := 0.0
	for n := 0; n < 1000; n++ {
		for i := 0; i < 4745; i++ {
			art := RandomArtifactOfSet("CrimsonWitchOfFlames", StrongboxBase4Chance)
			if art.MainStat == CritDmg && art.cv() >= 40 {
				godPieces++
			}
			if art.MainStat == CritRate && art.cv() >= 40 {
				godPieces++
			}
		}
	}
	t.Log(godPieces / 1000.0)
}

func TestAvgOfEMGoblets(t *testing.T) {
	count := 0.0
	runsWithoutEMGoblet := 0
	for n := 0; n < 1000; n++ {
		gotOne := false
		for i := 0; i < 210; i++ {
			art := RandomArtifactOfSet("ViridescentVenerer", StrongboxBase4Chance)
			if art.Slot == SlotGoblet && art.MainStat == ElementalMastery {
				count++
				gotOne = true
			}
		}
		if !gotOne {
			runsWithoutEMGoblet++
		}
	}
	t.Log(count / 1000.0)
	t.Log("Runs without EM Goblet:", runsWithoutEMGoblet)
}

func TestVVDomainRunsToGetEMGoblet(t *testing.T) {
	neededRuns := []int{}

	for i := 0; i < 10000; i++ {
		count := 0
		for {
			count++
			art := RandomArtifactFromDomain("VV", "Maidens")
			if art.Set == "VV" && art.MainStat == ElementalMastery && art.Slot == SlotGoblet {
				break
			}
		}
		neededRuns = append(neededRuns, count)
	}

	sort.Ints(neededRuns)

	t.Log("1% (PepeW):", neededRuns[100])
	t.Log("10% (Luckiest):", neededRuns[1000])
	t.Log("50% (Most people):", neededRuns[5000])
	t.Log("90% (Unluckiest):", neededRuns[9000])
	t.Log("99% (TrollDespair):", neededRuns[9900])
}

func TestStrongboxWithCertainMainAndSubs(t *testing.T) {
	neededRuns := []int{}

	for i := 0; i < 1000; i++ {
		count := 0
		for {
			count++
			art := RandomArtifactOfSet("CW", StrongboxBase4Chance)
			if art.MainStat == HPP && art.Slot == SlotSands {
				if !art.IsFourLiner {
					continue
				}

				wantedSubsCount := 0
				for _, sub := range art.SubStats {
					switch sub.Stat {
					case CritRate, CritDmg, ElementalMastery:
						wantedSubsCount++
					}
				}

				if wantedSubsCount == 3 {
					neededRuns = append(neededRuns, count)
					break
				}
			}
		}
	}

	sort.Ints(neededRuns)

	t.Log("10% (Luckiest):", neededRuns[100])
	t.Log("50% (Most people):", neededRuns[500])
	t.Log("90% (Unluckiest):", neededRuns[900])
}

func TestStrongbox(t *testing.T) {
	var artis []*Artifact
	for i := 0; i < 210; i++ {
		artis = append(artis, RandomArtifactOfSet("ViridescentVenerer", StrongboxBase4Chance))
	}
	export := ExportToGOOD(artis)
	b, err := json.Marshal(export)
	if err != nil {
		t.Error(err)
	}
	os.WriteFile("goodStrongbox_withDendro.json", b, 0755)
}
