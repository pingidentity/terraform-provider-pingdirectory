package datasecurityauditor

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
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &idleAccountDataSecurityAuditorResource{}
	_ resource.ResourceWithConfigure   = &idleAccountDataSecurityAuditorResource{}
	_ resource.ResourceWithImportState = &idleAccountDataSecurityAuditorResource{}
	_ resource.Resource                = &defaultIdleAccountDataSecurityAuditorResource{}
	_ resource.ResourceWithConfigure   = &defaultIdleAccountDataSecurityAuditorResource{}
	_ resource.ResourceWithImportState = &defaultIdleAccountDataSecurityAuditorResource{}
)

// Create a Idle Account Data Security Auditor resource
func NewIdleAccountDataSecurityAuditorResource() resource.Resource {
	return &idleAccountDataSecurityAuditorResource{}
}

func NewDefaultIdleAccountDataSecurityAuditorResource() resource.Resource {
	return &defaultIdleAccountDataSecurityAuditorResource{}
}

// idleAccountDataSecurityAuditorResource is the resource implementation.
type idleAccountDataSecurityAuditorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultIdleAccountDataSecurityAuditorResource is the resource implementation.
type defaultIdleAccountDataSecurityAuditorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *idleAccountDataSecurityAuditorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_idle_account_data_security_auditor"
}

func (r *defaultIdleAccountDataSecurityAuditorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_idle_account_data_security_auditor"
}

