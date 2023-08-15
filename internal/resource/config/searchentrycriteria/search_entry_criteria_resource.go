package searchentrycriteria

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	_ resource.Resource                = &searchEntryCriteriaResource{}
	_ resource.ResourceWithConfigure   = &searchEntryCriteriaResource{}
	_ resource.ResourceWithImportState = &searchEntryCriteriaResource{}
	_ resource.Resource                = &defaultSearchEntryCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultSearchEntryCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultSearchEntryCriteriaResource{}
)

// Create a Search Entry Criteria resource
func NewSearchEntryCriteriaResource() resource.Resource {
	return &searchEntryCriteriaResource{}
}

func NewDefaultSearchEntryCriteriaResource() resource.Resource {
	return &defaultSearchEntryCriteriaResource{}
}

// searchEntryCriteriaResource is the resource implementation.
type searchEntryCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSearchEntryCriteriaResource is the resource implementation.
type defaultSearchEntryCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *searchEntryCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_search_entry_criteria"
}

func (r *defaultSearchEntryCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_search_entry_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *searchEntryCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultSearchEntryCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type searchEntryCriteriaResourceModel struct {
	Id                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
	LastUpdated                       types.String `tfsdk:"last_updated"`
	Notifications                     types.Set    `tfsdk:"notifications"`
	RequiredActions                   types.Set    `tfsdk:"required_actions"`
	Type                              types.String `tfsdk:"type"`
	ExtensionClass                    types.String `tfsdk:"extension_class"`
	ExtensionArgument                 types.Set    `tfsdk:"extension_argument"`
	AllIncludedSearchEntryCriteria    types.Set    `tfsdk:"all_included_search_entry_criteria"`
	AnyIncludedSearchEntryCriteria    types.Set    `tfsdk:"any_included_search_entry_criteria"`
	NotAllIncludedSearchEntryCriteria types.Set    `tfsdk:"not_all_included_search_entry_criteria"`
	NoneIncludedSearchEntryCriteria   types.Set    `tfsdk:"none_included_search_entry_criteria"`
	RequestCriteria                   types.String `tfsdk:"request_criteria"`
	AllIncludedEntryControl           types.Set    `tfsdk:"all_included_entry_control"`
	AnyIncludedEntryControl           types.Set    `tfsdk:"any_included_entry_control"`
	NotAllIncludedEntryControl        types.Set    `tfsdk:"not_all_included_entry_control"`
	NoneIncludedEntryControl          types.Set    `tfsdk:"none_included_entry_control"`
	IncludedEntryBaseDN               types.Set    `tfsdk:"included_entry_base_dn"`
	ExcludedEntryBaseDN               types.Set    `tfsdk:"excluded_entry_base_dn"`
	AllIncludedEntryFilter            types.Set    `tfsdk:"all_included_entry_filter"`
	AnyIncludedEntryFilter            types.Set    `tfsdk:"any_included_entry_filter"`
	NotAllIncludedEntryFilter         types.Set    `tfsdk:"not_all_included_entry_filter"`
	NoneIncludedEntryFilter           types.Set    `tfsdk:"none_included_entry_filter"`
	AllIncludedEntryGroupDN           types.Set    `tfsdk:"all_included_entry_group_dn"`
	AnyIncludedEntryGroupDN           types.Set    `tfsdk:"any_included_entry_group_dn"`
	NotAllIncludedEntryGroupDN        types.Set    `tfsdk:"not_all_included_entry_group_dn"`
	NoneIncludedEntryGroupDN          types.Set    `tfsdk:"none_included_entry_group_dn"`
	Description                       types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *searchEntryCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	searchEntryCriteriaSchema(ctx, req, resp, false)
}

func (r *defaultSearchEntryCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	searchEntryCriteriaSchema(ctx, req, resp, true)
}

func searchEntryCriteriaSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Search Entry Criteria.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Search Entry Criteria resource. Options are ['simple', 'aggregate', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"simple", "aggregate", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Search Entry Criteria.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Search Entry Criteria. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"all_included_search_entry_criteria": schema.SetAttribute{
				Description: "Specifies a search entry criteria object that must match the associated search result entry in order to match the aggregate search entry criteria. If one or more all-included search entry criteria objects are provided, then a search result entry must match all of them in order to match the aggregate search entry criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"any_included_search_entry_criteria": schema.SetAttribute{
				Description: "Specifies a search entry criteria object that may match the associated search result entry in order to match the aggregate search entry criteria. If one or more any-included search entry criteria objects are provided, then a search result entry must match at least one of them in order to match the aggregate search entry criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"not_all_included_search_entry_criteria": schema.SetAttribute{
				Description: "Specifies a search entry criteria object that should not match the associated search result entry in order to match the aggregate search entry criteria. If one or more not-all-included search entry criteria objects are provided, then a search result entry must not match all of them (that is, it may match zero or more of them, but it must not match all of them) in order to match the aggregate search entry criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"none_included_search_entry_criteria": schema.SetAttribute{
				Description: "Specifies a search entry criteria object that must not match the associated search result entry in order to match the aggregate search entry criteria. If one or more none-included search entry criteria objects are provided, then a search result entry must not match any of them in order to match the aggregate search entry criteria.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"request_criteria": schema.StringAttribute{
				Description: "Specifies a request criteria object that must match the associated request for entries included in this Simple Search Entry Criteria. of them.",
				Optional:    true,
			},
			"all_included_entry_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must be present in search result entries included in this Simple Search Entry Criteria. If any control OIDs are provided, then the entry must contain all of those controls.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"any_included_entry_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that may be present in search result entries included in this Simple Search Entry Criteria. If any control OIDs are provided, then the entry must contain at least one of those controls.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"not_all_included_entry_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that should not be present in search result entries included in this Simple Search Entry Criteria. If any control OIDs are provided, then the entry must not contain at least one of those controls (that is, it may contain zero or more of those controls, but not all of them).",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"none_included_entry_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must not be present in search result entries included in this Simple Search Entry Criteria. If any control OIDs are provided, then the entry must not contain any of those controls.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"included_entry_base_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which entries included in this Simple Search Entry Criteria may exist.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"excluded_entry_base_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which entries included in this Simple Search Entry Criteria may not exist.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"all_included_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that must match search result entries included in this Simple Search Entry Criteria. Note that this matching will be performed against the entry that is actually returned to the client and may not reflect the complete entry stored in the server. If any filters are provided, then the returned entry must match all of those filters.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"any_included_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that may match search result entries included in this Simple Search Entry Criteria. Note that this matching will be performed against the entry that is actually returned to the client and may not reflect the complete entry stored in the server. If any filters are provided, then the entry must match at least one of those filters.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"not_all_included_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that should not match search result entries included in this Simple Search Entry Criteria. Note that this matching will be performed against the entry that is actually returned to the client and may not reflect the complete entry stored in the server. If any filters are provided, then the entry must not match at least one of those filters (that is, the entry may match zero or more of those filters, but not of all of them).",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"none_included_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that must not match search result entries included in this Simple Search Entry Criteria. Note that this matching will be performed against the entry that is actually returned to the client and may not reflect the complete entry stored in the server. If any filters are provided, then the entry must not match any of those filters.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"all_included_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the entry must be a member to be included in this Simple Search Entry Criteria. If any group DNs are provided, then the entry must be a member of all of them.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"any_included_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the entry may be a member to be included in this Simple Search Entry Criteria. If any group DNs are provided, then the entry must be a member of at least one of them.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"not_all_included_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the entry should not be a member to be included in this Simple Search Entry Criteria. If any group DNs are provided, then the entry must not be a member of at least one of them (that is, the entry may be a member of zero or more of the specified groups, but not of all of them).",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"none_included_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the entry must not be a member to be included in this Simple Search Entry Criteria. If any group DNs are provided, then the entry must not be a member of any of them.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Search Entry Criteria",
				Optional:    true,
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
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsSearchEntryCriteria() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("not_all_included_entry_control"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("none_included_entry_filter"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("not_all_included_entry_group_dn"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("any_included_search_entry_criteria"),
			path.MatchRoot("type"),
			[]string{"aggregate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("all_included_entry_filter"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("any_included_entry_filter"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("none_included_search_entry_criteria"),
			path.MatchRoot("type"),
			[]string{"aggregate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("any_included_entry_group_dn"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("none_included_entry_group_dn"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("any_included_entry_control"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("excluded_entry_base_dn"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("none_included_entry_control"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("all_included_search_entry_criteria"),
			path.MatchRoot("type"),
			[]string{"aggregate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("all_included_entry_control"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("request_criteria"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_argument"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_class"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("all_included_entry_group_dn"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("not_all_included_entry_filter"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("included_entry_base_dn"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("not_all_included_search_entry_criteria"),
			path.MatchRoot("type"),
			[]string{"aggregate"},
		),
	}
}

// Add config validators
func (r searchEntryCriteriaResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsSearchEntryCriteria()
}

// Add config validators
func (r defaultSearchEntryCriteriaResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsSearchEntryCriteria()
}

// Add optional fields to create request for simple search-entry-criteria
func addOptionalSimpleSearchEntryCriteriaFields(ctx context.Context, addRequest *client.AddSimpleSearchEntryCriteriaRequest, plan searchEntryCriteriaResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		addRequest.RequestCriteria = plan.RequestCriteria.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.AllIncludedEntryControl) {
		var slice []string
		plan.AllIncludedEntryControl.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedEntryControl = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedEntryControl) {
		var slice []string
		plan.AnyIncludedEntryControl.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedEntryControl = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedEntryControl) {
		var slice []string
		plan.NotAllIncludedEntryControl.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedEntryControl = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedEntryControl) {
		var slice []string
		plan.NoneIncludedEntryControl.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedEntryControl = slice
	}
	if internaltypes.IsDefined(plan.IncludedEntryBaseDN) {
		var slice []string
		plan.IncludedEntryBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.IncludedEntryBaseDN = slice
	}
	if internaltypes.IsDefined(plan.ExcludedEntryBaseDN) {
		var slice []string
		plan.ExcludedEntryBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedEntryBaseDN = slice
	}
	if internaltypes.IsDefined(plan.AllIncludedEntryFilter) {
		var slice []string
		plan.AllIncludedEntryFilter.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedEntryFilter = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedEntryFilter) {
		var slice []string
		plan.AnyIncludedEntryFilter.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedEntryFilter = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedEntryFilter) {
		var slice []string
		plan.NotAllIncludedEntryFilter.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedEntryFilter = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedEntryFilter) {
		var slice []string
		plan.NoneIncludedEntryFilter.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedEntryFilter = slice
	}
	if internaltypes.IsDefined(plan.AllIncludedEntryGroupDN) {
		var slice []string
		plan.AllIncludedEntryGroupDN.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedEntryGroupDN = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedEntryGroupDN) {
		var slice []string
		plan.AnyIncludedEntryGroupDN.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedEntryGroupDN = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedEntryGroupDN) {
		var slice []string
		plan.NotAllIncludedEntryGroupDN.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedEntryGroupDN = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedEntryGroupDN) {
		var slice []string
		plan.NoneIncludedEntryGroupDN.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedEntryGroupDN = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for aggregate search-entry-criteria
func addOptionalAggregateSearchEntryCriteriaFields(ctx context.Context, addRequest *client.AddAggregateSearchEntryCriteriaRequest, plan searchEntryCriteriaResourceModel) {
	if internaltypes.IsDefined(plan.AllIncludedSearchEntryCriteria) {
		var slice []string
		plan.AllIncludedSearchEntryCriteria.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedSearchEntryCriteria = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedSearchEntryCriteria) {
		var slice []string
		plan.AnyIncludedSearchEntryCriteria.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedSearchEntryCriteria = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedSearchEntryCriteria) {
		var slice []string
		plan.NotAllIncludedSearchEntryCriteria.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedSearchEntryCriteria = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedSearchEntryCriteria) {
		var slice []string
		plan.NoneIncludedSearchEntryCriteria.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedSearchEntryCriteria = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for third-party search-entry-criteria
func addOptionalThirdPartySearchEntryCriteriaFields(ctx context.Context, addRequest *client.AddThirdPartySearchEntryCriteriaRequest, plan searchEntryCriteriaResourceModel) {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateSearchEntryCriteriaUnknownValues(ctx context.Context, model *searchEntryCriteriaResourceModel) {
	if model.AllIncludedEntryGroupDN.ElementType(ctx) == nil {
		model.AllIncludedEntryGroupDN = types.SetNull(types.StringType)
	}
	if model.ExcludedEntryBaseDN.ElementType(ctx) == nil {
		model.ExcludedEntryBaseDN = types.SetNull(types.StringType)
	}
	if model.AnyIncludedSearchEntryCriteria.ElementType(ctx) == nil {
		model.AnyIncludedSearchEntryCriteria = types.SetNull(types.StringType)
	}
	if model.AnyIncludedEntryFilter.ElementType(ctx) == nil {
		model.AnyIncludedEntryFilter = types.SetNull(types.StringType)
	}
	if model.AllIncludedEntryControl.ElementType(ctx) == nil {
		model.AllIncludedEntryControl = types.SetNull(types.StringType)
	}
	if model.IncludedEntryBaseDN.ElementType(ctx) == nil {
		model.IncludedEntryBaseDN = types.SetNull(types.StringType)
	}
	if model.NoneIncludedSearchEntryCriteria.ElementType(ctx) == nil {
		model.NoneIncludedSearchEntryCriteria = types.SetNull(types.StringType)
	}
	if model.ExtensionArgument.ElementType(ctx) == nil {
		model.ExtensionArgument = types.SetNull(types.StringType)
	}
	if model.AnyIncludedEntryControl.ElementType(ctx) == nil {
		model.AnyIncludedEntryControl = types.SetNull(types.StringType)
	}
	if model.AllIncludedSearchEntryCriteria.ElementType(ctx) == nil {
		model.AllIncludedSearchEntryCriteria = types.SetNull(types.StringType)
	}
	if model.AllIncludedEntryFilter.ElementType(ctx) == nil {
		model.AllIncludedEntryFilter = types.SetNull(types.StringType)
	}
	if model.NoneIncludedEntryFilter.ElementType(ctx) == nil {
		model.NoneIncludedEntryFilter = types.SetNull(types.StringType)
	}
	if model.NoneIncludedEntryControl.ElementType(ctx) == nil {
		model.NoneIncludedEntryControl = types.SetNull(types.StringType)
	}
	if model.NotAllIncludedEntryGroupDN.ElementType(ctx) == nil {
		model.NotAllIncludedEntryGroupDN = types.SetNull(types.StringType)
	}
	if model.NotAllIncludedEntryControl.ElementType(ctx) == nil {
		model.NotAllIncludedEntryControl = types.SetNull(types.StringType)
	}
	if model.NotAllIncludedEntryFilter.ElementType(ctx) == nil {
		model.NotAllIncludedEntryFilter = types.SetNull(types.StringType)
	}
	if model.AnyIncludedEntryGroupDN.ElementType(ctx) == nil {
		model.AnyIncludedEntryGroupDN = types.SetNull(types.StringType)
	}
	if model.NotAllIncludedSearchEntryCriteria.ElementType(ctx) == nil {
		model.NotAllIncludedSearchEntryCriteria = types.SetNull(types.StringType)
	}
	if model.NoneIncludedEntryGroupDN.ElementType(ctx) == nil {
		model.NoneIncludedEntryGroupDN = types.SetNull(types.StringType)
	}
}

// Read a SimpleSearchEntryCriteriaResponse object into the model struct
func readSimpleSearchEntryCriteriaResponse(ctx context.Context, r *client.SimpleSearchEntryCriteriaResponse, state *searchEntryCriteriaResourceModel, expectedValues *searchEntryCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("simple")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.AllIncludedEntryControl = internaltypes.GetStringSet(r.AllIncludedEntryControl)
	state.AnyIncludedEntryControl = internaltypes.GetStringSet(r.AnyIncludedEntryControl)
	state.NotAllIncludedEntryControl = internaltypes.GetStringSet(r.NotAllIncludedEntryControl)
	state.NoneIncludedEntryControl = internaltypes.GetStringSet(r.NoneIncludedEntryControl)
	state.IncludedEntryBaseDN = internaltypes.GetStringSet(r.IncludedEntryBaseDN)
	state.ExcludedEntryBaseDN = internaltypes.GetStringSet(r.ExcludedEntryBaseDN)
	state.AllIncludedEntryFilter = internaltypes.GetStringSet(r.AllIncludedEntryFilter)
	state.AnyIncludedEntryFilter = internaltypes.GetStringSet(r.AnyIncludedEntryFilter)
	state.NotAllIncludedEntryFilter = internaltypes.GetStringSet(r.NotAllIncludedEntryFilter)
	state.NoneIncludedEntryFilter = internaltypes.GetStringSet(r.NoneIncludedEntryFilter)
	state.AllIncludedEntryGroupDN = internaltypes.GetStringSet(r.AllIncludedEntryGroupDN)
	state.AnyIncludedEntryGroupDN = internaltypes.GetStringSet(r.AnyIncludedEntryGroupDN)
	state.NotAllIncludedEntryGroupDN = internaltypes.GetStringSet(r.NotAllIncludedEntryGroupDN)
	state.NoneIncludedEntryGroupDN = internaltypes.GetStringSet(r.NoneIncludedEntryGroupDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSearchEntryCriteriaUnknownValues(ctx, state)
}

// Read a AggregateSearchEntryCriteriaResponse object into the model struct
func readAggregateSearchEntryCriteriaResponse(ctx context.Context, r *client.AggregateSearchEntryCriteriaResponse, state *searchEntryCriteriaResourceModel, expectedValues *searchEntryCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("aggregate")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AllIncludedSearchEntryCriteria = internaltypes.GetStringSet(r.AllIncludedSearchEntryCriteria)
	state.AnyIncludedSearchEntryCriteria = internaltypes.GetStringSet(r.AnyIncludedSearchEntryCriteria)
	state.NotAllIncludedSearchEntryCriteria = internaltypes.GetStringSet(r.NotAllIncludedSearchEntryCriteria)
	state.NoneIncludedSearchEntryCriteria = internaltypes.GetStringSet(r.NoneIncludedSearchEntryCriteria)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSearchEntryCriteriaUnknownValues(ctx, state)
}

// Read a ThirdPartySearchEntryCriteriaResponse object into the model struct
func readThirdPartySearchEntryCriteriaResponse(ctx context.Context, r *client.ThirdPartySearchEntryCriteriaResponse, state *searchEntryCriteriaResourceModel, expectedValues *searchEntryCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSearchEntryCriteriaUnknownValues(ctx, state)
}

// Create any update operations necessary to make the state match the plan
func createSearchEntryCriteriaOperations(plan searchEntryCriteriaResourceModel, state searchEntryCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedSearchEntryCriteria, state.AllIncludedSearchEntryCriteria, "all-included-search-entry-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedSearchEntryCriteria, state.AnyIncludedSearchEntryCriteria, "any-included-search-entry-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedSearchEntryCriteria, state.NotAllIncludedSearchEntryCriteria, "not-all-included-search-entry-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedSearchEntryCriteria, state.NoneIncludedSearchEntryCriteria, "none-included-search-entry-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.RequestCriteria, state.RequestCriteria, "request-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedEntryControl, state.AllIncludedEntryControl, "all-included-entry-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedEntryControl, state.AnyIncludedEntryControl, "any-included-entry-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedEntryControl, state.NotAllIncludedEntryControl, "not-all-included-entry-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedEntryControl, state.NoneIncludedEntryControl, "none-included-entry-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedEntryBaseDN, state.IncludedEntryBaseDN, "included-entry-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedEntryBaseDN, state.ExcludedEntryBaseDN, "excluded-entry-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedEntryFilter, state.AllIncludedEntryFilter, "all-included-entry-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedEntryFilter, state.AnyIncludedEntryFilter, "any-included-entry-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedEntryFilter, state.NotAllIncludedEntryFilter, "not-all-included-entry-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedEntryFilter, state.NoneIncludedEntryFilter, "none-included-entry-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedEntryGroupDN, state.AllIncludedEntryGroupDN, "all-included-entry-group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedEntryGroupDN, state.AnyIncludedEntryGroupDN, "any-included-entry-group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedEntryGroupDN, state.NotAllIncludedEntryGroupDN, "not-all-included-entry-group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedEntryGroupDN, state.NoneIncludedEntryGroupDN, "none-included-entry-group-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a simple search-entry-criteria
func (r *searchEntryCriteriaResource) CreateSimpleSearchEntryCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan searchEntryCriteriaResourceModel) (*searchEntryCriteriaResourceModel, error) {
	addRequest := client.NewAddSimpleSearchEntryCriteriaRequest(plan.Name.ValueString(),
		[]client.EnumsimpleSearchEntryCriteriaSchemaUrn{client.ENUMSIMPLESEARCHENTRYCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0SEARCH_ENTRY_CRITERIASIMPLE})
	addOptionalSimpleSearchEntryCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.SearchEntryCriteriaApi.AddSearchEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddSearchEntryCriteriaRequest(
		client.AddSimpleSearchEntryCriteriaRequestAsAddSearchEntryCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.SearchEntryCriteriaApi.AddSearchEntryCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Search Entry Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state searchEntryCriteriaResourceModel
	readSimpleSearchEntryCriteriaResponse(ctx, addResponse.SimpleSearchEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a aggregate search-entry-criteria
func (r *searchEntryCriteriaResource) CreateAggregateSearchEntryCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan searchEntryCriteriaResourceModel) (*searchEntryCriteriaResourceModel, error) {
	addRequest := client.NewAddAggregateSearchEntryCriteriaRequest(plan.Name.ValueString(),
		[]client.EnumaggregateSearchEntryCriteriaSchemaUrn{client.ENUMAGGREGATESEARCHENTRYCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0SEARCH_ENTRY_CRITERIAAGGREGATE})
	addOptionalAggregateSearchEntryCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.SearchEntryCriteriaApi.AddSearchEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddSearchEntryCriteriaRequest(
		client.AddAggregateSearchEntryCriteriaRequestAsAddSearchEntryCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.SearchEntryCriteriaApi.AddSearchEntryCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Search Entry Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state searchEntryCriteriaResourceModel
	readAggregateSearchEntryCriteriaResponse(ctx, addResponse.AggregateSearchEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party search-entry-criteria
func (r *searchEntryCriteriaResource) CreateThirdPartySearchEntryCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan searchEntryCriteriaResourceModel) (*searchEntryCriteriaResourceModel, error) {
	addRequest := client.NewAddThirdPartySearchEntryCriteriaRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartySearchEntryCriteriaSchemaUrn{client.ENUMTHIRDPARTYSEARCHENTRYCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0SEARCH_ENTRY_CRITERIATHIRD_PARTY},
		plan.ExtensionClass.ValueString())
	addOptionalThirdPartySearchEntryCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.SearchEntryCriteriaApi.AddSearchEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddSearchEntryCriteriaRequest(
		client.AddThirdPartySearchEntryCriteriaRequestAsAddSearchEntryCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.SearchEntryCriteriaApi.AddSearchEntryCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Search Entry Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state searchEntryCriteriaResourceModel
	readThirdPartySearchEntryCriteriaResponse(ctx, addResponse.ThirdPartySearchEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *searchEntryCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan searchEntryCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *searchEntryCriteriaResourceModel
	var err error
	if plan.Type.ValueString() == "simple" {
		state, err = r.CreateSimpleSearchEntryCriteria(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "aggregate" {
		state, err = r.CreateAggregateSearchEntryCriteria(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartySearchEntryCriteria(ctx, req, resp, plan)
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
func (r *defaultSearchEntryCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan searchEntryCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SearchEntryCriteriaApi.GetSearchEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Search Entry Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state searchEntryCriteriaResourceModel
	if readResponse.SimpleSearchEntryCriteriaResponse != nil {
		readSimpleSearchEntryCriteriaResponse(ctx, readResponse.SimpleSearchEntryCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AggregateSearchEntryCriteriaResponse != nil {
		readAggregateSearchEntryCriteriaResponse(ctx, readResponse.AggregateSearchEntryCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartySearchEntryCriteriaResponse != nil {
		readThirdPartySearchEntryCriteriaResponse(ctx, readResponse.ThirdPartySearchEntryCriteriaResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.SearchEntryCriteriaApi.UpdateSearchEntryCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createSearchEntryCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.SearchEntryCriteriaApi.UpdateSearchEntryCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Search Entry Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.SimpleSearchEntryCriteriaResponse != nil {
			readSimpleSearchEntryCriteriaResponse(ctx, updateResponse.SimpleSearchEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AggregateSearchEntryCriteriaResponse != nil {
			readAggregateSearchEntryCriteriaResponse(ctx, updateResponse.AggregateSearchEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartySearchEntryCriteriaResponse != nil {
			readThirdPartySearchEntryCriteriaResponse(ctx, updateResponse.ThirdPartySearchEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *searchEntryCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSearchEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultSearchEntryCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSearchEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readSearchEntryCriteria(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state searchEntryCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.SearchEntryCriteriaApi.GetSearchEntryCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Search Entry Criteria", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Search Entry Criteria", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.SimpleSearchEntryCriteriaResponse != nil {
		readSimpleSearchEntryCriteriaResponse(ctx, readResponse.SimpleSearchEntryCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AggregateSearchEntryCriteriaResponse != nil {
		readAggregateSearchEntryCriteriaResponse(ctx, readResponse.AggregateSearchEntryCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartySearchEntryCriteriaResponse != nil {
		readThirdPartySearchEntryCriteriaResponse(ctx, readResponse.ThirdPartySearchEntryCriteriaResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *searchEntryCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSearchEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSearchEntryCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSearchEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSearchEntryCriteria(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan searchEntryCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state searchEntryCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.SearchEntryCriteriaApi.UpdateSearchEntryCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createSearchEntryCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.SearchEntryCriteriaApi.UpdateSearchEntryCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Search Entry Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.SimpleSearchEntryCriteriaResponse != nil {
			readSimpleSearchEntryCriteriaResponse(ctx, updateResponse.SimpleSearchEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AggregateSearchEntryCriteriaResponse != nil {
			readAggregateSearchEntryCriteriaResponse(ctx, updateResponse.AggregateSearchEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartySearchEntryCriteriaResponse != nil {
			readThirdPartySearchEntryCriteriaResponse(ctx, updateResponse.ThirdPartySearchEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSearchEntryCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *searchEntryCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state searchEntryCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.SearchEntryCriteriaApi.DeleteSearchEntryCriteriaExecute(r.apiClient.SearchEntryCriteriaApi.DeleteSearchEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Search Entry Criteria", err, httpResp)
		return
	}
}

func (r *searchEntryCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSearchEntryCriteria(ctx, req, resp)
}

func (r *defaultSearchEntryCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSearchEntryCriteria(ctx, req, resp)
}

func importSearchEntryCriteria(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
