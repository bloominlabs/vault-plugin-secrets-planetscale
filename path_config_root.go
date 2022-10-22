package planetscale

import (
	"context"
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const configRootKey = "config/root"

func pathConfigToken(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: configRootKey,
		Fields: map[string]*framework.FieldSchema{
			"service_token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "service token to generate planetscale credentials with",
			},
			"service_token_name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "name of service token used to generate planetscale credentials with",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathConfigTokenRead,
			logical.CreateOperation: b.pathConfigTokenWrite,
			logical.UpdateOperation: b.pathConfigTokenWrite,
			logical.DeleteOperation: b.pathConfigTokenDelete,
		},

		ExistenceCheck: b.configTokenExistenceCheck,
	}
}

func (b *backend) configTokenExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := b.readConfigToken(ctx, req.Storage)
	if err != nil {
		return false, err
	}

	return entry != nil, nil
}

func (b *backend) readConfigToken(ctx context.Context, storage logical.Storage) (*rootConfig, error) {
	entry, err := storage.Get(ctx, configRootKey)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	conf := &rootConfig{}
	if err := entry.DecodeJSON(conf); err != nil {
		return nil, errwrap.Wrapf("error reading nomad access configuration: {{err}}", err)
	}

	return conf, nil
}

func (b *backend) pathConfigTokenRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	conf, err := b.readConfigToken(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if conf == nil {
		return logical.ErrorResponse(
			fmt.Sprintf("configuration does not exist. did you configure '%s'?", configRootKey),
		), nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"service_token":      conf.ServiceToken,
			"service_token_name": conf.ServiceTokenName,
		},
	}, nil
}

func (b *backend) pathConfigTokenWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	conf, err := b.readConfigToken(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if conf == nil {
		conf = &rootConfig{}
	}

	serviceToken, ok := data.GetOk("service_token")
	if !ok {
		return logical.ErrorResponse("Missing 'service_token' in configuration request"), nil
	}
	conf.ServiceToken = serviceToken.(string)

	serviceTokenName, ok := data.GetOk("service_token_name")
	if !ok {
		return logical.ErrorResponse("Missing 'service_token_name' in configuration request"), nil
	}
	conf.ServiceTokenName = serviceTokenName.(string)

	b.clearClients()

	entry, err := logical.StorageEntryJSON(configRootKey, conf)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathConfigTokenDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, configRootKey); err != nil {
		return nil, err
	}
	return nil, nil
}

type rootConfig struct {
	ServiceToken     string `json:"service_token,omitempty"`
	ServiceTokenName string `json:"service_token_name,omitempty"`
}

const pathConfigTokenHelpSyn = `
Configure [planetscale-go client](https://github.com/planetscale/planetscale-go) used by vault
`

const pathConfigTokenHelpDesc = `
Will configure this mount with the planetscale service token / service token name used by Vault for all planetscale
operations on this mount. 
`
