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
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &fileServerHttpServletExtensionResource{}
	_ resource.ResourceWithConfigure   = &fileServerHttpServletExtensionResource{}
	_ resource.ResourceWithImportState = &fileServerHttpServletExtensionResource{}
	_ resource.Resource                = &defaultFileServerHttpServletExtensionResource{}
	_ resource.ResourceWithConfigure   = &defaultFileServerHttpServletExtensionResource{}
	_ resource.ResourceWithImportState = &defaultFileServerHttpServletExtensionResource{}
)

// Create a File Server Http Servlet Extension resource
func NewFileServerHttpServletExtensionResource() resource.Resource {
	return &fileServerHttpServletExtensionResource{}
}

func NewDefaultFileServerHttpServletExtensionResource() resource.Resource {
	return &defaultFileServerHttpServletExtensionResource{}
}

// fileServerHttpServletExtensionResource is the resource implementation.
type fileServerHttpServletExtensionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultFileServerHttpServletExtensionResource is the resource implementation.
type defaultFileServerHttpServletExtensionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *fileServerHttpServletExtensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file_server_http_servlet_extension"
}

func (r *defaultFileServerHttpServletExtensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_file_server_http_servlet_extension"
}

// Configure adds the provider configured client to the resource.
func (r *fileServerHttpServletExtensionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultFileServerHttpServletExtensionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type fileServerHttpServletExtensionResourceModel struct {
	Id                                types.String `tfsdk:"id"`
	LastUpdated                       types.String `tfsdk:"last_updated"`
	Notifications                     types.Set    `tfsdk:"notifications"`
	RequiredActions                   types.Set    `tfsdk:"required_actions"`
	BaseContextPath                   types.String `tfsdk:"base_context_path"`
	DocumentRootDirectory             types.String `tfsdk:"document_root_directory"`
	EnableDirectoryIndexing           types.Bool   `tfsdk:"enable_directory_indexing"`
	IndexFile                         types.Set    `tfsdk:"index_file"`
	MimeTypesFile                     types.String `tfsdk:"mime_types_file"`
	DefaultMIMEType                   types.String `tfsdk:"default_mime_type"`
	RequireAuthentication             types.Bool   `tfsdk:"require_authentication"`
	AllowedAuthenticationType         types.Set    `tfsdk:"allowed_authentication_type"`
	AccessTokenValidator              types.Set    `tfsdk:"access_token_validator"`
	IdTokenValidator                  types.Set    `tfsdk:"id_token_validator"`
	RequireFileServletAccessPrivilege types.Bool   `tfsdk:"require_file_servlet_access_privilege"`
	RequireGroup                      types.Set    `tfsdk:"require_group"`
	IdentityMapper                    types.String `tfsdk:"identity_mapper"`
	Description                       types.String `tfsdk:"description"`
	CrossOriginPolicy                 types.String `tfsdk:"cross_origin_policy"`
	ResponseHeader                    types.Set    `tfsdk:"response_header"`
	CorrelationIDResponseHeader       types.String `tfsdk:"correlation_id_response_header"`
}

// GetSchema defines the schema for the resource.
func (r *fileServerHttpServletExtensionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	fileServerHttpServletExtensionSchema(ctx, req, resp, false)
}

func (r *defaultFileServerHttpServletExtensionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	fileServerHttpServletExtensionSchema(ctx, req, resp, true)
}

func fileServerHttpServletExtensionSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a File Server Http Servlet Extension.",
		Attributes: map[string]schema.Attribute{
			"base_context_path": schema.StringAttribute{
				Description: "Specifies the base context path that should be used by HTTP clients to reference content. The value must start with a forward slash and must represent a valid HTTP context path.",
				Required:    true,
			},
			"document_root_directory": schema.StringAttribute{
				Description: "Specifies the path to the directory on the local filesystem containing the files to be served by this File Server HTTP Servlet Extension. The path must exist, and it must be a directory.",
				Required:    true,
			},
			"enable_directory_indexing": schema.BoolAttribute{
				Description: "Indicates whether to generate a default HTML page with a listing of available files if the requested path refers to a directory rather than a file, and that directory does not contain an index file.",
				Optional:    true,
				Computed:    true,
			},
			"index_file": schema.SetAttribute{
				Description: "Specifies the name of a file whose contents may be returned to the client if the requested path refers to a directory rather than a file.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"mime_types_file": schema.StringAttribute{
				Description: "Specifies the path to a file that contains MIME type mappings that will be used to determine the appropriate value to return for the Content-Type header based on the extension of the requested file.",
				Optional:    true,
				Computed:    true,
			},
			"default_mime_type": schema.StringAttribute{
				Description: "Specifies the default MIME type to use for the Content-Type header when a mapping cannot be found.",
				Optional:    true,
				Computed:    true,
			},
			"require_authentication": schema.BoolAttribute{
				Description: "Indicates whether the servlet extension should only accept requests from authenticated clients.",
				Optional:    true,
				Computed:    true,
			},
			"allowed_authentication_type": schema.SetAttribute{
				Description: "The types of authentication that may be used to authenticate to the file servlet.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"access_token_validator": schema.SetAttribute{
				Description: "The access token validators that may be used to verify the authenticity of an OAuth 2.0 bearer token.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"id_token_validator": schema.SetAttribute{
				Description: "The ID token validators that may be used to verify the authenticity of an of an OpenID Connect ID token.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"require_file_servlet_access_privilege": schema.BoolAttribute{
				Description: "Indicates whether the servlet extension should only accept requests from authenticated clients that have the file-servlet-access privilege.",
				Optional:    true,
				Computed:    true,
			},
			"require_group": schema.SetAttribute{
				Description: "The DN of a group whose members will be permitted to access to the associated files. If multiple group DNs are configured, then anyone who is a member of at least one of those groups will be granted access.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"identity_mapper": schema.StringAttribute{
				Description: "The identity mapper that will be used to identify the entry with which a username is associated.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this HTTP Servlet Extension",
				Optional:    true,
			},
			"cross_origin_policy": schema.StringAttribute{
				Description: "The cross-origin request policy to use for the HTTP Servlet Extension.",
				Optional:    true,
				Computed:    true,
			},
			"response_header": schema.SetAttribute{
				Description: "Specifies HTTP header fields and values added to response headers for all requests.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"correlation_id_response_header": schema.StringAttribute{
				Description: "Specifies the name of the HTTP response header that will contain a correlation ID value. Example values are \"Correlation-Id\", \"X-Amzn-Trace-Id\", and \"X-Request-Id\".",
				Optional:    true,
				Computed:    true,
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalFileServerHttpServletExtensionFields(ctx context.Context, addRequest *client.AddFileServerHttpServletExtensionRequest, plan fileServerHttpServletExtensionResourceModel) error {
	if internaltypes.IsDefined(plan.EnableDirectoryIndexing) {
		boolVal := plan.EnableDirectoryIndexing.ValueBool()
		addRequest.EnableDirectoryIndexing = &boolVal
	}
	if internaltypes.IsDefined(plan.IndexFile) {
		var slice []string
		plan.IndexFile.ElementsAs(ctx, &slice, false)
		addRequest.IndexFile = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MimeTypesFile) {
		stringVal := plan.MimeTypesFile.ValueString()
		addRequest.MimeTypesFile = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DefaultMIMEType) {
		stringVal := plan.DefaultMIMEType.ValueString()
		addRequest.DefaultMIMEType = &stringVal
	}
	if internaltypes.IsDefined(plan.RequireAuthentication) {
		boolVal := plan.RequireAuthentication.ValueBool()
		addRequest.RequireAuthentication = &boolVal
	}
	if internaltypes.IsDefined(plan.AllowedAuthenticationType) {
		var slice []string
		plan.AllowedAuthenticationType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumhttpServletExtensionAllowedAuthenticationTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumhttpServletExtensionAllowedAuthenticationTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.AllowedAuthenticationType = enumSlice
	}
	if internaltypes.IsDefined(plan.AccessTokenValidator) {
		var slice []string
		plan.AccessTokenValidator.ElementsAs(ctx, &slice, false)
		addRequest.AccessTokenValidator = slice
	}
	if internaltypes.IsDefined(plan.IdTokenValidator) {
		var slice []string
		plan.IdTokenValidator.ElementsAs(ctx, &slice, false)
		addRequest.IdTokenValidator = slice
	}
	if internaltypes.IsDefined(plan.RequireFileServletAccessPrivilege) {
		boolVal := plan.RequireFileServletAccessPrivilege.ValueBool()
		addRequest.RequireFileServletAccessPrivilege = &boolVal
	}
	if internaltypes.IsDefined(plan.RequireGroup) {
		var slice []string
		plan.RequireGroup.ElementsAs(ctx, &slice, false)
		addRequest.RequireGroup = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IdentityMapper) {
		stringVal := plan.IdentityMapper.ValueString()
		addRequest.IdentityMapper = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CrossOriginPolicy) {
		stringVal := plan.CrossOriginPolicy.ValueString()
		addRequest.CrossOriginPolicy = &stringVal
	}
	if internaltypes.IsDefined(plan.ResponseHeader) {
		var slice []string
		plan.ResponseHeader.ElementsAs(ctx, &slice, false)
		addRequest.ResponseHeader = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CorrelationIDResponseHeader) {
		stringVal := plan.CorrelationIDResponseHeader.ValueString()
		addRequest.CorrelationIDResponseHeader = &stringVal
	}
	return nil
}

// Read a FileServerHttpServletExtensionResponse object into the model struct
func readFileServerHttpServletExtensionResponse(ctx context.Context, r *client.FileServerHttpServletExtensionResponse, state *fileServerHttpServletExtensionResourceModel, expectedValues *fileServerHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.DocumentRootDirectory = types.StringValue(r.DocumentRootDirectory)
	state.EnableDirectoryIndexing = internaltypes.BoolTypeOrNil(r.EnableDirectoryIndexing)
	state.IndexFile = internaltypes.GetStringSet(r.IndexFile)
	state.MimeTypesFile = internaltypes.StringTypeOrNil(r.MimeTypesFile, internaltypes.IsEmptyString(expectedValues.MimeTypesFile))
	state.DefaultMIMEType = internaltypes.StringTypeOrNil(r.DefaultMIMEType, internaltypes.IsEmptyString(expectedValues.DefaultMIMEType))
	state.RequireAuthentication = internaltypes.BoolTypeOrNil(r.RequireAuthentication)
	state.AllowedAuthenticationType = internaltypes.GetStringSet(
		client.StringSliceEnumhttpServletExtensionAllowedAuthenticationTypeProp(r.AllowedAuthenticationType))
	state.AccessTokenValidator = internaltypes.GetStringSet(r.AccessTokenValidator)
	state.IdTokenValidator = internaltypes.GetStringSet(r.IdTokenValidator)
	state.RequireFileServletAccessPrivilege = internaltypes.BoolTypeOrNil(r.RequireFileServletAccessPrivilege)
	state.RequireGroup = internaltypes.GetStringSet(r.RequireGroup)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, internaltypes.IsEmptyString(expectedValues.IdentityMapper))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, internaltypes.IsEmptyString(expectedValues.CrossOriginPolicy))
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, internaltypes.IsEmptyString(expectedValues.CorrelationIDResponseHeader))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createFileServerHttpServletExtensionOperations(plan fileServerHttpServletExtensionResourceModel, state fileServerHttpServletExtensionResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.BaseContextPath, state.BaseContextPath, "base-context-path")
	operations.AddStringOperationIfNecessary(&ops, plan.DocumentRootDirectory, state.DocumentRootDirectory, "document-root-directory")
	operations.AddBoolOperationIfNecessary(&ops, plan.EnableDirectoryIndexing, state.EnableDirectoryIndexing, "enable-directory-indexing")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IndexFile, state.IndexFile, "index-file")
	operations.AddStringOperationIfNecessary(&ops, plan.MimeTypesFile, state.MimeTypesFile, "mime-types-file")
	operations.AddStringOperationIfNecessary(&ops, plan.DefaultMIMEType, state.DefaultMIMEType, "default-mime-type")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequireAuthentication, state.RequireAuthentication, "require-authentication")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedAuthenticationType, state.AllowedAuthenticationType, "allowed-authentication-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AccessTokenValidator, state.AccessTokenValidator, "access-token-validator")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IdTokenValidator, state.IdTokenValidator, "id-token-validator")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequireFileServletAccessPrivilege, state.RequireFileServletAccessPrivilege, "require-file-servlet-access-privilege")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RequireGroup, state.RequireGroup, "require-group")
	operations.AddStringOperationIfNecessary(&ops, plan.IdentityMapper, state.IdentityMapper, "identity-mapper")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.CrossOriginPolicy, state.CrossOriginPolicy, "cross-origin-policy")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ResponseHeader, state.ResponseHeader, "response-header")
	operations.AddStringOperationIfNecessary(&ops, plan.CorrelationIDResponseHeader, state.CorrelationIDResponseHeader, "correlation-id-response-header")
	return ops
}

