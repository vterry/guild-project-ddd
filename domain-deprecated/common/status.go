package common

type Status int

const (
	Approved Status = iota
	Rejected
	Pending
	Canceled
)

func (c Status) String() string {
	return [...]string{"APPROVED", "REJECTED", "PENDING", "CANCELED"}[c]
}
