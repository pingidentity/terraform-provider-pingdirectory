package httpservletextension

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &httpServletExtensionResource{}
	_ resource.ResourceWithConfigure   = &httpServletExtensionResource{}
	_ resource.ResourceWithImportState = &httpServletExtensionResource{}
	_ resource.Resource                = &defaultHttpServletExtensionResource{}
	_ resource.ResourceWithConfigure   = &defaultHttpServletExtensionResource{}
	_ resource.ResourceWithImportState = &defaultHttpServletExtensionResource{}
)

// Create a Http Servlet Extension resource
func NewHttpServletExtensionResource() resource.Resource {
	return &httpServletExtensionResource{}
}

func NewDefaultHttpServletExtensionResource() resource.Resource {
	return &defaultHttpServletExtensionResource{}
}

// httpServletExtensionResource is the resource implementation.
type httpServletExtensionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultHttpServletExtensionResource is the resource implementation.
type defaultHttpServletExtensionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *httpServletExtensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_http_servlet_extension"
}

func (r *defaultHttpServletExtensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_http_servlet_extension"
}

// Configure adds the provider configured client to the resource.
func (r *httpServletExtensionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultHttpServletExtensionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type httpServletExtensionResourceModel struct {
	Id                                 types.String `tfsdk:"id"`
	Name                               types.String `tfsdk:"name"`
	LastUpdated                        types.String `tfsdk:"last_updated"`
	Notifications                      types.Set    `tfsdk:"notifications"`
	RequiredActions                    types.Set    `tfsdk:"required_actions"`
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

type defaultHttpServletExtensionResourceModel struct {
	Id                                 types.String `tfsdk:"id"`
	Name                               types.String `tfsdk:"name"`
	LastUpdated                        types.String `tfsdk:"last_updated"`
	Notifications                      types.Set    `tfsdk:"notifications"`
	RequiredActions                    types.Set    `tfsdk:"required_actions"`
	Type                               types.String `tfsdk:"type"`
	ExtensionClass                     types.String `tfsdk:"extension_class"`
	ExtensionArgument                  types.Set    `tfsdk:"extension_argument"`
	ScriptClass                        types.String `tfsdk:"script_class"`
	DocumentRootDirectory              types.String `tfsdk:"document_root_directory"`
	MapAccessTokensToLocalUsers        types.String `tfsdk:"map_access_tokens_to_local_users"`
	EnableDirectoryIndexing            types.Bool   `tfsdk:"enable_directory_indexing"`
	IndexFile                          types.Set    `tfsdk:"index_file"`
	MaxPageSize                        types.Int64  `tfsdk:"max_page_size"`
	SchemasEndpointObjectclass         types.Set    `tfsdk:"schemas_endpoint_objectclass"`
	DefaultOperationalAttribute        types.Set    `tfsdk:"default_operational_attribute"`
	RejectExpansionAttribute           types.Set    `tfsdk:"reject_expansion_attribute"`
	AlwaysUsePermissiveModify          types.Bool   `tfsdk:"always_use_permissive_modify"`
	AllowedControl                     types.Set    `tfsdk:"allowed_control"`
	ScriptArgument                     types.Set    `tfsdk:"script_argument"`
	OAuthTokenHandler                  types.String `tfsdk:"oauth_token_handler"`
	SwaggerEnabled                     types.Bool   `tfsdk:"swagger_enabled"`
	BearerTokenAuthEnabled             types.Bool   `tfsdk:"bearer_token_auth_enabled"`
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
	StaticContextPath                  types.String `tfsdk:"static_context_path"`
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
	StaticContentDirectory             types.String `tfsdk:"static_content_directory"`
	StaticCustomDirectory              types.String `tfsdk:"static_custom_directory"`
	TemplateDirectory                  types.Set    `tfsdk:"template_directory"`
	ExposeRequestAttributes            types.Bool   `tfsdk:"expose_request_attributes"`
	ExposeSessionAttributes            types.Bool   `tfsdk:"expose_session_attributes"`
	ExposeServerContext                types.Bool   `tfsdk:"expose_server_context"`
	AllowContextOverride               types.Bool   `tfsdk:"allow_context_override"`
	MimeTypesFile                      types.String `tfsdk:"mime_types_file"`
	DefaultMIMEType                    types.String `tfsdk:"default_mime_type"`
	CharacterEncoding                  types.String `tfsdk:"character_encoding"`
	IncludeInstanceNameLabel           types.Bool   `tfsdk:"include_instance_name_label"`
	StaticResponseHeader               types.Set    `tfsdk:"static_response_header"`
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
	AccessTokenScope                   types.String `tfsdk:"access_token_scope"`
	Audience                           types.String `tfsdk:"audience"`
	Description                        types.String `tfsdk:"description"`
	CrossOriginPolicy                  types.String `tfsdk:"cross_origin_policy"`
	ResponseHeader                     types.Set    `tfsdk:"response_header"`
	CorrelationIDResponseHeader        types.String `tfsdk:"correlation_id_response_header"`
}

// GetSchema defines the schema for the resource.
func (r *httpServletExtensionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	httpServletExtensionSchema(ctx, req, resp, false)
}

func (r *defaultHttpServletExtensionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	httpServletExtensionSchema(ctx, req, resp, true)
}

func httpServletExtensionSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Http Servlet Extension.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of HTTP Servlet Extension resource. Options are ['delegated-admin', 'quickstart', 'availability-state', 'prometheus-monitoring', 'velocity', 'consent', 'ldap-mapped-scim', 'groovy-scripted', 'file-server', 'config', 'scim2', 'directory-rest-api', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"quickstart", "availability-state", "prometheus-monitoring", "ldap-mapped-scim", "groovy-scripted", "file-server", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party HTTP Servlet Extension.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party HTTP Servlet Extension. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted HTTP Servlet Extension.",
				Optional:    true,
			},
			"document_root_directory": schema.StringAttribute{
				Description: "Specifies the path to the directory on the local filesystem containing the files to be served by this File Server HTTP Servlet Extension. The path must exist, and it must be a directory.",
				Optional:    true,
			},
			"enable_directory_indexing": schema.BoolAttribute{
				Description: "Indicates whether to generate a default HTML page with a listing of available files if the requested path refers to a directory rather than a file, and that directory does not contain an index file.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"index_file": schema.SetAttribute{
				Description: "Specifies the name of a file whose contents may be returned to the client if the requested path refers to a directory rather than a file.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted HTTP Servlet Extension. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"oauth_token_handler": schema.StringAttribute{
				Description: "Specifies the OAuth Token Handler implementation that should be used to validate OAuth 2.0 bearer tokens when they are included in a SCIM request.",
				Optional:    true,
			},
			"allowed_authentication_type": schema.SetAttribute{
				Description: "The types of authentication that may be used to authenticate to the file servlet.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"base_context_path": schema.StringAttribute{
				Description: "Specifies the base context path that HTTP clients should use to access this servlet. The value must start with a forward slash and must represent a valid HTTP context path.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id_token_validator": schema.SetAttribute{
				Description: "The ID token validators that may be used to verify the authenticity of an of an OpenID Connect ID token.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"require_file_servlet_access_privilege": schema.BoolAttribute{
				Description: "Indicates whether the servlet extension should only accept requests from authenticated clients that have the file-servlet-access privilege.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"require_group": schema.SetAttribute{
				Description: "The DN of a group whose members will be permitted to access to the associated files. If multiple group DNs are configured, then anyone who is a member of at least one of those groups will be granted access.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_mapping_file": schema.StringAttribute{
				Description: "The path to an XML file defining the resources supported by the SCIM interface and the SCIM-to-LDAP attribute mappings to use.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"include_ldap_objectclass": schema.SetAttribute{
				Description: "Specifies the LDAP object classes that should be exposed directly as SCIM resources.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"exclude_ldap_objectclass": schema.SetAttribute{
				Description: "Specifies the LDAP object classes that should be not be exposed directly as SCIM resources.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"include_ldap_base_dn": schema.SetAttribute{
				Description: "Specifies the base DNs for the branches of the DIT that should be exposed via the Identity Access API.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"exclude_ldap_base_dn": schema.SetAttribute{
				Description: "Specifies the base DNs for the branches of the DIT that should not be exposed via the Identity Access API.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"entity_tag_ldap_attribute": schema.StringAttribute{
				Description: "Specifies the LDAP attribute whose value should be used as the entity tag value to enable SCIM resource versioning support.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"temporary_directory": schema.StringAttribute{
				Description: "Specifies the location of the directory that is used to create temporary files containing SCIM request data.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"temporary_directory_permissions": schema.StringAttribute{
				Description: "Specifies the permissions that should be applied to the directory that is used to create temporary files.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"max_results": schema.Int64Attribute{
				Description: "The maximum number of resources that are returned in a response.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"bulk_max_operations": schema.Int64Attribute{
				Description: "The maximum number of operations that are permitted in a bulk request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"bulk_max_payload_size": schema.StringAttribute{
				Description: "The maximum payload size in bytes of a bulk request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bulk_max_concurrent_requests": schema.Int64Attribute{
				Description: "The maximum number of bulk requests that may be processed concurrently by the server. Any bulk request that would cause this limit to be exceeded is rejected with HTTP status code 503.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"debug_enabled": schema.BoolAttribute{
				Description: "Enables debug logging of the SCIM SDK. Debug messages will be forwarded to the Directory Server debug logger with the scope of com.unboundid.directory.server.extensions.scim.SCIMHTTPServletExtension.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"debug_level": schema.StringAttribute{
				Description: "The minimum debug level that should be used for messages to be logged.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"debug_type": schema.SetAttribute{
				Description: "The types of debug messages that should be logged.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"include_stack_trace": schema.BoolAttribute{
				Description: "Indicates whether a stack trace of the thread which called the debug method should be included in debug log messages.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"mime_types_file": schema.StringAttribute{
				Description: "Specifies the path to a file that contains MIME type mappings that will be used to determine the appropriate value to return for the Content-Type header based on the extension of the requested static content file.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"default_mime_type": schema.StringAttribute{
				Description: "Specifies the default value that will be used in the response's Content-Type header that indicates the type of content to return.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"include_instance_name_label": schema.BoolAttribute{
				Description: "Indicates whether generated metrics should include an \"instance\" label whose value is the instance name for this Directory Server instance.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"require_authentication": schema.BoolAttribute{
				Description: "Require authentication when accessing Velocity templates.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_product_name_label": schema.BoolAttribute{
				Description: "Indicates whether generated metrics should include a \"product\" label whose value is the product name for this Directory Server instance.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_location_name_label": schema.BoolAttribute{
				Description: "Indicates whether generated metrics should include a \"location\" label whose value is the location name for this Directory Server instance.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"always_include_monitor_entry_name_label": schema.BoolAttribute{
				Description: "Indicates whether generated metrics should always include a \"monitor_entry\" label whose value is the name of the monitor entry from which the metric was obtained.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_monitor_object_class_name_label": schema.BoolAttribute{
				Description: "Indicates whether generated metrics should include a \"monitor_object_class\" label whose value is the name of the object class for the monitor entry from which the metric was obtained.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_monitor_attribute_name_label": schema.BoolAttribute{
				Description: "Indicates whether generated metrics should include a \"monitor_attribute\" label whose value is the name of the monitor attribute from which the metric was obtained.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"label_name_value_pair": schema.SetAttribute{
				Description: "A set of name-value pairs for labels that should be included in all metrics exposed by this Directory Server instance.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"available_status_code": schema.Int64Attribute{
				Description: "Specifies the HTTP status code that the servlet should return if the server considers itself to be available.",
				Optional:    true,
			},
			"degraded_status_code": schema.Int64Attribute{
				Description: "Specifies the HTTP status code that the servlet should return if the server considers itself to be degraded.",
				Optional:    true,
			},
			"unavailable_status_code": schema.Int64Attribute{
				Description: "Specifies the HTTP status code that the servlet should return if the server considers itself to be unavailable.",
				Optional:    true,
			},
			"override_status_code": schema.Int64Attribute{
				Description: "Specifies a HTTP status code that the servlet should always return, regardless of the server's availability. If this value is defined, it will override the availability-based return codes.",
				Optional:    true,
			},
			"include_response_body": schema.BoolAttribute{
				Description: "Indicates whether the response should include a body that is a JSON object.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"additional_response_contents": schema.StringAttribute{
				Description: "A JSON-formatted string containing additional fields to be returned in the response body. For example, an additional-response-contents value of '{ \"key\": \"value\" }' would add the key and value to the root of the JSON response body.",
				Optional:    true,
			},
			"server": schema.StringAttribute{
				Description: "Specifies the PingFederate server to be configured.",
				Optional:    true,
			},
			"basic_auth_enabled": schema.BoolAttribute{
				Description: "Enables HTTP Basic authentication, using a username and password. The Identity Mapper specified by the identity-mapper property will be used to map the username to a DN.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"identity_mapper": schema.StringAttribute{
				Description: "Specifies the Identity Mapper that is to be used for associating user entries with basic authentication user names.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"access_token_validator": schema.SetAttribute{
				Description: "If specified, the Access Token Validator(s) that may be used to validate access tokens for requests submitted to this Delegated Admin HTTP Servlet Extension.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this HTTP Servlet Extension",
				Optional:    true,
			},
			"cross_origin_policy": schema.StringAttribute{
				Description: "The cross-origin request policy to use for the HTTP Servlet Extension.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"response_header": schema.SetAttribute{
				Description: "Specifies HTTP header fields and values added to response headers for all requests.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"correlation_id_response_header": schema.StringAttribute{
				Description: "Specifies the name of the HTTP response header that will contain a correlation ID value. Example values are \"Correlation-Id\", \"X-Amzn-Trace-Id\", and \"X-Request-Id\".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Validators = []validator.String{
			stringvalidator.OneOf([]string{"delegated-admin", "quickstart", "availability-state", "prometheus-monitoring", "velocity", "consent", "ldap-mapped-scim", "groovy-scripted", "file-server", "config", "scim2", "directory-rest-api", "third-party"}...),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		schemaDef.Attributes["map_access_tokens_to_local_users"] = schema.StringAttribute{
			Description: "Indicates whether the SCIM2 servlet should attempt to map the presented access token to a local user.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["max_page_size"] = schema.Int64Attribute{
			Description: "The maximum number of entries to be returned in one page of search results.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["schemas_endpoint_objectclass"] = schema.SetAttribute{
			Description: "The list of object classes which will be returned by the schemas endpoint.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["default_operational_attribute"] = schema.SetAttribute{
			Description: "A set of operational attributes that will be returned with entries by default.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["reject_expansion_attribute"] = schema.SetAttribute{
			Description: "A set of attributes which the client is not allowed to provide for the expand query parameters. This should be used for attributes that could either have a large number of values or that reference entries that are very large like groups.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["always_use_permissive_modify"] = schema.BoolAttribute{
			Description: "Indicates whether to always use permissive modify behavior for PATCH requests, even if the request did not include the permissive modify request control. Supported in PingDirectory product version 9.3.0.0+.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["allowed_control"] = schema.SetAttribute{
			Description: "Specifies the names of any request controls that should be allowed by the Directory REST API. Any request that contains a critical control not in this list will be rejected. Any non-critical request control which is not supported by the Directory REST API will be removed from the request.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["swagger_enabled"] = schema.BoolAttribute{
			Description: "Indicates whether the SCIM2 HTTP Servlet Extension will generate a Swagger specification document.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["bearer_token_auth_enabled"] = schema.BoolAttribute{
			Description: "Enables HTTP bearer token authentication.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["static_context_path"] = schema.StringAttribute{
			Description: "The path below the base context path by which static, non-template content such as images, CSS, and Javascript files are accessible.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["static_content_directory"] = schema.StringAttribute{
			Description: "Specifies the base directory in which static, non-template content such as images, CSS, and Javascript files are stored on the filesystem.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["static_custom_directory"] = schema.StringAttribute{
			Description: "Specifies the base directory in which custom static, non-template content such as images, CSS, and Javascript files are stored on the filesystem. Files in this directory will override those with the same name in the directory specified by the static-content-directory property.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["template_directory"] = schema.SetAttribute{
			Description: "Specifies an ordered list of directories in which to search for the template files.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["expose_request_attributes"] = schema.BoolAttribute{
			Description: "Specifies whether the HTTP request will be exposed to templates.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["expose_session_attributes"] = schema.BoolAttribute{
			Description: "Specifies whether the HTTP session will be exposed to templates.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["expose_server_context"] = schema.BoolAttribute{
			Description: "Specifies whether a server context will be exposed under context key 'ubid_server' for all template contexts.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["allow_context_override"] = schema.BoolAttribute{
			Description: "Indicates whether context providers may override existing context objects with new values.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["character_encoding"] = schema.StringAttribute{
			Description: "Specifies the value that will be used for all responses' Content-Type headers' charset parameter that indicates the character encoding of the document.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["static_response_header"] = schema.SetAttribute{
			Description: "Specifies HTTP header fields and values added to response headers for static content requests such as images and scripts.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["access_token_scope"] = schema.StringAttribute{
			Description: "The name of a scope that must be present in an access token accepted by the Delegated Admin HTTP Servlet Extension.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["audience"] = schema.StringAttribute{
			Description: "A string or URI that identifies the Delegated Admin HTTP Servlet Extension in the context of OAuth2 authorization.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		config.SetAllAttributesToOptionalAndComputed(&schemaDef)
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan
func (r *httpServletExtensionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanHttpServletExtension(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_http_servlet_extension")
}

func (r *defaultHttpServletExtensionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanHttpServletExtension(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_default_http_servlet_extension")
}

func modifyPlanHttpServletExtension(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, resourceName string) {
	var model defaultHttpServletExtensionResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.IdTokenValidator) && model.Type.ValueString() != "file-server" {
		resp.Diagnostics.AddError("Attribute 'id_token_validator' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'id_token_validator', the 'type' attribute must be one of ['file-server']")
	}
	if internaltypes.IsDefined(model.AccessTokenValidator) && model.Type.ValueString() != "delegated-admin" && model.Type.ValueString() != "file-server" && model.Type.ValueString() != "consent" && model.Type.ValueString() != "scim2" && model.Type.ValueString() != "directory-rest-api" {
		resp.Diagnostics.AddError("Attribute 'access_token_validator' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'access_token_validator', the 'type' attribute must be one of ['delegated-admin', 'file-server', 'consent', 'scim2', 'directory-rest-api']")
	}
	if internaltypes.IsDefined(model.StaticContentDirectory) && model.Type.ValueString() != "velocity" {
		resp.Diagnostics.AddError("Attribute 'static_content_directory' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'static_content_directory', the 'type' attribute must be one of ['velocity']")
	}
	if internaltypes.IsDefined(model.LabelNameValuePair) && model.Type.ValueString() != "prometheus-monitoring" {
		resp.Diagnostics.AddError("Attribute 'label_name_value_pair' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'label_name_value_pair', the 'type' attribute must be one of ['prometheus-monitoring']")
	}
	if internaltypes.IsDefined(model.IncludeLDAPBaseDN) && model.Type.ValueString() != "ldap-mapped-scim" {
		resp.Diagnostics.AddError("Attribute 'include_ldap_base_dn' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_ldap_base_dn', the 'type' attribute must be one of ['ldap-mapped-scim']")
	}
	if internaltypes.IsDefined(model.AdditionalResponseContents) && model.Type.ValueString() != "availability-state" {
		resp.Diagnostics.AddError("Attribute 'additional_response_contents' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'additional_response_contents', the 'type' attribute must be one of ['availability-state']")
	}
	if internaltypes.IsDefined(model.MapAccessTokensToLocalUsers) && model.Type.ValueString() != "scim2" {
		resp.Diagnostics.AddError("Attribute 'map_access_tokens_to_local_users' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'map_access_tokens_to_local_users', the 'type' attribute must be one of ['scim2']")
	}
	if internaltypes.IsDefined(model.ExposeSessionAttributes) && model.Type.ValueString() != "velocity" {
		resp.Diagnostics.AddError("Attribute 'expose_session_attributes' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'expose_session_attributes', the 'type' attribute must be one of ['velocity']")
	}
	if internaltypes.IsDefined(model.RequireFileServletAccessPrivilege) && model.Type.ValueString() != "file-server" {
		resp.Diagnostics.AddError("Attribute 'require_file_servlet_access_privilege' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'require_file_servlet_access_privilege', the 'type' attribute must be one of ['file-server']")
	}
	if internaltypes.IsDefined(model.DegradedStatusCode) && model.Type.ValueString() != "availability-state" {
		resp.Diagnostics.AddError("Attribute 'degraded_status_code' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'degraded_status_code', the 'type' attribute must be one of ['availability-state']")
	}
	if internaltypes.IsDefined(model.AvailableStatusCode) && model.Type.ValueString() != "availability-state" {
		resp.Diagnostics.AddError("Attribute 'available_status_code' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'available_status_code', the 'type' attribute must be one of ['availability-state']")
	}
	if internaltypes.IsDefined(model.RequireAuthentication) && model.Type.ValueString() != "file-server" && model.Type.ValueString() != "velocity" {
		resp.Diagnostics.AddError("Attribute 'require_authentication' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'require_authentication', the 'type' attribute must be one of ['file-server', 'velocity']")
	}
	if internaltypes.IsDefined(model.EnableDirectoryIndexing) && model.Type.ValueString() != "file-server" {
		resp.Diagnostics.AddError("Attribute 'enable_directory_indexing' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'enable_directory_indexing', the 'type' attribute must be one of ['file-server']")
	}
	if internaltypes.IsDefined(model.SwaggerEnabled) && model.Type.ValueString() != "scim2" {
		resp.Diagnostics.AddError("Attribute 'swagger_enabled' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'swagger_enabled', the 'type' attribute must be one of ['scim2']")
	}
	if internaltypes.IsDefined(model.IncludeMonitorAttributeNameLabel) && model.Type.ValueString() != "prometheus-monitoring" {
		resp.Diagnostics.AddError("Attribute 'include_monitor_attribute_name_label' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_monitor_attribute_name_label', the 'type' attribute must be one of ['prometheus-monitoring']")
	}
	if internaltypes.IsDefined(model.BearerTokenAuthEnabled) && model.Type.ValueString() != "consent" {
		resp.Diagnostics.AddError("Attribute 'bearer_token_auth_enabled' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'bearer_token_auth_enabled', the 'type' attribute must be one of ['consent']")
	}
	if internaltypes.IsDefined(model.Audience) && model.Type.ValueString() != "delegated-admin" && model.Type.ValueString() != "directory-rest-api" {
		resp.Diagnostics.AddError("Attribute 'audience' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'audience', the 'type' attribute must be one of ['delegated-admin', 'directory-rest-api']")
	}
	if internaltypes.IsDefined(model.ExposeRequestAttributes) && model.Type.ValueString() != "velocity" {
		resp.Diagnostics.AddError("Attribute 'expose_request_attributes' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'expose_request_attributes', the 'type' attribute must be one of ['velocity']")
	}
	if internaltypes.IsDefined(model.BulkMaxPayloadSize) && model.Type.ValueString() != "ldap-mapped-scim" {
		resp.Diagnostics.AddError("Attribute 'bulk_max_payload_size' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'bulk_max_payload_size', the 'type' attribute must be one of ['ldap-mapped-scim']")
	}
	if internaltypes.IsDefined(model.AllowedAuthenticationType) && model.Type.ValueString() != "file-server" {
		resp.Diagnostics.AddError("Attribute 'allowed_authentication_type' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'allowed_authentication_type', the 'type' attribute must be one of ['file-server']")
	}
	if internaltypes.IsDefined(model.StaticContextPath) && model.Type.ValueString() != "velocity" {
		resp.Diagnostics.AddError("Attribute 'static_context_path' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'static_context_path', the 'type' attribute must be one of ['velocity']")
	}
	if internaltypes.IsDefined(model.ResourceMappingFile) && model.Type.ValueString() != "ldap-mapped-scim" {
		resp.Diagnostics.AddError("Attribute 'resource_mapping_file' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'resource_mapping_file', the 'type' attribute must be one of ['ldap-mapped-scim']")
	}
	if internaltypes.IsDefined(model.RejectExpansionAttribute) && model.Type.ValueString() != "directory-rest-api" {
		resp.Diagnostics.AddError("Attribute 'reject_expansion_attribute' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'reject_expansion_attribute', the 'type' attribute must be one of ['directory-rest-api']")
	}
	if internaltypes.IsDefined(model.BulkMaxConcurrentRequests) && model.Type.ValueString() != "ldap-mapped-scim" {
		resp.Diagnostics.AddError("Attribute 'bulk_max_concurrent_requests' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'bulk_max_concurrent_requests', the 'type' attribute must be one of ['ldap-mapped-scim']")
	}
	if internaltypes.IsDefined(model.ExtensionClass) && model.Type.ValueString() != "third-party" {
		resp.Diagnostics.AddError("Attribute 'extension_class' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'extension_class', the 'type' attribute must be one of ['third-party']")
	}
	if internaltypes.IsDefined(model.ScriptClass) && model.Type.ValueString() != "groovy-scripted" {
		resp.Diagnostics.AddError("Attribute 'script_class' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'script_class', the 'type' attribute must be one of ['groovy-scripted']")
	}
	if internaltypes.IsDefined(model.IncludeLocationNameLabel) && model.Type.ValueString() != "prometheus-monitoring" {
		resp.Diagnostics.AddError("Attribute 'include_location_name_label' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_location_name_label', the 'type' attribute must be one of ['prometheus-monitoring']")
	}
	if internaltypes.IsDefined(model.IncludeResponseBody) && model.Type.ValueString() != "availability-state" {
		resp.Diagnostics.AddError("Attribute 'include_response_body' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_response_body', the 'type' attribute must be one of ['availability-state']")
	}
	if internaltypes.IsDefined(model.Server) && model.Type.ValueString() != "quickstart" {
		resp.Diagnostics.AddError("Attribute 'server' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'server', the 'type' attribute must be one of ['quickstart']")
	}
	if internaltypes.IsDefined(model.TemporaryDirectory) && model.Type.ValueString() != "ldap-mapped-scim" {
		resp.Diagnostics.AddError("Attribute 'temporary_directory' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'temporary_directory', the 'type' attribute must be one of ['ldap-mapped-scim']")
	}
	if internaltypes.IsDefined(model.AllowContextOverride) && model.Type.ValueString() != "velocity" {
		resp.Diagnostics.AddError("Attribute 'allow_context_override' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'allow_context_override', the 'type' attribute must be one of ['velocity']")
	}
	if internaltypes.IsDefined(model.TemporaryDirectoryPermissions) && model.Type.ValueString() != "ldap-mapped-scim" {
		resp.Diagnostics.AddError("Attribute 'temporary_directory_permissions' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'temporary_directory_permissions', the 'type' attribute must be one of ['ldap-mapped-scim']")
	}
	if internaltypes.IsDefined(model.MaxResults) && model.Type.ValueString() != "ldap-mapped-scim" {
		resp.Diagnostics.AddError("Attribute 'max_results' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'max_results', the 'type' attribute must be one of ['ldap-mapped-scim']")
	}
	if internaltypes.IsDefined(model.ExposeServerContext) && model.Type.ValueString() != "velocity" {
		resp.Diagnostics.AddError("Attribute 'expose_server_context' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'expose_server_context', the 'type' attribute must be one of ['velocity']")
	}
	if internaltypes.IsDefined(model.DebugEnabled) && model.Type.ValueString() != "scim2" && model.Type.ValueString() != "ldap-mapped-scim" {
		resp.Diagnostics.AddError("Attribute 'debug_enabled' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'debug_enabled', the 'type' attribute must be one of ['scim2', 'ldap-mapped-scim']")
	}
	if internaltypes.IsDefined(model.BaseContextPath) && model.Type.ValueString() != "availability-state" && model.Type.ValueString() != "prometheus-monitoring" && model.Type.ValueString() != "file-server" && model.Type.ValueString() != "velocity" && model.Type.ValueString() != "scim2" && model.Type.ValueString() != "ldap-mapped-scim" {
		resp.Diagnostics.AddError("Attribute 'base_context_path' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'base_context_path', the 'type' attribute must be one of ['availability-state', 'prometheus-monitoring', 'file-server', 'velocity', 'scim2', 'ldap-mapped-scim']")
	}
	if internaltypes.IsDefined(model.DocumentRootDirectory) && model.Type.ValueString() != "file-server" {
		resp.Diagnostics.AddError("Attribute 'document_root_directory' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'document_root_directory', the 'type' attribute must be one of ['file-server']")
	}
	if internaltypes.IsDefined(model.SchemasEndpointObjectclass) && model.Type.ValueString() != "directory-rest-api" {
		resp.Diagnostics.AddError("Attribute 'schemas_endpoint_objectclass' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'schemas_endpoint_objectclass', the 'type' attribute must be one of ['directory-rest-api']")
	}
	if internaltypes.IsDefined(model.IncludeProductNameLabel) && model.Type.ValueString() != "prometheus-monitoring" {
		resp.Diagnostics.AddError("Attribute 'include_product_name_label' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_product_name_label', the 'type' attribute must be one of ['prometheus-monitoring']")
	}
	if internaltypes.IsDefined(model.BasicAuthEnabled) && model.Type.ValueString() != "delegated-admin" && model.Type.ValueString() != "consent" && model.Type.ValueString() != "directory-rest-api" && model.Type.ValueString() != "ldap-mapped-scim" {
		resp.Diagnostics.AddError("Attribute 'basic_auth_enabled' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'basic_auth_enabled', the 'type' attribute must be one of ['delegated-admin', 'consent', 'directory-rest-api', 'ldap-mapped-scim']")
	}
	if internaltypes.IsDefined(model.IndexFile) && model.Type.ValueString() != "file-server" {
		resp.Diagnostics.AddError("Attribute 'index_file' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'index_file', the 'type' attribute must be one of ['file-server']")
	}
	if internaltypes.IsDefined(model.IncludeMonitorObjectClassNameLabel) && model.Type.ValueString() != "prometheus-monitoring" {
		resp.Diagnostics.AddError("Attribute 'include_monitor_object_class_name_label' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_monitor_object_class_name_label', the 'type' attribute must be one of ['prometheus-monitoring']")
	}
	if internaltypes.IsDefined(model.ExtensionArgument) && model.Type.ValueString() != "third-party" {
		resp.Diagnostics.AddError("Attribute 'extension_argument' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'extension_argument', the 'type' attribute must be one of ['third-party']")
	}
	if internaltypes.IsDefined(model.ScriptArgument) && model.Type.ValueString() != "groovy-scripted" {
		resp.Diagnostics.AddError("Attribute 'script_argument' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'script_argument', the 'type' attribute must be one of ['groovy-scripted']")
	}
	if internaltypes.IsDefined(model.AccessTokenScope) && model.Type.ValueString() != "delegated-admin" && model.Type.ValueString() != "directory-rest-api" {
		resp.Diagnostics.AddError("Attribute 'access_token_scope' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'access_token_scope', the 'type' attribute must be one of ['delegated-admin', 'directory-rest-api']")
	}
	if internaltypes.IsDefined(model.RequireGroup) && model.Type.ValueString() != "file-server" {
		resp.Diagnostics.AddError("Attribute 'require_group' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'require_group', the 'type' attribute must be one of ['file-server']")
	}
	if internaltypes.IsDefined(model.MimeTypesFile) && model.Type.ValueString() != "file-server" && model.Type.ValueString() != "velocity" {
		resp.Diagnostics.AddError("Attribute 'mime_types_file' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'mime_types_file', the 'type' attribute must be one of ['file-server', 'velocity']")
	}
	if internaltypes.IsDefined(model.IncludeStackTrace) && model.Type.ValueString() != "scim2" && model.Type.ValueString() != "ldap-mapped-scim" {
		resp.Diagnostics.AddError("Attribute 'include_stack_trace' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_stack_trace', the 'type' attribute must be one of ['scim2', 'ldap-mapped-scim']")
	}
	if internaltypes.IsDefined(model.IncludeInstanceNameLabel) && model.Type.ValueString() != "prometheus-monitoring" {
		resp.Diagnostics.AddError("Attribute 'include_instance_name_label' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_instance_name_label', the 'type' attribute must be one of ['prometheus-monitoring']")
	}
	if internaltypes.IsDefined(model.AlwaysUsePermissiveModify) && model.Type.ValueString() != "directory-rest-api" {
		resp.Diagnostics.AddError("Attribute 'always_use_permissive_modify' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'always_use_permissive_modify', the 'type' attribute must be one of ['directory-rest-api']")
	}
	if internaltypes.IsDefined(model.StaticCustomDirectory) && model.Type.ValueString() != "velocity" {
		resp.Diagnostics.AddError("Attribute 'static_custom_directory' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'static_custom_directory', the 'type' attribute must be one of ['velocity']")
	}
	if internaltypes.IsDefined(model.OverrideStatusCode) && model.Type.ValueString() != "availability-state" {
		resp.Diagnostics.AddError("Attribute 'override_status_code' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'override_status_code', the 'type' attribute must be one of ['availability-state']")
	}
	if internaltypes.IsDefined(model.DebugType) && model.Type.ValueString() != "scim2" && model.Type.ValueString() != "ldap-mapped-scim" {
		resp.Diagnostics.AddError("Attribute 'debug_type' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'debug_type', the 'type' attribute must be one of ['scim2', 'ldap-mapped-scim']")
	}
	if internaltypes.IsDefined(model.AlwaysIncludeMonitorEntryNameLabel) && model.Type.ValueString() != "prometheus-monitoring" {
		resp.Diagnostics.AddError("Attribute 'always_include_monitor_entry_name_label' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'always_include_monitor_entry_name_label', the 'type' attribute must be one of ['prometheus-monitoring']")
	}
	if internaltypes.IsDefined(model.EntityTagLDAPAttribute) && model.Type.ValueString() != "ldap-mapped-scim" {
		resp.Diagnostics.AddError("Attribute 'entity_tag_ldap_attribute' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'entity_tag_ldap_attribute', the 'type' attribute must be one of ['ldap-mapped-scim']")
	}
	if internaltypes.IsDefined(model.CharacterEncoding) && model.Type.ValueString() != "velocity" {
		resp.Diagnostics.AddError("Attribute 'character_encoding' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'character_encoding', the 'type' attribute must be one of ['velocity']")
	}
	if internaltypes.IsDefined(model.DefaultOperationalAttribute) && model.Type.ValueString() != "directory-rest-api" {
		resp.Diagnostics.AddError("Attribute 'default_operational_attribute' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'default_operational_attribute', the 'type' attribute must be one of ['directory-rest-api']")
	}
	if internaltypes.IsDefined(model.AllowedControl) && model.Type.ValueString() != "directory-rest-api" {
		resp.Diagnostics.AddError("Attribute 'allowed_control' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'allowed_control', the 'type' attribute must be one of ['directory-rest-api']")
	}
	if internaltypes.IsDefined(model.TemplateDirectory) && model.Type.ValueString() != "velocity" {
		resp.Diagnostics.AddError("Attribute 'template_directory' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'template_directory', the 'type' attribute must be one of ['velocity']")
	}
	if internaltypes.IsDefined(model.IdentityMapper) && model.Type.ValueString() != "delegated-admin" && model.Type.ValueString() != "file-server" && model.Type.ValueString() != "velocity" && model.Type.ValueString() != "consent" && model.Type.ValueString() != "config" && model.Type.ValueString() != "directory-rest-api" && model.Type.ValueString() != "ldap-mapped-scim" {
		resp.Diagnostics.AddError("Attribute 'identity_mapper' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'identity_mapper', the 'type' attribute must be one of ['delegated-admin', 'file-server', 'velocity', 'consent', 'config', 'directory-rest-api', 'ldap-mapped-scim']")
	}
	if internaltypes.IsDefined(model.DefaultMIMEType) && model.Type.ValueString() != "file-server" && model.Type.ValueString() != "velocity" {
		resp.Diagnostics.AddError("Attribute 'default_mime_type' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'default_mime_type', the 'type' attribute must be one of ['file-server', 'velocity']")
	}
	if internaltypes.IsDefined(model.StaticResponseHeader) && model.Type.ValueString() != "velocity" {
		resp.Diagnostics.AddError("Attribute 'static_response_header' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'static_response_header', the 'type' attribute must be one of ['velocity']")
	}
	if internaltypes.IsDefined(model.ExcludeLDAPObjectclass) && model.Type.ValueString() != "ldap-mapped-scim" {
		resp.Diagnostics.AddError("Attribute 'exclude_ldap_objectclass' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'exclude_ldap_objectclass', the 'type' attribute must be one of ['ldap-mapped-scim']")
	}
	if internaltypes.IsDefined(model.DebugLevel) && model.Type.ValueString() != "scim2" && model.Type.ValueString() != "ldap-mapped-scim" {
		resp.Diagnostics.AddError("Attribute 'debug_level' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'debug_level', the 'type' attribute must be one of ['scim2', 'ldap-mapped-scim']")
	}
	if internaltypes.IsDefined(model.IncludeLDAPObjectclass) && model.Type.ValueString() != "ldap-mapped-scim" {
		resp.Diagnostics.AddError("Attribute 'include_ldap_objectclass' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_ldap_objectclass', the 'type' attribute must be one of ['ldap-mapped-scim']")
	}
	if internaltypes.IsDefined(model.ExcludeLDAPBaseDN) && model.Type.ValueString() != "ldap-mapped-scim" {
		resp.Diagnostics.AddError("Attribute 'exclude_ldap_base_dn' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'exclude_ldap_base_dn', the 'type' attribute must be one of ['ldap-mapped-scim']")
	}
	if internaltypes.IsDefined(model.BulkMaxOperations) && model.Type.ValueString() != "ldap-mapped-scim" {
		resp.Diagnostics.AddError("Attribute 'bulk_max_operations' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'bulk_max_operations', the 'type' attribute must be one of ['ldap-mapped-scim']")
	}
	if internaltypes.IsDefined(model.MaxPageSize) && model.Type.ValueString() != "directory-rest-api" {
		resp.Diagnostics.AddError("Attribute 'max_page_size' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'max_page_size', the 'type' attribute must be one of ['directory-rest-api']")
	}
	if internaltypes.IsDefined(model.UnavailableStatusCode) && model.Type.ValueString() != "availability-state" {
		resp.Diagnostics.AddError("Attribute 'unavailable_status_code' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'unavailable_status_code', the 'type' attribute must be one of ['availability-state']")
	}
	if internaltypes.IsDefined(model.OAuthTokenHandler) && model.Type.ValueString() != "ldap-mapped-scim" {
		resp.Diagnostics.AddError("Attribute 'oauth_token_handler' not supported by pingdirectory_http_servlet_extension resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'oauth_token_handler', the 'type' attribute must be one of ['ldap-mapped-scim']")
	}
	compare, err := version.Compare(providerConfig.ProductVersion, version.PingDirectory9300)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	if internaltypes.IsDefined(model.AlwaysUsePermissiveModify) {
		resp.Diagnostics.AddError("Attribute 'always_use_permissive_modify' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
	compare, err = version.Compare(providerConfig.ProductVersion, version.PingDirectory9200)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	if internaltypes.IsDefined(model.Type) && model.Type.ValueString() == "prometheus-monitoring" {
		version.CheckResourceSupported(&resp.Diagnostics, version.PingDirectory9200,
			providerConfig.ProductVersion, resourceName+" with type \"prometheus_monitoring\"")
	}
}

// Add optional fields to create request for quickstart http-servlet-extension
func addOptionalQuickstartHttpServletExtensionFields(ctx context.Context, addRequest *client.AddQuickstartHttpServletExtensionRequest, plan httpServletExtensionResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Server) {
		addRequest.Server = plan.Server.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CrossOriginPolicy) {
		addRequest.CrossOriginPolicy = plan.CrossOriginPolicy.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.ResponseHeader) {
		var slice []string
		plan.ResponseHeader.ElementsAs(ctx, &slice, false)
		addRequest.ResponseHeader = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CorrelationIDResponseHeader) {
		addRequest.CorrelationIDResponseHeader = plan.CorrelationIDResponseHeader.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for availability-state http-servlet-extension
func addOptionalAvailabilityStateHttpServletExtensionFields(ctx context.Context, addRequest *client.AddAvailabilityStateHttpServletExtensionRequest, plan httpServletExtensionResourceModel) error {
	if internaltypes.IsDefined(plan.OverrideStatusCode) {
		addRequest.OverrideStatusCode = plan.OverrideStatusCode.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.IncludeResponseBody) {
		addRequest.IncludeResponseBody = plan.IncludeResponseBody.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AdditionalResponseContents) {
		addRequest.AdditionalResponseContents = plan.AdditionalResponseContents.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CrossOriginPolicy) {
		addRequest.CrossOriginPolicy = plan.CrossOriginPolicy.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.ResponseHeader) {
		var slice []string
		plan.ResponseHeader.ElementsAs(ctx, &slice, false)
		addRequest.ResponseHeader = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CorrelationIDResponseHeader) {
		addRequest.CorrelationIDResponseHeader = plan.CorrelationIDResponseHeader.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for prometheus-monitoring http-servlet-extension
func addOptionalPrometheusMonitoringHttpServletExtensionFields(ctx context.Context, addRequest *client.AddPrometheusMonitoringHttpServletExtensionRequest, plan httpServletExtensionResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BaseContextPath) {
		addRequest.BaseContextPath = plan.BaseContextPath.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludeInstanceNameLabel) {
		addRequest.IncludeInstanceNameLabel = plan.IncludeInstanceNameLabel.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeProductNameLabel) {
		addRequest.IncludeProductNameLabel = plan.IncludeProductNameLabel.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeLocationNameLabel) {
		addRequest.IncludeLocationNameLabel = plan.IncludeLocationNameLabel.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlwaysIncludeMonitorEntryNameLabel) {
		addRequest.AlwaysIncludeMonitorEntryNameLabel = plan.AlwaysIncludeMonitorEntryNameLabel.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeMonitorObjectClassNameLabel) {
		addRequest.IncludeMonitorObjectClassNameLabel = plan.IncludeMonitorObjectClassNameLabel.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeMonitorAttributeNameLabel) {
		addRequest.IncludeMonitorAttributeNameLabel = plan.IncludeMonitorAttributeNameLabel.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LabelNameValuePair) {
		var slice []string
		plan.LabelNameValuePair.ElementsAs(ctx, &slice, false)
		addRequest.LabelNameValuePair = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CrossOriginPolicy) {
		addRequest.CrossOriginPolicy = plan.CrossOriginPolicy.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.ResponseHeader) {
		var slice []string
		plan.ResponseHeader.ElementsAs(ctx, &slice, false)
		addRequest.ResponseHeader = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CorrelationIDResponseHeader) {
		addRequest.CorrelationIDResponseHeader = plan.CorrelationIDResponseHeader.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for ldap-mapped-scim http-servlet-extension
func addOptionalLdapMappedScimHttpServletExtensionFields(ctx context.Context, addRequest *client.AddLdapMappedScimHttpServletExtensionRequest, plan httpServletExtensionResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OAuthTokenHandler) {
		addRequest.OAuthTokenHandler = plan.OAuthTokenHandler.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.BasicAuthEnabled) {
		addRequest.BasicAuthEnabled = plan.BasicAuthEnabled.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IdentityMapper) {
		addRequest.IdentityMapper = plan.IdentityMapper.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ResourceMappingFile) {
		addRequest.ResourceMappingFile = plan.ResourceMappingFile.ValueStringPointer()
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
		addRequest.EntityTagLDAPAttribute = plan.EntityTagLDAPAttribute.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BaseContextPath) {
		addRequest.BaseContextPath = plan.BaseContextPath.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TemporaryDirectory) {
		addRequest.TemporaryDirectory = plan.TemporaryDirectory.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TemporaryDirectoryPermissions) {
		addRequest.TemporaryDirectoryPermissions = plan.TemporaryDirectoryPermissions.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.MaxResults) {
		addRequest.MaxResults = plan.MaxResults.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.BulkMaxOperations) {
		addRequest.BulkMaxOperations = plan.BulkMaxOperations.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BulkMaxPayloadSize) {
		addRequest.BulkMaxPayloadSize = plan.BulkMaxPayloadSize.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.BulkMaxConcurrentRequests) {
		addRequest.BulkMaxConcurrentRequests = plan.BulkMaxConcurrentRequests.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.DebugEnabled) {
		addRequest.DebugEnabled = plan.DebugEnabled.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DebugLevel) {
		debugLevel, err := client.NewEnumhttpServletExtensionDebugLevelPropFromValue(plan.DebugLevel.ValueString())
		if err != nil {
			return err
		}
		addRequest.DebugLevel = debugLevel
	}
	if internaltypes.IsDefined(plan.DebugType) {
		var slice []string
		plan.DebugType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumhttpServletExtensionDebugTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumhttpServletExtensionDebugTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DebugType = enumSlice
	}
	if internaltypes.IsDefined(plan.IncludeStackTrace) {
		addRequest.IncludeStackTrace = plan.IncludeStackTrace.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CrossOriginPolicy) {
		addRequest.CrossOriginPolicy = plan.CrossOriginPolicy.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.ResponseHeader) {
		var slice []string
		plan.ResponseHeader.ElementsAs(ctx, &slice, false)
		addRequest.ResponseHeader = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CorrelationIDResponseHeader) {
		addRequest.CorrelationIDResponseHeader = plan.CorrelationIDResponseHeader.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for groovy-scripted http-servlet-extension
func addOptionalGroovyScriptedHttpServletExtensionFields(ctx context.Context, addRequest *client.AddGroovyScriptedHttpServletExtensionRequest, plan httpServletExtensionResourceModel) error {
	if internaltypes.IsDefined(plan.ScriptArgument) {
		var slice []string
		plan.ScriptArgument.ElementsAs(ctx, &slice, false)
		addRequest.ScriptArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CrossOriginPolicy) {
		addRequest.CrossOriginPolicy = plan.CrossOriginPolicy.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.ResponseHeader) {
		var slice []string
		plan.ResponseHeader.ElementsAs(ctx, &slice, false)
		addRequest.ResponseHeader = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CorrelationIDResponseHeader) {
		addRequest.CorrelationIDResponseHeader = plan.CorrelationIDResponseHeader.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for file-server http-servlet-extension
func addOptionalFileServerHttpServletExtensionFields(ctx context.Context, addRequest *client.AddFileServerHttpServletExtensionRequest, plan httpServletExtensionResourceModel) error {
	if internaltypes.IsDefined(plan.EnableDirectoryIndexing) {
		addRequest.EnableDirectoryIndexing = plan.EnableDirectoryIndexing.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IndexFile) {
		var slice []string
		plan.IndexFile.ElementsAs(ctx, &slice, false)
		addRequest.IndexFile = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MimeTypesFile) {
		addRequest.MimeTypesFile = plan.MimeTypesFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DefaultMIMEType) {
		addRequest.DefaultMIMEType = plan.DefaultMIMEType.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RequireAuthentication) {
		addRequest.RequireAuthentication = plan.RequireAuthentication.ValueBoolPointer()
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
		addRequest.RequireFileServletAccessPrivilege = plan.RequireFileServletAccessPrivilege.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.RequireGroup) {
		var slice []string
		plan.RequireGroup.ElementsAs(ctx, &slice, false)
		addRequest.RequireGroup = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IdentityMapper) {
		addRequest.IdentityMapper = plan.IdentityMapper.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CrossOriginPolicy) {
		addRequest.CrossOriginPolicy = plan.CrossOriginPolicy.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.ResponseHeader) {
		var slice []string
		plan.ResponseHeader.ElementsAs(ctx, &slice, false)
		addRequest.ResponseHeader = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CorrelationIDResponseHeader) {
		addRequest.CorrelationIDResponseHeader = plan.CorrelationIDResponseHeader.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for third-party http-servlet-extension
func addOptionalThirdPartyHttpServletExtensionFields(ctx context.Context, addRequest *client.AddThirdPartyHttpServletExtensionRequest, plan httpServletExtensionResourceModel) error {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CrossOriginPolicy) {
		addRequest.CrossOriginPolicy = plan.CrossOriginPolicy.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.ResponseHeader) {
		var slice []string
		plan.ResponseHeader.ElementsAs(ctx, &slice, false)
		addRequest.ResponseHeader = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CorrelationIDResponseHeader) {
		addRequest.CorrelationIDResponseHeader = plan.CorrelationIDResponseHeader.ValueStringPointer()
	}
	return nil
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateHttpServletExtensionUnknownValues(ctx context.Context, model *httpServletExtensionResourceModel) {
	if model.ExcludeLDAPObjectclass.ElementType(ctx) == nil {
		model.ExcludeLDAPObjectclass = types.SetNull(types.StringType)
	}
	if model.ExcludeLDAPBaseDN.ElementType(ctx) == nil {
		model.ExcludeLDAPBaseDN = types.SetNull(types.StringType)
	}
	if model.IndexFile.ElementType(ctx) == nil {
		model.IndexFile = types.SetNull(types.StringType)
	}
	if model.LabelNameValuePair.ElementType(ctx) == nil {
		model.LabelNameValuePair = types.SetNull(types.StringType)
	}
	if model.ExtensionArgument.ElementType(ctx) == nil {
		model.ExtensionArgument = types.SetNull(types.StringType)
	}
	if model.DebugType.ElementType(ctx) == nil {
		model.DebugType = types.SetNull(types.StringType)
	}
	if model.IdTokenValidator.ElementType(ctx) == nil {
		model.IdTokenValidator = types.SetNull(types.StringType)
	}
	if model.AccessTokenValidator.ElementType(ctx) == nil {
		model.AccessTokenValidator = types.SetNull(types.StringType)
	}
	if model.ScriptArgument.ElementType(ctx) == nil {
		model.ScriptArgument = types.SetNull(types.StringType)
	}
	if model.IncludeLDAPObjectclass.ElementType(ctx) == nil {
		model.IncludeLDAPObjectclass = types.SetNull(types.StringType)
	}
	if model.RequireGroup.ElementType(ctx) == nil {
		model.RequireGroup = types.SetNull(types.StringType)
	}
	if model.AllowedAuthenticationType.ElementType(ctx) == nil {
		model.AllowedAuthenticationType = types.SetNull(types.StringType)
	}
	if model.IncludeLDAPBaseDN.ElementType(ctx) == nil {
		model.IncludeLDAPBaseDN = types.SetNull(types.StringType)
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateHttpServletExtensionUnknownValuesDefault(ctx context.Context, model *defaultHttpServletExtensionResourceModel) {
	if model.ExcludeLDAPObjectclass.ElementType(ctx) == nil {
		model.ExcludeLDAPObjectclass = types.SetNull(types.StringType)
	}
	if model.ExcludeLDAPBaseDN.ElementType(ctx) == nil {
		model.ExcludeLDAPBaseDN = types.SetNull(types.StringType)
	}
	if model.IndexFile.ElementType(ctx) == nil {
		model.IndexFile = types.SetNull(types.StringType)
	}
	if model.SchemasEndpointObjectclass.ElementType(ctx) == nil {
		model.SchemasEndpointObjectclass = types.SetNull(types.StringType)
	}
	if model.LabelNameValuePair.ElementType(ctx) == nil {
		model.LabelNameValuePair = types.SetNull(types.StringType)
	}
	if model.ExtensionArgument.ElementType(ctx) == nil {
		model.ExtensionArgument = types.SetNull(types.StringType)
	}
	if model.DebugType.ElementType(ctx) == nil {
		model.DebugType = types.SetNull(types.StringType)
	}
	if model.IdTokenValidator.ElementType(ctx) == nil {
		model.IdTokenValidator = types.SetNull(types.StringType)
	}
	if model.RejectExpansionAttribute.ElementType(ctx) == nil {
		model.RejectExpansionAttribute = types.SetNull(types.StringType)
	}
	if model.AllowedControl.ElementType(ctx) == nil {
		model.AllowedControl = types.SetNull(types.StringType)
	}
	if model.AccessTokenValidator.ElementType(ctx) == nil {
		model.AccessTokenValidator = types.SetNull(types.StringType)
	}
	if model.DefaultOperationalAttribute.ElementType(ctx) == nil {
		model.DefaultOperationalAttribute = types.SetNull(types.StringType)
	}
	if model.TemplateDirectory.ElementType(ctx) == nil {
		model.TemplateDirectory = types.SetNull(types.StringType)
	}
	if model.ScriptArgument.ElementType(ctx) == nil {
		model.ScriptArgument = types.SetNull(types.StringType)
	}
	if model.IncludeLDAPObjectclass.ElementType(ctx) == nil {
		model.IncludeLDAPObjectclass = types.SetNull(types.StringType)
	}
	if model.RequireGroup.ElementType(ctx) == nil {
		model.RequireGroup = types.SetNull(types.StringType)
	}
	if model.AllowedAuthenticationType.ElementType(ctx) == nil {
		model.AllowedAuthenticationType = types.SetNull(types.StringType)
	}
	if model.IncludeLDAPBaseDN.ElementType(ctx) == nil {
		model.IncludeLDAPBaseDN = types.SetNull(types.StringType)
	}
	if model.StaticResponseHeader.ElementType(ctx) == nil {
		model.StaticResponseHeader = types.SetNull(types.StringType)
	}
}

// Read a DelegatedAdminHttpServletExtensionResponse object into the model struct
func readDelegatedAdminHttpServletExtensionResponseDefault(ctx context.Context, r *client.DelegatedAdminHttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("delegated-admin")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.BasicAuthEnabled = internaltypes.BoolTypeOrNil(r.BasicAuthEnabled)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, internaltypes.IsEmptyString(expectedValues.IdentityMapper))
	state.AccessTokenValidator = internaltypes.GetStringSet(r.AccessTokenValidator)
	state.AccessTokenScope = internaltypes.StringTypeOrNil(r.AccessTokenScope, internaltypes.IsEmptyString(expectedValues.AccessTokenScope))
	state.Audience = internaltypes.StringTypeOrNil(r.Audience, internaltypes.IsEmptyString(expectedValues.Audience))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, internaltypes.IsEmptyString(expectedValues.CrossOriginPolicy))
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, internaltypes.IsEmptyString(expectedValues.CorrelationIDResponseHeader))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(ctx, state)
}

// Read a QuickstartHttpServletExtensionResponse object into the model struct
func readQuickstartHttpServletExtensionResponse(ctx context.Context, r *client.QuickstartHttpServletExtensionResponse, state *httpServletExtensionResourceModel, expectedValues *httpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("quickstart")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Server = internaltypes.StringTypeOrNil(r.Server, internaltypes.IsEmptyString(expectedValues.Server))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, internaltypes.IsEmptyString(expectedValues.CrossOriginPolicy))
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, internaltypes.IsEmptyString(expectedValues.CorrelationIDResponseHeader))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValues(ctx, state)
}

// Read a QuickstartHttpServletExtensionResponse object into the model struct
func readQuickstartHttpServletExtensionResponseDefault(ctx context.Context, r *client.QuickstartHttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("quickstart")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Server = internaltypes.StringTypeOrNil(r.Server, internaltypes.IsEmptyString(expectedValues.Server))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, internaltypes.IsEmptyString(expectedValues.CrossOriginPolicy))
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, internaltypes.IsEmptyString(expectedValues.CorrelationIDResponseHeader))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(ctx, state)
}

// Read a AvailabilityStateHttpServletExtensionResponse object into the model struct
func readAvailabilityStateHttpServletExtensionResponse(ctx context.Context, r *client.AvailabilityStateHttpServletExtensionResponse, state *httpServletExtensionResourceModel, expectedValues *httpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("availability-state")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.AvailableStatusCode = types.Int64Value(r.AvailableStatusCode)
	state.DegradedStatusCode = types.Int64Value(r.DegradedStatusCode)
	state.UnavailableStatusCode = types.Int64Value(r.UnavailableStatusCode)
	state.OverrideStatusCode = internaltypes.Int64TypeOrNil(r.OverrideStatusCode)
	state.IncludeResponseBody = internaltypes.BoolTypeOrNil(r.IncludeResponseBody)
	state.AdditionalResponseContents = internaltypes.StringTypeOrNil(r.AdditionalResponseContents, internaltypes.IsEmptyString(expectedValues.AdditionalResponseContents))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, internaltypes.IsEmptyString(expectedValues.CrossOriginPolicy))
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, internaltypes.IsEmptyString(expectedValues.CorrelationIDResponseHeader))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValues(ctx, state)
}

// Read a AvailabilityStateHttpServletExtensionResponse object into the model struct
func readAvailabilityStateHttpServletExtensionResponseDefault(ctx context.Context, r *client.AvailabilityStateHttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("availability-state")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.AvailableStatusCode = types.Int64Value(r.AvailableStatusCode)
	state.DegradedStatusCode = types.Int64Value(r.DegradedStatusCode)
	state.UnavailableStatusCode = types.Int64Value(r.UnavailableStatusCode)
	state.OverrideStatusCode = internaltypes.Int64TypeOrNil(r.OverrideStatusCode)
	state.IncludeResponseBody = internaltypes.BoolTypeOrNil(r.IncludeResponseBody)
	state.AdditionalResponseContents = internaltypes.StringTypeOrNil(r.AdditionalResponseContents, internaltypes.IsEmptyString(expectedValues.AdditionalResponseContents))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, internaltypes.IsEmptyString(expectedValues.CrossOriginPolicy))
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, internaltypes.IsEmptyString(expectedValues.CorrelationIDResponseHeader))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(ctx, state)
}

// Read a PrometheusMonitoringHttpServletExtensionResponse object into the model struct
func readPrometheusMonitoringHttpServletExtensionResponse(ctx context.Context, r *client.PrometheusMonitoringHttpServletExtensionResponse, state *httpServletExtensionResourceModel, expectedValues *httpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("prometheus-monitoring")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.IncludeInstanceNameLabel = internaltypes.BoolTypeOrNil(r.IncludeInstanceNameLabel)
	state.IncludeProductNameLabel = internaltypes.BoolTypeOrNil(r.IncludeProductNameLabel)
	state.IncludeLocationNameLabel = internaltypes.BoolTypeOrNil(r.IncludeLocationNameLabel)
	state.AlwaysIncludeMonitorEntryNameLabel = internaltypes.BoolTypeOrNil(r.AlwaysIncludeMonitorEntryNameLabel)
	state.IncludeMonitorObjectClassNameLabel = internaltypes.BoolTypeOrNil(r.IncludeMonitorObjectClassNameLabel)
	state.IncludeMonitorAttributeNameLabel = internaltypes.BoolTypeOrNil(r.IncludeMonitorAttributeNameLabel)
	state.LabelNameValuePair = internaltypes.GetStringSet(r.LabelNameValuePair)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, internaltypes.IsEmptyString(expectedValues.CrossOriginPolicy))
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, internaltypes.IsEmptyString(expectedValues.CorrelationIDResponseHeader))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValues(ctx, state)
}

// Read a PrometheusMonitoringHttpServletExtensionResponse object into the model struct
func readPrometheusMonitoringHttpServletExtensionResponseDefault(ctx context.Context, r *client.PrometheusMonitoringHttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("prometheus-monitoring")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.IncludeInstanceNameLabel = internaltypes.BoolTypeOrNil(r.IncludeInstanceNameLabel)
	state.IncludeProductNameLabel = internaltypes.BoolTypeOrNil(r.IncludeProductNameLabel)
	state.IncludeLocationNameLabel = internaltypes.BoolTypeOrNil(r.IncludeLocationNameLabel)
	state.AlwaysIncludeMonitorEntryNameLabel = internaltypes.BoolTypeOrNil(r.AlwaysIncludeMonitorEntryNameLabel)
	state.IncludeMonitorObjectClassNameLabel = internaltypes.BoolTypeOrNil(r.IncludeMonitorObjectClassNameLabel)
	state.IncludeMonitorAttributeNameLabel = internaltypes.BoolTypeOrNil(r.IncludeMonitorAttributeNameLabel)
	state.LabelNameValuePair = internaltypes.GetStringSet(r.LabelNameValuePair)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, internaltypes.IsEmptyString(expectedValues.CrossOriginPolicy))
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, internaltypes.IsEmptyString(expectedValues.CorrelationIDResponseHeader))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(ctx, state)
}

// Read a VelocityHttpServletExtensionResponse object into the model struct
func readVelocityHttpServletExtensionResponseDefault(ctx context.Context, r *client.VelocityHttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("velocity")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.StaticContextPath = internaltypes.StringTypeOrNil(r.StaticContextPath, internaltypes.IsEmptyString(expectedValues.StaticContextPath))
	state.StaticContentDirectory = internaltypes.StringTypeOrNil(r.StaticContentDirectory, internaltypes.IsEmptyString(expectedValues.StaticContentDirectory))
	state.StaticCustomDirectory = internaltypes.StringTypeOrNil(r.StaticCustomDirectory, internaltypes.IsEmptyString(expectedValues.StaticCustomDirectory))
	state.TemplateDirectory = internaltypes.GetStringSet(r.TemplateDirectory)
	state.ExposeRequestAttributes = internaltypes.BoolTypeOrNil(r.ExposeRequestAttributes)
	state.ExposeSessionAttributes = internaltypes.BoolTypeOrNil(r.ExposeSessionAttributes)
	state.ExposeServerContext = internaltypes.BoolTypeOrNil(r.ExposeServerContext)
	state.AllowContextOverride = internaltypes.BoolTypeOrNil(r.AllowContextOverride)
	state.MimeTypesFile = internaltypes.StringTypeOrNil(r.MimeTypesFile, internaltypes.IsEmptyString(expectedValues.MimeTypesFile))
	state.DefaultMIMEType = internaltypes.StringTypeOrNil(r.DefaultMIMEType, internaltypes.IsEmptyString(expectedValues.DefaultMIMEType))
	state.CharacterEncoding = internaltypes.StringTypeOrNil(r.CharacterEncoding, internaltypes.IsEmptyString(expectedValues.CharacterEncoding))
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.StaticResponseHeader = internaltypes.GetStringSet(r.StaticResponseHeader)
	state.RequireAuthentication = internaltypes.BoolTypeOrNil(r.RequireAuthentication)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, internaltypes.IsEmptyString(expectedValues.IdentityMapper))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, internaltypes.IsEmptyString(expectedValues.CrossOriginPolicy))
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, internaltypes.IsEmptyString(expectedValues.CorrelationIDResponseHeader))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(ctx, state)
}

// Read a ConsentHttpServletExtensionResponse object into the model struct
func readConsentHttpServletExtensionResponseDefault(ctx context.Context, r *client.ConsentHttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("consent")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.BearerTokenAuthEnabled = internaltypes.BoolTypeOrNil(r.BearerTokenAuthEnabled)
	state.BasicAuthEnabled = internaltypes.BoolTypeOrNil(r.BasicAuthEnabled)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, internaltypes.IsEmptyString(expectedValues.IdentityMapper))
	state.AccessTokenValidator = internaltypes.GetStringSet(r.AccessTokenValidator)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, internaltypes.IsEmptyString(expectedValues.CrossOriginPolicy))
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, internaltypes.IsEmptyString(expectedValues.CorrelationIDResponseHeader))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(ctx, state)
}

// Read a LdapMappedScimHttpServletExtensionResponse object into the model struct
func readLdapMappedScimHttpServletExtensionResponse(ctx context.Context, r *client.LdapMappedScimHttpServletExtensionResponse, state *httpServletExtensionResourceModel, expectedValues *httpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap-mapped-scim")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
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
	populateHttpServletExtensionUnknownValues(ctx, state)
}

// Read a LdapMappedScimHttpServletExtensionResponse object into the model struct
func readLdapMappedScimHttpServletExtensionResponseDefault(ctx context.Context, r *client.LdapMappedScimHttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap-mapped-scim")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
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
	populateHttpServletExtensionUnknownValuesDefault(ctx, state)
}

// Read a GroovyScriptedHttpServletExtensionResponse object into the model struct
func readGroovyScriptedHttpServletExtensionResponse(ctx context.Context, r *client.GroovyScriptedHttpServletExtensionResponse, state *httpServletExtensionResourceModel, expectedValues *httpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, internaltypes.IsEmptyString(expectedValues.CrossOriginPolicy))
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, internaltypes.IsEmptyString(expectedValues.CorrelationIDResponseHeader))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValues(ctx, state)
}

// Read a GroovyScriptedHttpServletExtensionResponse object into the model struct
func readGroovyScriptedHttpServletExtensionResponseDefault(ctx context.Context, r *client.GroovyScriptedHttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, internaltypes.IsEmptyString(expectedValues.CrossOriginPolicy))
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, internaltypes.IsEmptyString(expectedValues.CorrelationIDResponseHeader))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(ctx, state)
}

