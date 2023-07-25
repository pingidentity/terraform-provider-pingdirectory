package virtualattribute

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
	_ datasource.DataSource              = &virtualAttributeDataSource{}
	_ datasource.DataSourceWithConfigure = &virtualAttributeDataSource{}
)

// Create a Virtual Attribute data source
func NewVirtualAttributeDataSource() datasource.DataSource {
	return &virtualAttributeDataSource{}
}

// virtualAttributeDataSource is the datasource implementation.
type virtualAttributeDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *virtualAttributeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_attribute"
}

// Configure adds the provider configured client to the data source.
func (r *virtualAttributeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type virtualAttributeDataSourceModel struct {
	Id                                           types.String `tfsdk:"id"`
	Type                                         types.String `tfsdk:"type"`
	ExtensionClass                               types.String `tfsdk:"extension_class"`
	ExtensionArgument                            types.Set    `tfsdk:"extension_argument"`
	ScriptClass                                  types.String `tfsdk:"script_class"`
	AllowRetrievingMembership                    types.Bool   `tfsdk:"allow_retrieving_membership"`
	ScriptArgument                               types.Set    `tfsdk:"script_argument"`
	JoinSourceAttribute                          types.String `tfsdk:"join_source_attribute"`
	JoinTargetAttribute                          types.String `tfsdk:"join_target_attribute"`
	JoinMatchAll                                 types.Bool   `tfsdk:"join_match_all"`
	Value                                        types.Set    `tfsdk:"value"`
	ReferencedByAttribute                        types.Set    `tfsdk:"referenced_by_attribute"`
	ReferenceSearchBaseDN                        types.Set    `tfsdk:"reference_search_base_dn"`
	JoinDNAttribute                              types.String `tfsdk:"join_dn_attribute"`
	JoinBaseDNType                               types.String `tfsdk:"join_base_dn_type"`
	JoinCustomBaseDN                             types.String `tfsdk:"join_custom_base_dn"`
	JoinScope                                    types.String `tfsdk:"join_scope"`
	JoinSizeLimit                                types.Int64  `tfsdk:"join_size_limit"`
	JoinFilter                                   types.String `tfsdk:"join_filter"`
	JoinAttribute                                types.Set    `tfsdk:"join_attribute"`
	ValuePattern                                 types.Set    `tfsdk:"value_pattern"`
	ConflictBehavior                             types.String `tfsdk:"conflict_behavior"`
	DirectMembershipsOnly                        types.Bool   `tfsdk:"direct_memberships_only"`
	IncludedGroupFilter                          types.String `tfsdk:"included_group_filter"`
	RewriteSearchFilters                         types.String `tfsdk:"rewrite_search_filters"`
	SourceAttribute                              types.String `tfsdk:"source_attribute"`
	SourceEntryDNAttribute                       types.String `tfsdk:"source_entry_dn_attribute"`
	SourceEntryDNMap                             types.String `tfsdk:"source_entry_dn_map"`
	BypassAccessControlForSearches               types.Bool   `tfsdk:"bypass_access_control_for_searches"`
	Description                                  types.String `tfsdk:"description"`
	Enabled                                      types.Bool   `tfsdk:"enabled"`
	AttributeType                                types.String `tfsdk:"attribute_type"`
	BaseDN                                       types.Set    `tfsdk:"base_dn"`
	GroupDN                                      types.Set    `tfsdk:"group_dn"`
	Filter                                       types.Set    `tfsdk:"filter"`
	ClientConnectionPolicy                       types.Set    `tfsdk:"client_connection_policy"`
	RequireExplicitRequestByName                 types.Bool   `tfsdk:"require_explicit_request_by_name"`
	MultipleVirtualAttributeEvaluationOrderIndex types.Int64  `tfsdk:"multiple_virtual_attribute_evaluation_order_index"`
	MultipleVirtualAttributeMergeBehavior        types.String `tfsdk:"multiple_virtual_attribute_merge_behavior"`
	AllowIndexConflicts                          types.Bool   `tfsdk:"allow_index_conflicts"`
}

// GetSchema defines the schema for the datasource.
func (r *virtualAttributeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Describes a Virtual Attribute.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Name of this object.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of Virtual Attribute resource. Options are ['mirror', 'entry-checksum', 'member-of-server-group', 'constructed', 'is-member-of', 'custom', 'num-subordinates', 'reverse-dn-join', 'identify-references', 'user-defined', 'current-time', 'short-unique-id', 'entry-dn', 'has-subordinates', 'equality-join', 'groovy-scripted', 'instance-name', 'replication-state-detail', 'member', 'password-policy-state-json', 'subschema-subentry', 'dn-join', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Virtual Attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Virtual Attribute. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Virtual Attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_retrieving_membership": schema.BoolAttribute{
				Description: "Indicates whether to handle requests that request all values for the virtual attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Virtual Attribute. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"join_source_attribute": schema.StringAttribute{
				Description: "The attribute containing the value(s) in the source entry to use to identify related entries.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"join_target_attribute": schema.StringAttribute{
				Description: "The attribute in target entries whose value(s) match values of the source attribute in the source entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"join_match_all": schema.BoolAttribute{
				Description: "Indicates whether joined entries will be required to have all values for the source attribute, or only at least one of its values.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"value": schema.SetAttribute{
				Description: "Specifies the values to be included in the virtual attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"referenced_by_attribute": schema.SetAttribute{
				Description: "The name or OID of an attribute type whose values will be searched for references to the target entry. The attribute type must be defined in the server schema, must have a syntax of either \"distinguished name\" or \"name and optional UID\", and must be indexed for equality.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"reference_search_base_dn": schema.SetAttribute{
				Description: "The base DN that will be used when searching for references to the target entry. If no reference search base DN is specified, the default behavior will be to search below all public naming contexts configured in the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"join_dn_attribute": schema.StringAttribute{
				Description: "The attribute in related entries whose set of values must contain the DN of the search result entry to be joined with that entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"join_base_dn_type": schema.StringAttribute{
				Description: "Specifies how server should determine the base DN for the internal searches used to identify joined entries.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"join_custom_base_dn": schema.StringAttribute{
				Description: "The fixed, administrator-specified base DN for the internal searches used to identify joined entries.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"join_scope": schema.StringAttribute{
				Description: "The scope for searches used to identify joined entries.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"join_size_limit": schema.Int64Attribute{
				Description: "The maximum number of entries that may be joined with the source entry, which also corresponds to the maximum number of values that the virtual attribute provider will generate for an entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"join_filter": schema.StringAttribute{
				Description: "An optional filter that specifies additional criteria for identifying joined entries. If a join-filter value is specified, then only entries matching that filter (in addition to satisfying the other join criteria) will be joined with the search result entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"join_attribute": schema.SetAttribute{
				Description: "An optional set of the names of the attributes to include with joined entries.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"value_pattern": schema.SetAttribute{
				Description: "Specifies a pattern for constructing the virtual attribute value using fixed text and attribute values from the entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"conflict_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the server is to exhibit for entries that already contain one or more real values for the associated attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"direct_memberships_only": schema.BoolAttribute{
				Description: "Specifies whether to only include groups in which the user is directly associated with and the membership maybe modified via the group entry. Groups in which the user's membership is derived dynamically or through nested groups will not be included.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"included_group_filter": schema.StringAttribute{
				Description: "A search filter that will be used to identify which groups should be included in the values of the virtual attribute. With no value defined (which is the default behavior), all groups that contain the target user will be included.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"rewrite_search_filters": schema.StringAttribute{
				Description: "Search filters that include Is Member Of Virtual Attribute searches on dynamic groups can be updated to include the dynamic group filter in the search filter itself. This can allow the backend to more efficiently process the search filter by using attribute indexes sooner in the search processing.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"source_attribute": schema.StringAttribute{
				Description: "Specifies the source attribute containing the values to use for this virtual attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"source_entry_dn_attribute": schema.StringAttribute{
				Description: "Specifies the attribute containing the DN of another entry from which to obtain the source attribute providing the values for this virtual attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"source_entry_dn_map": schema.StringAttribute{
				Description: "Specifies a DN map that will be used to identify the entry from which to obtain the source attribute providing the values for this virtual attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"bypass_access_control_for_searches": schema.BoolAttribute{
				Description: "Indicates whether searches performed by this virtual attribute provider should be exempted from access control restrictions.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Virtual Attribute",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Virtual Attribute is enabled for use.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"attribute_type": schema.StringAttribute{
				Description: "Specifies the attribute type for the attribute whose values are to be dynamically assigned by the virtual attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"base_dn": schema.SetAttribute{
				Description: "Specifies the base DNs for the branches containing entries that are eligible to use this virtual attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"group_dn": schema.SetAttribute{
				Description: "Specifies the DNs of the groups whose members can be eligible to use this virtual attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"filter": schema.SetAttribute{
				Description: "Specifies the search filters to be applied against entries to determine if the virtual attribute is to be generated for those entries.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"client_connection_policy": schema.SetAttribute{
				Description: "Specifies a set of client connection policies for which this Virtual Attribute should be generated. If this is undefined, then this Virtual Attribute will always be generated. If it is associated with one or more client connection policies, then this Virtual Attribute will be generated only for operations requested by clients assigned to one of those client connection policies.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"require_explicit_request_by_name": schema.BoolAttribute{
				Description: "Indicates whether attributes of this type must be explicitly included by name in the list of requested attributes. Note that this will only apply to virtual attributes which are associated with an attribute type that is operational. It will be ignored for virtual attributes associated with a non-operational attribute type.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"multiple_virtual_attribute_evaluation_order_index": schema.Int64Attribute{
				Description: "Specifies the order in which virtual attribute definitions for the same attribute type will be evaluated when generating values for an entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"multiple_virtual_attribute_merge_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that will be exhibited for cases in which multiple virtual attribute definitions apply to the same multivalued attribute type. This will be ignored for single-valued attribute types.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_index_conflicts": schema.BoolAttribute{
				Description: "Indicates whether the server should allow creating or altering this virtual attribute definition even if it conflicts with one or more indexes defined in the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
}

// Read a MirrorVirtualAttributeResponse object into the model struct
func readMirrorVirtualAttributeResponseDataSource(ctx context.Context, r *client.MirrorVirtualAttributeResponse, state *virtualAttributeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("mirror")
	state.Id = types.StringValue(r.Id)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), false)
	state.SourceAttribute = types.StringValue(r.SourceAttribute)
	state.SourceEntryDNAttribute = internaltypes.StringTypeOrNil(r.SourceEntryDNAttribute, false)
	state.SourceEntryDNMap = internaltypes.StringTypeOrNil(r.SourceEntryDNMap, false)
	state.BypassAccessControlForSearches = internaltypes.BoolTypeOrNil(r.BypassAccessControlForSearches)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), false)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
}

