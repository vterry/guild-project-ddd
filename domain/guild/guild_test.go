package guild

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vterry/guild-project-ddd/domain/common"
	"github.com/vterry/guild-project-ddd/domain/guild/specs"
	"github.com/vterry/guild-project-ddd/domain/player"
	"github.com/vterry/guild-project-ddd/domain/player/valueobjects"
)

var (
	defaultGuildName = "GuildDefault"
	guildOwner, _    = player.NewPlayer("Guild Owner", valueobjects.Mage)
	guildInstance    *Guild
	guildOnce        sync.Once
)

func TestInvalidGuildInitialization(t *testing.T) {
	t.Cleanup(resetGuild)
	t.Run("test create new guild with empty name", func(t *testing.T) {
		_, err := CreateGuild("", guildOwner)
		assert.ErrorIs(t, err, specs.ErrInvalidGuildName)
	})

	t.Run("test create new guild large name", func(t *testing.T) {
		_, err := CreateGuild("GuildGuildGuildGuild", guildOwner)
		assert.ErrorIs(t, err, specs.ErrInvalidGuildName)
	})

	t.Run("test create new guild with special character (space)", func(t *testing.T) {
		_, err := CreateGuild("Guild Special", guildOwner)
		assert.ErrorIs(t, err, specs.ErrInvalidCharName)
	})

	t.Run("test create new guild with special character", func(t *testing.T) {
		_, err := CreateGuild("Guild@Spec14l", guildOwner)
		assert.ErrorIs(t, err, specs.ErrInvalidCharName)
	})

	t.Run("test create new guild with no guild owner", func(t *testing.T) {
		_, err := CreateGuild(defaultGuildName, nil)
		assert.ErrorIs(t, err, specs.ErrMustInformGuidOwner)
	})

	t.Run("test create new build being member of another", func(t *testing.T) {
		mockPlayer, _ := player.NewPlayer("Mock Player", valueobjects.Mage)
		mockPlayer.UpdateCurrentGuild("SampleGuild-2025-03-16-203640-n3Gkxo6R6")
		_, err := CreateGuild("ErrorGuild", mockPlayer)
		assert.ErrorIs(t, err, specs.ErrAnotherGuildMember)
	})
}

// Create a guild with correct name

func TestValidGuildInitialization(t *testing.T) {
	t.Cleanup(resetGuild)
	guild := getGuild()

	t.Run("Check guild was created", func(t *testing.T) {
		assert.NotNil(t, guild)
	})

	t.Run("Check guild ID is not nil", func(t *testing.T) {
		assert.NotNil(t, guild.ID())
	})

	t.Run("Check guild name", func(t *testing.T) {
		assert.Equal(t, defaultGuildName, guild.name)
	})

	t.Run("Check guild owner", func(t *testing.T) {
		assert.Equal(t, guildOwner.GetCurrentGuild(), guild.GuildID.ID())
	})

	t.Run("Check if guild has a vault", func(t *testing.T) {
		assert.NotNil(t, guild.vault)
	})

	t.Run("Check if guild has a treasure", func(t *testing.T) {
		assert.NotNil(t, guild.treasure)
	})

	t.Run("Check if guild has only 1 member", func(t *testing.T) {
		assert.Equal(t, 1, len(guild.players))
	})

	t.Run("Check if guild has no listed invite", func(t *testing.T) {
		assert.Equal(t, 0, len(guild.invites))
	})
}

