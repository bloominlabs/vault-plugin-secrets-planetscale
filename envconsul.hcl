# envconsul -config envconsul.hcl <your command>
upcase = true

consul {
  address = "https://nomad-servers.prod.stratos.host:8501"
  retry {
    attempts = 3
  }
}

vault {
  // don't renew token as its most likely user token 
  // (vault login -method=oidc role=developer) 
  // that cannot be renewed
  address     = "https://vault.prod.stratos.host:8200"
  renew_token = false
  retry {
    attempts = 3
  }
}

prefix {
  path   = "planetscale"
  format = "TIP_{{ key }}"
}

secret {
  path      = "infra/data/planetscale/auth0-oidc-app-keys"
  no_prefix = true
  format    = "AUTH0_{{ key }}"
}
