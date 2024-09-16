package usecase

import (
	"avito_api/internal/entities/organization"
	"avito_api/internal/entities/user"
	"errors"
	"log"
	"time"
)

var (
	ErrUserNotResponsible = errors.New("user is not responsible for this organization")
)

type organizationUseCase struct {
	orgRepo OrganizationRepo
}

func NewOrganizationUseCase(repo OrganizationRepo) OrganizationUseCase {
	return &organizationUseCase{
		orgRepo: repo,
	}
}

func (uc *organizationUseCase) GetUserOrganizationID(userID string) (string, error) {
	return uc.orgRepo.GetUserOrganizationID(userID)
}

func (uc *organizationUseCase) CreateOrganization(org *organization.Organization) (string, error) {
	const op = "usecase.organization.CreateOrganization:"

	org.CreatedAt = time.Now()
	org.UpdatedAt = time.Now()

	orgID, err := uc.orgRepo.Create(org)
	if err != nil {
		log.Println(op, err)
		return "", err
	}
	return orgID, nil
}

func (uc *organizationUseCase) GetResponsibleUsers(organizationID string) ([]user.User, error) {
	const op = "usesase.GetResponsibleUsers:"

	users, err := uc.orgRepo.GetResponsibleUsers(organizationID)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}
	return users, nil
}

func (uc *organizationUseCase) SetResponsibleUsers(organizationID string, userIDs []string) error {
	const op = "usesase.SetResponsibleUsers:"

	err := uc.orgRepo.SetResponsibleUsers(organizationID, userIDs)
	if err != nil {
		log.Println(op, err)
		return err
	}
	return nil
}

func (uc *organizationUseCase) GetOrganizationByName(name string) (*organization.Organization, error) {
	const op = "usesase.GetOrganizationByName:"

	org, err := uc.orgRepo.GetByName(name)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}
	return org, nil
}

func (uc *organizationUseCase) GetAllOrganizations() ([]*organization.Organization, error) {
	const op = "usesase.GetOrganizationByName:"

	allOrganizations, err := uc.orgRepo.GetAllOrganizations()
	if err != nil {
		log.Println(op, err)
		return nil, err
	}
	return allOrganizations, nil
}

type OrganizationResponsible struct {
	Organization *organization.Organization `json:"organization"`
	// OrganizationName string       `json:"organization_name"`
	ResponsibleUsers *[]user.User `json:"responsilbe_users"`
}

func (uc *organizationUseCase) GetResponsibleUsersForOrganization() ([]*OrganizationResponsible, error) {
	const op = "usesase.GetOrganizationsResponsibleUsers:"

	allOrganizations, err := uc.GetAllOrganizations()
	if err != nil {
		log.Println(op, err)
	}

	var orgsResponsibles []*OrganizationResponsible

	for _, org := range allOrganizations {

		var orgResponsible OrganizationResponsible

		users, err := uc.GetResponsibleUsers(org.ID)
		if err != nil {
			log.Println(op, err)
			return nil, err
		}
		// orgResponsible.OrganizationName = org.Name
		orgResponsible.Organization = org
		orgResponsible.ResponsibleUsers = &users
		orgsResponsibles = append(orgsResponsibles, &orgResponsible)
	}
	return orgsResponsibles, nil
}
