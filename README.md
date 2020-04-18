# d2blacksmith
Destiny 2 | Armor 2.0  - armor sets stats optimization


### Usage:
1. Ensure your .csv file headers are correct (see example-armors.csv). Please
note: enter armor pieces with *_base_* armor stats, *_do not_* include masterwork bonus, or armor mods.
1. Build the binary: `go build`
1. Run the program: `./d2blacksmith example-armors.csv`
#

### Example Output:
```
***
Totals -- Mob: 30 (0,low) | Res: 50 (0,low) | Rec: 45 (5,high) | Dis: 60 (0,low) | Int: 32 (2,low) | Str: 30 (0,low)
Tier: T(24)
Overflow: 7
Helm:
{
	"name": "Prodigal Helm",
	"type": "helmet",
	"rarity": "legendary",
	"mobility": 2,
	"resilience": 8,
	"recovery": 19,
	"discipline": 12,
	"intellect": 12,
	"strength": 6
}
Gauntlets:
{
	"name": "Iron Rememberance Gauntlets",
	"type": "gauntlets",
	"rarity": "legendary",
	"mobility": 12,
	"resilience": 10,
	"recovery": 10,
	"discipline": 13,
	"intellect": 2,
	"strength": 15
}
Chest:
{
	"name": "Righteous Plate",
	"type": "chest",
	"rarity": "legendary",
	"mobility": 6,
	"resilience": 22,
	"recovery": 6,
	"discipline": 16,
	"intellect": 9,
	"strength": 7
}
Legs:
{
	"name": "Dunemarchers",
	"type": "leg",
	"rarity": "exotic",
	"mobility": 10,
	"resilience": 10,
	"recovery": 10,
	"discipline": 19,
	"intellect": 9,
	"strength": 2
}
***

...
```
