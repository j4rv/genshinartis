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

func TestXiaoSets(t *testing.T) {
	piecesToGenerate := 1725
	repetitions := 50
	set := "VermillionHereafter" // MarechausseeHunter / VermillionHereafter
	minER := float32(140)

	c := character{
		level:   90,
		baseAtk: 349,
		weapon:  weaponHomaPassiveOff,
		bonusStats: map[stat]float32{
			ATK:             1050.8,           // Benny. 1203 for Aquila, 1050.8 for Sapwood.
			ATKP:            15,               // Tenacity, Noblesse, Pyro resonance, TTDS, etc
			BaseDMGIncrease: 208.27,           // Faru A4
			AnemoDMG:        32.4 + 95.2 + 15, // Faruzan, Xiao burst and Xiao A1
			CritRate:        19.2,             // Xiao main stat
			//CritDmg:         40,                // Faruzan c6
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
			ATK:            0.2,
			ATKP:           0.8,
			EnergyRecharge: 1,
			CritRate:       1,
			CritDmg:        1,
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

		artifactFilter := func(unfiltered []*Artifact) []*Artifact {
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
		}

		buildFilter := func(build map[artifactSlot]*Artifact) bool {
			setCount := 0
			var er float32
			for _, art := range build {
				if art.Set == artifactSet(set) {
					setCount++
				}
				for _, sub := range art.SubStats {
					if sub.Stat == EnergyRecharge {
						er += sub.Value
					}
				}
			}
			if setCount < 4 {
				return false
			}
			if 100+er < minER {
				return false
			}
			return true
		}

		_, bestTargetValue := config.findBest(artifactFilter, buildFilter)
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

func TestTimeToFarmTargetRV(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	var artis []*Artifact
	set1, set2 := "Emblem", "Shimenawa"
	targetRV := float32(26 * 0.85)
	minER := float32(100)
	iterations := 1000
	neededRuns := []int{}

	rvMultiplier := map[stat]float32{
		ATKP:     1,
		CritRate: 1,
		CritDmg:  1,
		ATK:      0.25,
	}

	validArtifact := func(art *Artifact) bool {
		if art.Slot == SlotSands {
			if art.MainStat != ATKP {
				return false
			}
		}
		if art.Slot == SlotGoblet {
			if !(art.MainStat == AnemoDMG) {
				return false
			}
		}
		if art.Slot == SlotCirclet {
			if !(art.MainStat == CritRate || art.MainStat == CritDmg) {
				return false
			}
		}
		return true
	}

	buildFilter := func(build map[artifactSlot]*Artifact) bool {
		setCount := 0
		var er float32
		for _, art := range build {
			if art.Set == artifactSet(set1) {
				setCount++
			}
			for _, sub := range art.SubStats {
				if sub.Stat == EnergyRecharge {
					er += sub.Value
				}
			}
		}
		if setCount < 4 {
			return false
		}
		if 100+er < minER {
			return false
		}
		return true
	}

	domainRuns := 0
	for i := 0; i < iterations; {
		// one domain run
		domainRuns++
		art := RandomArtifactFromDomain(set1, set2)
		if validArtifact(art) {
			artis = append(artis, art)
		}
		if rand.Float32() <= DomainExtraArtifactChance {
			art = RandomArtifactFromDomain(set1, set2)
			if validArtifact(art) {
				artis = append(artis, art)
			}
		}

		// some cleaning
		artis = RemoveTrashArtifacts(artis, rvMultiplier, 1)

		// check target RV
		_, rv := findHighestRV(artis, rvMultiplier, nil, buildFilter)
		//log.Printf("Iteration %d, %d domain runs done, current max rv: %f, target: %f\n", i, domainRuns, rv, targetRV)
		if rv >= targetRV {
			log.Printf("Iteration %d, %d domain runs needed\n", i, domainRuns)
			neededRuns = append(neededRuns, domainRuns)
			// reset for next iteration
			artis = nil
			domainRuns = 0
			i++
		}
	}

	sort.Ints(neededRuns)
	t.Log("1% (PepeW), domain runs needed:", neededRuns[10])
	t.Log("10% (Luckiest), domain runs needed:", neededRuns[100])
	t.Log("50% (Most people), domain runs needed:", neededRuns[500])
	t.Log("90% (Unluckiest), domain runs needed:", neededRuns[900])
	t.Log("99% (TrollDespair), domain runs needed:", neededRuns[990])
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
	//for i := 0; i < 862*2; i++ {
	//	arti := RandomArtifact(StrongboxBase4Chance)
	//	if arti.Set == "MarechausseeHunter" || arti.Set == "VermillionHereafter" {
	//		i--
	//		continue
	//	}
	//	artis = append(artis, arti)
	//}

	// Random pieces of a specific set
	for i := 0; i < 200; i++ {
		artis = append(artis, RandomArtifactOfSet("MarechausseeHunter", DomainBase4Chance))
	}

	subs := map[stat]float32{
		ATK:            0.2,
		ATKP:           1,
		EnergyRecharge: 0.5,
		CritRate:       0.8,
		CritDmg:        1,
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
