package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
)

type armor struct {
	Name       string `csv:"name" json:"name"`
	Type       string `csv:"type" json:"type"`
	Rarity     string `csv:"rarity" json:"rarity"`
	Mobility   int    `csv:"mobility" json:"mobility"`
	Resilience int    `csv:"resilience" json:"resilience"`
	Recovery   int    `csv:"recovery" json:"recovery"`
	Discipline int    `csv:"discipline" json:"discipline"`
	Intellect  int    `csv:"intellect" json:"intellect"`
	Strength   int    `csv:"strength" json:"strength"`
}

type inventory struct {
	Helmets   []*armor
	Gauntlets []*armor
	Chests    []*armor
	Legs      []*armor
}

func reader(filePath string) ([]*armor, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	gear := []*armor{}

	if err := gocsv.UnmarshalFile(f, &gear); err != nil {
		return nil, err
	}

	return gear, nil
}

func organizeGear(gear []*armor) inventory {
	inv := inventory{}
	for _, piece := range gear {
		switch piece.Type {
		case "helmet":
			inv.Helmets = append(inv.Helmets, piece)
		case "gauntlets":
			inv.Gauntlets = append(inv.Gauntlets, piece)
		case "chest":
			inv.Chests = append(inv.Chests, piece)
		case "leg":
			inv.Legs = append(inv.Legs, piece)
		default:
			break
		}
	}

	return inv
}

type remainder struct {
	val int
	low bool
}

type armorRemainders struct {
	Mobility   remainder
	Resilience remainder
	Recovery   remainder
	Discipline remainder
	Intellect  remainder
	Strength   remainder
}

type stats struct {
	totals     *armor
	helm       *armor
	gauntlets  *armor
	chest      *armor
	legs       *armor
	remainders armorRemainders
	tier       int
	over       int
}

func modulus(stat int) remainder {
	val := stat % 10
	return remainder{val: val, low: val < 5}
}

func sum(totals *armor) int {
	tier := int(totals.Mobility / 10)
	tier += int(totals.Resilience / 10)
	tier += int(totals.Recovery / 10)
	tier += int(totals.Discipline / 10)
	tier += int(totals.Intellect / 10)
	tier += int(totals.Strength / 10)
	return tier
}

func overflow(remainders armorRemainders) int {
	over := remainders.Mobility.val
	over += remainders.Resilience.val
	over += remainders.Recovery.val
	over += remainders.Discipline.val
	over += remainders.Intellect.val
	over += remainders.Strength.val
	return over
}

func parse(stat int, remainders remainder) string {
	if remainders.low {
		return fmt.Sprintf("%d (%d,low)", stat, remainders.val)
	}
	return fmt.Sprintf("%d (%d,high)", stat, remainders.val)
}

func printStats(remainders armorRemainders, totals *armor) {
	fmt.Printf(
		"Mob: %s | Res: %s | Rec: %s | Dis: %s | Int: %s | Str: %s\n",
		parse(totals.Mobility, remainders.Mobility),
		parse(totals.Resilience, remainders.Resilience),
		parse(totals.Recovery, remainders.Recovery),
		parse(totals.Discipline, remainders.Discipline),
		parse(totals.Intellect, remainders.Intellect),
		parse(totals.Strength, remainders.Strength),
	)
}

func printStatsFull(bundle stats) {
	fmt.Printf("*** \nTotals -- ")
	printStats(bundle.remainders, bundle.totals)
	fmt.Printf("Tier: T(%d)\n", bundle.tier)
	fmt.Printf("Overflow: %d\n", bundle.over)
	helmJSON, _ := json.MarshalIndent(bundle.helm, "", "\t")
	fmt.Printf("Helm:\n%s\n", string(helmJSON))
	gauntletsJSON, _ := json.MarshalIndent(bundle.gauntlets, "", "\t")
	fmt.Printf("Gauntlets:\n%s\n", string(gauntletsJSON))
	chestJSON, _ := json.MarshalIndent(bundle.chest, "", "\t")
	fmt.Printf("Chest:\n%s\n", string(chestJSON))
	legsJSON, _ := json.MarshalIndent(bundle.legs, "", "\t")
	fmt.Printf("Legs:\n%s\n", string(legsJSON))
	fmt.Println("***")
}

func traction(bundle *stats) {
	bundle.legs.Name += "*** TRACTION MOD APPLIED ***"
	bundle.legs.Mobility += 5
	bundle.totals.Mobility += 5
}

func analyze(bundle stats) {
	bundle.remainders = armorRemainders{
		Mobility:   modulus(bundle.totals.Mobility),
		Resilience: modulus(bundle.totals.Resilience),
		Recovery:   modulus(bundle.totals.Recovery),
		Discipline: modulus(bundle.totals.Discipline),
		Intellect:  modulus(bundle.totals.Intellect),
		Strength:   modulus(bundle.totals.Strength),
	}
	bundle.tier = sum(bundle.totals)
	bundle.over = overflow(bundle.remainders)

	if bundle.over <= 13 {
		printStatsFull(bundle)
	} else if bundle.tier > 23 {
		printStatsFull(bundle)
	} else if !bundle.remainders.Mobility.low {
		// TODO: traction func not working correctly
		// traction(&bundle)
		// analyze(bundle)
	}
}

func process(gear inventory) {
	helmets := gear.Helmets
	gauntlets := gear.Gauntlets
	chests := gear.Chests
	legs := gear.Legs

	for _, h := range helmets {
		for _, g := range gauntlets {
			for _, c := range chests {
				for _, l := range legs {
					bundle := stats{
						totals: &armor{
							Mobility:   h.Mobility + g.Mobility + c.Mobility + l.Mobility,
							Resilience: h.Resilience + g.Resilience + c.Resilience + l.Resilience,
							Recovery:   h.Recovery + g.Recovery + c.Recovery + l.Recovery,
							Discipline: h.Discipline + g.Discipline + c.Discipline + l.Discipline,
							Intellect:  h.Intellect + g.Intellect + c.Intellect + l.Intellect,
							Strength:   h.Strength + g.Strength + c.Strength + l.Strength,
						},
						helm:      h,
						gauntlets: g,
						chest:     c,
						legs:      l,
					}

					analyze(bundle)
				}
			}
		}
	}
}

func main() {
	fmt.Println("Hello, Guardian of the Light!")

	armors, err := reader(os.Args[1])
	if err != nil {
		panic(err)
	}

	gear := organizeGear(armors)

	process(gear)
}
