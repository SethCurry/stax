package stax

func AllTypes() []Type {
	return []Type{
		Type("Land"),
		Type("Creature"),
		Type("Artifact"),
		Type("Enchantment"),
		Type("Planeswalker"),
		Type("Instant"),
		Type("Sorcery"),
		Type("Battle"),
		Type("Dungeon"),
		Type("Conspiracy"),
		Type("Phenomenon"),
		Type("Plane"),
		Type("Scheme"),
		Type("Tribal"),
		Type("Vanguard"),
	}
}

func IsType(toTest string) bool {
	for _, v := range AllTypes() {
		if toTest == string(v) {
			return true
		}
	}

	return false
}
