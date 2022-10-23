# vault-plugin-secrets-planetscale

Generate @planetscale usernames and passwords using vault.

## Usage

### Setup Endpoint

1. Download and enable plugin locally (TODO)

2. Configure the plugin

   ```
   vault write /planetscale/config/root service_token=<service_token> service_token_name=<service_token_id>
   ```

3. Add one or more policies

### Configure Policies

```
vault write planetscale/roles/fjord organization=bloominlabs database=bloominlabs role=admin branch=main
```

you can then read from the role using

```
vault read /planetscale/creds/<role-name>
```

## Development

The provided [Earthfile] ([think makefile, but using
docker](https://earthly.dev)) is used to build, test, and publish the plugin.
See the build targets for more information. Common targets include

```bash
# build a local version of the plugin
$ earthly +build

# start vault and enable the plugin locally
earthly +dev
```

[vault]: https://www.vaultproject.io/
[planetscale]: https://planetscale.com/
[earthfile]: ./Earthfile
