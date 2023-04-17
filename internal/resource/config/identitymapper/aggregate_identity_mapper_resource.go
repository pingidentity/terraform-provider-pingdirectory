package identitymapper

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
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
	_ resource.Resource                = &aggregateIdentityMapperResource{}
	_ resource.ResourceWithConfigure   = &aggregateIdentityMapperResource{}
	_ resource.ResourceWithImportState = &aggregateIdentityMapperResource{}
	_ resource.Resource                = &defaultAggregateIdentityMapperResource{}
	_ resource.ResourceWithConfigure   = &defaultAggregateIdentityMapperResource{}
	_ resource.ResourceWithImportState = &defaultAggregateIdentityMapperResource{}
)

// Create a Aggregate Identity Mapper resource
func NewAggregateIdentityMapperResource() resource.Resource {
	return &aggregateIdentityMapperResource{}
}

func NewDefaultAggregateIdentityMapperResource() resource.Resource {
	return &defaultAggregateIdentityMapperResource{}
}

// aggregateIdentityMapperResource is the resource implementation.
type aggregateIdentityMapperResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultAggregateIdentityMapperResource is the resource implementation.
type defaultAggregateIdentityMapperResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *aggregateIdentityMapperResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aggregate_identity_mapper"
}

func (r *defaultAggregateIdentityMapperResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_aggregate_identity_mapper"
}

// Configure adds the provider configured client to the resource.
func (r *aggregateIdentityMapperResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultAggregateIdentityMapperResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type aggregateIdentityMapperResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	LastUpdated               types.String `tfsdk:"last_updated"`
	Notifications             types.Set    `tfsdk:"notifications"`
	RequiredActions           types.Set    `tfsdk:"required_actions"`
	AllIncludedIdentityMapper types.Set    `tfsdk:"all_included_identity_mapper"`
	AnyIncludedIdentityMapper types.Set    `tfsdk:"any_included_identity_mapper"`
	Description               types.String `tfsdk:"description"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *aggregateIdentityMapperResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	aggregateIdentityMapperSchema(ctx, req, resp, false)
}

func (r *defaultAggregateIdentityMapperResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	aggregateIdentityMapperSchema(ctx, req, resp, true)
}

func aggregateIdentityMapperSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Aggregate Identity Mapper.",
		Attributes: map[string]schema.Attribute{
			"all_included_identity_mapper": schema.SetAttribute{
				Description: "The set of identity mappers that must all match the target entry. Each identity mapper must uniquely match the same target entry. If any of the identity mappers match multiple entries, if any of them match zero entries, or if any of them match different entries, then the mapping will fail.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"any_included_identity_mapper": schema.SetAttribute{
				Description: "The set of identity mappers that will be used to identify the target entry. At least one identity mapper must uniquely match an entry. If multiple identity mappers match entries, then they must all uniquely match the same entry. If none of the identity mappers match any entries, if any of them match multiple entries, or if any of them match different entries, then the mapping will fail.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Identity Mapper",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Identity Mapper is enabled for use.",
				Required:    true,
			},
		},
	}
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add config validators
func (r aggregateIdentityMapperResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("any_included_identity_mapper"),
			path.MatchRoot("all_included_identity_mapper"),
		),
	}
}

// Add optional fields to create request
func addOptionalAggregateIdentityMapperFields(ctx context.Context, addRequest *client.AddAggregateIdentityMapperRequest, plan aggregateIdentityMapperResourceModel) {
	if internaltypes.IsDefined(plan.AllIncludedIdentityMapper) {
		var slice []string
		plan.AllIncludedIdentityMapper.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedIdentityMapper = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedIdentityMapper) {
		var slice []string
		plan.AnyIncludedIdentityMapper.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedIdentityMapper = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a AggregateIdentityMapperResponse object into the model struct
func readAggregateIdentityMapperResponse(ctx context.Context, r *client.AggregateIdentityMapperResponse, state *aggregateIdentityMapperResourceModel, expectedValues *aggregateIdentityMapperResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.AllIncludedIdentityMapper = internaltypes.GetStringSet(r.AllIncludedIdentityMapper)
	state.AnyIncludedIdentityMapper = internaltypes.GetStringSet(r.AnyIncludedIdentityMapper)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createAggregateIdentityMapperOperations(plan aggregateIdentityMapperResourceModel, state aggregateIdentityMapperResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedIdentityMapper, state.AllIncludedIdentityMapper, "all-included-identity-mapper")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedIdentityMapper, state.AnyIncludedIdentityMapper, "any-included-identity-mapper")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *aggregateIdentityMapperResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan aggregateIdentityMapperResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddAggregateIdentityMapperRequest(plan.Id.ValueString(),
		[]client.EnumaggregateIdentityMapperSchemaUrn{client.ENUMAGGREGATEIDENTITYMAPPERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0IDENTITY_MAPPERAGGREGATE},
		plan.Enabled.ValueBool())
	addOptionalAggregateIdentityMapperFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.IdentityMapperApi.AddIdentityMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddIdentityMapperRequest(
		client.AddAggregateIdentityMapperRequestAsAddIdentityMapperRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.IdentityMapperApi.AddIdentityMapperExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Aggregate Identity Mapper", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state aggregateIdentityMapperResourceModel
	readAggregateIdentityMapperResponse(ctx, addResponse.AggregateIdentityMapperResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultAggregateIdentityMapperResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan aggregateIdentityMapperResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.IdentityMapperApi.GetIdentityMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Aggregate Identity Mapper", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state aggregateIdentityMapperResourceModel
	readAggregateIdentityMapperResponse(ctx, readResponse.AggregateIdentityMapperResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.IdentityMapperApi.UpdateIdentityMapper(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createAggregateIdentityMapperOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.IdentityMapperApi.UpdateIdentityMapperExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Aggregate Identity Mapper", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAggregateIdentityMapperResponse(ctx, updateResponse.AggregateIdentityMapperResponse, &state, &plan, &resp.Diagnostics)
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
func (r *aggregateIdentityMapperResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAggregateIdentityMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAggregateIdentityMapperResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAggregateIdentityMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readAggregateIdentityMapper(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state aggregateIdentityMapperResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.IdentityMapperApi.GetIdentityMapper(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Aggregate Identity Mapper", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAggregateIdentityMapperResponse(ctx, readResponse.AggregateIdentityMapperResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *aggregateIdentityMapperResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAggregateIdentityMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAggregateIdentityMapperResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAggregateIdentityMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateAggregateIdentityMapper(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan aggregateIdentityMapperResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state aggregateIdentityMapperResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.IdentityMapperApi.UpdateIdentityMapper(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createAggregateIdentityMapperOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.IdentityMapperApi.UpdateIdentityMapperExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Aggregate Identity Mapper", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAggregateIdentityMapperResponse(ctx, updateResponse.AggregateIdentityMapperResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultAggregateIdentityMapperResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *aggregateIdentityMapperResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state aggregateIdentityMapperResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.IdentityMapperApi.DeleteIdentityMapperExecute(r.apiClient.IdentityMapperApi.DeleteIdentityMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Aggregate Identity Mapper", err, httpResp)
		return
	}
}

func (r *aggregateIdentityMapperResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAggregateIdentityMapper(ctx, req, resp)
}

func (r *defaultAggregateIdentityMapperResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAggregateIdentityMapper(ctx, req, resp)
}

func importAggregateIdentityMapper(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