// Read a FileServerHttpServletExtensionResponse object into the model struct
func readFileServerHttpServletExtensionResponse(ctx context.Context, r *client.FileServerHttpServletExtensionResponse, state *httpServletExtensionResourceModel, expectedValues *httpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-server")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
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
	populateHttpServletExtensionUnknownValues(ctx, state)
}

// Read a FileServerHttpServletExtensionResponse object into the model struct
func readFileServerHttpServletExtensionResponseDefault(ctx context.Context, r *client.FileServerHttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-server")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
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
	populateHttpServletExtensionUnknownValuesDefault(ctx, state)
}

// Read a ConfigHttpServletExtensionResponse object into the model struct
func readConfigHttpServletExtensionResponseDefault(ctx context.Context, r *client.ConfigHttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("config")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, internaltypes.IsEmptyString(expectedValues.IdentityMapper))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, internaltypes.IsEmptyString(expectedValues.CrossOriginPolicy))
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, internaltypes.IsEmptyString(expectedValues.CorrelationIDResponseHeader))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(ctx, state)
}

// Read a Scim2HttpServletExtensionResponse object into the model struct
func readScim2HttpServletExtensionResponseDefault(ctx context.Context, r *client.Scim2HttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("scim2")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.AccessTokenValidator = internaltypes.GetStringSet(r.AccessTokenValidator)
	state.MapAccessTokensToLocalUsers = internaltypes.StringTypeOrNil(
		client.StringPointerEnumhttpServletExtensionMapAccessTokensToLocalUsersProp(r.MapAccessTokensToLocalUsers), internaltypes.IsEmptyString(expectedValues.MapAccessTokensToLocalUsers))
	state.DebugEnabled = internaltypes.BoolTypeOrNil(r.DebugEnabled)
	state.DebugLevel = types.StringValue(r.DebugLevel.String())
	state.DebugType = internaltypes.GetStringSet(
		client.StringSliceEnumhttpServletExtensionDebugTypeProp(r.DebugType))
	state.IncludeStackTrace = types.BoolValue(r.IncludeStackTrace)
	state.SwaggerEnabled = internaltypes.BoolTypeOrNil(r.SwaggerEnabled)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, internaltypes.IsEmptyString(expectedValues.CrossOriginPolicy))
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, internaltypes.IsEmptyString(expectedValues.CorrelationIDResponseHeader))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(ctx, state)
}

