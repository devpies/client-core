package handlers

import (
	"github.com/devpies/devpie-client-core/users/domain/users"
	"github.com/devpies/devpie-client-core/users/platform/auth0"
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/devpies/devpie-client-core/users/platform/web"
	"github.com/pkg/errors"

	//"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
)

type Users struct {
	repo    database.Storer
	log     *log.Logger
	auth0   auth0.Auther
	origins string
	query   users.Querier
}

func (u *Users) RetrieveMe(w http.ResponseWriter, r *http.Request) error {
	var us users.User

	uid := u.auth0.GetUserByID(r.Context())

	if uid == "" {
		return web.NewRequestError(users.ErrNotFound, http.StatusNotFound)
	}
	us, err := u.query.RetrieveMe(r.Context(), u.repo, uid)
	if err != nil {
		switch err {
		case users.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case users.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for user %q", uid)
		}
	}

	return web.Respond(r.Context(), w, us, http.StatusOK)
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) error {
	var nu users.NewUser

	if err := web.Decode(r, &nu); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	// get auth0 management api token
	t, err := u.auth0.GetOrCreateToken()
	if err != nil {
		return err
	}

	// if user already exists update app metadata only
	us, err := u.query.RetrieveMeByAuthID(r.Context(), u.repo, nu.Auth0ID)
	if err == nil {
		// update app metadata for existing user
		if err = u.auth0.UpdateUserAppMetaData(t, nu.Auth0ID, us.ID); err != nil {
			switch err {
			case auth0.ErrInvalidID:
				return web.NewRequestError(err, http.StatusBadRequest)
			default:
				return errors.Wrapf(err,"failed to update user app metadata")
			}
		} // mock
		return web.Respond(r.Context(), w, us, http.StatusAccepted)
	}

	user, err := u.query.Create(r.Context(), u.repo, nu, time.Now())
	if err != nil {
		return err
	}

	// update app metadata for new user
	if err := u.auth0.UpdateUserAppMetaData(t, user.Auth0ID, user.ID); err != nil {
		switch err {
		case auth0.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err,"failed to update user app metadata")
		}
	}

	return web.Respond(r.Context(), w, user, http.StatusCreated)
}
