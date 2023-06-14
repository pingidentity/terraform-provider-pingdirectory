package passthroughauthenticationhandler

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
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &passThroughAuthenticationHandlerResource{}
	_ resource.ResourceWithConfigure   = &passThroughAuthenticationHandlerResource{}
	_ resource.ResourceWithImportState = &passThroughAuthenticationHandlerResource{}
	_ resource.Resource                = &defaultPassThroughAuthenticationHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultPassThroughAuthenticationHandlerResource{}
	_ resource.ResourceWithImportState = &defaultPassThroughAuthenticationHandlerResource{}
)

// Create a Pass Through Authentication Handler resource
func NewPassThroughAuthenticationHandlerResource() resource.Resource {
	return &passThroughAuthenticationHandlerResource{}
}

func NewDefaultPassThroughAuthenticationHandlerResource() resource.Resource {
	return &defaultPassThroughAuthenticationHandlerResource{}
}

// passThroughAuthenticationHandlerResource is the resource implementation.
type passThroughAuthenticationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultPassThroughAuthenticationHandlerResource is the resource implementation.
type defaultPassThroughAuthenticationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *passThroughAuthenticationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pass_through_authentication_handler"
}

func (r *defaultPassThroughAuthenticationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_pass_through_authentication_handler"
}

// Configure adds the provider configured client to the resource.
func (r *passThroughAuthenticationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultPassThroughAuthenticationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type passThroughAuthenticationHandlerResourceModel struct {
	Id                                 types.String `tfsdk:"id"`
	LastUpdated                        types.String `tfsdk:"last_updated"`
	Notifications                      types.Set    `tfsdk:"notifications"`
	RequiredActions                    types.Set    `tfsdk:"required_actions"`
	Type                               types.String `tfsdk:"type"`
	ExtensionClass                     types.String `tfsdk:"extension_class"`
	ExtensionArgument                  types.Set    `tfsdk:"extension_argument"`
	Server                             types.Set    `tfsdk:"server"`
	ServerAccessMode                   types.String `tfsdk:"server_access_mode"`
	DnMap                              types.Set    `tfsdk:"dn_map"`
	BindDNPattern                      types.String `tfsdk:"bind_dn_pattern"`
	SearchBaseDN                       types.String `tfsdk:"search_base_dn"`
	SearchFilterPattern                types.String `tfsdk:"search_filter_pattern"`
	InitialConnections                 types.Int64  `tfsdk:"initial_connections"`
	MaxConnections                     types.Int64  `tfsdk:"max_connections"`
	UseLocation                        types.Bool   `tfsdk:"use_location"`
	MaximumAllowedLocalResponseTime    types.String `tfsdk:"maximum_allowed_local_response_time"`
	MaximumAllowedNonlocalResponseTime types.String `tfsdk:"maximum_allowed_nonlocal_response_time"`
	UsePasswordPolicyControl           types.Bool   `tfsdk:"use_password_policy_control"`
	Description                        types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *passThroughAuthenticationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	passThroughAuthenticationHandlerSchema(ctx, req, resp, false)
}

func (r *defaultPassThroughAuthenticationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	passThroughAuthenticationHandlerSchema(ctx, req, resp, true)
}

func passThroughAuthenticationHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Pass Through Authentication Handler.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Pass Through Authentication Handler resource. Options are ['ldap', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"ldap", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Pass Through Authentication Handler.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Pass Through Authentication Handler. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"server": schema.SetAttribute{
				Description: "Specifies the LDAP external server(s) to which authentication attempts should be forwarded.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"server_access_mode": schema.StringAttribute{
				Description: "Specifies the manner in which external servers should be used for pass-through authentication attempts if multiple servers are defined.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dn_map": schema.SetAttribute{
				Description: "Specifies one or more DN mappings that may be used to transform bind DNs before attempting to bind to the external servers.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"bind_dn_pattern": schema.StringAttribute{
				Description: "A pattern to use to construct the bind DN for the simple bind request to send to the remote server. This may consist of a combination of static text and attribute values and other directives enclosed in curly braces.  For example, the value \"cn={cn},ou=People,dc=example,dc=com\" indicates that the remote bind DN should be constructed from the text \"cn=\" followed by the value of the local entry's cn attribute followed by the text \"ou=People,dc=example,dc=com\". If an attribute contains the value to use as the bind DN for pass-through authentication, then the pattern may simply be the name of that attribute in curly braces (e.g., if the seeAlso attribute contains the bind DN for the target user, then a bind DN pattern of \"{seeAlso}\" would be appropriate).  Note that a bind DN pattern can be used to construct a bind DN that is not actually a valid LDAP distinguished name. For example, if authentication is being passed through to a Microsoft Active Directory server, then a bind DN pattern could be used to construct a user principal name (UPN) as an alternative to a distinguished name.",
				Optional:    true,
			},
			"search_base_dn": schema.StringAttribute{
				Description: "The base DN to use when searching for the user entry using a filter constructed from the pattern defined in the search-filter-pattern property. If no base DN is specified, the null DN will be used as the search base DN.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"search_filter_pattern": schema.StringAttribute{
				Description: "A pattern to use to construct a filter to use when searching an external server for the entry of the user as whom to bind. For example, \"(mail={uid:ldapFilterEscape}@example.com)\" would construct a search filter to search for a user whose entry in the local server contains a uid attribute whose value appears before \"@example.com\" in the mail attribute in the external server. Note that the \"ldapFilterEscape\" modifier should almost always be used with attributes specified in the pattern.",
				Optional:    true,
			},
			"initial_connections": schema.Int64Attribute{
				Description: "Specifies the initial number of connections to establish to each external server against which authentication may be attempted.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"max_connections": schema.Int64Attribute{
				Description: "Specifies the maximum number of connections to maintain to each external server against which authentication may be attempted. This value must be greater than or equal to the value for the initial-connections property.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"use_location": schema.BoolAttribute{
				Description: "Indicates whether to take server locations into account when prioritizing the servers to use for pass-through authentication attempts.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"maximum_allowed_local_response_time": schema.StringAttribute{
				Description: "The maximum length of time to wait for a response from an external server in the same location as this Directory Server before considering it unavailable.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"maximum_allowed_nonlocal_response_time": schema.StringAttribute{
				Description: "The maximum length of time to wait for a response from an external server in a different location from this Directory Server before considering it unavailable.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"use_password_policy_control": schema.BoolAttribute{
				Description: "Indicates whether to include the password policy request control (as defined in draft-behera-ldap-password-policy-10) in bind requests sent to the external server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Pass Through Authentication Handler",
				Optional:    true,
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Validators = []validator.String{
			stringvalidator.OneOf([]string{"ldap", "third-party"}...),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAllAttributesToOptionalAndComputed(&schemaDef, []string{"id"})
	}
	config.AddCommonSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan
func (r *passThroughAuthenticationHandlerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanPassThroughAuthenticationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPassThroughAuthenticationHandlerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanPassThroughAuthenticationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanPassThroughAuthenticationHandler(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var model passThroughAuthenticationHandlerResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.Server) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'server' not supported by pingdirectory_pass_through_authentication_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'server', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.InitialConnections) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'initial_connections' not supported by pingdirectory_pass_through_authentication_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'initial_connections', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.UseLocation) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'use_location' not supported by pingdirectory_pass_through_authentication_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'use_location', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.ServerAccessMode) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'server_access_mode' not supported by pingdirectory_pass_through_authentication_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'server_access_mode', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.SearchFilterPattern) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'search_filter_pattern' not supported by pingdirectory_pass_through_authentication_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'search_filter_pattern', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.MaxConnections) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'max_connections' not supported by pingdirectory_pass_through_authentication_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'max_connections', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.SearchBaseDN) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'search_base_dn' not supported by pingdirectory_pass_through_authentication_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'search_base_dn', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.MaximumAllowedLocalResponseTime) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'maximum_allowed_local_response_time' not supported by pingdirectory_pass_through_authentication_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'maximum_allowed_local_response_time', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.DnMap) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'dn_map' not supported by pingdirectory_pass_through_authentication_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'dn_map', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.MaximumAllowedNonlocalResponseTime) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'maximum_allowed_nonlocal_response_time' not supported by pingdirectory_pass_through_authentication_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'maximum_allowed_nonlocal_response_time', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.ExtensionArgument) && model.Type.ValueString() != "third-party" {
		resp.Diagnostics.AddError("Attribute 'extension_argument' not supported by pingdirectory_pass_through_authentication_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'extension_argument', the 'type' attribute must be one of ['third-party']")
	}
	if internaltypes.IsDefined(model.ExtensionClass) && model.Type.ValueString() != "third-party" {
		resp.Diagnostics.AddError("Attribute 'extension_class' not supported by pingdirectory_pass_through_authentication_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'extension_class', the 'type' attribute must be one of ['third-party']")
	}
	if internaltypes.IsDefined(model.UsePasswordPolicyControl) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'use_password_policy_control' not supported by pingdirectory_pass_through_authentication_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'use_password_policy_control', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.BindDNPattern) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'bind_dn_pattern' not supported by pingdirectory_pass_through_authentication_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'bind_dn_pattern', the 'type' attribute must be one of ['ldap']")
	}
}