// Read a DirectoryRestApiHttpServletExtensionResponse object into the model struct
func readDirectoryRestApiHttpServletExtensionResponseDefault(ctx context.Context, r *client.DirectoryRestApiHttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("directory-rest-api")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.BasicAuthEnabled = internaltypes.BoolTypeOrNil(r.BasicAuthEnabled)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, internaltypes.IsEmptyString(expectedValues.IdentityMapper))
	state.AccessTokenValidator = internaltypes.GetStringSet(r.AccessTokenValidator)
	state.AccessTokenScope = internaltypes.StringTypeOrNil(r.AccessTokenScope, internaltypes.IsEmptyString(expectedValues.AccessTokenScope))
	state.Audience = internaltypes.StringTypeOrNil(r.Audience, internaltypes.IsEmptyString(expectedValues.Audience))
	state.MaxPageSize = internaltypes.Int64TypeOrNil(r.MaxPageSize)
	state.SchemasEndpointObjectclass = internaltypes.GetStringSet(r.SchemasEndpointObjectclass)
	state.DefaultOperationalAttribute = internaltypes.GetStringSet(r.DefaultOperationalAttribute)
	state.RejectExpansionAttribute = internaltypes.GetStringSet(r.RejectExpansionAttribute)
	state.AlwaysUsePermissiveModify = internaltypes.BoolTypeOrNil(r.AlwaysUsePermissiveModify)
	state.AllowedControl = internaltypes.GetStringSet(
		client.StringSliceEnumhttpServletExtensionAllowedControlProp(r.AllowedControl))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, internaltypes.IsEmptyString(expectedValues.CrossOriginPolicy))
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, internaltypes.IsEmptyString(expectedValues.CorrelationIDResponseHeader))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(ctx, state)
}

