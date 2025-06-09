package class

import (
	"fmt"
	"strings"
)

type Class int

const (
	Mage Class = iota
	Warrior
	Ranger
)

func (c Class) String() string {
	return [...]string{"MAGE", "WARRIOR", "RANGER"}[c]
}

func ParseClass(s string) (Class, error) {
	switch strings.ToUpper(s) {
	case "MAGE":
		return Mage, nil
	case "WARRIOR":
		return Warrior, nil
	case "RANGER":
		return Ranger, nil
	default:
		return Mage, fmt.Errorf("invalid class: %s", s)
	}
}
