// Copyright © 2025 Ping Identity Corporation

package delegatedadminattribute

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &delegatedAdminAttributesDataSource{}
	_ datasource.DataSourceWithConfigure = &delegatedAdminAttributesDataSource{}
)

// Create a Delegated Admin Attributes data source
func NewDelegatedAdminAttributesDataSource() datasource.DataSource {
	return &delegatedAdminAttributesDataSource{}
}

// delegatedAdminAttributesDataSource is the datasource implementation.
type delegatedAdminAttributesDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *delegatedAdminAttributesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_delegated_admin_attributes"
}

// Configure adds the provider configured client to the data source.
func (r *delegatedAdminAttributesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type delegatedAdminAttributesDataSourceModel struct {
	Id                   types.String `tfsdk:"id"`
	Filter               types.String `tfsdk:"filter"`
	Objects              types.Set    `tfsdk:"objects"`
	RestResourceTypeName types.String `tfsdk:"rest_resource_type_name"`
}

// GetSchema defines the schema for the datasource.
func (r *delegatedAdminAttributesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Delegated Admin Attribute objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"rest_resource_type_name": schema.StringAttribute{
				Description: "Name of the parent REST Resource Type",
				Required:    true,
			},
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Delegated Admin Attribute objects found in the configuration",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: internaltypes.ObjectsObjectType(),
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read resource information
func (r *delegatedAdminAttributesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state delegatedAdminAttributesDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.DelegatedAdminAttributeAPI.ListDelegatedAdminAttributes(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.RestResourceTypeName.ValueString())
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.DelegatedAdminAttributeAPI.ListDelegatedAdminAttributesExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Delegated Admin Attribute objects", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	objects := []attr.Value{}
	for _, response := range readResponse.Resources {
		attributes := map[string]attr.Value{}
		if response.CertificateDelegatedAdminAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.CertificateDelegatedAdminAttributeResponse.Id)
			attributes["type"] = types.StringValue("certificate")
		}
		if response.PhotoDelegatedAdminAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.PhotoDelegatedAdminAttributeResponse.Id)
			attributes["type"] = types.StringValue("photo")
		}
		if response.GenericDelegatedAdminAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.GenericDelegatedAdminAttributeResponse.Id)
			attributes["type"] = types.StringValue("generic")
		}
		obj, diags := types.ObjectValue(internaltypes.ObjectsAttrTypes(), attributes)
		resp.Diagnostics.Append(diags...)
		objects = append(objects, obj)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	state.Objects, diags = types.SetValue(internaltypes.ObjectsObjectType(), objects)
	resp.Diagnostics.Append(diags...)
	state.Id = types.StringValue("id")

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
