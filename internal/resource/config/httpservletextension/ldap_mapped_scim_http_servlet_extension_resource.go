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
	_ resource.Resource                = &ldapMappedScimHttpServletExtensionResource{}
	_ resource.ResourceWithConfigure   = &ldapMappedScimHttpServletExtensionResource{}
	_ resource.ResourceWithImportState = &ldapMappedScimHttpServletExtensionResource{}
)

// Create a Ldap Mapped Scim Http Servlet Extension resource
func NewLdapMappedScimHttpServletExtensionResource() resource.Resource {
	return &ldapMappedScimHttpServletExtensionResource{}
}

// ldapMappedScimHttpServletExtensionResource is the resource implementation.
type ldapMappedScimHttpServletExtensionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *ldapMappedScimHttpServletExtensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ldap_mapped_scim_http_servlet_extension"
}

// Configure adds the provider configured client to the resource.
func (r *ldapMappedScimHttpServletExtensionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type ldapMappedScimHttpServletExtensionResourceModel struct {
	Id                            types.String `tfsdk:"id"`
	LastUpdated                   types.String `tfsdk:"last_updated"`
	Notifications                 types.Set    `tfsdk:"notifications"`
	RequiredActions               types.Set    `tfsdk:"required_actions"`
	OAuthTokenHandler             types.String `tfsdk:"oauth_token_handler"`
	BasicAuthEnabled              types.Bool   `tfsdk:"basic_auth_enabled"`
	IdentityMapper                types.String `tfsdk:"identity_mapper"`
	ResourceMappingFile           types.String `tfsdk:"resource_mapping_file"`
	IncludeLDAPObjectclass        types.Set    `tfsdk:"include_ldap_objectclass"`
	ExcludeLDAPObjectclass        types.Set    `tfsdk:"exclude_ldap_objectclass"`
	IncludeLDAPBaseDN             types.Set    `tfsdk:"include_ldap_base_dn"`
	ExcludeLDAPBaseDN             types.Set    `tfsdk:"exclude_ldap_base_dn"`
	EntityTagLDAPAttribute        types.String `tfsdk:"entity_tag_ldap_attribute"`
	BaseContextPath               types.String `tfsdk:"base_context_path"`
	TemporaryDirectory            types.String `tfsdk:"temporary_directory"`
	TemporaryDirectoryPermissions types.String `tfsdk:"temporary_directory_permissions"`
	MaxResults                    types.Int64  `tfsdk:"max_results"`
	BulkMaxOperations             types.Int64  `tfsdk:"bulk_max_operations"`
	BulkMaxPayloadSize            types.String `tfsdk:"bulk_max_payload_size"`
	BulkMaxConcurrentRequests     types.Int64  `tfsdk:"bulk_max_concurrent_requests"`
	DebugEnabled                  types.Bool   `tfsdk:"debug_enabled"`
	DebugLevel                    types.String `tfsdk:"debug_level"`
	DebugType                     types.Set    `tfsdk:"debug_type"`
	IncludeStackTrace             types.Bool   `tfsdk:"include_stack_trace"`
	Description                   types.String `tfsdk:"description"`
	CrossOriginPolicy             types.String `tfsdk:"cross_origin_policy"`
	ResponseHeader                types.Set    `tfsdk:"response_header"`
	CorrelationIDResponseHeader   types.String `tfsdk:"correlation_id_response_header"`
}

// GetSchema defines the schema for the resource.
func (r *ldapMappedScimHttpServletExtensionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Ldap Mapped Scim Http Servlet Extension.",
		Attributes: map[string]schema.Attribute{
			"oauth_token_handler": schema.StringAttribute{
				Description: "Specifies the OAuth Token Handler implementation that should be used to validate OAuth 2.0 bearer tokens when they are included in a SCIM request.",
				Optional:    true,
			},
			"basic_auth_enabled": schema.BoolAttribute{
				Description: "Enables HTTP Basic authentication, using a username and password.",
				Optional:    true,
				Computed:    true,
			},
			"identity_mapper": schema.StringAttribute{
				Description: "Specifies the name of the identity mapper that is to be used to match the username included in the HTTP Basic authentication header to the corresponding user in the directory.",
				Optional:    true,
			},
			"resource_mapping_file": schema.StringAttribute{
				Description: "The path to an XML file defining the resources supported by the SCIM interface and the SCIM-to-LDAP attribute mappings to use.",
				Optional:    true,
				Computed:    true,
			},
			"include_ldap_objectclass": schema.SetAttribute{
				Description: "Specifies the LDAP object classes that should be exposed directly as SCIM resources.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"exclude_ldap_objectclass": schema.SetAttribute{
				Description: "Specifies the LDAP object classes that should be not be exposed directly as SCIM resources.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"include_ldap_base_dn": schema.SetAttribute{
				Description: "Specifies the base DNs for the branches of the DIT that should be exposed via the Identity Access API.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"exclude_ldap_base_dn": schema.SetAttribute{
				Description: "Specifies the base DNs for the branches of the DIT that should not be exposed via the Identity Access API.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"entity_tag_ldap_attribute": schema.StringAttribute{
				Description: "Specifies the LDAP attribute whose value should be used as the entity tag value to enable SCIM resource versioning support.",
				Optional:    true,
			},
			"base_context_path": schema.StringAttribute{
				Description: "The context path to use to access the SCIM interface. The value must start with a forward slash and must represent a valid HTTP context path.",
				Required:    true,
			},
			"temporary_directory": schema.StringAttribute{
				Description: "Specifies the location of the directory that is used to create temporary files containing SCIM request data.",
				Required:    true,
			},
			"temporary_directory_permissions": schema.StringAttribute{
				Description: "Specifies the permissions that should be applied to the directory that is used to create temporary files.",
				Required:    true,
			},
			"max_results": schema.Int64Attribute{
				Description: "The maximum number of resources that are returned in a response.",
				Optional:    true,
				Computed:    true,
			},
			"bulk_max_operations": schema.Int64Attribute{
				Description: "The maximum number of operations that are permitted in a bulk request.",
				Optional:    true,
				Computed:    true,
			},
			"bulk_max_payload_size": schema.StringAttribute{
				Description: "The maximum payload size in bytes of a bulk request.",
				Optional:    true,
				Computed:    true,
			},
			"bulk_max_concurrent_requests": schema.Int64Attribute{
				Description: "The maximum number of bulk requests that may be processed concurrently by the server. Any bulk request that would cause this limit to be exceeded is rejected with HTTP status code 503.",
				Optional:    true,
				Computed:    true,
			},
			"debug_enabled": schema.BoolAttribute{
				Description: "Enables debug logging of the SCIM SDK. Debug messages will be forwarded to the Directory Server debug logger with the scope of com.unboundid.directory.server.extensions.scim.SCIMHTTPServletExtension.",
				Optional:    true,
				Computed:    true,
			},
			"debug_level": schema.StringAttribute{
				Description: "The minimum debug level that should be used for messages to be logged.",
				Required:    true,
			},
			"debug_type": schema.SetAttribute{
				Description: "The types of debug messages that should be logged.",
				Required:    true,
				ElementType: types.StringType,
			},
			"include_stack_trace": schema.BoolAttribute{
				Description: "Indicates whether a stack trace of the thread which called the debug method should be included in debug log messages.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this HTTP Servlet Extension",
				Optional:    true,
			},
			"cross_origin_policy": schema.StringAttribute{
				Description: "The cross-origin request policy to use for the HTTP Servlet Extension.",
				Optional:    true,
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
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalLdapMappedScimHttpServletExtensionFields(ctx context.Context, addRequest *client.AddLdapMappedScimHttpServletExtensionRequest, plan ldapMappedScimHttpServletExtensionResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OAuthTokenHandler) {
		stringVal := plan.OAuthTokenHandler.ValueString()
		addRequest.OAuthTokenHandler = &stringVal
	}
	if internaltypes.IsDefined(plan.BasicAuthEnabled) {
		boolVal := plan.BasicAuthEnabled.ValueBool()
		addRequest.BasicAuthEnabled = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IdentityMapper) {
		stringVal := plan.IdentityMapper.ValueString()
		addRequest.IdentityMapper = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ResourceMappingFile) {
		stringVal := plan.ResourceMappingFile.ValueString()
		addRequest.ResourceMappingFile = &stringVal
	}
	if internaltypes.IsDefined(plan.IncludeLDAPObjectclass) {
		var slice []string
		plan.IncludeLDAPObjectclass.ElementsAs(ctx, &slice, false)
		addRequest.IncludeLDAPObjectclass = slice
	}
	if internaltypes.IsDefined(plan.ExcludeLDAPObjectclass) {
		var slice []string
		plan.ExcludeLDAPObjectclass.ElementsAs(ctx, &slice, false)
		addRequest.ExcludeLDAPObjectclass = slice
	}
	if internaltypes.IsDefined(plan.IncludeLDAPBaseDN) {
		var slice []string
		plan.IncludeLDAPBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.IncludeLDAPBaseDN = slice
	}
	if internaltypes.IsDefined(plan.ExcludeLDAPBaseDN) {
		var slice []string
		plan.ExcludeLDAPBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.ExcludeLDAPBaseDN = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EntityTagLDAPAttribute) {
		stringVal := plan.EntityTagLDAPAttribute.ValueString()
		addRequest.EntityTagLDAPAttribute = &stringVal
	}
	if internaltypes.IsDefined(plan.MaxResults) {
		intVal := int32(plan.MaxResults.ValueInt64())
		addRequest.MaxResults = &intVal
	}
	if internaltypes.IsDefined(plan.BulkMaxOperations) {
		intVal := int32(plan.BulkMaxOperations.ValueInt64())
		addRequest.BulkMaxOperations = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BulkMaxPayloadSize) {
		stringVal := plan.BulkMaxPayloadSize.ValueString()
		addRequest.BulkMaxPayloadSize = &stringVal
	}
	if internaltypes.IsDefined(plan.BulkMaxConcurrentRequests) {
		intVal := int32(plan.BulkMaxConcurrentRequests.ValueInt64())
		addRequest.BulkMaxConcurrentRequests = &intVal
	}
	if internaltypes.IsDefined(plan.DebugEnabled) {
		boolVal := plan.DebugEnabled.ValueBool()
		addRequest.DebugEnabled = &boolVal
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
}

// Read a LdapMappedScimHttpServletExtensionResponse object into the model struct
func readLdapMappedScimHttpServletExtensionResponse(ctx context.Context, r *client.LdapMappedScimHttpServletExtensionResponse, state *ldapMappedScimHttpServletExtensionResourceModel, expectedValues *ldapMappedScimHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.OAuthTokenHandler = internaltypes.StringTypeOrNil(r.OAuthTokenHandler, internaltypes.IsEmptyString(expectedValues.OAuthTokenHandler))
	state.BasicAuthEnabled = internaltypes.BoolTypeOrNil(r.BasicAuthEnabled)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, internaltypes.IsEmptyString(expectedValues.IdentityMapper))
	state.ResourceMappingFile = internaltypes.StringTypeOrNil(r.ResourceMappingFile, internaltypes.IsEmptyString(expectedValues.ResourceMappingFile))
	state.IncludeLDAPObjectclass = internaltypes.GetStringSet(r.IncludeLDAPObjectclass)
	state.ExcludeLDAPObjectclass = internaltypes.GetStringSet(r.ExcludeLDAPObjectclass)
	state.IncludeLDAPBaseDN = internaltypes.GetStringSet(r.IncludeLDAPBaseDN)
	state.ExcludeLDAPBaseDN = internaltypes.GetStringSet(r.ExcludeLDAPBaseDN)
	state.EntityTagLDAPAttribute = internaltypes.StringTypeOrNil(r.EntityTagLDAPAttribute, internaltypes.IsEmptyString(expectedValues.EntityTagLDAPAttribute))
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.TemporaryDirectory = types.StringValue(r.TemporaryDirectory)
	state.TemporaryDirectoryPermissions = types.StringValue(r.TemporaryDirectoryPermissions)
	state.MaxResults = internaltypes.Int64TypeOrNil(r.MaxResults)
	state.BulkMaxOperations = internaltypes.Int64TypeOrNil(r.BulkMaxOperations)
	state.BulkMaxPayloadSize = internaltypes.StringTypeOrNil(r.BulkMaxPayloadSize, internaltypes.IsEmptyString(expectedValues.BulkMaxPayloadSize))
	config.CheckMismatchedPDFormattedAttributes("bulk_max_payload_size",
		expectedValues.BulkMaxPayloadSize, state.BulkMaxPayloadSize, diagnostics)
	state.BulkMaxConcurrentRequests = internaltypes.Int64TypeOrNil(r.BulkMaxConcurrentRequests)
	state.DebugEnabled = internaltypes.BoolTypeOrNil(r.DebugEnabled)
	state.DebugLevel = types.StringValue(r.DebugLevel.String())
	state.DebugType = internaltypes.GetStringSet(
		client.StringSliceEnumhttpServletExtensionDebugTypeProp(r.DebugType))
	state.IncludeStackTrace = types.BoolValue(r.IncludeStackTrace)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, internaltypes.IsEmptyString(expectedValues.CrossOriginPolicy))
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, internaltypes.IsEmptyString(expectedValues.CorrelationIDResponseHeader))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createLdapMappedScimHttpServletExtensionOperations(plan ldapMappedScimHttpServletExtensionResourceModel, state ldapMappedScimHttpServletExtensionResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.OAuthTokenHandler, state.OAuthTokenHandler, "oauth-token-handler")
	operations.AddBoolOperationIfNecessary(&ops, plan.BasicAuthEnabled, state.BasicAuthEnabled, "basic-auth-enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.IdentityMapper, state.IdentityMapper, "identity-mapper")
	operations.AddStringOperationIfNecessary(&ops, plan.ResourceMappingFile, state.ResourceMappingFile, "resource-mapping-file")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeLDAPObjectclass, state.IncludeLDAPObjectclass, "include-ldap-objectclass")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeLDAPObjectclass, state.ExcludeLDAPObjectclass, "exclude-ldap-objectclass")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeLDAPBaseDN, state.IncludeLDAPBaseDN, "include-ldap-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeLDAPBaseDN, state.ExcludeLDAPBaseDN, "exclude-ldap-base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.EntityTagLDAPAttribute, state.EntityTagLDAPAttribute, "entity-tag-ldap-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.BaseContextPath, state.BaseContextPath, "base-context-path")
	operations.AddStringOperationIfNecessary(&ops, plan.TemporaryDirectory, state.TemporaryDirectory, "temporary-directory")
	operations.AddStringOperationIfNecessary(&ops, plan.TemporaryDirectoryPermissions, state.TemporaryDirectoryPermissions, "temporary-directory-permissions")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxResults, state.MaxResults, "max-results")
	operations.AddInt64OperationIfNecessary(&ops, plan.BulkMaxOperations, state.BulkMaxOperations, "bulk-max-operations")
	operations.AddStringOperationIfNecessary(&ops, plan.BulkMaxPayloadSize, state.BulkMaxPayloadSize, "bulk-max-payload-size")
	operations.AddInt64OperationIfNecessary(&ops, plan.BulkMaxConcurrentRequests, state.BulkMaxConcurrentRequests, "bulk-max-concurrent-requests")
	operations.AddBoolOperationIfNecessary(&ops, plan.DebugEnabled, state.DebugEnabled, "debug-enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.DebugLevel, state.DebugLevel, "debug-level")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DebugType, state.DebugType, "debug-type")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeStackTrace, state.IncludeStackTrace, "include-stack-trace")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.CrossOriginPolicy, state.CrossOriginPolicy, "cross-origin-policy")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ResponseHeader, state.ResponseHeader, "response-header")
	operations.AddStringOperationIfNecessary(&ops, plan.CorrelationIDResponseHeader, state.CorrelationIDResponseHeader, "correlation-id-response-header")
	return ops
}

