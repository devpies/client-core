package memberships

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/google/uuid"
)

// Error codes returned by failures to handle memberships.
var (
	ErrNotFound  = errors.New("membership not found")
	ErrInvalidID = errors.New("id provided was not a valid UUID")
)

func Create(ctx context.Context, repo database.DataStorer, nm NewMembership, now time.Time) (Membership, error) {
	m := Membership{
		ID:        uuid.New().String(),
		UserID:    nm.UserID,
		TeamID:    nm.TeamID,
		Role:      nm.Role,
		UpdatedAt: now.UTC(),
		CreatedAt: now.UTC(),
	}

	stmt := repo.Insert(
		"memberships",
	).SetMap(map[string]interface{}{
		"membership_id": m.ID,
		"user_id":       m.UserID,
		"team_id":       m.TeamID,
		"role":          m.Role,
		"updated_at":    m.UpdatedAt,
		"created_at":    m.CreatedAt,
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return m, err
	}

	return m, nil
}

func RetrieveMemberships(ctx context.Context, repo database.DataStorer, uid, tid string) ([]MembershipEnhanced, error) {
	var m []MembershipEnhanced

	if _, err := q.RetrieveMembership(ctx, repo, uid, tid); err != nil {
		return m, err
	}

	stmt := repo.Select(
		"membership_id",
		"user_id",
		"team_id",
		"email",
		"first_name",
		"last_name",
		"picture",
		"role",
		"m.updated_at",
		"m.created_at",
	).From(
		"memberships as m",
	).Where(sq.Eq{"team_id": "?"}).Join("users USING (user_id)")

	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: arguments (%v)", err, args)
	}

	if err := repo.SelectContext(ctx, &m, q, tid); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return m, nil
}

func RetrieveMembership(ctx context.Context, repo database.DataStorer, uid, tid string) (Membership, error) {
	var m Membership

	if _, err := uuid.Parse(tid); err != nil {
		return m, ErrInvalidID
	}

	if _, err := uuid.Parse(uid); err != nil {
		return m, ErrInvalidID
	}

	stmt := repo.Select(
		"membership_id",
		"user_id",
		"team_id",
		"role",
		"updated_at",
		"created_at",
	).From(
		"memberships",
	).Where(sq.Eq{"team_id": "?", "user_id": "?"})

	query, args, err := stmt.ToSql()
	if err != nil {
		return m, fmt.Errorf("%w: arguments (%v)", err, args)
	}
	err = repo.QueryRowxContext(ctx, q, tid, uid).StructScan(&m)
	if err != nil {
		if err == sql.ErrNoRows {
			return m, ErrNotFound
		}
		return m, err
	}

	return m, nil
}

func Update(ctx context.Context, repo database.DataStorer, tid string, update UpdateMembership, uid string, now time.Time) error {
	m, err := RetrieveMembership(ctx, repo, tid, uid)
	if err != nil {
		return err
	}

	if update.Role != nil {
		m.Role = *update.Role
	}

	stmt := repo.Update(
		"memberships",
	).SetMap(map[string]interface{}{
		"role":       m.Role,
		"updated_at": now.UTC(),
	}).Where(sq.Eq{"team_id": tid, "user_id": uid})

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func Delete(ctx context.Context, repo database.DataStorer, tid, uid string) (string, error) {
	if _, err := uuid.Parse(tid); err != nil {
		return "", ErrInvalidID
	}

	stmt := repo.Delete(
		"memberships",
	).Where(
		sq.Eq{"team_id": "?", "user_id": "?"},
	).Suffix(
		"RETURNING membership_id",
	)

	query, args, err := stmt.ToSql()
	if err != nil {
		return id, fmt.Errorf("%w: arguments (%v)", err, args)
	}

	row := repo.QueryRowxContext(ctx, q, tid, uid)
	var membershipID string

	if err = row.Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}
