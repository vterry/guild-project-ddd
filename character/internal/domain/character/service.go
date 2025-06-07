package character

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/vterry/ddd-study/character/internal/domain/common/class"
	"github.com/vterry/ddd-study/character/internal/domain/common/guild"
	"github.com/vterry/ddd-study/character/internal/domain/common/vault"
	"github.com/vterry/ddd-study/character/internal/domain/login"
	"github.com/vterry/ddd-study/character/internal/domain/playeritem"
)

var (
	ErrWhileCreation = errors.New("an error occur during creating")
)

type VaultService interface {
	CreateVault() (vault.VaultID, error)
}

type GuildService interface {
	LeaveGuild(characterId CharacterID) error
}

type CharacterService interface {
	CreateCharacter(userId string, email string, nickname string, class class.Class) error
	TransferItemTo(characterId CharacterID, playeritem playeritem.PlayerItem, quantity int, vaultId vault.VaultID) error
	TradeItem(origin CharacterID, playeritem playeritem.PlayerItem, quantity int, destiny CharacterID) error
	DepositGold(characterId CharacterID, quantity int, vaultId vault.VaultID) error
	PickItem() error
	DropItem(playerItemID playeritem.PlayerItemID, quantity int) error
	LeaveGuild(characterID CharacterID) error
}

type CharacterServiceImpl struct {
	vaultService VaultService
	repository   CharacterRepository
}

func NewCharacterService(repository CharacterRepository, vaultService VaultService) CharacterService {
	return &CharacterServiceImpl{
		vaultService: vaultService,
		repository:   repository,
	}
}

func (s *CharacterServiceImpl) CreateCharacter(userId string, email string, nickname string, class class.Class) error {
	login, err := login.NewLogin(userId, email)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrWhileCreation, err)
	}

	vaultId, err := s.vaultService.CreateVault()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrWhileCreation, err)
	}

	character, err := CreateNewCharacter(nickname, login, class, vaultId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrWhileCreation, err)
	}
	return s.repository.Save(*character)
}

func (s *CharacterServiceImpl) TransferItemTo(characterId CharacterID, playeritem playeritem.PlayerItem, quantity int, vaultId vault.VaultID) error {
	character, err := s.repository.FindCharacterById(characterId)
	if err != nil {
		return fmt.Errorf("failed to find character: %w", err)
	}

	// Verify the character owns the vault
	if !character.GetCurrentVaultId().Equals(vaultId) {
		return errors.New("character does not own this vault")
	}

	// Drop the item from character's inventory
	if err := character.DropItem(playeritem); err != nil {
		return fmt.Errorf("failed to drop item from inventory: %w", err)
	}

	// Save the updated character state
	if err := s.repository.Update(*character); err != nil {
		return fmt.Errorf("failed to update character: %w", err)
	}

	return nil
}

func (s *CharacterServiceImpl) TradeItem(origin CharacterID, playeritem playeritem.PlayerItem, quantity int, destiny CharacterID) error {
	// Get origin character
	originChar, err := s.repository.FindCharacterById(origin)
	if err != nil {
		return fmt.Errorf("failed to find origin character: %w", err)
	}

	// Get destination character
	destinyChar, err := s.repository.FindCharacterById(destiny)
	if err != nil {
		return fmt.Errorf("failed to find destination character: %w", err)
	}

	// Drop item from origin character
	if err := originChar.DropItem(playeritem); err != nil {
		return fmt.Errorf("failed to drop item from origin character: %w", err)
	}

	// Pick item with destination character
	if err := destinyChar.PickItem(playeritem); err != nil {
		// If picking fails, try to return the item to origin character
		_ = originChar.PickItem(playeritem)
		return fmt.Errorf("failed to pick item with destination character: %w", err)
	}

	// Save both characters' states
	if err := s.repository.Update(*originChar); err != nil {
		// Rollback: remove the item from destiny character
		_ = destinyChar.DropItem(playeritem)
		return fmt.Errorf("failed to update origin character: %w", err)
	}
	if err := s.repository.Update(*destinyChar); err != nil {
		return fmt.Errorf("failed to update destination character: %w", err)
	}
	return nil
}

func (s *CharacterServiceImpl) DepositGold(characterId CharacterID, quantity int, vaultId vault.VaultID) error {
	character, err := s.repository.FindCharacterById(characterId)
	if err != nil {
		return fmt.Errorf("failed to find character: %w", err)
	}

	// Verify the character owns the vault
	if !character.GetCurrentVaultId().Equals(vaultId) {
		return errors.New("character does not own this vault")
	}

	// Drop gold from character's inventory
	if err := character.DropGold(quantity); err != nil {
		return fmt.Errorf("failed to drop gold from inventory: %w", err)
	}

	// Save the updated character state
	if err := s.repository.Update(*character); err != nil {
		return fmt.Errorf("failed to update character: %w", err)
	}

	return nil
}

func (s *CharacterServiceImpl) PickItem() error {
	return errors.New("method not implemented - requires character ID and item parameters")
}

func (s *CharacterServiceImpl) DropItem(playerItemID playeritem.PlayerItemID, quantity int) error {
	return errors.New("method not implemented - requires character ID parameter")
}

func (s *CharacterServiceImpl) LeaveGuild(characterID CharacterID) error {
	character, err := s.repository.FindCharacterById(characterID)
	if err != nil {
		return fmt.Errorf("failed to find character: %w", err)
	}

	// Update guild info to empty guild
	character.UpdateGuildInfo(guild.New(uuid.Nil))

	// Save the updated character state
	if err := s.repository.Update(*character); err != nil {
		return fmt.Errorf("failed to update character: %w", err)
	}

	return nil
}