// Read a ThirdPartyHttpServletExtensionResponse object into the model struct
func readThirdPartyHttpServletExtensionResponse(ctx context.Context, r *client.ThirdPartyHttpServletExtensionResponse, state *httpServletExtensionResourceModel, expectedValues *httpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, internaltypes.IsEmptyString(expectedValues.CrossOriginPolicy))
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, internaltypes.IsEmptyString(expectedValues.CorrelationIDResponseHeader))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValues(ctx, state)
}

// Read a ThirdPartyHttpServletExtensionResponse object into the model struct
func readThirdPartyHttpServletExtensionResponseDefault(ctx context.Context, r *client.ThirdPartyHttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, internaltypes.IsEmptyString(expectedValues.CrossOriginPolicy))
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, internaltypes.IsEmptyString(expectedValues.CorrelationIDResponseHeader))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(ctx, state)
}

// Create any update operations necessary to make the state match the plan
func createHttpServletExtensionOperations(plan httpServletExtensionResourceModel, state httpServletExtensionResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.ScriptClass, state.ScriptClass, "script-class")
	operations.AddStringOperationIfNecessary(&ops, plan.DocumentRootDirectory, state.DocumentRootDirectory, "document-root-directory")
	operations.AddBoolOperationIfNecessary(&ops, plan.EnableDirectoryIndexing, state.EnableDirectoryIndexing, "enable-directory-indexing")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IndexFile, state.IndexFile, "index-file")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScriptArgument, state.ScriptArgument, "script-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.OAuthTokenHandler, state.OAuthTokenHandler, "oauth-token-handler")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedAuthenticationType, state.AllowedAuthenticationType, "allowed-authentication-type")
	operations.AddStringOperationIfNecessary(&ops, plan.BaseContextPath, state.BaseContextPath, "base-context-path")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IdTokenValidator, state.IdTokenValidator, "id-token-validator")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequireFileServletAccessPrivilege, state.RequireFileServletAccessPrivilege, "require-file-servlet-access-privilege")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RequireGroup, state.RequireGroup, "require-group")
	operations.AddStringOperationIfNecessary(&ops, plan.ResourceMappingFile, state.ResourceMappingFile, "resource-mapping-file")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeLDAPObjectclass, state.IncludeLDAPObjectclass, "include-ldap-objectclass")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeLDAPObjectclass, state.ExcludeLDAPObjectclass, "exclude-ldap-objectclass")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeLDAPBaseDN, state.IncludeLDAPBaseDN, "include-ldap-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeLDAPBaseDN, state.ExcludeLDAPBaseDN, "exclude-ldap-base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.EntityTagLDAPAttribute, state.EntityTagLDAPAttribute, "entity-tag-ldap-attribute")
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
	operations.AddStringOperationIfNecessary(&ops, plan.MimeTypesFile, state.MimeTypesFile, "mime-types-file")
	operations.AddStringOperationIfNecessary(&ops, plan.DefaultMIMEType, state.DefaultMIMEType, "default-mime-type")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeInstanceNameLabel, state.IncludeInstanceNameLabel, "include-instance-name-label")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequireAuthentication, state.RequireAuthentication, "require-authentication")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeProductNameLabel, state.IncludeProductNameLabel, "include-product-name-label")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeLocationNameLabel, state.IncludeLocationNameLabel, "include-location-name-label")
	operations.AddBoolOperationIfNecessary(&ops, plan.AlwaysIncludeMonitorEntryNameLabel, state.AlwaysIncludeMonitorEntryNameLabel, "always-include-monitor-entry-name-label")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeMonitorObjectClassNameLabel, state.IncludeMonitorObjectClassNameLabel, "include-monitor-object-class-name-label")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeMonitorAttributeNameLabel, state.IncludeMonitorAttributeNameLabel, "include-monitor-attribute-name-label")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.LabelNameValuePair, state.LabelNameValuePair, "label-name-value-pair")
	operations.AddInt64OperationIfNecessary(&ops, plan.AvailableStatusCode, state.AvailableStatusCode, "available-status-code")
	operations.AddInt64OperationIfNecessary(&ops, plan.DegradedStatusCode, state.DegradedStatusCode, "degraded-status-code")
	operations.AddInt64OperationIfNecessary(&ops, plan.UnavailableStatusCode, state.UnavailableStatusCode, "unavailable-status-code")
	operations.AddInt64OperationIfNecessary(&ops, plan.OverrideStatusCode, state.OverrideStatusCode, "override-status-code")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeResponseBody, state.IncludeResponseBody, "include-response-body")
	operations.AddStringOperationIfNecessary(&ops, plan.AdditionalResponseContents, state.AdditionalResponseContents, "additional-response-contents")
	operations.AddStringOperationIfNecessary(&ops, plan.Server, state.Server, "server")
	operations.AddBoolOperationIfNecessary(&ops, plan.BasicAuthEnabled, state.BasicAuthEnabled, "basic-auth-enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.IdentityMapper, state.IdentityMapper, "identity-mapper")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AccessTokenValidator, state.AccessTokenValidator, "access-token-validator")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.CrossOriginPolicy, state.CrossOriginPolicy, "cross-origin-policy")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ResponseHeader, state.ResponseHeader, "response-header")
	operations.AddStringOperationIfNecessary(&ops, plan.CorrelationIDResponseHeader, state.CorrelationIDResponseHeader, "correlation-id-response-header")
	return ops
}

