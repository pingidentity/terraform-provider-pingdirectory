package saslmechanismhandler

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
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
	_ resource.Resource                = &gssapiSaslMechanismHandlerResource{}
	_ resource.ResourceWithConfigure   = &gssapiSaslMechanismHandlerResource{}
	_ resource.ResourceWithImportState = &gssapiSaslMechanismHandlerResource{}
)

// Create a Gssapi Sasl Mechanism Handler resource
func NewGssapiSaslMechanismHandlerResource() resource.Resource {
	return &gssapiSaslMechanismHandlerResource{}
}

// gssapiSaslMechanismHandlerResource is the resource implementation.
type gssapiSaslMechanismHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *gssapiSaslMechanismHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_gssapi_sasl_mechanism_handler"
}

// Configure adds the provider configured client to the resource.
func (r *gssapiSaslMechanismHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type gssapiSaslMechanismHandlerResourceModel struct {
	Id                                   types.String `tfsdk:"id"`
	LastUpdated                          types.String `tfsdk:"last_updated"`
	Notifications                        types.Set    `tfsdk:"notifications"`
	RequiredActions                      types.Set    `tfsdk:"required_actions"`
	Realm                                types.String `tfsdk:"realm"`
	KdcAddress                           types.String `tfsdk:"kdc_address"`
	Keytab                               types.String `tfsdk:"keytab"`
	AllowNullServerFqdn                  types.Bool   `tfsdk:"allow_null_server_fqdn"`
	ServerFqdn                           types.String `tfsdk:"server_fqdn"`
	AllowedQualityOfProtection           types.Set    `tfsdk:"allowed_quality_of_protection"`
	IdentityMapper                       types.String `tfsdk:"identity_mapper"`
	AlternateAuthorizationIdentityMapper types.String `tfsdk:"alternate_authorization_identity_mapper"`
	KerberosServicePrincipal             types.String `tfsdk:"kerberos_service_principal"`
	GssapiRole                           types.String `tfsdk:"gssapi_role"`
	JaasConfigFile                       types.String `tfsdk:"jaas_config_file"`
	EnableDebug                          types.Bool   `tfsdk:"enable_debug"`
	Description                          types.String `tfsdk:"description"`
	Enabled                              types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *gssapiSaslMechanismHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Gssapi Sasl Mechanism Handler.",
		Attributes: map[string]schema.Attribute{
			"realm": schema.StringAttribute{
				Description: "Specifies the realm to be used for GSSAPI authentication.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"kdc_address": schema.StringAttribute{
				Description: "Specifies the address of the KDC that is to be used for Kerberos processing.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"keytab": schema.StringAttribute{
				Description: "Specifies the keytab file that should be used for Kerberos processing.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_null_server_fqdn": schema.BoolAttribute{
				Description: "Specifies whether or not to allow a null value for the server-fqdn.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"server_fqdn": schema.StringAttribute{
				Description: "Specifies the DNS-resolvable fully-qualified domain name for the system.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_quality_of_protection": schema.SetAttribute{
				Description: "Specifies the supported quality of protection (QoP) levels that clients will be permitted to request when performing GSSAPI authentication.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"identity_mapper": schema.StringAttribute{
				Description: "Specifies the name of the identity mapper that is to be used with this SASL mechanism handler to match the Kerberos principal included in the SASL bind request to the corresponding user in the directory.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"alternate_authorization_identity_mapper": schema.StringAttribute{
				Description: "Specifies the name of the identity mapper that is to be used with this SASL mechanism handler to map the alternate authorization identity (if provided, and if different from the Kerberos principal used as the authentication identity) to the corresponding user in the directory. If no value is specified, then the mapper specified in the identity-mapper configuration property will be used.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"kerberos_service_principal": schema.StringAttribute{
				Description: "Specifies the Kerberos service principal that the Directory Server will use to identify itself to the KDC.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gssapi_role": schema.StringAttribute{
				Description: "Specifies the role that should be declared for the server in the generated JAAS configuration file.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"jaas_config_file": schema.StringAttribute{
				Description: "Specifies the path to a JAAS (Java Authentication and Authorization Service) configuration file that provides the information that the JVM should use for Kerberos processing.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enable_debug": schema.BoolAttribute{
				Description: "Indicates whether to enable debugging for the Java GSSAPI provider. Debug information will be written to standard output, which should be captured in the server.out log file.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this SASL Mechanism Handler",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the SASL mechanism handler is enabled for use.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Read a GssapiSaslMechanismHandlerResponse object into the model struct
func readGssapiSaslMechanismHandlerResponse(ctx context.Context, r *client.GssapiSaslMechanismHandlerResponse, state *gssapiSaslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Realm = internaltypes.StringTypeOrNil(r.Realm, true)
	state.KdcAddress = internaltypes.StringTypeOrNil(r.KdcAddress, true)
	state.Keytab = internaltypes.StringTypeOrNil(r.Keytab, true)
	state.AllowNullServerFqdn = internaltypes.BoolTypeOrNil(r.AllowNullServerFqdn)
	state.ServerFqdn = internaltypes.StringTypeOrNil(r.ServerFqdn, true)
	state.AllowedQualityOfProtection = internaltypes.GetStringSet(
		client.StringSliceEnumsaslMechanismHandlerAllowedQualityOfProtectionProp(r.AllowedQualityOfProtection))
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.AlternateAuthorizationIdentityMapper = internaltypes.StringTypeOrNil(r.AlternateAuthorizationIdentityMapper, true)
	state.KerberosServicePrincipal = internaltypes.StringTypeOrNil(r.KerberosServicePrincipal, true)
	state.GssapiRole = internaltypes.StringTypeOrNil(
		client.StringPointerEnumsaslMechanismHandlerGssapiRoleProp(r.GssapiRole), true)
	state.JaasConfigFile = internaltypes.StringTypeOrNil(r.JaasConfigFile, true)
	state.EnableDebug = internaltypes.BoolTypeOrNil(r.EnableDebug)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createGssapiSaslMechanismHandlerOperations(plan gssapiSaslMechanismHandlerResourceModel, state gssapiSaslMechanismHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Realm, state.Realm, "realm")
	operations.AddStringOperationIfNecessary(&ops, plan.KdcAddress, state.KdcAddress, "kdc-address")
	operations.AddStringOperationIfNecessary(&ops, plan.Keytab, state.Keytab, "keytab")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowNullServerFqdn, state.AllowNullServerFqdn, "allow-null-server-fqdn")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerFqdn, state.ServerFqdn, "server-fqdn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedQualityOfProtection, state.AllowedQualityOfProtection, "allowed-quality-of-protection")
	operations.AddStringOperationIfNecessary(&ops, plan.IdentityMapper, state.IdentityMapper, "identity-mapper")
	operations.AddStringOperationIfNecessary(&ops, plan.AlternateAuthorizationIdentityMapper, state.AlternateAuthorizationIdentityMapper, "alternate-authorization-identity-mapper")
	operations.AddStringOperationIfNecessary(&ops, plan.KerberosServicePrincipal, state.KerberosServicePrincipal, "kerberos-service-principal")
	operations.AddStringOperationIfNecessary(&ops, plan.GssapiRole, state.GssapiRole, "gssapi-role")
	operations.AddStringOperationIfNecessary(&ops, plan.JaasConfigFile, state.JaasConfigFile, "jaas-config-file")
	operations.AddBoolOperationIfNecessary(&ops, plan.EnableDebug, state.EnableDebug, "enable-debug")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *gssapiSaslMechanismHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan gssapiSaslMechanismHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SaslMechanismHandlerApi.GetSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Gssapi Sasl Mechanism Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state gssapiSaslMechanismHandlerResourceModel
	readGssapiSaslMechanismHandlerResponse(ctx, readResponse.GssapiSaslMechanismHandlerResponse, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.SaslMechanismHandlerApi.UpdateSaslMechanismHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createGssapiSaslMechanismHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.SaslMechanismHandlerApi.UpdateSaslMechanismHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Gssapi Sasl Mechanism Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGssapiSaslMechanismHandlerResponse(ctx, updateResponse.GssapiSaslMechanismHandlerResponse, &state, &resp.Diagnostics)
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
func (r *gssapiSaslMechanismHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state gssapiSaslMechanismHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SaslMechanismHandlerApi.GetSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Gssapi Sasl Mechanism Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readGssapiSaslMechanismHandlerResponse(ctx, readResponse.GssapiSaslMechanismHandlerResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *gssapiSaslMechanismHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan gssapiSaslMechanismHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state gssapiSaslMechanismHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.SaslMechanismHandlerApi.UpdateSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createGssapiSaslMechanismHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.SaslMechanismHandlerApi.UpdateSaslMechanismHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Gssapi Sasl Mechanism Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGssapiSaslMechanismHandlerResponse(ctx, updateResponse.GssapiSaslMechanismHandlerResponse, &state, &resp.Diagnostics)
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
func (r *gssapiSaslMechanismHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *gssapiSaslMechanismHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