// Create a new resource
func (r *fileServerHttpServletExtensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan fileServerHttpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddFileServerHttpServletExtensionRequest(plan.Id.ValueString(),
		[]client.EnumfileServerHttpServletExtensionSchemaUrn{client.ENUMFILESERVERHTTPSERVLETEXTENSIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0HTTP_SERVLET_EXTENSIONFILE_SERVER},
		plan.BaseContextPath.ValueString(),
		plan.DocumentRootDirectory.ValueString())
	err := addOptionalFileServerHttpServletExtensionFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for File Server Http Servlet Extension", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.HttpServletExtensionApi.AddHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddHttpServletExtensionRequest(
		client.AddFileServerHttpServletExtensionRequestAsAddHttpServletExtensionRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.AddHttpServletExtensionExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the File Server Http Servlet Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state fileServerHttpServletExtensionResourceModel
	readFileServerHttpServletExtensionResponse(ctx, addResponse.FileServerHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultFileServerHttpServletExtensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan fileServerHttpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.GetHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the File Server Http Servlet Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state fileServerHttpServletExtensionResourceModel
	readFileServerHttpServletExtensionResponse(ctx, readResponse.FileServerHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtension(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createFileServerHttpServletExtensionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the File Server Http Servlet Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readFileServerHttpServletExtensionResponse(ctx, updateResponse.FileServerHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
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
func (r *fileServerHttpServletExtensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readFileServerHttpServletExtension(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultFileServerHttpServletExtensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readFileServerHttpServletExtension(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readFileServerHttpServletExtension(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state fileServerHttpServletExtensionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.HttpServletExtensionApi.GetHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the File Server Http Servlet Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readFileServerHttpServletExtensionResponse(ctx, readResponse.FileServerHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *fileServerHttpServletExtensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateFileServerHttpServletExtension(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultFileServerHttpServletExtensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateFileServerHttpServletExtension(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateFileServerHttpServletExtension(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan fileServerHttpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state fileServerHttpServletExtensionResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.HttpServletExtensionApi.UpdateHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createFileServerHttpServletExtensionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.HttpServletExtensionApi.UpdateHttpServletExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the File Server Http Servlet Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readFileServerHttpServletExtensionResponse(ctx, updateResponse.FileServerHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultFileServerHttpServletExtensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *fileServerHttpServletExtensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state fileServerHttpServletExtensionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.HttpServletExtensionApi.DeleteHttpServletExtensionExecute(r.apiClient.HttpServletExtensionApi.DeleteHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the File Server Http Servlet Extension", err, httpResp)
		return
	}
}

func (r *fileServerHttpServletExtensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importFileServerHttpServletExtension(ctx, req, resp)
}

func (r *defaultFileServerHttpServletExtensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importFileServerHttpServletExtension(ctx, req, resp)
}

func importFileServerHttpServletExtension(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
