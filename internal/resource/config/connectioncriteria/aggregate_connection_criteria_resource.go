package connectioncriteria

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
	_ resource.Resource                = &aggregateConnectionCriteriaResource{}
	_ resource.ResourceWithConfigure   = &aggregateConnectionCriteriaResource{}
	_ resource.ResourceWithImportState = &aggregateConnectionCriteriaResource{}
	_ resource.Resource                = &defaultAggregateConnectionCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultAggregateConnectionCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultAggregateConnectionCriteriaResource{}
)

// Create a Aggregate Connection Criteria resource
func NewAggregateConnectionCriteriaResource() resource.Resource {
	return &aggregateConnectionCriteriaResource{}
}

func NewDefaultAggregateConnectionCriteriaResource() resource.Resource {
	return &defaultAggregateConnectionCriteriaResource{}
}

// aggregateConnectionCriteriaResource is the resource implementation.
type aggregateConnectionCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultAggregateConnectionCriteriaResource is the resource implementation.
type defaultAggregateConnectionCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *aggregateConnectionCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aggregate_connection_criteria"
}

func (r *defaultAggregateConnectionCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_aggregate_connection_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *aggregateConnectionCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultAggregateConnectionCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type aggregateConnectionCriteriaResourceModel struct {
	Id                               types.String `tfsdk:"id"`
	LastUpdated                      types.String `tfsdk:"last_updated"`
	Notifications                    types.Set    `tfsdk:"notifications"`
	RequiredActions                  types.Set    `tfsdk:"required_actions"`
	AllIncludedConnectionCriteria    types.Set    `tfsdk:"all_included_connection_criteria"`
	AnyIncludedConnectionCriteria    types.Set    `tfsdk:"any_included_connection_criteria"`
	NotAllIncludedConnectionCriteria types.Set    `tfsdk:"not_all_included_connection_criteria"`
	NoneIncludedConnectionCriteria   types.Set    `tfsdk:"none_included_connection_criteria"`
	Description                      types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *aggregateConnectionCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	aggregateConnectionCriteriaSchema(ctx, req, resp, false)
}

func (r *defaultAggregateConnectionCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	aggregateConnectionCriteriaSchema(ctx, req, resp, true)
}

func aggregateConnectionCriteriaSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Aggregate Connection Criteria.",
		Attributes: map[string]schema.Attribute{
			"all_included_connection_criteria": schema.SetAttribute{
				Description: "Specifies a connection criteria object that must match the associated client connection in order to match the aggregate connection criteria. If one or more all-included connection criteria objects are provided, then a client connection must match all of them in order to match the aggregate connection criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"any_included_connection_criteria": schema.SetAttribute{
				Description: "Specifies a connection criteria object that may match the associated client connection in order to match the aggregate connection criteria. If one or more any-included connection criteria objects are provided, then a client connection must match at least one of them in order to match the aggregate connection criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"not_all_included_connection_criteria": schema.SetAttribute{
				Description: "Specifies a connection criteria object that should not match the associated client connection in order to match the aggregate connection criteria. If one or more not-all-included connection criteria objects are provided, then a client connection must not match all of them (that is, it may match zero or more of them, but it must not match all of them) in order to match the aggregate connection criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"none_included_connection_criteria": schema.SetAttribute{
				Description: "Specifies a connection criteria object that must not match the associated client connection in order to match the aggregate connection criteria. If one or more none-included connection criteria objects are provided, then a client connection must not match any of them in order to match the aggregate connection criteria.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Connection Criteria",
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
func addOptionalAggregateConnectionCriteriaFields(ctx context.Context, addRequest *client.AddAggregateConnectionCriteriaRequest, plan aggregateConnectionCriteriaResourceModel) {
	if internaltypes.IsDefined(plan.AllIncludedConnectionCriteria) {
		var slice []string
		plan.AllIncludedConnectionCriteria.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedConnectionCriteria = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedConnectionCriteria) {
		var slice []string
		plan.AnyIncludedConnectionCriteria.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedConnectionCriteria = slice
	}
	if internaltypes.IsDefined(plan.NotAllIncludedConnectionCriteria) {
		var slice []string
		plan.NotAllIncludedConnectionCriteria.ElementsAs(ctx, &slice, false)
		addRequest.NotAllIncludedConnectionCriteria = slice
	}
	if internaltypes.IsDefined(plan.NoneIncludedConnectionCriteria) {
		var slice []string
		plan.NoneIncludedConnectionCriteria.ElementsAs(ctx, &slice, false)
		addRequest.NoneIncludedConnectionCriteria = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
}

// Read a AggregateConnectionCriteriaResponse object into the model struct
func readAggregateConnectionCriteriaResponse(ctx context.Context, r *client.AggregateConnectionCriteriaResponse, state *aggregateConnectionCriteriaResourceModel, expectedValues *aggregateConnectionCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.AllIncludedConnectionCriteria = internaltypes.GetStringSet(r.AllIncludedConnectionCriteria)
	state.AnyIncludedConnectionCriteria = internaltypes.GetStringSet(r.AnyIncludedConnectionCriteria)
	state.NotAllIncludedConnectionCriteria = internaltypes.GetStringSet(r.NotAllIncludedConnectionCriteria)
	state.NoneIncludedConnectionCriteria = internaltypes.GetStringSet(r.NoneIncludedConnectionCriteria)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createAggregateConnectionCriteriaOperations(plan aggregateConnectionCriteriaResourceModel, state aggregateConnectionCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedConnectionCriteria, state.AllIncludedConnectionCriteria, "all-included-connection-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedConnectionCriteria, state.AnyIncludedConnectionCriteria, "any-included-connection-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NotAllIncludedConnectionCriteria, state.NotAllIncludedConnectionCriteria, "not-all-included-connection-criteria")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NoneIncludedConnectionCriteria, state.NoneIncludedConnectionCriteria, "none-included-connection-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *aggregateConnectionCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan aggregateConnectionCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddAggregateConnectionCriteriaRequest(plan.Id.ValueString(),
		[]client.EnumaggregateConnectionCriteriaSchemaUrn{client.ENUMAGGREGATECONNECTIONCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CONNECTION_CRITERIAAGGREGATE})
	addOptionalAggregateConnectionCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ConnectionCriteriaApi.AddConnectionCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddConnectionCriteriaRequest(
		client.AddAggregateConnectionCriteriaRequestAsAddConnectionCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ConnectionCriteriaApi.AddConnectionCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Aggregate Connection Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state aggregateConnectionCriteriaResourceModel
	readAggregateConnectionCriteriaResponse(ctx, addResponse.AggregateConnectionCriteriaResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultAggregateConnectionCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan aggregateConnectionCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ConnectionCriteriaApi.GetConnectionCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Aggregate Connection Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state aggregateConnectionCriteriaResourceModel
	readAggregateConnectionCriteriaResponse(ctx, readResponse.AggregateConnectionCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ConnectionCriteriaApi.UpdateConnectionCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createAggregateConnectionCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ConnectionCriteriaApi.UpdateConnectionCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Aggregate Connection Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAggregateConnectionCriteriaResponse(ctx, updateResponse.AggregateConnectionCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *aggregateConnectionCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAggregateConnectionCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAggregateConnectionCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAggregateConnectionCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readAggregateConnectionCriteria(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state aggregateConnectionCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ConnectionCriteriaApi.GetConnectionCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Aggregate Connection Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAggregateConnectionCriteriaResponse(ctx, readResponse.AggregateConnectionCriteriaResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *aggregateConnectionCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAggregateConnectionCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAggregateConnectionCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAggregateConnectionCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateAggregateConnectionCriteria(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan aggregateConnectionCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state aggregateConnectionCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ConnectionCriteriaApi.UpdateConnectionCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createAggregateConnectionCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ConnectionCriteriaApi.UpdateConnectionCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Aggregate Connection Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAggregateConnectionCriteriaResponse(ctx, updateResponse.AggregateConnectionCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultAggregateConnectionCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *aggregateConnectionCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state aggregateConnectionCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ConnectionCriteriaApi.DeleteConnectionCriteriaExecute(r.apiClient.ConnectionCriteriaApi.DeleteConnectionCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Aggregate Connection Criteria", err, httpResp)
		return
	}
}

func (r *aggregateConnectionCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAggregateConnectionCriteria(ctx, req, resp)
}

func (r *defaultAggregateConnectionCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAggregateConnectionCriteria(ctx, req, resp)
}

func importAggregateConnectionCriteria(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
