package logfilerotationlistener

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10100/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &logFileRotationListenerResource{}
	_ resource.ResourceWithConfigure   = &logFileRotationListenerResource{}
	_ resource.ResourceWithImportState = &logFileRotationListenerResource{}
	_ resource.Resource                = &defaultLogFileRotationListenerResource{}
	_ resource.ResourceWithConfigure   = &defaultLogFileRotationListenerResource{}
	_ resource.ResourceWithImportState = &defaultLogFileRotationListenerResource{}
)

// Create a Log File Rotation Listener resource
func NewLogFileRotationListenerResource() resource.Resource {
	return &logFileRotationListenerResource{}
}

func NewDefaultLogFileRotationListenerResource() resource.Resource {
	return &defaultLogFileRotationListenerResource{}
}

// logFileRotationListenerResource is the resource implementation.
type logFileRotationListenerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultLogFileRotationListenerResource is the resource implementation.
type defaultLogFileRotationListenerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *logFileRotationListenerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_file_rotation_listener"
}

func (r *defaultLogFileRotationListenerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_log_file_rotation_listener"
}

// Configure adds the provider configured client to the resource.
func (r *logFileRotationListenerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultLogFileRotationListenerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type logFileRotationListenerResourceModel struct {
	Id                                   types.String `tfsdk:"id"`
	Name                                 types.String `tfsdk:"name"`
	Notifications                        types.Set    `tfsdk:"notifications"`
	RequiredActions                      types.Set    `tfsdk:"required_actions"`
	Type                                 types.String `tfsdk:"type"`
	ExtensionClass                       types.String `tfsdk:"extension_class"`
	ExtensionArgument                    types.Set    `tfsdk:"extension_argument"`
	CopyToDirectory                      types.String `tfsdk:"copy_to_directory"`
	CompressOnCopy                       types.Bool   `tfsdk:"compress_on_copy"`
	OutputDirectory                      types.String `tfsdk:"output_directory"`
	AwsExternalServer                    types.String `tfsdk:"aws_external_server"`
	S3BucketName                         types.String `tfsdk:"s3_bucket_name"`
	TargetThroughputInMegabitsPerSecond  types.Int64  `tfsdk:"target_throughput_in_megabits_per_second"`
	MaximumConcurrentTransferConnections types.Int64  `tfsdk:"maximum_concurrent_transfer_connections"`
	MaximumFileCountToRetain             types.Int64  `tfsdk:"maximum_file_count_to_retain"`
	MaximumFileAgeToRetain               types.String `tfsdk:"maximum_file_age_to_retain"`
	FileRetentionPattern                 types.String `tfsdk:"file_retention_pattern"`
	Description                          types.String `tfsdk:"description"`
	Enabled                              types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *logFileRotationListenerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	logFileRotationListenerSchema(ctx, req, resp, false)
}

func (r *defaultLogFileRotationListenerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	logFileRotationListenerSchema(ctx, req, resp, true)
}

func logFileRotationListenerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Log File Rotation Listener.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Log File Rotation Listener resource. Options are ['upload-to-s3', 'summarize', 'copy', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"upload-to-s3", "summarize", "copy", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Log File Rotation Listener.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Log File Rotation Listener. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"copy_to_directory": schema.StringAttribute{
				Description: "The path to the directory to which log files should be copied. It must be different from the directory to which the log file is originally written, and administrators should ensure that the filesystem has sufficient space to hold files as they are copied.",
				Optional:    true,
			},
			"compress_on_copy": schema.BoolAttribute{
				Description: "Indicates whether the file should be gzip-compressed as it is copied into the destination directory.",
				Optional:    true,
				Computed:    true,
			},
			"output_directory": schema.StringAttribute{
				Description: "The path to the directory in which the summarize-access-log output should be written. If no value is provided, the output file will be written into the same directory as the rotated log file.",
				Optional:    true,
			},
			"aws_external_server": schema.StringAttribute{
				Description: "The external server with information to use when interacting with the AWS S3 service.",
				Optional:    true,
			},
			"s3_bucket_name": schema.StringAttribute{
				Description: "The name of the S3 bucket into which rotated log files should be copied.",
				Optional:    true,
			},
			"target_throughput_in_megabits_per_second": schema.Int64Attribute{
				Description: "The target throughput to attempt to achieve for data transfers to or from S3, in megabits per second.",
				Optional:    true,
			},
			"maximum_concurrent_transfer_connections": schema.Int64Attribute{
				Description: "The maximum number of concurrent connections that may be used when transferring data to or from S3.",
				Optional:    true,
			},
			"maximum_file_count_to_retain": schema.Int64Attribute{
				Description: "The maximum number of existing files matching the file retention pattern that should be retained in the S3 bucket after successfully uploading a newly rotated file.",
				Optional:    true,
			},
			"maximum_file_age_to_retain": schema.StringAttribute{
				Description: "The maximum length of time to retain files matching the file retention pattern that should be retained in the S3 bucket after successfully uploading a newly rotated file.",
				Optional:    true,
			},
			"file_retention_pattern": schema.StringAttribute{
				Description: "A regular expression pattern that will be used to identify which files are candidates for automatic removal based on the maximum-file-count-to-retain and maximum-file-age-to-retain properties. By default, all files in the bucket will be eligible for removal by retention processing.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log File Rotation Listener",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Log File Rotation Listener is enabled for use.",
				Required:    true,
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
	} else {
		// Add RequiresReplace modifier for read-only attributes
		extensionClassAttr := schemaDef.Attributes["extension_class"].(schema.StringAttribute)
		extensionClassAttr.PlanModifiers = append(extensionClassAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["extension_class"] = extensionClassAttr
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *logFileRotationListenerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanLogFileRotationListener(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_log_file_rotation_listener")
	var planModel, configModel logFileRotationListenerResourceModel
	req.Config.Get(ctx, &configModel)
	req.Plan.Get(ctx, &planModel)
	resourceType := planModel.Type.ValueString()
	anyDefaultsSet := false
	// Set defaults for upload-to-s3 type
	if resourceType == "upload-to-s3" {
		if !internaltypes.IsDefined(configModel.FileRetentionPattern) {
			defaultVal := types.StringValue(".*")
			if !planModel.FileRetentionPattern.Equal(defaultVal) {
				planModel.FileRetentionPattern = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for copy type
	if resourceType == "copy" {
		if !internaltypes.IsDefined(configModel.CompressOnCopy) {
			defaultVal := types.BoolValue(false)
			if !planModel.CompressOnCopy.Equal(defaultVal) {
				planModel.CompressOnCopy = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	if anyDefaultsSet {
		planModel.Notifications = types.SetUnknown(types.StringType)
		planModel.RequiredActions = types.SetUnknown(config.GetRequiredActionsObjectType())
	}
	planModel.setNotApplicableAttrsNull()
	resp.Plan.Set(ctx, &planModel)
}

func (r *defaultLogFileRotationListenerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanLogFileRotationListener(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_default_log_file_rotation_listener")
}

func modifyPlanLogFileRotationListener(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, resourceName string) {
	compare, err := version.Compare(providerConfig.ProductVersion, version.PingDirectory10000)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model logFileRotationListenerResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.Type) && model.Type.ValueString() == "upload-to-s3" {
		version.CheckResourceSupported(&resp.Diagnostics, version.PingDirectory10000,
			providerConfig.ProductVersion, resourceName+" with type \"upload_to_s3\"")
	}
}

func (model *logFileRotationListenerResourceModel) setNotApplicableAttrsNull() {
	resourceType := model.Type.ValueString()
	// Set any not applicable computed attributes to null for each type
	if resourceType == "upload-to-s3" {
		model.CompressOnCopy = types.BoolNull()
	}
	if resourceType == "summarize" {
		model.FileRetentionPattern = types.StringNull()
		model.CompressOnCopy = types.BoolNull()
	}
	if resourceType == "copy" {
		model.FileRetentionPattern = types.StringNull()
	}
	if resourceType == "third-party" {
		model.FileRetentionPattern = types.StringNull()
		model.CompressOnCopy = types.BoolNull()
	}
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsLogFileRotationListener() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("aws_external_server"),
			path.MatchRoot("type"),
			[]string{"upload-to-s3"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("s3_bucket_name"),
			path.MatchRoot("type"),
			[]string{"upload-to-s3"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("target_throughput_in_megabits_per_second"),
			path.MatchRoot("type"),
			[]string{"upload-to-s3"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("maximum_concurrent_transfer_connections"),
			path.MatchRoot("type"),
			[]string{"upload-to-s3"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("maximum_file_count_to_retain"),
			path.MatchRoot("type"),
			[]string{"upload-to-s3"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("maximum_file_age_to_retain"),
			path.MatchRoot("type"),
			[]string{"upload-to-s3"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("file_retention_pattern"),
			path.MatchRoot("type"),
			[]string{"upload-to-s3"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("output_directory"),
			path.MatchRoot("type"),
			[]string{"summarize"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("copy_to_directory"),
			path.MatchRoot("type"),
			[]string{"copy"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("compress_on_copy"),
			path.MatchRoot("type"),
			[]string{"copy"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_class"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_argument"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"upload-to-s3",
			[]path.Expression{path.MatchRoot("aws_external_server"), path.MatchRoot("s3_bucket_name")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"copy",
			[]path.Expression{path.MatchRoot("copy_to_directory")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party",
			[]path.Expression{path.MatchRoot("extension_class")},
		),
	}
}

// Add config validators
func (r logFileRotationListenerResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsLogFileRotationListener()
}

// Add config validators
func (r defaultLogFileRotationListenerResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsLogFileRotationListener()
}

// Add optional fields to create request for upload-to-s3 log-file-rotation-listener
func addOptionalUploadToS3LogFileRotationListenerFields(ctx context.Context, addRequest *client.AddUploadToS3LogFileRotationListenerRequest, plan logFileRotationListenerResourceModel) {
	if internaltypes.IsDefined(plan.TargetThroughputInMegabitsPerSecond) {
		addRequest.TargetThroughputInMegabitsPerSecond = plan.TargetThroughputInMegabitsPerSecond.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaximumConcurrentTransferConnections) {
		addRequest.MaximumConcurrentTransferConnections = plan.MaximumConcurrentTransferConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaximumFileCountToRetain) {
		addRequest.MaximumFileCountToRetain = plan.MaximumFileCountToRetain.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaximumFileAgeToRetain) {
		addRequest.MaximumFileAgeToRetain = plan.MaximumFileAgeToRetain.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.FileRetentionPattern) {
		addRequest.FileRetentionPattern = plan.FileRetentionPattern.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for summarize log-file-rotation-listener
func addOptionalSummarizeLogFileRotationListenerFields(ctx context.Context, addRequest *client.AddSummarizeLogFileRotationListenerRequest, plan logFileRotationListenerResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OutputDirectory) {
		addRequest.OutputDirectory = plan.OutputDirectory.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for copy log-file-rotation-listener
func addOptionalCopyLogFileRotationListenerFields(ctx context.Context, addRequest *client.AddCopyLogFileRotationListenerRequest, plan logFileRotationListenerResourceModel) {
	if internaltypes.IsDefined(plan.CompressOnCopy) {
		addRequest.CompressOnCopy = plan.CompressOnCopy.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for third-party log-file-rotation-listener
func addOptionalThirdPartyLogFileRotationListenerFields(ctx context.Context, addRequest *client.AddThirdPartyLogFileRotationListenerRequest, plan logFileRotationListenerResourceModel) {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateLogFileRotationListenerUnknownValues(model *logFileRotationListenerResourceModel) {
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *logFileRotationListenerResourceModel) populateAllComputedStringAttributes() {
	if model.MaximumFileAgeToRetain.IsUnknown() || model.MaximumFileAgeToRetain.IsNull() {
		model.MaximumFileAgeToRetain = types.StringValue("")
	}
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.OutputDirectory.IsUnknown() || model.OutputDirectory.IsNull() {
		model.OutputDirectory = types.StringValue("")
	}
	if model.ExtensionClass.IsUnknown() || model.ExtensionClass.IsNull() {
		model.ExtensionClass = types.StringValue("")
	}
	if model.AwsExternalServer.IsUnknown() || model.AwsExternalServer.IsNull() {
		model.AwsExternalServer = types.StringValue("")
	}
	if model.S3BucketName.IsUnknown() || model.S3BucketName.IsNull() {
		model.S3BucketName = types.StringValue("")
	}
	if model.CopyToDirectory.IsUnknown() || model.CopyToDirectory.IsNull() {
		model.CopyToDirectory = types.StringValue("")
	}
	if model.FileRetentionPattern.IsUnknown() || model.FileRetentionPattern.IsNull() {
		model.FileRetentionPattern = types.StringValue("")
	}
}

// Read a UploadToS3LogFileRotationListenerResponse object into the model struct
func readUploadToS3LogFileRotationListenerResponse(ctx context.Context, r *client.UploadToS3LogFileRotationListenerResponse, state *logFileRotationListenerResourceModel, expectedValues *logFileRotationListenerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("upload-to-s3")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AwsExternalServer = types.StringValue(r.AwsExternalServer)
	state.S3BucketName = types.StringValue(r.S3BucketName)
	state.TargetThroughputInMegabitsPerSecond = internaltypes.Int64TypeOrNil(r.TargetThroughputInMegabitsPerSecond)
	state.MaximumConcurrentTransferConnections = internaltypes.Int64TypeOrNil(r.MaximumConcurrentTransferConnections)
	state.MaximumFileCountToRetain = internaltypes.Int64TypeOrNil(r.MaximumFileCountToRetain)
	state.MaximumFileAgeToRetain = internaltypes.StringTypeOrNil(r.MaximumFileAgeToRetain, internaltypes.IsEmptyString(expectedValues.MaximumFileAgeToRetain))
	config.CheckMismatchedPDFormattedAttributes("maximum_file_age_to_retain",
		expectedValues.MaximumFileAgeToRetain, state.MaximumFileAgeToRetain, diagnostics)
	state.FileRetentionPattern = internaltypes.StringTypeOrNil(r.FileRetentionPattern, true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogFileRotationListenerUnknownValues(state)
}

// Read a SummarizeLogFileRotationListenerResponse object into the model struct
func readSummarizeLogFileRotationListenerResponse(ctx context.Context, r *client.SummarizeLogFileRotationListenerResponse, state *logFileRotationListenerResourceModel, expectedValues *logFileRotationListenerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("summarize")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.OutputDirectory = internaltypes.StringTypeOrNil(r.OutputDirectory, internaltypes.IsEmptyString(expectedValues.OutputDirectory))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogFileRotationListenerUnknownValues(state)
}

// Read a CopyLogFileRotationListenerResponse object into the model struct
func readCopyLogFileRotationListenerResponse(ctx context.Context, r *client.CopyLogFileRotationListenerResponse, state *logFileRotationListenerResourceModel, expectedValues *logFileRotationListenerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("copy")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.CopyToDirectory = types.StringValue(r.CopyToDirectory)
	state.CompressOnCopy = internaltypes.BoolTypeOrNil(r.CompressOnCopy)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogFileRotationListenerUnknownValues(state)
}

// Read a ThirdPartyLogFileRotationListenerResponse object into the model struct
func readThirdPartyLogFileRotationListenerResponse(ctx context.Context, r *client.ThirdPartyLogFileRotationListenerResponse, state *logFileRotationListenerResourceModel, expectedValues *logFileRotationListenerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogFileRotationListenerUnknownValues(state)
}

// Create any update operations necessary to make the state match the plan
func createLogFileRotationListenerOperations(plan logFileRotationListenerResourceModel, state logFileRotationListenerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.CopyToDirectory, state.CopyToDirectory, "copy-to-directory")
	operations.AddBoolOperationIfNecessary(&ops, plan.CompressOnCopy, state.CompressOnCopy, "compress-on-copy")
	operations.AddStringOperationIfNecessary(&ops, plan.OutputDirectory, state.OutputDirectory, "output-directory")
	operations.AddStringOperationIfNecessary(&ops, plan.AwsExternalServer, state.AwsExternalServer, "aws-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.S3BucketName, state.S3BucketName, "s3-bucket-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.TargetThroughputInMegabitsPerSecond, state.TargetThroughputInMegabitsPerSecond, "target-throughput-in-megabits-per-second")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumConcurrentTransferConnections, state.MaximumConcurrentTransferConnections, "maximum-concurrent-transfer-connections")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumFileCountToRetain, state.MaximumFileCountToRetain, "maximum-file-count-to-retain")
	operations.AddStringOperationIfNecessary(&ops, plan.MaximumFileAgeToRetain, state.MaximumFileAgeToRetain, "maximum-file-age-to-retain")
	operations.AddStringOperationIfNecessary(&ops, plan.FileRetentionPattern, state.FileRetentionPattern, "file-retention-pattern")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a upload-to-s3 log-file-rotation-listener
func (r *logFileRotationListenerResource) CreateUploadToS3LogFileRotationListener(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logFileRotationListenerResourceModel) (*logFileRotationListenerResourceModel, error) {
	addRequest := client.NewAddUploadToS3LogFileRotationListenerRequest([]client.EnumuploadToS3LogFileRotationListenerSchemaUrn{client.ENUMUPLOADTOS3LOGFILEROTATIONLISTENERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_FILE_ROTATION_LISTENERUPLOAD_TO_S3},
		plan.AwsExternalServer.ValueString(),
		plan.S3BucketName.ValueString(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	addOptionalUploadToS3LogFileRotationListenerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogFileRotationListenerAPI.AddLogFileRotationListener(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogFileRotationListenerRequest(
		client.AddUploadToS3LogFileRotationListenerRequestAsAddLogFileRotationListenerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogFileRotationListenerAPI.AddLogFileRotationListenerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log File Rotation Listener", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logFileRotationListenerResourceModel
	readUploadToS3LogFileRotationListenerResponse(ctx, addResponse.UploadToS3LogFileRotationListenerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a summarize log-file-rotation-listener
func (r *logFileRotationListenerResource) CreateSummarizeLogFileRotationListener(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logFileRotationListenerResourceModel) (*logFileRotationListenerResourceModel, error) {
	addRequest := client.NewAddSummarizeLogFileRotationListenerRequest([]client.EnumsummarizeLogFileRotationListenerSchemaUrn{client.ENUMSUMMARIZELOGFILEROTATIONLISTENERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_FILE_ROTATION_LISTENERSUMMARIZE},
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	addOptionalSummarizeLogFileRotationListenerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogFileRotationListenerAPI.AddLogFileRotationListener(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogFileRotationListenerRequest(
		client.AddSummarizeLogFileRotationListenerRequestAsAddLogFileRotationListenerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogFileRotationListenerAPI.AddLogFileRotationListenerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log File Rotation Listener", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logFileRotationListenerResourceModel
	readSummarizeLogFileRotationListenerResponse(ctx, addResponse.SummarizeLogFileRotationListenerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a copy log-file-rotation-listener
func (r *logFileRotationListenerResource) CreateCopyLogFileRotationListener(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logFileRotationListenerResourceModel) (*logFileRotationListenerResourceModel, error) {
	addRequest := client.NewAddCopyLogFileRotationListenerRequest([]client.EnumcopyLogFileRotationListenerSchemaUrn{client.ENUMCOPYLOGFILEROTATIONLISTENERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_FILE_ROTATION_LISTENERCOPY},
		plan.CopyToDirectory.ValueString(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	addOptionalCopyLogFileRotationListenerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogFileRotationListenerAPI.AddLogFileRotationListener(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogFileRotationListenerRequest(
		client.AddCopyLogFileRotationListenerRequestAsAddLogFileRotationListenerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogFileRotationListenerAPI.AddLogFileRotationListenerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log File Rotation Listener", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logFileRotationListenerResourceModel
	readCopyLogFileRotationListenerResponse(ctx, addResponse.CopyLogFileRotationListenerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party log-file-rotation-listener
func (r *logFileRotationListenerResource) CreateThirdPartyLogFileRotationListener(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logFileRotationListenerResourceModel) (*logFileRotationListenerResourceModel, error) {
	addRequest := client.NewAddThirdPartyLogFileRotationListenerRequest([]client.EnumthirdPartyLogFileRotationListenerSchemaUrn{client.ENUMTHIRDPARTYLOGFILEROTATIONLISTENERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_FILE_ROTATION_LISTENERTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	addOptionalThirdPartyLogFileRotationListenerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogFileRotationListenerAPI.AddLogFileRotationListener(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogFileRotationListenerRequest(
		client.AddThirdPartyLogFileRotationListenerRequestAsAddLogFileRotationListenerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogFileRotationListenerAPI.AddLogFileRotationListenerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log File Rotation Listener", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logFileRotationListenerResourceModel
	readThirdPartyLogFileRotationListenerResponse(ctx, addResponse.ThirdPartyLogFileRotationListenerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *logFileRotationListenerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan logFileRotationListenerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *logFileRotationListenerResourceModel
	var err error
	if plan.Type.ValueString() == "upload-to-s3" {
		state, err = r.CreateUploadToS3LogFileRotationListener(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "summarize" {
		state, err = r.CreateSummarizeLogFileRotationListener(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "copy" {
		state, err = r.CreateCopyLogFileRotationListener(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyLogFileRotationListener(ctx, req, resp, plan)
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
func (r *defaultLogFileRotationListenerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan logFileRotationListenerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogFileRotationListenerAPI.GetLogFileRotationListener(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Log File Rotation Listener", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state logFileRotationListenerResourceModel
	if readResponse.UploadToS3LogFileRotationListenerResponse != nil {
		readUploadToS3LogFileRotationListenerResponse(ctx, readResponse.UploadToS3LogFileRotationListenerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SummarizeLogFileRotationListenerResponse != nil {
		readSummarizeLogFileRotationListenerResponse(ctx, readResponse.SummarizeLogFileRotationListenerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CopyLogFileRotationListenerResponse != nil {
		readCopyLogFileRotationListenerResponse(ctx, readResponse.CopyLogFileRotationListenerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyLogFileRotationListenerResponse != nil {
		readThirdPartyLogFileRotationListenerResponse(ctx, readResponse.ThirdPartyLogFileRotationListenerResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogFileRotationListenerAPI.UpdateLogFileRotationListener(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createLogFileRotationListenerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogFileRotationListenerAPI.UpdateLogFileRotationListenerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Log File Rotation Listener", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.UploadToS3LogFileRotationListenerResponse != nil {
			readUploadToS3LogFileRotationListenerResponse(ctx, updateResponse.UploadToS3LogFileRotationListenerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SummarizeLogFileRotationListenerResponse != nil {
			readSummarizeLogFileRotationListenerResponse(ctx, updateResponse.SummarizeLogFileRotationListenerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CopyLogFileRotationListenerResponse != nil {
			readCopyLogFileRotationListenerResponse(ctx, updateResponse.CopyLogFileRotationListenerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyLogFileRotationListenerResponse != nil {
			readThirdPartyLogFileRotationListenerResponse(ctx, updateResponse.ThirdPartyLogFileRotationListenerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *logFileRotationListenerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLogFileRotationListener(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultLogFileRotationListenerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLogFileRotationListener(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readLogFileRotationListener(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state logFileRotationListenerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogFileRotationListenerAPI.GetLogFileRotationListener(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Log File Rotation Listener", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Log File Rotation Listener", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.UploadToS3LogFileRotationListenerResponse != nil {
		readUploadToS3LogFileRotationListenerResponse(ctx, readResponse.UploadToS3LogFileRotationListenerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SummarizeLogFileRotationListenerResponse != nil {
		readSummarizeLogFileRotationListenerResponse(ctx, readResponse.SummarizeLogFileRotationListenerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CopyLogFileRotationListenerResponse != nil {
		readCopyLogFileRotationListenerResponse(ctx, readResponse.CopyLogFileRotationListenerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyLogFileRotationListenerResponse != nil {
		readThirdPartyLogFileRotationListenerResponse(ctx, readResponse.ThirdPartyLogFileRotationListenerResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *logFileRotationListenerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLogFileRotationListener(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLogFileRotationListenerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLogFileRotationListener(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateLogFileRotationListener(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan logFileRotationListenerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state logFileRotationListenerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogFileRotationListenerAPI.UpdateLogFileRotationListener(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createLogFileRotationListenerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogFileRotationListenerAPI.UpdateLogFileRotationListenerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Log File Rotation Listener", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.UploadToS3LogFileRotationListenerResponse != nil {
			readUploadToS3LogFileRotationListenerResponse(ctx, updateResponse.UploadToS3LogFileRotationListenerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SummarizeLogFileRotationListenerResponse != nil {
			readSummarizeLogFileRotationListenerResponse(ctx, updateResponse.SummarizeLogFileRotationListenerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CopyLogFileRotationListenerResponse != nil {
			readCopyLogFileRotationListenerResponse(ctx, updateResponse.CopyLogFileRotationListenerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyLogFileRotationListenerResponse != nil {
			readThirdPartyLogFileRotationListenerResponse(ctx, updateResponse.ThirdPartyLogFileRotationListenerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultLogFileRotationListenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *logFileRotationListenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state logFileRotationListenerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogFileRotationListenerAPI.DeleteLogFileRotationListenerExecute(r.apiClient.LogFileRotationListenerAPI.DeleteLogFileRotationListener(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Log File Rotation Listener", err, httpResp)
		return
	}
}

func (r *logFileRotationListenerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLogFileRotationListener(ctx, req, resp)
}

func (r *defaultLogFileRotationListenerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLogFileRotationListener(ctx, req, resp)
}

func importLogFileRotationListener(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
