package bid

import "time"

type BidStatus string

type AuthorType string

type Decision string

const (
	Created   BidStatus = "Created"
	Published BidStatus = "Published"
	Canceled  BidStatus = "Canceled"
)

const (
	User         AuthorType = "User"
	Organization AuthorType = "Organization"
)

const (
	Approved Decision = "Approved"
	Rejected Decision = "Rejected"
)

type Bid struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	TenderID    string     `json:"tenderId"`
	AuthorType  AuthorType `json:"authorType"`
	AuthorID    string     `json:"authorId"`
	Status      BidStatus  `json:"status"`
	Decision    Decision   `json:"decision"`
	Version     int        `json:"version"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
