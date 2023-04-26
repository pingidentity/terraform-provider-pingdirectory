package webapplicationextension

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
	_ resource.Resource                = &genericWebApplicationExtensionResource{}
	_ resource.ResourceWithConfigure   = &genericWebApplicationExtensionResource{}
	_ resource.ResourceWithImportState = &genericWebApplicationExtensionResource{}
	_ resource.Resource                = &defaultGenericWebApplicationExtensionResource{}
	_ resource.ResourceWithConfigure   = &defaultGenericWebApplicationExtensionResource{}
	_ resource.ResourceWithImportState = &defaultGenericWebApplicationExtensionResource{}
)

// Create a Generic Web Application Extension resource
func NewGenericWebApplicationExtensionResource() resource.Resource {
	return &genericWebApplicationExtensionResource{}
}

func NewDefaultGenericWebApplicationExtensionResource() resource.Resource {
	return &defaultGenericWebApplicationExtensionResource{}
}

// genericWebApplicationExtensionResource is the resource implementation.
type genericWebApplicationExtensionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultGenericWebApplicationExtensionResource is the resource implementation.
type defaultGenericWebApplicationExtensionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *genericWebApplicationExtensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_generic_web_application_extension"
}

func (r *defaultGenericWebApplicationExtensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_generic_web_application_extension"
}

// Configure adds the provider configured client to the resource.
func (r *genericWebApplicationExtensionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultGenericWebApplicationExtensionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type genericWebApplicationExtensionResourceModel struct {
	Id                       types.String `tfsdk:"id"`
	LastUpdated              types.String `tfsdk:"last_updated"`
	Notifications            types.Set    `tfsdk:"notifications"`
	RequiredActions          types.Set    `tfsdk:"required_actions"`
	Description              types.String `tfsdk:"description"`
	BaseContextPath          types.String `tfsdk:"base_context_path"`
	WarFile                  types.String `tfsdk:"war_file"`
	DocumentRootDirectory    types.String `tfsdk:"document_root_directory"`
	DeploymentDescriptorFile types.String `tfsdk:"deployment_descriptor_file"`
	TemporaryDirectory       types.String `tfsdk:"temporary_directory"`
	InitParameter            types.Set    `tfsdk:"init_parameter"`
}

// GetSchema defines the schema for the resource.
func (r *genericWebApplicationExtensionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	genericWebApplicationExtensionSchema(ctx, req, resp, false)
}

func (r *defaultGenericWebApplicationExtensionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	genericWebApplicationExtensionSchema(ctx, req, resp, true)
}

func genericWebApplicationExtensionSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Generic Web Application Extension.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Description: "A description for this Web Application Extension",
				Optional:    true,
			},
			"base_context_path": schema.StringAttribute{
				Description: "Specifies the base context path that should be used by HTTP clients to reference content. The value must start with a forward slash and at least one additional character and must represent a valid HTTP context path.",
				Required:    true,
			},
			"war_file": schema.StringAttribute{
				Description: "Specifies the path to a standard web application archive (WAR) file.",
				Optional:    true,
			},
			"document_root_directory": schema.StringAttribute{
				Description: "Specifies the path to the directory on the local filesystem containing the files to be served by this Web Application Extension. The path must exist, and it must be a directory.",
				Optional:    true,
			},
			"deployment_descriptor_file": schema.StringAttribute{
				Description: "Specifies the path to the deployment descriptor file when used with document-root-directory.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"temporary_directory": schema.StringAttribute{
				Description: "Specifies the path to the directory that may be used to store temporary files such as extracted WAR files and compiled JSP files.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"init_parameter": schema.SetAttribute{
				Description: "Specifies an initialization parameter to pass into the web application during startup.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
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
func (r genericWebApplicationExtensionResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("war_file"),
			path.MatchRoot("document_root_directory"),
		),
	}
}

// Add optional fields to create request
func addOptionalGenericWebApplicationExtensionFields(ctx context.Context, addRequest *client.AddGenericWebApplicationExtensionRequest, plan genericWebApplicationExtensionResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.WarFile) {
		addRequest.WarFile = plan.WarFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DocumentRootDirectory) {
		addRequest.DocumentRootDirectory = plan.DocumentRootDirectory.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DeploymentDescriptorFile) {
		addRequest.DeploymentDescriptorFile = plan.DeploymentDescriptorFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TemporaryDirectory) {
		addRequest.TemporaryDirectory = plan.TemporaryDirectory.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InitParameter) {
		var slice []string
		plan.InitParameter.ElementsAs(ctx, &slice, false)
		addRequest.InitParameter = slice
	}
}

