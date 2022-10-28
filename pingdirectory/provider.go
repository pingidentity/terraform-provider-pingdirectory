package pingdirectory

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// pingdirectoryProviderModel maps provider schema data to a Go type.
//TODO add default user password to model
type pingdirectoryProviderModel struct {
	Host                types.String `tfsdk:"host"`
	Username            types.String `tfsdk:"username"`
	Password            types.String `tfsdk:"password"`
	DefaultUserPassword types.String `tfsdk:"default_user_password"`
}

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &pingdirectoryProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &pingdirectoryProvider{}
}

// pingdirectoryProvider is the provider implementation.
type pingdirectoryProvider struct{}

// Metadata returns the provider type name.
func (p *pingdirectoryProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pingdirectory"
}

// GetSchema defines the provider-level schema for configuration data.
func (p *pingdirectoryProvider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "PingDirectory POC Provider.",
		Attributes: map[string]tfsdk.Attribute{
			"host": {
				Description: "URI for PingDirectory LDAP port.",
				Type:        types.StringType,
				Required:    true,
			},
			"username": {
				Description: "Username for PingDirectory admin user.",
				Type:        types.StringType,
				Required:    true,
			},
			"password": {
				Description: "Password for PingDirectory admin user.",
				Type:        types.StringType,
				Required:    true,
				Sensitive:   true,
			},
			"default_user_password": {
				Description: "Default user password for created PingDirectory users.",
				Type:        types.StringType,
				Required:    true,
				Sensitive:   true,
			},
		},
	}, nil
}

// Configure prepares a PingDirectory LDAP client
func (p *pingdirectoryProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring PingDirectory client")

	// Retrieve provider data from configuration
	var config pingdirectoryProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown PingDirectory Host",
			"The provider cannot create the PingDirectory client as there is an unknown configuration value for the PingDirectory host. "+
				"Either target apply the source of the value first or set the value statically in the configuration.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown PingDirectory Username",
			"The provider cannot create the PingDirectory client as there is an unknown configuration value for the PingDirectory username. "+
				"Either target apply the source of the value first or set the value statically in the configuration.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown PingDirectory Password",
			"The provider cannot create the PingDirectory client as there is an unknown configuration value for the PingDirectory password. "+
				"Either target apply the source of the value first or set the value statically in the configuration.",
		)
	}

	if config.DefaultUserPassword.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("default_user_password"),
			"Unknown default PingDirectory user password",
			"The provider cannot create the PingDirectory client as there is an unknown configuration value for the default PingDirectory user password. "+
				"Either target apply the source of the value first or set the value statically in the configuration.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	var host = config.Host.Value
	var username = config.Username.Value
	var password = config.Password.Value
	var defaultUserPassword = config.DefaultUserPassword.Value

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing PingDirectory Host",
			"The provider cannot create the PingDirectory client as there is a missing or empty value for the PingDirectory host. "+
				"Set the host value in the configuration. "+
				"If it is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing PingDirectory Username",
			"The provider cannot create the PingDirectory client as there is a missing or empty value for the PingDirectory username. "+
				"Set the username value in the configuration. "+
				"If it is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing PingDirectory Password",
			"The provider cannot create the PingDirectory client as there is a missing or empty value for the PingDirectory password. "+
				"Set the password value in the configuration. "+
				"If it is already set, ensure the value is not empty.",
		)
	}

	if defaultUserPassword == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("default_user_password"),
			"Missing default PingDirectory user password",
			"The provider cannot create the PingDirectory client as there is a missing or empty value for the default PingDirectory user password. "+
				"Set the default user password value in the configuration. "+
				"If it is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Make the PingDirectory config info available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = config
	resp.ResourceData = config

	tflog.Info(ctx, "Configured PingDirectory client", map[string]interface{}{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *pingdirectoryProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *pingdirectoryProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUsersResource,
	}
}
