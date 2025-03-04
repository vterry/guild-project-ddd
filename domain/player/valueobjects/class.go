package valueobjects

type Class int

const (
	Mage Class = iota
	Warrior
	Ranger
)

func (c Class) String() string {
	return [...]string{"MAGE", "WARRIOR", "RANGER"}[c]
}
