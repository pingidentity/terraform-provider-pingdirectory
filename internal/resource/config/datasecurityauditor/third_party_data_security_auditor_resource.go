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
	_ resource.Resource                = &thirdPartyDataSecurityAuditorResource{}
	_ resource.ResourceWithConfigure   = &thirdPartyDataSecurityAuditorResource{}
	_ resource.ResourceWithImportState = &thirdPartyDataSecurityAuditorResource{}
	_ resource.Resource                = &defaultThirdPartyDataSecurityAuditorResource{}
	_ resource.ResourceWithConfigure   = &defaultThirdPartyDataSecurityAuditorResource{}
	_ resource.ResourceWithImportState = &defaultThirdPartyDataSecurityAuditorResource{}
)

// Create a Third Party Data Security Auditor resource
func NewThirdPartyDataSecurityAuditorResource() resource.Resource {
	return &thirdPartyDataSecurityAuditorResource{}
}

func NewDefaultThirdPartyDataSecurityAuditorResource() resource.Resource {
	return &defaultThirdPartyDataSecurityAuditorResource{}
}

// thirdPartyDataSecurityAuditorResource is the resource implementation.
type thirdPartyDataSecurityAuditorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultThirdPartyDataSecurityAuditorResource is the resource implementation.
type defaultThirdPartyDataSecurityAuditorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *thirdPartyDataSecurityAuditorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_third_party_data_security_auditor"
}

func (r *defaultThirdPartyDataSecurityAuditorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_third_party_data_security_auditor"
}

