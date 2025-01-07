package velocitytemplateloader

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &velocityTemplateLoaderResource{}
	_ resource.ResourceWithConfigure   = &velocityTemplateLoaderResource{}
	_ resource.ResourceWithImportState = &velocityTemplateLoaderResource{}
	_ resource.Resource                = &defaultVelocityTemplateLoaderResource{}
	_ resource.ResourceWithConfigure   = &defaultVelocityTemplateLoaderResource{}
	_ resource.ResourceWithImportState = &defaultVelocityTemplateLoaderResource{}
)

// Create a Velocity Template Loader resource
func NewVelocityTemplateLoaderResource() resource.Resource {
	return &velocityTemplateLoaderResource{}
}

func NewDefaultVelocityTemplateLoaderResource() resource.Resource {
	return &defaultVelocityTemplateLoaderResource{}
}

// velocityTemplateLoaderResource is the resource implementation.
type velocityTemplateLoaderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultVelocityTemplateLoaderResource is the resource implementation.
type defaultVelocityTemplateLoaderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *velocityTemplateLoaderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_velocity_template_loader"
}

func (r *defaultVelocityTemplateLoaderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_velocity_template_loader"
}

// Configure adds the provider configured client to the resource.
func (r *velocityTemplateLoaderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultVelocityTemplateLoaderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type velocityTemplateLoaderResourceModel struct {
	Id                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	Notifications            types.Set    `tfsdk:"notifications"`
	RequiredActions          types.Set    `tfsdk:"required_actions"`
	Type                     types.String `tfsdk:"type"`
	HttpServletExtensionName types.String `tfsdk:"http_servlet_extension_name"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	EvaluationOrderIndex     types.Int64  `tfsdk:"evaluation_order_index"`
	MimeTypeMatcher          types.String `tfsdk:"mime_type_matcher"`
	MimeType                 types.String `tfsdk:"mime_type"`
	TemplateSuffix           types.String `tfsdk:"template_suffix"`
	TemplateDirectory        types.String `tfsdk:"template_directory"`
}

// GetSchema defines the schema for the resource.
func (r *velocityTemplateLoaderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	velocityTemplateLoaderSchema(ctx, req, resp, false)
}

func (r *defaultVelocityTemplateLoaderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	velocityTemplateLoaderSchema(ctx, req, resp, true)
}

func velocityTemplateLoaderSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Velocity Template Loader.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Velocity Template Loader resource. Options are ['velocity-template-loader']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("velocity-template-loader"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"velocity-template-loader"}...),
				},
			},
			"http_servlet_extension_name": schema.StringAttribute{
				Description: "Name of the parent HTTP Servlet Extension",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Velocity Template Loader is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"evaluation_order_index": schema.Int64Attribute{
				Description: "This property determines the evaluation order for determining the correct Velocity Template Loader to load a template for generating content for a particular request.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(9999),
			},
			"mime_type_matcher": schema.StringAttribute{
				Description: "Specifies a media type for matching Accept request-header values.",
				Required:    true,
			},
			"mime_type": schema.StringAttribute{
				Description: "Specifies a the value that will be used in the response's Content-Type header that indicates the type of content to return.",
				Optional:    true,
			},
			"template_suffix": schema.StringAttribute{
				Description: "Specifies the suffix to append to the requested resource name when searching for the template file with which to form a response.",
				Optional:    true,
			},
			"template_directory": schema.StringAttribute{
				Description: "Specifies the directory in which to search for the template files.",
				Optional:    true,
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Optional = false
		typeAttr.Required = false
		typeAttr.Computed = true
		typeAttr.PlanModifiers = []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type", "http_servlet_extension_name"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add optional fields to create request for velocity-template-loader velocity-template-loader
func addOptionalVelocityTemplateLoaderFields(ctx context.Context, addRequest *client.AddVelocityTemplateLoaderRequest, plan velocityTemplateLoaderResourceModel) {
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EvaluationOrderIndex) {
		addRequest.EvaluationOrderIndex = plan.EvaluationOrderIndex.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MimeType) {
		addRequest.MimeType = plan.MimeType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TemplateSuffix) {
		addRequest.TemplateSuffix = plan.TemplateSuffix.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TemplateDirectory) {
		addRequest.TemplateDirectory = plan.TemplateDirectory.ValueStringPointer()
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *velocityTemplateLoaderResourceModel) populateAllComputedStringAttributes() {
	if model.MimeTypeMatcher.IsUnknown() || model.MimeTypeMatcher.IsNull() {
		model.MimeTypeMatcher = types.StringValue("")
	}
	if model.TemplateDirectory.IsUnknown() || model.TemplateDirectory.IsNull() {
		model.TemplateDirectory = types.StringValue("")
	}
	if model.TemplateSuffix.IsUnknown() || model.TemplateSuffix.IsNull() {
		model.TemplateSuffix = types.StringValue("")
	}
	if model.MimeType.IsUnknown() || model.MimeType.IsNull() {
		model.MimeType = types.StringValue("")
	}
}

// Read a VelocityTemplateLoaderResponse object into the model struct
func readVelocityTemplateLoaderResponse(ctx context.Context, r *client.VelocityTemplateLoaderResponse, state *velocityTemplateLoaderResourceModel, expectedValues *velocityTemplateLoaderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("velocity-template-loader")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = internaltypes.BoolTypeOrNil(r.Enabled)
	state.EvaluationOrderIndex = types.Int64Value(r.EvaluationOrderIndex)
	state.MimeTypeMatcher = types.StringValue(r.MimeTypeMatcher)
	state.MimeType = internaltypes.StringTypeOrNil(r.MimeType, internaltypes.IsEmptyString(expectedValues.MimeType))
	state.TemplateSuffix = internaltypes.StringTypeOrNil(r.TemplateSuffix, internaltypes.IsEmptyString(expectedValues.TemplateSuffix))
	state.TemplateDirectory = internaltypes.StringTypeOrNil(r.TemplateDirectory, internaltypes.IsEmptyString(expectedValues.TemplateDirectory))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *velocityTemplateLoaderResourceModel) setStateValuesNotReturnedByAPI(expectedValues *velocityTemplateLoaderResourceModel) {
	if !expectedValues.HttpServletExtensionName.IsUnknown() {
		state.HttpServletExtensionName = expectedValues.HttpServletExtensionName
	}
}

// Create any update operations necessary to make the state match the plan
func createVelocityTemplateLoaderOperations(plan velocityTemplateLoaderResourceModel, state velocityTemplateLoaderResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddInt64OperationIfNecessary(&ops, plan.EvaluationOrderIndex, state.EvaluationOrderIndex, "evaluation-order-index")
	operations.AddStringOperationIfNecessary(&ops, plan.MimeTypeMatcher, state.MimeTypeMatcher, "mime-type-matcher")
	operations.AddStringOperationIfNecessary(&ops, plan.MimeType, state.MimeType, "mime-type")
	operations.AddStringOperationIfNecessary(&ops, plan.TemplateSuffix, state.TemplateSuffix, "template-suffix")
	operations.AddStringOperationIfNecessary(&ops, plan.TemplateDirectory, state.TemplateDirectory, "template-directory")
	return ops
}

// Create a velocity-template-loader velocity-template-loader
func (r *velocityTemplateLoaderResource) CreateVelocityTemplateLoader(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan velocityTemplateLoaderResourceModel) (*velocityTemplateLoaderResourceModel, error) {
	addRequest := client.NewAddVelocityTemplateLoaderRequest(plan.MimeTypeMatcher.ValueString(),
		plan.Name.ValueString())
	addOptionalVelocityTemplateLoaderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VelocityTemplateLoaderAPI.AddVelocityTemplateLoader(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.HttpServletExtensionName.ValueString())
	apiAddRequest = apiAddRequest.AddVelocityTemplateLoaderRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.VelocityTemplateLoaderAPI.AddVelocityTemplateLoaderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Velocity Template Loader", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state velocityTemplateLoaderResourceModel
	readVelocityTemplateLoaderResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *velocityTemplateLoaderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan velocityTemplateLoaderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateVelocityTemplateLoader(ctx, req, resp, plan)
	if err != nil {
		return
	}

	// Populate Computed attribute values
	state.setStateValuesNotReturnedByAPI(&plan)
	// Set state to fully populated data
	diags = resp.State.Set(ctx, *state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultVelocityTemplateLoaderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan velocityTemplateLoaderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.VelocityTemplateLoaderAPI.GetVelocityTemplateLoader(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.HttpServletExtensionName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Velocity Template Loader", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state velocityTemplateLoaderResourceModel
	readVelocityTemplateLoaderResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.VelocityTemplateLoaderAPI.UpdateVelocityTemplateLoader(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.HttpServletExtensionName.ValueString())
	ops := createVelocityTemplateLoaderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.VelocityTemplateLoaderAPI.UpdateVelocityTemplateLoaderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Velocity Template Loader", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readVelocityTemplateLoaderResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *velocityTemplateLoaderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readVelocityTemplateLoader(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultVelocityTemplateLoaderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readVelocityTemplateLoader(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readVelocityTemplateLoader(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state velocityTemplateLoaderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.VelocityTemplateLoaderAPI.GetVelocityTemplateLoader(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString(), state.HttpServletExtensionName.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Velocity Template Loader", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Velocity Template Loader", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readVelocityTemplateLoaderResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *velocityTemplateLoaderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateVelocityTemplateLoader(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultVelocityTemplateLoaderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateVelocityTemplateLoader(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateVelocityTemplateLoader(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan velocityTemplateLoaderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state velocityTemplateLoaderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.VelocityTemplateLoaderAPI.UpdateVelocityTemplateLoader(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString(), plan.HttpServletExtensionName.ValueString())

	// Determine what update operations are necessary
	ops := createVelocityTemplateLoaderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.VelocityTemplateLoaderAPI.UpdateVelocityTemplateLoaderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Velocity Template Loader", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readVelocityTemplateLoaderResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// This config object is edit-only, so Terraform can't delete it.
// After running a delete, Terraform will just "forget" about this object and it can be managed elsewhere.
func (r *defaultVelocityTemplateLoaderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *velocityTemplateLoaderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state velocityTemplateLoaderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.VelocityTemplateLoaderAPI.DeleteVelocityTemplateLoaderExecute(r.apiClient.VelocityTemplateLoaderAPI.DeleteVelocityTemplateLoader(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.HttpServletExtensionName.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Velocity Template Loader", err, httpResp)
		return
	}
}

func (r *velocityTemplateLoaderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importVelocityTemplateLoader(ctx, req, resp)
}

func (r *defaultVelocityTemplateLoaderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importVelocityTemplateLoader(ctx, req, resp)
}

func importVelocityTemplateLoader(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [http-servlet-extension-name]/[velocity-template-loader-name]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("http_servlet_extension_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), split[1])...)
}
