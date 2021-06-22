package teams

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Error codes returned by failures to handle teams.
var (
	ErrNotFound  = errors.New("team not found")
	ErrInvalidID = errors.New("id provided was not a valid UUID")
)

func Create(ctx context.Context, repo database.DataStorer, nt NewTeam, uid string, now time.Time) (Team, error) {
	t := Team{
		ID:        uuid.New().String(),
		Name:      nt.Name,
		UserID:    uid,
		UpdatedAt: now.UTC(),
		CreateAt:  now.UTC(),
	}

	stmt := repo.Insert(
		"teams",
	).SetMap(map[string]interface{}{
		"team_id":    t.ID,
		"name":       t.Name,
		"user_id":    t.UserID,
		"updated_at": t.UpdatedAt,
		"created_at": t.CreateAt,
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return t, errors.Wrap(err, "inserting team")
	}

	return t, nil
}

func Retrieve(ctx context.Context, repo database.DataStorer, tid string) (Team, error) {
	var t Team

	if _, err := uuid.Parse(tid); err != nil {
		return t, ErrInvalidID
	}

	stmt := repo.Select(
		"team_id",
		"user_id",
		"name",
		"updated_at",
		"created_at",
	).From(
		"teams",
	).Where(sq.Eq{"team_id": "?"})

	q, args, err := stmt.ToSql()
	if err != nil {
		return t, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.GetContext(ctx, &t, q, tid); err != nil {
		if err == sql.ErrNoRows {
			return t, ErrNotFound
		}
		return t, err
	}

	return t, nil
}

func List(ctx context.Context, repo database.DataStorer, uid string) ([]Team, error) {
	var ts []Team

	if _, err := uuid.Parse(uid); err != nil {
		return ts, ErrInvalidID
	}

	stmt := repo.Select(
		"team_id",
		"user_id",
		"name",
		"updated_at",
		"created_at",
	).From(
		"teams",
	).Where("team_id IN (SELECT team_id FROM memberships WHERE user_id = ?)")

	q, args, err := stmt.ToSql()
	if err != nil {
		return ts, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.SelectContext(ctx, &ts, q, uid); err != nil {
		if err == sql.ErrNoRows {
			return ts, ErrNotFound
		}
		return ts, err
	}

	return ts, nil
}