func TestGuildInvitationFuncions(t *testing.T) {

	t.Run("Test GM inviting a player", func(t *testing.T) {
		t.Cleanup(resetGuild)
		guild := getGuild()

		guest, _ := player.NewPlayer("Guest", valueobjects.Ranger)
		invite, err := guild.InvitePlayer(guildOwner, guest)

		assert.Nil(t, invite)
		assert.Nil(t, err)

		_, isMember := guild.players[guest.ID()]
		assert.True(t, isMember)
		assert.True(t, guest.GetCurrentGuild() == guild.ID())
	})

	t.Run("Test guild's member inviting player", func(t *testing.T) {

		t.Cleanup(resetGuild)
		guild := getGuild()

		sender, _ := player.NewPlayer("Sender", valueobjects.Warrior)
		guest, _ := player.NewPlayer("Guest", valueobjects.Ranger)
		guild.InvitePlayer(guildOwner, sender)

		invite, _ := guild.InvitePlayer(sender, guest)
		_, test := guild.invites[invite.ID()]
		lenInvites := len(guild.invites)

		assert.True(t, test)

		_, err := guild.InvitePlayer(sender, guest)

		assert.ErrorIs(t, err, ErrAlreadyInvited)
		assert.Equal(t, lenInvites, len(guild.invites))

	})

	t.Run("Test where sender is not a guild's member", func(t *testing.T) {

		t.Cleanup(resetGuild)
		guild := getGuild()

		sender, _ := player.NewPlayer("Sender", valueobjects.Warrior)
		guest, _ := player.NewPlayer("Guest", valueobjects.Ranger)

		_, err := guild.InvitePlayer(sender, guest)

		assert.ErrorIs(t, err, specs.ErrCannotInvite)

	})

	t.Run("Test where guest is already guild's member", func(t *testing.T) {

		t.Cleanup(resetGuild)
		guild := getGuild()

		sender, _ := player.NewPlayer("Sender", valueobjects.Warrior)
		guest, _ := player.NewPlayer("Guest", valueobjects.Ranger)
		guild.InvitePlayer(guildOwner, sender)
		guild.InvitePlayer(guildOwner, guest)

		_, err := guild.InvitePlayer(sender, guest)

		assert.ErrorIs(t, err, specs.ErrPlayerIsAlreadyGuildMember)

	})

	t.Run("Test where guest is already member of another guild", func(t *testing.T) {

		t.Cleanup(resetGuild)
		guild := getGuild()

		sender, _ := player.NewPlayer("Sender", valueobjects.Warrior)
		guest, _ := player.NewPlayer("Guest", valueobjects.Ranger)
		guest.UpdateCurrentGuild("Sample-Guild-2025-03-16-203640-n3Gkxo6R6")
		guild.InvitePlayer(guildOwner, sender)

		_, err := guild.InvitePlayer(sender, guest)

		assert.ErrorIs(t, err, specs.ErrAnotherGuildMember)

	})

	t.Run("Test invite player with no room for new invites", func(t *testing.T) {

		t.Cleanup(resetGuild)
		guild := getGuild()

		sender, _ := player.NewPlayer("Sender", valueobjects.Warrior)
		guest, _ := player.NewPlayer("Guest", valueobjects.Ranger)
		guild.InvitePlayer(guildOwner, sender)

		fullfilGuildInviteList(guild, sender)

		_, err := guild.InvitePlayer(sender, guest)

		assert.ErrorIs(t, err, specs.ErrNoInviteAvailable)

	})

	t.Run("Test invite player with no room for new players", func(t *testing.T) {
		t.Cleanup(resetGuild)
		guild := getGuild()

		sender, _ := player.NewPlayer("Sender", valueobjects.Warrior)
		guest, _ := player.NewPlayer("Guest", valueobjects.Ranger)
		guild.InvitePlayer(guildOwner, sender)

		fullfilGuildCapacity(guild)

		_, err := guild.InvitePlayer(sender, guest)

		assert.ErrorIs(t, err, specs.ErrGuildAlreadyFull)

	})

}

