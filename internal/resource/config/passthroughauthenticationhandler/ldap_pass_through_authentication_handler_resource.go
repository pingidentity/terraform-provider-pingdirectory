package passthroughauthenticationhandler

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
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
	_ resource.Resource                = &ldapPassThroughAuthenticationHandlerResource{}
	_ resource.ResourceWithConfigure   = &ldapPassThroughAuthenticationHandlerResource{}
	_ resource.ResourceWithImportState = &ldapPassThroughAuthenticationHandlerResource{}
	_ resource.Resource                = &defaultLdapPassThroughAuthenticationHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultLdapPassThroughAuthenticationHandlerResource{}
	_ resource.ResourceWithImportState = &defaultLdapPassThroughAuthenticationHandlerResource{}
)

// Create a Ldap Pass Through Authentication Handler resource
func NewLdapPassThroughAuthenticationHandlerResource() resource.Resource {
	return &ldapPassThroughAuthenticationHandlerResource{}
}

func NewDefaultLdapPassThroughAuthenticationHandlerResource() resource.Resource {
	return &defaultLdapPassThroughAuthenticationHandlerResource{}
}

// ldapPassThroughAuthenticationHandlerResource is the resource implementation.
type ldapPassThroughAuthenticationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultLdapPassThroughAuthenticationHandlerResource is the resource implementation.
type defaultLdapPassThroughAuthenticationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *ldapPassThroughAuthenticationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ldap_pass_through_authentication_handler"
}

func (r *defaultLdapPassThroughAuthenticationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_ldap_pass_through_authentication_handler"
}

// Configure adds the provider configured client to the resource.
func (r *ldapPassThroughAuthenticationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultLdapPassThroughAuthenticationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type ldapPassThroughAuthenticationHandlerResourceModel struct {
	Id                                 types.String `tfsdk:"id"`
	LastUpdated                        types.String `tfsdk:"last_updated"`
	Notifications                      types.Set    `tfsdk:"notifications"`
	RequiredActions                    types.Set    `tfsdk:"required_actions"`
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
func (r *ldapPassThroughAuthenticationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ldapPassThroughAuthenticationHandlerSchema(ctx, req, resp, false)
}

func (r *defaultLdapPassThroughAuthenticationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ldapPassThroughAuthenticationHandlerSchema(ctx, req, resp, true)
}

func ldapPassThroughAuthenticationHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Ldap Pass Through Authentication Handler.",
		Attributes: map[string]schema.Attribute{
			"server": schema.SetAttribute{
				Description: "Specifies the LDAP external server(s) to which authentication attempts should be forwarded.",
				Required:    true,
				ElementType: types.StringType,
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
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
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
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalLdapPassThroughAuthenticationHandlerFields(ctx context.Context, addRequest *client.AddLdapPassThroughAuthenticationHandlerRequest, plan ldapPassThroughAuthenticationHandlerResourceModel) error {
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

// Read a LdapPassThroughAuthenticationHandlerResponse object into the model struct
func readLdapPassThroughAuthenticationHandlerResponse(ctx context.Context, r *client.LdapPassThroughAuthenticationHandlerResponse, state *ldapPassThroughAuthenticationHandlerResourceModel, expectedValues *ldapPassThroughAuthenticationHandlerResourceModel, diagnostics *diag.Diagnostics) {
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
}

// Create any update operations necessary to make the state match the plan
func createLdapPassThroughAuthenticationHandlerOperations(plan ldapPassThroughAuthenticationHandlerResourceModel, state ldapPassThroughAuthenticationHandlerResourceModel) []client.Operation {
	var ops []client.Operation
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

// Create a new resource
func (r *ldapPassThroughAuthenticationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ldapPassThroughAuthenticationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var ServerSlice []string
	plan.Server.ElementsAs(ctx, &ServerSlice, false)
	addRequest := client.NewAddLdapPassThroughAuthenticationHandlerRequest(plan.Id.ValueString(),
		[]client.EnumldapPassThroughAuthenticationHandlerSchemaUrn{client.ENUMLDAPPASSTHROUGHAUTHENTICATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASS_THROUGH_AUTHENTICATION_HANDLERLDAP},
		ServerSlice)
	err := addOptionalLdapPassThroughAuthenticationHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Ldap Pass Through Authentication Handler", err.Error())
		return
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Ldap Pass Through Authentication Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state ldapPassThroughAuthenticationHandlerResourceModel
	readLdapPassThroughAuthenticationHandlerResponse(ctx, addResponse.LdapPassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultLdapPassThroughAuthenticationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ldapPassThroughAuthenticationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PassThroughAuthenticationHandlerApi.GetPassThroughAuthenticationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ldap Pass Through Authentication Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state ldapPassThroughAuthenticationHandlerResourceModel
	readLdapPassThroughAuthenticationHandlerResponse(ctx, readResponse.LdapPassThroughAuthenticationHandlerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.PassThroughAuthenticationHandlerApi.UpdatePassThroughAuthenticationHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createLdapPassThroughAuthenticationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PassThroughAuthenticationHandlerApi.UpdatePassThroughAuthenticationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ldap Pass Through Authentication Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLdapPassThroughAuthenticationHandlerResponse(ctx, updateResponse.LdapPassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *ldapPassThroughAuthenticationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLdapPassThroughAuthenticationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLdapPassThroughAuthenticationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLdapPassThroughAuthenticationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readLdapPassThroughAuthenticationHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state ldapPassThroughAuthenticationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.PassThroughAuthenticationHandlerApi.GetPassThroughAuthenticationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ldap Pass Through Authentication Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readLdapPassThroughAuthenticationHandlerResponse(ctx, readResponse.LdapPassThroughAuthenticationHandlerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *ldapPassThroughAuthenticationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLdapPassThroughAuthenticationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLdapPassThroughAuthenticationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLdapPassThroughAuthenticationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateLdapPassThroughAuthenticationHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan ldapPassThroughAuthenticationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state ldapPassThroughAuthenticationHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.PassThroughAuthenticationHandlerApi.UpdatePassThroughAuthenticationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createLdapPassThroughAuthenticationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.PassThroughAuthenticationHandlerApi.UpdatePassThroughAuthenticationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ldap Pass Through Authentication Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLdapPassThroughAuthenticationHandlerResponse(ctx, updateResponse.LdapPassThroughAuthenticationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultLdapPassThroughAuthenticationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *ldapPassThroughAuthenticationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ldapPassThroughAuthenticationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PassThroughAuthenticationHandlerApi.DeletePassThroughAuthenticationHandlerExecute(r.apiClient.PassThroughAuthenticationHandlerApi.DeletePassThroughAuthenticationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Ldap Pass Through Authentication Handler", err, httpResp)
		return
	}
}

func (r *ldapPassThroughAuthenticationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLdapPassThroughAuthenticationHandler(ctx, req, resp)
}

func (r *defaultLdapPassThroughAuthenticationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLdapPassThroughAuthenticationHandler(ctx, req, resp)
}

func importLdapPassThroughAuthenticationHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