// Create a new resource
func (r *ldapMappedScimHttpServletExtensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ldapMappedScimHttpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	debugLevel, err := client.NewEnumhttpServletExtensionDebugLevelPropFromValue(plan.DebugLevel.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for DebugLevel", err.Error())
		return
	}
	var DebugTypeSlice []client.EnumhttpServletExtensionDebugTypeProp
	plan.DebugType.ElementsAs(ctx, &DebugTypeSlice, false)
	addRequest := client.NewAddLdapMappedScimHttpServletExtensionRequest(plan.Id.ValueString(),
		[]client.EnumldapMappedScimHttpServletExtensionSchemaUrn{client.ENUMLDAPMAPPEDSCIMHTTPSERVLETEXTENSIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0HTTP_SERVLET_EXTENSIONLDAP_MAPPED_SCIM},
		plan.BaseContextPath.ValueString(),
		plan.TemporaryDirectory.ValueString(),
		plan.TemporaryDirectoryPermissions.ValueString(),
		*debugLevel,
		DebugTypeSlice,
		plan.IncludeStackTrace.ValueBool())
	addOptionalLdapMappedScimHttpServletExtensionFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.HttpServletExtensionApi.AddHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddHttpServletExtensionRequest(
		client.AddLdapMappedScimHttpServletExtensionRequestAsAddHttpServletExtensionRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.AddHttpServletExtensionExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Ldap Mapped Scim Http Servlet Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state ldapMappedScimHttpServletExtensionResourceModel
	readLdapMappedScimHttpServletExtensionResponse(ctx, addResponse.LdapMappedScimHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *ldapMappedScimHttpServletExtensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ldapMappedScimHttpServletExtensionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.GetHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ldap Mapped Scim Http Servlet Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readLdapMappedScimHttpServletExtensionResponse(ctx, readResponse.LdapMappedScimHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *ldapMappedScimHttpServletExtensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan ldapMappedScimHttpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state ldapMappedScimHttpServletExtensionResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createLdapMappedScimHttpServletExtensionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ldap Mapped Scim Http Servlet Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLdapMappedScimHttpServletExtensionResponse(ctx, updateResponse.LdapMappedScimHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
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
func (r *ldapMappedScimHttpServletExtensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ldapMappedScimHttpServletExtensionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.HttpServletExtensionApi.DeleteHttpServletExtensionExecute(r.apiClient.HttpServletExtensionApi.DeleteHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Ldap Mapped Scim Http Servlet Extension", err, httpResp)
		return
	}
}

func (r *ldapMappedScimHttpServletExtensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
