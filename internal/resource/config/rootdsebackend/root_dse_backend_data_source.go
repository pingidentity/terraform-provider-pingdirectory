// Copyright © 2025 Ping Identity Corporation

package rootdsebackend

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &rootDseBackendDataSource{}
	_ datasource.DataSourceWithConfigure = &rootDseBackendDataSource{}
)

// Create a Root Dse Backend data source
func NewRootDseBackendDataSource() datasource.DataSource {
	return &rootDseBackendDataSource{}
}

// rootDseBackendDataSource is the datasource implementation.
type rootDseBackendDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *rootDseBackendDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_root_dse_backend"
}

// Configure adds the provider configured client to the data source.
func (r *rootDseBackendDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type rootDseBackendDataSourceModel struct {
	Id                            types.String `tfsdk:"id"`
	Type                          types.String `tfsdk:"type"`
	SubordinateBaseDN             types.Set    `tfsdk:"subordinate_base_dn"`
	AdditionalSupportedControlOID types.Set    `tfsdk:"additional_supported_control_oid"`
	ShowAllAttributes             types.Bool   `tfsdk:"show_all_attributes"`
	UseLegacyVendorVersion        types.Bool   `tfsdk:"use_legacy_vendor_version"`
}

// GetSchema defines the schema for the datasource.
func (r *rootDseBackendDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Root Dse Backend.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Root DSE Backend resource. Options are ['root-dse-backend']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"subordinate_base_dn": schema.SetAttribute{
				Description: "Specifies the set of base DNs used for singleLevel, wholeSubtree, and subordinateSubtree searches based at the root DSE.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"additional_supported_control_oid": schema.SetAttribute{
				Description: "Specifies an additional OID that should appear in the list of supportedControl values in the server's root DSE.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"show_all_attributes": schema.BoolAttribute{
				Description: "Indicates whether all attributes in the root DSE are to be treated like user attributes (and therefore returned to clients by default) regardless of the Directory Server schema configuration.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"use_legacy_vendor_version": schema.BoolAttribute{
				Description: "Indicates whether the server's root DSE should reflect current or legacy values for the vendorName and vendorVersion attributes.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a RootDseBackendResponse object into the model struct
func readRootDseBackendResponseDataSource(ctx context.Context, r *client.RootDseBackendResponse, state *rootDseBackendDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("root-dse-backend")
	// Placeholder id value required by test framework
	state.Id = types.StringValue("id")
	state.SubordinateBaseDN = internaltypes.GetStringSet(r.SubordinateBaseDN)
	state.AdditionalSupportedControlOID = internaltypes.GetStringSet(r.AdditionalSupportedControlOID)
	state.ShowAllAttributes = types.BoolValue(r.ShowAllAttributes)
	state.UseLegacyVendorVersion = internaltypes.BoolTypeOrNil(r.UseLegacyVendorVersion)
}

// Read resource information
func (r *rootDseBackendDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state rootDseBackendDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.RootDseBackendAPI.GetRootDseBackend(
		config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Root Dse Backend", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readRootDseBackendResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
