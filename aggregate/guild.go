package aggregate

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/entity"
	"github.com/vterry/guild-project-ddd/valueobjects"
)

var (
	ErrInvalidGuildName           = errors.New("guild name must by between 4 and 15 characters")
	ErrPlayerIsAlreadyGuildMember = errors.New("the is already member of a Guild")
	ErrAlreadyInvited             = errors.New("this player has already been invited")
	ErrInvalidOperation           = errors.New("invalid operation")
	ErrNoInviteAvailable          = errors.New("there are no room for new invitations")
	ErrGuildAlreadyFull           = errors.New("guild is already full")
	MaxPlayers                    = 50
	MaxInvites                    = 50
)

type Guild struct {
	valueobjects.GuildID
	name      string
	createdAt time.Time
	manageBy  *entity.Player
	vault     *entity.Vault
	treasure  *entity.Treasure
	players   []*entity.Player
	invites   []*entity.Invite
}

// TODO Garantir na camada de aplicação que Players não podem criar guilds no nome de outros.
func CreateGuild(guildName string, guildOwner *entity.Player) (Guild, error) {

	if len(guildName) < 4 && len(guildName) > 15 {
		return Guild{}, ErrInvalidGuildName
	}

	if guildOwner.CurrentGuild != nil {
		return Guild{}, ErrPlayerIsAlreadyGuildMember
	}

	players := make([]*entity.Player, 0, MaxPlayers)
	players = append(players, guildOwner)

	guild := Guild{
		GuildID:   valueobjects.NewGuildID(uuid.New()),
		name:      guildName,
		createdAt: time.Now(),
		manageBy:  guildOwner,
		vault:     initVault(),
		treasure:  initTreasure(),
		players:   players,
		invites:   make([]*entity.Invite, 0, MaxInvites),
	}

	guildOwner.CurrentGuild = &guild.GuildID

	return guild, nil
}

func (g *Guild) InvitePlayer(sender *entity.Player, player *entity.Player) (entity.Invite, error) {

	if len(g.invites)+1 > cap(g.invites) && len(g.players)+1 > cap(g.players) {
		return entity.Invite{}, ErrNoInviteAvailable
	}

	// Garante que o Player convidado não pertence a nenhuma outra guild.
	// TODO Vale extrair para uma função separada?
	if player.CurrentGuild != nil {
		return entity.Invite{}, ErrPlayerIsAlreadyGuildMember
	}

	// Se o GM enviar o convite, automaticamente o convidado fará parte da guild.
	// TODO Vale extrair para uma função separada?
	if sender.CurrentGuild.ID() == g.manageBy.ID() {
		g.AddPlayer(player)
	}

	// Valida se o jogador já não possui um convite pendente
	for _, invite := range g.invites {
		if invite.GetPlayerID() == player.PlayerID {
			return entity.Invite{}, ErrAlreadyInvited
		}
	}

	invite := entity.NewInvite(player.PlayerID, sender.PlayerID, g.GuildID)
	return invite, nil
}

// TODO tratar concorrência e formas de tornar essa manipulação mais segura e performática - Garantir que somente Roles especificas adicionem pessoas na guild
// Garantir que players já membros da guild não poderão ser adicionados novamente
// Garantir que o player adicionado tera o campo CurrentGuild atualizado
// Estudar forma segura de fazer essa manipulação levando em consideração concorrência
func (g *Guild) AddPlayer(player *entity.Player) (*Guild, error) {

	if player.CurrentGuild.ID() != uuid.Nil {
		return nil, ErrPlayerIsAlreadyGuildMember
	}

	player.CurrentGuild = &g.GuildID
	g.players = append(g.players, player)

	return g, nil
}

func initVault() *entity.Vault {
	vault := entity.Vault{
		VaultID:    valueobjects.NewVaultID(uuid.New()),
		Items:      []*entity.Item{},
		GoldAmount: 0,
	}

	return &vault
}

func initTreasure() *entity.Treasure {
	vault := entity.Treasure{
		TreasureID: valueobjects.NewTreasureID(uuid.New()),
		CashAmount: 0,
	}

	return &vault
}
