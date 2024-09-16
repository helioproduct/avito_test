package organization

import "time"

type OrganizationType string

const (
	IE  OrganizationType = "IE"
	LLC OrganizationType = "LLC"
	JSC OrganizationType = "JSC"
)

type Organization struct {
	ID          string           `json:"id"`          // UUID of the organization
	Name        string           `json:"name"`        // Name of the organization
	Description string           `json:"description"` // Description of the organization
	Type        OrganizationType `json:"type"`        // Organization type (IE, LLC, JSC)
	CreatedAt   time.Time        `json:"created_at"`  // Timestamp when the organization was created
	UpdatedAt   time.Time        `json:"updated_at"`  // Timestamp when the organization was last updated
}
