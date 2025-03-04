package guild

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/guild/valueobjects"
	"github.com/vterry/guild-project-ddd/domain/player"
	"github.com/vterry/guild-project-ddd/domain/treasure"
	"github.com/vterry/guild-project-ddd/domain/vault"
)

const (
	MAX_PLAYERS = 50
	MAX_INVITES = 50
)

type Guild struct {
	valueobjects.GuildID
	name      string
	createdAt time.Time
	manageBy  *player.Player
	vault     *vault.Vault
	treasure  *treasure.Treasure
	players   map[uuid.UUID]*player.Player
	invites   map[uuid.UUID]*Invite
	sync.Mutex
}

func CreateGuild(guildName string, guildOwner *player.Player) (*Guild, error) {

	if len(guildName) < 4 || len(guildName) > 15 {
		return nil, NewGuildError(ErrInvalidGuildName, nil)
	}

	if guildOwner == nil {
		return nil, NewGuildError(ErrMustInformGuidOwner, nil)
	}

	if guildOwner.GetCurrentGuild() != uuid.Nil {
		return nil, NewGuildError(ErrAnotherGuildMember, nil)
	}

	players := make(map[uuid.UUID]*player.Player, MAX_PLAYERS)
	players[guildOwner.ID()] = guildOwner

	guild := initializeGuild(guildName, guildOwner, players)
	guildOwner.UpdateCurrentGuild(guild.GuildID.ID())

	return guild, nil
}

func (g *Guild) InvitePlayer(sender *player.Player, guest *player.Player) (*Invite, error) {

	g.Lock()
	defer g.Unlock()

	if len(g.invites)+1 >= MAX_INVITES || len(g.players)+1 > MAX_PLAYERS {
		return nil, NewGuildError(ErrNoInviteAvailable, nil)
	}

	//Check if sender is a guild member
	if _, isMember := g.players[sender.ID()]; !isMember {
		return nil, NewGuildError(ErrCannotInvite, nil)
	}

	//Check if guest is a guild member
	if _, isMember := g.players[guest.ID()]; isMember {
		return nil, NewGuildError(ErrPlayerIsAlreadyGuildMember, nil)
	}

	// Garante que o Player convidado não pertence a nenhuma outra guild.
	if guest.GetCurrentGuild() != uuid.Nil {
		return nil, NewGuildError(ErrAnotherGuildMember, nil)
	}

	if sender.Equals(g.manageBy.PlayerID) {
		_, err := g.addPlayerUnsafe(guest)

		if err != nil {
			return nil, NewGuildError(ErrInvalidOperation, err)
		}

		return nil, nil
	}

	if g.isInvitedPlayer(guest) {
		return nil, NewGuildError(ErrAlreadyInvited, nil)
	}

	invite := NewInvite(guest.ID(), sender.ID(), g.ID())
	g.invites[invite.ID()] = invite

	return invite, nil
}

// TODO - Jogadores convidados podem recusar o convite
func (g *Guild) RejectInvite(invite *Invite) (*Invite, error) {
	g.Lock()
	defer g.Unlock()

	//Check is valid invite
	if _, isValid := g.invites[invite.ID()]; !isValid {
		return nil, NewGuildError(ErrInviteNotExistis, nil)
	}

	err := g.invites[invite.ID()].reject()

	if err != nil {
		return nil, NewGuildError(ErrInvalidOperation, err)
	}

	delete(g.invites, invite.ID())
	return invite, nil
}

// TODO - Garantir que somente GM ou Commanders podem cancelar convites
func (g *Guild) CancelInvite(invite *Invite) (*Invite, error) {
	g.Lock()
	defer g.Unlock()

	//Check is valid invite
	if _, isValid := g.invites[invite.ID()]; !isValid {
		return nil, NewGuildError(ErrInviteNotExistis, nil)
	}

	err := g.invites[invite.ID()].cancel()

	if err != nil {
		return nil, NewGuildError(ErrInvalidOperation, err)
	}

	delete(g.invites, invite.ID())

	return invite, nil
}

func (g *Guild) ApproveInvite(invite *Invite) (*Invite, error) {
	g.Lock()
	defer g.Unlock()

	//Check is valid invite
	if _, isValid := g.invites[invite.ID()]; !isValid {
		return nil, NewGuildError(ErrInviteNotExistis, nil)
	}

	err := g.invites[invite.ID()].approve()

	if err != nil {
		return nil, NewGuildError(ErrInvalidOperation, err)
	}

	delete(g.invites, invite.ID())

	return invite, nil
}

// TODO - Garantir que somente GM ou Commanders podem remover jogadores
func (g *Guild) AddPlayer(admin *player.Player, player *player.Player) (*Guild, error) {
	g.Lock()
	defer g.Unlock()

	if len(g.players) >= MAX_PLAYERS {
		return nil, NewGuildError(ErrGuildAlreadyFull, nil)
	}

	if _, isMember := g.players[player.ID()]; isMember {
		return nil, NewGuildError(ErrPlayerIsAlreadyGuildMember, nil)
	}

	if player.GetCurrentGuild() != uuid.Nil {
		return nil, NewGuildError(ErrAnotherGuildMember, nil)
	}

	player.UpdateCurrentGuild(g.ID())
	g.players[player.ID()] = player

	return g, nil
}

// TODO - Garantir que somente GM ou Commanders podem remover jogadores
func (g *Guild) RemovePlayer(admin *player.Player, player *player.Player) (*Guild, error) {
	g.Lock()
	defer g.Unlock()

	if _, isMember := g.players[player.ID()]; !isMember {
		return nil, NewGuildError(ErrNotGuildMember, nil)
	}

	delete(g.players, player.ID())
	return g, nil
}

func (g *Guild) LeaveGuild(player *player.Player) (*Guild, error) {
	g.Lock()
	defer g.Unlock()

	if _, isMember := g.players[player.ID()]; !isMember {
		return nil, NewGuildError(ErrNotGuildMember, nil)
	}

	delete(g.players, player.ID())
	return g, nil
}

func (g *Guild) addPlayerUnsafe(player *player.Player) (*Guild, error) {

	if len(g.players) >= MAX_PLAYERS {
		return nil, NewGuildError(ErrGuildAlreadyFull, nil)
	}

	if _, isMember := g.players[player.ID()]; isMember {
		return nil, NewGuildError(ErrPlayerIsAlreadyGuildMember, nil)
	}

	player.UpdateCurrentGuild(g.ID())
	g.players[player.ID()] = player

	return g, nil
}

func (g *Guild) isInvitedPlayer(player *player.Player) bool {
	for _, invite := range g.invites {
		if invite.GetPlayerID() == player.ID() {
			return true
		}
	}
	return false
}

func initializeGuild(guildName string, guildOwner *player.Player, players map[uuid.UUID]*player.Player) *Guild {
	return &Guild{
		GuildID:   valueobjects.NewGuildID(uuid.New()),
		name:      guildName,
		createdAt: time.Now(),
		manageBy:  guildOwner,
		vault:     vault.NewVault(),
		treasure:  treasure.NewTreasure(),
		players:   players,
		invites:   make(map[uuid.UUID]*Invite, MAX_INVITES),
	}
}
