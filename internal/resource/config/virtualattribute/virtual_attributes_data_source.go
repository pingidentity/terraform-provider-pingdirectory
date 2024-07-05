package virtualattribute

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10100/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &virtualAttributesDataSource{}
	_ datasource.DataSourceWithConfigure = &virtualAttributesDataSource{}
)

// Create a Virtual Attributes data source
func NewVirtualAttributesDataSource() datasource.DataSource {
	return &virtualAttributesDataSource{}
}

// virtualAttributesDataSource is the datasource implementation.
type virtualAttributesDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *virtualAttributesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_attributes"
}

// Configure adds the provider configured client to the data source.
func (r *virtualAttributesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type virtualAttributesDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *virtualAttributesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Virtual Attribute objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Virtual Attribute objects found in the configuration",
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
func (r *virtualAttributesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state virtualAttributesDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.VirtualAttributeAPI.ListVirtualAttributes(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.VirtualAttributeAPI.ListVirtualAttributesExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Virtual Attribute objects", err, httpResp)
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
		if response.MirrorVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.MirrorVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("mirror")
		}
		if response.EntryChecksumVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.EntryChecksumVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("entry-checksum")
		}
		if response.IsMemberOfVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.IsMemberOfVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("is-member-of")
		}
		if response.ReverseDnJoinVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.ReverseDnJoinVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("reverse-dn-join")
		}
		if response.IdentifyReferencesVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.IdentifyReferencesVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("identify-references")
		}
		if response.UserDefinedVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.UserDefinedVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("user-defined")
		}
		if response.ShortUniqueIdVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.ShortUniqueIdVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("short-unique-id")
		}
		if response.ExpandTimestampVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.ExpandTimestampVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("expand-timestamp")
		}
		if response.InstanceNameVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.InstanceNameVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("instance-name")
		}
		if response.MemberVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.MemberVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("member")
		}
		if response.PasswordPolicyStateJsonVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.PasswordPolicyStateJsonVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("password-policy-state-json")
		}
		if response.SubschemaSubentryVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.SubschemaSubentryVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("subschema-subentry")
		}
		if response.DnJoinVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.DnJoinVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("dn-join")
		}
		if response.LargeAttributeVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.LargeAttributeVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("large-attribute")
		}
		if response.ThirdPartyVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("third-party")
		}
		if response.MemberOfServerGroupVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.MemberOfServerGroupVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("member-of-server-group")
		}
		if response.ConstructedVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.ConstructedVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("constructed")
		}
		if response.CustomVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.CustomVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("custom")
		}
		if response.FileBasedVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.FileBasedVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("file-based")
		}
		if response.NumSubordinatesVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.NumSubordinatesVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("num-subordinates")
		}
		if response.CurrentTimeVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.CurrentTimeVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("current-time")
		}
		if response.EntryUuidVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.EntryUuidVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("entry-uuid")
		}
		if response.EntryDnVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.EntryDnVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("entry-dn")
		}
		if response.HasSubordinatesVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.HasSubordinatesVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("has-subordinates")
		}
		if response.ConfigModelVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.ConfigModelVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("config-model")
		}
		if response.EqualityJoinVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.EqualityJoinVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("equality-join")
		}
		if response.GroovyScriptedVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.GroovyScriptedVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("groovy-scripted")
		}
		if response.ReplicationStateDetailVirtualAttributeResponse != nil {
			attributes["id"] = types.StringValue(response.ReplicationStateDetailVirtualAttributeResponse.Id)
			attributes["type"] = types.StringValue("replication-state-detail")
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
