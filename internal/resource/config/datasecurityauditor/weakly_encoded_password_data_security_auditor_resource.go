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
	_ resource.Resource                = &weaklyEncodedPasswordDataSecurityAuditorResource{}
	_ resource.ResourceWithConfigure   = &weaklyEncodedPasswordDataSecurityAuditorResource{}
	_ resource.ResourceWithImportState = &weaklyEncodedPasswordDataSecurityAuditorResource{}
	_ resource.Resource                = &defaultWeaklyEncodedPasswordDataSecurityAuditorResource{}
	_ resource.ResourceWithConfigure   = &defaultWeaklyEncodedPasswordDataSecurityAuditorResource{}
	_ resource.ResourceWithImportState = &defaultWeaklyEncodedPasswordDataSecurityAuditorResource{}
)

// Create a Weakly Encoded Password Data Security Auditor resource
func NewWeaklyEncodedPasswordDataSecurityAuditorResource() resource.Resource {
	return &weaklyEncodedPasswordDataSecurityAuditorResource{}
}

func NewDefaultWeaklyEncodedPasswordDataSecurityAuditorResource() resource.Resource {
	return &defaultWeaklyEncodedPasswordDataSecurityAuditorResource{}
}

// weaklyEncodedPasswordDataSecurityAuditorResource is the resource implementation.
type weaklyEncodedPasswordDataSecurityAuditorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultWeaklyEncodedPasswordDataSecurityAuditorResource is the resource implementation.
type defaultWeaklyEncodedPasswordDataSecurityAuditorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *weaklyEncodedPasswordDataSecurityAuditorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_weakly_encoded_password_data_security_auditor"
}

func (r *defaultWeaklyEncodedPasswordDataSecurityAuditorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_weakly_encoded_password_data_security_auditor"
}

