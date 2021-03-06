package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
)

var (
	fileFlag       string
	guardianFlag   string
	masterworkFlag bool
	powerfulFlag   int
	overflowFlag   int
	tierFlag       int
	modsFlag       int
)

func init() {
	fileHelp := "file name of .csv file to process. example: -file:example-armors.csv"
	flag.StringVar(&fileFlag, "file", "", fileHelp)

	guardianHelp := "class type you want to process. example: -guardian=titan"
	flag.StringVar(&guardianFlag, "guardian", "", guardianHelp)

	masterworkHelp := "assume all pieces of armor sets are masterworked. example: -masterwork=true"
	flag.BoolVar(&masterworkFlag, "masterwork", false, masterworkHelp)

	powerfulHelp := "number of powerful friends mods to apply (no more than 2). example: -powerful=1"
	flag.IntVar(&powerfulFlag, "powerful", 0, powerfulHelp)

	overflowHelp := "(<) maximum of wasted stats to allow (default 10). example: -overflow=10"
	flag.IntVar(&overflowFlag, "overflow", 10, overflowHelp)

	tierHelp := "(>=) minimum of total stat tiers to find (default 38). example: -tier=38"
	flag.IntVar(&tierFlag, "tier", 38, tierHelp)

	modsHelp := "number of stat based mods to apply (no more than 5). example: -mods=5"
	flag.IntVar(&modsFlag, "mods", 0, modsHelp)
}

type armor struct {
	Name       string `csv:"Name" json:"name"`
	Type       string `csv:"Type" json:"type"`
	Rarity     string `csv:"Tier" json:"rarity"`
	Guardian   string `csv:"Equippable" json:"guardian"`
	Mobility   int    `csv:"Mobility (Base)" json:"mobility"`
	Resilience int    `csv:"Resilience (Base)" json:"resilience"`
	Recovery   int    `csv:"Recovery (Base)" json:"recovery"`
	Discipline int    `csv:"Discipline (Base)" json:"discipline"`
	Intellect  int    `csv:"Intellect (Base)" json:"intellect"`
	Strength   int    `csv:"Strength (Base)" json:"strength"`
}

type organized struct {
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

func organize(armors []*armor) organized {
	gear := organized{}
	for _, piece := range armors {
		if strings.ToLower(piece.Guardian) == guardianFlag {
			switch piece.Type {
			case "Helmet":
				gear.Helmets = append(gear.Helmets, piece)
			case "Gauntlets":
				gear.Gauntlets = append(gear.Gauntlets, piece)
			case "Chest Armor":
				gear.Chests = append(gear.Chests, piece)
			case "Leg Armor":
				gear.Legs = append(gear.Legs, piece)
			case "Titan Mark", "Hunter Cloak", "Warlock Bond":
				continue
			default:
				fmt.Println("error: corrupted value found in .csv file")
				os.Exit(1)
			}
		}
	}

	return gear
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
	fmt.Printf(
		"Tier: T(%d) == [(Base: %d) + (Mods: %d) + (Powerful Friends*%d: %d)]\n",
		bundle.tier,
		bundle.tier-modsFlag-(powerfulFlag*2),
		modsFlag,
		powerfulFlag,
		powerfulFlag*2,
	)
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

func exotic(set []string) int {
	count := 0

	for _, item := range set {
		if item == "Exotic" {
			count++
		}
	}

	return count
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
	bundle.tier = sum(bundle.totals) + modsFlag + (powerfulFlag * 2)
	if masterworkFlag {
		bundle.tier += 6
	}
	bundle.over = overflow(bundle.remainders)

	if bundle.over < overflowFlag && bundle.tier >= tierFlag {
		printStatsFull(bundle)
	}
}

func process(gear organized) {
	helmets := gear.Helmets
	gauntlets := gear.Gauntlets
	chests := gear.Chests
	legs := gear.Legs

	for _, h := range helmets {
		for _, g := range gauntlets {
			for _, c := range chests {
				for _, l := range legs {
					if exotic([]string{h.Rarity, g.Rarity, c.Rarity, l.Rarity}) > 1 {
						continue
					}

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
	required := []string{"file", "guardian"}
	classes := map[string]bool{
		"titan":   true,
		"hunter":  true,
		"warlock": true,
	}

	flag.Parse()
	guardianFlag = strings.ToLower(guardianFlag)

	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			fmt.Fprintf(os.Stderr, "missing required -%s flag\nusage: %s\n", req, flag.Lookup(req).Usage)
			os.Exit(2)
		}
	}

	if !classes[guardianFlag] {
		fmt.Fprintf(os.Stderr, "must provide valid guardian class\nusage: %s\n", flag.Lookup("guardian").Usage)
		os.Exit(2)
	}

	if powerfulFlag > 2 {
		fmt.Fprintf(
			os.Stderr,
			"value \"%s\" exceeded maximum for -powerful flag\nusage: %s\n",
			flag.Lookup("powerful").Value.String(),
			flag.Lookup("powerful").Usage,
		)
		os.Exit(2)
	}

	if modsFlag > 5 {
		fmt.Fprintf(
			os.Stderr,
			"value \"%s\" exceeded maximum for -mods flag\nusage: %s\n",
			flag.Lookup("mods").Value.String(),
			flag.Lookup("mods").Usage,
		)
		os.Exit(2)
	}

	fmt.Println("Hello, Guardian of the Light!")

	armors, err := reader(fileFlag)
	if err != nil {
		panic(err)
	}

	process(organize(armors))
}
