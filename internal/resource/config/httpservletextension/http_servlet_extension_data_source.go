package httpservletextension

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &httpServletExtensionDataSource{}
	_ datasource.DataSourceWithConfigure = &httpServletExtensionDataSource{}
)

// Create a Http Servlet Extension data source
func NewHttpServletExtensionDataSource() datasource.DataSource {
	return &httpServletExtensionDataSource{}
}

// httpServletExtensionDataSource is the datasource implementation.
type httpServletExtensionDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *httpServletExtensionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_http_servlet_extension"
}

// Configure adds the provider configured client to the data source.
func (r *httpServletExtensionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type httpServletExtensionDataSourceModel struct {
	Id                                 types.String `tfsdk:"id"`
	Type                               types.String `tfsdk:"type"`
	ExtensionClass                     types.String `tfsdk:"extension_class"`
	ExtensionArgument                  types.Set    `tfsdk:"extension_argument"`
	ScriptClass                        types.String `tfsdk:"script_class"`
	DocumentRootDirectory              types.String `tfsdk:"document_root_directory"`
	EnableDirectoryIndexing            types.Bool   `tfsdk:"enable_directory_indexing"`
	IndexFile                          types.Set    `tfsdk:"index_file"`
	ScriptArgument                     types.Set    `tfsdk:"script_argument"`
	OAuthTokenHandler                  types.String `tfsdk:"oauth_token_handler"`
	AllowedAuthenticationType          types.Set    `tfsdk:"allowed_authentication_type"`
	BaseContextPath                    types.String `tfsdk:"base_context_path"`
	IdTokenValidator                   types.Set    `tfsdk:"id_token_validator"`
	RequireFileServletAccessPrivilege  types.Bool   `tfsdk:"require_file_servlet_access_privilege"`
	RequireGroup                       types.Set    `tfsdk:"require_group"`
	ResourceMappingFile                types.String `tfsdk:"resource_mapping_file"`
	IncludeLDAPObjectclass             types.Set    `tfsdk:"include_ldap_objectclass"`
	ExcludeLDAPObjectclass             types.Set    `tfsdk:"exclude_ldap_objectclass"`
	IncludeLDAPBaseDN                  types.Set    `tfsdk:"include_ldap_base_dn"`
	ExcludeLDAPBaseDN                  types.Set    `tfsdk:"exclude_ldap_base_dn"`
	EntityTagLDAPAttribute             types.String `tfsdk:"entity_tag_ldap_attribute"`
	TemporaryDirectory                 types.String `tfsdk:"temporary_directory"`
	TemporaryDirectoryPermissions      types.String `tfsdk:"temporary_directory_permissions"`
	MaxResults                         types.Int64  `tfsdk:"max_results"`
	BulkMaxOperations                  types.Int64  `tfsdk:"bulk_max_operations"`
	BulkMaxPayloadSize                 types.String `tfsdk:"bulk_max_payload_size"`
	BulkMaxConcurrentRequests          types.Int64  `tfsdk:"bulk_max_concurrent_requests"`
	DebugEnabled                       types.Bool   `tfsdk:"debug_enabled"`
	DebugLevel                         types.String `tfsdk:"debug_level"`
	DebugType                          types.Set    `tfsdk:"debug_type"`
	IncludeStackTrace                  types.Bool   `tfsdk:"include_stack_trace"`
	MimeTypesFile                      types.String `tfsdk:"mime_types_file"`
	DefaultMIMEType                    types.String `tfsdk:"default_mime_type"`
	IncludeInstanceNameLabel           types.Bool   `tfsdk:"include_instance_name_label"`
	RequireAuthentication              types.Bool   `tfsdk:"require_authentication"`
	IncludeProductNameLabel            types.Bool   `tfsdk:"include_product_name_label"`
	IncludeLocationNameLabel           types.Bool   `tfsdk:"include_location_name_label"`
	AlwaysIncludeMonitorEntryNameLabel types.Bool   `tfsdk:"always_include_monitor_entry_name_label"`
	IncludeMonitorObjectClassNameLabel types.Bool   `tfsdk:"include_monitor_object_class_name_label"`
	IncludeMonitorAttributeNameLabel   types.Bool   `tfsdk:"include_monitor_attribute_name_label"`
	LabelNameValuePair                 types.Set    `tfsdk:"label_name_value_pair"`
	AvailableStatusCode                types.Int64  `tfsdk:"available_status_code"`
	DegradedStatusCode                 types.Int64  `tfsdk:"degraded_status_code"`
	UnavailableStatusCode              types.Int64  `tfsdk:"unavailable_status_code"`
	OverrideStatusCode                 types.Int64  `tfsdk:"override_status_code"`
	IncludeResponseBody                types.Bool   `tfsdk:"include_response_body"`
	AdditionalResponseContents         types.String `tfsdk:"additional_response_contents"`
	Server                             types.String `tfsdk:"server"`
	BasicAuthEnabled                   types.Bool   `tfsdk:"basic_auth_enabled"`
	IdentityMapper                     types.String `tfsdk:"identity_mapper"`
	AccessTokenValidator               types.Set    `tfsdk:"access_token_validator"`
	Description                        types.String `tfsdk:"description"`
	CrossOriginPolicy                  types.String `tfsdk:"cross_origin_policy"`
	ResponseHeader                     types.Set    `tfsdk:"response_header"`
	CorrelationIDResponseHeader        types.String `tfsdk:"correlation_id_response_header"`
}

// GetSchema defines the schema for the datasource.
func (r *httpServletExtensionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Describes a Http Servlet Extension.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Name of this object.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of HTTP Servlet Extension resource. Options are ['delegated-admin', 'quickstart', 'availability-state', 'prometheus-monitoring', 'velocity', 'consent', 'ldap-mapped-scim', 'groovy-scripted', 'file-server', 'config', 'scim2', 'directory-rest-api', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party HTTP Servlet Extension.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party HTTP Servlet Extension. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted HTTP Servlet Extension.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"document_root_directory": schema.StringAttribute{
				Description: "Specifies the path to the directory on the local filesystem containing the files to be served by this File Server HTTP Servlet Extension. The path must exist, and it must be a directory.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enable_directory_indexing": schema.BoolAttribute{
				Description: "Indicates whether to generate a default HTML page with a listing of available files if the requested path refers to a directory rather than a file, and that directory does not contain an index file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"index_file": schema.SetAttribute{
				Description: "Specifies the name of a file whose contents may be returned to the client if the requested path refers to a directory rather than a file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted HTTP Servlet Extension. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"oauth_token_handler": schema.StringAttribute{
				Description: "Specifies the OAuth Token Handler implementation that should be used to validate OAuth 2.0 bearer tokens when they are included in a SCIM request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allowed_authentication_type": schema.SetAttribute{
				Description: "The types of authentication that may be used to authenticate to the file servlet.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"base_context_path": schema.StringAttribute{
				Description: "Specifies the base context path that HTTP clients should use to access this servlet. The value must start with a forward slash and must represent a valid HTTP context path.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"id_token_validator": schema.SetAttribute{
				Description: "The ID token validators that may be used to verify the authenticity of an of an OpenID Connect ID token.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"require_file_servlet_access_privilege": schema.BoolAttribute{
				Description: "Indicates whether the servlet extension should only accept requests from authenticated clients that have the file-servlet-access privilege.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"require_group": schema.SetAttribute{
				Description: "The DN of a group whose members will be permitted to access to the associated files. If multiple group DNs are configured, then anyone who is a member of at least one of those groups will be granted access.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"resource_mapping_file": schema.StringAttribute{
				Description: "The path to an XML file defining the resources supported by the SCIM interface and the SCIM-to-LDAP attribute mappings to use.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_ldap_objectclass": schema.SetAttribute{
				Description: "Specifies the LDAP object classes that should be exposed directly as SCIM resources.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"exclude_ldap_objectclass": schema.SetAttribute{
				Description: "Specifies the LDAP object classes that should be not be exposed directly as SCIM resources.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"include_ldap_base_dn": schema.SetAttribute{
				Description: "Specifies the base DNs for the branches of the DIT that should be exposed via the Identity Access API.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"exclude_ldap_base_dn": schema.SetAttribute{
				Description: "Specifies the base DNs for the branches of the DIT that should not be exposed via the Identity Access API.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"entity_tag_ldap_attribute": schema.StringAttribute{
				Description: "Specifies the LDAP attribute whose value should be used as the entity tag value to enable SCIM resource versioning support.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"temporary_directory": schema.StringAttribute{
				Description: "Specifies the location of the directory that is used to create temporary files containing SCIM request data.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"temporary_directory_permissions": schema.StringAttribute{
				Description: "Specifies the permissions that should be applied to the directory that is used to create temporary files.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_results": schema.Int64Attribute{
				Description: "The maximum number of resources that are returned in a response.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"bulk_max_operations": schema.Int64Attribute{
				Description: "The maximum number of operations that are permitted in a bulk request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"bulk_max_payload_size": schema.StringAttribute{
				Description: "The maximum payload size in bytes of a bulk request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"bulk_max_concurrent_requests": schema.Int64Attribute{
				Description: "The maximum number of bulk requests that may be processed concurrently by the server. Any bulk request that would cause this limit to be exceeded is rejected with HTTP status code 503.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"debug_enabled": schema.BoolAttribute{
				Description: "Enables debug logging of the SCIM SDK. Debug messages will be forwarded to the Directory Server debug logger with the scope of com.unboundid.directory.server.extensions.scim.SCIMHTTPServletExtension.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"debug_level": schema.StringAttribute{
				Description: "The minimum debug level that should be used for messages to be logged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"debug_type": schema.SetAttribute{
				Description: "The types of debug messages that should be logged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"include_stack_trace": schema.BoolAttribute{
				Description: "Indicates whether a stack trace of the thread which called the debug method should be included in debug log messages.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"mime_types_file": schema.StringAttribute{
				Description: "Specifies the path to a file that contains MIME type mappings that will be used to determine the appropriate value to return for the Content-Type header based on the extension of the requested static content file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"default_mime_type": schema.StringAttribute{
				Description: "Specifies the default value that will be used in the response's Content-Type header that indicates the type of content to return.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_instance_name_label": schema.BoolAttribute{
				Description: "Indicates whether generated metrics should include an \"instance\" label whose value is the instance name for this Directory Server instance.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"require_authentication": schema.BoolAttribute{
				Description: "Require authentication when accessing Velocity templates.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_product_name_label": schema.BoolAttribute{
				Description: "Indicates whether generated metrics should include a \"product\" label whose value is the product name for this Directory Server instance.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_location_name_label": schema.BoolAttribute{
				Description: "Indicates whether generated metrics should include a \"location\" label whose value is the location name for this Directory Server instance.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"always_include_monitor_entry_name_label": schema.BoolAttribute{
				Description: "Indicates whether generated metrics should always include a \"monitor_entry\" label whose value is the name of the monitor entry from which the metric was obtained.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_monitor_object_class_name_label": schema.BoolAttribute{
				Description: "Indicates whether generated metrics should include a \"monitor_object_class\" label whose value is the name of the object class for the monitor entry from which the metric was obtained.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_monitor_attribute_name_label": schema.BoolAttribute{
				Description: "Indicates whether generated metrics should include a \"monitor_attribute\" label whose value is the name of the monitor attribute from which the metric was obtained.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"label_name_value_pair": schema.SetAttribute{
				Description: "A set of name-value pairs for labels that should be included in all metrics exposed by this Directory Server instance.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"available_status_code": schema.Int64Attribute{
				Description: "Specifies the HTTP status code that the servlet should return if the server considers itself to be available.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"degraded_status_code": schema.Int64Attribute{
				Description: "Specifies the HTTP status code that the servlet should return if the server considers itself to be degraded.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"unavailable_status_code": schema.Int64Attribute{
				Description: "Specifies the HTTP status code that the servlet should return if the server considers itself to be unavailable.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"override_status_code": schema.Int64Attribute{
				Description: "Specifies a HTTP status code that the servlet should always return, regardless of the server's availability. If this value is defined, it will override the availability-based return codes.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_response_body": schema.BoolAttribute{
				Description: "Indicates whether the response should include a body that is a JSON object.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"additional_response_contents": schema.StringAttribute{
				Description: "A JSON-formatted string containing additional fields to be returned in the response body. For example, an additional-response-contents value of '{ \"key\": \"value\" }' would add the key and value to the root of the JSON response body.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server": schema.StringAttribute{
				Description: "Specifies the PingFederate server to be configured.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"basic_auth_enabled": schema.BoolAttribute{
				Description: "Enables HTTP Basic authentication, using a username and password. The Identity Mapper specified by the identity-mapper property will be used to map the username to a DN.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"identity_mapper": schema.StringAttribute{
				Description: "Specifies the Identity Mapper that is to be used for associating user entries with basic authentication user names.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"access_token_validator": schema.SetAttribute{
				Description: "If specified, the Access Token Validator(s) that may be used to validate access tokens for requests submitted to this Delegated Admin HTTP Servlet Extension.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this HTTP Servlet Extension",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"cross_origin_policy": schema.StringAttribute{
				Description: "The cross-origin request policy to use for the HTTP Servlet Extension.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"response_header": schema.SetAttribute{
				Description: "Specifies HTTP header fields and values added to response headers for all requests.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"correlation_id_response_header": schema.StringAttribute{
				Description: "Specifies the name of the HTTP response header that will contain a correlation ID value. Example values are \"Correlation-Id\", \"X-Amzn-Trace-Id\", and \"X-Request-Id\".",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
}

// Read a QuickstartHttpServletExtensionResponse object into the model struct
func readQuickstartHttpServletExtensionResponseDataSource(ctx context.Context, r *client.QuickstartHttpServletExtensionResponse, state *httpServletExtensionDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("quickstart")
	state.Id = types.StringValue(r.Id)
	state.Server = internaltypes.StringTypeOrNil(r.Server, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, false)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, false)
}

// Read a AvailabilityStateHttpServletExtensionResponse object into the model struct
func readAvailabilityStateHttpServletExtensionResponseDataSource(ctx context.Context, r *client.AvailabilityStateHttpServletExtensionResponse, state *httpServletExtensionDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("availability-state")
	state.Id = types.StringValue(r.Id)
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.AvailableStatusCode = types.Int64Value(r.AvailableStatusCode)
	state.DegradedStatusCode = types.Int64Value(r.DegradedStatusCode)
	state.UnavailableStatusCode = types.Int64Value(r.UnavailableStatusCode)
	state.OverrideStatusCode = internaltypes.Int64TypeOrNil(r.OverrideStatusCode)
	state.IncludeResponseBody = internaltypes.BoolTypeOrNil(r.IncludeResponseBody)
	state.AdditionalResponseContents = internaltypes.StringTypeOrNil(r.AdditionalResponseContents, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, false)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, false)
}

// Read a PrometheusMonitoringHttpServletExtensionResponse object into the model struct
func readPrometheusMonitoringHttpServletExtensionResponseDataSource(ctx context.Context, r *client.PrometheusMonitoringHttpServletExtensionResponse, state *httpServletExtensionDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("prometheus-monitoring")
	state.Id = types.StringValue(r.Id)
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.IncludeInstanceNameLabel = internaltypes.BoolTypeOrNil(r.IncludeInstanceNameLabel)
	state.IncludeProductNameLabel = internaltypes.BoolTypeOrNil(r.IncludeProductNameLabel)
	state.IncludeLocationNameLabel = internaltypes.BoolTypeOrNil(r.IncludeLocationNameLabel)
	state.AlwaysIncludeMonitorEntryNameLabel = internaltypes.BoolTypeOrNil(r.AlwaysIncludeMonitorEntryNameLabel)
	state.IncludeMonitorObjectClassNameLabel = internaltypes.BoolTypeOrNil(r.IncludeMonitorObjectClassNameLabel)
	state.IncludeMonitorAttributeNameLabel = internaltypes.BoolTypeOrNil(r.IncludeMonitorAttributeNameLabel)
	state.LabelNameValuePair = internaltypes.GetStringSet(r.LabelNameValuePair)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, false)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, false)
}

// Read a LdapMappedScimHttpServletExtensionResponse object into the model struct
func readLdapMappedScimHttpServletExtensionResponseDataSource(ctx context.Context, r *client.LdapMappedScimHttpServletExtensionResponse, state *httpServletExtensionDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap-mapped-scim")
	state.Id = types.StringValue(r.Id)
	state.OAuthTokenHandler = internaltypes.StringTypeOrNil(r.OAuthTokenHandler, false)
	state.BasicAuthEnabled = internaltypes.BoolTypeOrNil(r.BasicAuthEnabled)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, false)
	state.ResourceMappingFile = internaltypes.StringTypeOrNil(r.ResourceMappingFile, false)
	state.IncludeLDAPObjectclass = internaltypes.GetStringSet(r.IncludeLDAPObjectclass)
	state.ExcludeLDAPObjectclass = internaltypes.GetStringSet(r.ExcludeLDAPObjectclass)
	state.IncludeLDAPBaseDN = internaltypes.GetStringSet(r.IncludeLDAPBaseDN)
	state.ExcludeLDAPBaseDN = internaltypes.GetStringSet(r.ExcludeLDAPBaseDN)
	state.EntityTagLDAPAttribute = internaltypes.StringTypeOrNil(r.EntityTagLDAPAttribute, false)
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.TemporaryDirectory = types.StringValue(r.TemporaryDirectory)
	state.TemporaryDirectoryPermissions = types.StringValue(r.TemporaryDirectoryPermissions)
	state.MaxResults = internaltypes.Int64TypeOrNil(r.MaxResults)
	state.BulkMaxOperations = internaltypes.Int64TypeOrNil(r.BulkMaxOperations)
	state.BulkMaxPayloadSize = internaltypes.StringTypeOrNil(r.BulkMaxPayloadSize, false)
	state.BulkMaxConcurrentRequests = internaltypes.Int64TypeOrNil(r.BulkMaxConcurrentRequests)
	state.DebugEnabled = internaltypes.BoolTypeOrNil(r.DebugEnabled)
	state.DebugLevel = types.StringValue(r.DebugLevel.String())
	state.DebugType = internaltypes.GetStringSet(
		client.StringSliceEnumhttpServletExtensionDebugTypeProp(r.DebugType))
	state.IncludeStackTrace = types.BoolValue(r.IncludeStackTrace)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, false)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, false)
}

// Read a GroovyScriptedHttpServletExtensionResponse object into the model struct
func readGroovyScriptedHttpServletExtensionResponseDataSource(ctx context.Context, r *client.GroovyScriptedHttpServletExtensionResponse, state *httpServletExtensionDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, false)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, false)
}

// Read a FileServerHttpServletExtensionResponse object into the model struct
func readFileServerHttpServletExtensionResponseDataSource(ctx context.Context, r *client.FileServerHttpServletExtensionResponse, state *httpServletExtensionDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-server")
	state.Id = types.StringValue(r.Id)
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.DocumentRootDirectory = types.StringValue(r.DocumentRootDirectory)
	state.EnableDirectoryIndexing = internaltypes.BoolTypeOrNil(r.EnableDirectoryIndexing)
	state.IndexFile = internaltypes.GetStringSet(r.IndexFile)
	state.MimeTypesFile = internaltypes.StringTypeOrNil(r.MimeTypesFile, false)
	state.DefaultMIMEType = internaltypes.StringTypeOrNil(r.DefaultMIMEType, false)
	state.RequireAuthentication = internaltypes.BoolTypeOrNil(r.RequireAuthentication)
	state.AllowedAuthenticationType = internaltypes.GetStringSet(
		client.StringSliceEnumhttpServletExtensionAllowedAuthenticationTypeProp(r.AllowedAuthenticationType))
	state.AccessTokenValidator = internaltypes.GetStringSet(r.AccessTokenValidator)
	state.IdTokenValidator = internaltypes.GetStringSet(r.IdTokenValidator)
	state.RequireFileServletAccessPrivilege = internaltypes.BoolTypeOrNil(r.RequireFileServletAccessPrivilege)
	state.RequireGroup = internaltypes.GetStringSet(r.RequireGroup)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, false)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, false)
}

// Read a ThirdPartyHttpServletExtensionResponse object into the model struct
func readThirdPartyHttpServletExtensionResponseDataSource(ctx context.Context, r *client.ThirdPartyHttpServletExtensionResponse, state *httpServletExtensionDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, false)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, false)
}

// Read resource information
func (r *httpServletExtensionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state httpServletExtensionDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.GetHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Http Servlet Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.QuickstartHttpServletExtensionResponse != nil {
		readQuickstartHttpServletExtensionResponseDataSource(ctx, readResponse.QuickstartHttpServletExtensionResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AvailabilityStateHttpServletExtensionResponse != nil {
		readAvailabilityStateHttpServletExtensionResponseDataSource(ctx, readResponse.AvailabilityStateHttpServletExtensionResponse, &state, &resp.Diagnostics)
	}
	if readResponse.PrometheusMonitoringHttpServletExtensionResponse != nil {
		readPrometheusMonitoringHttpServletExtensionResponseDataSource(ctx, readResponse.PrometheusMonitoringHttpServletExtensionResponse, &state, &resp.Diagnostics)
	}
	if readResponse.LdapMappedScimHttpServletExtensionResponse != nil {
		readLdapMappedScimHttpServletExtensionResponseDataSource(ctx, readResponse.LdapMappedScimHttpServletExtensionResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedHttpServletExtensionResponse != nil {
		readGroovyScriptedHttpServletExtensionResponseDataSource(ctx, readResponse.GroovyScriptedHttpServletExtensionResponse, &state, &resp.Diagnostics)
	}
	if readResponse.FileServerHttpServletExtensionResponse != nil {
		readFileServerHttpServletExtensionResponseDataSource(ctx, readResponse.FileServerHttpServletExtensionResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyHttpServletExtensionResponse != nil {
		readThirdPartyHttpServletExtensionResponseDataSource(ctx, readResponse.ThirdPartyHttpServletExtensionResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
