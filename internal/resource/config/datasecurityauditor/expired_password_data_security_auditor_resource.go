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
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &expiredPasswordDataSecurityAuditorResource{}
	_ resource.ResourceWithConfigure   = &expiredPasswordDataSecurityAuditorResource{}
	_ resource.ResourceWithImportState = &expiredPasswordDataSecurityAuditorResource{}
	_ resource.Resource                = &defaultExpiredPasswordDataSecurityAuditorResource{}
	_ resource.ResourceWithConfigure   = &defaultExpiredPasswordDataSecurityAuditorResource{}
	_ resource.ResourceWithImportState = &defaultExpiredPasswordDataSecurityAuditorResource{}
)

// Create a Expired Password Data Security Auditor resource
func NewExpiredPasswordDataSecurityAuditorResource() resource.Resource {
	return &expiredPasswordDataSecurityAuditorResource{}
}

func NewDefaultExpiredPasswordDataSecurityAuditorResource() resource.Resource {
	return &defaultExpiredPasswordDataSecurityAuditorResource{}
}

// expiredPasswordDataSecurityAuditorResource is the resource implementation.
type expiredPasswordDataSecurityAuditorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultExpiredPasswordDataSecurityAuditorResource is the resource implementation.
type defaultExpiredPasswordDataSecurityAuditorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *expiredPasswordDataSecurityAuditorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_expired_password_data_security_auditor"
}

func (r *defaultExpiredPasswordDataSecurityAuditorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_expired_password_data_security_auditor"
}

// Configure adds the provider configured client to the resource.
func (r *expiredPasswordDataSecurityAuditorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultExpiredPasswordDataSecurityAuditorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type expiredPasswordDataSecurityAuditorResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	LastUpdated           types.String `tfsdk:"last_updated"`
	Notifications         types.Set    `tfsdk:"notifications"`
	RequiredActions       types.Set    `tfsdk:"required_actions"`
	ReportFile            types.String `tfsdk:"report_file"`
	IncludeAttribute      types.Set    `tfsdk:"include_attribute"`
	PasswordEvaluationAge types.String `tfsdk:"password_evaluation_age"`
	Enabled               types.Bool   `tfsdk:"enabled"`
	AuditBackend          types.Set    `tfsdk:"audit_backend"`
	AuditSeverity         types.String `tfsdk:"audit_severity"`
}

// GetSchema defines the schema for the resource.
func (r *expiredPasswordDataSecurityAuditorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	expiredPasswordDataSecurityAuditorSchema(ctx, req, resp, false)
}

func (r *defaultExpiredPasswordDataSecurityAuditorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	expiredPasswordDataSecurityAuditorSchema(ctx, req, resp, true)
}

func expiredPasswordDataSecurityAuditorSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Expired Password Data Security Auditor.",
		Attributes: map[string]schema.Attribute{
			"report_file": schema.StringAttribute{
				Description: "Specifies the name of the detailed report file.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
			"password_evaluation_age": schema.StringAttribute{
				Description: "If set, the auditor will report all users with passwords older than the specified value even if password expiration is not enabled.",
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

// Add optional fields to create request
func addOptionalExpiredPasswordDataSecurityAuditorFields(ctx context.Context, addRequest *client.AddExpiredPasswordDataSecurityAuditorRequest, plan expiredPasswordDataSecurityAuditorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReportFile) {
		addRequest.ReportFile = plan.ReportFile.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAttribute) {
		var slice []string
		plan.IncludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.IncludeAttribute = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PasswordEvaluationAge) {
		addRequest.PasswordEvaluationAge = plan.PasswordEvaluationAge.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
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

// Read a ExpiredPasswordDataSecurityAuditorResponse object into the model struct
func readExpiredPasswordDataSecurityAuditorResponse(ctx context.Context, r *client.ExpiredPasswordDataSecurityAuditorResponse, state *expiredPasswordDataSecurityAuditorResourceModel, expectedValues *expiredPasswordDataSecurityAuditorResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.PasswordEvaluationAge = internaltypes.StringTypeOrNil(r.PasswordEvaluationAge, internaltypes.IsEmptyString(expectedValues.PasswordEvaluationAge))
	config.CheckMismatchedPDFormattedAttributes("password_evaluation_age",
		expectedValues.PasswordEvaluationAge, state.PasswordEvaluationAge, diagnostics)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), internaltypes.IsEmptyString(expectedValues.AuditSeverity))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createExpiredPasswordDataSecurityAuditorOperations(plan expiredPasswordDataSecurityAuditorResourceModel, state expiredPasswordDataSecurityAuditorResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ReportFile, state.ReportFile, "report-file")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeAttribute, state.IncludeAttribute, "include-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordEvaluationAge, state.PasswordEvaluationAge, "password-evaluation-age")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AuditBackend, state.AuditBackend, "audit-backend")
	operations.AddStringOperationIfNecessary(&ops, plan.AuditSeverity, state.AuditSeverity, "audit-severity")
	return ops
}

// Create a new resource
func (r *expiredPasswordDataSecurityAuditorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan expiredPasswordDataSecurityAuditorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddExpiredPasswordDataSecurityAuditorRequest(plan.Id.ValueString(),
		[]client.EnumexpiredPasswordDataSecurityAuditorSchemaUrn{client.ENUMEXPIREDPASSWORDDATASECURITYAUDITORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DATA_SECURITY_AUDITOREXPIRED_PASSWORD})
	err := addOptionalExpiredPasswordDataSecurityAuditorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Expired Password Data Security Auditor", err.Error())
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
		client.AddExpiredPasswordDataSecurityAuditorRequestAsAddDataSecurityAuditorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DataSecurityAuditorApi.AddDataSecurityAuditorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Expired Password Data Security Auditor", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state expiredPasswordDataSecurityAuditorResourceModel
	readExpiredPasswordDataSecurityAuditorResponse(ctx, addResponse.ExpiredPasswordDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultExpiredPasswordDataSecurityAuditorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan expiredPasswordDataSecurityAuditorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.DataSecurityAuditorApi.GetDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Expired Password Data Security Auditor", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state expiredPasswordDataSecurityAuditorResourceModel
	readExpiredPasswordDataSecurityAuditorResponse(ctx, readResponse.ExpiredPasswordDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.DataSecurityAuditorApi.UpdateDataSecurityAuditor(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createExpiredPasswordDataSecurityAuditorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.DataSecurityAuditorApi.UpdateDataSecurityAuditorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Expired Password Data Security Auditor", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readExpiredPasswordDataSecurityAuditorResponse(ctx, updateResponse.ExpiredPasswordDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *expiredPasswordDataSecurityAuditorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readExpiredPasswordDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultExpiredPasswordDataSecurityAuditorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readExpiredPasswordDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readExpiredPasswordDataSecurityAuditor(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state expiredPasswordDataSecurityAuditorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.DataSecurityAuditorApi.GetDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Expired Password Data Security Auditor", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readExpiredPasswordDataSecurityAuditorResponse(ctx, readResponse.ExpiredPasswordDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *expiredPasswordDataSecurityAuditorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateExpiredPasswordDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultExpiredPasswordDataSecurityAuditorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateExpiredPasswordDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateExpiredPasswordDataSecurityAuditor(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan expiredPasswordDataSecurityAuditorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state expiredPasswordDataSecurityAuditorResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.DataSecurityAuditorApi.UpdateDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createExpiredPasswordDataSecurityAuditorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.DataSecurityAuditorApi.UpdateDataSecurityAuditorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Expired Password Data Security Auditor", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readExpiredPasswordDataSecurityAuditorResponse(ctx, updateResponse.ExpiredPasswordDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultExpiredPasswordDataSecurityAuditorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *expiredPasswordDataSecurityAuditorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state expiredPasswordDataSecurityAuditorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.DataSecurityAuditorApi.DeleteDataSecurityAuditorExecute(r.apiClient.DataSecurityAuditorApi.DeleteDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Expired Password Data Security Auditor", err, httpResp)
		return
	}
}

func (r *expiredPasswordDataSecurityAuditorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importExpiredPasswordDataSecurityAuditor(ctx, req, resp)
}

func (r *defaultExpiredPasswordDataSecurityAuditorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importExpiredPasswordDataSecurityAuditor(ctx, req, resp)
}

func importExpiredPasswordDataSecurityAuditor(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