// Configure adds the provider configured client to the resource.
func (r *idleAccountDataSecurityAuditorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultIdleAccountDataSecurityAuditorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type idleAccountDataSecurityAuditorResourceModel struct {
	Id                                  types.String `tfsdk:"id"`
	LastUpdated                         types.String `tfsdk:"last_updated"`
	Notifications                       types.Set    `tfsdk:"notifications"`
	RequiredActions                     types.Set    `tfsdk:"required_actions"`
	ReportFile                          types.String `tfsdk:"report_file"`
	IdleAccountWarningInterval          types.String `tfsdk:"idle_account_warning_interval"`
	IdleAccountErrorInterval            types.String `tfsdk:"idle_account_error_interval"`
	NeverLoggedInAccountWarningInterval types.String `tfsdk:"never_logged_in_account_warning_interval"`
	NeverLoggedInAccountErrorInterval   types.String `tfsdk:"never_logged_in_account_error_interval"`
	Enabled                             types.Bool   `tfsdk:"enabled"`
	IncludeAttribute                    types.Set    `tfsdk:"include_attribute"`
	AuditBackend                        types.Set    `tfsdk:"audit_backend"`
	AuditSeverity                       types.String `tfsdk:"audit_severity"`
}

// GetSchema defines the schema for the resource.
func (r *idleAccountDataSecurityAuditorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	idleAccountDataSecurityAuditorSchema(ctx, req, resp, false)
}

func (r *defaultIdleAccountDataSecurityAuditorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	idleAccountDataSecurityAuditorSchema(ctx, req, resp, true)
}

func idleAccountDataSecurityAuditorSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Idle Account Data Security Auditor. Supported in PingDirectory product version 9.2.0.0+.",
		Attributes: map[string]schema.Attribute{
			"report_file": schema.StringAttribute{
				Description: "Specifies the name of the detailed report file.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"idle_account_warning_interval": schema.StringAttribute{
				Description: "The length of time to use as the warning interval for idle accounts. If the length of time since a user last authenticated is greater than the warning interval but less than the error interval (or if it is greater than the warning interval and no error interval is defined), then a warning will be generated for that account.",
				Required:    true,
			},
			"idle_account_error_interval": schema.StringAttribute{
				Description: "The length of time to use as the error interval for idle accounts. If the length of time since a user last authenticated is greater than the error interval, then an error will be generated for that account. If no error interval is defined, then only the warning interval will be used.",
				Optional:    true,
			},
			"never_logged_in_account_warning_interval": schema.StringAttribute{
				Description: "The length of time to use as the warning interval for accounts that do not appear to have authenticated. If this is not specified, then the idle account warning interval will be used.",
				Optional:    true,
			},
			"never_logged_in_account_error_interval": schema.StringAttribute{
				Description: "The length of time to use as the error interval for accounts that do not appear to have authenticated. If this is not specified, then the never-logged-in warning interval will be used. The idle account warning and error intervals will be used if no never-logged-in interval is configured.",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Data Security Auditor is enabled for use.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_attribute": schema.SetAttribute{
				Description: "Specifies the attributes from the audited entries that should be included detailed reports. By default, no attributes are included.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"audit_backend": schema.SetAttribute{
				Description: "Specifies which backends the data security auditor may be applied to. By default, the data security auditors will audit entries in all backend types that support data auditing (Local DB, LDIF, and Config File Handler).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"audit_severity": schema.StringAttribute{
				Description: "Specifies the severity of events to include in the report.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Validate that any version restrictions are met in the plan
func (r *idleAccountDataSecurityAuditorResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanIdleAccountDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_idle_account_data_security_auditor")
}

func (r *defaultIdleAccountDataSecurityAuditorResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanIdleAccountDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_default_idle_account_data_security_auditor")
}

func modifyPlanIdleAccountDataSecurityAuditor(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, resourceName string) {
	version.CheckResourceSupported(&resp.Diagnostics, version.PingDirectory9200,
		providerConfig.ProductVersion, resourceName)
}

// Add config validators
func (r idleAccountDataSecurityAuditorResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.Implies(
			path.MatchRoot("never_logged_in_account_error_interval"),
			path.MatchRoot("never_logged_in_account_warning_interval"),
		),
	}
}

// Add optional fields to create request
func addOptionalIdleAccountDataSecurityAuditorFields(ctx context.Context, addRequest *client.AddIdleAccountDataSecurityAuditorRequest, plan idleAccountDataSecurityAuditorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReportFile) {
		addRequest.ReportFile = plan.ReportFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IdleAccountErrorInterval) {
		addRequest.IdleAccountErrorInterval = plan.IdleAccountErrorInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.NeverLoggedInAccountWarningInterval) {
		addRequest.NeverLoggedInAccountWarningInterval = plan.NeverLoggedInAccountWarningInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.NeverLoggedInAccountErrorInterval) {
		addRequest.NeverLoggedInAccountErrorInterval = plan.NeverLoggedInAccountErrorInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAttribute) {
		var slice []string
		plan.IncludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.IncludeAttribute = slice
	}
	if internaltypes.IsDefined(plan.AuditBackend) {
		var slice []string
		plan.AuditBackend.ElementsAs(ctx, &slice, false)
		addRequest.AuditBackend = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuditSeverity) {
		auditSeverity, err := client.NewEnumdataSecurityAuditorAuditSeverityPropFromValue(plan.AuditSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuditSeverity = auditSeverity
	}
	return nil
}

// Read a IdleAccountDataSecurityAuditorResponse object into the model struct
func readIdleAccountDataSecurityAuditorResponse(ctx context.Context, r *client.IdleAccountDataSecurityAuditorResponse, state *idleAccountDataSecurityAuditorResourceModel, expectedValues *idleAccountDataSecurityAuditorResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.IdleAccountWarningInterval = types.StringValue(r.IdleAccountWarningInterval)
	config.CheckMismatchedPDFormattedAttributes("idle_account_warning_interval",
		expectedValues.IdleAccountWarningInterval, state.IdleAccountWarningInterval, diagnostics)
	state.IdleAccountErrorInterval = internaltypes.StringTypeOrNil(r.IdleAccountErrorInterval, internaltypes.IsEmptyString(expectedValues.IdleAccountErrorInterval))
	config.CheckMismatchedPDFormattedAttributes("idle_account_error_interval",
		expectedValues.IdleAccountErrorInterval, state.IdleAccountErrorInterval, diagnostics)
	state.NeverLoggedInAccountWarningInterval = internaltypes.StringTypeOrNil(r.NeverLoggedInAccountWarningInterval, internaltypes.IsEmptyString(expectedValues.NeverLoggedInAccountWarningInterval))
	config.CheckMismatchedPDFormattedAttributes("never_logged_in_account_warning_interval",
		expectedValues.NeverLoggedInAccountWarningInterval, state.NeverLoggedInAccountWarningInterval, diagnostics)
	state.NeverLoggedInAccountErrorInterval = internaltypes.StringTypeOrNil(r.NeverLoggedInAccountErrorInterval, internaltypes.IsEmptyString(expectedValues.NeverLoggedInAccountErrorInterval))
	config.CheckMismatchedPDFormattedAttributes("never_logged_in_account_error_interval",
		expectedValues.NeverLoggedInAccountErrorInterval, state.NeverLoggedInAccountErrorInterval, diagnostics)
	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), internaltypes.IsEmptyString(expectedValues.AuditSeverity))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createIdleAccountDataSecurityAuditorOperations(plan idleAccountDataSecurityAuditorResourceModel, state idleAccountDataSecurityAuditorResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ReportFile, state.ReportFile, "report-file")
	operations.AddStringOperationIfNecessary(&ops, plan.IdleAccountWarningInterval, state.IdleAccountWarningInterval, "idle-account-warning-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.IdleAccountErrorInterval, state.IdleAccountErrorInterval, "idle-account-error-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.NeverLoggedInAccountWarningInterval, state.NeverLoggedInAccountWarningInterval, "never-logged-in-account-warning-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.NeverLoggedInAccountErrorInterval, state.NeverLoggedInAccountErrorInterval, "never-logged-in-account-error-interval")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeAttribute, state.IncludeAttribute, "include-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AuditBackend, state.AuditBackend, "audit-backend")
	operations.AddStringOperationIfNecessary(&ops, plan.AuditSeverity, state.AuditSeverity, "audit-severity")
	return ops
}

