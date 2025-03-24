package guild

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/guild/specs"
	"github.com/vterry/guild-project-ddd/domain/guild/valueobjects"
	"github.com/vterry/guild-project-ddd/domain/player"
	"github.com/vterry/guild-project-ddd/domain/treasure"
	"github.com/vterry/guild-project-ddd/domain/vault"
)

const (
	MAX_PLAYERS = 50
	MAX_INVITES = 50
)

var (
	ErrNotGuildMember   = errors.New("player isnt a guild member")
	ErrAlreadyInvited   = errors.New("this player has already a pending invite")
	ErrInviteNotExistis = errors.New("invite no longer exists")
	ErrInvalidOperation = errors.New("invalid operation")
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

	if err := specs.ValidateGuildCreation(guildName, guildOwner); err != nil {
		return nil, err
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

	if err := specs.ValidateNewMember(sender, guest, g.players, len(g.invites)); err != nil {
		return nil, err
	}

	if g.isInvitedPlayer(guest) {
		return nil, ErrAlreadyInvited
	}

	invite := NewInvite(guest.ID(), sender.ID(), g.ID())
	g.invites[invite.ID()] = invite

	return invite, nil
}

func (g *Guild) LeaveGuild(player *player.Player) (*Guild, error) {
	g.Lock()
	defer g.Unlock()

	if _, isMember := g.players[player.ID()]; !isMember {
		return nil, ErrNotGuildMember
	}

	delete(g.players, player.ID())
	return g, nil
}

// TODO - Secure this using RBAC
// The invitee can reject the invitation, but It don have access to this, so
// I thinking to consider move this logic to a Domain Service where Player subscribe
// and sent te reject operation by sending another event.

func (g *Guild) rejectInvite(invite *Invite) (*Invite, error) {
	g.Lock()
	defer g.Unlock()

	if _, isValid := g.invites[invite.ID()]; !isValid {
		return nil, ErrInviteNotExistis
	}

	err := g.invites[invite.ID()].reject()

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidOperation, err)
	}

	delete(g.invites, invite.ID())
	return invite, nil
}

// TODO - Secure this using RBAC
// Managers can cancel a Invite
//Who sent the invite can cancel it

func (g *Guild) cancelInvite(invite *Invite) (*Invite, error) {
	g.Lock()
	defer g.Unlock()

	//Check is valid invite
	if _, isValid := g.invites[invite.ID()]; !isValid {
		return nil, ErrInviteNotExistis
	}

	err := g.invites[invite.ID()].cancel()

	if err != nil {
		return nil, ErrInvalidOperation
	}

	delete(g.invites, invite.ID())

	return invite, nil
}

// TODO - Secure this using RBAC
// onlu managers can approve

func (g *Guild) approveInvite(invite *Invite) (*Invite, error) {
	g.Lock()
	defer g.Unlock()

	//Check is valid invite
	if _, isValid := g.invites[invite.ID()]; !isValid {
		return nil, ErrInviteNotExistis
	}

	err := g.invites[invite.ID()].approve()

	if err != nil {
		return nil, ErrInvalidOperation
	}

	delete(g.invites, invite.ID())

	return invite, nil
}

// TODO - Secure this using RBAC
// Only managers can add players

func (g *Guild) addPlayer(admin *player.Player, player *player.Player) (*Guild, error) {
	g.Lock()
	defer g.Unlock()

	if err := specs.ValidateNewMember(admin, player, g.players, len(g.invites)); err != nil {
		return nil, err
	}

	player.UpdateCurrentGuild(g.ID())
	g.players[player.ID()] = player

	return g, nil
}

// TODO - Secure this using RBAC
// Only members can remove players

func (g *Guild) removePlayer(admin *player.Player, player *player.Player) (*Guild, error) {
	g.Lock()
	defer g.Unlock()

	if _, isMember := g.players[player.ID()]; !isMember {
		return nil, ErrNotGuildMember
	}

	delete(g.players, player.ID())
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
	createAt := time.Now()
	return &Guild{
		GuildID:   valueobjects.NewGuildID(guildName, createAt),
		name:      guildName,
		createdAt: createAt,
		manageBy:  guildOwner,
		vault:     vault.NewVault(guildName),
		treasure:  treasure.NewTreasure(guildName),
		players:   players,
		invites:   make(map[uuid.UUID]*Invite, MAX_INVITES),
	}
}

func (g *Guild) GetManager() *player.Player {
	return g.manageBy
}