func TestGuildInviteRejectFunctions(t *testing.T) {

	t.Run("Test reject a invite that not exists", func(t *testing.T) {

		t.Cleanup(resetGuild)
		guild := getGuild()

		sender, _ := player.NewPlayer("Sender", valueobjects.Warrior)
		guest, _ := player.NewPlayer("Guest", valueobjects.Ranger)
		mockInvite := NewInvite(guest.ID(), sender.ID(), guild.ID())

		_, err := guild.RejectInvite(mockInvite)
		assert.ErrorIs(t, err, ErrInviteNotExistis)

	})

	t.Run("Test reject a invite with an invalid state (!= Pending)", func(t *testing.T) {

		t.Cleanup(resetGuild)
		guild := getGuild()

		sender, _ := player.NewPlayer("Sender", valueobjects.Warrior)
		guest, _ := player.NewPlayer("Guest", valueobjects.Ranger)
		mockInvite := NewInvite(guest.ID(), sender.ID(), guild.ID())

		guild.invites[mockInvite.ID()] = mockInvite
		mockInvite.approve()

		_, err := guild.RejectInvite(mockInvite)
		assert.ErrorIs(t, err, ErrInvalidOperation, ErrInvalidInviteState)

	})

	t.Run("Test reject a invite in a valid state (== Pending)", func(t *testing.T) {
		t.Cleanup(resetGuild)
		guild := getGuild()

		sender, _ := player.NewPlayer("Sender", valueobjects.Warrior)
		guild.InvitePlayer(guildOwner, sender)

		guest, _ := player.NewPlayer("Guest", valueobjects.Ranger)
		invite, err := guild.InvitePlayer(sender, guest)

		assert.Nil(t, err)
		_, wasCreated := guild.invites[invite.ID()]
		assert.True(t, wasCreated)

		lenInv := len(guild.invites)

		_, rejErr := guild.RejectInvite(invite)
		assert.Nil(t, rejErr)
		assert.NotEqual(t, lenInv, len(guild.invites))
		assert.Equal(t, InviteStatus(common.Rejected), invite.CheckStatus())

	})
}

func TestGuildInviteCancelFunctions(t *testing.T) {

	t.Run("Test cancel a invite that not exists", func(t *testing.T) {

		t.Cleanup(resetGuild)
		guild := getGuild()

		sender, _ := player.NewPlayer("Sender", valueobjects.Warrior)
		guest, _ := player.NewPlayer("Guest", valueobjects.Ranger)
		mockInvite := NewInvite(guest.ID(), sender.ID(), guild.ID())

		_, err := guild.CancelInvite(mockInvite)
		assert.ErrorIs(t, err, ErrInviteNotExistis)

	})

	t.Run("Test cancel a invite with an invalid state (!= Pending)", func(t *testing.T) {
		t.Cleanup(resetGuild)
		guild := getGuild()

		sender, _ := player.NewPlayer("Sender", valueobjects.Warrior)
		guest, _ := player.NewPlayer("Guest", valueobjects.Ranger)
		mockInvite := NewInvite(guest.ID(), sender.ID(), guild.ID())

		guild.invites[mockInvite.ID()] = mockInvite
		mockInvite.approve()

		_, err := guild.CancelInvite(mockInvite)
		assert.ErrorIs(t, err, ErrInvalidOperation, ErrInvalidInviteState)

	})

	t.Run("Test cancel a invite in a valid state (== Pending)", func(t *testing.T) {
		t.Cleanup(resetGuild)
		guild := getGuild()

		sender, _ := player.NewPlayer("Sender", valueobjects.Warrior)
		guild.InvitePlayer(guildOwner, sender)

		guest, _ := player.NewPlayer("Guest", valueobjects.Ranger)
		invite, err := guild.InvitePlayer(sender, guest)

		assert.Nil(t, err)
		_, wasCreated := guild.invites[invite.ID()]
		assert.True(t, wasCreated)

		lenInv := len(guild.invites)

		_, rejErr := guild.CancelInvite(invite)
		assert.Nil(t, rejErr)
		assert.NotEqual(t, lenInv, len(guild.invites))
		assert.Equal(t, InviteStatus(common.Canceled), invite.CheckStatus())

	})
}