// Create a new resource
func (r *idleAccountDataSecurityAuditorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan idleAccountDataSecurityAuditorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddIdleAccountDataSecurityAuditorRequest(plan.Id.ValueString(),
		[]client.EnumidleAccountDataSecurityAuditorSchemaUrn{client.ENUMIDLEACCOUNTDATASECURITYAUDITORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DATA_SECURITY_AUDITORIDLE_ACCOUNT},
		plan.IdleAccountWarningInterval.ValueString())
	err := addOptionalIdleAccountDataSecurityAuditorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Idle Account Data Security Auditor", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DataSecurityAuditorApi.AddDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddDataSecurityAuditorRequest(
		client.AddIdleAccountDataSecurityAuditorRequestAsAddDataSecurityAuditorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DataSecurityAuditorApi.AddDataSecurityAuditorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Idle Account Data Security Auditor", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state idleAccountDataSecurityAuditorResourceModel
	readIdleAccountDataSecurityAuditorResponse(ctx, addResponse.IdleAccountDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultIdleAccountDataSecurityAuditorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan idleAccountDataSecurityAuditorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.DataSecurityAuditorApi.GetDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Idle Account Data Security Auditor", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state idleAccountDataSecurityAuditorResourceModel
	readIdleAccountDataSecurityAuditorResponse(ctx, readResponse.IdleAccountDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.DataSecurityAuditorApi.UpdateDataSecurityAuditor(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createIdleAccountDataSecurityAuditorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.DataSecurityAuditorApi.UpdateDataSecurityAuditorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Idle Account Data Security Auditor", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readIdleAccountDataSecurityAuditorResponse(ctx, updateResponse.IdleAccountDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *idleAccountDataSecurityAuditorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readIdleAccountDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultIdleAccountDataSecurityAuditorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readIdleAccountDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readIdleAccountDataSecurityAuditor(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state idleAccountDataSecurityAuditorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.DataSecurityAuditorApi.GetDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Idle Account Data Security Auditor", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readIdleAccountDataSecurityAuditorResponse(ctx, readResponse.IdleAccountDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *idleAccountDataSecurityAuditorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateIdleAccountDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultIdleAccountDataSecurityAuditorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateIdleAccountDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateIdleAccountDataSecurityAuditor(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan idleAccountDataSecurityAuditorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state idleAccountDataSecurityAuditorResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.DataSecurityAuditorApi.UpdateDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createIdleAccountDataSecurityAuditorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.DataSecurityAuditorApi.UpdateDataSecurityAuditorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Idle Account Data Security Auditor", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readIdleAccountDataSecurityAuditorResponse(ctx, updateResponse.IdleAccountDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultIdleAccountDataSecurityAuditorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *idleAccountDataSecurityAuditorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state idleAccountDataSecurityAuditorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.DataSecurityAuditorApi.DeleteDataSecurityAuditorExecute(r.apiClient.DataSecurityAuditorApi.DeleteDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Idle Account Data Security Auditor", err, httpResp)
		return
	}
}

func (r *idleAccountDataSecurityAuditorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importIdleAccountDataSecurityAuditor(ctx, req, resp)
}

func (r *defaultIdleAccountDataSecurityAuditorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importIdleAccountDataSecurityAuditor(ctx, req, resp)
}

func importIdleAccountDataSecurityAuditor(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
