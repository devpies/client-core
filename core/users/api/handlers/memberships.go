package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"github.com/devpies/devpie-client-core/users/domain/memberships"
	"github.com/devpies/devpie-client-core/users/platform/auth0"
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/devpies/devpie-client-core/users/platform/web"
	"github.com/devpies/devpie-client-events/go/events"
)

type Membership struct {
	repo  database.Storer
	log   *log.Logger
	auth0 *auth0.Auth0
	nats  *events.Client
	query MembershipQueries
}

type MembershipQueries struct {
	membership MembershipQuerier
}

type MembershipQuerier interface {
	Create(ctx context.Context, repo database.Storer, nm memberships.NewMembership, now time.Time) (memberships.Membership, error)
	RetrieveMemberships(ctx context.Context, repo database.Storer, uid, tid string) ([]memberships.MembershipEnhanced, error)
	RetrieveMembership(ctx context.Context, repo database.Storer, uid, tid string) (memberships.Membership, error)
	Update(ctx context.Context, repo database.Storer, tid string, update memberships.UpdateMembership, uid string, now time.Time) error
	Delete(ctx context.Context, repo database.Storer, tid, uid string) (string, error)
}

func (m *Membership) RetrieveMembers(w http.ResponseWriter, r *http.Request) error {
	uid := m.auth0.UserByID(r.Context())

	tid := chi.URLParam(r, "tid")

	ms, err := m.query.membership.RetrieveMemberships(r.Context(), m.repo, uid, tid)
	if err != nil {
		switch err {
		case memberships.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case memberships.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for team %q", tid)
		}
	}

	return web.Respond(r.Context(), w, ms, http.StatusOK)
}
