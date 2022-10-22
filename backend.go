package planetscale

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/planetscale/planetscale-go/planetscale"
)

// backend wraps the backend framework and adds a map for storing key value pairs
type backend struct {
	*framework.Backend

	planetscaleClientMutex sync.RWMutex
	planetscaleClient      *planetscale.Client
}

var _ logical.Factory = Factory

// Factory configures and returns Mock backends
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b, err := newBackend()
	if err != nil {
		return nil, err
	}

	if conf == nil {
		return nil, fmt.Errorf("configuration passed into backend is nil")
	}

	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}

	return b, nil
}

func newBackend() (*backend, error) {
	b := &backend{}

	b.Backend = &framework.Backend{
		Help:        strings.TrimSpace(backendHelp),
		BackendType: logical.TypeLogical,
		Paths: framework.PathAppend(
			b.paths(),
		),
		Secrets: []*framework.Secret{
			secretToken(b),
		},
		Invalidate: b.invalidate,
	}

	return b, nil
}

func (b *backend) paths() []*framework.Path {
	return []*framework.Path{
		pathConfigToken(b),
		pathCredsCreate(b),
		pathRoles(b),
		pathListRoles(b),
		pathConfigLease(b),
	}
}

func (b *backend) invalidate(ctx context.Context, key string) {
	switch {
	case key == configRootKey:
		b.clearClients()
	}
}

func (b *backend) clearClients() {
	b.planetscaleClientMutex.Lock()
	defer b.planetscaleClientMutex.Unlock()
	b.planetscaleClient = nil
}

const backendHelp = `
  The TIP secret plugin generates test users for use in stratos.host canaries / tests.

	After mounting this backend, an auth0 client (with the ability to generate
	TIP users) must be configured with the 'config/root' path and policies must
	be written using the "roles/" endpoints before any access keys can be
	generated.
`
