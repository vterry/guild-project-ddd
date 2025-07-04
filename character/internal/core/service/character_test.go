package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vterry/ddd-study/character/internal/core/domain/character"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/class"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/item"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/login"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/vault"
	"github.com/vterry/ddd-study/character/internal/core/domain/playeritem"
)

// MockVaultService is a mock implementation of VaultService
type MockVaultService struct {
	mock.Mock
}

func (m *MockVaultService) CreateVault() (vault.VaultID, error) {
	args := m.Called()
	return args.Get(0).(vault.VaultID), args.Error(1)
}

// MockCharacterRepository is a mock implementation of CharacterRepository
type MockCharacterRepository struct {
	mock.Mock
}

func (m *MockCharacterRepository) Save(ctx context.Context, character character.Character) error {
	args := m.Called(ctx, character)
	return args.Error(0)
}

func (m *MockCharacterRepository) FindCharacterById(ctx context.Context, id character.CharacterID) (*character.Character, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*character.Character), args.Error(1)
}

func (m *MockCharacterRepository) Update(ctx context.Context, character character.Character) error {
	args := m.Called(ctx, character)
	return args.Error(0)
}

// MockLogger is a stub implementation of logger.Logger for tests
// It does nothing

type MockLogger struct{}

func (m *MockLogger) Info(msg string, args ...interface{})  {}
func (m *MockLogger) Warn(msg string, args ...interface{})  {}
func (m *MockLogger) Error(msg string, args ...interface{}) {}
func (m *MockLogger) Debug(msg string, args ...interface{}) {}

