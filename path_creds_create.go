package planetscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/planetscale/planetscale-go/planetscale"
)

// maxTokenNameLength is the maximum length for the name of a Nomad access
// token
const maxTokenNameLength = 120

func createDisplayName(role string) string {
	lowerRole := strings.ToLower(role)

	name := fmt.Sprintf("vault-%s-%d", lowerRole, time.Now().UnixNano())

	// Note: if the given role name is sufficiently long, the UnixNano() portion
	// of the pseudo randomized token name is the part that gets trimmed off,
	// weakening it's randomness.
	if len(name) > maxTokenNameLength {
		name = name[:maxTokenNameLength]
	}

	return name
}

func pathCredsCreate(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("role"),
		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type:        framework.TypeString,
				Description: "Create a planetscale user from a role",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathCredsRead,
		},
	}
}

func (b *backend) pathCredsRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	role := d.Get("role").(string)

	roleEntry, err := b.roleRead(ctx, req.Storage, role)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("err while getting role configuration for '%s'. err: %s", role, err)), nil
	}
	if roleEntry == nil {
		return logical.ErrorResponse(fmt.Sprintf("could not find entry for role '%s', did you configure it?", role)), nil
	}

	c, err := b.client(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	lease, err := b.LeaseConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		lease = &configLease{}
	}

	ttl, _, err := framework.CalculateTTL(b.System(), 0, lease.TTL, 0, lease.MaxTTL, 0, time.Time{})
	if err != nil {
		return logical.ErrorResponse("failed to calculate ttl. err: %w", err), nil
	}

	password, err := c.Passwords.Create(ctx, &planetscale.DatabaseBranchPasswordRequest{
		Organization: roleEntry.Organization,
		Database:     roleEntry.Database,
		Branch:       roleEntry.Branch,
		Role:         roleEntry.Role,
		DisplayName:  createDisplayName(role),
	})
	if err != nil {
		return logical.ErrorResponse("failed to create planetscale password. err: %s", err), nil
	}

	// Use the helper to create the secret
	resp := b.Secret(SecretTokenType).Response(map[string]interface{}{
		"id": password.PublicID,
		// TODO: update the planetscale client to get this dynamically
		"host":         "us-west.connect.psdb.cloud",
		"username":     password.Username,
		"password":     password.PlainText,
		"database":     roleEntry.Database,
		"branch":       roleEntry.Branch,
		"organization": roleEntry.Organization,
		"role":         roleEntry.Role,
	}, map[string]interface{}{
		"id":           password.PublicID,
		"username":     password.Username,
		"password":     password.PlainText,
		"database":     roleEntry.Database,
		"branch":       roleEntry.Branch,
		"organization": roleEntry.Organization,
	})
	resp.Secret.TTL = ttl
	// TODO: reconsider this since we can regenerate the tokens with the username
	// / passwords; however, this may require some extra work where I don't
	// expect a user to be used for >= 24 hours.
	resp.Secret.MaxTTL = lease.MaxTTL
	return resp, nil
}
