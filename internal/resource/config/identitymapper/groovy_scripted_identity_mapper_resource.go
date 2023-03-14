package identitymapper

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
	_ resource.Resource                = &groovyScriptedIdentityMapperResource{}
	_ resource.ResourceWithConfigure   = &groovyScriptedIdentityMapperResource{}
	_ resource.ResourceWithImportState = &groovyScriptedIdentityMapperResource{}
	_ resource.Resource                = &defaultGroovyScriptedIdentityMapperResource{}
	_ resource.ResourceWithConfigure   = &defaultGroovyScriptedIdentityMapperResource{}
	_ resource.ResourceWithImportState = &defaultGroovyScriptedIdentityMapperResource{}
)

// Create a Groovy Scripted Identity Mapper resource
func NewGroovyScriptedIdentityMapperResource() resource.Resource {
	return &groovyScriptedIdentityMapperResource{}
}

func NewDefaultGroovyScriptedIdentityMapperResource() resource.Resource {
	return &defaultGroovyScriptedIdentityMapperResource{}
}

// groovyScriptedIdentityMapperResource is the resource implementation.
type groovyScriptedIdentityMapperResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultGroovyScriptedIdentityMapperResource is the resource implementation.
type defaultGroovyScriptedIdentityMapperResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *groovyScriptedIdentityMapperResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groovy_scripted_identity_mapper"
}

func (r *defaultGroovyScriptedIdentityMapperResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_groovy_scripted_identity_mapper"
}

// Configure adds the provider configured client to the resource.
func (r *groovyScriptedIdentityMapperResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultGroovyScriptedIdentityMapperResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type groovyScriptedIdentityMapperResourceModel struct {
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	ScriptClass     types.String `tfsdk:"script_class"`
	ScriptArgument  types.Set    `tfsdk:"script_argument"`
	Description     types.String `tfsdk:"description"`
	Enabled         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *groovyScriptedIdentityMapperResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	groovyScriptedIdentityMapperSchema(ctx, req, resp, false)
}

func (r *defaultGroovyScriptedIdentityMapperResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	groovyScriptedIdentityMapperSchema(ctx, req, resp, true)
}

func groovyScriptedIdentityMapperSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Groovy Scripted Identity Mapper.",
		Attributes: map[string]schema.Attribute{
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Identity Mapper.",
				Required:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Identity Mapper. Each configuration property should be given in the form 'name=value'.",
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

// Add optional fields to create request
func addOptionalGroovyScriptedIdentityMapperFields(ctx context.Context, addRequest *client.AddGroovyScriptedIdentityMapperRequest, plan groovyScriptedIdentityMapperResourceModel) {
	if internaltypes.IsDefined(plan.ScriptArgument) {
		var slice []string
		plan.ScriptArgument.ElementsAs(ctx, &slice, false)
		addRequest.ScriptArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
}

// Read a GroovyScriptedIdentityMapperResponse object into the model struct
func readGroovyScriptedIdentityMapperResponse(ctx context.Context, r *client.GroovyScriptedIdentityMapperResponse, state *groovyScriptedIdentityMapperResourceModel, expectedValues *groovyScriptedIdentityMapperResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createGroovyScriptedIdentityMapperOperations(plan groovyScriptedIdentityMapperResourceModel, state groovyScriptedIdentityMapperResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ScriptClass, state.ScriptClass, "script-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScriptArgument, state.ScriptArgument, "script-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *groovyScriptedIdentityMapperResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan groovyScriptedIdentityMapperResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddGroovyScriptedIdentityMapperRequest(plan.Id.ValueString(),
		[]client.EnumgroovyScriptedIdentityMapperSchemaUrn{client.ENUMGROOVYSCRIPTEDIDENTITYMAPPERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0IDENTITY_MAPPERGROOVY_SCRIPTED},
		plan.ScriptClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalGroovyScriptedIdentityMapperFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.IdentityMapperApi.AddIdentityMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddIdentityMapperRequest(
		client.AddGroovyScriptedIdentityMapperRequestAsAddIdentityMapperRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.IdentityMapperApi.AddIdentityMapperExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Groovy Scripted Identity Mapper", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state groovyScriptedIdentityMapperResourceModel
	readGroovyScriptedIdentityMapperResponse(ctx, addResponse.GroovyScriptedIdentityMapperResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultGroovyScriptedIdentityMapperResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan groovyScriptedIdentityMapperResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.IdentityMapperApi.GetIdentityMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Groovy Scripted Identity Mapper", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state groovyScriptedIdentityMapperResourceModel
	readGroovyScriptedIdentityMapperResponse(ctx, readResponse.GroovyScriptedIdentityMapperResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.IdentityMapperApi.UpdateIdentityMapper(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createGroovyScriptedIdentityMapperOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.IdentityMapperApi.UpdateIdentityMapperExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Groovy Scripted Identity Mapper", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGroovyScriptedIdentityMapperResponse(ctx, updateResponse.GroovyScriptedIdentityMapperResponse, &state, &plan, &resp.Diagnostics)
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
func (r *groovyScriptedIdentityMapperResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readGroovyScriptedIdentityMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultGroovyScriptedIdentityMapperResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readGroovyScriptedIdentityMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readGroovyScriptedIdentityMapper(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state groovyScriptedIdentityMapperResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.IdentityMapperApi.GetIdentityMapper(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Groovy Scripted Identity Mapper", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readGroovyScriptedIdentityMapperResponse(ctx, readResponse.GroovyScriptedIdentityMapperResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *groovyScriptedIdentityMapperResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateGroovyScriptedIdentityMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultGroovyScriptedIdentityMapperResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateGroovyScriptedIdentityMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateGroovyScriptedIdentityMapper(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan groovyScriptedIdentityMapperResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state groovyScriptedIdentityMapperResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.IdentityMapperApi.UpdateIdentityMapper(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createGroovyScriptedIdentityMapperOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.IdentityMapperApi.UpdateIdentityMapperExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Groovy Scripted Identity Mapper", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGroovyScriptedIdentityMapperResponse(ctx, updateResponse.GroovyScriptedIdentityMapperResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultGroovyScriptedIdentityMapperResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *groovyScriptedIdentityMapperResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state groovyScriptedIdentityMapperResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.IdentityMapperApi.DeleteIdentityMapperExecute(r.apiClient.IdentityMapperApi.DeleteIdentityMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Groovy Scripted Identity Mapper", err, httpResp)
		return
	}
}

func (r *groovyScriptedIdentityMapperResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importGroovyScriptedIdentityMapper(ctx, req, resp)
}

func (r *defaultGroovyScriptedIdentityMapperResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importGroovyScriptedIdentityMapper(ctx, req, resp)
}

func importGroovyScriptedIdentityMapper(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
