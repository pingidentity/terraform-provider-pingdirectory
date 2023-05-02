package azureauthenticationmethod

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
	   _ resource.Resource                = &clientSecretAzureAuthenticationMethodResource{}
	   _ resource.ResourceWithConfigure   = &clientSecretAzureAuthenticationMethodResource{}
	   _ resource.ResourceWithImportState = &clientSecretAzureAuthenticationMethodResource{}
	   _ resource.Resource                = &defaultClientSecretAzureAuthenticationMethodResource{}
	   _ resource.ResourceWithConfigure   = &defaultClientSecretAzureAuthenticationMethodResource{}
	   _ resource.ResourceWithImportState = &defaultClientSecretAzureAuthenticationMethodResource{}
)

// Create a Client Secret Azure Authentication Method resource
func NewClientSecretAzureAuthenticationMethodResource() resource.Resource {
	   return &clientSecretAzureAuthenticationMethodResource{}
}

func NewDefaultClientSecretAzureAuthenticationMethodResource() resource.Resource {
	   return &defaultClientSecretAzureAuthenticationMethodResource{}
}

// clientSecretAzureAuthenticationMethodResource is the resource implementation.
type clientSecretAzureAuthenticationMethodResource struct {
	   providerConfig internaltypes.ProviderConfiguration
	   apiClient      *client.APIClient
}

// defaultClientSecretAzureAuthenticationMethodResource is the resource implementation.
type defaultClientSecretAzureAuthenticationMethodResource struct {
	   providerConfig internaltypes.ProviderConfiguration
	   apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *clientSecretAzureAuthenticationMethodResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	   resp.TypeName = req.ProviderTypeName + "_client_secret_azure_authentication_method"
}

func (r *defaultClientSecretAzureAuthenticationMethodResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	   resp.TypeName = req.ProviderTypeName + "_default_client_secret_azure_authentication_method"
}

// Configure adds the provider configured client to the resource.
func (r *clientSecretAzureAuthenticationMethodResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
    r.providerConfig = providerCfg.ProviderConfig
    r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultClientSecretAzureAuthenticationMethodResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
    r.providerConfig = providerCfg.ProviderConfig
    r.apiClient = providerCfg.ApiClientV9200
}

type clientSecretAzureAuthenticationMethodResourceModel struct {
    Id              types.String `tfsdk:"id"`
    LastUpdated     types.String `tfsdk:"last_updated"`
    Notifications   types.Set    `tfsdk:"notifications"`
    RequiredActions types.Set    `tfsdk:"required_actions"`
    TenantID types.String `tfsdk:"tenant_id"`
    ClientID types.String `tfsdk:"client_id"`
    ClientSecret types.String `tfsdk:"client_secret"`
    Description types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *clientSecretAzureAuthenticationMethodResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    clientSecretAzureAuthenticationMethodSchema(ctx, req, resp, false)
}

func (r *defaultClientSecretAzureAuthenticationMethodResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    clientSecretAzureAuthenticationMethodSchema(ctx, req, resp, true)
}

func clientSecretAzureAuthenticationMethodSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
    schema := schema.Schema{
        Description: "Manages a Client Secret Azure Authentication Method.",
        Attributes: map[string]schema.Attribute{
            "tenant_id": schema.StringAttribute{
                Description: "The tenant ID to use to authenticate.",
                Required:    true,
            },
            "client_id": schema.StringAttribute{
                Description: "The client ID to use to authenticate.",
                Required:    true,
            },
            "client_secret": schema.StringAttribute{
                Description: "The client secret to use to authenticate.",
                Required:    true,
                Sensitive:    true,
            },
            "description": schema.StringAttribute{
                Description: "A description for this Azure Authentication Method",
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
func addOptionalClientSecretAzureAuthenticationMethodFields(ctx context.Context, addRequest *client.AddClientSecretAzureAuthenticationMethodRequest, plan clientSecretAzureAuthenticationMethodResourceModel) {
    // Empty strings are treated as equivalent to null
    if internaltypes.IsNonEmptyString(plan.Description) {
        addRequest.Description = plan.Description.ValueStringPointer()
    }
}

// Read a ClientSecretAzureAuthenticationMethodResponse object into the model struct
func readClientSecretAzureAuthenticationMethodResponse(ctx context.Context, r *client.ClientSecretAzureAuthenticationMethodResponse, state *clientSecretAzureAuthenticationMethodResourceModel, expectedValues *clientSecretAzureAuthenticationMethodResourceModel, diagnostics *diag.Diagnostics) {
    state.Id = types.StringValue(r.Id)
    state.TenantID = types.StringValue(r.TenantID)
    state.ClientID = types.StringValue(r.ClientID)
    // Obscured values aren't returned from the PD Configuration API - just use the expected value
    state.ClientSecret = expectedValues.ClientSecret
    state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
    state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createClientSecretAzureAuthenticationMethodOperations(plan clientSecretAzureAuthenticationMethodResourceModel, state clientSecretAzureAuthenticationMethodResourceModel) []client.Operation {
    var ops []client.Operation
    operations.AddStringOperationIfNecessary(&ops, plan.TenantID, state.TenantID, "tenant-id")
    operations.AddStringOperationIfNecessary(&ops, plan.ClientID, state.ClientID, "client-id")
    operations.AddStringOperationIfNecessary(&ops, plan.ClientSecret, state.ClientSecret, "client-secret")
    operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
    return ops
}

// Create a new resource
func (r *clientSecretAzureAuthenticationMethodResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Retrieve values from plan
    var plan clientSecretAzureAuthenticationMethodResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    addRequest := client.NewAddClientSecretAzureAuthenticationMethodRequest(plan.Id.ValueString(),
        []client.EnumclientSecretAzureAuthenticationMethodSchemaUrn{client.ENUMCLIENTSECRETAZUREAUTHENTICATIONMETHODSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0AZURE_AUTHENTICATION_METHODCLIENT_SECRET},
        plan.TenantID.ValueString(),
        plan.ClientID.ValueString(),
        plan.ClientSecret.ValueString())
    addOptionalClientSecretAzureAuthenticationMethodFields(ctx, addRequest, plan)
    // Log request JSON
    requestJson, err := addRequest.MarshalJSON()
    if err == nil {
        tflog.Debug(ctx, "Add request: "+string(requestJson))
    }
    apiAddRequest := r.apiClient.AzureAuthenticationMethodApi.AddAzureAuthenticationMethod(
        config.ProviderBasicAuthContext(ctx, r.providerConfig))
    apiAddRequest = apiAddRequest.AddAzureAuthenticationMethodRequest(
        client.AddClientSecretAzureAuthenticationMethodRequestAsAddAzureAuthenticationMethodRequest(addRequest))

    addResponse, httpResp, err := r.apiClient.AzureAuthenticationMethodApi.AddAzureAuthenticationMethodExecute(apiAddRequest)
    if err != nil {
        config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Client Secret Azure Authentication Method", err, httpResp)
        return
    }

    // Log response JSON
    responseJson, err := addResponse.MarshalJSON()
    if err == nil {
        tflog.Debug(ctx, "Add response: "+string(responseJson))
    }

    // Read the response into the state
    var state clientSecretAzureAuthenticationMethodResourceModel
    readClientSecretAzureAuthenticationMethodResponse(ctx, addResponse.ClientSecretAzureAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultClientSecretAzureAuthenticationMethodResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Retrieve values from plan
    var plan clientSecretAzureAuthenticationMethodResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    readResponse, httpResp, err := r.apiClient.AzureAuthenticationMethodApi.GetAzureAuthenticationMethod(
        config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
    if err != nil {
        config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Client Secret Azure Authentication Method", err, httpResp)
        return
    }

    // Log response JSON
    responseJson, err := readResponse.MarshalJSON()
    if err == nil {
        tflog.Debug(ctx, "Read response: "+string(responseJson))
    }

    // Read the existing configuration
    var state clientSecretAzureAuthenticationMethodResourceModel
    readClientSecretAzureAuthenticationMethodResponse(ctx, readResponse.ClientSecretAzureAuthenticationMethodResponse, &state, &state, &resp.Diagnostics)

    // Determine what changes are needed to match the plan
    updateRequest := r.apiClient.AzureAuthenticationMethodApi.UpdateAzureAuthenticationMethod(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	   ops := createClientSecretAzureAuthenticationMethodOperations(plan, state)
	   if len(ops) > 0 {
	   	   updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
	   	   // Log operations
	   	   operations.LogUpdateOperations(ctx, ops)

        updateResponse, httpResp, err := r.apiClient.AzureAuthenticationMethodApi.UpdateAzureAuthenticationMethodExecute(updateRequest)
        if err != nil {
            config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Client Secret Azure Authentication Method", err, httpResp)
            return
        }

	       // Log response JSON
        responseJson, err := updateResponse.MarshalJSON()
        if err == nil {
        	   tflog.Debug(ctx, "Update response: "+string(responseJson))
        }

	       // Read the response
    readClientSecretAzureAuthenticationMethodResponse(ctx, updateResponse.ClientSecretAzureAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
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
func (r *clientSecretAzureAuthenticationMethodResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    readClientSecretAzureAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultClientSecretAzureAuthenticationMethodResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    readClientSecretAzureAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readClientSecretAzureAuthenticationMethod(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
    // Get current state
  var state clientSecretAzureAuthenticationMethodResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    readResponse, httpResp, err := apiClient.AzureAuthenticationMethodApi.GetAzureAuthenticationMethod(
        config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
    if err != nil {
        config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Client Secret Azure Authentication Method", err, httpResp)
        return
    }

    // Log response JSON
    responseJson, err := readResponse.MarshalJSON()
    if err == nil {
        tflog.Debug(ctx, "Read response: "+string(responseJson))
    }

    // Read the response into the state
    readClientSecretAzureAuthenticationMethodResponse(ctx, readResponse.ClientSecretAzureAuthenticationMethodResponse, &state, &state, &resp.Diagnostics)

    // Set refreshed state
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Update a resource
func (r *clientSecretAzureAuthenticationMethodResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    updateClientSecretAzureAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultClientSecretAzureAuthenticationMethodResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    updateClientSecretAzureAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateClientSecretAzureAuthenticationMethod(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
  	 // Retrieve values from plan
	   var plan clientSecretAzureAuthenticationMethodResourceModel
	   diags := req.Plan.Get(ctx, &plan)
	   resp.Diagnostics.Append(diags...)
	   if resp.Diagnostics.HasError() {
	       return
	   }

	   // Get the current state to see how any attributes are changing
	   var state clientSecretAzureAuthenticationMethodResourceModel
	   req.State.Get(ctx, &state)
	   updateRequest := apiClient.AzureAuthenticationMethodApi.UpdateAzureAuthenticationMethod(
        config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	   // Determine what update operations are necessary
	   ops := createClientSecretAzureAuthenticationMethodOperations(plan, state)
	   if len(ops) > 0 {
	   	   updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
	   	   // Log operations
	   	   operations.LogUpdateOperations(ctx, ops)

        updateResponse, httpResp, err := apiClient.AzureAuthenticationMethodApi.UpdateAzureAuthenticationMethodExecute(updateRequest)
        if err != nil {
            config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Client Secret Azure Authentication Method", err, httpResp)
            return
        }

	       // Log response JSON
        responseJson, err := updateResponse.MarshalJSON()
        if err == nil {
        	   tflog.Debug(ctx, "Update response: "+string(responseJson))
        }

	       // Read the response
        readClientSecretAzureAuthenticationMethodResponse(ctx, updateResponse.ClientSecretAzureAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultClientSecretAzureAuthenticationMethodResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // No implementation necessary
}


func (r *clientSecretAzureAuthenticationMethodResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Retrieve values from state
    var state clientSecretAzureAuthenticationMethodResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
    	   return
    }

    httpResp, err := r.apiClient.AzureAuthenticationMethodApi.DeleteAzureAuthenticationMethodExecute(r.apiClient.AzureAuthenticationMethodApi.DeleteAzureAuthenticationMethod(
        config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
    if err != nil {
        config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Client Secret Azure Authentication Method", err, httpResp)
        return
    }
}

func (r *clientSecretAzureAuthenticationMethodResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    importClientSecretAzureAuthenticationMethod(ctx, req, resp)
}

func (r *defaultClientSecretAzureAuthenticationMethodResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    importClientSecretAzureAuthenticationMethod(ctx, req, resp)
}

func importClientSecretAzureAuthenticationMethod(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    // Retrieve import ID and save to id attribute
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

