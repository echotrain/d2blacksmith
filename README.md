# d2blacksmith
Destiny 2 | Armor 2.0  - armor sets stats optimization

### Usage:
1. Download your armor .csv file from DIM (example: destinyArmor.csv).
1. Make sure your armor .csv file is current, and is in the same current working directory as cloned repo.
1. Build the binary: `go build`
1. Run the program: `./d2blacksmith -file=destinyArmor.csv -guardian=Titan`

#

### Example Output:
```
***
Totals -- Mob: 20 (0,low) | Res: 30 (0,low) | Rec: 80 (0,low) | Dis: 65 (5,high) | Int: 30 (0,low) | Str: 33 (3,low)
Tier: T(40) == [(Base: 31) + (Mods: 5) + (Powerful Friends*2: 4)]
Overflow: 8
Helm:
{
	"name": "Tangled Web Hood",
	"type": "Helmet",
	"rarity": "Legendary",
	"guardian": "Warlock",
	"mobility": 2,
	"resilience": 2,
	"recovery": 28,
	"discipline": 21,
	"intellect": 7,
	"strength": 2
}
Gauntlets:
{
	"name": "Seventh Seraph Gloves",
	"type": "Gauntlets",
	"rarity": "Legendary",
	"guardian": "Warlock",
	"mobility": 6,
	"resilience": 6,
	"recovery": 22,
	"discipline": 9,
	"intellect": 12,
	"strength": 12
}
Chest:
{
	"name": "Righteous Robes",
	"type": "Chest Armor",
	"rarity": "Legendary",
	"guardian": "Warlock",
	"mobility": 10,
	"resilience": 12,
	"recovery": 11,
	"discipline": 16,
	"intellect": 9,
	"strength": 7
}
Legs:
{
	"name": "Boots of Ascendancy",
	"type": "Leg Armor",
	"rarity": "Legendary",
	"guardian": "Warlock",
	"mobility": 2,
	"resilience": 10,
	"recovery": 19,
	"discipline": 19,
	"intellect": 2,
	"strength": 12
}
***

...
```

#

###### Run the following for more information on CLI options: `./d2blacksmith -help`
