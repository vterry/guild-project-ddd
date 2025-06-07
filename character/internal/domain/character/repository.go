package character

type CharacterRepository interface {
	Save(Character) error
	FindCharacterById(CharacterID) (*Character, error)
	Update(Character) error
}