// Read a ConstructedVirtualAttributeResponse object into the model struct
func readConstructedVirtualAttributeResponseDataSource(ctx context.Context, r *client.ConstructedVirtualAttributeResponse, state *virtualAttributeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("constructed")
	state.Id = types.StringValue(r.Id)
	state.ValuePattern = internaltypes.GetStringSet(r.ValuePattern)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), false)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), false)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
}

// Read a IsMemberOfVirtualAttributeResponse object into the model struct
func readIsMemberOfVirtualAttributeResponseDataSource(ctx context.Context, r *client.IsMemberOfVirtualAttributeResponse, state *virtualAttributeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("is-member-of")
	state.Id = types.StringValue(r.Id)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), false)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.DirectMembershipsOnly = internaltypes.BoolTypeOrNil(r.DirectMembershipsOnly)
	state.IncludedGroupFilter = internaltypes.StringTypeOrNil(r.IncludedGroupFilter, false)
	state.RewriteSearchFilters = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeRewriteSearchFiltersProp(r.RewriteSearchFilters), false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), false)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
}

// Read a ReverseDnJoinVirtualAttributeResponse object into the model struct
func readReverseDnJoinVirtualAttributeResponseDataSource(ctx context.Context, r *client.ReverseDnJoinVirtualAttributeResponse, state *virtualAttributeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("reverse-dn-join")
	state.Id = types.StringValue(r.Id)
	state.JoinDNAttribute = types.StringValue(r.JoinDNAttribute)
	state.JoinBaseDNType = types.StringValue(r.JoinBaseDNType.String())
	state.JoinCustomBaseDN = internaltypes.StringTypeOrNil(r.JoinCustomBaseDN, false)
	state.JoinScope = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeJoinScopeProp(r.JoinScope), false)
	state.JoinSizeLimit = internaltypes.Int64TypeOrNil(r.JoinSizeLimit)
	state.JoinFilter = internaltypes.StringTypeOrNil(r.JoinFilter, false)
	state.JoinAttribute = internaltypes.GetStringSet(r.JoinAttribute)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), false)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), false)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
}

