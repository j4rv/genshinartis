package genshinartis

import (
	"encoding/json"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"testing"
	"time"
)

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

	subs := map[artifactStat]float32{
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
	for i := 0; i < 10000; i++ {
		artis = append(artis, RandomArtifact(DomainBase4Chance))
	}
	subs := map[artifactStat]float32{
		ATKP:             1,
		CritRate:         1,
		CritDmg:          1,
		ElementalMastery: 0.33,
		EnergyRecharge:   0.25,
		ATK:              0.25,
	}
	artis = RemoveTrashArtifacts(artis, subs, 10)
	export := ExportToGOOD(artis)
	b, err := json.Marshal(export)
	if err != nil {
		t.Error(err)
	}
	os.WriteFile("goodExportTest.json", b, 0755)
}

func TestExportToGOODEmblemHell(t *testing.T) {
	// 6 months of Emblem -> 1620 runs
	// 6 months of Emblem -> 4320 runs if max refreshing
	// 10000 resin -> 500 runs ~> 532 artifacts
	var artis []*Artifact
	for i := 0; i < 1620; i++ {
		artis = append(artis, RandomArtifactFromDomain("EmblemOfSeveredFate", "ShimenawasReminiscence"))
	}
	subs := map[artifactStat]float32{
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
	subWeights := map[artifactStat]float32{
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
				neededRuns = append(neededRuns, count)
				break
			}
		}
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