func TestCreateCharacter(t *testing.T) {
	validLogin := login.NewLoginID(uuid.New())

	tests := []struct {
		name       string
		loginID    login.LoginID
		nickname   string
		class      class.Class
		setupMocks func(*MockVaultService, *MockCharacterRepository)
		wantErr    bool
	}{
		{
			name:     "successful character creation",
			loginID:  validLogin,
			nickname: "TestChar",
			class:    class.Warrior,
			setupMocks: func(vs *MockVaultService, cr *MockCharacterRepository) {
				vs.On("CreateVault").Return(vault.NewVaultID(uuid.New()), nil)
				cr.On("Save", mock.Anything, mock.AnythingOfType("Character")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "invalid login id",
			loginID:  login.LoginID{},
			nickname: "TestChar",
			class:    class.Warrior,
			setupMocks: func(vs *MockVaultService, cr *MockCharacterRepository) {
				// Even though we expect validation to fail, we should still set up the mock
				// in case the validation changes in the future
				vs.On("CreateVault").Return(vault.NewVaultID(uuid.New()), nil)
			},
			wantErr: true,
		},
		{
			name:     "vault creation failure",
			loginID:  validLogin,
			nickname: "TestChar",
			class:    class.Warrior,
			setupMocks: func(vs *MockVaultService, cr *MockCharacterRepository) {
				vs.On("CreateVault").Return(vault.NewVaultID(uuid.Nil), assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVaultService := new(MockVaultService)
			mockRepo := new(MockCharacterRepository)
			tt.setupMocks(mockVaultService, mockRepo)

			service := NewCharacterService(mockRepo, mockVaultService, &MockLogger{})
			err := service.CreateCharacter(context.Background(), tt.loginID, tt.nickname, tt.class)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockVaultService.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTransferItemTo(t *testing.T) {
	characterID := character.NewCharacterID(uuid.New())
	vaultID := vault.NewVaultID(uuid.New())
	itemID := item.NewItemID(uuid.New())
	testItem, _ := playeritem.NewPlayerItem(itemID, "Test Item", 1)

	tests := []struct {
		name       string
		setupMocks func(*MockCharacterRepository)
		wantErr    bool
	}{
		{
			name: "successful item transfer",
			setupMocks: func(cr *MockCharacterRepository) {
				loginID := login.NewLoginID(uuid.New())
				character, _ := character.CreateNewCharacter("TestChar", loginID, class.Warrior, vaultID)
				// Add item to character's inventory
				_ = character.PickItem(*testItem)
				cr.On("FindCharacterById", mock.Anything, characterID).Return(character, nil)
				cr.On("Update", mock.Anything, mock.AnythingOfType("Character")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "character not found",
			setupMocks: func(cr *MockCharacterRepository) {
				cr.On("FindCharacterById", mock.Anything, characterID).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "wrong vault",
			setupMocks: func(cr *MockCharacterRepository) {
				loginID := login.NewLoginID(uuid.New())
				character, _ := character.CreateNewCharacter("TestChar", loginID, class.Warrior, vault.NewVaultID(uuid.New()))
				cr.On("FindCharacterById", mock.Anything, characterID).Return(character, nil)
			},
			wantErr: true,
		},
		{
			name: "failed item drop",
			setupMocks: func(cr *MockCharacterRepository) {
				loginID := login.NewLoginID(uuid.New())
				character, _ := character.CreateNewCharacter("TestChar", loginID, class.Warrior, vaultID)
				// Don't add item to inventory, so drop will fail
				cr.On("FindCharacterById", mock.Anything, characterID).Return(character, nil)
			},
			wantErr: true,
		},
		{
			name: "failed character update",
			setupMocks: func(cr *MockCharacterRepository) {
				loginID := login.NewLoginID(uuid.New())
				character, _ := character.CreateNewCharacter("TestChar", loginID, class.Warrior, vaultID)
				_ = character.PickItem(*testItem)
				cr.On("FindCharacterById", mock.Anything, characterID).Return(character, nil)
				cr.On("Update", mock.Anything, mock.AnythingOfType("Character")).Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCharacterRepository)
			tt.setupMocks(mockRepo)

			service := NewCharacterService(mockRepo, nil, &MockLogger{})
			err := service.TransferItemTo(context.Background(), characterID, *testItem, 1, vaultID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTradeItem(t *testing.T) {
	originID := character.NewCharacterID(uuid.New())
	destinyID := character.NewCharacterID(uuid.New())
	itemID := item.NewItemID(uuid.New())
	testItem, _ := playeritem.NewPlayerItem(itemID, "Test Item", 1)

	tests := []struct {
		name       string
		setupMocks func(*MockCharacterRepository)
		wantErr    bool
	}{
		{
			name: "successful trade",
			setupMocks: func(cr *MockCharacterRepository) {
				loginID := login.NewLoginID(uuid.New())
				originChar, _ := character.CreateNewCharacter("OriginChar", loginID, class.Warrior, vault.NewVaultID(uuid.New()))
				destinyChar, _ := character.CreateNewCharacter("DestinyChar", loginID, class.Warrior, vault.NewVaultID(uuid.New()))
				// Add item to origin character's inventory
				_ = originChar.PickItem(*testItem)
				cr.On("FindCharacterById", mock.Anything, originID).Return(originChar, nil)
				cr.On("FindCharacterById", mock.Anything, destinyID).Return(destinyChar, nil)
				cr.On("Update", mock.Anything, mock.AnythingOfType("Character")).Return(nil).Times(2)
			},
			wantErr: false,
		},
		{
			name: "origin character not found",
			setupMocks: func(cr *MockCharacterRepository) {
				cr.On("FindCharacterById", mock.Anything, originID).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "destiny character not found",
			setupMocks: func(cr *MockCharacterRepository) {
				loginID := login.NewLoginID(uuid.New())
				originChar, _ := character.CreateNewCharacter("OriginChar", loginID, class.Warrior, vault.NewVaultID(uuid.New()))
				cr.On("FindCharacterById", mock.Anything, originID).Return(originChar, nil)
				cr.On("FindCharacterById", mock.Anything, destinyID).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed item drop from origin",
			setupMocks: func(cr *MockCharacterRepository) {
				loginID := login.NewLoginID(uuid.New())
				originChar, _ := character.CreateNewCharacter("OriginChar", loginID, class.Warrior, vault.NewVaultID(uuid.New()))
				destinyChar, _ := character.CreateNewCharacter("DestinyChar", loginID, class.Warrior, vault.NewVaultID(uuid.New()))
				// Don't add item to origin character's inventory
				cr.On("FindCharacterById", mock.Anything, originID).Return(originChar, nil)
				cr.On("FindCharacterById", mock.Anything, destinyID).Return(destinyChar, nil)
			},
			wantErr: true,
		},
		{
			name: "failed item pick with destiny",
			setupMocks: func(cr *MockCharacterRepository) {
				loginID := login.NewLoginID(uuid.New())
				originChar, _ := character.CreateNewCharacter("OriginChar", loginID, class.Warrior, vault.NewVaultID(uuid.New()))
				destinyChar, _ := character.CreateNewCharacter("DestinyChar", loginID, class.Warrior, vault.NewVaultID(uuid.New()))
				// Fill destiny character's inventory
				for i := 0; i < 10; i++ {
					item, _ := playeritem.NewPlayerItem(item.NewItemID(uuid.New()), fmt.Sprintf("Item%d", i), 1)
					_ = destinyChar.PickItem(*item)
				}
				_ = originChar.PickItem(*testItem)
				cr.On("FindCharacterById", mock.Anything, originID).Return(originChar, nil)
				cr.On("FindCharacterById", mock.Anything, destinyID).Return(destinyChar, nil)
			},
			wantErr: true,
		},
		{
			name: "failed origin character update",
			setupMocks: func(cr *MockCharacterRepository) {
				loginID := login.NewLoginID(uuid.New())
				originChar, _ := character.CreateNewCharacter("OriginChar", loginID, class.Warrior, vault.NewVaultID(uuid.New()))
				destinyChar, _ := character.CreateNewCharacter("DestinyChar", loginID, class.Warrior, vault.NewVaultID(uuid.New()))
				_ = originChar.PickItem(*testItem)
				cr.On("FindCharacterById", mock.Anything, originID).Return(originChar, nil)
				cr.On("FindCharacterById", mock.Anything, destinyID).Return(destinyChar, nil)
				cr.On("Update", mock.Anything, mock.AnythingOfType("Character")).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed destiny character update",
			setupMocks: func(cr *MockCharacterRepository) {
				loginID := login.NewLoginID(uuid.New())
				originChar, _ := character.CreateNewCharacter("OriginChar", loginID, class.Warrior, vault.NewVaultID(uuid.New()))
				destinyChar, _ := character.CreateNewCharacter("DestinyChar", loginID, class.Warrior, vault.NewVaultID(uuid.New()))
				_ = originChar.PickItem(*testItem)
				cr.On("FindCharacterById", mock.Anything, originID).Return(originChar, nil)
				cr.On("FindCharacterById", mock.Anything, destinyID).Return(destinyChar, nil)
				cr.On("Update", mock.Anything, mock.AnythingOfType("Character")).Return(nil).Once()
				cr.On("Update", mock.Anything, mock.AnythingOfType("Character")).Return(assert.AnError).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCharacterRepository)
			tt.setupMocks(mockRepo)

			service := NewCharacterService(mockRepo, nil, &MockLogger{})
			err := service.TradeItem(context.Background(), originID, *testItem, 1, destinyID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDepositGold(t *testing.T) {
	characterID := character.NewCharacterID(uuid.New())
	vaultID := vault.NewVaultID(uuid.New())

	tests := []struct {
		name       string
		quantity   int
		setupMocks func(*MockCharacterRepository)
		wantErr    bool
	}{
		{
			name:     "successful gold deposit",
			quantity: 100,
			setupMocks: func(cr *MockCharacterRepository) {
				loginID := login.NewLoginID(uuid.New())
				character, _ := character.CreateNewCharacter("TestChar", loginID, class.Warrior, vaultID)
				// Add gold to character's inventory
				err := character.PickGold(200) // Add more than we want to deposit
				assert.NoError(t, err)
				cr.On("FindCharacterById", mock.Anything, characterID).Return(character, nil)
				cr.On("Update", mock.Anything, mock.AnythingOfType("Character")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "character not found",
			quantity: 100,
			setupMocks: func(cr *MockCharacterRepository) {
				cr.On("FindCharacterById", mock.Anything, characterID).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
		{
			name:     "wrong vault",
			quantity: 100,
			setupMocks: func(cr *MockCharacterRepository) {
				loginID := login.NewLoginID(uuid.New())
				character, _ := character.CreateNewCharacter("TestChar", loginID, class.Warrior, vault.NewVaultID(uuid.New()))
				cr.On("FindCharacterById", mock.Anything, characterID).Return(character, nil)
			},
			wantErr: true,
		},
		{
			name:     "failed gold withdrawal",
			quantity: 100,
			setupMocks: func(cr *MockCharacterRepository) {
				loginID := login.NewLoginID(uuid.New())
				character, _ := character.CreateNewCharacter("TestChar", loginID, class.Warrior, vaultID)
				// Don't add any gold to character's inventory
				cr.On("FindCharacterById", mock.Anything, characterID).Return(character, nil)
			},
			wantErr: true,
		},
		{
			name:     "failed character update",
			quantity: 100,
			setupMocks: func(cr *MockCharacterRepository) {
				loginID := login.NewLoginID(uuid.New())
				character, _ := character.CreateNewCharacter("TestChar", loginID, class.Warrior, vaultID)
				err := character.PickGold(200)
				assert.NoError(t, err)
				cr.On("FindCharacterById", mock.Anything, characterID).Return(character, nil)
				cr.On("Update", mock.Anything, mock.AnythingOfType("Character")).Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCharacterRepository)
			tt.setupMocks(mockRepo)

			service := NewCharacterService(mockRepo, nil, &MockLogger{})
			err := service.DepositGold(context.Background(), characterID, tt.quantity, vaultID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLeaveGuild(t *testing.T) {
	characterID := character.NewCharacterID(uuid.New())

	tests := []struct {
		name       string
		setupMocks func(*MockCharacterRepository)
		wantErr    bool
	}{
		{
			name: "successful guild leave",
			setupMocks: func(cr *MockCharacterRepository) {
				loginID := login.NewLoginID(uuid.New())
				character, _ := character.CreateNewCharacter("TestChar", loginID, class.Warrior, vault.NewVaultID(uuid.New()))
				cr.On("FindCharacterById", mock.Anything, characterID).Return(character, nil)
				cr.On("Update", mock.Anything, mock.AnythingOfType("Character")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "character not found",
			setupMocks: func(cr *MockCharacterRepository) {
				cr.On("FindCharacterById", mock.Anything, characterID).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCharacterRepository)
			tt.setupMocks(mockRepo)

			service := NewCharacterService(mockRepo, nil, &MockLogger{})
			err := service.LeaveGuild(context.Background(), characterID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestNotImplementedMethods(t *testing.T) {
	service := NewCharacterService(nil, nil, &MockLogger{})

	t.Run("PickItem", func(t *testing.T) {
		err := service.PickItem()
		assert.Error(t, err)
		assert.Equal(t, "method not implemented - requires character ID and item parameters", err.Error())
	})

	t.Run("DropItem", func(t *testing.T) {
		err := service.DropItem(context.Background(), playeritem.NewPlayerItemID(uuid.New()), 1)
		assert.Error(t, err)
		assert.Equal(t, "method not implemented - requires character ID parameter", err.Error())
	})
}