func TestGuildInviteApproveFunctions(t *testing.T) {

	t.Run("Test approve a invite that not exists", func(t *testing.T) {
		t.Cleanup(resetGuild)
		guild := getGuild()

		sender, _ := player.NewPlayer("Sender", valueobjects.Warrior)
		guest, _ := player.NewPlayer("Guest", valueobjects.Ranger)
		mockInvite := NewInvite(guest.ID(), sender.ID(), guild.ID())

		_, err := guild.ApproveInvite(mockInvite)
		assert.ErrorIs(t, err, ErrInviteNotExistis)

	})

	t.Run("Test approve a invite with an invalid state (!= Pending)", func(t *testing.T) {
		t.Cleanup(resetGuild)
		guild := getGuild()

		sender, _ := player.NewPlayer("Sender", valueobjects.Warrior)
		guest, _ := player.NewPlayer("Guest", valueobjects.Ranger)
		mockInvite := NewInvite(guest.ID(), sender.ID(), guild.ID())

		guild.invites[mockInvite.ID()] = mockInvite
		mockInvite.reject()

		_, err := guild.ApproveInvite(mockInvite)
		assert.ErrorIs(t, err, ErrInvalidOperation, ErrInvalidInviteState)

	})

	t.Run("Test approve a invite in a valid state (== Pending)", func(t *testing.T) {
		t.Cleanup(resetGuild)
		guild := getGuild()

		sender, _ := player.NewPlayer("Sender", valueobjects.Warrior)
		guild.InvitePlayer(guildOwner, sender)

		guest, _ := player.NewPlayer("Guest", valueobjects.Ranger)
		invite, err := guild.InvitePlayer(sender, guest)

		assert.Nil(t, err)
		_, wasCreated := guild.invites[invite.ID()]
		assert.True(t, wasCreated)

		lenInv := len(guild.invites)

		_, rejErr := guild.ApproveInvite(invite)
		assert.Nil(t, rejErr)
		assert.NotEqual(t, lenInv, len(guild.invites))
		assert.Equal(t, InviteStatus(common.Approved), invite.CheckStatus())

	})
}

func TestAddToGuildFunctions(t *testing.T) {
	t.Run("Test adding player of another guild", func(t *testing.T) {
		t.Cleanup(resetGuild)
		guild := getGuild()

		mockPlayer, _ := player.NewPlayer("MockPlayer", valueobjects.Mage)
		mockPlayer.UpdateCurrentGuild("Sample-Guild-2025-03-16-203640-n3Gkxo6R6")

		_, err := guild.AddPlayer(guildOwner, mockPlayer)
		assert.ErrorIs(t, err, specs.ErrAnotherGuildMember)
	})

	t.Run("Test adding player who is already a guild's member", func(t *testing.T) {
		t.Cleanup(resetGuild)
		guild := getGuild()

		guest, _ := player.NewPlayer("MockPlayer", valueobjects.Mage)
		guild.InvitePlayer(guildOwner, guest)

		_, err := guild.AddPlayer(guildOwner, guest)
		assert.ErrorIs(t, err, specs.ErrPlayerIsAlreadyGuildMember)
	})

	t.Run("Test adding player with no room for new players", func(t *testing.T) {
		t.Cleanup(resetGuild)
		guild := getGuild()

		guest, _ := player.NewPlayer("MockPlayer", valueobjects.Mage)

		fullfilGuildCapacity(guild)

		_, err := guild.AddPlayer(guildOwner, guest)
		assert.ErrorIs(t, err, specs.ErrGuildAlreadyFull)
	})

	t.Run("Test adding new player", func(t *testing.T) {
		t.Cleanup(resetGuild)
		guild := getGuild()

		guest, _ := player.NewPlayer("MockPlayer", valueobjects.Mage)

		_, err := guild.AddPlayer(guildOwner, guest)
		assert.Nil(t, err)

		_, isMember := guild.players[guest.ID()]
		assert.True(t, isMember)
	})
}

