/*
 * Copyright (c) 2019 Andreas Signer <asigner@gmail.com>
 *
 * This file is part of Duplikator.
 *
 * Duplikator is free software: you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Duplikator is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Duplikator.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"context"

	"github.com/asig/duplikator/edam"
)

type throttlingUserStore struct {
	us edam.UserStore
}

func (t throttlingUserStore) CheckVersion(ctx context.Context, clientName string, edamVersionMajor int16, edamVersionMinor int16) (r bool, err error) {
	for {
		res, err := t.us.CheckVersion(ctx, clientName, edamVersionMajor, edamVersionMinor)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingUserStore) GetBootstrapInfo(ctx context.Context, locale string) (r *edam.BootstrapInfo, err error) {
	for {
		res, err := t.us.GetBootstrapInfo(ctx, locale)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingUserStore) AuthenticateLongSession(ctx context.Context, username string, password string, consumerKey string, consumerSecret string, deviceIdentifier string, deviceDescription string, supportsTwoFactor bool) (r *edam.AuthenticationResult_, err error) {
	for {
		res, err := t.us.AuthenticateLongSession(ctx, username, password, consumerSecret, consumerSecret, deviceIdentifier, deviceDescription, supportsTwoFactor)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingUserStore) CompleteTwoFactorAuthentication(ctx context.Context, authenticationToken string, oneTimeCode string, deviceIdentifier string, deviceDescription string) (r *edam.AuthenticationResult_, err error) {
	for {
		res, err := t.us.CompleteTwoFactorAuthentication(ctx, authenticationToken, oneTimeCode, deviceIdentifier, deviceDescription)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingUserStore) RevokeLongSession(ctx context.Context, authenticationToken string) (err error) {
	for {
		err = t.us.RevokeLongSession(ctx, authenticationToken)
		if maybeThrottle(err) {
			continue
		}
		return err
	}
}

func (t throttlingUserStore) AuthenticateToBusiness(ctx context.Context, authenticationToken string) (r *edam.AuthenticationResult_, err error) {
	for {
		res, err := t.us.AuthenticateToBusiness(ctx, authenticationToken)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingUserStore) GetUser(ctx context.Context, authenticationToken string) (r *edam.User, err error) {
	for {
		res, err := t.us.GetUser(ctx, authenticationToken)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingUserStore) GetPublicUserInfo(ctx context.Context, username string) (r *edam.PublicUserInfo, err error) {
	for {
		res, err := t.us.GetPublicUserInfo(ctx, username)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingUserStore) GetUserUrls(ctx context.Context, authenticationToken string) (r *edam.UserUrls, err error) {
	for {
		res, err := t.us.GetUserUrls(ctx, authenticationToken)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingUserStore) InviteToBusiness(ctx context.Context, authenticationToken string, emailAddress string) (err error) {
	for {
		err = t.us.InviteToBusiness(ctx, authenticationToken, emailAddress)
		if maybeThrottle(err) {
			continue;
		}
		return err
	}
}

func (t throttlingUserStore) RemoveFromBusiness(ctx context.Context, authenticationToken string, emailAddress string) (err error) {
	for {
		err = t.us.RemoveFromBusiness(ctx, authenticationToken, emailAddress)
		if maybeThrottle(err) {
			continue
		}
		return err
	}
}

func (t throttlingUserStore) UpdateBusinessUserIdentifier(ctx context.Context, authenticationToken string, oldEmailAddress string, newEmailAddress string) (err error) {
	for {
		err = t.us.UpdateBusinessUserIdentifier(ctx, authenticationToken, oldEmailAddress, newEmailAddress)
		if maybeThrottle(err) {
			continue
		}
		return err
	}
}

func (t throttlingUserStore) ListBusinessUsers(ctx context.Context, authenticationToken string) (r []*edam.UserProfile, err error) {
	for {
		res, err := t.us.ListBusinessUsers(ctx, authenticationToken)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingUserStore) ListBusinessInvitations(ctx context.Context, authenticationToken string, includeRequestedInvitations bool) (r []*edam.BusinessInvitation, err error) {
	for {
		res, err := t.us.ListBusinessInvitations(ctx, authenticationToken, includeRequestedInvitations)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}

func (t throttlingUserStore) GetAccountLimits(ctx context.Context, serviceLevel edam.ServiceLevel) (r *edam.AccountLimits, err error) {
	for {
		res, err := t.us.GetAccountLimits(ctx, serviceLevel)
		if maybeThrottle(err) {
			continue
		}
		return res, err
	}
}