// Read a GenericWebApplicationExtensionResponse object into the model struct
func readGenericWebApplicationExtensionResponse(ctx context.Context, r *client.GenericWebApplicationExtensionResponse, state *genericWebApplicationExtensionResourceModel, expectedValues *genericWebApplicationExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.WarFile = internaltypes.StringTypeOrNil(r.WarFile, internaltypes.IsEmptyString(expectedValues.WarFile))
	state.DocumentRootDirectory = internaltypes.StringTypeOrNil(r.DocumentRootDirectory, internaltypes.IsEmptyString(expectedValues.DocumentRootDirectory))
	state.DeploymentDescriptorFile = internaltypes.StringTypeOrNil(r.DeploymentDescriptorFile, internaltypes.IsEmptyString(expectedValues.DeploymentDescriptorFile))
	state.TemporaryDirectory = internaltypes.StringTypeOrNil(r.TemporaryDirectory, internaltypes.IsEmptyString(expectedValues.TemporaryDirectory))
	state.InitParameter = internaltypes.GetStringSet(r.InitParameter)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createGenericWebApplicationExtensionOperations(plan genericWebApplicationExtensionResourceModel, state genericWebApplicationExtensionResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.BaseContextPath, state.BaseContextPath, "base-context-path")
	operations.AddStringOperationIfNecessary(&ops, plan.WarFile, state.WarFile, "war-file")
	operations.AddStringOperationIfNecessary(&ops, plan.DocumentRootDirectory, state.DocumentRootDirectory, "document-root-directory")
	operations.AddStringOperationIfNecessary(&ops, plan.DeploymentDescriptorFile, state.DeploymentDescriptorFile, "deployment-descriptor-file")
	operations.AddStringOperationIfNecessary(&ops, plan.TemporaryDirectory, state.TemporaryDirectory, "temporary-directory")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.InitParameter, state.InitParameter, "init-parameter")
	return ops
}

// Create a new resource
func (r *genericWebApplicationExtensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan genericWebApplicationExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddGenericWebApplicationExtensionRequest(plan.Id.ValueString(),
		[]client.EnumgenericWebApplicationExtensionSchemaUrn{client.ENUMGENERICWEBAPPLICATIONEXTENSIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0WEB_APPLICATION_EXTENSIONGENERIC},
		plan.BaseContextPath.ValueString())
	addOptionalGenericWebApplicationExtensionFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.WebApplicationExtensionApi.AddWebApplicationExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddGenericWebApplicationExtensionRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.WebApplicationExtensionApi.AddWebApplicationExtensionExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Generic Web Application Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state genericWebApplicationExtensionResourceModel
	readGenericWebApplicationExtensionResponse(ctx, addResponse.GenericWebApplicationExtensionResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultGenericWebApplicationExtensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan genericWebApplicationExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.WebApplicationExtensionApi.GetWebApplicationExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Generic Web Application Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state genericWebApplicationExtensionResourceModel
	readGenericWebApplicationExtensionResponse(ctx, readResponse.GenericWebApplicationExtensionResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.WebApplicationExtensionApi.UpdateWebApplicationExtension(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createGenericWebApplicationExtensionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.WebApplicationExtensionApi.UpdateWebApplicationExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Generic Web Application Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGenericWebApplicationExtensionResponse(ctx, updateResponse.GenericWebApplicationExtensionResponse, &state, &plan, &resp.Diagnostics)
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
func (r *genericWebApplicationExtensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readGenericWebApplicationExtension(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultGenericWebApplicationExtensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readGenericWebApplicationExtension(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readGenericWebApplicationExtension(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state genericWebApplicationExtensionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.WebApplicationExtensionApi.GetWebApplicationExtension(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Generic Web Application Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readGenericWebApplicationExtensionResponse(ctx, readResponse.GenericWebApplicationExtensionResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *genericWebApplicationExtensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateGenericWebApplicationExtension(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultGenericWebApplicationExtensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateGenericWebApplicationExtension(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateGenericWebApplicationExtension(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan genericWebApplicationExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state genericWebApplicationExtensionResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.WebApplicationExtensionApi.UpdateWebApplicationExtension(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createGenericWebApplicationExtensionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.WebApplicationExtensionApi.UpdateWebApplicationExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Generic Web Application Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGenericWebApplicationExtensionResponse(ctx, updateResponse.GenericWebApplicationExtensionResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultGenericWebApplicationExtensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *genericWebApplicationExtensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state genericWebApplicationExtensionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.WebApplicationExtensionApi.DeleteWebApplicationExtensionExecute(r.apiClient.WebApplicationExtensionApi.DeleteWebApplicationExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Generic Web Application Extension", err, httpResp)
		return
	}
}

func (r *genericWebApplicationExtensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importGenericWebApplicationExtension(ctx, req, resp)
}

func (r *defaultGenericWebApplicationExtensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importGenericWebApplicationExtension(ctx, req, resp)
}

func importGenericWebApplicationExtension(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