// Add optional fields to create request for ldap pass-through-authentication-handler
func addOptionalLdapPassThroughAuthenticationHandlerFields(ctx context.Context, addRequest *client.AddLdapPassThroughAuthenticationHandlerRequest, plan passThroughAuthenticationHandlerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ServerAccessMode) {
		serverAccessMode, err := client.NewEnumpassThroughAuthenticationHandlerServerAccessModePropFromValue(plan.ServerAccessMode.ValueString())
		if err != nil {
			return err
		}
		addRequest.ServerAccessMode = serverAccessMode
	}
	if internaltypes.IsDefined(plan.DnMap) {
		var slice []string
		plan.DnMap.ElementsAs(ctx, &slice, false)
		addRequest.DnMap = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BindDNPattern) {
		addRequest.BindDNPattern = plan.BindDNPattern.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchBaseDN) {
		addRequest.SearchBaseDN = plan.SearchBaseDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchFilterPattern) {
		addRequest.SearchFilterPattern = plan.SearchFilterPattern.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InitialConnections) {
		addRequest.InitialConnections = plan.InitialConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaxConnections) {
		addRequest.MaxConnections = plan.MaxConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.UseLocation) {
		addRequest.UseLocation = plan.UseLocation.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaximumAllowedLocalResponseTime) {
		addRequest.MaximumAllowedLocalResponseTime = plan.MaximumAllowedLocalResponseTime.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaximumAllowedNonlocalResponseTime) {
		addRequest.MaximumAllowedNonlocalResponseTime = plan.MaximumAllowedNonlocalResponseTime.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.UsePasswordPolicyControl) {
		addRequest.UsePasswordPolicyControl = plan.UsePasswordPolicyControl.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for third-party pass-through-authentication-handler
func addOptionalThirdPartyPassThroughAuthenticationHandlerFields(ctx context.Context, addRequest *client.AddThirdPartyPassThroughAuthenticationHandlerRequest, plan passThroughAuthenticationHandlerResourceModel) error {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Populate any sets that have a nil ElementType, to avoid a nil pointer when setting the state
func populatePassThroughAuthenticationHandlerNilSets(ctx context.Context, model *passThroughAuthenticationHandlerResourceModel) {
	if model.DnMap.ElementType(ctx) == nil {
		model.DnMap = types.SetNull(types.StringType)
	}
	if model.Server.ElementType(ctx) == nil {
		model.Server = types.SetNull(types.StringType)
	}
	if model.ExtensionArgument.ElementType(ctx) == nil {
		model.ExtensionArgument = types.SetNull(types.StringType)
	}
}

// Read a LdapPassThroughAuthenticationHandlerResponse object into the model struct
func readLdapPassThroughAuthenticationHandlerResponse(ctx context.Context, r *client.LdapPassThroughAuthenticationHandlerResponse, state *passThroughAuthenticationHandlerResourceModel, expectedValues *passThroughAuthenticationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap")
	state.Id = types.StringValue(r.Id)
	state.Server = internaltypes.GetStringSet(r.Server)
	state.ServerAccessMode = types.StringValue(r.ServerAccessMode.String())
	state.DnMap = internaltypes.GetStringSet(r.DnMap)
	state.BindDNPattern = internaltypes.StringTypeOrNil(r.BindDNPattern, internaltypes.IsEmptyString(expectedValues.BindDNPattern))
	state.SearchBaseDN = internaltypes.StringTypeOrNil(r.SearchBaseDN, internaltypes.IsEmptyString(expectedValues.SearchBaseDN))
	state.SearchFilterPattern = internaltypes.StringTypeOrNil(r.SearchFilterPattern, internaltypes.IsEmptyString(expectedValues.SearchFilterPattern))
	state.InitialConnections = types.Int64Value(r.InitialConnections)
	state.MaxConnections = types.Int64Value(r.MaxConnections)
	state.UseLocation = internaltypes.BoolTypeOrNil(r.UseLocation)
	state.MaximumAllowedLocalResponseTime = internaltypes.StringTypeOrNil(r.MaximumAllowedLocalResponseTime, internaltypes.IsEmptyString(expectedValues.MaximumAllowedLocalResponseTime))
	config.CheckMismatchedPDFormattedAttributes("maximum_allowed_local_response_time",
		expectedValues.MaximumAllowedLocalResponseTime, state.MaximumAllowedLocalResponseTime, diagnostics)
	state.MaximumAllowedNonlocalResponseTime = internaltypes.StringTypeOrNil(r.MaximumAllowedNonlocalResponseTime, internaltypes.IsEmptyString(expectedValues.MaximumAllowedNonlocalResponseTime))
	config.CheckMismatchedPDFormattedAttributes("maximum_allowed_nonlocal_response_time",
		expectedValues.MaximumAllowedNonlocalResponseTime, state.MaximumAllowedNonlocalResponseTime, diagnostics)
	state.UsePasswordPolicyControl = internaltypes.BoolTypeOrNil(r.UsePasswordPolicyControl)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePassThroughAuthenticationHandlerNilSets(ctx, state)
}

// Read a ThirdPartyPassThroughAuthenticationHandlerResponse object into the model struct
func readThirdPartyPassThroughAuthenticationHandlerResponse(ctx context.Context, r *client.ThirdPartyPassThroughAuthenticationHandlerResponse, state *passThroughAuthenticationHandlerResourceModel, expectedValues *passThroughAuthenticationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePassThroughAuthenticationHandlerNilSets(ctx, state)
}

// Create any update operations necessary to make the state match the plan
func createPassThroughAuthenticationHandlerOperations(plan passThroughAuthenticationHandlerResourceModel, state passThroughAuthenticationHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.Server, state.Server, "server")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerAccessMode, state.ServerAccessMode, "server-access-mode")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DnMap, state.DnMap, "dn-map")
	operations.AddStringOperationIfNecessary(&ops, plan.BindDNPattern, state.BindDNPattern, "bind-dn-pattern")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchBaseDN, state.SearchBaseDN, "search-base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchFilterPattern, state.SearchFilterPattern, "search-filter-pattern")
	operations.AddInt64OperationIfNecessary(&ops, plan.InitialConnections, state.InitialConnections, "initial-connections")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxConnections, state.MaxConnections, "max-connections")
	operations.AddBoolOperationIfNecessary(&ops, plan.UseLocation, state.UseLocation, "use-location")
	operations.AddStringOperationIfNecessary(&ops, plan.MaximumAllowedLocalResponseTime, state.MaximumAllowedLocalResponseTime, "maximum-allowed-local-response-time")
	operations.AddStringOperationIfNecessary(&ops, plan.MaximumAllowedNonlocalResponseTime, state.MaximumAllowedNonlocalResponseTime, "maximum-allowed-nonlocal-response-time")
	operations.AddBoolOperationIfNecessary(&ops, plan.UsePasswordPolicyControl, state.UsePasswordPolicyControl, "use-password-policy-control")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a ldap pass-through-authentication-handler
