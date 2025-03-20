package guild

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/common"
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

	params := specs.CreateGuildParams{
		GuildName:  guildName,
		GuildOwner: guildOwner,
	}

	spec := specs.NewCreateGuildSpecification()

	if err := spec(common.Base[specs.CreateGuildParams]{Entity: &params}); err != nil {
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

	params := specs.NewMemberParams{
		InviteSender:     sender,
		GuestPlayer:      guest,
		GuildMembers:     g.players,
		GuildInvitesSize: len(g.invites),
	}

	spec := specs.NewGuildMemberSpecification()

	if err := spec(common.Base[specs.NewMemberParams]{Entity: &params}); err != nil {
		return nil, err
	}

	if sender.Equals(g.manageBy.PlayerID) {
		_, err := g.addPlayerUnsafe(guest)

		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidOperation, err)
		}

		return nil, nil
	}

	if g.isInvitedPlayer(guest) {
		return nil, ErrAlreadyInvited
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
		return nil, ErrInviteNotExistis
	}

	err := g.invites[invite.ID()].reject()

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidOperation, err)
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
		return nil, ErrInviteNotExistis
	}

	err := g.invites[invite.ID()].cancel()

	if err != nil {
		return nil, ErrInvalidOperation
	}

	delete(g.invites, invite.ID())

	return invite, nil
}

func (g *Guild) ApproveInvite(invite *Invite) (*Invite, error) {
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

// TODO - Garantir que somente GM ou Commanders podem remover jogadores
func (g *Guild) AddPlayer(admin *player.Player, player *player.Player) (*Guild, error) {
	g.Lock()
	defer g.Unlock()

	params := specs.NewMemberParams{
		InviteSender:     admin,
		GuestPlayer:      player,
		GuildMembers:     g.players,
		GuildInvitesSize: len(g.invites),
	}

	spec := specs.NewGuildMemberSpecification()

	if err := spec(common.Base[specs.NewMemberParams]{Entity: &params}); err != nil {
		return nil, err
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
		return nil, ErrNotGuildMember
	}

	delete(g.players, player.ID())
	return g, nil
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

func (g *Guild) addPlayerUnsafe(player *player.Player) (*Guild, error) {
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