// Configure adds the provider configured client to the resource.
func (r *thirdPartyDataSecurityAuditorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultThirdPartyDataSecurityAuditorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type thirdPartyDataSecurityAuditorResourceModel struct {
	Id                types.String `tfsdk:"id"`
	LastUpdated       types.String `tfsdk:"last_updated"`
	Notifications     types.Set    `tfsdk:"notifications"`
	RequiredActions   types.Set    `tfsdk:"required_actions"`
	ExtensionClass    types.String `tfsdk:"extension_class"`
	ExtensionArgument types.Set    `tfsdk:"extension_argument"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	ReportFile        types.String `tfsdk:"report_file"`
	IncludeAttribute  types.Set    `tfsdk:"include_attribute"`
	AuditBackend      types.Set    `tfsdk:"audit_backend"`
	AuditSeverity     types.String `tfsdk:"audit_severity"`
}

// GetSchema defines the schema for the resource.
func (r *thirdPartyDataSecurityAuditorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	thirdPartyDataSecurityAuditorSchema(ctx, req, resp, false)
}

func (r *defaultThirdPartyDataSecurityAuditorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	thirdPartyDataSecurityAuditorSchema(ctx, req, resp, true)
}

func thirdPartyDataSecurityAuditorSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Third Party Data Security Auditor. Supported in PingDirectory product version 9.2.0.0+.",
		Attributes: map[string]schema.Attribute{
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Data Security Auditor.",
				Required:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Data Security Auditor. Each configuration property should be given in the form 'name=value'.",
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
			"report_file": schema.StringAttribute{
				Description: "Specifies the name of the detailed report file.",
				Required:    true,
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
func (r *thirdPartyDataSecurityAuditorResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanThirdPartyDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_third_party_data_security_auditor")
}

func (r *defaultThirdPartyDataSecurityAuditorResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanThirdPartyDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_default_third_party_data_security_auditor")
}

func modifyPlanThirdPartyDataSecurityAuditor(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, resourceName string) {
	version.CheckResourceSupported(&resp.Diagnostics, version.PingDirectory9200,
		providerConfig.ProductVersion, resourceName)
}

// Add optional fields to create request
func addOptionalThirdPartyDataSecurityAuditorFields(ctx context.Context, addRequest *client.AddThirdPartyDataSecurityAuditorRequest, plan thirdPartyDataSecurityAuditorResourceModel) error {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
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

// Read a ThirdPartyDataSecurityAuditorResponse object into the model struct
func readThirdPartyDataSecurityAuditorResponse(ctx context.Context, r *client.ThirdPartyDataSecurityAuditorResponse, state *thirdPartyDataSecurityAuditorResourceModel, expectedValues *thirdPartyDataSecurityAuditorResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), internaltypes.IsEmptyString(expectedValues.AuditSeverity))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createThirdPartyDataSecurityAuditorOperations(plan thirdPartyDataSecurityAuditorResourceModel, state thirdPartyDataSecurityAuditorResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.ReportFile, state.ReportFile, "report-file")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeAttribute, state.IncludeAttribute, "include-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AuditBackend, state.AuditBackend, "audit-backend")
	operations.AddStringOperationIfNecessary(&ops, plan.AuditSeverity, state.AuditSeverity, "audit-severity")
	return ops
}

// Create a new resource
func (r *thirdPartyDataSecurityAuditorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan thirdPartyDataSecurityAuditorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddThirdPartyDataSecurityAuditorRequest(plan.Id.ValueString(),
		[]client.EnumthirdPartyDataSecurityAuditorSchemaUrn{client.ENUMTHIRDPARTYDATASECURITYAUDITORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DATA_SECURITY_AUDITORTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.ReportFile.ValueString())
	err := addOptionalThirdPartyDataSecurityAuditorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Third Party Data Security Auditor", err.Error())
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
		client.AddThirdPartyDataSecurityAuditorRequestAsAddDataSecurityAuditorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DataSecurityAuditorApi.AddDataSecurityAuditorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Third Party Data Security Auditor", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state thirdPartyDataSecurityAuditorResourceModel
	readThirdPartyDataSecurityAuditorResponse(ctx, addResponse.ThirdPartyDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultThirdPartyDataSecurityAuditorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan thirdPartyDataSecurityAuditorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.DataSecurityAuditorApi.GetDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Third Party Data Security Auditor", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state thirdPartyDataSecurityAuditorResourceModel
	readThirdPartyDataSecurityAuditorResponse(ctx, readResponse.ThirdPartyDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.DataSecurityAuditorApi.UpdateDataSecurityAuditor(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createThirdPartyDataSecurityAuditorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.DataSecurityAuditorApi.UpdateDataSecurityAuditorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Third Party Data Security Auditor", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readThirdPartyDataSecurityAuditorResponse(ctx, updateResponse.ThirdPartyDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *thirdPartyDataSecurityAuditorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readThirdPartyDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultThirdPartyDataSecurityAuditorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readThirdPartyDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readThirdPartyDataSecurityAuditor(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state thirdPartyDataSecurityAuditorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.DataSecurityAuditorApi.GetDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Third Party Data Security Auditor", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readThirdPartyDataSecurityAuditorResponse(ctx, readResponse.ThirdPartyDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *thirdPartyDataSecurityAuditorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateThirdPartyDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultThirdPartyDataSecurityAuditorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateThirdPartyDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateThirdPartyDataSecurityAuditor(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan thirdPartyDataSecurityAuditorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state thirdPartyDataSecurityAuditorResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.DataSecurityAuditorApi.UpdateDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createThirdPartyDataSecurityAuditorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.DataSecurityAuditorApi.UpdateDataSecurityAuditorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Third Party Data Security Auditor", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readThirdPartyDataSecurityAuditorResponse(ctx, updateResponse.ThirdPartyDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultThirdPartyDataSecurityAuditorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *thirdPartyDataSecurityAuditorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state thirdPartyDataSecurityAuditorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.DataSecurityAuditorApi.DeleteDataSecurityAuditorExecute(r.apiClient.DataSecurityAuditorApi.DeleteDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Third Party Data Security Auditor", err, httpResp)
		return
	}
}

func (r *thirdPartyDataSecurityAuditorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importThirdPartyDataSecurityAuditor(ctx, req, resp)
}

func (r *defaultThirdPartyDataSecurityAuditorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importThirdPartyDataSecurityAuditor(ctx, req, resp)
}

func importThirdPartyDataSecurityAuditor(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