// Read a IdentifyReferencesVirtualAttributeResponse object into the model struct
func readIdentifyReferencesVirtualAttributeResponseDataSource(ctx context.Context, r *client.IdentifyReferencesVirtualAttributeResponse, state *virtualAttributeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("identify-references")
	state.Id = types.StringValue(r.Id)
	state.ReferencedByAttribute = internaltypes.GetStringSet(r.ReferencedByAttribute)
	state.ReferenceSearchBaseDN = internaltypes.GetStringSet(r.ReferenceSearchBaseDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), false)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), false)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
}

// Read a UserDefinedVirtualAttributeResponse object into the model struct
func readUserDefinedVirtualAttributeResponseDataSource(ctx context.Context, r *client.UserDefinedVirtualAttributeResponse, state *virtualAttributeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("user-defined")
	state.Id = types.StringValue(r.Id)
	state.Value = internaltypes.GetStringSet(r.Value)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), false)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), false)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
}

// Read a EntryDnVirtualAttributeResponse object into the model struct
func readEntryDnVirtualAttributeResponseDataSource(ctx context.Context, r *client.EntryDnVirtualAttributeResponse, state *virtualAttributeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("entry-dn")
	state.Id = types.StringValue(r.Id)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), false)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), false)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
}

// Read a EqualityJoinVirtualAttributeResponse object into the model struct
func readEqualityJoinVirtualAttributeResponseDataSource(ctx context.Context, r *client.EqualityJoinVirtualAttributeResponse, state *virtualAttributeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("equality-join")
	state.Id = types.StringValue(r.Id)
	state.JoinSourceAttribute = types.StringValue(r.JoinSourceAttribute)
	state.JoinTargetAttribute = types.StringValue(r.JoinTargetAttribute)
	state.JoinMatchAll = internaltypes.BoolTypeOrNil(r.JoinMatchAll)
	state.JoinBaseDNType = types.StringValue(r.JoinBaseDNType.String())
	state.JoinCustomBaseDN = internaltypes.StringTypeOrNil(r.JoinCustomBaseDN, false)
	state.JoinScope = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeJoinScopeProp(r.JoinScope), false)
	state.JoinSizeLimit = internaltypes.Int64TypeOrNil(r.JoinSizeLimit)
	state.JoinFilter = internaltypes.StringTypeOrNil(r.JoinFilter, false)
	state.JoinAttribute = internaltypes.GetStringSet(r.JoinAttribute)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), false)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), false)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
}

