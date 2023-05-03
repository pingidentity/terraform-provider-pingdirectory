package logfieldmapping

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &errorLogFieldMappingResource{}
	_ resource.ResourceWithConfigure   = &errorLogFieldMappingResource{}
	_ resource.ResourceWithImportState = &errorLogFieldMappingResource{}
	_ resource.Resource                = &defaultErrorLogFieldMappingResource{}
	_ resource.ResourceWithConfigure   = &defaultErrorLogFieldMappingResource{}
	_ resource.ResourceWithImportState = &defaultErrorLogFieldMappingResource{}
)

// Create a Error Log Field Mapping resource
func NewErrorLogFieldMappingResource() resource.Resource {
	return &errorLogFieldMappingResource{}
}

func NewDefaultErrorLogFieldMappingResource() resource.Resource {
	return &defaultErrorLogFieldMappingResource{}
}

// errorLogFieldMappingResource is the resource implementation.
type errorLogFieldMappingResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultErrorLogFieldMappingResource is the resource implementation.
type defaultErrorLogFieldMappingResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *errorLogFieldMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_error_log_field_mapping"
}

func (r *defaultErrorLogFieldMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_error_log_field_mapping"
}

// Configure adds the provider configured client to the resource.
func (r *errorLogFieldMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultErrorLogFieldMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type errorLogFieldMappingResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	LastUpdated          types.String `tfsdk:"last_updated"`
	Notifications        types.Set    `tfsdk:"notifications"`
	RequiredActions      types.Set    `tfsdk:"required_actions"`
	LogFieldTimestamp    types.String `tfsdk:"log_field_timestamp"`
	LogFieldProductName  types.String `tfsdk:"log_field_product_name"`
	LogFieldInstanceName types.String `tfsdk:"log_field_instance_name"`
	LogFieldStartupid    types.String `tfsdk:"log_field_startupid"`
	LogFieldCategory     types.String `tfsdk:"log_field_category"`
	LogFieldSeverity     types.String `tfsdk:"log_field_severity"`
	LogFieldMessageID    types.String `tfsdk:"log_field_message_id"`
	LogFieldMessage      types.String `tfsdk:"log_field_message"`
	Description          types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *errorLogFieldMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	errorLogFieldMappingSchema(ctx, req, resp, false)
}

func (r *defaultErrorLogFieldMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	errorLogFieldMappingSchema(ctx, req, resp, true)
}

func errorLogFieldMappingSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Error Log Field Mapping.",
		Attributes: map[string]schema.Attribute{
			"log_field_timestamp": schema.StringAttribute{
				Description: "The time that the log message was generated.",
				Optional:    true,
			},
			"log_field_product_name": schema.StringAttribute{
				Description: "The name for this Directory Server product, which may be used to identify which product was used to log the message if multiple products log to the same database table.",
				Optional:    true,
			},
			"log_field_instance_name": schema.StringAttribute{
				Description: "A name that uniquely identifies this Directory Server instance, which may be used to identify which instance was used to log the message if multiple server instances log to the same database table.",
				Optional:    true,
			},
			"log_field_startupid": schema.StringAttribute{
				Description: "The startup ID for the Directory Server. A different value will be generated each time the server is started.",
				Optional:    true,
			},
			"log_field_category": schema.StringAttribute{
				Description: "The category for the log message.",
				Optional:    true,
			},
			"log_field_severity": schema.StringAttribute{
				Description: "The severity for the log message.",
				Optional:    true,
			},
			"log_field_message_id": schema.StringAttribute{
				Description: "The numeric value which uniquely identifies the type of message.",
				Optional:    true,
			},
			"log_field_message": schema.StringAttribute{
				Description: "The text of the log message.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log Field Mapping",
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
func addOptionalErrorLogFieldMappingFields(ctx context.Context, addRequest *client.AddErrorLogFieldMappingRequest, plan errorLogFieldMappingResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldTimestamp) {
		addRequest.LogFieldTimestamp = plan.LogFieldTimestamp.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldProductName) {
		addRequest.LogFieldProductName = plan.LogFieldProductName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldInstanceName) {
		addRequest.LogFieldInstanceName = plan.LogFieldInstanceName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldStartupid) {
		addRequest.LogFieldStartupid = plan.LogFieldStartupid.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldCategory) {
		addRequest.LogFieldCategory = plan.LogFieldCategory.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldSeverity) {
		addRequest.LogFieldSeverity = plan.LogFieldSeverity.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldMessageID) {
		addRequest.LogFieldMessageID = plan.LogFieldMessageID.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldMessage) {
		addRequest.LogFieldMessage = plan.LogFieldMessage.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a ErrorLogFieldMappingResponse object into the model struct
func readErrorLogFieldMappingResponse(ctx context.Context, r *client.ErrorLogFieldMappingResponse, state *errorLogFieldMappingResourceModel, expectedValues *errorLogFieldMappingResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.LogFieldTimestamp = internaltypes.StringTypeOrNil(r.LogFieldTimestamp, internaltypes.IsEmptyString(expectedValues.LogFieldTimestamp))
	state.LogFieldProductName = internaltypes.StringTypeOrNil(r.LogFieldProductName, internaltypes.IsEmptyString(expectedValues.LogFieldProductName))
	state.LogFieldInstanceName = internaltypes.StringTypeOrNil(r.LogFieldInstanceName, internaltypes.IsEmptyString(expectedValues.LogFieldInstanceName))
	state.LogFieldStartupid = internaltypes.StringTypeOrNil(r.LogFieldStartupid, internaltypes.IsEmptyString(expectedValues.LogFieldStartupid))
	state.LogFieldCategory = internaltypes.StringTypeOrNil(r.LogFieldCategory, internaltypes.IsEmptyString(expectedValues.LogFieldCategory))
	state.LogFieldSeverity = internaltypes.StringTypeOrNil(r.LogFieldSeverity, internaltypes.IsEmptyString(expectedValues.LogFieldSeverity))
	state.LogFieldMessageID = internaltypes.StringTypeOrNil(r.LogFieldMessageID, internaltypes.IsEmptyString(expectedValues.LogFieldMessageID))
	state.LogFieldMessage = internaltypes.StringTypeOrNil(r.LogFieldMessage, internaltypes.IsEmptyString(expectedValues.LogFieldMessage))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createErrorLogFieldMappingOperations(plan errorLogFieldMappingResourceModel, state errorLogFieldMappingResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldTimestamp, state.LogFieldTimestamp, "log-field-timestamp")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldProductName, state.LogFieldProductName, "log-field-product-name")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldInstanceName, state.LogFieldInstanceName, "log-field-instance-name")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldStartupid, state.LogFieldStartupid, "log-field-startupid")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldCategory, state.LogFieldCategory, "log-field-category")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldSeverity, state.LogFieldSeverity, "log-field-severity")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldMessageID, state.LogFieldMessageID, "log-field-message-id")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldMessage, state.LogFieldMessage, "log-field-message")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *errorLogFieldMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan errorLogFieldMappingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddErrorLogFieldMappingRequest(plan.Id.ValueString(),
		[]client.EnumerrorLogFieldMappingSchemaUrn{client.ENUMERRORLOGFIELDMAPPINGSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_FIELD_MAPPINGERROR})
	addOptionalErrorLogFieldMappingFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogFieldMappingApi.AddLogFieldMapping(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogFieldMappingRequest(
		client.AddErrorLogFieldMappingRequestAsAddLogFieldMappingRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogFieldMappingApi.AddLogFieldMappingExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Error Log Field Mapping", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state errorLogFieldMappingResourceModel
	readErrorLogFieldMappingResponse(ctx, addResponse.ErrorLogFieldMappingResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultErrorLogFieldMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan errorLogFieldMappingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogFieldMappingApi.GetLogFieldMapping(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Error Log Field Mapping", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state errorLogFieldMappingResourceModel
	readErrorLogFieldMappingResponse(ctx, readResponse.ErrorLogFieldMappingResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogFieldMappingApi.UpdateLogFieldMapping(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createErrorLogFieldMappingOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogFieldMappingApi.UpdateLogFieldMappingExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Error Log Field Mapping", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readErrorLogFieldMappingResponse(ctx, updateResponse.ErrorLogFieldMappingResponse, &state, &plan, &resp.Diagnostics)
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
func (r *errorLogFieldMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readErrorLogFieldMapping(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultErrorLogFieldMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readErrorLogFieldMapping(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readErrorLogFieldMapping(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state errorLogFieldMappingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogFieldMappingApi.GetLogFieldMapping(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Error Log Field Mapping", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readErrorLogFieldMappingResponse(ctx, readResponse.ErrorLogFieldMappingResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *errorLogFieldMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateErrorLogFieldMapping(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultErrorLogFieldMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateErrorLogFieldMapping(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateErrorLogFieldMapping(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan errorLogFieldMappingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state errorLogFieldMappingResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogFieldMappingApi.UpdateLogFieldMapping(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createErrorLogFieldMappingOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogFieldMappingApi.UpdateLogFieldMappingExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Error Log Field Mapping", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readErrorLogFieldMappingResponse(ctx, updateResponse.ErrorLogFieldMappingResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultErrorLogFieldMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *errorLogFieldMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state errorLogFieldMappingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogFieldMappingApi.DeleteLogFieldMappingExecute(r.apiClient.LogFieldMappingApi.DeleteLogFieldMapping(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Error Log Field Mapping", err, httpResp)
		return
	}
}

func (r *errorLogFieldMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importErrorLogFieldMapping(ctx, req, resp)
}

func (r *defaultErrorLogFieldMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importErrorLogFieldMapping(ctx, req, resp)
}

func importErrorLogFieldMapping(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
