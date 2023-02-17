package httpservletextension

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
	client "github.com/pingidentity/pingdirectory-go-client/v9100"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &velocityHttpServletExtensionResource{}
	_ resource.ResourceWithConfigure   = &velocityHttpServletExtensionResource{}
	_ resource.ResourceWithImportState = &velocityHttpServletExtensionResource{}
)

// Create a Velocity Http Servlet Extension resource
func NewVelocityHttpServletExtensionResource() resource.Resource {
	return &velocityHttpServletExtensionResource{}
}

// velocityHttpServletExtensionResource is the resource implementation.
type velocityHttpServletExtensionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *velocityHttpServletExtensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_velocity_http_servlet_extension"
}

// Configure adds the provider configured client to the resource.
func (r *velocityHttpServletExtensionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type velocityHttpServletExtensionResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	LastUpdated                 types.String `tfsdk:"last_updated"`
	Notifications               types.Set    `tfsdk:"notifications"`
	RequiredActions             types.Set    `tfsdk:"required_actions"`
	BaseContextPath             types.String `tfsdk:"base_context_path"`
	StaticContextPath           types.String `tfsdk:"static_context_path"`
	StaticContentDirectory      types.String `tfsdk:"static_content_directory"`
	StaticCustomDirectory       types.String `tfsdk:"static_custom_directory"`
	TemplateDirectory           types.Set    `tfsdk:"template_directory"`
	ExposeRequestAttributes     types.Bool   `tfsdk:"expose_request_attributes"`
	ExposeSessionAttributes     types.Bool   `tfsdk:"expose_session_attributes"`
	ExposeServerContext         types.Bool   `tfsdk:"expose_server_context"`
	AllowContextOverride        types.Bool   `tfsdk:"allow_context_override"`
	MimeTypesFile               types.String `tfsdk:"mime_types_file"`
	DefaultMIMEType             types.String `tfsdk:"default_mime_type"`
	CharacterEncoding           types.String `tfsdk:"character_encoding"`
	ResponseHeader              types.Set    `tfsdk:"response_header"`
	StaticResponseHeader        types.Set    `tfsdk:"static_response_header"`
	RequireAuthentication       types.Bool   `tfsdk:"require_authentication"`
	IdentityMapper              types.String `tfsdk:"identity_mapper"`
	Description                 types.String `tfsdk:"description"`
	CrossOriginPolicy           types.String `tfsdk:"cross_origin_policy"`
	CorrelationIDResponseHeader types.String `tfsdk:"correlation_id_response_header"`
}