func (r *passThroughAuthenticationHandlerResource) CreateLdapPassThroughAuthenticationHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passThroughAuthenticationHandlerResourceModel) (*passThroughAuthenticationHandlerResourceModel, error) {
	var ServerSlice []string
	plan.Server.ElementsAs(ctx, &ServerSlice, false)
	addRequest := client.NewAddLdapPassThroughAuthenticationHandlerRequest(plan.Id.ValueString(),
		[]client.EnumldapPassThroughAuthenticationHandlerSchemaUrn{client.ENUMLDAPPASSTHROUGHAUTHENTICATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASS_THROUGH_AUTHENTICATION_HANDLERLDAP},
		ServerSlice)
	err := addOptionalLdapPassThroughAuthenticationHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Pass Through Authentication Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PassThroughAuthenticationHandlerApi.AddPassThroughAuthenticationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPassThroughAuthenticationHandlerRequest(
		client.AddLdapPassThroughAuthenticationHandlerRequestAsAddPassThroughAuthenticationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PassThroughAuthenticationHandlerApi.AddPassThroughAuthenticationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Pass Through Authentication Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passThroughAuthenticationHandlerResourceModel
	readLdapPassThroughAuthenticationHandlerResponse(ctx, addResponse.LdapPassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party pass-through-authentication-handler
func (r *passThroughAuthenticationHandlerResource) CreateThirdPartyPassThroughAuthenticationHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passThroughAuthenticationHandlerResourceModel) (*passThroughAuthenticationHandlerResourceModel, error) {
	addRequest := client.NewAddThirdPartyPassThroughAuthenticationHandlerRequest(plan.Id.ValueString(),
		[]client.EnumthirdPartyPassThroughAuthenticationHandlerSchemaUrn{client.ENUMTHIRDPARTYPASSTHROUGHAUTHENTICATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASS_THROUGH_AUTHENTICATION_HANDLERTHIRD_PARTY},
		plan.ExtensionClass.ValueString())
	err := addOptionalThirdPartyPassThroughAuthenticationHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Pass Through Authentication Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PassThroughAuthenticationHandlerApi.AddPassThroughAuthenticationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPassThroughAuthenticationHandlerRequest(
		client.AddThirdPartyPassThroughAuthenticationHandlerRequestAsAddPassThroughAuthenticationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PassThroughAuthenticationHandlerApi.AddPassThroughAuthenticationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Pass Through Authentication Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passThroughAuthenticationHandlerResourceModel
	readThirdPartyPassThroughAuthenticationHandlerResponse(ctx, addResponse.ThirdPartyPassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *passThroughAuthenticationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan passThroughAuthenticationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *passThroughAuthenticationHandlerResourceModel
	var err error
	if plan.Type.ValueString() == "ldap" {
		state, err = r.CreateLdapPassThroughAuthenticationHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyPassThroughAuthenticationHandler(ctx, req, resp, plan)
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
func (r *defaultPassThroughAuthenticationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan passThroughAuthenticationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PassThroughAuthenticationHandlerApi.GetPassThroughAuthenticationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Pass Through Authentication Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state passThroughAuthenticationHandlerResourceModel
	if plan.Type.ValueString() == "ldap" {
		readLdapPassThroughAuthenticationHandlerResponse(ctx, readResponse.LdapPassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "third-party" {
		readThirdPartyPassThroughAuthenticationHandlerResponse(ctx, readResponse.ThirdPartyPassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.PassThroughAuthenticationHandlerApi.UpdatePassThroughAuthenticationHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createPassThroughAuthenticationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PassThroughAuthenticationHandlerApi.UpdatePassThroughAuthenticationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Pass Through Authentication Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "ldap" {
			readLdapPassThroughAuthenticationHandlerResponse(ctx, updateResponse.LdapPassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party" {
			readThirdPartyPassThroughAuthenticationHandlerResponse(ctx, updateResponse.ThirdPartyPassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *passThroughAuthenticationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPassThroughAuthenticationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPassThroughAuthenticationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPassThroughAuthenticationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readPassThroughAuthenticationHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state passThroughAuthenticationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.PassThroughAuthenticationHandlerApi.GetPassThroughAuthenticationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Pass Through Authentication Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.LdapPassThroughAuthenticationHandlerResponse != nil {
		readLdapPassThroughAuthenticationHandlerResponse(ctx, readResponse.LdapPassThroughAuthenticationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyPassThroughAuthenticationHandlerResponse != nil {
		readThirdPartyPassThroughAuthenticationHandlerResponse(ctx, readResponse.ThirdPartyPassThroughAuthenticationHandlerResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *passThroughAuthenticationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePassThroughAuthenticationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPassThroughAuthenticationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePassThroughAuthenticationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updatePassThroughAuthenticationHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan passThroughAuthenticationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state passThroughAuthenticationHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.PassThroughAuthenticationHandlerApi.UpdatePassThroughAuthenticationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createPassThroughAuthenticationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.PassThroughAuthenticationHandlerApi.UpdatePassThroughAuthenticationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Pass Through Authentication Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "ldap" {
			readLdapPassThroughAuthenticationHandlerResponse(ctx, updateResponse.LdapPassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party" {
			readThirdPartyPassThroughAuthenticationHandlerResponse(ctx, updateResponse.ThirdPartyPassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultPassThroughAuthenticationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *passThroughAuthenticationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state passThroughAuthenticationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PassThroughAuthenticationHandlerApi.DeletePassThroughAuthenticationHandlerExecute(r.apiClient.PassThroughAuthenticationHandlerApi.DeletePassThroughAuthenticationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Pass Through Authentication Handler", err, httpResp)
		return
	}
}

func (r *passThroughAuthenticationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPassThroughAuthenticationHandler(ctx, req, resp)
}

func (r *defaultPassThroughAuthenticationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPassThroughAuthenticationHandler(ctx, req, resp)
}

func importPassThroughAuthenticationHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
