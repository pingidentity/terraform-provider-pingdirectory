package accesstokenvalidator

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &pingFederateAccessTokenValidatorResource{}
	_ resource.ResourceWithConfigure   = &pingFederateAccessTokenValidatorResource{}
	_ resource.ResourceWithImportState = &pingFederateAccessTokenValidatorResource{}
	_ resource.Resource                = &defaultPingFederateAccessTokenValidatorResource{}
	_ resource.ResourceWithConfigure   = &defaultPingFederateAccessTokenValidatorResource{}
	_ resource.ResourceWithImportState = &defaultPingFederateAccessTokenValidatorResource{}
)

// Create a Ping Federate Access Token Validator resource
func NewPingFederateAccessTokenValidatorResource() resource.Resource {
	return &pingFederateAccessTokenValidatorResource{}
}

func NewDefaultPingFederateAccessTokenValidatorResource() resource.Resource {
	return &defaultPingFederateAccessTokenValidatorResource{}
}

// pingFederateAccessTokenValidatorResource is the resource implementation.
type pingFederateAccessTokenValidatorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultPingFederateAccessTokenValidatorResource is the resource implementation.
type defaultPingFederateAccessTokenValidatorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *pingFederateAccessTokenValidatorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ping_federate_access_token_validator"
}

func (r *defaultPingFederateAccessTokenValidatorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_ping_federate_access_token_validator"
}