func TestRemoveFromGuildFunctions(t *testing.T) {
	t.Run("Test remove player who isn't a guild's member", func(t *testing.T) {
		t.Cleanup(resetGuild)
		guild := getGuild()

		mockPlayer, _ := player.NewPlayer("MockPlayer", valueobjects.Mage)
		mockPlayer.UpdateCurrentGuild("Sample-Guild-2025-03-16-203640-n3Gkxo6R6")

		_, err := guild.RemovePlayer(guildOwner, mockPlayer)
		assert.ErrorIs(t, err, ErrNotGuildMember)
	})

	t.Run("Test remove a guild's member", func(t *testing.T) {
		t.Cleanup(resetGuild)
		guild := getGuild()

		guest, _ := player.NewPlayer("MockPlayer", valueobjects.Mage)
		guild.AddPlayer(guildOwner, guest)
		lenPlayer := len(guild.players)

		_, err := guild.RemovePlayer(guildOwner, guest)
		assert.Nil(t, err)

		_, isMember := guild.players[guest.ID()]
		assert.False(t, isMember)
		assert.Equal(t, lenPlayer-1, len(guild.players))
	})
}

func TestLeaveFromGuildFunctions(t *testing.T) {
	t.Run("Test leave not being a guild's member", func(t *testing.T) {
		t.Cleanup(resetGuild)
		guild := getGuild()

		mockPlayer, _ := player.NewPlayer("MockPlayer", valueobjects.Mage)
		mockPlayer.UpdateCurrentGuild("Sample-Guild-2025-03-16-203640-n3Gkxo6R6")

		_, err := guild.LeaveGuild(mockPlayer)
		assert.ErrorIs(t, err, ErrNotGuildMember)
	})

	t.Run("Test leave guild", func(t *testing.T) {
		t.Cleanup(resetGuild)
		guild := getGuild()

		guest, _ := player.NewPlayer("MockPlayer", valueobjects.Mage)
		guild.AddPlayer(guildOwner, guest)
		lenPlayer := len(guild.players)

		_, err := guild.LeaveGuild(guest)
		assert.Nil(t, err)

		_, isMember := guild.players[guest.ID()]
		assert.False(t, isMember)
		assert.Equal(t, lenPlayer-1, len(guild.players))
	})
}

func TestLeaveGuildFunctions(t *testing.T) {
	t.Run("Test remove player who isn't a guild's member", func(t *testing.T) {
		t.Cleanup(resetGuild)
		guild := getGuild()

		mockPlayer, _ := player.NewPlayer("MockPlayer", valueobjects.Mage)
		mockPlayer.UpdateCurrentGuild("Sample-Guild-2025-03-16-203640-n3Gkxo6R6")

		_, err := guild.RemovePlayer(guildOwner, mockPlayer)
		assert.ErrorIs(t, err, ErrNotGuildMember)
	})

	t.Run("Test remove a guild's member", func(t *testing.T) {
		t.Cleanup(resetGuild)
		guild := getGuild()

		guest, _ := player.NewPlayer("MockPlayer", valueobjects.Mage)
		guild.AddPlayer(guildOwner, guest)
		lenPlayer := len(guild.players)

		_, err := guild.RemovePlayer(guildOwner, guest)
		assert.Nil(t, err)

		_, isMember := guild.players[guest.ID()]
		assert.False(t, isMember)
		assert.Equal(t, lenPlayer-1, len(guild.players))
	})
}

func fullfilGuildInviteList(g *Guild, regPlayer *player.Player) {
	for i := 0; i < MAX_INVITES; i++ {
		name := fmt.Sprintf("%s%d", "Player", i)
		player, _ := player.NewPlayer(name, valueobjects.Mage)
		g.InvitePlayer(regPlayer, player)
	}
}

func fullfilGuildCapacity(g *Guild) {
	for i := 0; i < MAX_PLAYERS; i++ {
		name := fmt.Sprintf("%s%d", "Player", i)
		player, _ := player.NewPlayer(name, valueobjects.Mage)
		g.InvitePlayer(guildOwner, player)
	}
}

func resetGuild() {
	guildOwner.UpdateCurrentGuild("")
	guildInstance, _ = CreateGuild(defaultGuildName, guildOwner)
}

func getGuild() *Guild {
	guildOnce.Do(func() {
		guildOwner.UpdateCurrentGuild("")
		guildInstance, _ = CreateGuild(defaultGuildName, guildOwner)
	})
	return guildInstance
}
