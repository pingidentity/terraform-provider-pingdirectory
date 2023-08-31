package httpservletextension

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
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
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultHttpServletExtensionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type httpServletExtensionResourceModel struct {
	Id                                 types.String `tfsdk:"id"`
	Name                               types.String `tfsdk:"name"`
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party HTTP Servlet Extension. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
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
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
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
				Description:         "When the `type` attribute is set to  one of [`availability-state`, `prometheus-monitoring`]: Specifies the base context path that HTTP clients should use to access this servlet. The value must start with a forward slash and must represent a valid HTTP context path. When the `type` attribute is set to `velocity`: The context path to use to access all template-based and static content. The value must start with a forward slash and must represent a valid HTTP context path. When the `type` attribute is set to `ldap-mapped-scim`: The context path to use to access the SCIM interface. The value must start with a forward slash and must represent a valid HTTP context path. When the `type` attribute is set to `file-server`: Specifies the base context path that should be used by HTTP clients to reference content. The value must start with a forward slash and must represent a valid HTTP context path. When the `type` attribute is set to `scim2`: The context path to use to access the SCIM 2.0 interface. The value must start with a forward slash and must represent a valid HTTP context path.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`availability-state`, `prometheus-monitoring`]: Specifies the base context path that HTTP clients should use to access this servlet. The value must start with a forward slash and must represent a valid HTTP context path.\n  - `velocity`: The context path to use to access all template-based and static content. The value must start with a forward slash and must represent a valid HTTP context path.\n  - `ldap-mapped-scim`: The context path to use to access the SCIM interface. The value must start with a forward slash and must represent a valid HTTP context path.\n  - `file-server`: Specifies the base context path that should be used by HTTP clients to reference content. The value must start with a forward slash and must represent a valid HTTP context path.\n  - `scim2`: The context path to use to access the SCIM 2.0 interface. The value must start with a forward slash and must represent a valid HTTP context path.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id_token_validator": schema.SetAttribute{
				Description: "The ID token validators that may be used to verify the authenticity of an of an OpenID Connect ID token.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
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
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
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
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"exclude_ldap_objectclass": schema.SetAttribute{
				Description: "Specifies the LDAP object classes that should be not be exposed directly as SCIM resources.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"include_ldap_base_dn": schema.SetAttribute{
				Description: "Specifies the base DNs for the branches of the DIT that should be exposed via the Identity Access API.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"exclude_ldap_base_dn": schema.SetAttribute{
				Description: "Specifies the base DNs for the branches of the DIT that should not be exposed via the Identity Access API.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"entity_tag_ldap_attribute": schema.StringAttribute{
				Description: "Specifies the LDAP attribute whose value should be used as the entity tag value to enable SCIM resource versioning support.",
				Optional:    true,
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
				Description:         "When the `type` attribute is set to `ldap-mapped-scim`: Enables debug logging of the SCIM SDK. Debug messages will be forwarded to the Directory Server debug logger with the scope of com.unboundid.directory.server.extensions.scim.SCIMHTTPServletExtension. When the `type` attribute is set to `scim2`: Enables debug logging of the SCIM 2.0 SDK. Debug messages will be forwarded to the Directory Server debug logger with the scope of com.unboundid.directory.broker.http.scim2.extension.SCIM2HTTPServletExtension.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `ldap-mapped-scim`: Enables debug logging of the SCIM SDK. Debug messages will be forwarded to the Directory Server debug logger with the scope of com.unboundid.directory.server.extensions.scim.SCIMHTTPServletExtension.\n  - `scim2`: Enables debug logging of the SCIM 2.0 SDK. Debug messages will be forwarded to the Directory Server debug logger with the scope of com.unboundid.directory.broker.http.scim2.extension.SCIM2HTTPServletExtension.",
				Optional:            true,
				Computed:            true,
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
				Description:         "When the `type` attribute is set to `velocity`: Specifies the path to a file that contains MIME type mappings that will be used to determine the appropriate value to return for the Content-Type header based on the extension of the requested static content file. When the `type` attribute is set to `file-server`: Specifies the path to a file that contains MIME type mappings that will be used to determine the appropriate value to return for the Content-Type header based on the extension of the requested file.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `velocity`: Specifies the path to a file that contains MIME type mappings that will be used to determine the appropriate value to return for the Content-Type header based on the extension of the requested static content file.\n  - `file-server`: Specifies the path to a file that contains MIME type mappings that will be used to determine the appropriate value to return for the Content-Type header based on the extension of the requested file.",
				Optional:            true,
			},
			"default_mime_type": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `velocity`: Specifies the default value that will be used in the response's Content-Type header that indicates the type of content to return. When the `type` attribute is set to `file-server`: Specifies the default MIME type to use for the Content-Type header when a mapping cannot be found.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `velocity`: Specifies the default value that will be used in the response's Content-Type header that indicates the type of content to return.\n  - `file-server`: Specifies the default MIME type to use for the Content-Type header when a mapping cannot be found.",
				Optional:            true,
				Computed:            true,
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
				Description:         "When the `type` attribute is set to `velocity`: Require authentication when accessing Velocity templates. When the `type` attribute is set to `file-server`: Indicates whether the servlet extension should only accept requests from authenticated clients.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `velocity`: Require authentication when accessing Velocity templates.\n  - `file-server`: Indicates whether the servlet extension should only accept requests from authenticated clients.",
				Optional:            true,
				Computed:            true,
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
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
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
				Description:         "When the `type` attribute is set to  one of [`delegated-admin`, `consent`, `directory-rest-api`]: Enables HTTP Basic authentication, using a username and password. The Identity Mapper specified by the identity-mapper property will be used to map the username to a DN. When the `type` attribute is set to `ldap-mapped-scim`: Enables HTTP Basic authentication, using a username and password.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`delegated-admin`, `consent`, `directory-rest-api`]: Enables HTTP Basic authentication, using a username and password. The Identity Mapper specified by the identity-mapper property will be used to map the username to a DN.\n  - `ldap-mapped-scim`: Enables HTTP Basic authentication, using a username and password.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"identity_mapper": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `delegated-admin`: Specifies the Identity Mapper that is to be used for associating user entries with basic authentication user names. When the `type` attribute is set to `velocity`: Specifies the name of the identity mapper that is to be used for associating basic authentication credentials with user entries. When the `type` attribute is set to `consent`: Specifies the Identity Mapper that is to be used for associating basic authentication usernames with DNs. When the `type` attribute is set to `ldap-mapped-scim`: Specifies the name of the identity mapper that is to be used to match the username included in the HTTP Basic authentication header to the corresponding user in the directory. When the `type` attribute is set to `file-server`: The identity mapper that will be used to identify the entry with which a username is associated. When the `type` attribute is set to `config`: Specifies the name of the identity mapper that is to be used for associating user entries with basic authentication user names. When the `type` attribute is set to `directory-rest-api`: Specifies the Identity Mapper that is to be used for associating user entries with basic authentication usernames.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `delegated-admin`: Specifies the Identity Mapper that is to be used for associating user entries with basic authentication user names.\n  - `velocity`: Specifies the name of the identity mapper that is to be used for associating basic authentication credentials with user entries.\n  - `consent`: Specifies the Identity Mapper that is to be used for associating basic authentication usernames with DNs.\n  - `ldap-mapped-scim`: Specifies the name of the identity mapper that is to be used to match the username included in the HTTP Basic authentication header to the corresponding user in the directory.\n  - `file-server`: The identity mapper that will be used to identify the entry with which a username is associated.\n  - `config`: Specifies the name of the identity mapper that is to be used for associating user entries with basic authentication user names.\n  - `directory-rest-api`: Specifies the Identity Mapper that is to be used for associating user entries with basic authentication usernames.",
				Optional:            true,
			},
			"access_token_validator": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `delegated-admin`: If specified, the Access Token Validator(s) that may be used to validate access tokens for requests submitted to this Delegated Admin HTTP Servlet Extension. When the `type` attribute is set to `consent`: If specified, the Access Token Validator(s) that may be used to validate access tokens for requests submitted to this Consent HTTP Servlet Extension. When the `type` attribute is set to `file-server`: The access token validators that may be used to verify the authenticity of an OAuth 2.0 bearer token. When the `type` attribute is set to `scim2`: If specified, the Access Token Validator(s) that may be used to validate access tokens for requests submitted to this SCIM2 HTTP Servlet Extension. When the `type` attribute is set to `directory-rest-api`: If specified, the Access Token Validator(s) that may be used to validate access tokens for requests submitted to this Directory REST API HTTP Servlet Extension.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `delegated-admin`: If specified, the Access Token Validator(s) that may be used to validate access tokens for requests submitted to this Delegated Admin HTTP Servlet Extension.\n  - `consent`: If specified, the Access Token Validator(s) that may be used to validate access tokens for requests submitted to this Consent HTTP Servlet Extension.\n  - `file-server`: The access token validators that may be used to verify the authenticity of an OAuth 2.0 bearer token.\n  - `scim2`: If specified, the Access Token Validator(s) that may be used to validate access tokens for requests submitted to this SCIM2 HTTP Servlet Extension.\n  - `directory-rest-api`: If specified, the Access Token Validator(s) that may be used to validate access tokens for requests submitted to this Directory REST API HTTP Servlet Extension.",
				Optional:            true,
				Computed:            true,
				Default:             internaltypes.EmptySetDefault(types.StringType),
				ElementType:         types.StringType,
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
				Description:         "When the `type` attribute is set to  one of [`delegated-admin`, `quickstart`, `availability-state`, `prometheus-monitoring`, `consent`, `ldap-mapped-scim`, `groovy-scripted`, `file-server`, `config`, `scim2`, `directory-rest-api`, `third-party`]: Specifies HTTP header fields and values added to response headers for all requests. When the `type` attribute is set to `velocity`: Specifies HTTP header fields and values added to response headers for all template page requests.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`delegated-admin`, `quickstart`, `availability-state`, `prometheus-monitoring`, `consent`, `ldap-mapped-scim`, `groovy-scripted`, `file-server`, `config`, `scim2`, `directory-rest-api`, `third-party`]: Specifies HTTP header fields and values added to response headers for all requests.\n  - `velocity`: Specifies HTTP header fields and values added to response headers for all template page requests.",
				Optional:            true,
				Computed:            true,
				Default:             internaltypes.EmptySetDefault(types.StringType),
				ElementType:         types.StringType,
			},
			"correlation_id_response_header": schema.StringAttribute{
				Description: "Specifies the name of the HTTP response header that will contain a correlation ID value. Example values are \"Correlation-Id\", \"X-Amzn-Trace-Id\", and \"X-Request-Id\".",
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
		typeAttr.Validators = []validator.String{
			stringvalidator.OneOf([]string{"delegated-admin", "quickstart", "availability-state", "prometheus-monitoring", "velocity", "consent", "ldap-mapped-scim", "groovy-scripted", "file-server", "config", "scim2", "directory-rest-api", "third-party"}...),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		schemaDef.Attributes["map_access_tokens_to_local_users"] = schema.StringAttribute{
			Description: "Indicates whether the SCIM2 servlet should attempt to map the presented access token to a local user.",
		}
		schemaDef.Attributes["max_page_size"] = schema.Int64Attribute{
			Description: "The maximum number of entries to be returned in one page of search results.",
		}
		schemaDef.Attributes["schemas_endpoint_objectclass"] = schema.SetAttribute{
			Description: "The list of object classes which will be returned by the schemas endpoint.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["default_operational_attribute"] = schema.SetAttribute{
			Description: "A set of operational attributes that will be returned with entries by default.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["reject_expansion_attribute"] = schema.SetAttribute{
			Description: "A set of attributes which the client is not allowed to provide for the expand query parameters. This should be used for attributes that could either have a large number of values or that reference entries that are very large like groups.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["always_use_permissive_modify"] = schema.BoolAttribute{
			Description: "Supported in PingDirectory product version 9.3.0.0+. Indicates whether to always use permissive modify behavior for PATCH requests, even if the request did not include the permissive modify request control.",
		}
		schemaDef.Attributes["allowed_control"] = schema.SetAttribute{
			Description: "Specifies the names of any request controls that should be allowed by the Directory REST API. Any request that contains a critical control not in this list will be rejected. Any non-critical request control which is not supported by the Directory REST API will be removed from the request.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["swagger_enabled"] = schema.BoolAttribute{
			Description: "Indicates whether the SCIM2 HTTP Servlet Extension will generate a Swagger specification document.",
		}
		schemaDef.Attributes["bearer_token_auth_enabled"] = schema.BoolAttribute{
			Description: "Enables HTTP bearer token authentication.",
		}
		schemaDef.Attributes["static_context_path"] = schema.StringAttribute{
			Description: "The path below the base context path by which static, non-template content such as images, CSS, and Javascript files are accessible.",
		}
		schemaDef.Attributes["static_content_directory"] = schema.StringAttribute{
			Description: "Specifies the base directory in which static, non-template content such as images, CSS, and Javascript files are stored on the filesystem.",
		}
		schemaDef.Attributes["static_custom_directory"] = schema.StringAttribute{
			Description: "Specifies the base directory in which custom static, non-template content such as images, CSS, and Javascript files are stored on the filesystem. Files in this directory will override those with the same name in the directory specified by the static-content-directory property.",
		}
		schemaDef.Attributes["template_directory"] = schema.SetAttribute{
			Description: "Specifies an ordered list of directories in which to search for the template files.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["expose_request_attributes"] = schema.BoolAttribute{
			Description: "Specifies whether the HTTP request will be exposed to templates.",
		}
		schemaDef.Attributes["expose_session_attributes"] = schema.BoolAttribute{
			Description: "Specifies whether the HTTP session will be exposed to templates.",
		}
		schemaDef.Attributes["expose_server_context"] = schema.BoolAttribute{
			Description: "Specifies whether a server context will be exposed under context key 'ubid_server' for all template contexts.",
		}
		schemaDef.Attributes["allow_context_override"] = schema.BoolAttribute{
			Description: "Indicates whether context providers may override existing context objects with new values.",
		}
		schemaDef.Attributes["character_encoding"] = schema.StringAttribute{
			Description: "Specifies the value that will be used for all responses' Content-Type headers' charset parameter that indicates the character encoding of the document.",
		}
		schemaDef.Attributes["static_response_header"] = schema.SetAttribute{
			Description: "Specifies HTTP header fields and values added to response headers for static content requests such as images and scripts.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["access_token_scope"] = schema.StringAttribute{
			Description:         "When the `type` attribute is set to `delegated-admin`: The name of a scope that must be present in an access token accepted by the Delegated Admin HTTP Servlet Extension. When the `type` attribute is set to `directory-rest-api`: The name of a scope that must be present in an access token accepted by the Directory REST API HTTP Servlet Extension.",
			MarkdownDescription: "When the `type` attribute is set to:\n  - `delegated-admin`: The name of a scope that must be present in an access token accepted by the Delegated Admin HTTP Servlet Extension.\n  - `directory-rest-api`: The name of a scope that must be present in an access token accepted by the Directory REST API HTTP Servlet Extension.",
		}
		schemaDef.Attributes["audience"] = schema.StringAttribute{
			Description:         "When the `type` attribute is set to `delegated-admin`: A string or URI that identifies the Delegated Admin HTTP Servlet Extension in the context of OAuth2 authorization. When the `type` attribute is set to `directory-rest-api`: A string or URI that identifies the Directory REST API HTTP Servlet Extension in the context of OAuth2 authorization.",
			MarkdownDescription: "When the `type` attribute is set to:\n  - `delegated-admin`: A string or URI that identifies the Delegated Admin HTTP Servlet Extension in the context of OAuth2 authorization.\n  - `directory-rest-api`: A string or URI that identifies the Directory REST API HTTP Servlet Extension in the context of OAuth2 authorization.",
		}
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *httpServletExtensionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanHttpServletExtension(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_http_servlet_extension")
	var model httpServletExtensionResourceModel
	req.Plan.Get(ctx, &model)
	resourceType := model.Type.ValueString()
	// Set defaults for availability-state type
	if resourceType == "availability-state" {
		if !internaltypes.IsDefined(model.IncludeResponseBody) {
			model.IncludeResponseBody = types.BoolValue(true)
		}
	}
	// Set defaults for prometheus-monitoring type
	if resourceType == "prometheus-monitoring" {
		if !internaltypes.IsDefined(model.BaseContextPath) {
			model.BaseContextPath = types.StringValue("/metrics")
		}
		if !internaltypes.IsDefined(model.IncludeInstanceNameLabel) {
			model.IncludeInstanceNameLabel = types.BoolValue(false)
		}
		if !internaltypes.IsDefined(model.IncludeProductNameLabel) {
			model.IncludeProductNameLabel = types.BoolValue(false)
		}
		if !internaltypes.IsDefined(model.IncludeLocationNameLabel) {
			model.IncludeLocationNameLabel = types.BoolValue(false)
		}
		if !internaltypes.IsDefined(model.AlwaysIncludeMonitorEntryNameLabel) {
			model.AlwaysIncludeMonitorEntryNameLabel = types.BoolValue(false)
		}
		if !internaltypes.IsDefined(model.IncludeMonitorObjectClassNameLabel) {
			model.IncludeMonitorObjectClassNameLabel = types.BoolValue(false)
		}
		if !internaltypes.IsDefined(model.IncludeMonitorAttributeNameLabel) {
			model.IncludeMonitorAttributeNameLabel = types.BoolValue(false)
		}
	}
	// Set defaults for ldap-mapped-scim type
	if resourceType == "ldap-mapped-scim" {
		if !internaltypes.IsDefined(model.BasicAuthEnabled) {
			model.BasicAuthEnabled = types.BoolValue(true)
		}
		if !internaltypes.IsDefined(model.ResourceMappingFile) {
			model.ResourceMappingFile = types.StringValue("config/scim-resources.xml")
		}
		if !internaltypes.IsDefined(model.BaseContextPath) {
			model.BaseContextPath = types.StringValue("/")
		}
		if !internaltypes.IsDefined(model.TemporaryDirectory) {
			model.TemporaryDirectory = types.StringValue("scim-data-tmp")
		}
		if !internaltypes.IsDefined(model.TemporaryDirectoryPermissions) {
			model.TemporaryDirectoryPermissions = types.StringValue("700")
		}
		if !internaltypes.IsDefined(model.MaxResults) {
			model.MaxResults = types.Int64Value(100)
		}
		if !internaltypes.IsDefined(model.BulkMaxOperations) {
			model.BulkMaxOperations = types.Int64Value(10000)
		}
		if !internaltypes.IsDefined(model.BulkMaxPayloadSize) {
			model.BulkMaxPayloadSize = types.StringValue("10 MB")
		}
		if !internaltypes.IsDefined(model.BulkMaxConcurrentRequests) {
			model.BulkMaxConcurrentRequests = types.Int64Value(10)
		}
		if !internaltypes.IsDefined(model.DebugEnabled) {
			model.DebugEnabled = types.BoolValue(false)
		}
		if !internaltypes.IsDefined(model.DebugLevel) {
			model.DebugLevel = types.StringValue("info")
		}
		if !internaltypes.IsDefined(model.DebugType) {
			model.DebugType, _ = types.SetValue(types.StringType, []attr.Value{types.StringValue("coding-error"), types.StringValue("exception")})
		}
		if !internaltypes.IsDefined(model.IncludeStackTrace) {
			model.IncludeStackTrace = types.BoolValue(false)
		}
	}
	// Set defaults for file-server type
	if resourceType == "file-server" {
		if !internaltypes.IsDefined(model.EnableDirectoryIndexing) {
			model.EnableDirectoryIndexing = types.BoolValue(false)
		}
		if !internaltypes.IsDefined(model.IndexFile) {
			model.IndexFile, _ = types.SetValue(types.StringType, []attr.Value{types.StringValue("index.html"), types.StringValue("index.htm")})
		}
		if !internaltypes.IsDefined(model.DefaultMIMEType) {
			model.DefaultMIMEType = types.StringValue("application/octet-stream")
		}
		if !internaltypes.IsDefined(model.RequireAuthentication) {
			model.RequireAuthentication = types.BoolValue(false)
		}
		if !internaltypes.IsDefined(model.AllowedAuthenticationType) {
			model.AllowedAuthenticationType, _ = types.SetValue(types.StringType, []attr.Value{types.StringValue("basic"), types.StringValue("access-token"), types.StringValue("id-token")})
		}
		if !internaltypes.IsDefined(model.RequireFileServletAccessPrivilege) {
			model.RequireFileServletAccessPrivilege = types.BoolValue(false)
		}
	}
	resp.Plan.Set(ctx, &model)
}

func (r *defaultHttpServletExtensionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanHttpServletExtension(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_default_http_servlet_extension")
}

func modifyPlanHttpServletExtension(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, resourceName string) {
	compare, err := version.Compare(providerConfig.ProductVersion, version.PingDirectory9300)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model defaultHttpServletExtensionResourceModel
	req.Plan.Get(ctx, &model)
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

// Add config validators that apply to both default_ and non-default_
func configValidatorsHttpServletExtension() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherValidator(
			path.MatchRoot("type"),
			[]string{"ldap-mapped-scim"},
			resourcevalidator.Conflicting(
				path.MatchRoot("include_ldap_objectclass"),
				path.MatchRoot("exclude_ldap_objectclass"),
			),
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("basic_auth_enabled"),
			path.MatchRoot("type"),
			[]string{"delegated-admin", "consent", "ldap-mapped-scim", "directory-rest-api"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("identity_mapper"),
			path.MatchRoot("type"),
			[]string{"delegated-admin", "velocity", "consent", "ldap-mapped-scim", "file-server", "config", "directory-rest-api"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("access_token_validator"),
			path.MatchRoot("type"),
			[]string{"delegated-admin", "consent", "file-server", "scim2", "directory-rest-api"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("server"),
			path.MatchRoot("type"),
			[]string{"quickstart"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("base_context_path"),
			path.MatchRoot("type"),
			[]string{"availability-state", "prometheus-monitoring", "velocity", "ldap-mapped-scim", "file-server", "scim2"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("available_status_code"),
			path.MatchRoot("type"),
			[]string{"availability-state"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("degraded_status_code"),
			path.MatchRoot("type"),
			[]string{"availability-state"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("unavailable_status_code"),
			path.MatchRoot("type"),
			[]string{"availability-state"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("override_status_code"),
			path.MatchRoot("type"),
			[]string{"availability-state"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_response_body"),
			path.MatchRoot("type"),
			[]string{"availability-state"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("additional_response_contents"),
			path.MatchRoot("type"),
			[]string{"availability-state"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_instance_name_label"),
			path.MatchRoot("type"),
			[]string{"prometheus-monitoring"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_product_name_label"),
			path.MatchRoot("type"),
			[]string{"prometheus-monitoring"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_location_name_label"),
			path.MatchRoot("type"),
			[]string{"prometheus-monitoring"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("always_include_monitor_entry_name_label"),
			path.MatchRoot("type"),
			[]string{"prometheus-monitoring"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_monitor_object_class_name_label"),
			path.MatchRoot("type"),
			[]string{"prometheus-monitoring"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_monitor_attribute_name_label"),
			path.MatchRoot("type"),
			[]string{"prometheus-monitoring"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("label_name_value_pair"),
			path.MatchRoot("type"),
			[]string{"prometheus-monitoring"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("mime_types_file"),
			path.MatchRoot("type"),
			[]string{"velocity", "file-server"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("default_mime_type"),
			path.MatchRoot("type"),
			[]string{"velocity", "file-server"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("require_authentication"),
			path.MatchRoot("type"),
			[]string{"velocity", "file-server"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("oauth_token_handler"),
			path.MatchRoot("type"),
			[]string{"ldap-mapped-scim"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("resource_mapping_file"),
			path.MatchRoot("type"),
			[]string{"ldap-mapped-scim"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_ldap_objectclass"),
			path.MatchRoot("type"),
			[]string{"ldap-mapped-scim"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("exclude_ldap_objectclass"),
			path.MatchRoot("type"),
			[]string{"ldap-mapped-scim"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_ldap_base_dn"),
			path.MatchRoot("type"),
			[]string{"ldap-mapped-scim"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("exclude_ldap_base_dn"),
			path.MatchRoot("type"),
			[]string{"ldap-mapped-scim"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("entity_tag_ldap_attribute"),
			path.MatchRoot("type"),
			[]string{"ldap-mapped-scim"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("temporary_directory"),
			path.MatchRoot("type"),
			[]string{"ldap-mapped-scim"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("temporary_directory_permissions"),
			path.MatchRoot("type"),
			[]string{"ldap-mapped-scim"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("max_results"),
			path.MatchRoot("type"),
			[]string{"ldap-mapped-scim"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("bulk_max_operations"),
			path.MatchRoot("type"),
			[]string{"ldap-mapped-scim"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("bulk_max_payload_size"),
			path.MatchRoot("type"),
			[]string{"ldap-mapped-scim"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("bulk_max_concurrent_requests"),
			path.MatchRoot("type"),
			[]string{"ldap-mapped-scim"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("debug_enabled"),
			path.MatchRoot("type"),
			[]string{"ldap-mapped-scim", "scim2"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("debug_level"),
			path.MatchRoot("type"),
			[]string{"ldap-mapped-scim", "scim2"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("debug_type"),
			path.MatchRoot("type"),
			[]string{"ldap-mapped-scim", "scim2"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_stack_trace"),
			path.MatchRoot("type"),
			[]string{"ldap-mapped-scim", "scim2"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("script_class"),
			path.MatchRoot("type"),
			[]string{"groovy-scripted"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("script_argument"),
			path.MatchRoot("type"),
			[]string{"groovy-scripted"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("document_root_directory"),
			path.MatchRoot("type"),
			[]string{"file-server"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("enable_directory_indexing"),
			path.MatchRoot("type"),
			[]string{"file-server"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("index_file"),
			path.MatchRoot("type"),
			[]string{"file-server"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("allowed_authentication_type"),
			path.MatchRoot("type"),
			[]string{"file-server"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("id_token_validator"),
			path.MatchRoot("type"),
			[]string{"file-server"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("require_file_servlet_access_privilege"),
			path.MatchRoot("type"),
			[]string{"file-server"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("require_group"),
			path.MatchRoot("type"),
			[]string{"file-server"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_class"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_argument"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"availability-state",
			[]path.Expression{path.MatchRoot("base_context_path"), path.MatchRoot("available_status_code"), path.MatchRoot("degraded_status_code"), path.MatchRoot("unavailable_status_code")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"file-server",
			[]path.Expression{path.MatchRoot("base_context_path"), path.MatchRoot("document_root_directory")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"groovy-scripted",
			[]path.Expression{path.MatchRoot("script_class")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party",
			[]path.Expression{path.MatchRoot("extension_class")},
		),
	}
}

// Add config validators
func (r httpServletExtensionResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsHttpServletExtension()
}

// Add config validators
func (r defaultHttpServletExtensionResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	validators := []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("access_token_scope"),
			path.MatchRoot("type"),
			[]string{"delegated-admin", "directory-rest-api"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("audience"),
			path.MatchRoot("type"),
			[]string{"delegated-admin", "directory-rest-api"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("static_context_path"),
			path.MatchRoot("type"),
			[]string{"velocity"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("static_content_directory"),
			path.MatchRoot("type"),
			[]string{"velocity"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("static_custom_directory"),
			path.MatchRoot("type"),
			[]string{"velocity"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("template_directory"),
			path.MatchRoot("type"),
			[]string{"velocity"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("expose_request_attributes"),
			path.MatchRoot("type"),
			[]string{"velocity"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("expose_session_attributes"),
			path.MatchRoot("type"),
			[]string{"velocity"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("expose_server_context"),
			path.MatchRoot("type"),
			[]string{"velocity"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("allow_context_override"),
			path.MatchRoot("type"),
			[]string{"velocity"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("character_encoding"),
			path.MatchRoot("type"),
			[]string{"velocity"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("static_response_header"),
			path.MatchRoot("type"),
			[]string{"velocity"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("bearer_token_auth_enabled"),
			path.MatchRoot("type"),
			[]string{"consent"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("map_access_tokens_to_local_users"),
			path.MatchRoot("type"),
			[]string{"scim2"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("swagger_enabled"),
			path.MatchRoot("type"),
			[]string{"scim2"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("max_page_size"),
			path.MatchRoot("type"),
			[]string{"directory-rest-api"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("schemas_endpoint_objectclass"),
			path.MatchRoot("type"),
			[]string{"directory-rest-api"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("default_operational_attribute"),
			path.MatchRoot("type"),
			[]string{"directory-rest-api"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("reject_expansion_attribute"),
			path.MatchRoot("type"),
			[]string{"directory-rest-api"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("always_use_permissive_modify"),
			path.MatchRoot("type"),
			[]string{"directory-rest-api"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("allowed_control"),
			path.MatchRoot("type"),
			[]string{"directory-rest-api"},
		),
	}
	return append(configValidatorsHttpServletExtension(), validators...)
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
func populateHttpServletExtensionUnknownValues(model *httpServletExtensionResourceModel) {
	if model.ExcludeLDAPObjectclass.IsUnknown() || model.ExcludeLDAPObjectclass.IsNull() {
		model.ExcludeLDAPObjectclass, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExcludeLDAPBaseDN.IsUnknown() || model.ExcludeLDAPBaseDN.IsNull() {
		model.ExcludeLDAPBaseDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IndexFile.IsUnknown() || model.IndexFile.IsNull() {
		model.IndexFile, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.LabelNameValuePair.IsUnknown() || model.LabelNameValuePair.IsNull() {
		model.LabelNameValuePair, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.DebugType.IsUnknown() || model.DebugType.IsNull() {
		model.DebugType, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IdTokenValidator.IsUnknown() || model.IdTokenValidator.IsNull() {
		model.IdTokenValidator, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AccessTokenValidator.IsUnknown() || model.AccessTokenValidator.IsNull() {
		model.AccessTokenValidator, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ScriptArgument.IsUnknown() || model.ScriptArgument.IsNull() {
		model.ScriptArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IncludeLDAPObjectclass.IsUnknown() || model.IncludeLDAPObjectclass.IsNull() {
		model.IncludeLDAPObjectclass, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.RequireGroup.IsUnknown() || model.RequireGroup.IsNull() {
		model.RequireGroup, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AllowedAuthenticationType.IsUnknown() || model.AllowedAuthenticationType.IsNull() {
		model.AllowedAuthenticationType, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IncludeLDAPBaseDN.IsUnknown() || model.IncludeLDAPBaseDN.IsNull() {
		model.IncludeLDAPBaseDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.BaseContextPath.IsUnknown() || model.BaseContextPath.IsNull() {
		model.BaseContextPath = types.StringValue("")
	}
	if model.ResourceMappingFile.IsUnknown() || model.ResourceMappingFile.IsNull() {
		model.ResourceMappingFile = types.StringValue("")
	}
	if model.DefaultMIMEType.IsUnknown() || model.DefaultMIMEType.IsNull() {
		model.DefaultMIMEType = types.StringValue("")
	}
	if model.TemporaryDirectory.IsUnknown() || model.TemporaryDirectory.IsNull() {
		model.TemporaryDirectory = types.StringValue("")
	}
	if model.TemporaryDirectoryPermissions.IsUnknown() || model.TemporaryDirectoryPermissions.IsNull() {
		model.TemporaryDirectoryPermissions = types.StringValue("")
	}
	if model.BulkMaxPayloadSize.IsUnknown() || model.BulkMaxPayloadSize.IsNull() {
		model.BulkMaxPayloadSize = types.StringValue("")
	}
	if model.DebugLevel.IsUnknown() || model.DebugLevel.IsNull() {
		model.DebugLevel = types.StringValue("")
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateHttpServletExtensionUnknownValuesDefault(model *defaultHttpServletExtensionResourceModel) {
	if model.ExcludeLDAPObjectclass.IsUnknown() || model.ExcludeLDAPObjectclass.IsNull() {
		model.ExcludeLDAPObjectclass, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExcludeLDAPBaseDN.IsUnknown() || model.ExcludeLDAPBaseDN.IsNull() {
		model.ExcludeLDAPBaseDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IndexFile.IsUnknown() || model.IndexFile.IsNull() {
		model.IndexFile, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.SchemasEndpointObjectclass.IsUnknown() || model.SchemasEndpointObjectclass.IsNull() {
		model.SchemasEndpointObjectclass, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.LabelNameValuePair.IsUnknown() || model.LabelNameValuePair.IsNull() {
		model.LabelNameValuePair, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.DebugType.IsUnknown() || model.DebugType.IsNull() {
		model.DebugType, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IdTokenValidator.IsUnknown() || model.IdTokenValidator.IsNull() {
		model.IdTokenValidator, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.RejectExpansionAttribute.IsUnknown() || model.RejectExpansionAttribute.IsNull() {
		model.RejectExpansionAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AllowedControl.IsUnknown() || model.AllowedControl.IsNull() {
		model.AllowedControl, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AccessTokenValidator.IsUnknown() || model.AccessTokenValidator.IsNull() {
		model.AccessTokenValidator, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.DefaultOperationalAttribute.IsUnknown() || model.DefaultOperationalAttribute.IsNull() {
		model.DefaultOperationalAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.TemplateDirectory.IsUnknown() || model.TemplateDirectory.IsNull() {
		model.TemplateDirectory, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ScriptArgument.IsUnknown() || model.ScriptArgument.IsNull() {
		model.ScriptArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IncludeLDAPObjectclass.IsUnknown() || model.IncludeLDAPObjectclass.IsNull() {
		model.IncludeLDAPObjectclass, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.RequireGroup.IsUnknown() || model.RequireGroup.IsNull() {
		model.RequireGroup, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AllowedAuthenticationType.IsUnknown() || model.AllowedAuthenticationType.IsNull() {
		model.AllowedAuthenticationType, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IncludeLDAPBaseDN.IsUnknown() || model.IncludeLDAPBaseDN.IsNull() {
		model.IncludeLDAPBaseDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.StaticResponseHeader.IsUnknown() || model.StaticResponseHeader.IsNull() {
		model.StaticResponseHeader, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ResourceMappingFile.IsUnknown() || model.ResourceMappingFile.IsNull() {
		model.ResourceMappingFile = types.StringValue("")
	}
	if model.StaticContextPath.IsUnknown() || model.StaticContextPath.IsNull() {
		model.StaticContextPath = types.StringValue("")
	}
	if model.Server.IsUnknown() || model.Server.IsNull() {
		model.Server = types.StringValue("")
	}
	if model.DefaultMIMEType.IsUnknown() || model.DefaultMIMEType.IsNull() {
		model.DefaultMIMEType = types.StringValue("")
	}
	if model.Audience.IsUnknown() || model.Audience.IsNull() {
		model.Audience = types.StringValue("")
	}
	if model.ExtensionClass.IsUnknown() || model.ExtensionClass.IsNull() {
		model.ExtensionClass = types.StringValue("")
	}
	if model.OAuthTokenHandler.IsUnknown() || model.OAuthTokenHandler.IsNull() {
		model.OAuthTokenHandler = types.StringValue("")
	}
	if model.StaticContentDirectory.IsUnknown() || model.StaticContentDirectory.IsNull() {
		model.StaticContentDirectory = types.StringValue("")
	}
	if model.MimeTypesFile.IsUnknown() || model.MimeTypesFile.IsNull() {
		model.MimeTypesFile = types.StringValue("")
	}
	if model.TemporaryDirectory.IsUnknown() || model.TemporaryDirectory.IsNull() {
		model.TemporaryDirectory = types.StringValue("")
	}
	if model.AccessTokenScope.IsUnknown() || model.AccessTokenScope.IsNull() {
		model.AccessTokenScope = types.StringValue("")
	}
	if model.EntityTagLDAPAttribute.IsUnknown() || model.EntityTagLDAPAttribute.IsNull() {
		model.EntityTagLDAPAttribute = types.StringValue("")
	}
	if model.BaseContextPath.IsUnknown() || model.BaseContextPath.IsNull() {
		model.BaseContextPath = types.StringValue("")
	}
	if model.MapAccessTokensToLocalUsers.IsUnknown() || model.MapAccessTokensToLocalUsers.IsNull() {
		model.MapAccessTokensToLocalUsers = types.StringValue("")
	}
	if model.IdentityMapper.IsUnknown() || model.IdentityMapper.IsNull() {
		model.IdentityMapper = types.StringValue("")
	}
	if model.DocumentRootDirectory.IsUnknown() || model.DocumentRootDirectory.IsNull() {
		model.DocumentRootDirectory = types.StringValue("")
	}
	if model.CharacterEncoding.IsUnknown() || model.CharacterEncoding.IsNull() {
		model.CharacterEncoding = types.StringValue("")
	}
	if model.AdditionalResponseContents.IsUnknown() || model.AdditionalResponseContents.IsNull() {
		model.AdditionalResponseContents = types.StringValue("")
	}
	if model.TemporaryDirectoryPermissions.IsUnknown() || model.TemporaryDirectoryPermissions.IsNull() {
		model.TemporaryDirectoryPermissions = types.StringValue("")
	}
	if model.ScriptClass.IsUnknown() || model.ScriptClass.IsNull() {
		model.ScriptClass = types.StringValue("")
	}
	if model.StaticCustomDirectory.IsUnknown() || model.StaticCustomDirectory.IsNull() {
		model.StaticCustomDirectory = types.StringValue("")
	}
	if model.BulkMaxPayloadSize.IsUnknown() || model.BulkMaxPayloadSize.IsNull() {
		model.BulkMaxPayloadSize = types.StringValue("")
	}
	if model.DebugLevel.IsUnknown() || model.DebugLevel.IsNull() {
		model.DebugLevel = types.StringValue("")
	}
}

// Read a DelegatedAdminHttpServletExtensionResponse object into the model struct
func readDelegatedAdminHttpServletExtensionResponseDefault(ctx context.Context, r *client.DelegatedAdminHttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("delegated-admin")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.BasicAuthEnabled = internaltypes.BoolTypeOrNil(r.BasicAuthEnabled)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, true)
	state.AccessTokenValidator = internaltypes.GetStringSet(r.AccessTokenValidator)
	state.AccessTokenScope = internaltypes.StringTypeOrNil(r.AccessTokenScope, true)
	state.Audience = internaltypes.StringTypeOrNil(r.Audience, true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, true)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(state)
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
	populateHttpServletExtensionUnknownValues(state)
}

// Read a QuickstartHttpServletExtensionResponse object into the model struct
func readQuickstartHttpServletExtensionResponseDefault(ctx context.Context, r *client.QuickstartHttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("quickstart")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Server = internaltypes.StringTypeOrNil(r.Server, true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, true)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(state)
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
	populateHttpServletExtensionUnknownValues(state)
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
	state.AdditionalResponseContents = internaltypes.StringTypeOrNil(r.AdditionalResponseContents, true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, true)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(state)
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
	populateHttpServletExtensionUnknownValues(state)
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
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, true)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(state)
}

// Read a VelocityHttpServletExtensionResponse object into the model struct
func readVelocityHttpServletExtensionResponseDefault(ctx context.Context, r *client.VelocityHttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("velocity")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
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
	populateHttpServletExtensionUnknownValuesDefault(state)
}

// Read a ConsentHttpServletExtensionResponse object into the model struct
func readConsentHttpServletExtensionResponseDefault(ctx context.Context, r *client.ConsentHttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("consent")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.BearerTokenAuthEnabled = internaltypes.BoolTypeOrNil(r.BearerTokenAuthEnabled)
	state.BasicAuthEnabled = internaltypes.BoolTypeOrNil(r.BasicAuthEnabled)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, true)
	state.AccessTokenValidator = internaltypes.GetStringSet(r.AccessTokenValidator)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, true)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(state)
}

// Read a LdapMappedScimHttpServletExtensionResponse object into the model struct
func readLdapMappedScimHttpServletExtensionResponse(ctx context.Context, r *client.LdapMappedScimHttpServletExtensionResponse, state *httpServletExtensionResourceModel, expectedValues *httpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap-mapped-scim")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.OAuthTokenHandler = internaltypes.StringTypeOrNil(r.OAuthTokenHandler, internaltypes.IsEmptyString(expectedValues.OAuthTokenHandler))
	state.BasicAuthEnabled = internaltypes.BoolTypeOrNil(r.BasicAuthEnabled)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, internaltypes.IsEmptyString(expectedValues.IdentityMapper))
	state.ResourceMappingFile = internaltypes.StringTypeOrNil(r.ResourceMappingFile, true)
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
	state.BulkMaxPayloadSize = internaltypes.StringTypeOrNil(r.BulkMaxPayloadSize, true)
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
	populateHttpServletExtensionUnknownValues(state)
}

// Read a LdapMappedScimHttpServletExtensionResponse object into the model struct
func readLdapMappedScimHttpServletExtensionResponseDefault(ctx context.Context, r *client.LdapMappedScimHttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap-mapped-scim")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.OAuthTokenHandler = internaltypes.StringTypeOrNil(r.OAuthTokenHandler, true)
	state.BasicAuthEnabled = internaltypes.BoolTypeOrNil(r.BasicAuthEnabled)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, true)
	state.ResourceMappingFile = internaltypes.StringTypeOrNil(r.ResourceMappingFile, true)
	state.IncludeLDAPObjectclass = internaltypes.GetStringSet(r.IncludeLDAPObjectclass)
	state.ExcludeLDAPObjectclass = internaltypes.GetStringSet(r.ExcludeLDAPObjectclass)
	state.IncludeLDAPBaseDN = internaltypes.GetStringSet(r.IncludeLDAPBaseDN)
	state.ExcludeLDAPBaseDN = internaltypes.GetStringSet(r.ExcludeLDAPBaseDN)
	state.EntityTagLDAPAttribute = internaltypes.StringTypeOrNil(r.EntityTagLDAPAttribute, true)
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.TemporaryDirectory = types.StringValue(r.TemporaryDirectory)
	state.TemporaryDirectoryPermissions = types.StringValue(r.TemporaryDirectoryPermissions)
	state.MaxResults = internaltypes.Int64TypeOrNil(r.MaxResults)
	state.BulkMaxOperations = internaltypes.Int64TypeOrNil(r.BulkMaxOperations)
	state.BulkMaxPayloadSize = internaltypes.StringTypeOrNil(r.BulkMaxPayloadSize, true)
	config.CheckMismatchedPDFormattedAttributes("bulk_max_payload_size",
		expectedValues.BulkMaxPayloadSize, state.BulkMaxPayloadSize, diagnostics)
	state.BulkMaxConcurrentRequests = internaltypes.Int64TypeOrNil(r.BulkMaxConcurrentRequests)
	state.DebugEnabled = internaltypes.BoolTypeOrNil(r.DebugEnabled)
	state.DebugLevel = types.StringValue(r.DebugLevel.String())
	state.DebugType = internaltypes.GetStringSet(
		client.StringSliceEnumhttpServletExtensionDebugTypeProp(r.DebugType))
	state.IncludeStackTrace = types.BoolValue(r.IncludeStackTrace)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, true)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(state)
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
	populateHttpServletExtensionUnknownValues(state)
}

// Read a GroovyScriptedHttpServletExtensionResponse object into the model struct
func readGroovyScriptedHttpServletExtensionResponseDefault(ctx context.Context, r *client.GroovyScriptedHttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, true)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(state)
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
	state.DefaultMIMEType = internaltypes.StringTypeOrNil(r.DefaultMIMEType, true)
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
	populateHttpServletExtensionUnknownValues(state)
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
	state.MimeTypesFile = internaltypes.StringTypeOrNil(r.MimeTypesFile, true)
	state.DefaultMIMEType = internaltypes.StringTypeOrNil(r.DefaultMIMEType, true)
	state.RequireAuthentication = internaltypes.BoolTypeOrNil(r.RequireAuthentication)
	state.AllowedAuthenticationType = internaltypes.GetStringSet(
		client.StringSliceEnumhttpServletExtensionAllowedAuthenticationTypeProp(r.AllowedAuthenticationType))
	state.AccessTokenValidator = internaltypes.GetStringSet(r.AccessTokenValidator)
	state.IdTokenValidator = internaltypes.GetStringSet(r.IdTokenValidator)
	state.RequireFileServletAccessPrivilege = internaltypes.BoolTypeOrNil(r.RequireFileServletAccessPrivilege)
	state.RequireGroup = internaltypes.GetStringSet(r.RequireGroup)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, true)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(state)
}

// Read a ConfigHttpServletExtensionResponse object into the model struct
func readConfigHttpServletExtensionResponseDefault(ctx context.Context, r *client.ConfigHttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("config")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, true)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(state)
}

// Read a Scim2HttpServletExtensionResponse object into the model struct
func readScim2HttpServletExtensionResponseDefault(ctx context.Context, r *client.Scim2HttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("scim2")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.AccessTokenValidator = internaltypes.GetStringSet(r.AccessTokenValidator)
	state.MapAccessTokensToLocalUsers = internaltypes.StringTypeOrNil(
		client.StringPointerEnumhttpServletExtensionMapAccessTokensToLocalUsersProp(r.MapAccessTokensToLocalUsers), true)
	state.DebugEnabled = internaltypes.BoolTypeOrNil(r.DebugEnabled)
	state.DebugLevel = types.StringValue(r.DebugLevel.String())
	state.DebugType = internaltypes.GetStringSet(
		client.StringSliceEnumhttpServletExtensionDebugTypeProp(r.DebugType))
	state.IncludeStackTrace = types.BoolValue(r.IncludeStackTrace)
	state.SwaggerEnabled = internaltypes.BoolTypeOrNil(r.SwaggerEnabled)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, true)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(state)
}

// Read a DirectoryRestApiHttpServletExtensionResponse object into the model struct
func readDirectoryRestApiHttpServletExtensionResponseDefault(ctx context.Context, r *client.DirectoryRestApiHttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("directory-rest-api")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.BasicAuthEnabled = internaltypes.BoolTypeOrNil(r.BasicAuthEnabled)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, true)
	state.AccessTokenValidator = internaltypes.GetStringSet(r.AccessTokenValidator)
	state.AccessTokenScope = internaltypes.StringTypeOrNil(r.AccessTokenScope, true)
	state.Audience = internaltypes.StringTypeOrNil(r.Audience, true)
	state.MaxPageSize = internaltypes.Int64TypeOrNil(r.MaxPageSize)
	state.SchemasEndpointObjectclass = internaltypes.GetStringSet(r.SchemasEndpointObjectclass)
	state.DefaultOperationalAttribute = internaltypes.GetStringSet(r.DefaultOperationalAttribute)
	state.RejectExpansionAttribute = internaltypes.GetStringSet(r.RejectExpansionAttribute)
	state.AlwaysUsePermissiveModify = internaltypes.BoolTypeOrNil(r.AlwaysUsePermissiveModify)
	state.AllowedControl = internaltypes.GetStringSet(
		client.StringSliceEnumhttpServletExtensionAllowedControlProp(r.AllowedControl))
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, true)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(state)
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
	populateHttpServletExtensionUnknownValues(state)
}

// Read a ThirdPartyHttpServletExtensionResponse object into the model struct
func readThirdPartyHttpServletExtensionResponseDefault(ctx context.Context, r *client.ThirdPartyHttpServletExtensionResponse, state *defaultHttpServletExtensionResourceModel, expectedValues *defaultHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, true)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateHttpServletExtensionUnknownValuesDefault(state)
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
	if readResponse.DelegatedAdminHttpServletExtensionResponse != nil {
		readDelegatedAdminHttpServletExtensionResponseDefault(ctx, readResponse.DelegatedAdminHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.QuickstartHttpServletExtensionResponse != nil {
		readQuickstartHttpServletExtensionResponseDefault(ctx, readResponse.QuickstartHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AvailabilityStateHttpServletExtensionResponse != nil {
		readAvailabilityStateHttpServletExtensionResponseDefault(ctx, readResponse.AvailabilityStateHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PrometheusMonitoringHttpServletExtensionResponse != nil {
		readPrometheusMonitoringHttpServletExtensionResponseDefault(ctx, readResponse.PrometheusMonitoringHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.VelocityHttpServletExtensionResponse != nil {
		readVelocityHttpServletExtensionResponseDefault(ctx, readResponse.VelocityHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ConsentHttpServletExtensionResponse != nil {
		readConsentHttpServletExtensionResponseDefault(ctx, readResponse.ConsentHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LdapMappedScimHttpServletExtensionResponse != nil {
		readLdapMappedScimHttpServletExtensionResponseDefault(ctx, readResponse.LdapMappedScimHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedHttpServletExtensionResponse != nil {
		readGroovyScriptedHttpServletExtensionResponseDefault(ctx, readResponse.GroovyScriptedHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileServerHttpServletExtensionResponse != nil {
		readFileServerHttpServletExtensionResponseDefault(ctx, readResponse.FileServerHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)
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
	if readResponse.ThirdPartyHttpServletExtensionResponse != nil {
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
		if updateResponse.DelegatedAdminHttpServletExtensionResponse != nil {
			readDelegatedAdminHttpServletExtensionResponseDefault(ctx, updateResponse.DelegatedAdminHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.QuickstartHttpServletExtensionResponse != nil {
			readQuickstartHttpServletExtensionResponseDefault(ctx, updateResponse.QuickstartHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AvailabilityStateHttpServletExtensionResponse != nil {
			readAvailabilityStateHttpServletExtensionResponseDefault(ctx, updateResponse.AvailabilityStateHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PrometheusMonitoringHttpServletExtensionResponse != nil {
			readPrometheusMonitoringHttpServletExtensionResponseDefault(ctx, updateResponse.PrometheusMonitoringHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.VelocityHttpServletExtensionResponse != nil {
			readVelocityHttpServletExtensionResponseDefault(ctx, updateResponse.VelocityHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConsentHttpServletExtensionResponse != nil {
			readConsentHttpServletExtensionResponseDefault(ctx, updateResponse.ConsentHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LdapMappedScimHttpServletExtensionResponse != nil {
			readLdapMappedScimHttpServletExtensionResponseDefault(ctx, updateResponse.LdapMappedScimHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedHttpServletExtensionResponse != nil {
			readGroovyScriptedHttpServletExtensionResponseDefault(ctx, updateResponse.GroovyScriptedHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileServerHttpServletExtensionResponse != nil {
			readFileServerHttpServletExtensionResponseDefault(ctx, updateResponse.FileServerHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConfigHttpServletExtensionResponse != nil {
			readConfigHttpServletExtensionResponseDefault(ctx, updateResponse.ConfigHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Scim2HttpServletExtensionResponse != nil {
			readScim2HttpServletExtensionResponseDefault(ctx, updateResponse.Scim2HttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DirectoryRestApiHttpServletExtensionResponse != nil {
			readDirectoryRestApiHttpServletExtensionResponseDefault(ctx, updateResponse.DirectoryRestApiHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyHttpServletExtensionResponse != nil {
			readThirdPartyHttpServletExtensionResponseDefault(ctx, updateResponse.ThirdPartyHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
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
		if httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Http Servlet Extension", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Http Servlet Extension", err, httpResp)
		}
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
		if updateResponse.QuickstartHttpServletExtensionResponse != nil {
			readQuickstartHttpServletExtensionResponse(ctx, updateResponse.QuickstartHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AvailabilityStateHttpServletExtensionResponse != nil {
			readAvailabilityStateHttpServletExtensionResponse(ctx, updateResponse.AvailabilityStateHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PrometheusMonitoringHttpServletExtensionResponse != nil {
			readPrometheusMonitoringHttpServletExtensionResponse(ctx, updateResponse.PrometheusMonitoringHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LdapMappedScimHttpServletExtensionResponse != nil {
			readLdapMappedScimHttpServletExtensionResponse(ctx, updateResponse.LdapMappedScimHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedHttpServletExtensionResponse != nil {
			readGroovyScriptedHttpServletExtensionResponse(ctx, updateResponse.GroovyScriptedHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileServerHttpServletExtensionResponse != nil {
			readFileServerHttpServletExtensionResponse(ctx, updateResponse.FileServerHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyHttpServletExtensionResponse != nil {
			readThirdPartyHttpServletExtensionResponse(ctx, updateResponse.ThirdPartyHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
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
		if updateResponse.DelegatedAdminHttpServletExtensionResponse != nil {
			readDelegatedAdminHttpServletExtensionResponseDefault(ctx, updateResponse.DelegatedAdminHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.QuickstartHttpServletExtensionResponse != nil {
			readQuickstartHttpServletExtensionResponseDefault(ctx, updateResponse.QuickstartHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AvailabilityStateHttpServletExtensionResponse != nil {
			readAvailabilityStateHttpServletExtensionResponseDefault(ctx, updateResponse.AvailabilityStateHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PrometheusMonitoringHttpServletExtensionResponse != nil {
			readPrometheusMonitoringHttpServletExtensionResponseDefault(ctx, updateResponse.PrometheusMonitoringHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.VelocityHttpServletExtensionResponse != nil {
			readVelocityHttpServletExtensionResponseDefault(ctx, updateResponse.VelocityHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConsentHttpServletExtensionResponse != nil {
			readConsentHttpServletExtensionResponseDefault(ctx, updateResponse.ConsentHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LdapMappedScimHttpServletExtensionResponse != nil {
			readLdapMappedScimHttpServletExtensionResponseDefault(ctx, updateResponse.LdapMappedScimHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedHttpServletExtensionResponse != nil {
			readGroovyScriptedHttpServletExtensionResponseDefault(ctx, updateResponse.GroovyScriptedHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileServerHttpServletExtensionResponse != nil {
			readFileServerHttpServletExtensionResponseDefault(ctx, updateResponse.FileServerHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConfigHttpServletExtensionResponse != nil {
			readConfigHttpServletExtensionResponseDefault(ctx, updateResponse.ConfigHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Scim2HttpServletExtensionResponse != nil {
			readScim2HttpServletExtensionResponseDefault(ctx, updateResponse.Scim2HttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DirectoryRestApiHttpServletExtensionResponse != nil {
			readDirectoryRestApiHttpServletExtensionResponseDefault(ctx, updateResponse.DirectoryRestApiHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyHttpServletExtensionResponse != nil {
			readThirdPartyHttpServletExtensionResponseDefault(ctx, updateResponse.ThirdPartyHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
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
	if err != nil && httpResp.StatusCode != 404 {
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