// Create any update operations necessary to make the state match the plan
func createHttpServletExtensionOperationsDefault(plan defaultHttpServletExtensionResourceModel, state defaultHttpServletExtensionResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.ScriptClass, state.ScriptClass, "script-class")
	operations.AddStringOperationIfNecessary(&ops, plan.DocumentRootDirectory, state.DocumentRootDirectory, "document-root-directory")
	operations.AddStringOperationIfNecessary(&ops, plan.MapAccessTokensToLocalUsers, state.MapAccessTokensToLocalUsers, "map-access-tokens-to-local-users")
	operations.AddBoolOperationIfNecessary(&ops, plan.EnableDirectoryIndexing, state.EnableDirectoryIndexing, "enable-directory-indexing")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IndexFile, state.IndexFile, "index-file")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxPageSize, state.MaxPageSize, "max-page-size")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SchemasEndpointObjectclass, state.SchemasEndpointObjectclass, "schemas-endpoint-objectclass")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DefaultOperationalAttribute, state.DefaultOperationalAttribute, "default-operational-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RejectExpansionAttribute, state.RejectExpansionAttribute, "reject-expansion-attribute")
	operations.AddBoolOperationIfNecessary(&ops, plan.AlwaysUsePermissiveModify, state.AlwaysUsePermissiveModify, "always-use-permissive-modify")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedControl, state.AllowedControl, "allowed-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScriptArgument, state.ScriptArgument, "script-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.OAuthTokenHandler, state.OAuthTokenHandler, "oauth-token-handler")
	operations.AddBoolOperationIfNecessary(&ops, plan.SwaggerEnabled, state.SwaggerEnabled, "swagger-enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.BearerTokenAuthEnabled, state.BearerTokenAuthEnabled, "bearer-token-auth-enabled")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedAuthenticationType, state.AllowedAuthenticationType, "allowed-authentication-type")
	operations.AddStringOperationIfNecessary(&ops, plan.BaseContextPath, state.BaseContextPath, "base-context-path")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IdTokenValidator, state.IdTokenValidator, "id-token-validator")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequireFileServletAccessPrivilege, state.RequireFileServletAccessPrivilege, "require-file-servlet-access-privilege")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RequireGroup, state.RequireGroup, "require-group")
	operations.AddStringOperationIfNecessary(&ops, plan.ResourceMappingFile, state.ResourceMappingFile, "resource-mapping-file")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeLDAPObjectclass, state.IncludeLDAPObjectclass, "include-ldap-objectclass")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeLDAPObjectclass, state.ExcludeLDAPObjectclass, "exclude-ldap-objectclass")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeLDAPBaseDN, state.IncludeLDAPBaseDN, "include-ldap-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeLDAPBaseDN, state.ExcludeLDAPBaseDN, "exclude-ldap-base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.EntityTagLDAPAttribute, state.EntityTagLDAPAttribute, "entity-tag-ldap-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.StaticContextPath, state.StaticContextPath, "static-context-path")
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
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeInstanceNameLabel, state.IncludeInstanceNameLabel, "include-instance-name-label")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.StaticResponseHeader, state.StaticResponseHeader, "static-response-header")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequireAuthentication, state.RequireAuthentication, "require-authentication")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeProductNameLabel, state.IncludeProductNameLabel, "include-product-name-label")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeLocationNameLabel, state.IncludeLocationNameLabel, "include-location-name-label")
	operations.AddBoolOperationIfNecessary(&ops, plan.AlwaysIncludeMonitorEntryNameLabel, state.AlwaysIncludeMonitorEntryNameLabel, "always-include-monitor-entry-name-label")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeMonitorObjectClassNameLabel, state.IncludeMonitorObjectClassNameLabel, "include-monitor-object-class-name-label")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeMonitorAttributeNameLabel, state.IncludeMonitorAttributeNameLabel, "include-monitor-attribute-name-label")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.LabelNameValuePair, state.LabelNameValuePair, "label-name-value-pair")
	operations.AddInt64OperationIfNecessary(&ops, plan.AvailableStatusCode, state.AvailableStatusCode, "available-status-code")
	operations.AddInt64OperationIfNecessary(&ops, plan.DegradedStatusCode, state.DegradedStatusCode, "degraded-status-code")
	operations.AddInt64OperationIfNecessary(&ops, plan.UnavailableStatusCode, state.UnavailableStatusCode, "unavailable-status-code")
	operations.AddInt64OperationIfNecessary(&ops, plan.OverrideStatusCode, state.OverrideStatusCode, "override-status-code")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeResponseBody, state.IncludeResponseBody, "include-response-body")
	operations.AddStringOperationIfNecessary(&ops, plan.AdditionalResponseContents, state.AdditionalResponseContents, "additional-response-contents")
	operations.AddStringOperationIfNecessary(&ops, plan.Server, state.Server, "server")
	operations.AddBoolOperationIfNecessary(&ops, plan.BasicAuthEnabled, state.BasicAuthEnabled, "basic-auth-enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.IdentityMapper, state.IdentityMapper, "identity-mapper")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AccessTokenValidator, state.AccessTokenValidator, "access-token-validator")
	operations.AddStringOperationIfNecessary(&ops, plan.AccessTokenScope, state.AccessTokenScope, "access-token-scope")
	operations.AddStringOperationIfNecessary(&ops, plan.Audience, state.Audience, "audience")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.CrossOriginPolicy, state.CrossOriginPolicy, "cross-origin-policy")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ResponseHeader, state.ResponseHeader, "response-header")
	operations.AddStringOperationIfNecessary(&ops, plan.CorrelationIDResponseHeader, state.CorrelationIDResponseHeader, "correlation-id-response-header")
	return ops
}

