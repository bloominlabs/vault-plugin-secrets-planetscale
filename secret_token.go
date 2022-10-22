package planetscale

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/planetscale/planetscale-go/planetscale"
)

const (
	SecretTokenType = "token"
)

func secretToken(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretTokenType,
		Fields: map[string]*framework.FieldSchema{
			"id": {
				Type:        framework.TypeString,
				Description: "ID of the planetscale password",
			},
			"username": {
				Type:        framework.TypeString,
				Description: "username to connect to planetscale database with",
			},
			"password": {
				Type:        framework.TypeString,
				Description: "password to connect to the planetscale database with",
			},
			"organization": {
				Type:        framework.TypeString,
				Description: "planetscale organization the password was created for",
			},
			"branch": {
				Type:        framework.TypeString,
				Description: "planetscale database branch the password was created for",
			},
			"database": {
				Type:        framework.TypeString,
				Description: "planetscale database the password was created for",
			},
		},

		Revoke: b.secretTokenRevoke,
		Renew:  b.secretTokenRenew,
	}
}

func (b *backend) secretTokenRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	lease, err := b.LeaseConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		lease = &configLease{}
	}

	// TODO: since we have the username and password, we can regenerate the
	// access token after 24 hours; however, defaulting to 24 hours now since
	// JWTs returned from auth0 expire in 24 hours. Main concern is even though
	// we can regenerate the accessToken, I'm not sure how we expose that to the
	// user.
	resp := &logical.Response{Secret: req.Secret}
	resp.Secret.TTL = lease.TTL
	resp.Secret.MaxTTL = time.Hour * 24
	return resp, nil
}

func (b *backend) secretTokenRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	c, err := b.client(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, fmt.Errorf("error getting planetscale client")
	}

	id, ok := req.Secret.InternalData["id"]
	if !ok {
		return nil, fmt.Errorf("'id' is missing on the lease")
	}
	username, ok := req.Secret.InternalData["username"]
	if !ok {
		return nil, fmt.Errorf("'username' is missing on the lease")
	}
	database, ok := req.Secret.InternalData["database"]
	if !ok {
		return nil, fmt.Errorf("'database' is missing on the lease")
	}
	organization, ok := req.Secret.InternalData["organization"]
	if !ok {
		return nil, fmt.Errorf("'organization' is missing on the lease")
	}
	branch, ok := req.Secret.InternalData["branch"]
	if !ok {
		return nil, fmt.Errorf("'branch' is missing on the lease")
	}

	b.Logger().Info(fmt.Sprintf("deleting planetscale password. id: %s, username: %s", id, username))
	err = b.planetscaleClient.Passwords.Delete(ctx, &planetscale.DeleteDatabaseBranchPasswordRequest{
		PasswordId:   id.(string),
		Organization: organization.(string),
		Database:     database.(string),
		Branch:       branch.(string),
	})
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("failed to delete planetscale password. id: %s, username: %s, err: %s", id, username, err)), nil
	}

	return nil, nil
}
