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
	   _ resource.Resource                = &usernamePasswordAzureAuthenticationMethodResource{}
	   _ resource.ResourceWithConfigure   = &usernamePasswordAzureAuthenticationMethodResource{}
	   _ resource.ResourceWithImportState = &usernamePasswordAzureAuthenticationMethodResource{}
	   _ resource.Resource                = &defaultUsernamePasswordAzureAuthenticationMethodResource{}
	   _ resource.ResourceWithConfigure   = &defaultUsernamePasswordAzureAuthenticationMethodResource{}
	   _ resource.ResourceWithImportState = &defaultUsernamePasswordAzureAuthenticationMethodResource{}
)

// Create a Username Password Azure Authentication Method resource
func NewUsernamePasswordAzureAuthenticationMethodResource() resource.Resource {
	   return &usernamePasswordAzureAuthenticationMethodResource{}
}

func NewDefaultUsernamePasswordAzureAuthenticationMethodResource() resource.Resource {
	   return &defaultUsernamePasswordAzureAuthenticationMethodResource{}
}

// usernamePasswordAzureAuthenticationMethodResource is the resource implementation.
type usernamePasswordAzureAuthenticationMethodResource struct {
	   providerConfig internaltypes.ProviderConfiguration
	   apiClient      *client.APIClient
}

// defaultUsernamePasswordAzureAuthenticationMethodResource is the resource implementation.
type defaultUsernamePasswordAzureAuthenticationMethodResource struct {
	   providerConfig internaltypes.ProviderConfiguration
	   apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *usernamePasswordAzureAuthenticationMethodResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	   resp.TypeName = req.ProviderTypeName + "_username_password_azure_authentication_method"
}

func (r *defaultUsernamePasswordAzureAuthenticationMethodResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	   resp.TypeName = req.ProviderTypeName + "_default_username_password_azure_authentication_method"
}

// Configure adds the provider configured client to the resource.
func (r *usernamePasswordAzureAuthenticationMethodResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
    r.providerConfig = providerCfg.ProviderConfig
    r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultUsernamePasswordAzureAuthenticationMethodResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
    r.providerConfig = providerCfg.ProviderConfig
    r.apiClient = providerCfg.ApiClientV9200
}

type usernamePasswordAzureAuthenticationMethodResourceModel struct {
    Id              types.String `tfsdk:"id"`
    LastUpdated     types.String `tfsdk:"last_updated"`
    Notifications   types.Set    `tfsdk:"notifications"`
    RequiredActions types.Set    `tfsdk:"required_actions"`
    TenantID types.String `tfsdk:"tenant_id"`
    ClientID types.String `tfsdk:"client_id"`
    Username types.String `tfsdk:"username"`
    Password types.String `tfsdk:"password"`
    Description types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *usernamePasswordAzureAuthenticationMethodResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    usernamePasswordAzureAuthenticationMethodSchema(ctx, req, resp, false)
}

func (r *defaultUsernamePasswordAzureAuthenticationMethodResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    usernamePasswordAzureAuthenticationMethodSchema(ctx, req, resp, true)
}

func usernamePasswordAzureAuthenticationMethodSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
    schema := schema.Schema{
        Description: "Manages a Username Password Azure Authentication Method.",
        Attributes: map[string]schema.Attribute{
            "tenant_id": schema.StringAttribute{
                Description: "The tenant ID to use to authenticate.",
                Required:    true,
            },
            "client_id": schema.StringAttribute{
                Description: "The client ID to use to authenticate.",
                Required:    true,
            },
            "username": schema.StringAttribute{
                Description: "The username for the user to authenticate.",
                Required:    true,
            },
            "password": schema.StringAttribute{
                Description: "The password for the user to authenticate.",
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
func addOptionalUsernamePasswordAzureAuthenticationMethodFields(ctx context.Context, addRequest *client.AddUsernamePasswordAzureAuthenticationMethodRequest, plan usernamePasswordAzureAuthenticationMethodResourceModel) {
    // Empty strings are treated as equivalent to null
    if internaltypes.IsNonEmptyString(plan.Description) {
        addRequest.Description = plan.Description.ValueStringPointer()
    }
}

// Read a UsernamePasswordAzureAuthenticationMethodResponse object into the model struct
func readUsernamePasswordAzureAuthenticationMethodResponse(ctx context.Context, r *client.UsernamePasswordAzureAuthenticationMethodResponse, state *usernamePasswordAzureAuthenticationMethodResourceModel, expectedValues *usernamePasswordAzureAuthenticationMethodResourceModel, diagnostics *diag.Diagnostics) {
    state.Id = types.StringValue(r.Id)
    state.TenantID = types.StringValue(r.TenantID)
    state.ClientID = types.StringValue(r.ClientID)
    state.Username = types.StringValue(r.Username)
    // Obscured values aren't returned from the PD Configuration API - just use the expected value
    state.Password = expectedValues.Password
    state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
    state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createUsernamePasswordAzureAuthenticationMethodOperations(plan usernamePasswordAzureAuthenticationMethodResourceModel, state usernamePasswordAzureAuthenticationMethodResourceModel) []client.Operation {
    var ops []client.Operation
    operations.AddStringOperationIfNecessary(&ops, plan.TenantID, state.TenantID, "tenant-id")
    operations.AddStringOperationIfNecessary(&ops, plan.ClientID, state.ClientID, "client-id")
    operations.AddStringOperationIfNecessary(&ops, plan.Username, state.Username, "username")
    operations.AddStringOperationIfNecessary(&ops, plan.Password, state.Password, "password")
    operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
    return ops
}

// Create a new resource
func (r *usernamePasswordAzureAuthenticationMethodResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Retrieve values from plan
    var plan usernamePasswordAzureAuthenticationMethodResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    addRequest := client.NewAddUsernamePasswordAzureAuthenticationMethodRequest(plan.Id.ValueString(),
        []client.EnumusernamePasswordAzureAuthenticationMethodSchemaUrn{client.ENUMUSERNAMEPASSWORDAZUREAUTHENTICATIONMETHODSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0AZURE_AUTHENTICATION_METHODUSERNAME_PASSWORD},
        plan.TenantID.ValueString(),
        plan.ClientID.ValueString(),
        plan.Username.ValueString(),
        plan.Password.ValueString())
    addOptionalUsernamePasswordAzureAuthenticationMethodFields(ctx, addRequest, plan)
    // Log request JSON
    requestJson, err := addRequest.MarshalJSON()
    if err == nil {
        tflog.Debug(ctx, "Add request: "+string(requestJson))
    }
    apiAddRequest := r.apiClient.AzureAuthenticationMethodApi.AddAzureAuthenticationMethod(
        config.ProviderBasicAuthContext(ctx, r.providerConfig))
    apiAddRequest = apiAddRequest.AddAzureAuthenticationMethodRequest(
        client.AddUsernamePasswordAzureAuthenticationMethodRequestAsAddAzureAuthenticationMethodRequest(addRequest))

    addResponse, httpResp, err := r.apiClient.AzureAuthenticationMethodApi.AddAzureAuthenticationMethodExecute(apiAddRequest)
    if err != nil {
        config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Username Password Azure Authentication Method", err, httpResp)
        return
    }

    // Log response JSON
    responseJson, err := addResponse.MarshalJSON()
    if err == nil {
        tflog.Debug(ctx, "Add response: "+string(responseJson))
    }

    // Read the response into the state
    var state usernamePasswordAzureAuthenticationMethodResourceModel
    readUsernamePasswordAzureAuthenticationMethodResponse(ctx, addResponse.UsernamePasswordAzureAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultUsernamePasswordAzureAuthenticationMethodResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Retrieve values from plan
    var plan usernamePasswordAzureAuthenticationMethodResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    readResponse, httpResp, err := r.apiClient.AzureAuthenticationMethodApi.GetAzureAuthenticationMethod(
        config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
    if err != nil {
        config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Username Password Azure Authentication Method", err, httpResp)
        return
    }

    // Log response JSON
    responseJson, err := readResponse.MarshalJSON()
    if err == nil {
        tflog.Debug(ctx, "Read response: "+string(responseJson))
    }

    // Read the existing configuration
    var state usernamePasswordAzureAuthenticationMethodResourceModel
    readUsernamePasswordAzureAuthenticationMethodResponse(ctx, readResponse.UsernamePasswordAzureAuthenticationMethodResponse, &state, &state, &resp.Diagnostics)

    // Determine what changes are needed to match the plan
    updateRequest := r.apiClient.AzureAuthenticationMethodApi.UpdateAzureAuthenticationMethod(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	   ops := createUsernamePasswordAzureAuthenticationMethodOperations(plan, state)
	   if len(ops) > 0 {
	   	   updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
	   	   // Log operations
	   	   operations.LogUpdateOperations(ctx, ops)

        updateResponse, httpResp, err := r.apiClient.AzureAuthenticationMethodApi.UpdateAzureAuthenticationMethodExecute(updateRequest)
        if err != nil {
            config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Username Password Azure Authentication Method", err, httpResp)
            return
        }

	       // Log response JSON
        responseJson, err := updateResponse.MarshalJSON()
        if err == nil {
        	   tflog.Debug(ctx, "Update response: "+string(responseJson))
        }

	       // Read the response
    readUsernamePasswordAzureAuthenticationMethodResponse(ctx, updateResponse.UsernamePasswordAzureAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
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
func (r *usernamePasswordAzureAuthenticationMethodResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    readUsernamePasswordAzureAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultUsernamePasswordAzureAuthenticationMethodResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    readUsernamePasswordAzureAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readUsernamePasswordAzureAuthenticationMethod(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
    // Get current state
  var state usernamePasswordAzureAuthenticationMethodResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    readResponse, httpResp, err := apiClient.AzureAuthenticationMethodApi.GetAzureAuthenticationMethod(
        config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
    if err != nil {
        config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Username Password Azure Authentication Method", err, httpResp)
        return
    }

    // Log response JSON
    responseJson, err := readResponse.MarshalJSON()
    if err == nil {
        tflog.Debug(ctx, "Read response: "+string(responseJson))
    }

    // Read the response into the state
    readUsernamePasswordAzureAuthenticationMethodResponse(ctx, readResponse.UsernamePasswordAzureAuthenticationMethodResponse, &state, &state, &resp.Diagnostics)

    // Set refreshed state
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Update a resource
func (r *usernamePasswordAzureAuthenticationMethodResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    updateUsernamePasswordAzureAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultUsernamePasswordAzureAuthenticationMethodResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    updateUsernamePasswordAzureAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateUsernamePasswordAzureAuthenticationMethod(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
  	 // Retrieve values from plan
	   var plan usernamePasswordAzureAuthenticationMethodResourceModel
	   diags := req.Plan.Get(ctx, &plan)
	   resp.Diagnostics.Append(diags...)
	   if resp.Diagnostics.HasError() {
	       return
	   }

	   // Get the current state to see how any attributes are changing
	   var state usernamePasswordAzureAuthenticationMethodResourceModel
	   req.State.Get(ctx, &state)
	   updateRequest := apiClient.AzureAuthenticationMethodApi.UpdateAzureAuthenticationMethod(
        config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	   // Determine what update operations are necessary
	   ops := createUsernamePasswordAzureAuthenticationMethodOperations(plan, state)
	   if len(ops) > 0 {
	   	   updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
	   	   // Log operations
	   	   operations.LogUpdateOperations(ctx, ops)

        updateResponse, httpResp, err := apiClient.AzureAuthenticationMethodApi.UpdateAzureAuthenticationMethodExecute(updateRequest)
        if err != nil {
            config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Username Password Azure Authentication Method", err, httpResp)
            return
        }

	       // Log response JSON
        responseJson, err := updateResponse.MarshalJSON()
        if err == nil {
        	   tflog.Debug(ctx, "Update response: "+string(responseJson))
        }

	       // Read the response
        readUsernamePasswordAzureAuthenticationMethodResponse(ctx, updateResponse.UsernamePasswordAzureAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultUsernamePasswordAzureAuthenticationMethodResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // No implementation necessary
}


func (r *usernamePasswordAzureAuthenticationMethodResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Retrieve values from state
    var state usernamePasswordAzureAuthenticationMethodResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
    	   return
    }

    httpResp, err := r.apiClient.AzureAuthenticationMethodApi.DeleteAzureAuthenticationMethodExecute(r.apiClient.AzureAuthenticationMethodApi.DeleteAzureAuthenticationMethod(
        config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
    if err != nil {
        config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Username Password Azure Authentication Method", err, httpResp)
        return
    }
}

func (r *usernamePasswordAzureAuthenticationMethodResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    importUsernamePasswordAzureAuthenticationMethod(ctx, req, resp)
}

func (r *defaultUsernamePasswordAzureAuthenticationMethodResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    importUsernamePasswordAzureAuthenticationMethod(ctx, req, resp)
}

func importUsernamePasswordAzureAuthenticationMethod(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    // Retrieve import ID and save to id attribute
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

