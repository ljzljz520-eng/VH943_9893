package storage

import (
	"context"
	"database/sql"

	"campgear/internal/domain"
)

func (r *Repository) InsertStaff(ctx context.Context, member domain.StaffMember) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO staff_members (id, name, role, active) VALUES (?, ?, ?, ?)`, member.ID, member.Name, member.Role, boolInt(member.Active))
	return err
}

func (r *Repository) GetStaff(ctx context.Context, id string) (domain.StaffMember, error) {
	var member domain.StaffMember
	var active int64
	err := r.db.QueryRowContext(ctx, `SELECT id, name, role, active FROM staff_members WHERE id=?`, id).Scan(&member.ID, &member.Name, &member.Role, &active)
	if err != nil {
		return domain.StaffMember{}, wrapNotFound(err)
	}
	member.Active = intBool(active)
	return member, nil
}

func (r *Repository) ListStaff(ctx context.Context, activeOnly bool) ([]domain.StaffMember, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, role, active FROM staff_members WHERE (?=0 OR active=1) ORDER BY name, id`, boolInt(activeOnly))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.StaffMember, 0)
	for rows.Next() {
		var member domain.StaffMember
		var active int64
		if err := rows.Scan(&member.ID, &member.Name, &member.Role, &active); err != nil {
			return nil, err
		}
		member.Active = intBool(active)
		result = append(result, member)
	}
	return result, rows.Err()
}

func (r *Repository) SetStaffActive(ctx context.Context, id string, active bool) error {
	result, err := r.db.ExecContext(ctx, `UPDATE staff_members SET active=? WHERE id=?`, boolInt(active), id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return sql.ErrNoRows
	}
	return nil
}
