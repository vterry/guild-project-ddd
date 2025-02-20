package aggregate

import (
	"errors"
	"fmt"
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

	if guildOwner.GetCurrentGuild() != uuid.Nil {
		return Guild{}, ErrPlayerIsAlreadyGuildMember
	}

	players := make([]*entity.Player, 0, MaxPlayers)
	players = append(players, guildOwner)

	guild := Guild{
		GuildID:   valueobjects.NewGuildID(uuid.New()),
		name:      guildName,
		createdAt: time.Now(),
		manageBy:  guildOwner,
		vault:     entity.NewVault(),
		treasure:  entity.NewTreasure(),
		players:   players,
		invites:   make([]*entity.Invite, 0, MaxInvites),
	}

	guildOwner.UpdateCurrentGuild(&guild.GuildID)

	return guild, nil
}

func (g *Guild) InvitePlayer(sender *entity.Player, player *entity.Player) (entity.Invite, error) {

	if len(g.invites)+1 > cap(g.invites) && len(g.players)+1 > cap(g.players) {
		return entity.Invite{}, ErrNoInviteAvailable
	}

	// Garante que o Player convidado não pertence a nenhuma outra guild.
	// TODO Vale extrair para uma função separada?
	if player.GetCurrentGuild() != uuid.Nil {
		return entity.Invite{}, ErrPlayerIsAlreadyGuildMember
	}

	// Se o GM enviar o convite, automaticamente o convidado fará parte da guild.
	// TODO Vale extrair para uma função separada?
	if sender.GetCurrentGuild() == g.manageBy.ID() {
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

	if player.GetCurrentGuild() != uuid.Nil {
		return nil, ErrPlayerIsAlreadyGuildMember
	}

	player.UpdateCurrentGuild(&g.GuildID)
	g.players = append(g.players, player)

	return g, nil
}

func (g *Guild) RemovePlayer(player *entity.Player) (*Guild, error) {
	if player.GetCurrentGuild() == g.ID() {
		return g, nil
	}
	return nil, ErrInvalidOperation
}

func (g *Guild) LeaveGuild(player *entity.Player) (*Guild, error) {
	if player.GetCurrentGuild() == g.ID() {
		player.UpdateCurrentGuild(nil)
		return g, nil
	}
	return nil, ErrInvalidOperation
}

// TODO - Melhorar -- better error handler
func (g Guild) AddItem(i *entity.Item) error {
	g.vault.Items = append(g.vault.Items, i)
	return nil
}

// TODO - Melhorar -- better error handler
func (g Guild) AddGold(gold int) error {
	g.vault.GoldAmount += gold
	return nil
}

func (g Guild) Print() {
	fmt.Printf("Guild ID: %v \n Guild Name: %v \n Guild Master: %v \n Guild Vault ID: %v \n Members Size: %d \n Items Stored: %d \n", g.GuildID.ID(), g.name, g.manageBy, g.vault.VaultID.ID(), len(g.players), len(g.vault.Items))
}
