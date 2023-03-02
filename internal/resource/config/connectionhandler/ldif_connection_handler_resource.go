package connectionhandler

import (
	"context"
	"time"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ldifConnectionHandlerResource{}
	_ resource.ResourceWithConfigure   = &ldifConnectionHandlerResource{}
	_ resource.ResourceWithImportState = &ldifConnectionHandlerResource{}
	_ resource.Resource                = &defaultLdifConnectionHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultLdifConnectionHandlerResource{}
	_ resource.ResourceWithImportState = &defaultLdifConnectionHandlerResource{}
)

// Create a Ldif Connection Handler resource
func NewLdifConnectionHandlerResource() resource.Resource {
	return &ldifConnectionHandlerResource{}
}

func NewDefaultLdifConnectionHandlerResource() resource.Resource {
	return &defaultLdifConnectionHandlerResource{}
}

// ldifConnectionHandlerResource is the resource implementation.
type ldifConnectionHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultLdifConnectionHandlerResource is the resource implementation.
type defaultLdifConnectionHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *ldifConnectionHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ldif_connection_handler"
}

func (r *defaultLdifConnectionHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_ldif_connection_handler"
}

// Configure adds the provider configured client to the resource.
func (r *ldifConnectionHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultLdifConnectionHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type ldifConnectionHandlerResourceModel struct {
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	AllowedClient   types.Set    `tfsdk:"allowed_client"`
	DeniedClient    types.Set    `tfsdk:"denied_client"`
	LdifDirectory   types.String `tfsdk:"ldif_directory"`
	PollInterval    types.String `tfsdk:"poll_interval"`
	Description     types.String `tfsdk:"description"`
	Enabled         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *ldifConnectionHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ldifConnectionHandlerSchema(ctx, req, resp, false)
}

func (r *defaultLdifConnectionHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ldifConnectionHandlerSchema(ctx, req, resp, true)
}

func ldifConnectionHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Ldif Connection Handler.",
		Attributes: map[string]schema.Attribute{
			"allowed_client": schema.SetAttribute{
				Description: "Specifies a set of address masks that determines the addresses of the clients that are allowed to establish connections to this connection handler.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"denied_client": schema.SetAttribute{
				Description: "Specifies a set of address masks that determines the addresses of the clients that are not allowed to establish connections to this connection handler.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"ldif_directory": schema.StringAttribute{
				Description: "Specifies the path to the directory in which the LDIF files should be placed.",
				Optional:    true,
				Computed:    true,
			},
			"poll_interval": schema.StringAttribute{
				Description: "Specifies how frequently the LDIF connection handler should check the LDIF directory to determine whether a new LDIF file has been added.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Connection Handler",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Connection Handler is enabled.",
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

// Add optional fields to create request
func addOptionalLdifConnectionHandlerFields(ctx context.Context, addRequest *client.AddLdifConnectionHandlerRequest, plan ldifConnectionHandlerResourceModel) {
	if internaltypes.IsDefined(plan.AllowedClient) {
		var slice []string
		plan.AllowedClient.ElementsAs(ctx, &slice, false)
		addRequest.AllowedClient = slice
	}
	if internaltypes.IsDefined(plan.DeniedClient) {
		var slice []string
		plan.DeniedClient.ElementsAs(ctx, &slice, false)
		addRequest.DeniedClient = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LdifDirectory) {
		stringVal := plan.LdifDirectory.ValueString()
		addRequest.LdifDirectory = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PollInterval) {
		stringVal := plan.PollInterval.ValueString()
		addRequest.PollInterval = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
}

// Read a LdifConnectionHandlerResponse object into the model struct
func readLdifConnectionHandlerResponse(ctx context.Context, r *client.LdifConnectionHandlerResponse, state *ldifConnectionHandlerResourceModel, expectedValues *ldifConnectionHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.AllowedClient = internaltypes.GetStringSet(r.AllowedClient)
	state.DeniedClient = internaltypes.GetStringSet(r.DeniedClient)
	state.LdifDirectory = types.StringValue(r.LdifDirectory)
	state.PollInterval = types.StringValue(r.PollInterval)
	config.CheckMismatchedPDFormattedAttributes("poll_interval",
		expectedValues.PollInterval, state.PollInterval, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createLdifConnectionHandlerOperations(plan ldifConnectionHandlerResourceModel, state ldifConnectionHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedClient, state.AllowedClient, "allowed-client")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DeniedClient, state.DeniedClient, "denied-client")
	operations.AddStringOperationIfNecessary(&ops, plan.LdifDirectory, state.LdifDirectory, "ldif-directory")
	operations.AddStringOperationIfNecessary(&ops, plan.PollInterval, state.PollInterval, "poll-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *ldifConnectionHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ldifConnectionHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddLdifConnectionHandlerRequest(plan.Id.ValueString(),
		[]client.EnumldifConnectionHandlerSchemaUrn{client.ENUMLDIFCONNECTIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CONNECTION_HANDLERLDIF},
		plan.Enabled.ValueBool())
	addOptionalLdifConnectionHandlerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ConnectionHandlerApi.AddConnectionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddConnectionHandlerRequest(
		client.AddLdifConnectionHandlerRequestAsAddConnectionHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ConnectionHandlerApi.AddConnectionHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Ldif Connection Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state ldifConnectionHandlerResourceModel
	readLdifConnectionHandlerResponse(ctx, addResponse.LdifConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultLdifConnectionHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ldifConnectionHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ConnectionHandlerApi.GetConnectionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ldif Connection Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state ldifConnectionHandlerResourceModel
	readLdifConnectionHandlerResponse(ctx, readResponse.LdifConnectionHandlerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ConnectionHandlerApi.UpdateConnectionHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createLdifConnectionHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ConnectionHandlerApi.UpdateConnectionHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ldif Connection Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLdifConnectionHandlerResponse(ctx, updateResponse.LdifConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *ldifConnectionHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLdifConnectionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLdifConnectionHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLdifConnectionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readLdifConnectionHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state ldifConnectionHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ConnectionHandlerApi.GetConnectionHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ldif Connection Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readLdifConnectionHandlerResponse(ctx, readResponse.LdifConnectionHandlerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *ldifConnectionHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLdifConnectionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLdifConnectionHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLdifConnectionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateLdifConnectionHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan ldifConnectionHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state ldifConnectionHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ConnectionHandlerApi.UpdateConnectionHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createLdifConnectionHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ConnectionHandlerApi.UpdateConnectionHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ldif Connection Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLdifConnectionHandlerResponse(ctx, updateResponse.LdifConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultLdifConnectionHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *ldifConnectionHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ldifConnectionHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ConnectionHandlerApi.DeleteConnectionHandlerExecute(r.apiClient.ConnectionHandlerApi.DeleteConnectionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Ldif Connection Handler", err, httpResp)
		return
	}
}

func (r *ldifConnectionHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLdifConnectionHandler(ctx, req, resp)
}

func (r *defaultLdifConnectionHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLdifConnectionHandler(ctx, req, resp)
}

func importLdifConnectionHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
