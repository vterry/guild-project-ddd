package valueobjects

type GuildID struct {
	BaseID[string]
}

func (pID GuildID) Equals(otherID GuildID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
