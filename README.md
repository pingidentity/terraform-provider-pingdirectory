# PingDirectory Terraform Provider

The PingDirectory Terraform provider is a plugin for [Terraform](https://www.terraform.io/) that supports the management of PingDirectory configuration. This provider is maintained internally by the Ping Identity team.

# Disclaimer - Provider in Development

The PingDirectory Terraform provider is still in development, and breaking changes are likely. As such, it is not yet published on the Terraform registry.

## Requirements
* Terraform 1.1+
* Go 1.18+

# Using the PingDirectory Terraform Provider

The provider can be used to manage PingDirectory servers via the PingDirectory Configuration API. It can replace configuration management that is normally done using the `dsconfig` command-line tool or through dsconfig batch files.

The following example configures a `Location` object and updates the PingDirectory `Global Configuration`.

```
provider "pingdirectory" {
  username = "cn=administrator"
  password = "2FederateM0re"
  https_host = "https://localhost:1443"
  # Warning: The insecure_trust_all_tls attribute configures the provider to trust any certificate presented by the PingDirectory server.
  insecure_trust_all_tls = true
}

resource "pingdirectory_location" "mylocation" {
  id = "MyLocation"
  description = "My description"
}

resource "pingdirectory_global_configuration" "global" {
  location = "Docker"
  encrypt_data = true
}
```

Applying this Terraform configuration file will create a `Location` on the PingDirectory server managed by Terraform, and will update the PingDirectory `Global Configuration` to match the defined resource. Terraform can then manage this configuration, rather than using dsconfig commands such as:

```
dsconfig create-location --location-name MyLocation --set "description:My description"
dsconfig set-location-prop --location-name MyLocation --set "description:My changed description"
dsconfig delete-location --location-name MyLocation
dsconfig set-global-configuration-prop --set encrypt-data:true
```

The provider represents each different configuration object in PingDirectory as a separate resource. The attributes of each resource align with the attributes managed by `dsconfig`.

See the [examples](examples/) directory for more examples using the provider.

## Useful Links

* [Discuss the PingDirectory Terraform Provider](https://support.pingidentity.com/s/topic/0TO1W000000IF30WAG/pingdevops)
* [Ping Identity Home](https://www.pingidentity.com/en.html)

Extended documentation can be found at:
* [PingDirectory Documentation](https://docs.pingidentity.com/r/en-us/pingdirectory-92/pd_ds_landing_page)
* [Ping Identity Developer Portal](https://developer.pingidentity.com/en.html)
* Provider documentation coming soon

## Contributing

We appreciate your help! To contribute through logging issues or creating pull requests, please read the [contribution guidelines](CONTRIBUTING.md)
