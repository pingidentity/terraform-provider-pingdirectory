package accesscontrolhandler

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
	_ datasource.DataSource              = &accessControlHandlerDataSource{}
	_ datasource.DataSourceWithConfigure = &accessControlHandlerDataSource{}
)

// Create a Access Control Handler data source
func NewAccessControlHandlerDataSource() datasource.DataSource {
	return &accessControlHandlerDataSource{}
}

// accessControlHandlerDataSource is the datasource implementation.
type accessControlHandlerDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *accessControlHandlerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_control_handler"
}

// Configure adds the provider configured client to the data source.
func (r *accessControlHandlerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type accessControlHandlerDataSourceModel struct {
	// Id field required for acceptance testing framework
	Id                    types.String `tfsdk:"id"`
	GlobalACI             types.Set    `tfsdk:"global_aci"`
	AllowedBindControl    types.Set    `tfsdk:"allowed_bind_control"`
	AllowedBindControlOID types.Set    `tfsdk:"allowed_bind_control_oid"`
	Enabled               types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *accessControlHandlerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Describes a Access Control Handler.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Name of this object.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"global_aci": schema.SetAttribute{
				Description: "Defines global access control rules.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"allowed_bind_control": schema.SetAttribute{
				Description: "Specifies a set of controls that clients should be allowed to include in bind requests. As bind requests are evaluated as the unauthenticated user, any controls included in this set will be permitted for any bind attempt. If you wish to grant permission for any bind controls not listed here, then the allowed-bind-control-oid property may be used to accomplish that.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"allowed_bind_control_oid": schema.SetAttribute{
				Description: "Specifies the OIDs of any additional controls (not covered by the allowed-bind-control property) that should be permitted in bind requests.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Access Control Handler is enabled. If set to FALSE, then no access control is enforced, and any client (including unauthenticated or anonymous clients) could be allowed to perform any operation if not subject to other restrictions, such as those enforced by the privilege subsystem.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
}

// Read a DseeCompatAccessControlHandlerResponse object into the model struct
func readDseeCompatAccessControlHandlerResponseDataSource(ctx context.Context, r *client.DseeCompatAccessControlHandlerResponse, state *accessControlHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	// Placeholder id value required by test framework
	state.Id = types.StringValue("id")
	state.GlobalACI = internaltypes.GetStringSet(r.GlobalACI)
	state.AllowedBindControl = internaltypes.GetStringSet(
		client.StringSliceEnumaccessControlHandlerAllowedBindControlProp(r.AllowedBindControl))
	state.AllowedBindControlOID = internaltypes.GetStringSet(r.AllowedBindControlOID)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read resource information
func (r *accessControlHandlerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state accessControlHandlerDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AccessControlHandlerApi.GetAccessControlHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Access Control Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDseeCompatAccessControlHandlerResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
