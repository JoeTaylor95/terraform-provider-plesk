package main

import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
    "git.home.taylor.am/terraform-provider/terraform-provider-plesk/provider"
)

func main() {
    plugin.Serve(&plugin.ServeOpts{
        ProviderFunc: provider.Provider,
    })
}
