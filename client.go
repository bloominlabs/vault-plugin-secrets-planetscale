package planetscale

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/planetscale/planetscale-go/planetscale"
)

func (b *backend) client(ctx context.Context, s logical.Storage) (*planetscale.Client, error) {
	b.planetscaleClientMutex.RLock()
	if b.planetscaleClient != nil {
		b.planetscaleClientMutex.RUnlock()
		return b.planetscaleClient, nil
	}

	// upgrade the lock since we need to create the client
	b.planetscaleClientMutex.RUnlock()
	b.planetscaleClientMutex.Lock()
	defer b.planetscaleClientMutex.Unlock()

	// check the client again, in the event that a client was being created while we waited for Lock()
	if b.planetscaleClient != nil {
		return b.planetscaleClient, nil
	}

	conf, err := b.readConfigToken(ctx, s)
	if err != nil {
		return nil, err
	}

	b.Logger().Warn("tokens", "name", conf.ServiceTokenName, "token", conf.ServiceToken)
	client, err := planetscale.NewClient(
		planetscale.WithServiceToken(conf.ServiceTokenName, conf.ServiceToken),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to intialize planetscale client: %w", err)
	}

	b.planetscaleClient = client
	return b.planetscaleClient, nil
}
