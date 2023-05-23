package searchentrycriteria

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &simpleSearchEntryCriteriaResource{}
	_ resource.ResourceWithConfigure   = &simpleSearchEntryCriteriaResource{}
	_ resource.ResourceWithImportState = &simpleSearchEntryCriteriaResource{}
	_ resource.Resource                = &defaultSimpleSearchEntryCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultSimpleSearchEntryCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultSimpleSearchEntryCriteriaResource{}
)

// Create a Simple Search Entry Criteria resource
func NewSimpleSearchEntryCriteriaResource() resource.Resource {
	return &simpleSearchEntryCriteriaResource{}
}

func NewDefaultSimpleSearchEntryCriteriaResource() resource.Resource {
	return &defaultSimpleSearchEntryCriteriaResource{}
}

// simpleSearchEntryCriteriaResource is the resource implementation.
type simpleSearchEntryCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSimpleSearchEntryCriteriaResource is the resource implementation.
type defaultSimpleSearchEntryCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *simpleSearchEntryCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_simple_search_entry_criteria"
}

func (r *defaultSimpleSearchEntryCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_simple_search_entry_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *simpleSearchEntryCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultSimpleSearchEntryCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type simpleSearchEntryCriteriaResourceModel struct {
	Id                         types.String `tfsdk:"id"`
	LastUpdated                types.String `tfsdk:"last_updated"`
	Notifications              types.Set    `tfsdk:"notifications"`
	RequiredActions            types.Set    `tfsdk:"required_actions"`
	RequestCriteria            types.String `tfsdk:"request_criteria"`
	AllIncludedEntryControl    types.Set    `tfsdk:"all_included_entry_control"`
	AnyIncludedEntryControl    types.Set    `tfsdk:"any_included_entry_control"`
	NotAllIncludedEntryControl types.Set    `tfsdk:"not_all_included_entry_control"`
	NoneIncludedEntryControl   types.Set    `tfsdk:"none_included_entry_control"`
	IncludedEntryBaseDN        types.Set    `tfsdk:"included_entry_base_dn"`
	ExcludedEntryBaseDN        types.Set    `tfsdk:"excluded_entry_base_dn"`
	AllIncludedEntryFilter     types.Set    `tfsdk:"all_included_entry_filter"`
	AnyIncludedEntryFilter     types.Set    `tfsdk:"any_included_entry_filter"`
	NotAllIncludedEntryFilter  types.Set    `tfsdk:"not_all_included_entry_filter"`
	NoneIncludedEntryFilter    types.Set    `tfsdk:"none_included_entry_filter"`
	AllIncludedEntryGroupDN    types.Set    `tfsdk:"all_included_entry_group_dn"`
	AnyIncludedEntryGroupDN    types.Set    `tfsdk:"any_included_entry_group_dn"`
	NotAllIncludedEntryGroupDN types.Set    `tfsdk:"not_all_included_entry_group_dn"`
	NoneIncludedEntryGroupDN   types.Set    `tfsdk:"none_included_entry_group_dn"`
	Description                types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *simpleSearchEntryCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	simpleSearchEntryCriteriaSchema(ctx, req, resp, false)
}

func (r *defaultSimpleSearchEntryCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	simpleSearchEntryCriteriaSchema(ctx, req, resp, true)
}

func simpleSearchEntryCriteriaSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Simple Search Entry Criteria.",
		Attributes: map[string]schema.Attribute{
			"request_criteria": schema.StringAttribute{
				Description: "Specifies a request criteria object that must match the associated request for entries included in this Simple Search Entry Criteria. of them.",
				Optional:    true,
			},
			"all_included_entry_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must be present in search result entries included in this Simple Search Entry Criteria. If any control OIDs are provided, then the entry must contain all of those controls.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"any_included_entry_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that may be present in search result entries included in this Simple Search Entry Criteria. If any control OIDs are provided, then the entry must contain at least one of those controls.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"not_all_included_entry_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that should not be present in search result entries included in this Simple Search Entry Criteria. If any control OIDs are provided, then the entry must not contain at least one of those controls (that is, it may contain zero or more of those controls, but not all of them).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"none_included_entry_control": schema.SetAttribute{
				Description: "Specifies the OID of a control that must not be present in search result entries included in this Simple Search Entry Criteria. If any control OIDs are provided, then the entry must not contain any of those controls.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"included_entry_base_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which entries included in this Simple Search Entry Criteria may exist.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"excluded_entry_base_dn": schema.SetAttribute{
				Description: "Specifies a base DN below which entries included in this Simple Search Entry Criteria may not exist.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"all_included_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that must match search result entries included in this Simple Search Entry Criteria. Note that this matching will be performed against the entry that is actually returned to the client and may not reflect the complete entry stored in the server. If any filters are provided, then the returned entry must match all of those filters.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"any_included_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that may match search result entries included in this Simple Search Entry Criteria. Note that this matching will be performed against the entry that is actually returned to the client and may not reflect the complete entry stored in the server. If any filters are provided, then the entry must match at least one of those filters.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"not_all_included_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that should not match search result entries included in this Simple Search Entry Criteria. Note that this matching will be performed against the entry that is actually returned to the client and may not reflect the complete entry stored in the server. If any filters are provided, then the entry must not match at least one of those filters (that is, the entry may match zero or more of those filters, but not of all of them).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"none_included_entry_filter": schema.SetAttribute{
				Description: "Specifies a search filter that must not match search result entries included in this Simple Search Entry Criteria. Note that this matching will be performed against the entry that is actually returned to the client and may not reflect the complete entry stored in the server. If any filters are provided, then the entry must not match any of those filters.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"all_included_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the entry must be a member to be included in this Simple Search Entry Criteria. If any group DNs are provided, then the entry must be a member of all of them.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"any_included_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the entry may be a member to be included in this Simple Search Entry Criteria. If any group DNs are provided, then the entry must be a member of at least one of them.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"not_all_included_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the entry should not be a member to be included in this Simple Search Entry Criteria. If any group DNs are provided, then the entry must not be a member of at least one of them (that is, the entry may be a member of zero or more of the specified groups, but not of all of them).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"none_included_entry_group_dn": schema.SetAttribute{
				Description: "Specifies the DN of a group in which the user associated with the entry must not be a member to be included in this Simple Search Entry Criteria. If any group DNs are provided, then the entry must not be a member of any of them.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Search Entry Criteria",
				Optional:    true,
			},
		},
	}
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalSimpleSearchEntryCriteriaFields(ctx context.Context, addRequest *client.AddSimpleSearchEntryCriteriaRequest, plan simpleSearchEntryCriteriaResourceModel) {
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

// Read a SimpleSearchEntryCriteriaResponse object into the model struct
func readSimpleSearchEntryCriteriaResponse(ctx context.Context, r *client.SimpleSearchEntryCriteriaResponse, state *simpleSearchEntryCriteriaResourceModel, expectedValues *simpleSearchEntryCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
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
}

// Create any update operations necessary to make the state match the plan
func createSimpleSearchEntryCriteriaOperations(plan simpleSearchEntryCriteriaResourceModel, state simpleSearchEntryCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
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

// Create a new resource
func (r *simpleSearchEntryCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan simpleSearchEntryCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddSimpleSearchEntryCriteriaRequest(plan.Id.ValueString(),
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Simple Search Entry Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state simpleSearchEntryCriteriaResourceModel
	readSimpleSearchEntryCriteriaResponse(ctx, addResponse.SimpleSearchEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultSimpleSearchEntryCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan simpleSearchEntryCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SearchEntryCriteriaApi.GetSearchEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Simple Search Entry Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state simpleSearchEntryCriteriaResourceModel
	readSimpleSearchEntryCriteriaResponse(ctx, readResponse.SimpleSearchEntryCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.SearchEntryCriteriaApi.UpdateSearchEntryCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createSimpleSearchEntryCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.SearchEntryCriteriaApi.UpdateSearchEntryCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Simple Search Entry Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSimpleSearchEntryCriteriaResponse(ctx, updateResponse.SimpleSearchEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *simpleSearchEntryCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSimpleSearchEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSimpleSearchEntryCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSimpleSearchEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSimpleSearchEntryCriteria(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state simpleSearchEntryCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.SearchEntryCriteriaApi.GetSearchEntryCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Simple Search Entry Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSimpleSearchEntryCriteriaResponse(ctx, readResponse.SimpleSearchEntryCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *simpleSearchEntryCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSimpleSearchEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSimpleSearchEntryCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSimpleSearchEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSimpleSearchEntryCriteria(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan simpleSearchEntryCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state simpleSearchEntryCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.SearchEntryCriteriaApi.UpdateSearchEntryCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createSimpleSearchEntryCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.SearchEntryCriteriaApi.UpdateSearchEntryCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Simple Search Entry Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSimpleSearchEntryCriteriaResponse(ctx, updateResponse.SimpleSearchEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSimpleSearchEntryCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *simpleSearchEntryCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state simpleSearchEntryCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.SearchEntryCriteriaApi.DeleteSearchEntryCriteriaExecute(r.apiClient.SearchEntryCriteriaApi.DeleteSearchEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Simple Search Entry Criteria", err, httpResp)
		return
	}
}

func (r *simpleSearchEntryCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSimpleSearchEntryCriteria(ctx, req, resp)
}

func (r *defaultSimpleSearchEntryCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSimpleSearchEntryCriteria(ctx, req, resp)
}

func importSimpleSearchEntryCriteria(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
