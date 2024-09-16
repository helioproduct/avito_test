package tender

import "time"

type ServiceType string
type StatusType string

const (
	Construction ServiceType = "Construction"
	Delivery     ServiceType = "Delivery"
	Manufacture  ServiceType = "Manufacture"
)

const (
	Created   StatusType = "Created"
	Published StatusType = "Published"
	Closed    StatusType = "Closed"
)

type Tender struct {
	ID             string      `json:"id"`          // UUID for tender ID
	Name           string      `json:"name"`        // Name of the tender
	Description    string      `json:"description"` // Description of the tender
	ServiceType    ServiceType `json:"serviceType"` // Can be 'Construction', 'Delivery', 'Manufacture'
	OrganizationID string
	CreatorID      string
	Status         StatusType `json:"status"`    // Tender status: CREATED, PUBLISHED, CLOSED
	CurrentVersion int        `json:"version"`   // The current version of the tender
	CreatedAt      time.Time  `json:"createdAt"` // Timestamp of tender creation
	UpdatedAt      time.Time  `json:"updatedAt"` // Timestamp of the last update
}