// Read a GroovyScriptedVirtualAttributeResponse object into the model struct
func readGroovyScriptedVirtualAttributeResponseDataSource(ctx context.Context, r *client.GroovyScriptedVirtualAttributeResponse, state *virtualAttributeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), false)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), false)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
}

// Read a MemberVirtualAttributeResponse object into the model struct
func readMemberVirtualAttributeResponseDataSource(ctx context.Context, r *client.MemberVirtualAttributeResponse, state *virtualAttributeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("member")
	state.Id = types.StringValue(r.Id)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), false)
	state.AllowRetrievingMembership = types.BoolValue(r.AllowRetrievingMembership)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), false)
}

// Read a PasswordPolicyStateJsonVirtualAttributeResponse object into the model struct
func readPasswordPolicyStateJsonVirtualAttributeResponseDataSource(ctx context.Context, r *client.PasswordPolicyStateJsonVirtualAttributeResponse, state *virtualAttributeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("password-policy-state-json")
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
}

// Read a DnJoinVirtualAttributeResponse object into the model struct
func readDnJoinVirtualAttributeResponseDataSource(ctx context.Context, r *client.DnJoinVirtualAttributeResponse, state *virtualAttributeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("dn-join")
	state.Id = types.StringValue(r.Id)
	state.JoinDNAttribute = types.StringValue(r.JoinDNAttribute)
	state.JoinBaseDNType = types.StringValue(r.JoinBaseDNType.String())
	state.JoinCustomBaseDN = internaltypes.StringTypeOrNil(r.JoinCustomBaseDN, false)
	state.JoinScope = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeJoinScopeProp(r.JoinScope), false)
	state.JoinSizeLimit = internaltypes.Int64TypeOrNil(r.JoinSizeLimit)
	state.JoinFilter = internaltypes.StringTypeOrNil(r.JoinFilter, false)
	state.JoinAttribute = internaltypes.GetStringSet(r.JoinAttribute)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), false)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), false)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
}