// GetSchema defines the schema for the resource.
func (r *velocityHttpServletExtensionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Velocity Http Servlet Extension.",
		Attributes: map[string]schema.Attribute{
			"base_context_path": schema.StringAttribute{
				Description: "The context path to use to access all template-based and static content. The value must start with a forward slash and must represent a valid HTTP context path.",
				Optional:    true,
				Computed:    true,
			},
			"static_context_path": schema.StringAttribute{
				Description: "The path below the base context path by which static, non-template content such as images, CSS, and Javascript files are accessible.",
				Optional:    true,
				Computed:    true,
			},
			"static_content_directory": schema.StringAttribute{
				Description: "Specifies the base directory in which static, non-template content such as images, CSS, and Javascript files are stored on the filesystem.",
				Optional:    true,
				Computed:    true,
			},
			"static_custom_directory": schema.StringAttribute{
				Description: "Specifies the base directory in which custom static, non-template content such as images, CSS, and Javascript files are stored on the filesystem. Files in this directory will override those with the same name in the directory specified by the static-content-directory property.",
				Optional:    true,
				Computed:    true,
			},
			"template_directory": schema.SetAttribute{
				Description: "Specifies an ordered list of directories in which to search for the template files.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"expose_request_attributes": schema.BoolAttribute{
				Description: "Specifies whether the HTTP request will be exposed to templates.",
				Optional:    true,
				Computed:    true,
			},
			"expose_session_attributes": schema.BoolAttribute{
				Description: "Specifies whether the HTTP session will be exposed to templates.",
				Optional:    true,
				Computed:    true,
			},
			"expose_server_context": schema.BoolAttribute{
				Description: "Specifies whether a server context will be exposed under context key 'ubid_server' for all template contexts.",
				Optional:    true,
				Computed:    true,
			},
			"allow_context_override": schema.BoolAttribute{
				Description: "Indicates whether context providers may override existing context objects with new values.",
				Optional:    true,
				Computed:    true,
			},
			"mime_types_file": schema.StringAttribute{
				Description: "Specifies the path to a file that contains MIME type mappings that will be used to determine the appropriate value to return for the Content-Type header based on the extension of the requested static content file.",
				Optional:    true,
				Computed:    true,
			},
			"default_mime_type": schema.StringAttribute{
				Description: "Specifies the default value that will be used in the response's Content-Type header that indicates the type of content to return.",
				Optional:    true,
				Computed:    true,
			},
			"character_encoding": schema.StringAttribute{
				Description: "Specifies the value that will be used for all responses' Content-Type headers' charset parameter that indicates the character encoding of the document.",
				Optional:    true,
				Computed:    true,
			},
			"response_header": schema.SetAttribute{
				Description: "Specifies HTTP header fields and values added to response headers for all template page requests.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"static_response_header": schema.SetAttribute{
				Description: "Specifies HTTP header fields and values added to response headers for static content requests such as images and scripts.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"require_authentication": schema.BoolAttribute{
				Description: "Require authentication when accessing Velocity templates.",
				Optional:    true,
				Computed:    true,
			},
			"identity_mapper": schema.StringAttribute{
				Description: "Specifies the name of the identity mapper that is to be used for associating basic authentication credentials with user entries.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this HTTP Servlet Extension",
				Optional:    true,
				Computed:    true,
			},
			"cross_origin_policy": schema.StringAttribute{
				Description: "The cross-origin request policy to use for the HTTP Servlet Extension.",
				Optional:    true,
				Computed:    true,
			},
			"correlation_id_response_header": schema.StringAttribute{
				Description: "Specifies the name of the HTTP response header that will contain a correlation ID value. Example values are \"Correlation-Id\", \"X-Amzn-Trace-Id\", and \"X-Request-Id\".",
				Optional:    true,
				Computed:    true,
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Read a VelocityHttpServletExtensionResponse object into the model struct
func readVelocityHttpServletExtensionResponse(ctx context.Context, r *client.VelocityHttpServletExtensionResponse, state *velocityHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.StaticContextPath = internaltypes.StringTypeOrNil(r.StaticContextPath, true)
	state.StaticContentDirectory = internaltypes.StringTypeOrNil(r.StaticContentDirectory, true)
	state.StaticCustomDirectory = internaltypes.StringTypeOrNil(r.StaticCustomDirectory, true)
	state.TemplateDirectory = internaltypes.GetStringSet(r.TemplateDirectory)
	state.ExposeRequestAttributes = internaltypes.BoolTypeOrNil(r.ExposeRequestAttributes)
	state.ExposeSessionAttributes = internaltypes.BoolTypeOrNil(r.ExposeSessionAttributes)
	state.ExposeServerContext = internaltypes.BoolTypeOrNil(r.ExposeServerContext)
	state.AllowContextOverride = internaltypes.BoolTypeOrNil(r.AllowContextOverride)
	state.MimeTypesFile = internaltypes.StringTypeOrNil(r.MimeTypesFile, true)
	state.DefaultMIMEType = internaltypes.StringTypeOrNil(r.DefaultMIMEType, true)
	state.CharacterEncoding = internaltypes.StringTypeOrNil(r.CharacterEncoding, true)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.StaticResponseHeader = internaltypes.GetStringSet(r.StaticResponseHeader)
	state.RequireAuthentication = internaltypes.BoolTypeOrNil(r.RequireAuthentication)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, true)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createVelocityHttpServletExtensionOperations(plan velocityHttpServletExtensionResourceModel, state velocityHttpServletExtensionResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.BaseContextPath, state.BaseContextPath, "base-context-path")
	operations.AddStringOperationIfNecessary(&ops, plan.StaticContextPath, state.StaticContextPath, "static-context-path")
	operations.AddStringOperationIfNecessary(&ops, plan.StaticContentDirectory, state.StaticContentDirectory, "static-content-directory")
	operations.AddStringOperationIfNecessary(&ops, plan.StaticCustomDirectory, state.StaticCustomDirectory, "static-custom-directory")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.TemplateDirectory, state.TemplateDirectory, "template-directory")
	operations.AddBoolOperationIfNecessary(&ops, plan.ExposeRequestAttributes, state.ExposeRequestAttributes, "expose-request-attributes")
	operations.AddBoolOperationIfNecessary(&ops, plan.ExposeSessionAttributes, state.ExposeSessionAttributes, "expose-session-attributes")
	operations.AddBoolOperationIfNecessary(&ops, plan.ExposeServerContext, state.ExposeServerContext, "expose-server-context")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowContextOverride, state.AllowContextOverride, "allow-context-override")
	operations.AddStringOperationIfNecessary(&ops, plan.MimeTypesFile, state.MimeTypesFile, "mime-types-file")
	operations.AddStringOperationIfNecessary(&ops, plan.DefaultMIMEType, state.DefaultMIMEType, "default-mime-type")
	operations.AddStringOperationIfNecessary(&ops, plan.CharacterEncoding, state.CharacterEncoding, "character-encoding")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ResponseHeader, state.ResponseHeader, "response-header")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.StaticResponseHeader, state.StaticResponseHeader, "static-response-header")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequireAuthentication, state.RequireAuthentication, "require-authentication")
	operations.AddStringOperationIfNecessary(&ops, plan.IdentityMapper, state.IdentityMapper, "identity-mapper")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.CrossOriginPolicy, state.CrossOriginPolicy, "cross-origin-policy")
	operations.AddStringOperationIfNecessary(&ops, plan.CorrelationIDResponseHeader, state.CorrelationIDResponseHeader, "correlation-id-response-header")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *velocityHttpServletExtensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan velocityHttpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.GetHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Velocity Http Servlet Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state velocityHttpServletExtensionResourceModel
	readVelocityHttpServletExtensionResponse(ctx, readResponse.VelocityHttpServletExtensionResponse, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtension(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createVelocityHttpServletExtensionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Velocity Http Servlet Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readVelocityHttpServletExtensionResponse(ctx, updateResponse.VelocityHttpServletExtensionResponse, &state, &resp.Diagnostics)
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
func (r *velocityHttpServletExtensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state velocityHttpServletExtensionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.GetHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Velocity Http Servlet Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readVelocityHttpServletExtensionResponse(ctx, readResponse.VelocityHttpServletExtensionResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *velocityHttpServletExtensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan velocityHttpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state velocityHttpServletExtensionResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createVelocityHttpServletExtensionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Velocity Http Servlet Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readVelocityHttpServletExtensionResponse(ctx, updateResponse.VelocityHttpServletExtensionResponse, &state, &resp.Diagnostics)
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
func (r *velocityHttpServletExtensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *velocityHttpServletExtensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
