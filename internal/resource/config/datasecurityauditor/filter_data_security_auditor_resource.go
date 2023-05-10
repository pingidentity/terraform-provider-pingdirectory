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
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &filterDataSecurityAuditorResource{}
	_ resource.ResourceWithConfigure   = &filterDataSecurityAuditorResource{}
	_ resource.ResourceWithImportState = &filterDataSecurityAuditorResource{}
	_ resource.Resource                = &defaultFilterDataSecurityAuditorResource{}
	_ resource.ResourceWithConfigure   = &defaultFilterDataSecurityAuditorResource{}
	_ resource.ResourceWithImportState = &defaultFilterDataSecurityAuditorResource{}
)

// Create a Filter Data Security Auditor resource
func NewFilterDataSecurityAuditorResource() resource.Resource {
	return &filterDataSecurityAuditorResource{}
}

func NewDefaultFilterDataSecurityAuditorResource() resource.Resource {
	return &defaultFilterDataSecurityAuditorResource{}
}

// filterDataSecurityAuditorResource is the resource implementation.
type filterDataSecurityAuditorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultFilterDataSecurityAuditorResource is the resource implementation.
type defaultFilterDataSecurityAuditorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *filterDataSecurityAuditorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filter_data_security_auditor"
}

func (r *defaultFilterDataSecurityAuditorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_filter_data_security_auditor"
}

// Configure adds the provider configured client to the resource.
func (r *filterDataSecurityAuditorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultFilterDataSecurityAuditorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type filterDataSecurityAuditorResourceModel struct {
	Id               types.String `tfsdk:"id"`
	LastUpdated      types.String `tfsdk:"last_updated"`
	Notifications    types.Set    `tfsdk:"notifications"`
	RequiredActions  types.Set    `tfsdk:"required_actions"`
	ReportFile       types.String `tfsdk:"report_file"`
	Filter           types.Set    `tfsdk:"filter"`
	Enabled          types.Bool   `tfsdk:"enabled"`
	IncludeAttribute types.Set    `tfsdk:"include_attribute"`
	AuditBackend     types.Set    `tfsdk:"audit_backend"`
	AuditSeverity    types.String `tfsdk:"audit_severity"`
}

// GetSchema defines the schema for the resource.
func (r *filterDataSecurityAuditorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	filterDataSecurityAuditorSchema(ctx, req, resp, false)
}

func (r *defaultFilterDataSecurityAuditorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	filterDataSecurityAuditorSchema(ctx, req, resp, true)
}

func filterDataSecurityAuditorSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Filter Data Security Auditor. Supported in PingDirectory product version 9.2.0.0+.",
		Attributes: map[string]schema.Attribute{
			"report_file": schema.StringAttribute{
				Description: "Specifies the name of the detailed report file.",
				Required:    true,
			},
			"filter": schema.SetAttribute{
				Description: "The filter to use to identify entries that should be reported. Multiple filters may be configured, and each reported entry will indicate which of these filter(s) matched that entry.",
				Required:    true,
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

// Validate that any version restrictions are met in the plan
func (r *filterDataSecurityAuditorResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanFilterDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_filter_data_security_auditor")
}

func (r *defaultFilterDataSecurityAuditorResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanFilterDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_default_filter_data_security_auditor")
}

func modifyPlanFilterDataSecurityAuditor(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, resourceName string) {
	version.CheckResourceSupported(&resp.Diagnostics, version.PingDirectory9200,
		providerConfig.ProductVersion, resourceName)
}

// Add optional fields to create request
func addOptionalFilterDataSecurityAuditorFields(ctx context.Context, addRequest *client.AddFilterDataSecurityAuditorRequest, plan filterDataSecurityAuditorResourceModel) error {
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

// Read a FilterDataSecurityAuditorResponse object into the model struct
func readFilterDataSecurityAuditorResponse(ctx context.Context, r *client.FilterDataSecurityAuditorResponse, state *filterDataSecurityAuditorResourceModel, expectedValues *filterDataSecurityAuditorResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), internaltypes.IsEmptyString(expectedValues.AuditSeverity))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createFilterDataSecurityAuditorOperations(plan filterDataSecurityAuditorResourceModel, state filterDataSecurityAuditorResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ReportFile, state.ReportFile, "report-file")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.Filter, state.Filter, "filter")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeAttribute, state.IncludeAttribute, "include-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AuditBackend, state.AuditBackend, "audit-backend")
	operations.AddStringOperationIfNecessary(&ops, plan.AuditSeverity, state.AuditSeverity, "audit-severity")
	return ops
}

// Create a new resource
func (r *filterDataSecurityAuditorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan filterDataSecurityAuditorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var FilterSlice []string
	plan.Filter.ElementsAs(ctx, &FilterSlice, false)
	addRequest := client.NewAddFilterDataSecurityAuditorRequest(plan.Id.ValueString(),
		[]client.EnumfilterDataSecurityAuditorSchemaUrn{client.ENUMFILTERDATASECURITYAUDITORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DATA_SECURITY_AUDITORFILTER},
		plan.ReportFile.ValueString(),
		FilterSlice)
	err := addOptionalFilterDataSecurityAuditorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Filter Data Security Auditor", err.Error())
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
		client.AddFilterDataSecurityAuditorRequestAsAddDataSecurityAuditorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DataSecurityAuditorApi.AddDataSecurityAuditorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Filter Data Security Auditor", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state filterDataSecurityAuditorResourceModel
	readFilterDataSecurityAuditorResponse(ctx, addResponse.FilterDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultFilterDataSecurityAuditorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan filterDataSecurityAuditorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.DataSecurityAuditorApi.GetDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Filter Data Security Auditor", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state filterDataSecurityAuditorResourceModel
	readFilterDataSecurityAuditorResponse(ctx, readResponse.FilterDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.DataSecurityAuditorApi.UpdateDataSecurityAuditor(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createFilterDataSecurityAuditorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.DataSecurityAuditorApi.UpdateDataSecurityAuditorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Filter Data Security Auditor", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readFilterDataSecurityAuditorResponse(ctx, updateResponse.FilterDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *filterDataSecurityAuditorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readFilterDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultFilterDataSecurityAuditorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readFilterDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readFilterDataSecurityAuditor(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state filterDataSecurityAuditorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.DataSecurityAuditorApi.GetDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Filter Data Security Auditor", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readFilterDataSecurityAuditorResponse(ctx, readResponse.FilterDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *filterDataSecurityAuditorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateFilterDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultFilterDataSecurityAuditorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateFilterDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateFilterDataSecurityAuditor(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan filterDataSecurityAuditorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state filterDataSecurityAuditorResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.DataSecurityAuditorApi.UpdateDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createFilterDataSecurityAuditorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.DataSecurityAuditorApi.UpdateDataSecurityAuditorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Filter Data Security Auditor", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readFilterDataSecurityAuditorResponse(ctx, updateResponse.FilterDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultFilterDataSecurityAuditorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *filterDataSecurityAuditorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state filterDataSecurityAuditorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.DataSecurityAuditorApi.DeleteDataSecurityAuditorExecute(r.apiClient.DataSecurityAuditorApi.DeleteDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Filter Data Security Auditor", err, httpResp)
		return
	}
}

func (r *filterDataSecurityAuditorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importFilterDataSecurityAuditor(ctx, req, resp)
}

func (r *defaultFilterDataSecurityAuditorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importFilterDataSecurityAuditor(ctx, req, resp)
}

func importFilterDataSecurityAuditor(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