// Configure adds the provider configured client to the resource.
func (r *pingFederateAccessTokenValidatorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultPingFederateAccessTokenValidatorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type pingFederateAccessTokenValidatorResourceModel struct {
	Id                             types.String `tfsdk:"id"`
	LastUpdated                    types.String `tfsdk:"last_updated"`
	Notifications                  types.Set    `tfsdk:"notifications"`
	RequiredActions                types.Set    `tfsdk:"required_actions"`
	ClientID                       types.String `tfsdk:"client_id"`
	ClientSecret                   types.String `tfsdk:"client_secret"`
	ClientSecretPassphraseProvider types.String `tfsdk:"client_secret_passphrase_provider"`
	IncludeAudParameter            types.Bool   `tfsdk:"include_aud_parameter"`
	AccessTokenManagerID           types.String `tfsdk:"access_token_manager_id"`
	EndpointCacheRefresh           types.String `tfsdk:"endpoint_cache_refresh"`
	EvaluationOrderIndex           types.Int64  `tfsdk:"evaluation_order_index"`
	AuthorizationServer            types.String `tfsdk:"authorization_server"`
	IdentityMapper                 types.String `tfsdk:"identity_mapper"`
	SubjectClaimName               types.String `tfsdk:"subject_claim_name"`
	Description                    types.String `tfsdk:"description"`
	Enabled                        types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *pingFederateAccessTokenValidatorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	pingFederateAccessTokenValidatorSchema(ctx, req, resp, false)
}

func (r *defaultPingFederateAccessTokenValidatorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	pingFederateAccessTokenValidatorSchema(ctx, req, resp, true)
}

func pingFederateAccessTokenValidatorSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Ping Federate Access Token Validator.",
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Description: "The client identifier to use when authenticating to the PingFederate authorization server.",
				Required:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "The client secret to use when authenticating to the PingFederate authorization server.",
				Optional:    true,
				Sensitive:   true,
			},
			"client_secret_passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider for obtaining the client secret to use when authenticating to the PingFederate authorization server.",
				Optional:    true,
			},
			"include_aud_parameter": schema.BoolAttribute{
				Description: "Whether to include the incoming request URL as the \"aud\" parameter when calling the PingFederate introspection endpoint. This property is ignored if the access-token-manager-id property is set.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"access_token_manager_id": schema.StringAttribute{
				Description: "The Access Token Manager instance ID to specify when calling the PingFederate introspection endpoint. If this property is set the include-aud-parameter property is ignored.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"endpoint_cache_refresh": schema.StringAttribute{
				Description: "How often the Access Token Validator should refresh its stored value of the PingFederate server's token introspection endpoint.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"evaluation_order_index": schema.Int64Attribute{
				Description: "When multiple Ping Federate Access Token Validators are defined for a single Directory Server, this property determines the evaluation order for determining the correct validator class for an access token received by the Directory Server. Values of this property must be unique among all Ping Federate Access Token Validators defined within Directory Server but not necessarily contiguous. Ping Federate Access Token Validators with a smaller value will be evaluated first to determine if they are able to validate the access token.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"authorization_server": schema.StringAttribute{
				Description: "Specifies the external server that will be used to aid in validating access tokens. In most cases this will be the Authorization Server that minted the token.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"identity_mapper": schema.StringAttribute{
				Description: "Specifies the name of the Identity Mapper that should be used for associating user entries with Bearer token subject names. The claim name from which to obtain the subject (i.e. the currently logged-in user) may be configured using the subject-claim-name property.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subject_claim_name": schema.StringAttribute{
				Description: "The name of the token claim that contains the subject, i.e. the logged-in user in an access token. This property goes hand-in-hand with the identity-mapper property and tells the Identity Mapper which field to use to look up the user entry on the server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Access Token Validator",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Access Token Validator is enabled for use in Directory Server.",
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
func addOptionalPingFederateAccessTokenValidatorFields(ctx context.Context, addRequest *client.AddPingFederateAccessTokenValidatorRequest, plan pingFederateAccessTokenValidatorResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ClientSecret) {
		stringVal := plan.ClientSecret.ValueString()
		addRequest.ClientSecret = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ClientSecretPassphraseProvider) {
		stringVal := plan.ClientSecretPassphraseProvider.ValueString()
		addRequest.ClientSecretPassphraseProvider = &stringVal
	}
	if internaltypes.IsDefined(plan.IncludeAudParameter) {
		boolVal := plan.IncludeAudParameter.ValueBool()
		addRequest.IncludeAudParameter = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccessTokenManagerID) {
		stringVal := plan.AccessTokenManagerID.ValueString()
		addRequest.AccessTokenManagerID = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EndpointCacheRefresh) {
		stringVal := plan.EndpointCacheRefresh.ValueString()
		addRequest.EndpointCacheRefresh = &stringVal
	}
	if internaltypes.IsDefined(plan.EvaluationOrderIndex) {
		intVal := int32(plan.EvaluationOrderIndex.ValueInt64())
		addRequest.EvaluationOrderIndex = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuthorizationServer) {
		stringVal := plan.AuthorizationServer.ValueString()
		addRequest.AuthorizationServer = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IdentityMapper) {
		stringVal := plan.IdentityMapper.ValueString()
		addRequest.IdentityMapper = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SubjectClaimName) {
		stringVal := plan.SubjectClaimName.ValueString()
		addRequest.SubjectClaimName = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
}

// Read a PingFederateAccessTokenValidatorResponse object into the model struct
func readPingFederateAccessTokenValidatorResponse(ctx context.Context, r *client.PingFederateAccessTokenValidatorResponse, state *pingFederateAccessTokenValidatorResourceModel, expectedValues *pingFederateAccessTokenValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ClientID = types.StringValue(r.ClientID)
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.ClientSecret = expectedValues.ClientSecret
	state.ClientSecretPassphraseProvider = internaltypes.StringTypeOrNil(r.ClientSecretPassphraseProvider, internaltypes.IsEmptyString(expectedValues.ClientSecretPassphraseProvider))
	state.IncludeAudParameter = internaltypes.BoolTypeOrNil(r.IncludeAudParameter)
	state.AccessTokenManagerID = internaltypes.StringTypeOrNil(r.AccessTokenManagerID, internaltypes.IsEmptyString(expectedValues.AccessTokenManagerID))
	state.EndpointCacheRefresh = internaltypes.StringTypeOrNil(r.EndpointCacheRefresh, internaltypes.IsEmptyString(expectedValues.EndpointCacheRefresh))
	config.CheckMismatchedPDFormattedAttributes("endpoint_cache_refresh",
		expectedValues.EndpointCacheRefresh, state.EndpointCacheRefresh, diagnostics)
	state.EvaluationOrderIndex = types.Int64Value(int64(r.EvaluationOrderIndex))
	state.AuthorizationServer = internaltypes.StringTypeOrNil(r.AuthorizationServer, internaltypes.IsEmptyString(expectedValues.AuthorizationServer))
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, internaltypes.IsEmptyString(expectedValues.IdentityMapper))
	state.SubjectClaimName = internaltypes.StringTypeOrNil(r.SubjectClaimName, internaltypes.IsEmptyString(expectedValues.SubjectClaimName))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createPingFederateAccessTokenValidatorOperations(plan pingFederateAccessTokenValidatorResourceModel, state pingFederateAccessTokenValidatorResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ClientID, state.ClientID, "client-id")
	operations.AddStringOperationIfNecessary(&ops, plan.ClientSecret, state.ClientSecret, "client-secret")
	operations.AddStringOperationIfNecessary(&ops, plan.ClientSecretPassphraseProvider, state.ClientSecretPassphraseProvider, "client-secret-passphrase-provider")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeAudParameter, state.IncludeAudParameter, "include-aud-parameter")
	operations.AddStringOperationIfNecessary(&ops, plan.AccessTokenManagerID, state.AccessTokenManagerID, "access-token-manager-id")
	operations.AddStringOperationIfNecessary(&ops, plan.EndpointCacheRefresh, state.EndpointCacheRefresh, "endpoint-cache-refresh")
	operations.AddInt64OperationIfNecessary(&ops, plan.EvaluationOrderIndex, state.EvaluationOrderIndex, "evaluation-order-index")
	operations.AddStringOperationIfNecessary(&ops, plan.AuthorizationServer, state.AuthorizationServer, "authorization-server")
	operations.AddStringOperationIfNecessary(&ops, plan.IdentityMapper, state.IdentityMapper, "identity-mapper")
	operations.AddStringOperationIfNecessary(&ops, plan.SubjectClaimName, state.SubjectClaimName, "subject-claim-name")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *pingFederateAccessTokenValidatorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan pingFederateAccessTokenValidatorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddPingFederateAccessTokenValidatorRequest(plan.Id.ValueString(),
		[]client.EnumpingFederateAccessTokenValidatorSchemaUrn{client.ENUMPINGFEDERATEACCESSTOKENVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ACCESS_TOKEN_VALIDATORPING_FEDERATE},
		plan.ClientID.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalPingFederateAccessTokenValidatorFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AccessTokenValidatorApi.AddAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAccessTokenValidatorRequest(
		client.AddPingFederateAccessTokenValidatorRequestAsAddAccessTokenValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AccessTokenValidatorApi.AddAccessTokenValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Ping Federate Access Token Validator", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pingFederateAccessTokenValidatorResourceModel
	readPingFederateAccessTokenValidatorResponse(ctx, addResponse.PingFederateAccessTokenValidatorResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultPingFederateAccessTokenValidatorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan pingFederateAccessTokenValidatorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AccessTokenValidatorApi.GetAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ping Federate Access Token Validator", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state pingFederateAccessTokenValidatorResourceModel
	readPingFederateAccessTokenValidatorResponse(ctx, readResponse.PingFederateAccessTokenValidatorResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.AccessTokenValidatorApi.UpdateAccessTokenValidator(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createPingFederateAccessTokenValidatorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AccessTokenValidatorApi.UpdateAccessTokenValidatorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ping Federate Access Token Validator", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readPingFederateAccessTokenValidatorResponse(ctx, updateResponse.PingFederateAccessTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *pingFederateAccessTokenValidatorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPingFederateAccessTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPingFederateAccessTokenValidatorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPingFederateAccessTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readPingFederateAccessTokenValidator(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state pingFederateAccessTokenValidatorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.AccessTokenValidatorApi.GetAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ping Federate Access Token Validator", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readPingFederateAccessTokenValidatorResponse(ctx, readResponse.PingFederateAccessTokenValidatorResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *pingFederateAccessTokenValidatorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePingFederateAccessTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPingFederateAccessTokenValidatorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePingFederateAccessTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updatePingFederateAccessTokenValidator(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan pingFederateAccessTokenValidatorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state pingFederateAccessTokenValidatorResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.AccessTokenValidatorApi.UpdateAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createPingFederateAccessTokenValidatorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.AccessTokenValidatorApi.UpdateAccessTokenValidatorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ping Federate Access Token Validator", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readPingFederateAccessTokenValidatorResponse(ctx, updateResponse.PingFederateAccessTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultPingFederateAccessTokenValidatorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *pingFederateAccessTokenValidatorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state pingFederateAccessTokenValidatorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.AccessTokenValidatorApi.DeleteAccessTokenValidatorExecute(r.apiClient.AccessTokenValidatorApi.DeleteAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Ping Federate Access Token Validator", err, httpResp)
		return
	}
}

func (r *pingFederateAccessTokenValidatorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPingFederateAccessTokenValidator(ctx, req, resp)
}

func (r *defaultPingFederateAccessTokenValidatorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPingFederateAccessTokenValidator(ctx, req, resp)
}

func importPingFederateAccessTokenValidator(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