// Read a ThirdPartyVirtualAttributeResponse object into the model struct
func readThirdPartyVirtualAttributeResponseDataSource(ctx context.Context, r *client.ThirdPartyVirtualAttributeResponse, state *virtualAttributeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), false)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), false)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
}

// Read resource information
func (r *virtualAttributeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state virtualAttributeDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.VirtualAttributeApi.GetVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Virtual Attribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.MirrorVirtualAttributeResponse != nil {
		readMirrorVirtualAttributeResponseDataSource(ctx, readResponse.MirrorVirtualAttributeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ConstructedVirtualAttributeResponse != nil {
		readConstructedVirtualAttributeResponseDataSource(ctx, readResponse.ConstructedVirtualAttributeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.IsMemberOfVirtualAttributeResponse != nil {
		readIsMemberOfVirtualAttributeResponseDataSource(ctx, readResponse.IsMemberOfVirtualAttributeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ReverseDnJoinVirtualAttributeResponse != nil {
		readReverseDnJoinVirtualAttributeResponseDataSource(ctx, readResponse.ReverseDnJoinVirtualAttributeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.IdentifyReferencesVirtualAttributeResponse != nil {
		readIdentifyReferencesVirtualAttributeResponseDataSource(ctx, readResponse.IdentifyReferencesVirtualAttributeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.UserDefinedVirtualAttributeResponse != nil {
		readUserDefinedVirtualAttributeResponseDataSource(ctx, readResponse.UserDefinedVirtualAttributeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.EntryDnVirtualAttributeResponse != nil {
		readEntryDnVirtualAttributeResponseDataSource(ctx, readResponse.EntryDnVirtualAttributeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.EqualityJoinVirtualAttributeResponse != nil {
		readEqualityJoinVirtualAttributeResponseDataSource(ctx, readResponse.EqualityJoinVirtualAttributeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedVirtualAttributeResponse != nil {
		readGroovyScriptedVirtualAttributeResponseDataSource(ctx, readResponse.GroovyScriptedVirtualAttributeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.MemberVirtualAttributeResponse != nil {
		readMemberVirtualAttributeResponseDataSource(ctx, readResponse.MemberVirtualAttributeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.PasswordPolicyStateJsonVirtualAttributeResponse != nil {
		readPasswordPolicyStateJsonVirtualAttributeResponseDataSource(ctx, readResponse.PasswordPolicyStateJsonVirtualAttributeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DnJoinVirtualAttributeResponse != nil {
		readDnJoinVirtualAttributeResponseDataSource(ctx, readResponse.DnJoinVirtualAttributeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyVirtualAttributeResponse != nil {
		readThirdPartyVirtualAttributeResponseDataSource(ctx, readResponse.ThirdPartyVirtualAttributeResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
