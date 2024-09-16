package usecase

import (
	"avito_api/internal/entities/bid"
	"avito_api/internal/entities/organization"
	"avito_api/internal/entities/tender"
	"avito_api/internal/entities/user"
)

type (
	BidUseCase interface {
		Create(b *bid.Bid) (*bid.Bid, error)
		GetBidStatus(bidID, username string) (bid.BidStatus, error)
		ChangeStatus(bidID string, username string, newStatus bid.BidStatus) (*bid.Bid, error)
		Edit(updatedBid *bid.Bid, username string) (*bid.Bid, error)
		GetBidByID(bidID string) (*bid.Bid, error)
		ChangeDecision(bidID string, username string, decision bid.Decision) (*bid.Bid, error)
		GetBidsForTender(tenderID string, username string, limit, offset int) ([]*bid.Bid, error)
		GetBidsByUser(username string, limit, offset int) ([]*bid.Bid, error)
	}

	BidRepository interface {
		CreateBid(bid *bid.Bid) (string, error)
		GetBidByID(bidID string) (*bid.Bid, error)
		EditBid(updatedBid *bid.Bid) error
		ChangeStatus(bidID, username string, newStatus bid.BidStatus) error
		ChangeDecision(bidID, username string, decision bid.Decision) error
		GetBids(tenderID string, username string, limit, offset int) ([]*bid.Bid, error)
		GetBidsByUser(username string, limit, offset int) ([]*bid.Bid, error)
	}

	UserUseCase interface {
		GetUserByID(id string) (*user.User, error)
		GetUserByUsername(username string) (*user.User, error)
		GetAllUsers() ([]*user.User, error)
	}

	UserRepository interface {
		CreateUser(user *user.User) error
		GetUserByID(id string) (*user.User, error)
		GetUserByUsername(username string) (*user.User, error)
		GetAllUsers() ([]*user.User, error)
	}

	TenderUseCase interface {
		CreateTender(t *tender.Tender, username string) (string, error)
		GetStatus(tenderID, username string) (tender.StatusType, error)
		ChangeStatus(tenderID string, username string, newStatus tender.StatusType) error
		EditTender(updatedTender *tender.Tender, username string) error
		GetTenders(limit, offset int, serviceType string) ([]*tender.Tender, error)
		GetMyTenders(username string, limit, offset int) ([]*tender.Tender, error)
		GetUserID(username string) (string, error)
		GetByID(tenderID, username string) (*tender.Tender, error)
	}

	TenderRepository interface {
		CreateTender(tender *tender.Tender, username string) (string, error)
		GetStatus(tenderID string) (tender.StatusType, error)
		GetByID(tenderID string) (*tender.Tender, error)
		GetVersionByID(tenderID string, version int) (*tender.Tender, error)
		ChangeStatus(tenderID string, username string, newStatus tender.StatusType) error
		EditTender(updatedTender *tender.Tender, username string) error
		GetTendersByResponsibleUser(username string, limit, offset int, serviceType string) ([]*tender.Tender, error)
		GetPublishedTenders(limit, offset int, serviceType string) ([]*tender.Tender, error)
		GetOwnedTenders(limit, offset int, username string) ([]*tender.Tender, error)
	}

	OrganizationUseCase interface {
		CreateOrganization(org *organization.Organization) (string, error)
		GetResponsibleUsers(organizationID string) ([]user.User, error)
		SetResponsibleUsers(organizationID string, userIDs []string) error
		GetOrganizationByName(name string) (*organization.Organization, error)
		GetAllOrganizations() ([]*organization.Organization, error)
		GetResponsibleUsersForOrganization() ([]*OrganizationResponsible, error)
		GetUserOrganizationID(userID string) (string, error)
	}

	OrganizationRepository interface {
		Create(organization *organization.Organization) (string, error)
		GetByName(name string) (*organization.Organization, error)
		GetResponsibleUsers(organizationID string) ([]user.User, error)
		SetResponsibleUsers(organizationID string, userIDs []string) error
		IsUserResponsibleForOrganization(username, organizationID string) (bool, error)
		GetAllOrganizations() ([]*organization.Organization, error)
		GetUserOrganizationID(userID string) (string, error)
	}
)