// Create a quickstart http-servlet-extension
func (r *httpServletExtensionResource) CreateQuickstartHttpServletExtension(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan httpServletExtensionResourceModel) (*httpServletExtensionResourceModel, error) {
	addRequest := client.NewAddQuickstartHttpServletExtensionRequest(plan.Name.ValueString(),
		[]client.EnumquickstartHttpServletExtensionSchemaUrn{client.ENUMQUICKSTARTHTTPSERVLETEXTENSIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0HTTP_SERVLET_EXTENSIONQUICKSTART})
	err := addOptionalQuickstartHttpServletExtensionFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Http Servlet Extension", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.HttpServletExtensionApi.AddHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddHttpServletExtensionRequest(
		client.AddQuickstartHttpServletExtensionRequestAsAddHttpServletExtensionRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.AddHttpServletExtensionExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Http Servlet Extension", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state httpServletExtensionResourceModel
	readQuickstartHttpServletExtensionResponse(ctx, addResponse.QuickstartHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a groovy-scripted http-servlet-extension
func (r *httpServletExtensionResource) CreateGroovyScriptedHttpServletExtension(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan httpServletExtensionResourceModel) (*httpServletExtensionResourceModel, error) {
	addRequest := client.NewAddGroovyScriptedHttpServletExtensionRequest(plan.Name.ValueString(),
		[]client.EnumgroovyScriptedHttpServletExtensionSchemaUrn{client.ENUMGROOVYSCRIPTEDHTTPSERVLETEXTENSIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0HTTP_SERVLET_EXTENSIONGROOVY_SCRIPTED},
		plan.ScriptClass.ValueString())
	err := addOptionalGroovyScriptedHttpServletExtensionFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Http Servlet Extension", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.HttpServletExtensionApi.AddHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddHttpServletExtensionRequest(
		client.AddGroovyScriptedHttpServletExtensionRequestAsAddHttpServletExtensionRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.AddHttpServletExtensionExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Http Servlet Extension", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state httpServletExtensionResourceModel
	readGroovyScriptedHttpServletExtensionResponse(ctx, addResponse.GroovyScriptedHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a availability-state http-servlet-extension
func (r *httpServletExtensionResource) CreateAvailabilityStateHttpServletExtension(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan httpServletExtensionResourceModel) (*httpServletExtensionResourceModel, error) {
	addRequest := client.NewAddAvailabilityStateHttpServletExtensionRequest(plan.Name.ValueString(),
		[]client.EnumavailabilityStateHttpServletExtensionSchemaUrn{client.ENUMAVAILABILITYSTATEHTTPSERVLETEXTENSIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0HTTP_SERVLET_EXTENSIONAVAILABILITY_STATE},
		plan.BaseContextPath.ValueString(),
		plan.AvailableStatusCode.ValueInt64(),
		plan.DegradedStatusCode.ValueInt64(),
		plan.UnavailableStatusCode.ValueInt64())
	err := addOptionalAvailabilityStateHttpServletExtensionFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Http Servlet Extension", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.HttpServletExtensionApi.AddHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddHttpServletExtensionRequest(
		client.AddAvailabilityStateHttpServletExtensionRequestAsAddHttpServletExtensionRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.AddHttpServletExtensionExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Http Servlet Extension", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state httpServletExtensionResourceModel
	readAvailabilityStateHttpServletExtensionResponse(ctx, addResponse.AvailabilityStateHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a prometheus-monitoring http-servlet-extension
func (r *httpServletExtensionResource) CreatePrometheusMonitoringHttpServletExtension(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan httpServletExtensionResourceModel) (*httpServletExtensionResourceModel, error) {
	addRequest := client.NewAddPrometheusMonitoringHttpServletExtensionRequest(plan.Name.ValueString(),
		[]client.EnumprometheusMonitoringHttpServletExtensionSchemaUrn{client.ENUMPROMETHEUSMONITORINGHTTPSERVLETEXTENSIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0HTTP_SERVLET_EXTENSIONPROMETHEUS_MONITORING})
	err := addOptionalPrometheusMonitoringHttpServletExtensionFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Http Servlet Extension", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.HttpServletExtensionApi.AddHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddHttpServletExtensionRequest(
		client.AddPrometheusMonitoringHttpServletExtensionRequestAsAddHttpServletExtensionRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.AddHttpServletExtensionExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Http Servlet Extension", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state httpServletExtensionResourceModel
	readPrometheusMonitoringHttpServletExtensionResponse(ctx, addResponse.PrometheusMonitoringHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a file-server http-servlet-extension
func (r *httpServletExtensionResource) CreateFileServerHttpServletExtension(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan httpServletExtensionResourceModel) (*httpServletExtensionResourceModel, error) {
	addRequest := client.NewAddFileServerHttpServletExtensionRequest(plan.Name.ValueString(),
		[]client.EnumfileServerHttpServletExtensionSchemaUrn{client.ENUMFILESERVERHTTPSERVLETEXTENSIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0HTTP_SERVLET_EXTENSIONFILE_SERVER},
		plan.BaseContextPath.ValueString(),
		plan.DocumentRootDirectory.ValueString())
	err := addOptionalFileServerHttpServletExtensionFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Http Servlet Extension", err.Error())
		return nil, err
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Http Servlet Extension", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state httpServletExtensionResourceModel
	readFileServerHttpServletExtensionResponse(ctx, addResponse.FileServerHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a ldap-mapped-scim http-servlet-extension
func (r *httpServletExtensionResource) CreateLdapMappedScimHttpServletExtension(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan httpServletExtensionResourceModel) (*httpServletExtensionResourceModel, error) {
	addRequest := client.NewAddLdapMappedScimHttpServletExtensionRequest(plan.Name.ValueString(),
		[]client.EnumldapMappedScimHttpServletExtensionSchemaUrn{client.ENUMLDAPMAPPEDSCIMHTTPSERVLETEXTENSIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0HTTP_SERVLET_EXTENSIONLDAP_MAPPED_SCIM})
	err := addOptionalLdapMappedScimHttpServletExtensionFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Http Servlet Extension", err.Error())
		return nil, err
	}
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Http Servlet Extension", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state httpServletExtensionResourceModel
	readLdapMappedScimHttpServletExtensionResponse(ctx, addResponse.LdapMappedScimHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party http-servlet-extension
func (r *httpServletExtensionResource) CreateThirdPartyHttpServletExtension(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan httpServletExtensionResourceModel) (*httpServletExtensionResourceModel, error) {
	addRequest := client.NewAddThirdPartyHttpServletExtensionRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyHttpServletExtensionSchemaUrn{client.ENUMTHIRDPARTYHTTPSERVLETEXTENSIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0HTTP_SERVLET_EXTENSIONTHIRD_PARTY},
		plan.ExtensionClass.ValueString())
	err := addOptionalThirdPartyHttpServletExtensionFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Http Servlet Extension", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.HttpServletExtensionApi.AddHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddHttpServletExtensionRequest(
		client.AddThirdPartyHttpServletExtensionRequestAsAddHttpServletExtensionRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.AddHttpServletExtensionExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Http Servlet Extension", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state httpServletExtensionResourceModel
	readThirdPartyHttpServletExtensionResponse(ctx, addResponse.ThirdPartyHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *httpServletExtensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan httpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *httpServletExtensionResourceModel
	var err error
	if plan.Type.ValueString() == "quickstart" {
		state, err = r.CreateQuickstartHttpServletExtension(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "groovy-scripted" {
		state, err = r.CreateGroovyScriptedHttpServletExtension(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "availability-state" {
		state, err = r.CreateAvailabilityStateHttpServletExtension(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "prometheus-monitoring" {
		state, err = r.CreatePrometheusMonitoringHttpServletExtension(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "file-server" {
		state, err = r.CreateFileServerHttpServletExtension(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "ldap-mapped-scim" {
		state, err = r.CreateLdapMappedScimHttpServletExtension(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyHttpServletExtension(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

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
func (r *defaultHttpServletExtensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan defaultHttpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.GetHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Http Servlet Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state defaultHttpServletExtensionResourceModel
	if plan.Type.ValueString() == "delegated-admin" {
		readDelegatedAdminHttpServletExtensionResponseDefault(ctx, readResponse.DelegatedAdminHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "quickstart" {
		readQuickstartHttpServletExtensionResponseDefault(ctx, readResponse.QuickstartHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "availability-state" {
		readAvailabilityStateHttpServletExtensionResponseDefault(ctx, readResponse.AvailabilityStateHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "prometheus-monitoring" {
		readPrometheusMonitoringHttpServletExtensionResponseDefault(ctx, readResponse.PrometheusMonitoringHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "velocity" {
		readVelocityHttpServletExtensionResponseDefault(ctx, readResponse.VelocityHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "consent" {
		readConsentHttpServletExtensionResponseDefault(ctx, readResponse.ConsentHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "ldap-mapped-scim" {
		readLdapMappedScimHttpServletExtensionResponseDefault(ctx, readResponse.LdapMappedScimHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "groovy-scripted" {
		readGroovyScriptedHttpServletExtensionResponseDefault(ctx, readResponse.GroovyScriptedHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "file-server" {
		readFileServerHttpServletExtensionResponseDefault(ctx, readResponse.FileServerHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "config" {
		readConfigHttpServletExtensionResponseDefault(ctx, readResponse.ConfigHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "scim2" {
		readScim2HttpServletExtensionResponseDefault(ctx, readResponse.Scim2HttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "directory-rest-api" {
		readDirectoryRestApiHttpServletExtensionResponseDefault(ctx, readResponse.DirectoryRestApiHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "third-party" {
		readThirdPartyHttpServletExtensionResponseDefault(ctx, readResponse.ThirdPartyHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtension(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createHttpServletExtensionOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Http Servlet Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "delegated-admin" {
			readDelegatedAdminHttpServletExtensionResponseDefault(ctx, updateResponse.DelegatedAdminHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "quickstart" {
			readQuickstartHttpServletExtensionResponseDefault(ctx, updateResponse.QuickstartHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "availability-state" {
			readAvailabilityStateHttpServletExtensionResponseDefault(ctx, updateResponse.AvailabilityStateHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "prometheus-monitoring" {
			readPrometheusMonitoringHttpServletExtensionResponseDefault(ctx, updateResponse.PrometheusMonitoringHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "velocity" {
			readVelocityHttpServletExtensionResponseDefault(ctx, updateResponse.VelocityHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "consent" {
			readConsentHttpServletExtensionResponseDefault(ctx, updateResponse.ConsentHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "ldap-mapped-scim" {
			readLdapMappedScimHttpServletExtensionResponseDefault(ctx, updateResponse.LdapMappedScimHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "groovy-scripted" {
			readGroovyScriptedHttpServletExtensionResponseDefault(ctx, updateResponse.GroovyScriptedHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "file-server" {
			readFileServerHttpServletExtensionResponseDefault(ctx, updateResponse.FileServerHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "config" {
			readConfigHttpServletExtensionResponseDefault(ctx, updateResponse.ConfigHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "scim2" {
			readScim2HttpServletExtensionResponseDefault(ctx, updateResponse.Scim2HttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "directory-rest-api" {
			readDirectoryRestApiHttpServletExtensionResponseDefault(ctx, updateResponse.DirectoryRestApiHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party" {
			readThirdPartyHttpServletExtensionResponseDefault(ctx, updateResponse.ThirdPartyHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *httpServletExtensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state httpServletExtensionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.GetHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
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
		readQuickstartHttpServletExtensionResponse(ctx, readResponse.QuickstartHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AvailabilityStateHttpServletExtensionResponse != nil {
		readAvailabilityStateHttpServletExtensionResponse(ctx, readResponse.AvailabilityStateHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PrometheusMonitoringHttpServletExtensionResponse != nil {
		readPrometheusMonitoringHttpServletExtensionResponse(ctx, readResponse.PrometheusMonitoringHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LdapMappedScimHttpServletExtensionResponse != nil {
		readLdapMappedScimHttpServletExtensionResponse(ctx, readResponse.LdapMappedScimHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedHttpServletExtensionResponse != nil {
		readGroovyScriptedHttpServletExtensionResponse(ctx, readResponse.GroovyScriptedHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileServerHttpServletExtensionResponse != nil {
		readFileServerHttpServletExtensionResponse(ctx, readResponse.FileServerHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyHttpServletExtensionResponse != nil {
		readThirdPartyHttpServletExtensionResponse(ctx, readResponse.ThirdPartyHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *defaultHttpServletExtensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state defaultHttpServletExtensionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.GetHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
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
	if readResponse.DelegatedAdminHttpServletExtensionResponse != nil {
		readDelegatedAdminHttpServletExtensionResponseDefault(ctx, readResponse.DelegatedAdminHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.VelocityHttpServletExtensionResponse != nil {
		readVelocityHttpServletExtensionResponseDefault(ctx, readResponse.VelocityHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ConsentHttpServletExtensionResponse != nil {
		readConsentHttpServletExtensionResponseDefault(ctx, readResponse.ConsentHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ConfigHttpServletExtensionResponse != nil {
		readConfigHttpServletExtensionResponseDefault(ctx, readResponse.ConfigHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Scim2HttpServletExtensionResponse != nil {
		readScim2HttpServletExtensionResponseDefault(ctx, readResponse.Scim2HttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DirectoryRestApiHttpServletExtensionResponse != nil {
		readDirectoryRestApiHttpServletExtensionResponseDefault(ctx, readResponse.DirectoryRestApiHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *httpServletExtensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan httpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state httpServletExtensionResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createHttpServletExtensionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Http Servlet Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "quickstart" {
			readQuickstartHttpServletExtensionResponse(ctx, updateResponse.QuickstartHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "availability-state" {
			readAvailabilityStateHttpServletExtensionResponse(ctx, updateResponse.AvailabilityStateHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "prometheus-monitoring" {
			readPrometheusMonitoringHttpServletExtensionResponse(ctx, updateResponse.PrometheusMonitoringHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "ldap-mapped-scim" {
			readLdapMappedScimHttpServletExtensionResponse(ctx, updateResponse.LdapMappedScimHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "groovy-scripted" {
			readGroovyScriptedHttpServletExtensionResponse(ctx, updateResponse.GroovyScriptedHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "file-server" {
			readFileServerHttpServletExtensionResponse(ctx, updateResponse.FileServerHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party" {
			readThirdPartyHttpServletExtensionResponse(ctx, updateResponse.ThirdPartyHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
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

func (r *defaultHttpServletExtensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan defaultHttpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state defaultHttpServletExtensionResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createHttpServletExtensionOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Http Servlet Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "delegated-admin" {
			readDelegatedAdminHttpServletExtensionResponseDefault(ctx, updateResponse.DelegatedAdminHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "quickstart" {
			readQuickstartHttpServletExtensionResponseDefault(ctx, updateResponse.QuickstartHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "availability-state" {
			readAvailabilityStateHttpServletExtensionResponseDefault(ctx, updateResponse.AvailabilityStateHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "prometheus-monitoring" {
			readPrometheusMonitoringHttpServletExtensionResponseDefault(ctx, updateResponse.PrometheusMonitoringHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "velocity" {
			readVelocityHttpServletExtensionResponseDefault(ctx, updateResponse.VelocityHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "consent" {
			readConsentHttpServletExtensionResponseDefault(ctx, updateResponse.ConsentHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "ldap-mapped-scim" {
			readLdapMappedScimHttpServletExtensionResponseDefault(ctx, updateResponse.LdapMappedScimHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "groovy-scripted" {
			readGroovyScriptedHttpServletExtensionResponseDefault(ctx, updateResponse.GroovyScriptedHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "file-server" {
			readFileServerHttpServletExtensionResponseDefault(ctx, updateResponse.FileServerHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "config" {
			readConfigHttpServletExtensionResponseDefault(ctx, updateResponse.ConfigHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "scim2" {
			readScim2HttpServletExtensionResponseDefault(ctx, updateResponse.Scim2HttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "directory-rest-api" {
			readDirectoryRestApiHttpServletExtensionResponseDefault(ctx, updateResponse.DirectoryRestApiHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party" {
			readThirdPartyHttpServletExtensionResponseDefault(ctx, updateResponse.ThirdPartyHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *defaultHttpServletExtensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *httpServletExtensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state httpServletExtensionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.HttpServletExtensionApi.DeleteHttpServletExtensionExecute(r.apiClient.HttpServletExtensionApi.DeleteHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Http Servlet Extension", err, httpResp)
		return
	}
}

func (r *httpServletExtensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importHttpServletExtension(ctx, req, resp)
}

func (r *defaultHttpServletExtensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importHttpServletExtension(ctx, req, resp)
}

func importHttpServletExtension(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
