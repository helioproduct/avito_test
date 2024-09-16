package usecase

import (
	"avito_api/internal/entities/tender"
	"errors"
	"reflect"

	"log"
	"time"
)

var (
	ErrNoSuchTender       = errors.New("tender does not exists")
	ErrInvalidServiceType = errors.New("invalid service type")
)

type tenderUseCase struct {
	tenderRepo TenderRepository
	orgRepo    OrganizationRepository
	userRepo   UserRepository
}

func NewTenderUseCase(tenderRepo TenderRepository, orgRepo OrganizationRepository,
	userRepo UserRepository) TenderUseCase {

	return &tenderUseCase{
		tenderRepo: tenderRepo,
		orgRepo:    orgRepo,
		userRepo:   userRepo,
	}
}

func (uc *tenderUseCase) CreateTender(t *tender.Tender, username string) (string, error) {
	const op = "usecase.tenders.CreateTender:"

	responsible, err := uc.orgRepo.IsUserResponsibleForOrganization(username, t.OrganizationID)
	if err != nil {
		log.Println(op, err)
		return "", err
	}
	if !responsible {
		return "", ErrUserNotResponsible
	}

	t.Status = tender.Created
	t.CurrentVersion = 1
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	return uc.tenderRepo.CreateTender(t, username)
}

func (uc *tenderUseCase) GetStatus(tenderID, username string) (tender.StatusType, error) {
	const op = "usecase.tender.GetStatu:"

	tender, err := uc.tenderRepo.GetByID(tenderID)
	if err != nil || tender == nil {
		log.Println(op, err)
		return "", ErrNoSuchTender
	}

	responsible, err := uc.orgRepo.IsUserResponsibleForOrganization(username, tender.OrganizationID)
	if err != nil {
		log.Println(op, err)
		return "", err
	}
	if !responsible {
		return "", ErrUserNotResponsible
	}
	return tender.Status, nil
}

func (uc *tenderUseCase) ChangeStatus(tenderID string, username string, newStatus tender.StatusType) error {
	const op = "usecase.tender.ChangeStatus:"

	currentTender, err := uc.tenderRepo.GetByID(tenderID)
	if err != nil {
		log.Println(op, err)
		return err
	}

	responsible, err := uc.orgRepo.IsUserResponsibleForOrganization(username, currentTender.OrganizationID)
	if err != nil {
		log.Println(op, err)
		return err
	}
	if !responsible {
		return ErrUserNotResponsible
	}

	if currentTender.Status == newStatus {
		return errors.New("tender is already in the desired status")
	}
	return uc.tenderRepo.ChangeStatus(tenderID, username, newStatus)
}

// Edit: Name, Description, ServiceType
func (uc *tenderUseCase) EditTender(updatedTender *tender.Tender, username string) error {
	const op = "usecase.tender.EditTender:"

	currentTender, err := uc.tenderRepo.GetByID(updatedTender.ID)
	if err != nil {
		log.Println(op, err)
		return ErrNoSuchTender
	}

	responsible, err := uc.orgRepo.IsUserResponsibleForOrganization(username, currentTender.OrganizationID)
	if err != nil {
		log.Println(op, err)
		return err
	}
	if !responsible {
		return ErrUserNotResponsible
	}

	updateTenderStruct(currentTender, *updatedTender)
	currentTender.CurrentVersion += 1
	currentTender.UpdatedAt = time.Now()

	return uc.tenderRepo.EditTender(currentTender, username)
}

// if limit=-1 returns all published tenders
func (uc *tenderUseCase) GetTenders(limit, offset int, serviceType string) ([]*tender.Tender, error) {
	const op = "usecase.tender.GetTenders:"

	publishedTenders, err := uc.tenderRepo.GetPublishedTenders(limit, offset, serviceType)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}
	return publishedTenders, nil
}

func (uc *tenderUseCase) GetMyTenders(username string, limit, offset int) ([]*tender.Tender, error) {
	const op = "usecase.tender.GetMyTenders:"
	tenders, err := uc.tenderRepo.GetOwnedTenders(limit, offset, username)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}
	return tenders, nil
}

func (uc *tenderUseCase) GetUserID(username string) (string, error) {
	const op = "usecase.tender.GetUserID:"

	user, err := uc.userRepo.GetUserByUsername(username)
	if err != nil || user == nil {
		log.Println(op, err)
		return "", err
	}
	return user.ID, nil
}

func (uc *tenderUseCase) GetByID(tenderID, username string) (*tender.Tender, error) {
	const op = "usecase.tender.GetByID:"
	availableTenders, err := uc.GetMyTenders(username, -1, 0)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}

	for _, t := range availableTenders {
		if t.ID == tenderID {
			return t, nil
		}
	}
	return nil, ErrNoSuchTender
}

func updateTenderStruct(a *tender.Tender, b tender.Tender) {
	vA := reflect.ValueOf(a).Elem()
	vB := reflect.ValueOf(b)

	for i := 0; i < vA.NumField(); i++ {
		fieldA := vA.Field(i)
		fieldB := vB.Field(i)

		if fieldA.CanSet() {
			switch fieldA.Kind() {
			case reflect.String:
				if fieldB.String() != "" {
					fieldA.SetString(fieldB.String())
				}
			case reflect.Int:
				if fieldB.Int() != 0 {
					fieldA.SetInt(fieldB.Int())
				}
			case reflect.Bool:
				fieldA.SetBool(fieldB.Bool())
			}
		}
	}
}
