package config

import (
	"context"
	"time"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &locationResource{}
	_ resource.ResourceWithConfigure   = &locationResource{}
	_ resource.ResourceWithImportState = &locationResource{}
)

// Create a Location resource
func NewLocationResource() resource.Resource {
	return &locationResource{}
}

// locationResource is the resource implementation.
type locationResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// locationResourceModel maps the resource schema data.
type locationResourceModel struct {
	// Id field required for acceptance testing framework
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
}

// Metadata returns the resource type name.
func (r *locationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_location"
}

// GetSchema defines the schema for the resource.
func (r *locationResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	schema := tfsdk.Schema{
		Description: "Manages a Location.",
		Attributes: map[string]tfsdk.Attribute{
			"name": {
				Description: "Name of the Location.",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"description": {
				Description: "A description for this Location.",
				Type:        types.StringType,
				Optional:    true,
			},
		},
	}
	AddCommonSchema(&schema)
	return schema, nil
}

// Configure adds the provider configured client to the resource.
func (r *locationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Add optional fields to create request
func addOptionalLocationFields(addRequest *client.AddLocationRequest, plan locationResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
}

// Create a new resource
func (r *locationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan locationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddLocationRequest(plan.Name.ValueString())
	addOptionalLocationFields(addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LocationApi.AddLocation(ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLocationRequest(*addRequest)

	locationResponse, httpResp, err := r.apiClient.LocationApi.AddLocationExecute(apiAddRequest)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Location", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := locationResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state locationResourceModel
	readLocationResponse(ctx, locationResponse, &state, &plan)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read a LocationResponse object into the model struct
func readLocationResponse(ctx context.Context, r *client.LocationResponse, state *locationResourceModel, expectedValues *locationResourceModel) {
	// Placeholder Id value for acceptance test framework
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	// If a plan was provided and is using an empty string, use that for a nil string in the response.
	// To PingDirectory, nil and empty string is equivalent, but to Terraform they are distinct. So we
	// just want to match whatever is in the plan here.
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20)
}

// Read resource information
func (r *locationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state locationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	locationResponse, httpResp, err := r.apiClient.LocationApi.GetLocation(ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Location", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := locationResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readLocationResponse(ctx, locationResponse, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create any update operations necessary to make the state match the plan
func createLocationOperations(plan locationResourceModel, state locationResourceModel) []client.Operation {
	var ops []client.Operation

	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Update a resource
func (r *locationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan locationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state locationResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.LocationApi.UpdateLocation(ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createLocationOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		locationResponse, httpResp, err := r.apiClient.LocationApi.UpdateLocationExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Location", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := locationResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLocationResponse(ctx, locationResponse, &state, &plan)
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
func (r *locationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state locationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LocationApi.DeleteLocationExecute(r.apiClient.LocationApi.DeleteLocation(ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Location", err, httpResp)
		return
	}
}

func (r *locationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to Name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
