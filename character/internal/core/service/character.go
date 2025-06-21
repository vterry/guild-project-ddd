package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/vterry/ddd-study/character/internal/core/domain/character"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/class"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/guild"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/login"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/vault"
	"github.com/vterry/ddd-study/character/internal/core/domain/playeritem"
	"github.com/vterry/ddd-study/character/internal/core/ports/input/service"
	"github.com/vterry/ddd-study/character/internal/core/ports/output/gateway"
	"github.com/vterry/ddd-study/character/internal/core/ports/output/repository"
)

var (
	ErrWhileCreation = errors.New("an error occur during creating")
)

type CharacterServiceImpl struct {
	vaultGateway        gateway.Vault
	characterRepository repository.CharacterRepository
}

func NewCharacterService(characterRepository repository.CharacterRepository, vaultGateway gateway.Vault) service.CharacterService {
	return &CharacterServiceImpl{
		vaultGateway:        vaultGateway,
		characterRepository: characterRepository,
	}
}

func (s *CharacterServiceImpl) CreateCharacter(ctx context.Context, loginId login.LoginID, nickname string, class class.Class) error {

	vaultId, err := s.vaultGateway.CreateVault()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrWhileCreation, err)
	}

	character, err := character.CreateNewCharacter(nickname, loginId, class, vaultId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrWhileCreation, err)
	}

	return s.characterRepository.Save(ctx, *character)
}

func (s *CharacterServiceImpl) TransferItemTo(ctx context.Context, characterId character.CharacterID, playeritem playeritem.PlayerItem, quantity int, vaultId vault.VaultID) error {
	character, err := s.characterRepository.FindCharacterById(ctx, characterId)
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
	if err := s.characterRepository.Update(ctx, *character); err != nil {
		return fmt.Errorf("failed to update character: %w", err)
	}

	return nil
}

func (s *CharacterServiceImpl) TradeItem(ctx context.Context, origin character.CharacterID, playeritem playeritem.PlayerItem, quantity int, destiny character.CharacterID) error {
	// Get origin character
	originChar, err := s.characterRepository.FindCharacterById(ctx, origin)
	if err != nil {
		return fmt.Errorf("failed to find origin character: %w", err)
	}

	// Get destination character
	destinyChar, err := s.characterRepository.FindCharacterById(ctx, destiny)
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
	if err := s.characterRepository.Update(ctx, *originChar); err != nil {
		// Rollback: remove the item from destiny character
		_ = destinyChar.DropItem(playeritem)
		return fmt.Errorf("failed to update origin character: %w", err)
	}
	if err := s.characterRepository.Update(ctx, *destinyChar); err != nil {
		return fmt.Errorf("failed to update destination character: %w", err)
	}
	return nil
}

func (s *CharacterServiceImpl) DepositGold(ctx context.Context, characterId character.CharacterID, quantity int, vaultId vault.VaultID) error {
	character, err := s.characterRepository.FindCharacterById(ctx, characterId)
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
	if err := s.characterRepository.Update(ctx, *character); err != nil {
		return fmt.Errorf("failed to update character: %w", err)
	}

	return nil
}

func (s *CharacterServiceImpl) PickItem() error {
	return errors.New("method not implemented - requires character ID and item parameters")
}

func (s *CharacterServiceImpl) DropItem(ctx context.Context, playerItemID playeritem.PlayerItemID, quantity int) error {
	return errors.New("method not implemented - requires character ID parameter")
}

func (s *CharacterServiceImpl) LeaveGuild(ctx context.Context, characterID character.CharacterID) error {
	character, err := s.characterRepository.FindCharacterById(ctx, characterID)
	if err != nil {
		return fmt.Errorf("failed to find character: %w", err)
	}

	// Update guild info to empty guild
	character.UpdateGuildInfo(guild.NewGuildID(uuid.Nil))

	// Save the updated character state
	if err := s.characterRepository.Update(ctx, *character); err != nil {
		return fmt.Errorf("failed to update character: %w", err)
	}

	return nil
}