// Configure adds the provider configured client to the resource.
func (r *weaklyEncodedPasswordDataSecurityAuditorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultWeaklyEncodedPasswordDataSecurityAuditorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type weaklyEncodedPasswordDataSecurityAuditorResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	LastUpdated               types.String `tfsdk:"last_updated"`
	Notifications             types.Set    `tfsdk:"notifications"`
	RequiredActions           types.Set    `tfsdk:"required_actions"`
	ReportFile                types.String `tfsdk:"report_file"`
	WeakPasswordStorageScheme types.Set    `tfsdk:"weak_password_storage_scheme"`
	WeakCryptEncoding         types.Set    `tfsdk:"weak_crypt_encoding"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
	IncludeAttribute          types.Set    `tfsdk:"include_attribute"`
	AuditBackend              types.Set    `tfsdk:"audit_backend"`
	AuditSeverity             types.String `tfsdk:"audit_severity"`
}

// GetSchema defines the schema for the resource.
func (r *weaklyEncodedPasswordDataSecurityAuditorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	weaklyEncodedPasswordDataSecurityAuditorSchema(ctx, req, resp, false)
}

func (r *defaultWeaklyEncodedPasswordDataSecurityAuditorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	weaklyEncodedPasswordDataSecurityAuditorSchema(ctx, req, resp, true)
}

func weaklyEncodedPasswordDataSecurityAuditorSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Weakly Encoded Password Data Security Auditor.",
		Attributes: map[string]schema.Attribute{
			"report_file": schema.StringAttribute{
				Description: "Specifies the name of the detailed report file.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"weak_password_storage_scheme": schema.SetAttribute{
				Description: "The password storage schemes that are considered weak. Users with any of the specified password storage schemes will be included in the report.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"weak_crypt_encoding": schema.SetAttribute{
				Description: "Reporting on users with passwords encoded using the Crypt Password Storage scheme may be further limited by selecting one or more encoding mechanisms that are considered weak.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
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

// Add optional fields to create request
func addOptionalWeaklyEncodedPasswordDataSecurityAuditorFields(ctx context.Context, addRequest *client.AddWeaklyEncodedPasswordDataSecurityAuditorRequest, plan weaklyEncodedPasswordDataSecurityAuditorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReportFile) {
		addRequest.ReportFile = plan.ReportFile.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.WeakPasswordStorageScheme) {
		var slice []string
		plan.WeakPasswordStorageScheme.ElementsAs(ctx, &slice, false)
		addRequest.WeakPasswordStorageScheme = slice
	}
	if internaltypes.IsDefined(plan.WeakCryptEncoding) {
		var slice []string
		plan.WeakCryptEncoding.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumdataSecurityAuditorWeakCryptEncodingProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumdataSecurityAuditorWeakCryptEncodingPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.WeakCryptEncoding = enumSlice
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

// Read a WeaklyEncodedPasswordDataSecurityAuditorResponse object into the model struct
func readWeaklyEncodedPasswordDataSecurityAuditorResponse(ctx context.Context, r *client.WeaklyEncodedPasswordDataSecurityAuditorResponse, state *weaklyEncodedPasswordDataSecurityAuditorResourceModel, expectedValues *weaklyEncodedPasswordDataSecurityAuditorResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.WeakPasswordStorageScheme = internaltypes.GetStringSet(r.WeakPasswordStorageScheme)
	state.WeakCryptEncoding = internaltypes.GetStringSet(
		client.StringSliceEnumdataSecurityAuditorWeakCryptEncodingProp(r.WeakCryptEncoding))
	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), internaltypes.IsEmptyString(expectedValues.AuditSeverity))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createWeaklyEncodedPasswordDataSecurityAuditorOperations(plan weaklyEncodedPasswordDataSecurityAuditorResourceModel, state weaklyEncodedPasswordDataSecurityAuditorResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ReportFile, state.ReportFile, "report-file")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.WeakPasswordStorageScheme, state.WeakPasswordStorageScheme, "weak-password-storage-scheme")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.WeakCryptEncoding, state.WeakCryptEncoding, "weak-crypt-encoding")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeAttribute, state.IncludeAttribute, "include-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AuditBackend, state.AuditBackend, "audit-backend")
	operations.AddStringOperationIfNecessary(&ops, plan.AuditSeverity, state.AuditSeverity, "audit-severity")
	return ops
}

// Create a new resource
func (r *weaklyEncodedPasswordDataSecurityAuditorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan weaklyEncodedPasswordDataSecurityAuditorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddWeaklyEncodedPasswordDataSecurityAuditorRequest(plan.Id.ValueString(),
		[]client.EnumweaklyEncodedPasswordDataSecurityAuditorSchemaUrn{client.ENUMWEAKLYENCODEDPASSWORDDATASECURITYAUDITORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DATA_SECURITY_AUDITORWEAKLY_ENCODED_PASSWORD})
	err := addOptionalWeaklyEncodedPasswordDataSecurityAuditorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Weakly Encoded Password Data Security Auditor", err.Error())
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
		client.AddWeaklyEncodedPasswordDataSecurityAuditorRequestAsAddDataSecurityAuditorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DataSecurityAuditorApi.AddDataSecurityAuditorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Weakly Encoded Password Data Security Auditor", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state weaklyEncodedPasswordDataSecurityAuditorResourceModel
	readWeaklyEncodedPasswordDataSecurityAuditorResponse(ctx, addResponse.WeaklyEncodedPasswordDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultWeaklyEncodedPasswordDataSecurityAuditorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan weaklyEncodedPasswordDataSecurityAuditorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.DataSecurityAuditorApi.GetDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Weakly Encoded Password Data Security Auditor", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state weaklyEncodedPasswordDataSecurityAuditorResourceModel
	readWeaklyEncodedPasswordDataSecurityAuditorResponse(ctx, readResponse.WeaklyEncodedPasswordDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.DataSecurityAuditorApi.UpdateDataSecurityAuditor(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createWeaklyEncodedPasswordDataSecurityAuditorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.DataSecurityAuditorApi.UpdateDataSecurityAuditorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Weakly Encoded Password Data Security Auditor", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readWeaklyEncodedPasswordDataSecurityAuditorResponse(ctx, updateResponse.WeaklyEncodedPasswordDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *weaklyEncodedPasswordDataSecurityAuditorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readWeaklyEncodedPasswordDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultWeaklyEncodedPasswordDataSecurityAuditorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readWeaklyEncodedPasswordDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readWeaklyEncodedPasswordDataSecurityAuditor(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state weaklyEncodedPasswordDataSecurityAuditorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.DataSecurityAuditorApi.GetDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Weakly Encoded Password Data Security Auditor", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readWeaklyEncodedPasswordDataSecurityAuditorResponse(ctx, readResponse.WeaklyEncodedPasswordDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *weaklyEncodedPasswordDataSecurityAuditorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateWeaklyEncodedPasswordDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultWeaklyEncodedPasswordDataSecurityAuditorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateWeaklyEncodedPasswordDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateWeaklyEncodedPasswordDataSecurityAuditor(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan weaklyEncodedPasswordDataSecurityAuditorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state weaklyEncodedPasswordDataSecurityAuditorResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.DataSecurityAuditorApi.UpdateDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createWeaklyEncodedPasswordDataSecurityAuditorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.DataSecurityAuditorApi.UpdateDataSecurityAuditorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Weakly Encoded Password Data Security Auditor", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readWeaklyEncodedPasswordDataSecurityAuditorResponse(ctx, updateResponse.WeaklyEncodedPasswordDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultWeaklyEncodedPasswordDataSecurityAuditorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *weaklyEncodedPasswordDataSecurityAuditorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state weaklyEncodedPasswordDataSecurityAuditorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.DataSecurityAuditorApi.DeleteDataSecurityAuditorExecute(r.apiClient.DataSecurityAuditorApi.DeleteDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Weakly Encoded Password Data Security Auditor", err, httpResp)
		return
	}
}

func (r *weaklyEncodedPasswordDataSecurityAuditorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importWeaklyEncodedPasswordDataSecurityAuditor(ctx, req, resp)
}

func (r *defaultWeaklyEncodedPasswordDataSecurityAuditorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importWeaklyEncodedPasswordDataSecurityAuditor(ctx, req, resp)
}

func importWeaklyEncodedPasswordDataSecurityAuditor(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
