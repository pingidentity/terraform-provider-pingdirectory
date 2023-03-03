package externalserver

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &conjurExternalServerResource{}
	_ resource.ResourceWithConfigure   = &conjurExternalServerResource{}
	_ resource.ResourceWithImportState = &conjurExternalServerResource{}
	_ resource.Resource                = &defaultConjurExternalServerResource{}
	_ resource.ResourceWithConfigure   = &defaultConjurExternalServerResource{}
	_ resource.ResourceWithImportState = &defaultConjurExternalServerResource{}
)

// Create a Conjur External Server resource
func NewConjurExternalServerResource() resource.Resource {
	return &conjurExternalServerResource{}
}

func NewDefaultConjurExternalServerResource() resource.Resource {
	return &defaultConjurExternalServerResource{}
}

// conjurExternalServerResource is the resource implementation.
type conjurExternalServerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultConjurExternalServerResource is the resource implementation.
type defaultConjurExternalServerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *conjurExternalServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_conjur_external_server"
}

func (r *defaultConjurExternalServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_conjur_external_server"
}

// Configure adds the provider configured client to the resource.
func (r *conjurExternalServerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultConjurExternalServerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type conjurExternalServerResourceModel struct {
	Id                         types.String `tfsdk:"id"`
	LastUpdated                types.String `tfsdk:"last_updated"`
	Notifications              types.Set    `tfsdk:"notifications"`
	RequiredActions            types.Set    `tfsdk:"required_actions"`
	ConjurServerBaseURI        types.Set    `tfsdk:"conjur_server_base_uri"`
	ConjurAuthenticationMethod types.String `tfsdk:"conjur_authentication_method"`
	ConjurAccountName          types.String `tfsdk:"conjur_account_name"`
	TrustStoreFile             types.String `tfsdk:"trust_store_file"`
	TrustStorePin              types.String `tfsdk:"trust_store_pin"`
	TrustStoreType             types.String `tfsdk:"trust_store_type"`
	Description                types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *conjurExternalServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	conjurExternalServerSchema(ctx, req, resp, false)
}

func (r *defaultConjurExternalServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	conjurExternalServerSchema(ctx, req, resp, true)
}

func conjurExternalServerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Conjur External Server.",
		Attributes: map[string]schema.Attribute{
			"conjur_server_base_uri": schema.SetAttribute{
				Description: "The base URL needed to access the CyberArk Conjur server. The base URL should consist of the protocol (\"http\" or \"https\"), the server address (resolvable name or IP address), and the port number. For example, \"https://conjur.example.com:8443/\".",
				Required:    true,
				ElementType: types.StringType,
			},
			"conjur_authentication_method": schema.StringAttribute{
				Description: "The mechanism used to authenticate to the Conjur server.",
				Required:    true,
			},
			"conjur_account_name": schema.StringAttribute{
				Description: "The name of the account with which the desired secrets are associated.",
				Required:    true,
			},
			"trust_store_file": schema.StringAttribute{
				Description: "The path to a file containing the information needed to trust the certificate presented by the Conjur servers.",
				Optional:    true,
			},
			"trust_store_pin": schema.StringAttribute{
				Description: "The PIN needed to access the contents of the trust store. This is only required if a trust store file is required, and if that trust store requires a PIN to access its contents.",
				Optional:    true,
				Sensitive:   true,
			},
			"trust_store_type": schema.StringAttribute{
				Description: "The store type for the specified trust store file. The value should likely be one of \"JKS\", \"PKCS12\", or \"BCFKS\".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this External Server",
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
func addOptionalConjurExternalServerFields(ctx context.Context, addRequest *client.AddConjurExternalServerRequest, plan conjurExternalServerResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustStoreFile) {
		stringVal := plan.TrustStoreFile.ValueString()
		addRequest.TrustStoreFile = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustStorePin) {
		stringVal := plan.TrustStorePin.ValueString()
		addRequest.TrustStorePin = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustStoreType) {
		stringVal := plan.TrustStoreType.ValueString()
		addRequest.TrustStoreType = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
}

// Read a ConjurExternalServerResponse object into the model struct
func readConjurExternalServerResponse(ctx context.Context, r *client.ConjurExternalServerResponse, state *conjurExternalServerResourceModel, expectedValues *conjurExternalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ConjurServerBaseURI = internaltypes.GetStringSet(r.ConjurServerBaseURI)
	state.ConjurAuthenticationMethod = types.StringValue(r.ConjurAuthenticationMethod)
	state.ConjurAccountName = types.StringValue(r.ConjurAccountName)
	state.TrustStoreFile = internaltypes.StringTypeOrNil(r.TrustStoreFile, internaltypes.IsEmptyString(expectedValues.TrustStoreFile))
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.TrustStorePin = expectedValues.TrustStorePin
	state.TrustStoreType = internaltypes.StringTypeOrNil(r.TrustStoreType, internaltypes.IsEmptyString(expectedValues.TrustStoreType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createConjurExternalServerOperations(plan conjurExternalServerResourceModel, state conjurExternalServerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ConjurServerBaseURI, state.ConjurServerBaseURI, "conjur-server-base-uri")
	operations.AddStringOperationIfNecessary(&ops, plan.ConjurAuthenticationMethod, state.ConjurAuthenticationMethod, "conjur-authentication-method")
	operations.AddStringOperationIfNecessary(&ops, plan.ConjurAccountName, state.ConjurAccountName, "conjur-account-name")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStoreFile, state.TrustStoreFile, "trust-store-file")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStorePin, state.TrustStorePin, "trust-store-pin")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStoreType, state.TrustStoreType, "trust-store-type")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *conjurExternalServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan conjurExternalServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var ConjurServerBaseURISlice []string
	plan.ConjurServerBaseURI.ElementsAs(ctx, &ConjurServerBaseURISlice, false)
	addRequest := client.NewAddConjurExternalServerRequest(plan.Id.ValueString(),
		[]client.EnumconjurExternalServerSchemaUrn{client.ENUMCONJUREXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERCONJUR},
		ConjurServerBaseURISlice,
		plan.ConjurAuthenticationMethod.ValueString(),
		plan.ConjurAccountName.ValueString())
	addOptionalConjurExternalServerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExternalServerApi.AddExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExternalServerRequest(
		client.AddConjurExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerApi.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Conjur External Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state conjurExternalServerResourceModel
	readConjurExternalServerResponse(ctx, addResponse.ConjurExternalServerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultConjurExternalServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan conjurExternalServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ExternalServerApi.GetExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Conjur External Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state conjurExternalServerResourceModel
	readConjurExternalServerResponse(ctx, readResponse.ConjurExternalServerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ExternalServerApi.UpdateExternalServer(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createConjurExternalServerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ExternalServerApi.UpdateExternalServerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Conjur External Server", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readConjurExternalServerResponse(ctx, updateResponse.ConjurExternalServerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *conjurExternalServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readConjurExternalServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultConjurExternalServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readConjurExternalServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readConjurExternalServer(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state conjurExternalServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ExternalServerApi.GetExternalServer(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Conjur External Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readConjurExternalServerResponse(ctx, readResponse.ConjurExternalServerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *conjurExternalServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateConjurExternalServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultConjurExternalServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateConjurExternalServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateConjurExternalServer(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan conjurExternalServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state conjurExternalServerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ExternalServerApi.UpdateExternalServer(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createConjurExternalServerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ExternalServerApi.UpdateExternalServerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Conjur External Server", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readConjurExternalServerResponse(ctx, updateResponse.ConjurExternalServerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultConjurExternalServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *conjurExternalServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state conjurExternalServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ExternalServerApi.DeleteExternalServerExecute(r.apiClient.ExternalServerApi.DeleteExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Conjur External Server", err, httpResp)
		return
	}
}

func (r *conjurExternalServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importConjurExternalServer(ctx, req, resp)
}

func (r *defaultConjurExternalServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importConjurExternalServer(ctx, req, resp)
}

func importConjurExternalServer(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
