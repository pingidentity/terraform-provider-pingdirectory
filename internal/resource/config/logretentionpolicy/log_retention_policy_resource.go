package logretentionpolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10000/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &logRetentionPolicyResource{}
	_ resource.ResourceWithConfigure   = &logRetentionPolicyResource{}
	_ resource.ResourceWithImportState = &logRetentionPolicyResource{}
	_ resource.Resource                = &defaultLogRetentionPolicyResource{}
	_ resource.ResourceWithConfigure   = &defaultLogRetentionPolicyResource{}
	_ resource.ResourceWithImportState = &defaultLogRetentionPolicyResource{}
)

// Create a Log Retention Policy resource
func NewLogRetentionPolicyResource() resource.Resource {
	return &logRetentionPolicyResource{}
}

func NewDefaultLogRetentionPolicyResource() resource.Resource {
	return &defaultLogRetentionPolicyResource{}
}

// logRetentionPolicyResource is the resource implementation.
type logRetentionPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultLogRetentionPolicyResource is the resource implementation.
type defaultLogRetentionPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *logRetentionPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_retention_policy"
}

func (r *defaultLogRetentionPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_log_retention_policy"
}

// Configure adds the provider configured client to the resource.
func (r *logRetentionPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultLogRetentionPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type logRetentionPolicyResourceModel struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	Type            types.String `tfsdk:"type"`
	DiskSpaceUsed   types.String `tfsdk:"disk_space_used"`
	FreeDiskSpace   types.String `tfsdk:"free_disk_space"`
	NumberOfFiles   types.Int64  `tfsdk:"number_of_files"`
	RetainDuration  types.String `tfsdk:"retain_duration"`
	Description     types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *logRetentionPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	logRetentionPolicySchema(ctx, req, resp, false)
}

func (r *defaultLogRetentionPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	logRetentionPolicySchema(ctx, req, resp, true)
}

func logRetentionPolicySchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Log Retention Policy.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Log Retention Policy resource. Options are ['time-limit', 'never-delete', 'file-count', 'free-disk-space', 'size-limit']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"time-limit", "never-delete", "file-count", "free-disk-space", "size-limit"}...),
				},
			},
			"disk_space_used": schema.StringAttribute{
				Description: "Specifies the maximum total disk space used by the log files.",
				Optional:    true,
			},
			"free_disk_space": schema.StringAttribute{
				Description: "Specifies the minimum amount of free disk space that should be available on the file system on which the archived log files are stored.",
				Optional:    true,
			},
			"number_of_files": schema.Int64Attribute{
				Description: "Specifies the number of archived log files to retain before the oldest ones are cleaned.",
				Optional:    true,
			},
			"retain_duration": schema.StringAttribute{
				Description: "Specifies the desired minimum length of time that each log file should be retained.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log Retention Policy",
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
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsLogRetentionPolicy() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("retain_duration"),
			path.MatchRoot("type"),
			[]string{"time-limit"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("number_of_files"),
			path.MatchRoot("type"),
			[]string{"file-count"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("free_disk_space"),
			path.MatchRoot("type"),
			[]string{"free-disk-space"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("disk_space_used"),
			path.MatchRoot("type"),
			[]string{"size-limit"},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"time-limit",
			[]path.Expression{path.MatchRoot("retain_duration")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"file-count",
			[]path.Expression{path.MatchRoot("number_of_files")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"free-disk-space",
			[]path.Expression{path.MatchRoot("free_disk_space")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"size-limit",
			[]path.Expression{path.MatchRoot("disk_space_used")},
		),
	}
}

// Add config validators
func (r logRetentionPolicyResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsLogRetentionPolicy()
}

// Add config validators
func (r defaultLogRetentionPolicyResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsLogRetentionPolicy()
}

// Add optional fields to create request for time-limit log-retention-policy
func addOptionalTimeLimitLogRetentionPolicyFields(ctx context.Context, addRequest *client.AddTimeLimitLogRetentionPolicyRequest, plan logRetentionPolicyResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for never-delete log-retention-policy
func addOptionalNeverDeleteLogRetentionPolicyFields(ctx context.Context, addRequest *client.AddNeverDeleteLogRetentionPolicyRequest, plan logRetentionPolicyResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for file-count log-retention-policy
func addOptionalFileCountLogRetentionPolicyFields(ctx context.Context, addRequest *client.AddFileCountLogRetentionPolicyRequest, plan logRetentionPolicyResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for free-disk-space log-retention-policy
func addOptionalFreeDiskSpaceLogRetentionPolicyFields(ctx context.Context, addRequest *client.AddFreeDiskSpaceLogRetentionPolicyRequest, plan logRetentionPolicyResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for size-limit log-retention-policy
func addOptionalSizeLimitLogRetentionPolicyFields(ctx context.Context, addRequest *client.AddSizeLimitLogRetentionPolicyRequest, plan logRetentionPolicyResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *logRetentionPolicyResourceModel) populateAllComputedStringAttributes() {
	if model.RetainDuration.IsUnknown() || model.RetainDuration.IsNull() {
		model.RetainDuration = types.StringValue("")
	}
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.DiskSpaceUsed.IsUnknown() || model.DiskSpaceUsed.IsNull() {
		model.DiskSpaceUsed = types.StringValue("")
	}
	if model.FreeDiskSpace.IsUnknown() || model.FreeDiskSpace.IsNull() {
		model.FreeDiskSpace = types.StringValue("")
	}
}

// Read a TimeLimitLogRetentionPolicyResponse object into the model struct
func readTimeLimitLogRetentionPolicyResponse(ctx context.Context, r *client.TimeLimitLogRetentionPolicyResponse, state *logRetentionPolicyResourceModel, expectedValues *logRetentionPolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("time-limit")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.RetainDuration = types.StringValue(r.RetainDuration)
	config.CheckMismatchedPDFormattedAttributes("retain_duration",
		expectedValues.RetainDuration, state.RetainDuration, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Read a NeverDeleteLogRetentionPolicyResponse object into the model struct
func readNeverDeleteLogRetentionPolicyResponse(ctx context.Context, r *client.NeverDeleteLogRetentionPolicyResponse, state *logRetentionPolicyResourceModel, expectedValues *logRetentionPolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("never-delete")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Read a FileCountLogRetentionPolicyResponse object into the model struct
func readFileCountLogRetentionPolicyResponse(ctx context.Context, r *client.FileCountLogRetentionPolicyResponse, state *logRetentionPolicyResourceModel, expectedValues *logRetentionPolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-count")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.NumberOfFiles = types.Int64Value(r.NumberOfFiles)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Read a FreeDiskSpaceLogRetentionPolicyResponse object into the model struct
func readFreeDiskSpaceLogRetentionPolicyResponse(ctx context.Context, r *client.FreeDiskSpaceLogRetentionPolicyResponse, state *logRetentionPolicyResourceModel, expectedValues *logRetentionPolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("free-disk-space")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.FreeDiskSpace = types.StringValue(r.FreeDiskSpace)
	config.CheckMismatchedPDFormattedAttributes("free_disk_space",
		expectedValues.FreeDiskSpace, state.FreeDiskSpace, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Read a SizeLimitLogRetentionPolicyResponse object into the model struct
func readSizeLimitLogRetentionPolicyResponse(ctx context.Context, r *client.SizeLimitLogRetentionPolicyResponse, state *logRetentionPolicyResourceModel, expectedValues *logRetentionPolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("size-limit")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.DiskSpaceUsed = types.StringValue(r.DiskSpaceUsed)
	config.CheckMismatchedPDFormattedAttributes("disk_space_used",
		expectedValues.DiskSpaceUsed, state.DiskSpaceUsed, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createLogRetentionPolicyOperations(plan logRetentionPolicyResourceModel, state logRetentionPolicyResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.DiskSpaceUsed, state.DiskSpaceUsed, "disk-space-used")
	operations.AddStringOperationIfNecessary(&ops, plan.FreeDiskSpace, state.FreeDiskSpace, "free-disk-space")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumberOfFiles, state.NumberOfFiles, "number-of-files")
	operations.AddStringOperationIfNecessary(&ops, plan.RetainDuration, state.RetainDuration, "retain-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a time-limit log-retention-policy
func (r *logRetentionPolicyResource) CreateTimeLimitLogRetentionPolicy(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logRetentionPolicyResourceModel) (*logRetentionPolicyResourceModel, error) {
	addRequest := client.NewAddTimeLimitLogRetentionPolicyRequest([]client.EnumtimeLimitLogRetentionPolicySchemaUrn{client.ENUMTIMELIMITLOGRETENTIONPOLICYSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_RETENTION_POLICYTIME_LIMIT},
		plan.RetainDuration.ValueString(),
		plan.Name.ValueString())
	addOptionalTimeLimitLogRetentionPolicyFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogRetentionPolicyAPI.AddLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogRetentionPolicyRequest(
		client.AddTimeLimitLogRetentionPolicyRequestAsAddLogRetentionPolicyRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogRetentionPolicyAPI.AddLogRetentionPolicyExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Retention Policy", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logRetentionPolicyResourceModel
	readTimeLimitLogRetentionPolicyResponse(ctx, addResponse.TimeLimitLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a never-delete log-retention-policy
func (r *logRetentionPolicyResource) CreateNeverDeleteLogRetentionPolicy(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logRetentionPolicyResourceModel) (*logRetentionPolicyResourceModel, error) {
	addRequest := client.NewAddNeverDeleteLogRetentionPolicyRequest([]client.EnumneverDeleteLogRetentionPolicySchemaUrn{client.ENUMNEVERDELETELOGRETENTIONPOLICYSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_RETENTION_POLICYNEVER_DELETE},
		plan.Name.ValueString())
	addOptionalNeverDeleteLogRetentionPolicyFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogRetentionPolicyAPI.AddLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogRetentionPolicyRequest(
		client.AddNeverDeleteLogRetentionPolicyRequestAsAddLogRetentionPolicyRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogRetentionPolicyAPI.AddLogRetentionPolicyExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Retention Policy", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logRetentionPolicyResourceModel
	readNeverDeleteLogRetentionPolicyResponse(ctx, addResponse.NeverDeleteLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a file-count log-retention-policy
func (r *logRetentionPolicyResource) CreateFileCountLogRetentionPolicy(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logRetentionPolicyResourceModel) (*logRetentionPolicyResourceModel, error) {
	addRequest := client.NewAddFileCountLogRetentionPolicyRequest([]client.EnumfileCountLogRetentionPolicySchemaUrn{client.ENUMFILECOUNTLOGRETENTIONPOLICYSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_RETENTION_POLICYFILE_COUNT},
		plan.NumberOfFiles.ValueInt64(),
		plan.Name.ValueString())
	addOptionalFileCountLogRetentionPolicyFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogRetentionPolicyAPI.AddLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogRetentionPolicyRequest(
		client.AddFileCountLogRetentionPolicyRequestAsAddLogRetentionPolicyRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogRetentionPolicyAPI.AddLogRetentionPolicyExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Retention Policy", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logRetentionPolicyResourceModel
	readFileCountLogRetentionPolicyResponse(ctx, addResponse.FileCountLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a free-disk-space log-retention-policy
func (r *logRetentionPolicyResource) CreateFreeDiskSpaceLogRetentionPolicy(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logRetentionPolicyResourceModel) (*logRetentionPolicyResourceModel, error) {
	addRequest := client.NewAddFreeDiskSpaceLogRetentionPolicyRequest([]client.EnumfreeDiskSpaceLogRetentionPolicySchemaUrn{client.ENUMFREEDISKSPACELOGRETENTIONPOLICYSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_RETENTION_POLICYFREE_DISK_SPACE},
		plan.FreeDiskSpace.ValueString(),
		plan.Name.ValueString())
	addOptionalFreeDiskSpaceLogRetentionPolicyFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogRetentionPolicyAPI.AddLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogRetentionPolicyRequest(
		client.AddFreeDiskSpaceLogRetentionPolicyRequestAsAddLogRetentionPolicyRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogRetentionPolicyAPI.AddLogRetentionPolicyExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Retention Policy", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logRetentionPolicyResourceModel
	readFreeDiskSpaceLogRetentionPolicyResponse(ctx, addResponse.FreeDiskSpaceLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a size-limit log-retention-policy
func (r *logRetentionPolicyResource) CreateSizeLimitLogRetentionPolicy(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logRetentionPolicyResourceModel) (*logRetentionPolicyResourceModel, error) {
	addRequest := client.NewAddSizeLimitLogRetentionPolicyRequest([]client.EnumsizeLimitLogRetentionPolicySchemaUrn{client.ENUMSIZELIMITLOGRETENTIONPOLICYSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_RETENTION_POLICYSIZE_LIMIT},
		plan.DiskSpaceUsed.ValueString(),
		plan.Name.ValueString())
	addOptionalSizeLimitLogRetentionPolicyFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogRetentionPolicyAPI.AddLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogRetentionPolicyRequest(
		client.AddSizeLimitLogRetentionPolicyRequestAsAddLogRetentionPolicyRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogRetentionPolicyAPI.AddLogRetentionPolicyExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Retention Policy", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logRetentionPolicyResourceModel
	readSizeLimitLogRetentionPolicyResponse(ctx, addResponse.SizeLimitLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *logRetentionPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan logRetentionPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *logRetentionPolicyResourceModel
	var err error
	if plan.Type.ValueString() == "time-limit" {
		state, err = r.CreateTimeLimitLogRetentionPolicy(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "never-delete" {
		state, err = r.CreateNeverDeleteLogRetentionPolicy(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "file-count" {
		state, err = r.CreateFileCountLogRetentionPolicy(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "free-disk-space" {
		state, err = r.CreateFreeDiskSpaceLogRetentionPolicy(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "size-limit" {
		state, err = r.CreateSizeLimitLogRetentionPolicy(ctx, req, resp, plan)
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
func (r *defaultLogRetentionPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan logRetentionPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogRetentionPolicyAPI.GetLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Log Retention Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state logRetentionPolicyResourceModel
	if readResponse.TimeLimitLogRetentionPolicyResponse != nil {
		readTimeLimitLogRetentionPolicyResponse(ctx, readResponse.TimeLimitLogRetentionPolicyResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.NeverDeleteLogRetentionPolicyResponse != nil {
		readNeverDeleteLogRetentionPolicyResponse(ctx, readResponse.NeverDeleteLogRetentionPolicyResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileCountLogRetentionPolicyResponse != nil {
		readFileCountLogRetentionPolicyResponse(ctx, readResponse.FileCountLogRetentionPolicyResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FreeDiskSpaceLogRetentionPolicyResponse != nil {
		readFreeDiskSpaceLogRetentionPolicyResponse(ctx, readResponse.FreeDiskSpaceLogRetentionPolicyResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SizeLimitLogRetentionPolicyResponse != nil {
		readSizeLimitLogRetentionPolicyResponse(ctx, readResponse.SizeLimitLogRetentionPolicyResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogRetentionPolicyAPI.UpdateLogRetentionPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createLogRetentionPolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogRetentionPolicyAPI.UpdateLogRetentionPolicyExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Log Retention Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.TimeLimitLogRetentionPolicyResponse != nil {
			readTimeLimitLogRetentionPolicyResponse(ctx, updateResponse.TimeLimitLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.NeverDeleteLogRetentionPolicyResponse != nil {
			readNeverDeleteLogRetentionPolicyResponse(ctx, updateResponse.NeverDeleteLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileCountLogRetentionPolicyResponse != nil {
			readFileCountLogRetentionPolicyResponse(ctx, updateResponse.FileCountLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FreeDiskSpaceLogRetentionPolicyResponse != nil {
			readFreeDiskSpaceLogRetentionPolicyResponse(ctx, updateResponse.FreeDiskSpaceLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SizeLimitLogRetentionPolicyResponse != nil {
			readSizeLimitLogRetentionPolicyResponse(ctx, updateResponse.SizeLimitLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
	}

	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *logRetentionPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLogRetentionPolicy(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultLogRetentionPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLogRetentionPolicy(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readLogRetentionPolicy(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state logRetentionPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogRetentionPolicyAPI.GetLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Log Retention Policy", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Log Retention Policy", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.TimeLimitLogRetentionPolicyResponse != nil {
		readTimeLimitLogRetentionPolicyResponse(ctx, readResponse.TimeLimitLogRetentionPolicyResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.NeverDeleteLogRetentionPolicyResponse != nil {
		readNeverDeleteLogRetentionPolicyResponse(ctx, readResponse.NeverDeleteLogRetentionPolicyResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileCountLogRetentionPolicyResponse != nil {
		readFileCountLogRetentionPolicyResponse(ctx, readResponse.FileCountLogRetentionPolicyResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FreeDiskSpaceLogRetentionPolicyResponse != nil {
		readFreeDiskSpaceLogRetentionPolicyResponse(ctx, readResponse.FreeDiskSpaceLogRetentionPolicyResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SizeLimitLogRetentionPolicyResponse != nil {
		readSizeLimitLogRetentionPolicyResponse(ctx, readResponse.SizeLimitLogRetentionPolicyResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *logRetentionPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLogRetentionPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLogRetentionPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLogRetentionPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateLogRetentionPolicy(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan logRetentionPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state logRetentionPolicyResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogRetentionPolicyAPI.UpdateLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createLogRetentionPolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogRetentionPolicyAPI.UpdateLogRetentionPolicyExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Log Retention Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.TimeLimitLogRetentionPolicyResponse != nil {
			readTimeLimitLogRetentionPolicyResponse(ctx, updateResponse.TimeLimitLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.NeverDeleteLogRetentionPolicyResponse != nil {
			readNeverDeleteLogRetentionPolicyResponse(ctx, updateResponse.NeverDeleteLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileCountLogRetentionPolicyResponse != nil {
			readFileCountLogRetentionPolicyResponse(ctx, updateResponse.FileCountLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FreeDiskSpaceLogRetentionPolicyResponse != nil {
			readFreeDiskSpaceLogRetentionPolicyResponse(ctx, updateResponse.FreeDiskSpaceLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SizeLimitLogRetentionPolicyResponse != nil {
			readSizeLimitLogRetentionPolicyResponse(ctx, updateResponse.SizeLimitLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultLogRetentionPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *logRetentionPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state logRetentionPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogRetentionPolicyAPI.DeleteLogRetentionPolicyExecute(r.apiClient.LogRetentionPolicyAPI.DeleteLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Log Retention Policy", err, httpResp)
		return
	}
}

func (r *logRetentionPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLogRetentionPolicy(ctx, req, resp)
}

func (r *defaultLogRetentionPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLogRetentionPolicy(ctx, req, resp)
}

func importLogRetentionPolicy(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
