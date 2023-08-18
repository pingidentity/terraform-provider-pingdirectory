package interserverauthenticationinfo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &interServerAuthenticationInfoDataSource{}
	_ datasource.DataSourceWithConfigure = &interServerAuthenticationInfoDataSource{}
)

// Create a Inter Server Authentication Info data source
func NewInterServerAuthenticationInfoDataSource() datasource.DataSource {
	return &interServerAuthenticationInfoDataSource{}
}

// interServerAuthenticationInfoDataSource is the datasource implementation.
type interServerAuthenticationInfoDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *interServerAuthenticationInfoDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inter_server_authentication_info"
}

// Configure adds the provider configured client to the data source.
func (r *interServerAuthenticationInfoDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type interServerAuthenticationInfoDataSourceModel struct {
	Id                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	Type                       types.String `tfsdk:"type"`
	ServerInstanceListenerName types.String `tfsdk:"server_instance_listener_name"`
	ServerInstanceName         types.String `tfsdk:"server_instance_name"`
	AuthenticationType         types.String `tfsdk:"authentication_type"`
	BindDN                     types.String `tfsdk:"bind_dn"`
	Username                   types.String `tfsdk:"username"`
	Password                   types.String `tfsdk:"password"`
	Purpose                    types.Set    `tfsdk:"purpose"`
}

// GetSchema defines the schema for the datasource.
func (r *interServerAuthenticationInfoDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Inter Server Authentication Info.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Inter Server Authentication Info resource. Options are ['password', 'certificate']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_instance_listener_name": schema.StringAttribute{
				Description: "Name of the parent Server Instance Listener",
				Required:    true,
			},
			"server_instance_name": schema.StringAttribute{
				Description: "Name of the parent Server Instance",
				Required:    true,
			},
			"authentication_type": schema.StringAttribute{
				Description: "Identifies the type of password authentication that will be used.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"bind_dn": schema.StringAttribute{
				Description: "A DN of the username that should be used for the bind request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "The username that should be used for the bind request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password for the username or bind-dn.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"purpose": schema.SetAttribute{
				Description: "Identifies the purpose of this Inter Server Authentication Info.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a PasswordInterServerAuthenticationInfoResponse object into the model struct
func readPasswordInterServerAuthenticationInfoResponseDataSource(ctx context.Context, r *client.PasswordInterServerAuthenticationInfoResponse, state *interServerAuthenticationInfoDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("password")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AuthenticationType = internaltypes.StringTypeOrNil(
		client.StringPointerEnuminterServerAuthenticationInfoAuthenticationTypeProp(r.AuthenticationType), false)
	state.BindDN = internaltypes.StringTypeOrNil(r.BindDN, false)
	state.Username = internaltypes.StringTypeOrNil(r.Username, false)
	state.Purpose = internaltypes.GetStringSet(
		client.StringSliceEnuminterServerAuthenticationInfoPurposeProp(r.Purpose))
}

// Read a CertificateInterServerAuthenticationInfoResponse object into the model struct
func readCertificateInterServerAuthenticationInfoResponseDataSource(ctx context.Context, r *client.CertificateInterServerAuthenticationInfoResponse, state *interServerAuthenticationInfoDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("certificate")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Purpose = internaltypes.GetStringSet(
		client.StringSliceEnuminterServerAuthenticationInfoPurposeProp(r.Purpose))
}

// Read resource information
func (r *interServerAuthenticationInfoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state interServerAuthenticationInfoDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.InterServerAuthenticationInfoApi.GetInterServerAuthenticationInfo(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.ServerInstanceListenerName.ValueString(), state.ServerInstanceName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Inter Server Authentication Info", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.PasswordInterServerAuthenticationInfoResponse != nil {
		readPasswordInterServerAuthenticationInfoResponseDataSource(ctx, readResponse.PasswordInterServerAuthenticationInfoResponse, &state, &resp.Diagnostics)
	}
	if readResponse.CertificateInterServerAuthenticationInfoResponse != nil {
		readCertificateInterServerAuthenticationInfoResponseDataSource(ctx, readResponse.CertificateInterServerAuthenticationInfoResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
