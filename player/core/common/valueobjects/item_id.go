package valueobjects

type ItemID struct {
	BaseID[string]
}

func (i ItemID) Equals(otherID ItemID) bool {
	return i.BaseID.Equals(otherID.BaseID)
}
