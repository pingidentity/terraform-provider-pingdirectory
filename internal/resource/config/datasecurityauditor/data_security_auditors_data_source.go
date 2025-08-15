// Copyright © 2025 Ping Identity Corporation

package datasecurityauditor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dataSecurityAuditorsDataSource{}
	_ datasource.DataSourceWithConfigure = &dataSecurityAuditorsDataSource{}
)

// Create a Data Security Auditors data source
func NewDataSecurityAuditorsDataSource() datasource.DataSource {
	return &dataSecurityAuditorsDataSource{}
}

// dataSecurityAuditorsDataSource is the datasource implementation.
type dataSecurityAuditorsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *dataSecurityAuditorsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_security_auditors"
}

// Configure adds the provider configured client to the data source.
func (r *dataSecurityAuditorsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type dataSecurityAuditorsDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *dataSecurityAuditorsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Data Security Auditor objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Data Security Auditor objects found in the configuration",
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
func (r *dataSecurityAuditorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state dataSecurityAuditorsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.DataSecurityAuditorAPI.ListDataSecurityAuditors(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.DataSecurityAuditorAPI.ListDataSecurityAuditorsExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Data Security Auditor objects", err, httpResp)
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
		if response.ExpiredPasswordDataSecurityAuditorResponse != nil {
			attributes["id"] = types.StringValue(response.ExpiredPasswordDataSecurityAuditorResponse.Id)
			attributes["type"] = types.StringValue("expired-password")
		}
		if response.IdleAccountDataSecurityAuditorResponse != nil {
			attributes["id"] = types.StringValue(response.IdleAccountDataSecurityAuditorResponse.Id)
			attributes["type"] = types.StringValue("idle-account")
		}
		if response.DisabledAccountDataSecurityAuditorResponse != nil {
			attributes["id"] = types.StringValue(response.DisabledAccountDataSecurityAuditorResponse.Id)
			attributes["type"] = types.StringValue("disabled-account")
		}
		if response.WeaklyEncodedPasswordDataSecurityAuditorResponse != nil {
			attributes["id"] = types.StringValue(response.WeaklyEncodedPasswordDataSecurityAuditorResponse.Id)
			attributes["type"] = types.StringValue("weakly-encoded-password")
		}
		if response.PrivilegeDataSecurityAuditorResponse != nil {
			attributes["id"] = types.StringValue(response.PrivilegeDataSecurityAuditorResponse.Id)
			attributes["type"] = types.StringValue("privilege")
		}
		if response.AccountUsabilityIssuesDataSecurityAuditorResponse != nil {
			attributes["id"] = types.StringValue(response.AccountUsabilityIssuesDataSecurityAuditorResponse.Id)
			attributes["type"] = types.StringValue("account-usability-issues")
		}
		if response.LockedAccountDataSecurityAuditorResponse != nil {
			attributes["id"] = types.StringValue(response.LockedAccountDataSecurityAuditorResponse.Id)
			attributes["type"] = types.StringValue("locked-account")
		}
		if response.FilterDataSecurityAuditorResponse != nil {
			attributes["id"] = types.StringValue(response.FilterDataSecurityAuditorResponse.Id)
			attributes["type"] = types.StringValue("filter")
		}
		if response.AccountValidityWindowDataSecurityAuditorResponse != nil {
			attributes["id"] = types.StringValue(response.AccountValidityWindowDataSecurityAuditorResponse.Id)
			attributes["type"] = types.StringValue("account-validity-window")
		}
		if response.MultiplePasswordDataSecurityAuditorResponse != nil {
			attributes["id"] = types.StringValue(response.MultiplePasswordDataSecurityAuditorResponse.Id)
			attributes["type"] = types.StringValue("multiple-password")
		}
		if response.DeprecatedPasswordStorageSchemeDataSecurityAuditorResponse != nil {
			attributes["id"] = types.StringValue(response.DeprecatedPasswordStorageSchemeDataSecurityAuditorResponse.Id)
			attributes["type"] = types.StringValue("deprecated-password-storage-scheme")
		}
		if response.NonexistentPasswordPolicyDataSecurityAuditorResponse != nil {
			attributes["id"] = types.StringValue(response.NonexistentPasswordPolicyDataSecurityAuditorResponse.Id)
			attributes["type"] = types.StringValue("nonexistent-password-policy")
		}
		if response.AccessControlDataSecurityAuditorResponse != nil {
			attributes["id"] = types.StringValue(response.AccessControlDataSecurityAuditorResponse.Id)
			attributes["type"] = types.StringValue("access-control")
		}
		if response.ThirdPartyDataSecurityAuditorResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyDataSecurityAuditorResponse.Id)
			attributes["type"] = types.StringValue("third-party")
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
