package virtualattribute

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &virtualAttributeResource{}
	_ resource.ResourceWithConfigure   = &virtualAttributeResource{}
	_ resource.ResourceWithImportState = &virtualAttributeResource{}
	_ resource.Resource                = &defaultVirtualAttributeResource{}
	_ resource.ResourceWithConfigure   = &defaultVirtualAttributeResource{}
	_ resource.ResourceWithImportState = &defaultVirtualAttributeResource{}
)

// Create a Virtual Attribute resource
func NewVirtualAttributeResource() resource.Resource {
	return &virtualAttributeResource{}
}

func NewDefaultVirtualAttributeResource() resource.Resource {
	return &defaultVirtualAttributeResource{}
}

// virtualAttributeResource is the resource implementation.
type virtualAttributeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultVirtualAttributeResource is the resource implementation.
type defaultVirtualAttributeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *virtualAttributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_attribute"
}

func (r *defaultVirtualAttributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_virtual_attribute"
}

// Configure adds the provider configured client to the resource.
func (r *virtualAttributeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultVirtualAttributeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type virtualAttributeResourceModel struct {
	Id                                           types.String `tfsdk:"id"`
	Name                                         types.String `tfsdk:"name"`
	LastUpdated                                  types.String `tfsdk:"last_updated"`
	Notifications                                types.Set    `tfsdk:"notifications"`
	RequiredActions                              types.Set    `tfsdk:"required_actions"`
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

type defaultVirtualAttributeResourceModel struct {
	Id                                           types.String `tfsdk:"id"`
	Name                                         types.String `tfsdk:"name"`
	LastUpdated                                  types.String `tfsdk:"last_updated"`
	Notifications                                types.Set    `tfsdk:"notifications"`
	RequiredActions                              types.Set    `tfsdk:"required_actions"`
	Type                                         types.String `tfsdk:"type"`
	ExtensionClass                               types.String `tfsdk:"extension_class"`
	ExtensionArgument                            types.Set    `tfsdk:"extension_argument"`
	ScriptClass                                  types.String `tfsdk:"script_class"`
	AllowRetrievingMembership                    types.Bool   `tfsdk:"allow_retrieving_membership"`
	ScriptArgument                               types.Set    `tfsdk:"script_argument"`
	JoinSourceAttribute                          types.String `tfsdk:"join_source_attribute"`
	JoinTargetAttribute                          types.String `tfsdk:"join_target_attribute"`
	JoinMatchAll                                 types.Bool   `tfsdk:"join_match_all"`
	SequenceNumberAttribute                      types.String `tfsdk:"sequence_number_attribute"`
	Value                                        types.Set    `tfsdk:"value"`
	ReferencedByAttribute                        types.Set    `tfsdk:"referenced_by_attribute"`
	ReturnUtcTime                                types.Bool   `tfsdk:"return_utc_time"`
	IncludeMilliseconds                          types.Bool   `tfsdk:"include_milliseconds"`
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
	ExcludeOperationalAttributes                 types.Bool   `tfsdk:"exclude_operational_attributes"`
	ExcludedAttribute                            types.Set    `tfsdk:"excluded_attribute"`
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

// GetSchema defines the schema for the resource.
func (r *virtualAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	virtualAttributeSchema(ctx, req, resp, false)
}

func (r *defaultVirtualAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	virtualAttributeSchema(ctx, req, resp, true)
}

func virtualAttributeSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Virtual Attribute.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Virtual Attribute resource. Options are ['mirror', 'entry-checksum', 'member-of-server-group', 'constructed', 'is-member-of', 'custom', 'num-subordinates', 'reverse-dn-join', 'identify-references', 'user-defined', 'current-time', 'short-unique-id', 'entry-dn', 'has-subordinates', 'equality-join', 'groovy-scripted', 'instance-name', 'replication-state-detail', 'member', 'password-policy-state-json', 'subschema-subentry', 'dn-join', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"mirror", "constructed", "is-member-of", "reverse-dn-join", "identify-references", "user-defined", "entry-dn", "equality-join", "groovy-scripted", "member", "password-policy-state-json", "dn-join", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Virtual Attribute.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Virtual Attribute. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Virtual Attribute.",
				Optional:    true,
			},
			"allow_retrieving_membership": schema.BoolAttribute{
				Description: "Indicates whether to handle requests that request all values for the virtual attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Virtual Attribute. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"join_source_attribute": schema.StringAttribute{
				Description: "The attribute containing the value(s) in the source entry to use to identify related entries.",
				Optional:    true,
			},
			"join_target_attribute": schema.StringAttribute{
				Description: "The attribute in target entries whose value(s) match values of the source attribute in the source entry.",
				Optional:    true,
			},
			"join_match_all": schema.BoolAttribute{
				Description: "Indicates whether joined entries will be required to have all values for the source attribute, or only at least one of its values.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"value": schema.SetAttribute{
				Description: "Specifies the values to be included in the virtual attribute.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"referenced_by_attribute": schema.SetAttribute{
				Description: "The name or OID of an attribute type whose values will be searched for references to the target entry. The attribute type must be defined in the server schema, must have a syntax of either \"distinguished name\" or \"name and optional UID\", and must be indexed for equality.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"reference_search_base_dn": schema.SetAttribute{
				Description: "The base DN that will be used when searching for references to the target entry. If no reference search base DN is specified, the default behavior will be to search below all public naming contexts configured in the server.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"join_dn_attribute": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `reverse-dn-join`: The attribute in related entries whose set of values must contain the DN of the search result entry to be joined with that entry. When the `type` attribute is set to `dn-join`: The attribute whose values are the DNs of the entries to be joined with the search result entry.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `reverse-dn-join`: The attribute in related entries whose set of values must contain the DN of the search result entry to be joined with that entry.\n  - `dn-join`: The attribute whose values are the DNs of the entries to be joined with the search result entry.",
				Optional:            true,
			},
			"join_base_dn_type": schema.StringAttribute{
				Description: "Specifies how server should determine the base DN for the internal searches used to identify joined entries.",
				Optional:    true,
			},
			"join_custom_base_dn": schema.StringAttribute{
				Description: "The fixed, administrator-specified base DN for the internal searches used to identify joined entries.",
				Optional:    true,
			},
			"join_scope": schema.StringAttribute{
				Description: "The scope for searches used to identify joined entries.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"join_size_limit": schema.Int64Attribute{
				Description: "The maximum number of entries that may be joined with the source entry, which also corresponds to the maximum number of values that the virtual attribute provider will generate for an entry.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"join_filter": schema.StringAttribute{
				Description: "An optional filter that specifies additional criteria for identifying joined entries. If a join-filter value is specified, then only entries matching that filter (in addition to satisfying the other join criteria) will be joined with the search result entry.",
				Optional:    true,
			},
			"join_attribute": schema.SetAttribute{
				Description: "An optional set of the names of the attributes to include with joined entries.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"value_pattern": schema.SetAttribute{
				Description: "Specifies a pattern for constructing the virtual attribute value using fixed text and attribute values from the entry.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"conflict_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the server is to exhibit for entries that already contain one or more real values for the associated attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"direct_memberships_only": schema.BoolAttribute{
				Description: "Specifies whether to only include groups in which the user is directly associated with and the membership maybe modified via the group entry. Groups in which the user's membership is derived dynamically or through nested groups will not be included.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"included_group_filter": schema.StringAttribute{
				Description: "A search filter that will be used to identify which groups should be included in the values of the virtual attribute. With no value defined (which is the default behavior), all groups that contain the target user will be included.",
				Optional:    true,
			},
			"rewrite_search_filters": schema.StringAttribute{
				Description: "Search filters that include Is Member Of Virtual Attribute searches on dynamic groups can be updated to include the dynamic group filter in the search filter itself. This can allow the backend to more efficiently process the search filter by using attribute indexes sooner in the search processing.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_attribute": schema.StringAttribute{
				Description: "Specifies the source attribute containing the values to use for this virtual attribute.",
				Optional:    true,
			},
			"source_entry_dn_attribute": schema.StringAttribute{
				Description: "Specifies the attribute containing the DN of another entry from which to obtain the source attribute providing the values for this virtual attribute.",
				Optional:    true,
			},
			"source_entry_dn_map": schema.StringAttribute{
				Description: "Specifies a DN map that will be used to identify the entry from which to obtain the source attribute providing the values for this virtual attribute.",
				Optional:    true,
			},
			"bypass_access_control_for_searches": schema.BoolAttribute{
				Description: "Indicates whether searches performed by this virtual attribute provider should be exempted from access control restrictions.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Virtual Attribute",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Virtual Attribute is enabled for use.",
				Required:    true,
			},
			"attribute_type": schema.StringAttribute{
				Description: "Specifies the attribute type for the attribute whose values are to be dynamically assigned by the virtual attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"base_dn": schema.SetAttribute{
				Description: "Specifies the base DNs for the branches containing entries that are eligible to use this virtual attribute.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"group_dn": schema.SetAttribute{
				Description: "Specifies the DNs of the groups whose members can be eligible to use this virtual attribute.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"filter": schema.SetAttribute{
				Description: "Specifies the search filters to be applied against entries to determine if the virtual attribute is to be generated for those entries.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"client_connection_policy": schema.SetAttribute{
				Description: "Specifies a set of client connection policies for which this Virtual Attribute should be generated. If this is undefined, then this Virtual Attribute will always be generated. If it is associated with one or more client connection policies, then this Virtual Attribute will be generated only for operations requested by clients assigned to one of those client connection policies.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"require_explicit_request_by_name": schema.BoolAttribute{
				Description: "Indicates whether attributes of this type must be explicitly included by name in the list of requested attributes. Note that this will only apply to virtual attributes which are associated with an attribute type that is operational. It will be ignored for virtual attributes associated with a non-operational attribute type.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"multiple_virtual_attribute_evaluation_order_index": schema.Int64Attribute{
				Description: "Specifies the order in which virtual attribute definitions for the same attribute type will be evaluated when generating values for an entry.",
				Optional:    true,
			},
			"multiple_virtual_attribute_merge_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that will be exhibited for cases in which multiple virtual attribute definitions apply to the same multivalued attribute type. This will be ignored for single-valued attribute types.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_index_conflicts": schema.BoolAttribute{
				Description: "Indicates whether the server should allow creating or altering this virtual attribute definition even if it conflicts with one or more indexes defined in the server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Optional = false
		typeAttr.Required = false
		typeAttr.Computed = true
		typeAttr.PlanModifiers = []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		}
		typeAttr.Validators = []validator.String{
			stringvalidator.OneOf([]string{"mirror", "entry-checksum", "member-of-server-group", "constructed", "is-member-of", "custom", "num-subordinates", "reverse-dn-join", "identify-references", "user-defined", "current-time", "short-unique-id", "entry-dn", "has-subordinates", "equality-join", "groovy-scripted", "instance-name", "replication-state-detail", "member", "password-policy-state-json", "subschema-subentry", "dn-join", "third-party"}...),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		schemaDef.Attributes["sequence_number_attribute"] = schema.StringAttribute{
			Description: "Specifies the name or OID of the attribute which contains the sequence number from which unique identifiers are generated. The attribute should have Integer syntax or a String syntax permitting integer values. If this property is modified then the filter property should be updated accordingly so that only entries containing the sequence number attribute are eligible to have a value generated for this virtual attribute.",
		}
		schemaDef.Attributes["return_utc_time"] = schema.BoolAttribute{
			Description: "Indicates whether to return current time in UTC.",
		}
		schemaDef.Attributes["include_milliseconds"] = schema.BoolAttribute{
			Description: "Indicates whether the current time includes millisecond precision.",
		}
		schemaDef.Attributes["exclude_operational_attributes"] = schema.BoolAttribute{
			Description: "Indicates whether all operational attributes should be excluded from the generated checksum.",
		}
		schemaDef.Attributes["excluded_attribute"] = schema.SetAttribute{
			Description: "Specifies the attributes that should be excluded from the checksum calculation.",
			ElementType: types.StringType,
		}
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsVirtualAttribute() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherValidator(
			path.MatchRoot("type"),
			[]string{"mirror"},
			resourcevalidator.Conflicting(
				path.MatchRoot("source_entry_dn_attribute"),
				path.MatchRoot("source_entry_dn_map"),
			),
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("conflict_behavior"),
			path.MatchRoot("type"),
			[]string{"mirror", "entry-checksum", "member-of-server-group", "constructed", "is-member-of", "custom", "num-subordinates", "reverse-dn-join", "identify-references", "user-defined", "current-time", "entry-dn", "has-subordinates", "equality-join", "groovy-scripted", "instance-name", "member", "subschema-subentry", "dn-join", "third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("source_attribute"),
			path.MatchRoot("type"),
			[]string{"mirror"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("source_entry_dn_attribute"),
			path.MatchRoot("type"),
			[]string{"mirror"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("source_entry_dn_map"),
			path.MatchRoot("type"),
			[]string{"mirror"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("bypass_access_control_for_searches"),
			path.MatchRoot("type"),
			[]string{"mirror"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("description"),
			path.MatchRoot("type"),
			[]string{"mirror", "entry-checksum", "member-of-server-group", "constructed", "is-member-of", "custom", "num-subordinates", "reverse-dn-join", "identify-references", "user-defined", "current-time", "short-unique-id", "entry-dn", "has-subordinates", "equality-join", "groovy-scripted", "instance-name", "member", "password-policy-state-json", "subschema-subentry", "dn-join", "third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("attribute_type"),
			path.MatchRoot("type"),
			[]string{"mirror", "entry-checksum", "member-of-server-group", "constructed", "is-member-of", "custom", "num-subordinates", "reverse-dn-join", "identify-references", "user-defined", "current-time", "short-unique-id", "entry-dn", "has-subordinates", "equality-join", "groovy-scripted", "instance-name", "member", "subschema-subentry", "dn-join", "third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("base_dn"),
			path.MatchRoot("type"),
			[]string{"mirror", "entry-checksum", "member-of-server-group", "constructed", "is-member-of", "custom", "num-subordinates", "reverse-dn-join", "identify-references", "user-defined", "current-time", "short-unique-id", "entry-dn", "has-subordinates", "equality-join", "groovy-scripted", "instance-name", "member", "password-policy-state-json", "subschema-subentry", "dn-join", "third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("group_dn"),
			path.MatchRoot("type"),
			[]string{"mirror", "entry-checksum", "member-of-server-group", "constructed", "is-member-of", "custom", "num-subordinates", "reverse-dn-join", "identify-references", "user-defined", "current-time", "short-unique-id", "entry-dn", "has-subordinates", "equality-join", "groovy-scripted", "instance-name", "member", "password-policy-state-json", "subschema-subentry", "dn-join", "third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("filter"),
			path.MatchRoot("type"),
			[]string{"mirror", "entry-checksum", "member-of-server-group", "constructed", "is-member-of", "custom", "num-subordinates", "reverse-dn-join", "identify-references", "user-defined", "current-time", "short-unique-id", "entry-dn", "has-subordinates", "equality-join", "groovy-scripted", "instance-name", "member", "password-policy-state-json", "subschema-subentry", "dn-join", "third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("client_connection_policy"),
			path.MatchRoot("type"),
			[]string{"mirror", "entry-checksum", "member-of-server-group", "constructed", "is-member-of", "custom", "num-subordinates", "reverse-dn-join", "identify-references", "user-defined", "current-time", "short-unique-id", "entry-dn", "has-subordinates", "equality-join", "groovy-scripted", "instance-name", "member", "password-policy-state-json", "subschema-subentry", "dn-join", "third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("require_explicit_request_by_name"),
			path.MatchRoot("type"),
			[]string{"mirror", "entry-checksum", "member-of-server-group", "constructed", "is-member-of", "custom", "num-subordinates", "reverse-dn-join", "identify-references", "user-defined", "current-time", "entry-dn", "has-subordinates", "equality-join", "groovy-scripted", "instance-name", "replication-state-detail", "member", "password-policy-state-json", "subschema-subentry", "dn-join", "third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("multiple_virtual_attribute_evaluation_order_index"),
			path.MatchRoot("type"),
			[]string{"mirror", "entry-checksum", "member-of-server-group", "constructed", "is-member-of", "custom", "num-subordinates", "reverse-dn-join", "identify-references", "user-defined", "current-time", "short-unique-id", "entry-dn", "has-subordinates", "equality-join", "groovy-scripted", "instance-name", "member", "password-policy-state-json", "subschema-subentry", "dn-join", "third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("multiple_virtual_attribute_merge_behavior"),
			path.MatchRoot("type"),
			[]string{"mirror", "entry-checksum", "member-of-server-group", "constructed", "is-member-of", "custom", "num-subordinates", "reverse-dn-join", "identify-references", "user-defined", "current-time", "short-unique-id", "entry-dn", "has-subordinates", "equality-join", "groovy-scripted", "instance-name", "member", "subschema-subentry", "dn-join", "third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("allow_index_conflicts"),
			path.MatchRoot("type"),
			[]string{"mirror", "entry-checksum", "member-of-server-group", "constructed", "is-member-of", "custom", "num-subordinates", "reverse-dn-join", "identify-references", "user-defined", "current-time", "short-unique-id", "entry-dn", "has-subordinates", "equality-join", "groovy-scripted", "instance-name", "member", "subschema-subentry", "dn-join", "third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("value_pattern"),
			path.MatchRoot("type"),
			[]string{"constructed"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("direct_memberships_only"),
			path.MatchRoot("type"),
			[]string{"is-member-of"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("included_group_filter"),
			path.MatchRoot("type"),
			[]string{"is-member-of"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("rewrite_search_filters"),
			path.MatchRoot("type"),
			[]string{"is-member-of"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("join_dn_attribute"),
			path.MatchRoot("type"),
			[]string{"reverse-dn-join", "dn-join"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("join_base_dn_type"),
			path.MatchRoot("type"),
			[]string{"reverse-dn-join", "equality-join", "dn-join"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("join_custom_base_dn"),
			path.MatchRoot("type"),
			[]string{"reverse-dn-join", "equality-join", "dn-join"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("join_scope"),
			path.MatchRoot("type"),
			[]string{"reverse-dn-join", "equality-join", "dn-join"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("join_size_limit"),
			path.MatchRoot("type"),
			[]string{"reverse-dn-join", "equality-join", "dn-join"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("join_filter"),
			path.MatchRoot("type"),
			[]string{"reverse-dn-join", "equality-join", "dn-join"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("join_attribute"),
			path.MatchRoot("type"),
			[]string{"reverse-dn-join", "equality-join", "dn-join"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("referenced_by_attribute"),
			path.MatchRoot("type"),
			[]string{"identify-references"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("reference_search_base_dn"),
			path.MatchRoot("type"),
			[]string{"identify-references"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("value"),
			path.MatchRoot("type"),
			[]string{"user-defined"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("join_source_attribute"),
			path.MatchRoot("type"),
			[]string{"equality-join"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("join_target_attribute"),
			path.MatchRoot("type"),
			[]string{"equality-join"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("join_match_all"),
			path.MatchRoot("type"),
			[]string{"equality-join"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("script_class"),
			path.MatchRoot("type"),
			[]string{"groovy-scripted"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("script_argument"),
			path.MatchRoot("type"),
			[]string{"groovy-scripted"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("allow_retrieving_membership"),
			path.MatchRoot("type"),
			[]string{"member"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_class"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_argument"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
	}
}

// Add config validators
func (r virtualAttributeResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsVirtualAttribute()
}

// Add config validators
func (r defaultVirtualAttributeResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	validators := []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("exclude_operational_attributes"),
			path.MatchRoot("type"),
			[]string{"entry-checksum"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("excluded_attribute"),
			path.MatchRoot("type"),
			[]string{"entry-checksum"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("return_utc_time"),
			path.MatchRoot("type"),
			[]string{"current-time"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_milliseconds"),
			path.MatchRoot("type"),
			[]string{"current-time"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("sequence_number_attribute"),
			path.MatchRoot("type"),
			[]string{"short-unique-id"},
		),
	}
	return append(configValidatorsVirtualAttribute(), validators...)
}

// Add optional fields to create request for mirror virtual-attribute
func addOptionalMirrorVirtualAttributeFields(ctx context.Context, addRequest *client.AddMirrorVirtualAttributeRequest, plan virtualAttributeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConflictBehavior) {
		conflictBehavior, err := client.NewEnumvirtualAttributeConflictBehaviorPropFromValue(plan.ConflictBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConflictBehavior = conflictBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SourceEntryDNAttribute) {
		addRequest.SourceEntryDNAttribute = plan.SourceEntryDNAttribute.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SourceEntryDNMap) {
		addRequest.SourceEntryDNMap = plan.SourceEntryDNMap.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.BypassAccessControlForSearches) {
		addRequest.BypassAccessControlForSearches = plan.BypassAccessControlForSearches.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.BaseDN) {
		var slice []string
		plan.BaseDN.ElementsAs(ctx, &slice, false)
		addRequest.BaseDN = slice
	}
	if internaltypes.IsDefined(plan.GroupDN) {
		var slice []string
		plan.GroupDN.ElementsAs(ctx, &slice, false)
		addRequest.GroupDN = slice
	}
	if internaltypes.IsDefined(plan.Filter) {
		var slice []string
		plan.Filter.ElementsAs(ctx, &slice, false)
		addRequest.Filter = slice
	}
	if internaltypes.IsDefined(plan.ClientConnectionPolicy) {
		var slice []string
		plan.ClientConnectionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.ClientConnectionPolicy = slice
	}
	if internaltypes.IsDefined(plan.RequireExplicitRequestByName) {
		addRequest.RequireExplicitRequestByName = plan.RequireExplicitRequestByName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MultipleVirtualAttributeEvaluationOrderIndex) {
		addRequest.MultipleVirtualAttributeEvaluationOrderIndex = plan.MultipleVirtualAttributeEvaluationOrderIndex.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MultipleVirtualAttributeMergeBehavior) {
		multipleVirtualAttributeMergeBehavior, err := client.NewEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorPropFromValue(plan.MultipleVirtualAttributeMergeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.MultipleVirtualAttributeMergeBehavior = multipleVirtualAttributeMergeBehavior
	}
	if internaltypes.IsDefined(plan.AllowIndexConflicts) {
		addRequest.AllowIndexConflicts = plan.AllowIndexConflicts.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for constructed virtual-attribute
func addOptionalConstructedVirtualAttributeFields(ctx context.Context, addRequest *client.AddConstructedVirtualAttributeRequest, plan virtualAttributeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.BaseDN) {
		var slice []string
		plan.BaseDN.ElementsAs(ctx, &slice, false)
		addRequest.BaseDN = slice
	}
	if internaltypes.IsDefined(plan.GroupDN) {
		var slice []string
		plan.GroupDN.ElementsAs(ctx, &slice, false)
		addRequest.GroupDN = slice
	}
	if internaltypes.IsDefined(plan.Filter) {
		var slice []string
		plan.Filter.ElementsAs(ctx, &slice, false)
		addRequest.Filter = slice
	}
	if internaltypes.IsDefined(plan.ClientConnectionPolicy) {
		var slice []string
		plan.ClientConnectionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.ClientConnectionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConflictBehavior) {
		conflictBehavior, err := client.NewEnumvirtualAttributeConflictBehaviorPropFromValue(plan.ConflictBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConflictBehavior = conflictBehavior
	}
	if internaltypes.IsDefined(plan.RequireExplicitRequestByName) {
		addRequest.RequireExplicitRequestByName = plan.RequireExplicitRequestByName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MultipleVirtualAttributeEvaluationOrderIndex) {
		addRequest.MultipleVirtualAttributeEvaluationOrderIndex = plan.MultipleVirtualAttributeEvaluationOrderIndex.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MultipleVirtualAttributeMergeBehavior) {
		multipleVirtualAttributeMergeBehavior, err := client.NewEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorPropFromValue(plan.MultipleVirtualAttributeMergeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.MultipleVirtualAttributeMergeBehavior = multipleVirtualAttributeMergeBehavior
	}
	if internaltypes.IsDefined(plan.AllowIndexConflicts) {
		addRequest.AllowIndexConflicts = plan.AllowIndexConflicts.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for is-member-of virtual-attribute
func addOptionalIsMemberOfVirtualAttributeFields(ctx context.Context, addRequest *client.AddIsMemberOfVirtualAttributeRequest, plan virtualAttributeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConflictBehavior) {
		conflictBehavior, err := client.NewEnumvirtualAttributeConflictBehaviorPropFromValue(plan.ConflictBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConflictBehavior = conflictBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AttributeType) {
		addRequest.AttributeType = plan.AttributeType.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.DirectMembershipsOnly) {
		addRequest.DirectMembershipsOnly = plan.DirectMembershipsOnly.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IncludedGroupFilter) {
		addRequest.IncludedGroupFilter = plan.IncludedGroupFilter.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RewriteSearchFilters) {
		rewriteSearchFilters, err := client.NewEnumvirtualAttributeRewriteSearchFiltersPropFromValue(plan.RewriteSearchFilters.ValueString())
		if err != nil {
			return err
		}
		addRequest.RewriteSearchFilters = rewriteSearchFilters
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.BaseDN) {
		var slice []string
		plan.BaseDN.ElementsAs(ctx, &slice, false)
		addRequest.BaseDN = slice
	}
	if internaltypes.IsDefined(plan.GroupDN) {
		var slice []string
		plan.GroupDN.ElementsAs(ctx, &slice, false)
		addRequest.GroupDN = slice
	}
	if internaltypes.IsDefined(plan.Filter) {
		var slice []string
		plan.Filter.ElementsAs(ctx, &slice, false)
		addRequest.Filter = slice
	}
	if internaltypes.IsDefined(plan.ClientConnectionPolicy) {
		var slice []string
		plan.ClientConnectionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.ClientConnectionPolicy = slice
	}
	if internaltypes.IsDefined(plan.RequireExplicitRequestByName) {
		addRequest.RequireExplicitRequestByName = plan.RequireExplicitRequestByName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MultipleVirtualAttributeEvaluationOrderIndex) {
		addRequest.MultipleVirtualAttributeEvaluationOrderIndex = plan.MultipleVirtualAttributeEvaluationOrderIndex.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MultipleVirtualAttributeMergeBehavior) {
		multipleVirtualAttributeMergeBehavior, err := client.NewEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorPropFromValue(plan.MultipleVirtualAttributeMergeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.MultipleVirtualAttributeMergeBehavior = multipleVirtualAttributeMergeBehavior
	}
	if internaltypes.IsDefined(plan.AllowIndexConflicts) {
		addRequest.AllowIndexConflicts = plan.AllowIndexConflicts.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for reverse-dn-join virtual-attribute
func addOptionalReverseDnJoinVirtualAttributeFields(ctx context.Context, addRequest *client.AddReverseDnJoinVirtualAttributeRequest, plan virtualAttributeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.JoinCustomBaseDN) {
		addRequest.JoinCustomBaseDN = plan.JoinCustomBaseDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.JoinScope) {
		joinScope, err := client.NewEnumvirtualAttributeJoinScopePropFromValue(plan.JoinScope.ValueString())
		if err != nil {
			return err
		}
		addRequest.JoinScope = joinScope
	}
	if internaltypes.IsDefined(plan.JoinSizeLimit) {
		addRequest.JoinSizeLimit = plan.JoinSizeLimit.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.JoinFilter) {
		addRequest.JoinFilter = plan.JoinFilter.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.JoinAttribute) {
		var slice []string
		plan.JoinAttribute.ElementsAs(ctx, &slice, false)
		addRequest.JoinAttribute = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.BaseDN) {
		var slice []string
		plan.BaseDN.ElementsAs(ctx, &slice, false)
		addRequest.BaseDN = slice
	}
	if internaltypes.IsDefined(plan.GroupDN) {
		var slice []string
		plan.GroupDN.ElementsAs(ctx, &slice, false)
		addRequest.GroupDN = slice
	}
	if internaltypes.IsDefined(plan.Filter) {
		var slice []string
		plan.Filter.ElementsAs(ctx, &slice, false)
		addRequest.Filter = slice
	}
	if internaltypes.IsDefined(plan.ClientConnectionPolicy) {
		var slice []string
		plan.ClientConnectionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.ClientConnectionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConflictBehavior) {
		conflictBehavior, err := client.NewEnumvirtualAttributeConflictBehaviorPropFromValue(plan.ConflictBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConflictBehavior = conflictBehavior
	}
	if internaltypes.IsDefined(plan.RequireExplicitRequestByName) {
		addRequest.RequireExplicitRequestByName = plan.RequireExplicitRequestByName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MultipleVirtualAttributeEvaluationOrderIndex) {
		addRequest.MultipleVirtualAttributeEvaluationOrderIndex = plan.MultipleVirtualAttributeEvaluationOrderIndex.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MultipleVirtualAttributeMergeBehavior) {
		multipleVirtualAttributeMergeBehavior, err := client.NewEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorPropFromValue(plan.MultipleVirtualAttributeMergeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.MultipleVirtualAttributeMergeBehavior = multipleVirtualAttributeMergeBehavior
	}
	if internaltypes.IsDefined(plan.AllowIndexConflicts) {
		addRequest.AllowIndexConflicts = plan.AllowIndexConflicts.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for identify-references virtual-attribute
func addOptionalIdentifyReferencesVirtualAttributeFields(ctx context.Context, addRequest *client.AddIdentifyReferencesVirtualAttributeRequest, plan virtualAttributeResourceModel) error {
	if internaltypes.IsDefined(plan.ReferenceSearchBaseDN) {
		var slice []string
		plan.ReferenceSearchBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.ReferenceSearchBaseDN = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.BaseDN) {
		var slice []string
		plan.BaseDN.ElementsAs(ctx, &slice, false)
		addRequest.BaseDN = slice
	}
	if internaltypes.IsDefined(plan.GroupDN) {
		var slice []string
		plan.GroupDN.ElementsAs(ctx, &slice, false)
		addRequest.GroupDN = slice
	}
	if internaltypes.IsDefined(plan.Filter) {
		var slice []string
		plan.Filter.ElementsAs(ctx, &slice, false)
		addRequest.Filter = slice
	}
	if internaltypes.IsDefined(plan.ClientConnectionPolicy) {
		var slice []string
		plan.ClientConnectionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.ClientConnectionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConflictBehavior) {
		conflictBehavior, err := client.NewEnumvirtualAttributeConflictBehaviorPropFromValue(plan.ConflictBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConflictBehavior = conflictBehavior
	}
	if internaltypes.IsDefined(plan.RequireExplicitRequestByName) {
		addRequest.RequireExplicitRequestByName = plan.RequireExplicitRequestByName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MultipleVirtualAttributeEvaluationOrderIndex) {
		addRequest.MultipleVirtualAttributeEvaluationOrderIndex = plan.MultipleVirtualAttributeEvaluationOrderIndex.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MultipleVirtualAttributeMergeBehavior) {
		multipleVirtualAttributeMergeBehavior, err := client.NewEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorPropFromValue(plan.MultipleVirtualAttributeMergeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.MultipleVirtualAttributeMergeBehavior = multipleVirtualAttributeMergeBehavior
	}
	if internaltypes.IsDefined(plan.AllowIndexConflicts) {
		addRequest.AllowIndexConflicts = plan.AllowIndexConflicts.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for user-defined virtual-attribute
func addOptionalUserDefinedVirtualAttributeFields(ctx context.Context, addRequest *client.AddUserDefinedVirtualAttributeRequest, plan virtualAttributeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.BaseDN) {
		var slice []string
		plan.BaseDN.ElementsAs(ctx, &slice, false)
		addRequest.BaseDN = slice
	}
	if internaltypes.IsDefined(plan.GroupDN) {
		var slice []string
		plan.GroupDN.ElementsAs(ctx, &slice, false)
		addRequest.GroupDN = slice
	}
	if internaltypes.IsDefined(plan.Filter) {
		var slice []string
		plan.Filter.ElementsAs(ctx, &slice, false)
		addRequest.Filter = slice
	}
	if internaltypes.IsDefined(plan.ClientConnectionPolicy) {
		var slice []string
		plan.ClientConnectionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.ClientConnectionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConflictBehavior) {
		conflictBehavior, err := client.NewEnumvirtualAttributeConflictBehaviorPropFromValue(plan.ConflictBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConflictBehavior = conflictBehavior
	}
	if internaltypes.IsDefined(plan.RequireExplicitRequestByName) {
		addRequest.RequireExplicitRequestByName = plan.RequireExplicitRequestByName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MultipleVirtualAttributeEvaluationOrderIndex) {
		addRequest.MultipleVirtualAttributeEvaluationOrderIndex = plan.MultipleVirtualAttributeEvaluationOrderIndex.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MultipleVirtualAttributeMergeBehavior) {
		multipleVirtualAttributeMergeBehavior, err := client.NewEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorPropFromValue(plan.MultipleVirtualAttributeMergeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.MultipleVirtualAttributeMergeBehavior = multipleVirtualAttributeMergeBehavior
	}
	if internaltypes.IsDefined(plan.AllowIndexConflicts) {
		addRequest.AllowIndexConflicts = plan.AllowIndexConflicts.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for entry-dn virtual-attribute
func addOptionalEntryDnVirtualAttributeFields(ctx context.Context, addRequest *client.AddEntryDnVirtualAttributeRequest, plan virtualAttributeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConflictBehavior) {
		conflictBehavior, err := client.NewEnumvirtualAttributeConflictBehaviorPropFromValue(plan.ConflictBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConflictBehavior = conflictBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AttributeType) {
		addRequest.AttributeType = plan.AttributeType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.BaseDN) {
		var slice []string
		plan.BaseDN.ElementsAs(ctx, &slice, false)
		addRequest.BaseDN = slice
	}
	if internaltypes.IsDefined(plan.GroupDN) {
		var slice []string
		plan.GroupDN.ElementsAs(ctx, &slice, false)
		addRequest.GroupDN = slice
	}
	if internaltypes.IsDefined(plan.Filter) {
		var slice []string
		plan.Filter.ElementsAs(ctx, &slice, false)
		addRequest.Filter = slice
	}
	if internaltypes.IsDefined(plan.ClientConnectionPolicy) {
		var slice []string
		plan.ClientConnectionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.ClientConnectionPolicy = slice
	}
	if internaltypes.IsDefined(plan.RequireExplicitRequestByName) {
		addRequest.RequireExplicitRequestByName = plan.RequireExplicitRequestByName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MultipleVirtualAttributeEvaluationOrderIndex) {
		addRequest.MultipleVirtualAttributeEvaluationOrderIndex = plan.MultipleVirtualAttributeEvaluationOrderIndex.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MultipleVirtualAttributeMergeBehavior) {
		multipleVirtualAttributeMergeBehavior, err := client.NewEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorPropFromValue(plan.MultipleVirtualAttributeMergeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.MultipleVirtualAttributeMergeBehavior = multipleVirtualAttributeMergeBehavior
	}
	if internaltypes.IsDefined(plan.AllowIndexConflicts) {
		addRequest.AllowIndexConflicts = plan.AllowIndexConflicts.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for equality-join virtual-attribute
func addOptionalEqualityJoinVirtualAttributeFields(ctx context.Context, addRequest *client.AddEqualityJoinVirtualAttributeRequest, plan virtualAttributeResourceModel) error {
	if internaltypes.IsDefined(plan.JoinMatchAll) {
		addRequest.JoinMatchAll = plan.JoinMatchAll.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.JoinCustomBaseDN) {
		addRequest.JoinCustomBaseDN = plan.JoinCustomBaseDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.JoinScope) {
		joinScope, err := client.NewEnumvirtualAttributeJoinScopePropFromValue(plan.JoinScope.ValueString())
		if err != nil {
			return err
		}
		addRequest.JoinScope = joinScope
	}
	if internaltypes.IsDefined(plan.JoinSizeLimit) {
		addRequest.JoinSizeLimit = plan.JoinSizeLimit.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.JoinFilter) {
		addRequest.JoinFilter = plan.JoinFilter.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.JoinAttribute) {
		var slice []string
		plan.JoinAttribute.ElementsAs(ctx, &slice, false)
		addRequest.JoinAttribute = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.BaseDN) {
		var slice []string
		plan.BaseDN.ElementsAs(ctx, &slice, false)
		addRequest.BaseDN = slice
	}
	if internaltypes.IsDefined(plan.GroupDN) {
		var slice []string
		plan.GroupDN.ElementsAs(ctx, &slice, false)
		addRequest.GroupDN = slice
	}
	if internaltypes.IsDefined(plan.Filter) {
		var slice []string
		plan.Filter.ElementsAs(ctx, &slice, false)
		addRequest.Filter = slice
	}
	if internaltypes.IsDefined(plan.ClientConnectionPolicy) {
		var slice []string
		plan.ClientConnectionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.ClientConnectionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConflictBehavior) {
		conflictBehavior, err := client.NewEnumvirtualAttributeConflictBehaviorPropFromValue(plan.ConflictBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConflictBehavior = conflictBehavior
	}
	if internaltypes.IsDefined(plan.RequireExplicitRequestByName) {
		addRequest.RequireExplicitRequestByName = plan.RequireExplicitRequestByName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MultipleVirtualAttributeEvaluationOrderIndex) {
		addRequest.MultipleVirtualAttributeEvaluationOrderIndex = plan.MultipleVirtualAttributeEvaluationOrderIndex.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MultipleVirtualAttributeMergeBehavior) {
		multipleVirtualAttributeMergeBehavior, err := client.NewEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorPropFromValue(plan.MultipleVirtualAttributeMergeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.MultipleVirtualAttributeMergeBehavior = multipleVirtualAttributeMergeBehavior
	}
	if internaltypes.IsDefined(plan.AllowIndexConflicts) {
		addRequest.AllowIndexConflicts = plan.AllowIndexConflicts.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for groovy-scripted virtual-attribute
func addOptionalGroovyScriptedVirtualAttributeFields(ctx context.Context, addRequest *client.AddGroovyScriptedVirtualAttributeRequest, plan virtualAttributeResourceModel) error {
	if internaltypes.IsDefined(plan.ScriptArgument) {
		var slice []string
		plan.ScriptArgument.ElementsAs(ctx, &slice, false)
		addRequest.ScriptArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.BaseDN) {
		var slice []string
		plan.BaseDN.ElementsAs(ctx, &slice, false)
		addRequest.BaseDN = slice
	}
	if internaltypes.IsDefined(plan.GroupDN) {
		var slice []string
		plan.GroupDN.ElementsAs(ctx, &slice, false)
		addRequest.GroupDN = slice
	}
	if internaltypes.IsDefined(plan.Filter) {
		var slice []string
		plan.Filter.ElementsAs(ctx, &slice, false)
		addRequest.Filter = slice
	}
	if internaltypes.IsDefined(plan.ClientConnectionPolicy) {
		var slice []string
		plan.ClientConnectionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.ClientConnectionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConflictBehavior) {
		conflictBehavior, err := client.NewEnumvirtualAttributeConflictBehaviorPropFromValue(plan.ConflictBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConflictBehavior = conflictBehavior
	}
	if internaltypes.IsDefined(plan.RequireExplicitRequestByName) {
		addRequest.RequireExplicitRequestByName = plan.RequireExplicitRequestByName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MultipleVirtualAttributeEvaluationOrderIndex) {
		addRequest.MultipleVirtualAttributeEvaluationOrderIndex = plan.MultipleVirtualAttributeEvaluationOrderIndex.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MultipleVirtualAttributeMergeBehavior) {
		multipleVirtualAttributeMergeBehavior, err := client.NewEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorPropFromValue(plan.MultipleVirtualAttributeMergeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.MultipleVirtualAttributeMergeBehavior = multipleVirtualAttributeMergeBehavior
	}
	if internaltypes.IsDefined(plan.AllowIndexConflicts) {
		addRequest.AllowIndexConflicts = plan.AllowIndexConflicts.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for member virtual-attribute
func addOptionalMemberVirtualAttributeFields(ctx context.Context, addRequest *client.AddMemberVirtualAttributeRequest, plan virtualAttributeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConflictBehavior) {
		conflictBehavior, err := client.NewEnumvirtualAttributeConflictBehaviorPropFromValue(plan.ConflictBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConflictBehavior = conflictBehavior
	}
	if internaltypes.IsDefined(plan.AllowRetrievingMembership) {
		addRequest.AllowRetrievingMembership = plan.AllowRetrievingMembership.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.Filter) {
		var slice []string
		plan.Filter.ElementsAs(ctx, &slice, false)
		addRequest.Filter = slice
	}
	if internaltypes.IsDefined(plan.AllowIndexConflicts) {
		addRequest.AllowIndexConflicts = plan.AllowIndexConflicts.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.BaseDN) {
		var slice []string
		plan.BaseDN.ElementsAs(ctx, &slice, false)
		addRequest.BaseDN = slice
	}
	if internaltypes.IsDefined(plan.GroupDN) {
		var slice []string
		plan.GroupDN.ElementsAs(ctx, &slice, false)
		addRequest.GroupDN = slice
	}
	if internaltypes.IsDefined(plan.ClientConnectionPolicy) {
		var slice []string
		plan.ClientConnectionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.ClientConnectionPolicy = slice
	}
	if internaltypes.IsDefined(plan.RequireExplicitRequestByName) {
		addRequest.RequireExplicitRequestByName = plan.RequireExplicitRequestByName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MultipleVirtualAttributeEvaluationOrderIndex) {
		addRequest.MultipleVirtualAttributeEvaluationOrderIndex = plan.MultipleVirtualAttributeEvaluationOrderIndex.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MultipleVirtualAttributeMergeBehavior) {
		multipleVirtualAttributeMergeBehavior, err := client.NewEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorPropFromValue(plan.MultipleVirtualAttributeMergeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.MultipleVirtualAttributeMergeBehavior = multipleVirtualAttributeMergeBehavior
	}
	return nil
}

// Add optional fields to create request for password-policy-state-json virtual-attribute
func addOptionalPasswordPolicyStateJsonVirtualAttributeFields(ctx context.Context, addRequest *client.AddPasswordPolicyStateJsonVirtualAttributeRequest, plan virtualAttributeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.BaseDN) {
		var slice []string
		plan.BaseDN.ElementsAs(ctx, &slice, false)
		addRequest.BaseDN = slice
	}
	if internaltypes.IsDefined(plan.GroupDN) {
		var slice []string
		plan.GroupDN.ElementsAs(ctx, &slice, false)
		addRequest.GroupDN = slice
	}
	if internaltypes.IsDefined(plan.Filter) {
		var slice []string
		plan.Filter.ElementsAs(ctx, &slice, false)
		addRequest.Filter = slice
	}
	if internaltypes.IsDefined(plan.ClientConnectionPolicy) {
		var slice []string
		plan.ClientConnectionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.ClientConnectionPolicy = slice
	}
	if internaltypes.IsDefined(plan.RequireExplicitRequestByName) {
		addRequest.RequireExplicitRequestByName = plan.RequireExplicitRequestByName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MultipleVirtualAttributeEvaluationOrderIndex) {
		addRequest.MultipleVirtualAttributeEvaluationOrderIndex = plan.MultipleVirtualAttributeEvaluationOrderIndex.ValueInt64Pointer()
	}
	return nil
}

// Add optional fields to create request for dn-join virtual-attribute
func addOptionalDnJoinVirtualAttributeFields(ctx context.Context, addRequest *client.AddDnJoinVirtualAttributeRequest, plan virtualAttributeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.JoinCustomBaseDN) {
		addRequest.JoinCustomBaseDN = plan.JoinCustomBaseDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.JoinScope) {
		joinScope, err := client.NewEnumvirtualAttributeJoinScopePropFromValue(plan.JoinScope.ValueString())
		if err != nil {
			return err
		}
		addRequest.JoinScope = joinScope
	}
	if internaltypes.IsDefined(plan.JoinSizeLimit) {
		addRequest.JoinSizeLimit = plan.JoinSizeLimit.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.JoinFilter) {
		addRequest.JoinFilter = plan.JoinFilter.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.JoinAttribute) {
		var slice []string
		plan.JoinAttribute.ElementsAs(ctx, &slice, false)
		addRequest.JoinAttribute = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.BaseDN) {
		var slice []string
		plan.BaseDN.ElementsAs(ctx, &slice, false)
		addRequest.BaseDN = slice
	}
	if internaltypes.IsDefined(plan.GroupDN) {
		var slice []string
		plan.GroupDN.ElementsAs(ctx, &slice, false)
		addRequest.GroupDN = slice
	}
	if internaltypes.IsDefined(plan.Filter) {
		var slice []string
		plan.Filter.ElementsAs(ctx, &slice, false)
		addRequest.Filter = slice
	}
	if internaltypes.IsDefined(plan.ClientConnectionPolicy) {
		var slice []string
		plan.ClientConnectionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.ClientConnectionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConflictBehavior) {
		conflictBehavior, err := client.NewEnumvirtualAttributeConflictBehaviorPropFromValue(plan.ConflictBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConflictBehavior = conflictBehavior
	}
	if internaltypes.IsDefined(plan.RequireExplicitRequestByName) {
		addRequest.RequireExplicitRequestByName = plan.RequireExplicitRequestByName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MultipleVirtualAttributeEvaluationOrderIndex) {
		addRequest.MultipleVirtualAttributeEvaluationOrderIndex = plan.MultipleVirtualAttributeEvaluationOrderIndex.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MultipleVirtualAttributeMergeBehavior) {
		multipleVirtualAttributeMergeBehavior, err := client.NewEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorPropFromValue(plan.MultipleVirtualAttributeMergeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.MultipleVirtualAttributeMergeBehavior = multipleVirtualAttributeMergeBehavior
	}
	if internaltypes.IsDefined(plan.AllowIndexConflicts) {
		addRequest.AllowIndexConflicts = plan.AllowIndexConflicts.ValueBoolPointer()
	}
	return nil
}

// Add optional fields to create request for third-party virtual-attribute
func addOptionalThirdPartyVirtualAttributeFields(ctx context.Context, addRequest *client.AddThirdPartyVirtualAttributeRequest, plan virtualAttributeResourceModel) error {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.BaseDN) {
		var slice []string
		plan.BaseDN.ElementsAs(ctx, &slice, false)
		addRequest.BaseDN = slice
	}
	if internaltypes.IsDefined(plan.GroupDN) {
		var slice []string
		plan.GroupDN.ElementsAs(ctx, &slice, false)
		addRequest.GroupDN = slice
	}
	if internaltypes.IsDefined(plan.Filter) {
		var slice []string
		plan.Filter.ElementsAs(ctx, &slice, false)
		addRequest.Filter = slice
	}
	if internaltypes.IsDefined(plan.ClientConnectionPolicy) {
		var slice []string
		plan.ClientConnectionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.ClientConnectionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConflictBehavior) {
		conflictBehavior, err := client.NewEnumvirtualAttributeConflictBehaviorPropFromValue(plan.ConflictBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConflictBehavior = conflictBehavior
	}
	if internaltypes.IsDefined(plan.RequireExplicitRequestByName) {
		addRequest.RequireExplicitRequestByName = plan.RequireExplicitRequestByName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MultipleVirtualAttributeEvaluationOrderIndex) {
		addRequest.MultipleVirtualAttributeEvaluationOrderIndex = plan.MultipleVirtualAttributeEvaluationOrderIndex.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MultipleVirtualAttributeMergeBehavior) {
		multipleVirtualAttributeMergeBehavior, err := client.NewEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorPropFromValue(plan.MultipleVirtualAttributeMergeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.MultipleVirtualAttributeMergeBehavior = multipleVirtualAttributeMergeBehavior
	}
	if internaltypes.IsDefined(plan.AllowIndexConflicts) {
		addRequest.AllowIndexConflicts = plan.AllowIndexConflicts.ValueBoolPointer()
	}
	return nil
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateVirtualAttributeUnknownValues(model *virtualAttributeResourceModel) {
	if model.ValuePattern.IsUnknown() || model.ValuePattern.IsNull() {
		model.ValuePattern, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ScriptArgument.IsUnknown() || model.ScriptArgument.IsNull() {
		model.ScriptArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.Filter.IsUnknown() || model.Filter.IsNull() {
		model.Filter, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.BaseDN.IsUnknown() || model.BaseDN.IsNull() {
		model.BaseDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.Value.IsUnknown() || model.Value.IsNull() {
		model.Value, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.GroupDN.IsUnknown() || model.GroupDN.IsNull() {
		model.GroupDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ReferenceSearchBaseDN.IsUnknown() || model.ReferenceSearchBaseDN.IsNull() {
		model.ReferenceSearchBaseDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.JoinAttribute.IsUnknown() || model.JoinAttribute.IsNull() {
		model.JoinAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ClientConnectionPolicy.IsUnknown() || model.ClientConnectionPolicy.IsNull() {
		model.ClientConnectionPolicy, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ReferencedByAttribute.IsUnknown() || model.ReferencedByAttribute.IsNull() {
		model.ReferencedByAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.JoinScope.IsUnknown() || model.JoinScope.IsNull() {
		model.JoinScope = types.StringValue("")
	}
	if model.ConflictBehavior.IsUnknown() || model.ConflictBehavior.IsNull() {
		model.ConflictBehavior = types.StringValue("")
	}
	if model.AttributeType.IsUnknown() || model.AttributeType.IsNull() {
		model.AttributeType = types.StringValue("")
	}
	if model.RewriteSearchFilters.IsUnknown() || model.RewriteSearchFilters.IsNull() {
		model.RewriteSearchFilters = types.StringValue("")
	}
	if model.MultipleVirtualAttributeMergeBehavior.IsUnknown() || model.MultipleVirtualAttributeMergeBehavior.IsNull() {
		model.MultipleVirtualAttributeMergeBehavior = types.StringValue("")
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateVirtualAttributeUnknownValuesDefault(model *defaultVirtualAttributeResourceModel) {
	if model.ValuePattern.IsUnknown() || model.ValuePattern.IsNull() {
		model.ValuePattern, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ScriptArgument.IsUnknown() || model.ScriptArgument.IsNull() {
		model.ScriptArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.Filter.IsUnknown() || model.Filter.IsNull() {
		model.Filter, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.BaseDN.IsUnknown() || model.BaseDN.IsNull() {
		model.BaseDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.Value.IsUnknown() || model.Value.IsNull() {
		model.Value, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.GroupDN.IsUnknown() || model.GroupDN.IsNull() {
		model.GroupDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ReferenceSearchBaseDN.IsUnknown() || model.ReferenceSearchBaseDN.IsNull() {
		model.ReferenceSearchBaseDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExcludedAttribute.IsUnknown() || model.ExcludedAttribute.IsNull() {
		model.ExcludedAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.JoinAttribute.IsUnknown() || model.JoinAttribute.IsNull() {
		model.JoinAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ClientConnectionPolicy.IsUnknown() || model.ClientConnectionPolicy.IsNull() {
		model.ClientConnectionPolicy, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ReferencedByAttribute.IsUnknown() || model.ReferencedByAttribute.IsNull() {
		model.ReferencedByAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.IncludedGroupFilter.IsUnknown() || model.IncludedGroupFilter.IsNull() {
		model.IncludedGroupFilter = types.StringValue("")
	}
	if model.JoinScope.IsUnknown() || model.JoinScope.IsNull() {
		model.JoinScope = types.StringValue("")
	}
	if model.ConflictBehavior.IsUnknown() || model.ConflictBehavior.IsNull() {
		model.ConflictBehavior = types.StringValue("")
	}
	if model.ExtensionClass.IsUnknown() || model.ExtensionClass.IsNull() {
		model.ExtensionClass = types.StringValue("")
	}
	if model.SourceEntryDNMap.IsUnknown() || model.SourceEntryDNMap.IsNull() {
		model.SourceEntryDNMap = types.StringValue("")
	}
	if model.RewriteSearchFilters.IsUnknown() || model.RewriteSearchFilters.IsNull() {
		model.RewriteSearchFilters = types.StringValue("")
	}
	if model.JoinDNAttribute.IsUnknown() || model.JoinDNAttribute.IsNull() {
		model.JoinDNAttribute = types.StringValue("")
	}
	if model.SourceEntryDNAttribute.IsUnknown() || model.SourceEntryDNAttribute.IsNull() {
		model.SourceEntryDNAttribute = types.StringValue("")
	}
	if model.SourceAttribute.IsUnknown() || model.SourceAttribute.IsNull() {
		model.SourceAttribute = types.StringValue("")
	}
	if model.SequenceNumberAttribute.IsUnknown() || model.SequenceNumberAttribute.IsNull() {
		model.SequenceNumberAttribute = types.StringValue("")
	}
	if model.JoinTargetAttribute.IsUnknown() || model.JoinTargetAttribute.IsNull() {
		model.JoinTargetAttribute = types.StringValue("")
	}
	if model.AttributeType.IsUnknown() || model.AttributeType.IsNull() {
		model.AttributeType = types.StringValue("")
	}
	if model.JoinCustomBaseDN.IsUnknown() || model.JoinCustomBaseDN.IsNull() {
		model.JoinCustomBaseDN = types.StringValue("")
	}
	if model.JoinBaseDNType.IsUnknown() || model.JoinBaseDNType.IsNull() {
		model.JoinBaseDNType = types.StringValue("")
	}
	if model.JoinSourceAttribute.IsUnknown() || model.JoinSourceAttribute.IsNull() {
		model.JoinSourceAttribute = types.StringValue("")
	}
	if model.ScriptClass.IsUnknown() || model.ScriptClass.IsNull() {
		model.ScriptClass = types.StringValue("")
	}
	if model.MultipleVirtualAttributeMergeBehavior.IsUnknown() || model.MultipleVirtualAttributeMergeBehavior.IsNull() {
		model.MultipleVirtualAttributeMergeBehavior = types.StringValue("")
	}
	if model.JoinFilter.IsUnknown() || model.JoinFilter.IsNull() {
		model.JoinFilter = types.StringValue("")
	}
}

// Read a MirrorVirtualAttributeResponse object into the model struct
func readMirrorVirtualAttributeResponse(ctx context.Context, r *client.MirrorVirtualAttributeResponse, state *virtualAttributeResourceModel, expectedValues *virtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("mirror")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.SourceAttribute = types.StringValue(r.SourceAttribute)
	state.SourceEntryDNAttribute = internaltypes.StringTypeOrNil(r.SourceEntryDNAttribute, internaltypes.IsEmptyString(expectedValues.SourceEntryDNAttribute))
	state.SourceEntryDNMap = internaltypes.StringTypeOrNil(r.SourceEntryDNMap, internaltypes.IsEmptyString(expectedValues.SourceEntryDNMap))
	state.BypassAccessControlForSearches = internaltypes.BoolTypeOrNil(r.BypassAccessControlForSearches)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValues(state)
}

// Read a MirrorVirtualAttributeResponse object into the model struct
func readMirrorVirtualAttributeResponseDefault(ctx context.Context, r *client.MirrorVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("mirror")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.SourceAttribute = types.StringValue(r.SourceAttribute)
	state.SourceEntryDNAttribute = internaltypes.StringTypeOrNil(r.SourceEntryDNAttribute, true)
	state.SourceEntryDNMap = internaltypes.StringTypeOrNil(r.SourceEntryDNMap, true)
	state.BypassAccessControlForSearches = internaltypes.BoolTypeOrNil(r.BypassAccessControlForSearches)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Read a EntryChecksumVirtualAttributeResponse object into the model struct
func readEntryChecksumVirtualAttributeResponseDefault(ctx context.Context, r *client.EntryChecksumVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("entry-checksum")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.ExcludeOperationalAttributes = internaltypes.BoolTypeOrNil(r.ExcludeOperationalAttributes)
	state.ExcludedAttribute = internaltypes.GetStringSet(r.ExcludedAttribute)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Read a MemberOfServerGroupVirtualAttributeResponse object into the model struct
func readMemberOfServerGroupVirtualAttributeResponseDefault(ctx context.Context, r *client.MemberOfServerGroupVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("member-of-server-group")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Read a ConstructedVirtualAttributeResponse object into the model struct
func readConstructedVirtualAttributeResponse(ctx context.Context, r *client.ConstructedVirtualAttributeResponse, state *virtualAttributeResourceModel, expectedValues *virtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("constructed")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ValuePattern = internaltypes.GetStringSet(r.ValuePattern)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValues(state)
}

// Read a ConstructedVirtualAttributeResponse object into the model struct
func readConstructedVirtualAttributeResponseDefault(ctx context.Context, r *client.ConstructedVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("constructed")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ValuePattern = internaltypes.GetStringSet(r.ValuePattern)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Read a IsMemberOfVirtualAttributeResponse object into the model struct
func readIsMemberOfVirtualAttributeResponse(ctx context.Context, r *client.IsMemberOfVirtualAttributeResponse, state *virtualAttributeResourceModel, expectedValues *virtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("is-member-of")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.DirectMembershipsOnly = internaltypes.BoolTypeOrNil(r.DirectMembershipsOnly)
	state.IncludedGroupFilter = internaltypes.StringTypeOrNil(r.IncludedGroupFilter, internaltypes.IsEmptyString(expectedValues.IncludedGroupFilter))
	state.RewriteSearchFilters = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeRewriteSearchFiltersProp(r.RewriteSearchFilters), true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValues(state)
}

// Read a IsMemberOfVirtualAttributeResponse object into the model struct
func readIsMemberOfVirtualAttributeResponseDefault(ctx context.Context, r *client.IsMemberOfVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("is-member-of")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.DirectMembershipsOnly = internaltypes.BoolTypeOrNil(r.DirectMembershipsOnly)
	state.IncludedGroupFilter = internaltypes.StringTypeOrNil(r.IncludedGroupFilter, true)
	state.RewriteSearchFilters = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeRewriteSearchFiltersProp(r.RewriteSearchFilters), true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Read a CustomVirtualAttributeResponse object into the model struct
func readCustomVirtualAttributeResponseDefault(ctx context.Context, r *client.CustomVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("custom")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Read a NumSubordinatesVirtualAttributeResponse object into the model struct
func readNumSubordinatesVirtualAttributeResponseDefault(ctx context.Context, r *client.NumSubordinatesVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("num-subordinates")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Read a ReverseDnJoinVirtualAttributeResponse object into the model struct
func readReverseDnJoinVirtualAttributeResponse(ctx context.Context, r *client.ReverseDnJoinVirtualAttributeResponse, state *virtualAttributeResourceModel, expectedValues *virtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("reverse-dn-join")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.JoinDNAttribute = types.StringValue(r.JoinDNAttribute)
	state.JoinBaseDNType = types.StringValue(r.JoinBaseDNType.String())
	state.JoinCustomBaseDN = internaltypes.StringTypeOrNil(r.JoinCustomBaseDN, internaltypes.IsEmptyString(expectedValues.JoinCustomBaseDN))
	state.JoinScope = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeJoinScopeProp(r.JoinScope), true)
	state.JoinSizeLimit = internaltypes.Int64TypeOrNil(r.JoinSizeLimit)
	state.JoinFilter = internaltypes.StringTypeOrNil(r.JoinFilter, internaltypes.IsEmptyString(expectedValues.JoinFilter))
	state.JoinAttribute = internaltypes.GetStringSet(r.JoinAttribute)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValues(state)
}

// Read a ReverseDnJoinVirtualAttributeResponse object into the model struct
func readReverseDnJoinVirtualAttributeResponseDefault(ctx context.Context, r *client.ReverseDnJoinVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("reverse-dn-join")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.JoinDNAttribute = types.StringValue(r.JoinDNAttribute)
	state.JoinBaseDNType = types.StringValue(r.JoinBaseDNType.String())
	state.JoinCustomBaseDN = internaltypes.StringTypeOrNil(r.JoinCustomBaseDN, true)
	state.JoinScope = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeJoinScopeProp(r.JoinScope), true)
	state.JoinSizeLimit = internaltypes.Int64TypeOrNil(r.JoinSizeLimit)
	state.JoinFilter = internaltypes.StringTypeOrNil(r.JoinFilter, true)
	state.JoinAttribute = internaltypes.GetStringSet(r.JoinAttribute)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Read a IdentifyReferencesVirtualAttributeResponse object into the model struct
func readIdentifyReferencesVirtualAttributeResponse(ctx context.Context, r *client.IdentifyReferencesVirtualAttributeResponse, state *virtualAttributeResourceModel, expectedValues *virtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("identify-references")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReferencedByAttribute = internaltypes.GetStringSet(r.ReferencedByAttribute)
	state.ReferenceSearchBaseDN = internaltypes.GetStringSet(r.ReferenceSearchBaseDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValues(state)
}

// Read a IdentifyReferencesVirtualAttributeResponse object into the model struct
func readIdentifyReferencesVirtualAttributeResponseDefault(ctx context.Context, r *client.IdentifyReferencesVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("identify-references")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReferencedByAttribute = internaltypes.GetStringSet(r.ReferencedByAttribute)
	state.ReferenceSearchBaseDN = internaltypes.GetStringSet(r.ReferenceSearchBaseDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Read a UserDefinedVirtualAttributeResponse object into the model struct
func readUserDefinedVirtualAttributeResponse(ctx context.Context, r *client.UserDefinedVirtualAttributeResponse, state *virtualAttributeResourceModel, expectedValues *virtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("user-defined")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Value = internaltypes.GetStringSet(r.Value)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValues(state)
}

// Read a UserDefinedVirtualAttributeResponse object into the model struct
func readUserDefinedVirtualAttributeResponseDefault(ctx context.Context, r *client.UserDefinedVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("user-defined")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Value = internaltypes.GetStringSet(r.Value)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Read a CurrentTimeVirtualAttributeResponse object into the model struct
func readCurrentTimeVirtualAttributeResponseDefault(ctx context.Context, r *client.CurrentTimeVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("current-time")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.ReturnUtcTime = internaltypes.BoolTypeOrNil(r.ReturnUtcTime)
	state.IncludeMilliseconds = internaltypes.BoolTypeOrNil(r.IncludeMilliseconds)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Read a ShortUniqueIdVirtualAttributeResponse object into the model struct
func readShortUniqueIdVirtualAttributeResponseDefault(ctx context.Context, r *client.ShortUniqueIdVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("short-unique-id")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SequenceNumberAttribute = types.StringValue(r.SequenceNumberAttribute)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Read a EntryDnVirtualAttributeResponse object into the model struct
func readEntryDnVirtualAttributeResponse(ctx context.Context, r *client.EntryDnVirtualAttributeResponse, state *virtualAttributeResourceModel, expectedValues *virtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("entry-dn")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValues(state)
}

// Read a EntryDnVirtualAttributeResponse object into the model struct
func readEntryDnVirtualAttributeResponseDefault(ctx context.Context, r *client.EntryDnVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("entry-dn")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Read a HasSubordinatesVirtualAttributeResponse object into the model struct
func readHasSubordinatesVirtualAttributeResponseDefault(ctx context.Context, r *client.HasSubordinatesVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("has-subordinates")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Read a EqualityJoinVirtualAttributeResponse object into the model struct
func readEqualityJoinVirtualAttributeResponse(ctx context.Context, r *client.EqualityJoinVirtualAttributeResponse, state *virtualAttributeResourceModel, expectedValues *virtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("equality-join")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.JoinSourceAttribute = types.StringValue(r.JoinSourceAttribute)
	state.JoinTargetAttribute = types.StringValue(r.JoinTargetAttribute)
	state.JoinMatchAll = internaltypes.BoolTypeOrNil(r.JoinMatchAll)
	state.JoinBaseDNType = types.StringValue(r.JoinBaseDNType.String())
	state.JoinCustomBaseDN = internaltypes.StringTypeOrNil(r.JoinCustomBaseDN, internaltypes.IsEmptyString(expectedValues.JoinCustomBaseDN))
	state.JoinScope = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeJoinScopeProp(r.JoinScope), true)
	state.JoinSizeLimit = internaltypes.Int64TypeOrNil(r.JoinSizeLimit)
	state.JoinFilter = internaltypes.StringTypeOrNil(r.JoinFilter, internaltypes.IsEmptyString(expectedValues.JoinFilter))
	state.JoinAttribute = internaltypes.GetStringSet(r.JoinAttribute)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValues(state)
}

// Read a EqualityJoinVirtualAttributeResponse object into the model struct
func readEqualityJoinVirtualAttributeResponseDefault(ctx context.Context, r *client.EqualityJoinVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("equality-join")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.JoinSourceAttribute = types.StringValue(r.JoinSourceAttribute)
	state.JoinTargetAttribute = types.StringValue(r.JoinTargetAttribute)
	state.JoinMatchAll = internaltypes.BoolTypeOrNil(r.JoinMatchAll)
	state.JoinBaseDNType = types.StringValue(r.JoinBaseDNType.String())
	state.JoinCustomBaseDN = internaltypes.StringTypeOrNil(r.JoinCustomBaseDN, true)
	state.JoinScope = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeJoinScopeProp(r.JoinScope), true)
	state.JoinSizeLimit = internaltypes.Int64TypeOrNil(r.JoinSizeLimit)
	state.JoinFilter = internaltypes.StringTypeOrNil(r.JoinFilter, true)
	state.JoinAttribute = internaltypes.GetStringSet(r.JoinAttribute)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Read a GroovyScriptedVirtualAttributeResponse object into the model struct
func readGroovyScriptedVirtualAttributeResponse(ctx context.Context, r *client.GroovyScriptedVirtualAttributeResponse, state *virtualAttributeResourceModel, expectedValues *virtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValues(state)
}

// Read a GroovyScriptedVirtualAttributeResponse object into the model struct
func readGroovyScriptedVirtualAttributeResponseDefault(ctx context.Context, r *client.GroovyScriptedVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Read a InstanceNameVirtualAttributeResponse object into the model struct
func readInstanceNameVirtualAttributeResponseDefault(ctx context.Context, r *client.InstanceNameVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("instance-name")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Read a ReplicationStateDetailVirtualAttributeResponse object into the model struct
func readReplicationStateDetailVirtualAttributeResponseDefault(ctx context.Context, r *client.ReplicationStateDetailVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("replication-state-detail")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Read a MemberVirtualAttributeResponse object into the model struct
func readMemberVirtualAttributeResponse(ctx context.Context, r *client.MemberVirtualAttributeResponse, state *virtualAttributeResourceModel, expectedValues *virtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("member")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.AllowRetrievingMembership = types.BoolValue(r.AllowRetrievingMembership)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValues(state)
}

// Read a MemberVirtualAttributeResponse object into the model struct
func readMemberVirtualAttributeResponseDefault(ctx context.Context, r *client.MemberVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("member")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.AllowRetrievingMembership = types.BoolValue(r.AllowRetrievingMembership)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Read a PasswordPolicyStateJsonVirtualAttributeResponse object into the model struct
func readPasswordPolicyStateJsonVirtualAttributeResponse(ctx context.Context, r *client.PasswordPolicyStateJsonVirtualAttributeResponse, state *virtualAttributeResourceModel, expectedValues *virtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("password-policy-state-json")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValues(state)
}

// Read a PasswordPolicyStateJsonVirtualAttributeResponse object into the model struct
func readPasswordPolicyStateJsonVirtualAttributeResponseDefault(ctx context.Context, r *client.PasswordPolicyStateJsonVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("password-policy-state-json")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Read a SubschemaSubentryVirtualAttributeResponse object into the model struct
func readSubschemaSubentryVirtualAttributeResponseDefault(ctx context.Context, r *client.SubschemaSubentryVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("subschema-subentry")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Read a DnJoinVirtualAttributeResponse object into the model struct
func readDnJoinVirtualAttributeResponse(ctx context.Context, r *client.DnJoinVirtualAttributeResponse, state *virtualAttributeResourceModel, expectedValues *virtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("dn-join")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.JoinDNAttribute = types.StringValue(r.JoinDNAttribute)
	state.JoinBaseDNType = types.StringValue(r.JoinBaseDNType.String())
	state.JoinCustomBaseDN = internaltypes.StringTypeOrNil(r.JoinCustomBaseDN, internaltypes.IsEmptyString(expectedValues.JoinCustomBaseDN))
	state.JoinScope = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeJoinScopeProp(r.JoinScope), true)
	state.JoinSizeLimit = internaltypes.Int64TypeOrNil(r.JoinSizeLimit)
	state.JoinFilter = internaltypes.StringTypeOrNil(r.JoinFilter, internaltypes.IsEmptyString(expectedValues.JoinFilter))
	state.JoinAttribute = internaltypes.GetStringSet(r.JoinAttribute)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValues(state)
}

// Read a DnJoinVirtualAttributeResponse object into the model struct
func readDnJoinVirtualAttributeResponseDefault(ctx context.Context, r *client.DnJoinVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("dn-join")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.JoinDNAttribute = types.StringValue(r.JoinDNAttribute)
	state.JoinBaseDNType = types.StringValue(r.JoinBaseDNType.String())
	state.JoinCustomBaseDN = internaltypes.StringTypeOrNil(r.JoinCustomBaseDN, true)
	state.JoinScope = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeJoinScopeProp(r.JoinScope), true)
	state.JoinSizeLimit = internaltypes.Int64TypeOrNil(r.JoinSizeLimit)
	state.JoinFilter = internaltypes.StringTypeOrNil(r.JoinFilter, true)
	state.JoinAttribute = internaltypes.GetStringSet(r.JoinAttribute)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Read a ThirdPartyVirtualAttributeResponse object into the model struct
func readThirdPartyVirtualAttributeResponse(ctx context.Context, r *client.ThirdPartyVirtualAttributeResponse, state *virtualAttributeResourceModel, expectedValues *virtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValues(state)
}

// Read a ThirdPartyVirtualAttributeResponse object into the model struct
func readThirdPartyVirtualAttributeResponseDefault(ctx context.Context, r *client.ThirdPartyVirtualAttributeResponse, state *defaultVirtualAttributeResourceModel, expectedValues *defaultVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), true)
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVirtualAttributeUnknownValuesDefault(state)
}

// Create any update operations necessary to make the state match the plan
func createVirtualAttributeOperations(plan virtualAttributeResourceModel, state virtualAttributeResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.ScriptClass, state.ScriptClass, "script-class")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowRetrievingMembership, state.AllowRetrievingMembership, "allow-retrieving-membership")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScriptArgument, state.ScriptArgument, "script-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.JoinSourceAttribute, state.JoinSourceAttribute, "join-source-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.JoinTargetAttribute, state.JoinTargetAttribute, "join-target-attribute")
	operations.AddBoolOperationIfNecessary(&ops, plan.JoinMatchAll, state.JoinMatchAll, "join-match-all")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.Value, state.Value, "value")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ReferencedByAttribute, state.ReferencedByAttribute, "referenced-by-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ReferenceSearchBaseDN, state.ReferenceSearchBaseDN, "reference-search-base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.JoinDNAttribute, state.JoinDNAttribute, "join-dn-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.JoinBaseDNType, state.JoinBaseDNType, "join-base-dn-type")
	operations.AddStringOperationIfNecessary(&ops, plan.JoinCustomBaseDN, state.JoinCustomBaseDN, "join-custom-base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.JoinScope, state.JoinScope, "join-scope")
	operations.AddInt64OperationIfNecessary(&ops, plan.JoinSizeLimit, state.JoinSizeLimit, "join-size-limit")
	operations.AddStringOperationIfNecessary(&ops, plan.JoinFilter, state.JoinFilter, "join-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.JoinAttribute, state.JoinAttribute, "join-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ValuePattern, state.ValuePattern, "value-pattern")
	operations.AddStringOperationIfNecessary(&ops, plan.ConflictBehavior, state.ConflictBehavior, "conflict-behavior")
	operations.AddBoolOperationIfNecessary(&ops, plan.DirectMembershipsOnly, state.DirectMembershipsOnly, "direct-memberships-only")
	operations.AddStringOperationIfNecessary(&ops, plan.IncludedGroupFilter, state.IncludedGroupFilter, "included-group-filter")
	operations.AddStringOperationIfNecessary(&ops, plan.RewriteSearchFilters, state.RewriteSearchFilters, "rewrite-search-filters")
	operations.AddStringOperationIfNecessary(&ops, plan.SourceAttribute, state.SourceAttribute, "source-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.SourceEntryDNAttribute, state.SourceEntryDNAttribute, "source-entry-dn-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.SourceEntryDNMap, state.SourceEntryDNMap, "source-entry-dn-map")
	operations.AddBoolOperationIfNecessary(&ops, plan.BypassAccessControlForSearches, state.BypassAccessControlForSearches, "bypass-access-control-for-searches")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.AttributeType, state.AttributeType, "attribute-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.GroupDN, state.GroupDN, "group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.Filter, state.Filter, "filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ClientConnectionPolicy, state.ClientConnectionPolicy, "client-connection-policy")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequireExplicitRequestByName, state.RequireExplicitRequestByName, "require-explicit-request-by-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.MultipleVirtualAttributeEvaluationOrderIndex, state.MultipleVirtualAttributeEvaluationOrderIndex, "multiple-virtual-attribute-evaluation-order-index")
	operations.AddStringOperationIfNecessary(&ops, plan.MultipleVirtualAttributeMergeBehavior, state.MultipleVirtualAttributeMergeBehavior, "multiple-virtual-attribute-merge-behavior")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowIndexConflicts, state.AllowIndexConflicts, "allow-index-conflicts")
	return ops
}

// Create any update operations necessary to make the state match the plan
func createVirtualAttributeOperationsDefault(plan defaultVirtualAttributeResourceModel, state defaultVirtualAttributeResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.ScriptClass, state.ScriptClass, "script-class")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowRetrievingMembership, state.AllowRetrievingMembership, "allow-retrieving-membership")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScriptArgument, state.ScriptArgument, "script-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.JoinSourceAttribute, state.JoinSourceAttribute, "join-source-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.JoinTargetAttribute, state.JoinTargetAttribute, "join-target-attribute")
	operations.AddBoolOperationIfNecessary(&ops, plan.JoinMatchAll, state.JoinMatchAll, "join-match-all")
	operations.AddStringOperationIfNecessary(&ops, plan.SequenceNumberAttribute, state.SequenceNumberAttribute, "sequence-number-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.Value, state.Value, "value")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ReferencedByAttribute, state.ReferencedByAttribute, "referenced-by-attribute")
	operations.AddBoolOperationIfNecessary(&ops, plan.ReturnUtcTime, state.ReturnUtcTime, "return-utc-time")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeMilliseconds, state.IncludeMilliseconds, "include-milliseconds")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ReferenceSearchBaseDN, state.ReferenceSearchBaseDN, "reference-search-base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.JoinDNAttribute, state.JoinDNAttribute, "join-dn-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.JoinBaseDNType, state.JoinBaseDNType, "join-base-dn-type")
	operations.AddStringOperationIfNecessary(&ops, plan.JoinCustomBaseDN, state.JoinCustomBaseDN, "join-custom-base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.JoinScope, state.JoinScope, "join-scope")
	operations.AddInt64OperationIfNecessary(&ops, plan.JoinSizeLimit, state.JoinSizeLimit, "join-size-limit")
	operations.AddStringOperationIfNecessary(&ops, plan.JoinFilter, state.JoinFilter, "join-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.JoinAttribute, state.JoinAttribute, "join-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ValuePattern, state.ValuePattern, "value-pattern")
	operations.AddStringOperationIfNecessary(&ops, plan.ConflictBehavior, state.ConflictBehavior, "conflict-behavior")
	operations.AddBoolOperationIfNecessary(&ops, plan.DirectMembershipsOnly, state.DirectMembershipsOnly, "direct-memberships-only")
	operations.AddStringOperationIfNecessary(&ops, plan.IncludedGroupFilter, state.IncludedGroupFilter, "included-group-filter")
	operations.AddStringOperationIfNecessary(&ops, plan.RewriteSearchFilters, state.RewriteSearchFilters, "rewrite-search-filters")
	operations.AddStringOperationIfNecessary(&ops, plan.SourceAttribute, state.SourceAttribute, "source-attribute")
	operations.AddBoolOperationIfNecessary(&ops, plan.ExcludeOperationalAttributes, state.ExcludeOperationalAttributes, "exclude-operational-attributes")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedAttribute, state.ExcludedAttribute, "excluded-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.SourceEntryDNAttribute, state.SourceEntryDNAttribute, "source-entry-dn-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.SourceEntryDNMap, state.SourceEntryDNMap, "source-entry-dn-map")
	operations.AddBoolOperationIfNecessary(&ops, plan.BypassAccessControlForSearches, state.BypassAccessControlForSearches, "bypass-access-control-for-searches")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.AttributeType, state.AttributeType, "attribute-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.GroupDN, state.GroupDN, "group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.Filter, state.Filter, "filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ClientConnectionPolicy, state.ClientConnectionPolicy, "client-connection-policy")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequireExplicitRequestByName, state.RequireExplicitRequestByName, "require-explicit-request-by-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.MultipleVirtualAttributeEvaluationOrderIndex, state.MultipleVirtualAttributeEvaluationOrderIndex, "multiple-virtual-attribute-evaluation-order-index")
	operations.AddStringOperationIfNecessary(&ops, plan.MultipleVirtualAttributeMergeBehavior, state.MultipleVirtualAttributeMergeBehavior, "multiple-virtual-attribute-merge-behavior")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowIndexConflicts, state.AllowIndexConflicts, "allow-index-conflicts")
	return ops
}

// Create a mirror virtual-attribute
func (r *virtualAttributeResource) CreateMirrorVirtualAttribute(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan virtualAttributeResourceModel) (*virtualAttributeResourceModel, error) {
	addRequest := client.NewAddMirrorVirtualAttributeRequest(plan.Name.ValueString(),
		[]client.EnummirrorVirtualAttributeSchemaUrn{client.ENUMMIRRORVIRTUALATTRIBUTESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0VIRTUAL_ATTRIBUTEMIRROR},
		plan.SourceAttribute.ValueString(),
		plan.Enabled.ValueBool(),
		plan.AttributeType.ValueString())
	err := addOptionalMirrorVirtualAttributeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Virtual Attribute", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VirtualAttributeApi.AddVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddVirtualAttributeRequest(
		client.AddMirrorVirtualAttributeRequestAsAddVirtualAttributeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.VirtualAttributeApi.AddVirtualAttributeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Virtual Attribute", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state virtualAttributeResourceModel
	readMirrorVirtualAttributeResponse(ctx, addResponse.MirrorVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a constructed virtual-attribute
func (r *virtualAttributeResource) CreateConstructedVirtualAttribute(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan virtualAttributeResourceModel) (*virtualAttributeResourceModel, error) {
	var ValuePatternSlice []string
	plan.ValuePattern.ElementsAs(ctx, &ValuePatternSlice, false)
	addRequest := client.NewAddConstructedVirtualAttributeRequest(plan.Name.ValueString(),
		[]client.EnumconstructedVirtualAttributeSchemaUrn{client.ENUMCONSTRUCTEDVIRTUALATTRIBUTESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0VIRTUAL_ATTRIBUTECONSTRUCTED},
		ValuePatternSlice,
		plan.Enabled.ValueBool(),
		plan.AttributeType.ValueString())
	err := addOptionalConstructedVirtualAttributeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Virtual Attribute", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VirtualAttributeApi.AddVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddVirtualAttributeRequest(
		client.AddConstructedVirtualAttributeRequestAsAddVirtualAttributeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.VirtualAttributeApi.AddVirtualAttributeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Virtual Attribute", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state virtualAttributeResourceModel
	readConstructedVirtualAttributeResponse(ctx, addResponse.ConstructedVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a is-member-of virtual-attribute
func (r *virtualAttributeResource) CreateIsMemberOfVirtualAttribute(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan virtualAttributeResourceModel) (*virtualAttributeResourceModel, error) {
	addRequest := client.NewAddIsMemberOfVirtualAttributeRequest(plan.Name.ValueString(),
		[]client.EnumisMemberOfVirtualAttributeSchemaUrn{client.ENUMISMEMBEROFVIRTUALATTRIBUTESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0VIRTUAL_ATTRIBUTEIS_MEMBER_OF},
		plan.Enabled.ValueBool())
	err := addOptionalIsMemberOfVirtualAttributeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Virtual Attribute", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VirtualAttributeApi.AddVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddVirtualAttributeRequest(
		client.AddIsMemberOfVirtualAttributeRequestAsAddVirtualAttributeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.VirtualAttributeApi.AddVirtualAttributeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Virtual Attribute", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state virtualAttributeResourceModel
	readIsMemberOfVirtualAttributeResponse(ctx, addResponse.IsMemberOfVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a reverse-dn-join virtual-attribute
func (r *virtualAttributeResource) CreateReverseDnJoinVirtualAttribute(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan virtualAttributeResourceModel) (*virtualAttributeResourceModel, error) {
	joinBaseDNType, err := client.NewEnumvirtualAttributeJoinBaseDNTypePropFromValue(plan.JoinBaseDNType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for JoinBaseDNType", err.Error())
		return nil, err
	}
	addRequest := client.NewAddReverseDnJoinVirtualAttributeRequest(plan.Name.ValueString(),
		[]client.EnumreverseDnJoinVirtualAttributeSchemaUrn{client.ENUMREVERSEDNJOINVIRTUALATTRIBUTESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0VIRTUAL_ATTRIBUTEREVERSE_DN_JOIN},
		plan.JoinDNAttribute.ValueString(),
		*joinBaseDNType,
		plan.Enabled.ValueBool(),
		plan.AttributeType.ValueString())
	err = addOptionalReverseDnJoinVirtualAttributeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Virtual Attribute", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VirtualAttributeApi.AddVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddVirtualAttributeRequest(
		client.AddReverseDnJoinVirtualAttributeRequestAsAddVirtualAttributeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.VirtualAttributeApi.AddVirtualAttributeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Virtual Attribute", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state virtualAttributeResourceModel
	readReverseDnJoinVirtualAttributeResponse(ctx, addResponse.ReverseDnJoinVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a identify-references virtual-attribute
func (r *virtualAttributeResource) CreateIdentifyReferencesVirtualAttribute(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan virtualAttributeResourceModel) (*virtualAttributeResourceModel, error) {
	var ReferencedByAttributeSlice []string
	plan.ReferencedByAttribute.ElementsAs(ctx, &ReferencedByAttributeSlice, false)
	addRequest := client.NewAddIdentifyReferencesVirtualAttributeRequest(plan.Name.ValueString(),
		[]client.EnumidentifyReferencesVirtualAttributeSchemaUrn{client.ENUMIDENTIFYREFERENCESVIRTUALATTRIBUTESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0VIRTUAL_ATTRIBUTEIDENTIFY_REFERENCES},
		ReferencedByAttributeSlice,
		plan.Enabled.ValueBool(),
		plan.AttributeType.ValueString())
	err := addOptionalIdentifyReferencesVirtualAttributeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Virtual Attribute", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VirtualAttributeApi.AddVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddVirtualAttributeRequest(
		client.AddIdentifyReferencesVirtualAttributeRequestAsAddVirtualAttributeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.VirtualAttributeApi.AddVirtualAttributeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Virtual Attribute", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state virtualAttributeResourceModel
	readIdentifyReferencesVirtualAttributeResponse(ctx, addResponse.IdentifyReferencesVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a user-defined virtual-attribute
func (r *virtualAttributeResource) CreateUserDefinedVirtualAttribute(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan virtualAttributeResourceModel) (*virtualAttributeResourceModel, error) {
	var ValueSlice []string
	plan.Value.ElementsAs(ctx, &ValueSlice, false)
	addRequest := client.NewAddUserDefinedVirtualAttributeRequest(plan.Name.ValueString(),
		[]client.EnumuserDefinedVirtualAttributeSchemaUrn{client.ENUMUSERDEFINEDVIRTUALATTRIBUTESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0VIRTUAL_ATTRIBUTEUSER_DEFINED},
		ValueSlice,
		plan.Enabled.ValueBool(),
		plan.AttributeType.ValueString())
	err := addOptionalUserDefinedVirtualAttributeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Virtual Attribute", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VirtualAttributeApi.AddVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddVirtualAttributeRequest(
		client.AddUserDefinedVirtualAttributeRequestAsAddVirtualAttributeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.VirtualAttributeApi.AddVirtualAttributeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Virtual Attribute", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state virtualAttributeResourceModel
	readUserDefinedVirtualAttributeResponse(ctx, addResponse.UserDefinedVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a entry-dn virtual-attribute
func (r *virtualAttributeResource) CreateEntryDnVirtualAttribute(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan virtualAttributeResourceModel) (*virtualAttributeResourceModel, error) {
	addRequest := client.NewAddEntryDnVirtualAttributeRequest(plan.Name.ValueString(),
		[]client.EnumentryDnVirtualAttributeSchemaUrn{client.ENUMENTRYDNVIRTUALATTRIBUTESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0VIRTUAL_ATTRIBUTEENTRY_DN},
		plan.Enabled.ValueBool())
	err := addOptionalEntryDnVirtualAttributeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Virtual Attribute", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VirtualAttributeApi.AddVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddVirtualAttributeRequest(
		client.AddEntryDnVirtualAttributeRequestAsAddVirtualAttributeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.VirtualAttributeApi.AddVirtualAttributeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Virtual Attribute", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state virtualAttributeResourceModel
	readEntryDnVirtualAttributeResponse(ctx, addResponse.EntryDnVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a equality-join virtual-attribute
func (r *virtualAttributeResource) CreateEqualityJoinVirtualAttribute(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan virtualAttributeResourceModel) (*virtualAttributeResourceModel, error) {
	joinBaseDNType, err := client.NewEnumvirtualAttributeJoinBaseDNTypePropFromValue(plan.JoinBaseDNType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for JoinBaseDNType", err.Error())
		return nil, err
	}
	addRequest := client.NewAddEqualityJoinVirtualAttributeRequest(plan.Name.ValueString(),
		[]client.EnumequalityJoinVirtualAttributeSchemaUrn{client.ENUMEQUALITYJOINVIRTUALATTRIBUTESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0VIRTUAL_ATTRIBUTEEQUALITY_JOIN},
		plan.JoinSourceAttribute.ValueString(),
		plan.JoinTargetAttribute.ValueString(),
		*joinBaseDNType,
		plan.Enabled.ValueBool(),
		plan.AttributeType.ValueString())
	err = addOptionalEqualityJoinVirtualAttributeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Virtual Attribute", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VirtualAttributeApi.AddVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddVirtualAttributeRequest(
		client.AddEqualityJoinVirtualAttributeRequestAsAddVirtualAttributeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.VirtualAttributeApi.AddVirtualAttributeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Virtual Attribute", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state virtualAttributeResourceModel
	readEqualityJoinVirtualAttributeResponse(ctx, addResponse.EqualityJoinVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a groovy-scripted virtual-attribute
func (r *virtualAttributeResource) CreateGroovyScriptedVirtualAttribute(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan virtualAttributeResourceModel) (*virtualAttributeResourceModel, error) {
	addRequest := client.NewAddGroovyScriptedVirtualAttributeRequest(plan.Name.ValueString(),
		[]client.EnumgroovyScriptedVirtualAttributeSchemaUrn{client.ENUMGROOVYSCRIPTEDVIRTUALATTRIBUTESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0VIRTUAL_ATTRIBUTEGROOVY_SCRIPTED},
		plan.ScriptClass.ValueString(),
		plan.Enabled.ValueBool(),
		plan.AttributeType.ValueString())
	err := addOptionalGroovyScriptedVirtualAttributeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Virtual Attribute", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VirtualAttributeApi.AddVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddVirtualAttributeRequest(
		client.AddGroovyScriptedVirtualAttributeRequestAsAddVirtualAttributeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.VirtualAttributeApi.AddVirtualAttributeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Virtual Attribute", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state virtualAttributeResourceModel
	readGroovyScriptedVirtualAttributeResponse(ctx, addResponse.GroovyScriptedVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a member virtual-attribute
func (r *virtualAttributeResource) CreateMemberVirtualAttribute(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan virtualAttributeResourceModel) (*virtualAttributeResourceModel, error) {
	addRequest := client.NewAddMemberVirtualAttributeRequest(plan.Name.ValueString(),
		[]client.EnummemberVirtualAttributeSchemaUrn{client.ENUMMEMBERVIRTUALATTRIBUTESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0VIRTUAL_ATTRIBUTEMEMBER},
		plan.Enabled.ValueBool(),
		plan.AttributeType.ValueString())
	err := addOptionalMemberVirtualAttributeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Virtual Attribute", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VirtualAttributeApi.AddVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddVirtualAttributeRequest(
		client.AddMemberVirtualAttributeRequestAsAddVirtualAttributeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.VirtualAttributeApi.AddVirtualAttributeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Virtual Attribute", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state virtualAttributeResourceModel
	readMemberVirtualAttributeResponse(ctx, addResponse.MemberVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a password-policy-state-json virtual-attribute
func (r *virtualAttributeResource) CreatePasswordPolicyStateJsonVirtualAttribute(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan virtualAttributeResourceModel) (*virtualAttributeResourceModel, error) {
	addRequest := client.NewAddPasswordPolicyStateJsonVirtualAttributeRequest(plan.Name.ValueString(),
		[]client.EnumpasswordPolicyStateJsonVirtualAttributeSchemaUrn{client.ENUMPASSWORDPOLICYSTATEJSONVIRTUALATTRIBUTESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0VIRTUAL_ATTRIBUTEPASSWORD_POLICY_STATE_JSON},
		plan.Enabled.ValueBool())
	err := addOptionalPasswordPolicyStateJsonVirtualAttributeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Virtual Attribute", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VirtualAttributeApi.AddVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddVirtualAttributeRequest(
		client.AddPasswordPolicyStateJsonVirtualAttributeRequestAsAddVirtualAttributeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.VirtualAttributeApi.AddVirtualAttributeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Virtual Attribute", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state virtualAttributeResourceModel
	readPasswordPolicyStateJsonVirtualAttributeResponse(ctx, addResponse.PasswordPolicyStateJsonVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a dn-join virtual-attribute
func (r *virtualAttributeResource) CreateDnJoinVirtualAttribute(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan virtualAttributeResourceModel) (*virtualAttributeResourceModel, error) {
	joinBaseDNType, err := client.NewEnumvirtualAttributeJoinBaseDNTypePropFromValue(plan.JoinBaseDNType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for JoinBaseDNType", err.Error())
		return nil, err
	}
	addRequest := client.NewAddDnJoinVirtualAttributeRequest(plan.Name.ValueString(),
		[]client.EnumdnJoinVirtualAttributeSchemaUrn{client.ENUMDNJOINVIRTUALATTRIBUTESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0VIRTUAL_ATTRIBUTEDN_JOIN},
		plan.JoinDNAttribute.ValueString(),
		*joinBaseDNType,
		plan.Enabled.ValueBool(),
		plan.AttributeType.ValueString())
	err = addOptionalDnJoinVirtualAttributeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Virtual Attribute", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VirtualAttributeApi.AddVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddVirtualAttributeRequest(
		client.AddDnJoinVirtualAttributeRequestAsAddVirtualAttributeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.VirtualAttributeApi.AddVirtualAttributeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Virtual Attribute", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state virtualAttributeResourceModel
	readDnJoinVirtualAttributeResponse(ctx, addResponse.DnJoinVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party virtual-attribute
func (r *virtualAttributeResource) CreateThirdPartyVirtualAttribute(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan virtualAttributeResourceModel) (*virtualAttributeResourceModel, error) {
	addRequest := client.NewAddThirdPartyVirtualAttributeRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyVirtualAttributeSchemaUrn{client.ENUMTHIRDPARTYVIRTUALATTRIBUTESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0VIRTUAL_ATTRIBUTETHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool(),
		plan.AttributeType.ValueString())
	err := addOptionalThirdPartyVirtualAttributeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Virtual Attribute", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VirtualAttributeApi.AddVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddVirtualAttributeRequest(
		client.AddThirdPartyVirtualAttributeRequestAsAddVirtualAttributeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.VirtualAttributeApi.AddVirtualAttributeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Virtual Attribute", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state virtualAttributeResourceModel
	readThirdPartyVirtualAttributeResponse(ctx, addResponse.ThirdPartyVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *virtualAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan virtualAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *virtualAttributeResourceModel
	var err error
	if plan.Type.ValueString() == "mirror" {
		state, err = r.CreateMirrorVirtualAttribute(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "constructed" {
		state, err = r.CreateConstructedVirtualAttribute(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "is-member-of" {
		state, err = r.CreateIsMemberOfVirtualAttribute(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "reverse-dn-join" {
		state, err = r.CreateReverseDnJoinVirtualAttribute(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "identify-references" {
		state, err = r.CreateIdentifyReferencesVirtualAttribute(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "user-defined" {
		state, err = r.CreateUserDefinedVirtualAttribute(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "entry-dn" {
		state, err = r.CreateEntryDnVirtualAttribute(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "equality-join" {
		state, err = r.CreateEqualityJoinVirtualAttribute(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "groovy-scripted" {
		state, err = r.CreateGroovyScriptedVirtualAttribute(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "member" {
		state, err = r.CreateMemberVirtualAttribute(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "password-policy-state-json" {
		state, err = r.CreatePasswordPolicyStateJsonVirtualAttribute(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "dn-join" {
		state, err = r.CreateDnJoinVirtualAttribute(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyVirtualAttribute(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, *state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultVirtualAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan defaultVirtualAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.VirtualAttributeApi.GetVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Virtual Attribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state defaultVirtualAttributeResourceModel
	if readResponse.MirrorVirtualAttributeResponse != nil {
		readMirrorVirtualAttributeResponseDefault(ctx, readResponse.MirrorVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.EntryChecksumVirtualAttributeResponse != nil {
		readEntryChecksumVirtualAttributeResponseDefault(ctx, readResponse.EntryChecksumVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.MemberOfServerGroupVirtualAttributeResponse != nil {
		readMemberOfServerGroupVirtualAttributeResponseDefault(ctx, readResponse.MemberOfServerGroupVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ConstructedVirtualAttributeResponse != nil {
		readConstructedVirtualAttributeResponseDefault(ctx, readResponse.ConstructedVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.IsMemberOfVirtualAttributeResponse != nil {
		readIsMemberOfVirtualAttributeResponseDefault(ctx, readResponse.IsMemberOfVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CustomVirtualAttributeResponse != nil {
		readCustomVirtualAttributeResponseDefault(ctx, readResponse.CustomVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.NumSubordinatesVirtualAttributeResponse != nil {
		readNumSubordinatesVirtualAttributeResponseDefault(ctx, readResponse.NumSubordinatesVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ReverseDnJoinVirtualAttributeResponse != nil {
		readReverseDnJoinVirtualAttributeResponseDefault(ctx, readResponse.ReverseDnJoinVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.IdentifyReferencesVirtualAttributeResponse != nil {
		readIdentifyReferencesVirtualAttributeResponseDefault(ctx, readResponse.IdentifyReferencesVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.UserDefinedVirtualAttributeResponse != nil {
		readUserDefinedVirtualAttributeResponseDefault(ctx, readResponse.UserDefinedVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CurrentTimeVirtualAttributeResponse != nil {
		readCurrentTimeVirtualAttributeResponseDefault(ctx, readResponse.CurrentTimeVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ShortUniqueIdVirtualAttributeResponse != nil {
		readShortUniqueIdVirtualAttributeResponseDefault(ctx, readResponse.ShortUniqueIdVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.EntryDnVirtualAttributeResponse != nil {
		readEntryDnVirtualAttributeResponseDefault(ctx, readResponse.EntryDnVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.HasSubordinatesVirtualAttributeResponse != nil {
		readHasSubordinatesVirtualAttributeResponseDefault(ctx, readResponse.HasSubordinatesVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.EqualityJoinVirtualAttributeResponse != nil {
		readEqualityJoinVirtualAttributeResponseDefault(ctx, readResponse.EqualityJoinVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedVirtualAttributeResponse != nil {
		readGroovyScriptedVirtualAttributeResponseDefault(ctx, readResponse.GroovyScriptedVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.InstanceNameVirtualAttributeResponse != nil {
		readInstanceNameVirtualAttributeResponseDefault(ctx, readResponse.InstanceNameVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ReplicationStateDetailVirtualAttributeResponse != nil {
		readReplicationStateDetailVirtualAttributeResponseDefault(ctx, readResponse.ReplicationStateDetailVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.MemberVirtualAttributeResponse != nil {
		readMemberVirtualAttributeResponseDefault(ctx, readResponse.MemberVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PasswordPolicyStateJsonVirtualAttributeResponse != nil {
		readPasswordPolicyStateJsonVirtualAttributeResponseDefault(ctx, readResponse.PasswordPolicyStateJsonVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SubschemaSubentryVirtualAttributeResponse != nil {
		readSubschemaSubentryVirtualAttributeResponseDefault(ctx, readResponse.SubschemaSubentryVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DnJoinVirtualAttributeResponse != nil {
		readDnJoinVirtualAttributeResponseDefault(ctx, readResponse.DnJoinVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyVirtualAttributeResponse != nil {
		readThirdPartyVirtualAttributeResponseDefault(ctx, readResponse.ThirdPartyVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.VirtualAttributeApi.UpdateVirtualAttribute(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createVirtualAttributeOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.VirtualAttributeApi.UpdateVirtualAttributeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Virtual Attribute", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.MirrorVirtualAttributeResponse != nil {
			readMirrorVirtualAttributeResponseDefault(ctx, updateResponse.MirrorVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.EntryChecksumVirtualAttributeResponse != nil {
			readEntryChecksumVirtualAttributeResponseDefault(ctx, updateResponse.EntryChecksumVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.MemberOfServerGroupVirtualAttributeResponse != nil {
			readMemberOfServerGroupVirtualAttributeResponseDefault(ctx, updateResponse.MemberOfServerGroupVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConstructedVirtualAttributeResponse != nil {
			readConstructedVirtualAttributeResponseDefault(ctx, updateResponse.ConstructedVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.IsMemberOfVirtualAttributeResponse != nil {
			readIsMemberOfVirtualAttributeResponseDefault(ctx, updateResponse.IsMemberOfVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CustomVirtualAttributeResponse != nil {
			readCustomVirtualAttributeResponseDefault(ctx, updateResponse.CustomVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.NumSubordinatesVirtualAttributeResponse != nil {
			readNumSubordinatesVirtualAttributeResponseDefault(ctx, updateResponse.NumSubordinatesVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ReverseDnJoinVirtualAttributeResponse != nil {
			readReverseDnJoinVirtualAttributeResponseDefault(ctx, updateResponse.ReverseDnJoinVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.IdentifyReferencesVirtualAttributeResponse != nil {
			readIdentifyReferencesVirtualAttributeResponseDefault(ctx, updateResponse.IdentifyReferencesVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.UserDefinedVirtualAttributeResponse != nil {
			readUserDefinedVirtualAttributeResponseDefault(ctx, updateResponse.UserDefinedVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CurrentTimeVirtualAttributeResponse != nil {
			readCurrentTimeVirtualAttributeResponseDefault(ctx, updateResponse.CurrentTimeVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ShortUniqueIdVirtualAttributeResponse != nil {
			readShortUniqueIdVirtualAttributeResponseDefault(ctx, updateResponse.ShortUniqueIdVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.EntryDnVirtualAttributeResponse != nil {
			readEntryDnVirtualAttributeResponseDefault(ctx, updateResponse.EntryDnVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.HasSubordinatesVirtualAttributeResponse != nil {
			readHasSubordinatesVirtualAttributeResponseDefault(ctx, updateResponse.HasSubordinatesVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.EqualityJoinVirtualAttributeResponse != nil {
			readEqualityJoinVirtualAttributeResponseDefault(ctx, updateResponse.EqualityJoinVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedVirtualAttributeResponse != nil {
			readGroovyScriptedVirtualAttributeResponseDefault(ctx, updateResponse.GroovyScriptedVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.InstanceNameVirtualAttributeResponse != nil {
			readInstanceNameVirtualAttributeResponseDefault(ctx, updateResponse.InstanceNameVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ReplicationStateDetailVirtualAttributeResponse != nil {
			readReplicationStateDetailVirtualAttributeResponseDefault(ctx, updateResponse.ReplicationStateDetailVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.MemberVirtualAttributeResponse != nil {
			readMemberVirtualAttributeResponseDefault(ctx, updateResponse.MemberVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PasswordPolicyStateJsonVirtualAttributeResponse != nil {
			readPasswordPolicyStateJsonVirtualAttributeResponseDefault(ctx, updateResponse.PasswordPolicyStateJsonVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SubschemaSubentryVirtualAttributeResponse != nil {
			readSubschemaSubentryVirtualAttributeResponseDefault(ctx, updateResponse.SubschemaSubentryVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DnJoinVirtualAttributeResponse != nil {
			readDnJoinVirtualAttributeResponseDefault(ctx, updateResponse.DnJoinVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyVirtualAttributeResponse != nil {
			readThirdPartyVirtualAttributeResponseDefault(ctx, updateResponse.ThirdPartyVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *virtualAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state virtualAttributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.VirtualAttributeApi.GetVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Virtual Attribute", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Virtual Attribute", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.MirrorVirtualAttributeResponse != nil {
		readMirrorVirtualAttributeResponse(ctx, readResponse.MirrorVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ConstructedVirtualAttributeResponse != nil {
		readConstructedVirtualAttributeResponse(ctx, readResponse.ConstructedVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.IsMemberOfVirtualAttributeResponse != nil {
		readIsMemberOfVirtualAttributeResponse(ctx, readResponse.IsMemberOfVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ReverseDnJoinVirtualAttributeResponse != nil {
		readReverseDnJoinVirtualAttributeResponse(ctx, readResponse.ReverseDnJoinVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.IdentifyReferencesVirtualAttributeResponse != nil {
		readIdentifyReferencesVirtualAttributeResponse(ctx, readResponse.IdentifyReferencesVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.UserDefinedVirtualAttributeResponse != nil {
		readUserDefinedVirtualAttributeResponse(ctx, readResponse.UserDefinedVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.EntryDnVirtualAttributeResponse != nil {
		readEntryDnVirtualAttributeResponse(ctx, readResponse.EntryDnVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.EqualityJoinVirtualAttributeResponse != nil {
		readEqualityJoinVirtualAttributeResponse(ctx, readResponse.EqualityJoinVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedVirtualAttributeResponse != nil {
		readGroovyScriptedVirtualAttributeResponse(ctx, readResponse.GroovyScriptedVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.MemberVirtualAttributeResponse != nil {
		readMemberVirtualAttributeResponse(ctx, readResponse.MemberVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PasswordPolicyStateJsonVirtualAttributeResponse != nil {
		readPasswordPolicyStateJsonVirtualAttributeResponse(ctx, readResponse.PasswordPolicyStateJsonVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DnJoinVirtualAttributeResponse != nil {
		readDnJoinVirtualAttributeResponse(ctx, readResponse.DnJoinVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyVirtualAttributeResponse != nil {
		readThirdPartyVirtualAttributeResponse(ctx, readResponse.ThirdPartyVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *defaultVirtualAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state defaultVirtualAttributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.VirtualAttributeApi.GetVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
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
	if readResponse.EntryChecksumVirtualAttributeResponse != nil {
		readEntryChecksumVirtualAttributeResponseDefault(ctx, readResponse.EntryChecksumVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.MemberOfServerGroupVirtualAttributeResponse != nil {
		readMemberOfServerGroupVirtualAttributeResponseDefault(ctx, readResponse.MemberOfServerGroupVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CustomVirtualAttributeResponse != nil {
		readCustomVirtualAttributeResponseDefault(ctx, readResponse.CustomVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.NumSubordinatesVirtualAttributeResponse != nil {
		readNumSubordinatesVirtualAttributeResponseDefault(ctx, readResponse.NumSubordinatesVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CurrentTimeVirtualAttributeResponse != nil {
		readCurrentTimeVirtualAttributeResponseDefault(ctx, readResponse.CurrentTimeVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ShortUniqueIdVirtualAttributeResponse != nil {
		readShortUniqueIdVirtualAttributeResponseDefault(ctx, readResponse.ShortUniqueIdVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.HasSubordinatesVirtualAttributeResponse != nil {
		readHasSubordinatesVirtualAttributeResponseDefault(ctx, readResponse.HasSubordinatesVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.InstanceNameVirtualAttributeResponse != nil {
		readInstanceNameVirtualAttributeResponseDefault(ctx, readResponse.InstanceNameVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ReplicationStateDetailVirtualAttributeResponse != nil {
		readReplicationStateDetailVirtualAttributeResponseDefault(ctx, readResponse.ReplicationStateDetailVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SubschemaSubentryVirtualAttributeResponse != nil {
		readSubschemaSubentryVirtualAttributeResponseDefault(ctx, readResponse.SubschemaSubentryVirtualAttributeResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *virtualAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan virtualAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state virtualAttributeResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.VirtualAttributeApi.UpdateVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createVirtualAttributeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.VirtualAttributeApi.UpdateVirtualAttributeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Virtual Attribute", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.MirrorVirtualAttributeResponse != nil {
			readMirrorVirtualAttributeResponse(ctx, updateResponse.MirrorVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConstructedVirtualAttributeResponse != nil {
			readConstructedVirtualAttributeResponse(ctx, updateResponse.ConstructedVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.IsMemberOfVirtualAttributeResponse != nil {
			readIsMemberOfVirtualAttributeResponse(ctx, updateResponse.IsMemberOfVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ReverseDnJoinVirtualAttributeResponse != nil {
			readReverseDnJoinVirtualAttributeResponse(ctx, updateResponse.ReverseDnJoinVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.IdentifyReferencesVirtualAttributeResponse != nil {
			readIdentifyReferencesVirtualAttributeResponse(ctx, updateResponse.IdentifyReferencesVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.UserDefinedVirtualAttributeResponse != nil {
			readUserDefinedVirtualAttributeResponse(ctx, updateResponse.UserDefinedVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.EntryDnVirtualAttributeResponse != nil {
			readEntryDnVirtualAttributeResponse(ctx, updateResponse.EntryDnVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.EqualityJoinVirtualAttributeResponse != nil {
			readEqualityJoinVirtualAttributeResponse(ctx, updateResponse.EqualityJoinVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedVirtualAttributeResponse != nil {
			readGroovyScriptedVirtualAttributeResponse(ctx, updateResponse.GroovyScriptedVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.MemberVirtualAttributeResponse != nil {
			readMemberVirtualAttributeResponse(ctx, updateResponse.MemberVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PasswordPolicyStateJsonVirtualAttributeResponse != nil {
			readPasswordPolicyStateJsonVirtualAttributeResponse(ctx, updateResponse.PasswordPolicyStateJsonVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DnJoinVirtualAttributeResponse != nil {
			readDnJoinVirtualAttributeResponse(ctx, updateResponse.DnJoinVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyVirtualAttributeResponse != nil {
			readThirdPartyVirtualAttributeResponse(ctx, updateResponse.ThirdPartyVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *defaultVirtualAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan defaultVirtualAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state defaultVirtualAttributeResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.VirtualAttributeApi.UpdateVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createVirtualAttributeOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.VirtualAttributeApi.UpdateVirtualAttributeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Virtual Attribute", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.MirrorVirtualAttributeResponse != nil {
			readMirrorVirtualAttributeResponseDefault(ctx, updateResponse.MirrorVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.EntryChecksumVirtualAttributeResponse != nil {
			readEntryChecksumVirtualAttributeResponseDefault(ctx, updateResponse.EntryChecksumVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.MemberOfServerGroupVirtualAttributeResponse != nil {
			readMemberOfServerGroupVirtualAttributeResponseDefault(ctx, updateResponse.MemberOfServerGroupVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConstructedVirtualAttributeResponse != nil {
			readConstructedVirtualAttributeResponseDefault(ctx, updateResponse.ConstructedVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.IsMemberOfVirtualAttributeResponse != nil {
			readIsMemberOfVirtualAttributeResponseDefault(ctx, updateResponse.IsMemberOfVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CustomVirtualAttributeResponse != nil {
			readCustomVirtualAttributeResponseDefault(ctx, updateResponse.CustomVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.NumSubordinatesVirtualAttributeResponse != nil {
			readNumSubordinatesVirtualAttributeResponseDefault(ctx, updateResponse.NumSubordinatesVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ReverseDnJoinVirtualAttributeResponse != nil {
			readReverseDnJoinVirtualAttributeResponseDefault(ctx, updateResponse.ReverseDnJoinVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.IdentifyReferencesVirtualAttributeResponse != nil {
			readIdentifyReferencesVirtualAttributeResponseDefault(ctx, updateResponse.IdentifyReferencesVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.UserDefinedVirtualAttributeResponse != nil {
			readUserDefinedVirtualAttributeResponseDefault(ctx, updateResponse.UserDefinedVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CurrentTimeVirtualAttributeResponse != nil {
			readCurrentTimeVirtualAttributeResponseDefault(ctx, updateResponse.CurrentTimeVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ShortUniqueIdVirtualAttributeResponse != nil {
			readShortUniqueIdVirtualAttributeResponseDefault(ctx, updateResponse.ShortUniqueIdVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.EntryDnVirtualAttributeResponse != nil {
			readEntryDnVirtualAttributeResponseDefault(ctx, updateResponse.EntryDnVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.HasSubordinatesVirtualAttributeResponse != nil {
			readHasSubordinatesVirtualAttributeResponseDefault(ctx, updateResponse.HasSubordinatesVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.EqualityJoinVirtualAttributeResponse != nil {
			readEqualityJoinVirtualAttributeResponseDefault(ctx, updateResponse.EqualityJoinVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedVirtualAttributeResponse != nil {
			readGroovyScriptedVirtualAttributeResponseDefault(ctx, updateResponse.GroovyScriptedVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.InstanceNameVirtualAttributeResponse != nil {
			readInstanceNameVirtualAttributeResponseDefault(ctx, updateResponse.InstanceNameVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ReplicationStateDetailVirtualAttributeResponse != nil {
			readReplicationStateDetailVirtualAttributeResponseDefault(ctx, updateResponse.ReplicationStateDetailVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.MemberVirtualAttributeResponse != nil {
			readMemberVirtualAttributeResponseDefault(ctx, updateResponse.MemberVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PasswordPolicyStateJsonVirtualAttributeResponse != nil {
			readPasswordPolicyStateJsonVirtualAttributeResponseDefault(ctx, updateResponse.PasswordPolicyStateJsonVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SubschemaSubentryVirtualAttributeResponse != nil {
			readSubschemaSubentryVirtualAttributeResponseDefault(ctx, updateResponse.SubschemaSubentryVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DnJoinVirtualAttributeResponse != nil {
			readDnJoinVirtualAttributeResponseDefault(ctx, updateResponse.DnJoinVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyVirtualAttributeResponse != nil {
			readThirdPartyVirtualAttributeResponseDefault(ctx, updateResponse.ThirdPartyVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// This config object is edit-only, so Terraform can't delete it.
// After running a delete, Terraform will just "forget" about this object and it can be managed elsewhere.
func (r *defaultVirtualAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *virtualAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state virtualAttributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.VirtualAttributeApi.DeleteVirtualAttributeExecute(r.apiClient.VirtualAttributeApi.DeleteVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Virtual Attribute", err, httpResp)
		return
	}
}

func (r *virtualAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importVirtualAttribute(ctx, req, resp)
}

func (r *defaultVirtualAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importVirtualAttribute(ctx, req, resp)
}

func importVirtualAttribute(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
