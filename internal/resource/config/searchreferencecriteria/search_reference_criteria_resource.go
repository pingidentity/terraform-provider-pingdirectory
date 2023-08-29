package searchreferencecriteria

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	_ resource.Resource                = &searchReferenceCriteriaResource{}
	_ resource.ResourceWithConfigure   = &searchReferenceCriteriaResource{}
	_ resource.ResourceWithImportState = &searchReferenceCriteriaResource{}
	_ resource.Resource                = &defaultSearchReferenceCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultSearchReferenceCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultSearchReferenceCriteriaResource{}
)

// Create a Search Reference Criteria resource
func NewSearchReferenceCriteriaResource() resource.Resource {
	return &searchReferenceCriteriaResource{}
}

func NewDefaultSearchReferenceCriteriaResource() resource.Resource {
	return &defaultSearchReferenceCriteriaResource{}
}

// searchReferenceCriteriaResource is the resource implementation.
type searchReferenceCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSearchReferenceCriteriaResource is the resource implementation.
type defaultSearchReferenceCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *searchReferenceCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_search_reference_criteria"
}

func (r *defaultSearchReferenceCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_search_reference_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *searchReferenceCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultSearchReferenceCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type searchReferenceCriteriaResourceModel struct {
	Id                                    types.String `tfsdk:"id"`
	Name                                  types.String `tfsdk:"name"`
	LastUpdated                           types.String `tfsdk:"last_updated"`
	Notifications                         types.Set    `tfsdk:"notifications"`
	RequiredActions                       types.Set    `tfsdk:"required_actions"`
	Type                                  types.String `tfsdk:"type"`
	ExtensionClass                        types.String `tfsdk:"extension_class"`
	ExtensionArgument                     types.Set    `tfsdk:"extension_argument"`
	AllIncludedSearchReferenceCriteria    types.Set    `tfsdk:"all_included_search_reference_criteria"`
	AnyIncludedSearchReferenceCriteria    types.Set    `tfsdk:"any_included_search_reference_criteria"`
	NotAllIncludedSearchReferenceCriteria types.Set    `tfsdk:"not_all_included_search_reference_criteria"`
	NoneIncludedSearchReferenceCriteria   types.Set    `tfsdk:"none_included_search_reference_criteria"`
	RequestCriteria                       types.String `tfsdk:"request_criteria"`
	AllIncludedReferenceControl           types.Set    `tfsdk:"all_included_reference_control"`
	AnyIncludedReferenceControl           types.Set    `tfsdk:"any_included_reference_control"`
	NotAllIncludedReferenceControl        types.Set    `tfsdk:"not_all_included_reference_control"`
	NoneIncludedReferenceControl          types.Set    `tfsdk:"none_included_reference_control"`
	Description                           types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *searchReferenceCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	searchReferenceCriteriaSchema(ctx, req, resp, false)
}

func (r *defaultSearchReferenceCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	searchReferenceCriteriaSchema(ctx, req, resp, true)
}

func searchReferenceCriteriaSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Search Reference Criteria.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Search Reference Criteria resource. Options are ['simple', 'aggregate', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"simple", "aggregate", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Search Reference Criteria.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Search Reference Criteria. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"all_included_search_reference_criteria": schema.SetAttribute{
				Description: "Specifies a search reference criteria object that must match the associated search result reference in order to match the aggregate search reference criteria. If one or more all-included search reference criteria objects are provided, then a search result reference must match all of them in order to match the aggregate search reference criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"any_included_search_reference_criteria": schema.SetAttribute{
				Description: "Specifies a search reference criteria object that may match the associated search result reference in order to match the aggregate search reference criteria. If one or more any-included search reference criteria objects are provided, then a search result reference must match at least one of them in order to match the aggregate search reference criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"not_all_included_search_reference_criteria": schema.SetAttribute{
				Description: "Specifies a search reference criteria object that should not match the associated search result reference in order to match the aggregate search reference criteria. If one or more not-all-included search reference criteria objects are provided, then a search result reference must not match all of them (that is, it may match zero or more of them, but it must not match all of them) in order to match the aggregate search reference criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"none_included_search_reference_criteria": schema.SetAttribute{
				Description: "Specifies a search reference criteria object that must not match the associated search result reference in order to match the aggregate search reference criteria. If one or more none-included search reference criteria objects are provided, then a search result reference must not match any of them in order to match the aggregate search reference criteria.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"request_criteria": schema.StringAttribute{
				Description: "Specifies a request criteria object that must match the associated request for references included in this Simple Search Reference Criteria.",
				Optional:    true,
			},
			"all_included_reference_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must be present in search result references included in this Simple Search Reference Criteria. If any control OIDs are provided, then the reference must contain all of those controls.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"any_included_reference_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that may be present in search result references included in this Simple Search Reference Criteria. If any control OIDs are provided, then the reference must contain at least one of those controls.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"not_all_included_reference_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that should not be present in search result references included in this Simple Search Reference Criteria. If any control OIDs are provided, then the reference must not contain at least one of those controls (that is, it may contain zero or more of those controls, but not all of them).",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"none_included_reference_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must not be present in search result references included in this Simple Search Reference Criteria. If any control OIDs are provided, then the reference must not contain any of those controls.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Search Reference Criteria",
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
func configValidatorsSearchReferenceCriteria() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("request_criteria"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("all_included_reference_control"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("any_included_reference_control"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("not_all_included_reference_control"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("none_included_reference_control"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("all_included_search_reference_criteria"),
			path.MatchRoot("type"),
			[]string{"aggregate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("any_included_search_reference_criteria"),
			path.MatchRoot("type"),
			[]string{"aggregate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("not_all_included_search_reference_criteria"),
			path.MatchRoot("type"),
			[]string{"aggregate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("none_included_search_reference_criteria"),
			path.MatchRoot("type"),
			[]string{"aggregate"},
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
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party",
			[]path.Expression{path.MatchRoot("extension_class")},
		),
	}
}

// Add config validators
func (r searchReferenceCriteriaResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsSearchReferenceCriteria()
}

// Add config validators
func (r defaultSearchReferenceCriteriaResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsSearchReferenceCriteria()
}

// Add optional fields to create request for simple search-reference-criteria
func addOptionalSimpleSearchReferenceCriteriaFields(ctx context.Context, addRequest *client.AddSimpleSearchReferenceCriteriaRequest, plan searchReferenceCriteriaResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		addRequest.RequestCriteria = plan.RequestCriteria.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.AllIncludedReferenceControl) {
		var slice []string
		plan.AllIncludedReferenceControl.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedReferenceControl = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedReferenceControl) {
		var slice []string
		plan.AnyIncludedReferenceControl.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedReferenceControl = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedReferenceControl) {
		var slice []string
		plan.NotAllIncludedReferenceControl.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedReferenceControl = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedReferenceControl) {
		var slice []string
		plan.NoneIncludedReferenceControl.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedReferenceControl = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for aggregate search-reference-criteria
func addOptionalAggregateSearchReferenceCriteriaFields(ctx context.Context, addRequest *client.AddAggregateSearchReferenceCriteriaRequest, plan searchReferenceCriteriaResourceModel) {
	if internaltypes.IsDefined(plan.AllIncludedSearchReferenceCriteria) {
		var slice []string
		plan.AllIncludedSearchReferenceCriteria.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedSearchReferenceCriteria = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedSearchReferenceCriteria) {
		var slice []string
		plan.AnyIncludedSearchReferenceCriteria.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedSearchReferenceCriteria = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedSearchReferenceCriteria) {
		var slice []string
		plan.NotAllIncludedSearchReferenceCriteria.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedSearchReferenceCriteria = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedSearchReferenceCriteria) {
		var slice []string
		plan.NoneIncludedSearchReferenceCriteria.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedSearchReferenceCriteria = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for third-party search-reference-criteria
func addOptionalThirdPartySearchReferenceCriteriaFields(ctx context.Context, addRequest *client.AddThirdPartySearchReferenceCriteriaRequest, plan searchReferenceCriteriaResourceModel) {
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
func populateSearchReferenceCriteriaUnknownValues(model *searchReferenceCriteriaResourceModel) {
	if model.AllIncludedReferenceControl.IsUnknown() || model.AllIncludedReferenceControl.IsNull() {
		model.AllIncludedReferenceControl, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AnyIncludedSearchReferenceCriteria.IsUnknown() || model.AnyIncludedSearchReferenceCriteria.IsNull() {
		model.AnyIncludedSearchReferenceCriteria, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.NotAllIncludedSearchReferenceCriteria.IsUnknown() || model.NotAllIncludedSearchReferenceCriteria.IsNull() {
		model.NotAllIncludedSearchReferenceCriteria, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.NoneIncludedReferenceControl.IsUnknown() || model.NoneIncludedReferenceControl.IsNull() {
		model.NoneIncludedReferenceControl, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.NotAllIncludedReferenceControl.IsUnknown() || model.NotAllIncludedReferenceControl.IsNull() {
		model.NotAllIncludedReferenceControl, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.NoneIncludedSearchReferenceCriteria.IsUnknown() || model.NoneIncludedSearchReferenceCriteria.IsNull() {
		model.NoneIncludedSearchReferenceCriteria, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AnyIncludedReferenceControl.IsUnknown() || model.AnyIncludedReferenceControl.IsNull() {
		model.AnyIncludedReferenceControl, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AllIncludedSearchReferenceCriteria.IsUnknown() || model.AllIncludedSearchReferenceCriteria.IsNull() {
		model.AllIncludedSearchReferenceCriteria, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *searchReferenceCriteriaResourceModel) populateAllComputedStringAttributes() {
	if model.RequestCriteria.IsUnknown() || model.RequestCriteria.IsNull() {
		model.RequestCriteria = types.StringValue("")
	}
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.ExtensionClass.IsUnknown() || model.ExtensionClass.IsNull() {
		model.ExtensionClass = types.StringValue("")
	}
}

// Read a SimpleSearchReferenceCriteriaResponse object into the model struct
func readSimpleSearchReferenceCriteriaResponse(ctx context.Context, r *client.SimpleSearchReferenceCriteriaResponse, state *searchReferenceCriteriaResourceModel, expectedValues *searchReferenceCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("simple")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.AllIncludedReferenceControl = internaltypes.GetStringSet(r.AllIncludedReferenceControl)
	state.AnyIncludedReferenceControl = internaltypes.GetStringSet(r.AnyIncludedReferenceControl)
	state.NotAllIncludedReferenceControl = internaltypes.GetStringSet(r.NotAllIncludedReferenceControl)
	state.NoneIncludedReferenceControl = internaltypes.GetStringSet(r.NoneIncludedReferenceControl)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSearchReferenceCriteriaUnknownValues(state)
}

// Read a AggregateSearchReferenceCriteriaResponse object into the model struct
func readAggregateSearchReferenceCriteriaResponse(ctx context.Context, r *client.AggregateSearchReferenceCriteriaResponse, state *searchReferenceCriteriaResourceModel, expectedValues *searchReferenceCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("aggregate")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AllIncludedSearchReferenceCriteria = internaltypes.GetStringSet(r.AllIncludedSearchReferenceCriteria)
	state.AnyIncludedSearchReferenceCriteria = internaltypes.GetStringSet(r.AnyIncludedSearchReferenceCriteria)
	state.NotAllIncludedSearchReferenceCriteria = internaltypes.GetStringSet(r.NotAllIncludedSearchReferenceCriteria)
	state.NoneIncludedSearchReferenceCriteria = internaltypes.GetStringSet(r.NoneIncludedSearchReferenceCriteria)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSearchReferenceCriteriaUnknownValues(state)
}

// Read a ThirdPartySearchReferenceCriteriaResponse object into the model struct
func readThirdPartySearchReferenceCriteriaResponse(ctx context.Context, r *client.ThirdPartySearchReferenceCriteriaResponse, state *searchReferenceCriteriaResourceModel, expectedValues *searchReferenceCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateSearchReferenceCriteriaUnknownValues(state)
}

// Create any update operations necessary to make the state match the plan
func createSearchReferenceCriteriaOperations(plan searchReferenceCriteriaResourceModel, state searchReferenceCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedSearchReferenceCriteria, state.AllIncludedSearchReferenceCriteria, "all-included-search-reference-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedSearchReferenceCriteria, state.AnyIncludedSearchReferenceCriteria, "any-included-search-reference-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedSearchReferenceCriteria, state.NotAllIncludedSearchReferenceCriteria, "not-all-included-search-reference-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedSearchReferenceCriteria, state.NoneIncludedSearchReferenceCriteria, "none-included-search-reference-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.RequestCriteria, state.RequestCriteria, "request-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedReferenceControl, state.AllIncludedReferenceControl, "all-included-reference-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedReferenceControl, state.AnyIncludedReferenceControl, "any-included-reference-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedReferenceControl, state.NotAllIncludedReferenceControl, "not-all-included-reference-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedReferenceControl, state.NoneIncludedReferenceControl, "none-included-reference-control")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a simple search-reference-criteria
func (r *searchReferenceCriteriaResource) CreateSimpleSearchReferenceCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan searchReferenceCriteriaResourceModel) (*searchReferenceCriteriaResourceModel, error) {
	addRequest := client.NewAddSimpleSearchReferenceCriteriaRequest(plan.Name.ValueString(),
		[]client.EnumsimpleSearchReferenceCriteriaSchemaUrn{client.ENUMSIMPLESEARCHREFERENCECRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0SEARCH_REFERENCE_CRITERIASIMPLE})
	addOptionalSimpleSearchReferenceCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.SearchReferenceCriteriaApi.AddSearchReferenceCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddSearchReferenceCriteriaRequest(
		client.AddSimpleSearchReferenceCriteriaRequestAsAddSearchReferenceCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.SearchReferenceCriteriaApi.AddSearchReferenceCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Search Reference Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state searchReferenceCriteriaResourceModel
	readSimpleSearchReferenceCriteriaResponse(ctx, addResponse.SimpleSearchReferenceCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a aggregate search-reference-criteria
func (r *searchReferenceCriteriaResource) CreateAggregateSearchReferenceCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan searchReferenceCriteriaResourceModel) (*searchReferenceCriteriaResourceModel, error) {
	addRequest := client.NewAddAggregateSearchReferenceCriteriaRequest(plan.Name.ValueString(),
		[]client.EnumaggregateSearchReferenceCriteriaSchemaUrn{client.ENUMAGGREGATESEARCHREFERENCECRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0SEARCH_REFERENCE_CRITERIAAGGREGATE})
	addOptionalAggregateSearchReferenceCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.SearchReferenceCriteriaApi.AddSearchReferenceCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddSearchReferenceCriteriaRequest(
		client.AddAggregateSearchReferenceCriteriaRequestAsAddSearchReferenceCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.SearchReferenceCriteriaApi.AddSearchReferenceCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Search Reference Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state searchReferenceCriteriaResourceModel
	readAggregateSearchReferenceCriteriaResponse(ctx, addResponse.AggregateSearchReferenceCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party search-reference-criteria
func (r *searchReferenceCriteriaResource) CreateThirdPartySearchReferenceCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan searchReferenceCriteriaResourceModel) (*searchReferenceCriteriaResourceModel, error) {
	addRequest := client.NewAddThirdPartySearchReferenceCriteriaRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartySearchReferenceCriteriaSchemaUrn{client.ENUMTHIRDPARTYSEARCHREFERENCECRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0SEARCH_REFERENCE_CRITERIATHIRD_PARTY},
		plan.ExtensionClass.ValueString())
	addOptionalThirdPartySearchReferenceCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.SearchReferenceCriteriaApi.AddSearchReferenceCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddSearchReferenceCriteriaRequest(
		client.AddThirdPartySearchReferenceCriteriaRequestAsAddSearchReferenceCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.SearchReferenceCriteriaApi.AddSearchReferenceCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Search Reference Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state searchReferenceCriteriaResourceModel
	readThirdPartySearchReferenceCriteriaResponse(ctx, addResponse.ThirdPartySearchReferenceCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *searchReferenceCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan searchReferenceCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *searchReferenceCriteriaResourceModel
	var err error
	if plan.Type.ValueString() == "simple" {
		state, err = r.CreateSimpleSearchReferenceCriteria(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "aggregate" {
		state, err = r.CreateAggregateSearchReferenceCriteria(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartySearchReferenceCriteria(ctx, req, resp, plan)
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
func (r *defaultSearchReferenceCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan searchReferenceCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SearchReferenceCriteriaApi.GetSearchReferenceCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Search Reference Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state searchReferenceCriteriaResourceModel
	if readResponse.SimpleSearchReferenceCriteriaResponse != nil {
		readSimpleSearchReferenceCriteriaResponse(ctx, readResponse.SimpleSearchReferenceCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AggregateSearchReferenceCriteriaResponse != nil {
		readAggregateSearchReferenceCriteriaResponse(ctx, readResponse.AggregateSearchReferenceCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartySearchReferenceCriteriaResponse != nil {
		readThirdPartySearchReferenceCriteriaResponse(ctx, readResponse.ThirdPartySearchReferenceCriteriaResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.SearchReferenceCriteriaApi.UpdateSearchReferenceCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createSearchReferenceCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.SearchReferenceCriteriaApi.UpdateSearchReferenceCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Search Reference Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.SimpleSearchReferenceCriteriaResponse != nil {
			readSimpleSearchReferenceCriteriaResponse(ctx, updateResponse.SimpleSearchReferenceCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AggregateSearchReferenceCriteriaResponse != nil {
			readAggregateSearchReferenceCriteriaResponse(ctx, updateResponse.AggregateSearchReferenceCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartySearchReferenceCriteriaResponse != nil {
			readThirdPartySearchReferenceCriteriaResponse(ctx, updateResponse.ThirdPartySearchReferenceCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	}

	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *searchReferenceCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSearchReferenceCriteria(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultSearchReferenceCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSearchReferenceCriteria(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readSearchReferenceCriteria(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state searchReferenceCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.SearchReferenceCriteriaApi.GetSearchReferenceCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Search Reference Criteria", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Search Reference Criteria", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.SimpleSearchReferenceCriteriaResponse != nil {
		readSimpleSearchReferenceCriteriaResponse(ctx, readResponse.SimpleSearchReferenceCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AggregateSearchReferenceCriteriaResponse != nil {
		readAggregateSearchReferenceCriteriaResponse(ctx, readResponse.AggregateSearchReferenceCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartySearchReferenceCriteriaResponse != nil {
		readThirdPartySearchReferenceCriteriaResponse(ctx, readResponse.ThirdPartySearchReferenceCriteriaResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *searchReferenceCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSearchReferenceCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSearchReferenceCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSearchReferenceCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSearchReferenceCriteria(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan searchReferenceCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state searchReferenceCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.SearchReferenceCriteriaApi.UpdateSearchReferenceCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createSearchReferenceCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.SearchReferenceCriteriaApi.UpdateSearchReferenceCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Search Reference Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.SimpleSearchReferenceCriteriaResponse != nil {
			readSimpleSearchReferenceCriteriaResponse(ctx, updateResponse.SimpleSearchReferenceCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AggregateSearchReferenceCriteriaResponse != nil {
			readAggregateSearchReferenceCriteriaResponse(ctx, updateResponse.AggregateSearchReferenceCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartySearchReferenceCriteriaResponse != nil {
			readThirdPartySearchReferenceCriteriaResponse(ctx, updateResponse.ThirdPartySearchReferenceCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSearchReferenceCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *searchReferenceCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state searchReferenceCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.SearchReferenceCriteriaApi.DeleteSearchReferenceCriteriaExecute(r.apiClient.SearchReferenceCriteriaApi.DeleteSearchReferenceCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Search Reference Criteria", err, httpResp)
		return
	}
}

func (r *searchReferenceCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSearchReferenceCriteria(ctx, req, resp)
}

func (r *defaultSearchReferenceCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSearchReferenceCriteria(ctx, req, resp)
}

func importSearchReferenceCriteria(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
