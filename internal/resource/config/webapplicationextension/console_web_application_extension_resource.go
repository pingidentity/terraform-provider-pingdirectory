package webapplicationextension

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
	_ resource.Resource                = &consoleWebApplicationExtensionResource{}
	_ resource.ResourceWithConfigure   = &consoleWebApplicationExtensionResource{}
	_ resource.ResourceWithImportState = &consoleWebApplicationExtensionResource{}
)

// Create a Console Web Application Extension resource
func NewConsoleWebApplicationExtensionResource() resource.Resource {
	return &consoleWebApplicationExtensionResource{}
}

// consoleWebApplicationExtensionResource is the resource implementation.
type consoleWebApplicationExtensionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *consoleWebApplicationExtensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_console_web_application_extension"
}

// Configure adds the provider configured client to the resource.
func (r *consoleWebApplicationExtensionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type consoleWebApplicationExtensionResourceModel struct {
	Id                                  types.String `tfsdk:"id"`
	LastUpdated                         types.String `tfsdk:"last_updated"`
	Notifications                       types.Set    `tfsdk:"notifications"`
	RequiredActions                     types.Set    `tfsdk:"required_actions"`
	SsoEnabled                          types.Bool   `tfsdk:"sso_enabled"`
	OidcClientID                        types.String `tfsdk:"oidc_client_id"`
	OidcClientSecret                    types.String `tfsdk:"oidc_client_secret"`
	OidcClientSecretPassphraseProvider  types.String `tfsdk:"oidc_client_secret_passphrase_provider"`
	OidcIssuerURL                       types.String `tfsdk:"oidc_issuer_url"`
	OidcTrustStoreFile                  types.String `tfsdk:"oidc_trust_store_file"`
	OidcTrustStoreType                  types.String `tfsdk:"oidc_trust_store_type"`
	OidcTrustStorePinPassphraseProvider types.String `tfsdk:"oidc_trust_store_pin_passphrase_provider"`
	OidcStrictHostnameVerification      types.Bool   `tfsdk:"oidc_strict_hostname_verification"`
	OidcTrustAll                        types.Bool   `tfsdk:"oidc_trust_all"`
	LdapServer                          types.String `tfsdk:"ldap_server"`
	TrustStoreFile                      types.String `tfsdk:"trust_store_file"`
	TrustStoreType                      types.String `tfsdk:"trust_store_type"`
	TrustStorePinPassphraseProvider     types.String `tfsdk:"trust_store_pin_passphrase_provider"`
	LogFile                             types.String `tfsdk:"log_file"`
	Complexity                          types.String `tfsdk:"complexity"`
	Description                         types.String `tfsdk:"description"`
	BaseContextPath                     types.String `tfsdk:"base_context_path"`
	WarFile                             types.String `tfsdk:"war_file"`
	DocumentRootDirectory               types.String `tfsdk:"document_root_directory"`
	DeploymentDescriptorFile            types.String `tfsdk:"deployment_descriptor_file"`
	TemporaryDirectory                  types.String `tfsdk:"temporary_directory"`
	InitParameter                       types.Set    `tfsdk:"init_parameter"`
}

// GetSchema defines the schema for the resource.
func (r *consoleWebApplicationExtensionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Console Web Application Extension.",
		Attributes: map[string]schema.Attribute{
			"sso_enabled": schema.BoolAttribute{
				Description: "Indicates that SSO login into the Administrative Console is enabled.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"oidc_client_id": schema.StringAttribute{
				Description: "The client ID to use when authenticating to the OpenID Connect provider.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"oidc_client_secret": schema.StringAttribute{
				Description: "The client secret to use when authenticating to the OpenID Connect provider.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Sensitive: true,
			},
			"oidc_client_secret_passphrase_provider": schema.StringAttribute{
				Description: "A passphrase provider that may be used to obtain the client secret to use when authenticating to the OpenID Connect provider.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"oidc_issuer_url": schema.StringAttribute{
				Description: "The issuer URL of the OpenID Connect provider.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"oidc_trust_store_file": schema.StringAttribute{
				Description: "Specifies the path to the truststore file used by this application to evaluate OIDC provider certificates. If this field is left blank, the default JVM trust store will be used.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"oidc_trust_store_type": schema.StringAttribute{
				Description: "Specifies the format for the data in the OIDC trust store file.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"oidc_trust_store_pin_passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider that may be used to obtain the PIN for the trust store used with OIDC providers. This is only required if a trust store file is required, and if that trust store requires a PIN to access its contents.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"oidc_strict_hostname_verification": schema.BoolAttribute{
				Description: "Controls whether or not hostname verification is performed, which checks if the hostname of the OIDC provider matches the name(s) stored inside the certificate it provides. This property should only be set to false for testing purposes.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"oidc_trust_all": schema.BoolAttribute{
				Description: "Controls whether or not this application will always trust any certificate that is presented to it, regardless of its contents. This property should only be set to true for testing purposes.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ldap_server": schema.StringAttribute{
				Description: "The LDAP URL used to connect to the managed server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"trust_store_file": schema.StringAttribute{
				Description: "Specifies the path to the truststore file, which is used by this application to establish trust of managed servers.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"trust_store_type": schema.StringAttribute{
				Description: "Specifies the format for the data in the trust store file.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"trust_store_pin_passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider that may be used to obtain the PIN for the trust store used with managed LDAP servers. This is only required if a trust store file is required, and if that trust store requires a PIN to access its contents.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"log_file": schema.StringAttribute{
				Description: "The path to the log file for the web application.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"complexity": schema.StringAttribute{
				Description: "Specifies the maximum complexity level for managed configuration elements.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Web Application Extension",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"base_context_path": schema.StringAttribute{
				Description: "Specifies the base context path that should be used by HTTP clients to reference content. The value must start with a forward slash and at least one additional character and must represent a valid HTTP context path.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"war_file": schema.StringAttribute{
				Description: "Specifies the path to a standard web application archive (WAR) file.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"document_root_directory": schema.StringAttribute{
				Description: "Specifies the path to the directory on the local filesystem containing the files to be served by this Web Application Extension. The path must exist, and it must be a directory.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"deployment_descriptor_file": schema.StringAttribute{
				Description: "Specifies the path to the deployment descriptor file when used with document-root-directory.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"temporary_directory": schema.StringAttribute{
				Description: "Specifies the path to the directory that may be used to store temporary files such as extracted WAR files and compiled JSP files.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"init_parameter": schema.SetAttribute{
				Description: "Specifies an initialization parameter to pass into the web application during startup.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Read a ConsoleWebApplicationExtensionResponse object into the model struct
func readConsoleWebApplicationExtensionResponse(ctx context.Context, r *client.ConsoleWebApplicationExtensionResponse, state *consoleWebApplicationExtensionResourceModel, expectedValues *consoleWebApplicationExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.SsoEnabled = internaltypes.BoolTypeOrNil(r.SsoEnabled)
	state.OidcClientID = internaltypes.StringTypeOrNil(r.OidcClientID, true)
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.OidcClientSecret = expectedValues.OidcClientSecret
	state.OidcClientSecretPassphraseProvider = internaltypes.StringTypeOrNil(r.OidcClientSecretPassphraseProvider, true)
	state.OidcIssuerURL = internaltypes.StringTypeOrNil(r.OidcIssuerURL, true)
	state.OidcTrustStoreFile = internaltypes.StringTypeOrNil(r.OidcTrustStoreFile, true)
	state.OidcTrustStoreType = internaltypes.StringTypeOrNil(r.OidcTrustStoreType, true)
	state.OidcTrustStorePinPassphraseProvider = internaltypes.StringTypeOrNil(r.OidcTrustStorePinPassphraseProvider, true)
	state.OidcStrictHostnameVerification = internaltypes.BoolTypeOrNil(r.OidcStrictHostnameVerification)
	state.OidcTrustAll = internaltypes.BoolTypeOrNil(r.OidcTrustAll)
	state.LdapServer = internaltypes.StringTypeOrNil(r.LdapServer, true)
	state.TrustStoreFile = internaltypes.StringTypeOrNil(r.TrustStoreFile, true)
	state.TrustStoreType = internaltypes.StringTypeOrNil(r.TrustStoreType, true)
	state.TrustStorePinPassphraseProvider = internaltypes.StringTypeOrNil(r.TrustStorePinPassphraseProvider, true)
	state.LogFile = internaltypes.StringTypeOrNil(r.LogFile, true)
	state.Complexity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumwebApplicationExtensionComplexityProp(r.Complexity), true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.WarFile = internaltypes.StringTypeOrNil(r.WarFile, true)
	state.DocumentRootDirectory = internaltypes.StringTypeOrNil(r.DocumentRootDirectory, true)
	state.DeploymentDescriptorFile = internaltypes.StringTypeOrNil(r.DeploymentDescriptorFile, true)
	state.TemporaryDirectory = internaltypes.StringTypeOrNil(r.TemporaryDirectory, true)
	state.InitParameter = internaltypes.GetStringSet(r.InitParameter)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createConsoleWebApplicationExtensionOperations(plan consoleWebApplicationExtensionResourceModel, state consoleWebApplicationExtensionResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddBoolOperationIfNecessary(&ops, plan.SsoEnabled, state.SsoEnabled, "sso-enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.OidcClientID, state.OidcClientID, "oidc-client-id")
	operations.AddStringOperationIfNecessary(&ops, plan.OidcClientSecret, state.OidcClientSecret, "oidc-client-secret")
	operations.AddStringOperationIfNecessary(&ops, plan.OidcClientSecretPassphraseProvider, state.OidcClientSecretPassphraseProvider, "oidc-client-secret-passphrase-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.OidcIssuerURL, state.OidcIssuerURL, "oidc-issuer-url")
	operations.AddStringOperationIfNecessary(&ops, plan.OidcTrustStoreFile, state.OidcTrustStoreFile, "oidc-trust-store-file")
	operations.AddStringOperationIfNecessary(&ops, plan.OidcTrustStoreType, state.OidcTrustStoreType, "oidc-trust-store-type")
	operations.AddStringOperationIfNecessary(&ops, plan.OidcTrustStorePinPassphraseProvider, state.OidcTrustStorePinPassphraseProvider, "oidc-trust-store-pin-passphrase-provider")
	operations.AddBoolOperationIfNecessary(&ops, plan.OidcStrictHostnameVerification, state.OidcStrictHostnameVerification, "oidc-strict-hostname-verification")
	operations.AddBoolOperationIfNecessary(&ops, plan.OidcTrustAll, state.OidcTrustAll, "oidc-trust-all")
	operations.AddStringOperationIfNecessary(&ops, plan.LdapServer, state.LdapServer, "ldap-server")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStoreFile, state.TrustStoreFile, "trust-store-file")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStoreType, state.TrustStoreType, "trust-store-type")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStorePinPassphraseProvider, state.TrustStorePinPassphraseProvider, "trust-store-pin-passphrase-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFile, state.LogFile, "log-file")
	operations.AddStringOperationIfNecessary(&ops, plan.Complexity, state.Complexity, "complexity")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.BaseContextPath, state.BaseContextPath, "base-context-path")
	operations.AddStringOperationIfNecessary(&ops, plan.WarFile, state.WarFile, "war-file")
	operations.AddStringOperationIfNecessary(&ops, plan.DocumentRootDirectory, state.DocumentRootDirectory, "document-root-directory")
	operations.AddStringOperationIfNecessary(&ops, plan.DeploymentDescriptorFile, state.DeploymentDescriptorFile, "deployment-descriptor-file")
	operations.AddStringOperationIfNecessary(&ops, plan.TemporaryDirectory, state.TemporaryDirectory, "temporary-directory")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.InitParameter, state.InitParameter, "init-parameter")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *consoleWebApplicationExtensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan consoleWebApplicationExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.WebApplicationExtensionApi.GetWebApplicationExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Console Web Application Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state consoleWebApplicationExtensionResourceModel
	readConsoleWebApplicationExtensionResponse(ctx, readResponse.ConsoleWebApplicationExtensionResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.WebApplicationExtensionApi.UpdateWebApplicationExtension(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createConsoleWebApplicationExtensionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.WebApplicationExtensionApi.UpdateWebApplicationExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Console Web Application Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readConsoleWebApplicationExtensionResponse(ctx, updateResponse.ConsoleWebApplicationExtensionResponse, &state, &plan, &resp.Diagnostics)
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
func (r *consoleWebApplicationExtensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state consoleWebApplicationExtensionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.WebApplicationExtensionApi.GetWebApplicationExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Console Web Application Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readConsoleWebApplicationExtensionResponse(ctx, readResponse.ConsoleWebApplicationExtensionResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *consoleWebApplicationExtensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan consoleWebApplicationExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state consoleWebApplicationExtensionResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.WebApplicationExtensionApi.UpdateWebApplicationExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createConsoleWebApplicationExtensionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.WebApplicationExtensionApi.UpdateWebApplicationExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Console Web Application Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readConsoleWebApplicationExtensionResponse(ctx, updateResponse.ConsoleWebApplicationExtensionResponse, &state, &plan, &resp.Diagnostics)
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
func (r *consoleWebApplicationExtensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *consoleWebApplicationExtensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
