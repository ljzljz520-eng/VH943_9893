package staff

import (
	"context"
	"fmt"
	"strings"

	"campgear/internal/domain"
	"campgear/internal/storage"
)

const (
	RoleManager    = "manager"
	RoleRentalDesk = "rental_desk"
	RoleWarehouse  = "warehouse"
	RoleTechnician = "technician"
)

type Service struct{ repo *storage.Repository }

func NewService(repo *storage.Repository) *Service { return &Service{repo: repo} }

func (s *Service) Enroll(ctx context.Context, id, name, role string) (domain.StaffMember, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(name) == "" {
		return domain.StaffMember{}, fmt.Errorf("staff id and name are required")
	}
	if !ValidRole(role) {
		return domain.StaffMember{}, fmt.Errorf("unsupported staff role")
	}
	member := domain.StaffMember{ID: id, Name: strings.TrimSpace(name), Role: role, Active: true}
	if err := s.repo.InsertStaff(ctx, member); err != nil {
		return domain.StaffMember{}, err
	}
	return member, nil
}

func (s *Service) Deactivate(ctx context.Context, id string) error {
	member, err := s.repo.GetStaff(ctx, id)
	if err != nil {
		return err
	}
	if !member.Active {
		return fmt.Errorf("staff member already inactive")
	}
	return s.repo.SetStaffActive(ctx, id, false)
}

func (s *Service) Get(ctx context.Context, id string) (domain.StaffMember, error) {
	return s.repo.GetStaff(ctx, id)
}

func (s *Service) List(ctx context.Context, activeOnly bool) ([]domain.StaffMember, error) {
	return s.repo.ListStaff(ctx, activeOnly)
}

func (s *Service) CanPerform(ctx context.Context, id, action string) (bool, error) {
	member, err := s.Get(ctx, id)
	if err != nil {
		return false, err
	}
	if !member.Active {
		return false, nil
	}
	switch action {
	case "catalog", "rental":
		return member.Role == RoleManager || member.Role == RoleRentalDesk || member.Role == RoleWarehouse, nil
	case "maintenance":
		return member.Role == RoleManager || member.Role == RoleTechnician, nil
	default:
		return false, fmt.Errorf("unsupported action")
	}
}

func ValidRole(role string) bool {
	switch role {
	case RoleManager, RoleRentalDesk, RoleWarehouse, RoleTechnician:
		return true
	default:
		return false
	}
}
