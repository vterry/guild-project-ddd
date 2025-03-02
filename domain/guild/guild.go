package guild

import (
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/guild/valueobjects"
	"github.com/vterry/guild-project-ddd/domain/player"
	"github.com/vterry/guild-project-ddd/domain/treasure"
	"github.com/vterry/guild-project-ddd/domain/vault"
)

var (
	ErrInvalidGuildName           = errors.New("guild name must by between 4 and 15 characters")
	ErrPlayerIsAlreadyGuildMember = errors.New("the is already member of a Guild")
	ErrAlreadyInvited             = errors.New("this player has already been invited")
	ErrInvalidOperation           = errors.New("invalid operation")
	ErrNoInviteAvailable          = errors.New("there are no room for new invitations")
	ErrGuildAlreadyFull           = errors.New("guild is already full")
	MAX_PLAYERS                   = 50
	MAX_INVITES                   = 50
)

type Guild struct {
	valueobjects.GuildID
	name      string
	createdAt time.Time
	manageBy  *player.Player
	vault     *vault.Vault
	treasure  *treasure.Treasure
	players   []*player.Player
	invites   []*Invite
	sync.Mutex
}

func CreateGuild(guildName string, guildOwner *player.Player) (*Guild, error) {

	if len(guildName) < 4 || len(guildName) > 15 {
		return nil, ErrInvalidGuildName
	}

	if guildOwner.GetCurrentGuild() != uuid.Nil {
		return nil, ErrPlayerIsAlreadyGuildMember
	}

	players := make([]*player.Player, 0, MAX_PLAYERS)
	players = append(players, guildOwner)

	guild := initializeGuild(guildName, guildOwner, players)
	guildOwner.UpdateCurrentGuild(guild.GuildID.ID())

	return guild, nil
}

func (g *Guild) InvitePlayer(sender *player.Player, player *player.Player) (*Invite, error) {

	g.Lock()
	defer g.Unlock()

	if len(g.invites)+1 > cap(g.invites) || len(g.players)+1 > cap(g.players) {
		return nil, ErrNoInviteAvailable
	}

	// Garante que o Player convidado não pertence a nenhuma outra guild.
	// TODO Vale extrair para uma função separada?
	if player.GetCurrentGuild() != uuid.Nil {
		return nil, ErrPlayerIsAlreadyGuildMember
	}

	// Se o GM enviar o convite, automaticamente o convidado fará parte da guild.
	// TODO Vale extrair para uma função separada?
	if sender.GetCurrentGuild() == g.manageBy.ID() {
		g.AddPlayer(player)
	}

	if g.isInvitedPlayer(player) {
		return nil, ErrAlreadyInvited
	}

	invite := NewInvite(player.PlayerID, sender.PlayerID, g.GuildID)
	g.invites = append(g.invites, invite)

	return invite, nil
}

// TODO - Jogadores convidados podem recusar o convite
func (g *Guild) RejectInvite(invite *Invite) (*Invite, error) {
	g.Lock()
	defer g.Unlock()

	inviteIndex := slices.IndexFunc(g.invites, func(i *Invite) bool {
		return i.InviteID.Equals(invite.InviteID)
	})

	if inviteIndex == -1 {
		return nil, ErrInvalidOperation
	}

	g.invites[inviteIndex].reject()

	g.invites = append(g.invites[:inviteIndex], g.invites[inviteIndex+1:]...)

	return invite, nil
}

// TODO - Garantir que somente GM ou Commanders podem cancelar convites
func (g *Guild) CancelInvite(invite *Invite) (*Invite, error) {
	g.Lock()
	defer g.Unlock()

	inviteIndex := slices.IndexFunc(g.invites, func(i *Invite) bool {
		return i.InviteID.Equals(invite.InviteID)
	})

	if inviteIndex == -1 {
		return nil, ErrInvalidOperation
	}

	g.invites[inviteIndex].cancel()

	g.invites = append(g.invites[:inviteIndex], g.invites[inviteIndex+1:]...)

	return invite, nil
}

func (g *Guild) ApproveInvite(invite *Invite) (*Invite, error) {
	g.Lock()
	defer g.Unlock()

	inviteIndex := slices.IndexFunc(g.invites, func(i *Invite) bool {
		return i.InviteID.Equals(invite.InviteID)
	})

	if inviteIndex == -1 {
		return nil, ErrInvalidOperation
	}

	g.invites[inviteIndex].approve()

	g.invites = append(g.invites[:inviteIndex], g.invites[inviteIndex+1:]...)

	return invite, nil
}

// TODO tratar concorrência e formas de tornar essa manipulação mais segura e performática - Garantir que somente Roles especificas adicionem pessoas na guild
// Garantir que players já membros da guild não poderão ser adicionados novamente
// Garantir que o player adicionado tera o campo CurrentGuild atualizado
// Estudar forma segura de fazer essa manipulação levando em consideração concorrência
func (g *Guild) AddPlayer(player *player.Player) (*Guild, error) {
	g.Lock()
	defer g.Unlock()

	if player.GetCurrentGuild() != uuid.Nil {
		return nil, ErrPlayerIsAlreadyGuildMember
	}

	player.UpdateCurrentGuild(g.GuildID.ID())
	g.players = append(g.players, player)

	return g, nil
}

// TODO - Garantir que somente GM ou Commanders podem remover jogadores
func (g *Guild) RemovePlayer(player *player.Player) (*Guild, error) {
	g.Lock()
	defer g.Unlock()

	if player.GetCurrentGuild() == g.ID() {
		return g, ErrInvalidOperation
	}

	playerIndex := slices.Index(g.players, player)
	if playerIndex == -1 {
		return nil, ErrInvalidOperation
	}

	g.players = append(g.players[:playerIndex], g.players[playerIndex+1:]...)
	return nil, ErrInvalidOperation
}

func (g *Guild) LeaveGuild(player *player.Player) (*Guild, error) {
	g.Lock()
	defer g.Unlock()

	if player.GetCurrentGuild() == g.ID() {
		return g, ErrInvalidOperation
	}

	playerIndex := slices.Index(g.players, player)
	if playerIndex == -1 {
		return nil, ErrInvalidOperation
	}

	g.players = append(g.players[:playerIndex], g.players[playerIndex+1:]...)
	return nil, ErrInvalidOperation
}

func (g *Guild) isInvitedPlayer(player *player.Player) bool {
	for _, invite := range g.invites {
		if invite.GetPlayerID() == player.PlayerID {
			return true
		}
	}
	return false
}

func initializeGuild(guildName string, guildOwner *player.Player, players []*player.Player) *Guild {
	return &Guild{
		GuildID:   valueobjects.NewGuildID(uuid.New()),
		name:      guildName,
		createdAt: time.Now(),
		manageBy:  guildOwner,
		vault:     vault.NewVault(),
		treasure:  treasure.NewTreasure(),
		players:   players,
		invites:   make([]*Invite, 0, MAX_INVITES),
	}
}
