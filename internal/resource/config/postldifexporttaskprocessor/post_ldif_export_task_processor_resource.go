package postldifexporttaskprocessor

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
	client "github.com/pingidentity/pingdirectory-go-client/v10000/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &postLdifExportTaskProcessorResource{}
	_ resource.ResourceWithConfigure   = &postLdifExportTaskProcessorResource{}
	_ resource.ResourceWithImportState = &postLdifExportTaskProcessorResource{}
	_ resource.Resource                = &defaultPostLdifExportTaskProcessorResource{}
	_ resource.ResourceWithConfigure   = &defaultPostLdifExportTaskProcessorResource{}
	_ resource.ResourceWithImportState = &defaultPostLdifExportTaskProcessorResource{}
)

// Create a Post Ldif Export Task Processor resource
func NewPostLdifExportTaskProcessorResource() resource.Resource {
	return &postLdifExportTaskProcessorResource{}
}

func NewDefaultPostLdifExportTaskProcessorResource() resource.Resource {
	return &defaultPostLdifExportTaskProcessorResource{}
}

// postLdifExportTaskProcessorResource is the resource implementation.
type postLdifExportTaskProcessorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultPostLdifExportTaskProcessorResource is the resource implementation.
type defaultPostLdifExportTaskProcessorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *postLdifExportTaskProcessorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_post_ldif_export_task_processor"
}

func (r *defaultPostLdifExportTaskProcessorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_post_ldif_export_task_processor"
}

// Configure adds the provider configured client to the resource.
func (r *postLdifExportTaskProcessorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultPostLdifExportTaskProcessorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type postLdifExportTaskProcessorResourceModel struct {
	Id                                   types.String `tfsdk:"id"`
	Name                                 types.String `tfsdk:"name"`
	Notifications                        types.Set    `tfsdk:"notifications"`
	RequiredActions                      types.Set    `tfsdk:"required_actions"`
	Type                                 types.String `tfsdk:"type"`
	ExtensionClass                       types.String `tfsdk:"extension_class"`
	ExtensionArgument                    types.Set    `tfsdk:"extension_argument"`
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
func (r *postLdifExportTaskProcessorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	postLdifExportTaskProcessorSchema(ctx, req, resp, false)
}

func (r *defaultPostLdifExportTaskProcessorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	postLdifExportTaskProcessorSchema(ctx, req, resp, true)
}

func postLdifExportTaskProcessorSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Post Ldif Export Task Processor. Supported in PingDirectory product version 10.0.0.0+.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Post LDIF Export Task Processor resource. Options are ['upload-to-s3', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"upload-to-s3", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Post LDIF Export Task Processor.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Post LDIF Export Task Processor. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"aws_external_server": schema.StringAttribute{
				Description: "The external server with information to use when interacting with the AWS S3 service.",
				Optional:    true,
			},
			"s3_bucket_name": schema.StringAttribute{
				Description: "The name of the S3 bucket into which LDIF files should be copied.",
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
				Description: "The maximum number of existing files matching the file retention pattern that should be retained in the S3 bucket after successfully uploading a newly exported file.",
				Optional:    true,
			},
			"maximum_file_age_to_retain": schema.StringAttribute{
				Description: "The maximum length of time to retain files matching the file retention pattern that should be retained in the S3 bucket after successfully uploading a newly exported file.",
				Optional:    true,
			},
			"file_retention_pattern": schema.StringAttribute{
				Description: "A regular expression pattern that will be used to identify which files are candidates for automatic removal based on the maximum-file-count-to-retain and maximum-file-age-to-retain properties. By default, all files in the bucket will be eligible for removal by retention processing.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Post LDIF Export Task Processor",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Post LDIF Export Task Processor is enabled for use.",
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
func (r *postLdifExportTaskProcessorResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanPostLdifExportTaskProcessor(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_post_ldif_export_task_processor")
	var planModel, configModel postLdifExportTaskProcessorResourceModel
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
	if anyDefaultsSet {
		planModel.Notifications = types.SetUnknown(types.StringType)
		planModel.RequiredActions = types.SetUnknown(config.GetRequiredActionsObjectType())
	}
	planModel.setNotApplicableAttrsNull()
	resp.Plan.Set(ctx, &planModel)
}

func (r *defaultPostLdifExportTaskProcessorResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanPostLdifExportTaskProcessor(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_default_post_ldif_export_task_processor")
}

func modifyPlanPostLdifExportTaskProcessor(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, resourceName string) {
	version.CheckResourceSupported(&resp.Diagnostics, version.PingDirectory10000,
		providerConfig.ProductVersion, resourceName)
}

func (model *postLdifExportTaskProcessorResourceModel) setNotApplicableAttrsNull() {
	resourceType := model.Type.ValueString()
	// Set any not applicable computed attributes to null for each type
	if resourceType == "third-party" {
		model.FileRetentionPattern = types.StringNull()
	}
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsPostLdifExportTaskProcessor() []resource.ConfigValidator {
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
			"third-party",
			[]path.Expression{path.MatchRoot("extension_class")},
		),
	}
}

// Add config validators
func (r postLdifExportTaskProcessorResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsPostLdifExportTaskProcessor()
}

// Add config validators
func (r defaultPostLdifExportTaskProcessorResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsPostLdifExportTaskProcessor()
}

// Add optional fields to create request for upload-to-s3 post-ldif-export-task-processor
func addOptionalUploadToS3PostLdifExportTaskProcessorFields(ctx context.Context, addRequest *client.AddUploadToS3PostLdifExportTaskProcessorRequest, plan postLdifExportTaskProcessorResourceModel) {
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

// Add optional fields to create request for third-party post-ldif-export-task-processor
func addOptionalThirdPartyPostLdifExportTaskProcessorFields(ctx context.Context, addRequest *client.AddThirdPartyPostLdifExportTaskProcessorRequest, plan postLdifExportTaskProcessorResourceModel) {
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
func populatePostLdifExportTaskProcessorUnknownValues(model *postLdifExportTaskProcessorResourceModel) {
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *postLdifExportTaskProcessorResourceModel) populateAllComputedStringAttributes() {
	if model.MaximumFileAgeToRetain.IsUnknown() || model.MaximumFileAgeToRetain.IsNull() {
		model.MaximumFileAgeToRetain = types.StringValue("")
	}
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
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
	if model.FileRetentionPattern.IsUnknown() || model.FileRetentionPattern.IsNull() {
		model.FileRetentionPattern = types.StringValue("")
	}
}

// Read a UploadToS3PostLdifExportTaskProcessorResponse object into the model struct
func readUploadToS3PostLdifExportTaskProcessorResponse(ctx context.Context, r *client.UploadToS3PostLdifExportTaskProcessorResponse, state *postLdifExportTaskProcessorResourceModel, expectedValues *postLdifExportTaskProcessorResourceModel, diagnostics *diag.Diagnostics) {
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
	populatePostLdifExportTaskProcessorUnknownValues(state)
}

// Read a ThirdPartyPostLdifExportTaskProcessorResponse object into the model struct
func readThirdPartyPostLdifExportTaskProcessorResponse(ctx context.Context, r *client.ThirdPartyPostLdifExportTaskProcessorResponse, state *postLdifExportTaskProcessorResourceModel, expectedValues *postLdifExportTaskProcessorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePostLdifExportTaskProcessorUnknownValues(state)
}

// Create any update operations necessary to make the state match the plan
func createPostLdifExportTaskProcessorOperations(plan postLdifExportTaskProcessorResourceModel, state postLdifExportTaskProcessorResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
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

// Create a upload-to-s3 post-ldif-export-task-processor
func (r *postLdifExportTaskProcessorResource) CreateUploadToS3PostLdifExportTaskProcessor(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan postLdifExportTaskProcessorResourceModel) (*postLdifExportTaskProcessorResourceModel, error) {
	addRequest := client.NewAddUploadToS3PostLdifExportTaskProcessorRequest([]client.EnumuploadToS3PostLdifExportTaskProcessorSchemaUrn{client.ENUMUPLOADTOS3POSTLDIFEXPORTTASKPROCESSORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0POST_LDIF_EXPORT_TASK_PROCESSORUPLOAD_TO_S3},
		plan.AwsExternalServer.ValueString(),
		plan.S3BucketName.ValueString(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	addOptionalUploadToS3PostLdifExportTaskProcessorFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PostLdifExportTaskProcessorAPI.AddPostLdifExportTaskProcessor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPostLdifExportTaskProcessorRequest(
		client.AddUploadToS3PostLdifExportTaskProcessorRequestAsAddPostLdifExportTaskProcessorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PostLdifExportTaskProcessorAPI.AddPostLdifExportTaskProcessorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Post Ldif Export Task Processor", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state postLdifExportTaskProcessorResourceModel
	readUploadToS3PostLdifExportTaskProcessorResponse(ctx, addResponse.UploadToS3PostLdifExportTaskProcessorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party post-ldif-export-task-processor
func (r *postLdifExportTaskProcessorResource) CreateThirdPartyPostLdifExportTaskProcessor(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan postLdifExportTaskProcessorResourceModel) (*postLdifExportTaskProcessorResourceModel, error) {
	addRequest := client.NewAddThirdPartyPostLdifExportTaskProcessorRequest([]client.EnumthirdPartyPostLdifExportTaskProcessorSchemaUrn{client.ENUMTHIRDPARTYPOSTLDIFEXPORTTASKPROCESSORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0POST_LDIF_EXPORT_TASK_PROCESSORTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	addOptionalThirdPartyPostLdifExportTaskProcessorFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PostLdifExportTaskProcessorAPI.AddPostLdifExportTaskProcessor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPostLdifExportTaskProcessorRequest(
		client.AddThirdPartyPostLdifExportTaskProcessorRequestAsAddPostLdifExportTaskProcessorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PostLdifExportTaskProcessorAPI.AddPostLdifExportTaskProcessorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Post Ldif Export Task Processor", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state postLdifExportTaskProcessorResourceModel
	readThirdPartyPostLdifExportTaskProcessorResponse(ctx, addResponse.ThirdPartyPostLdifExportTaskProcessorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *postLdifExportTaskProcessorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan postLdifExportTaskProcessorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *postLdifExportTaskProcessorResourceModel
	var err error
	if plan.Type.ValueString() == "upload-to-s3" {
		state, err = r.CreateUploadToS3PostLdifExportTaskProcessor(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyPostLdifExportTaskProcessor(ctx, req, resp, plan)
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
func (r *defaultPostLdifExportTaskProcessorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan postLdifExportTaskProcessorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PostLdifExportTaskProcessorAPI.GetPostLdifExportTaskProcessor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Post Ldif Export Task Processor", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state postLdifExportTaskProcessorResourceModel
	if readResponse.UploadToS3PostLdifExportTaskProcessorResponse != nil {
		readUploadToS3PostLdifExportTaskProcessorResponse(ctx, readResponse.UploadToS3PostLdifExportTaskProcessorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyPostLdifExportTaskProcessorResponse != nil {
		readThirdPartyPostLdifExportTaskProcessorResponse(ctx, readResponse.ThirdPartyPostLdifExportTaskProcessorResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.PostLdifExportTaskProcessorAPI.UpdatePostLdifExportTaskProcessor(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createPostLdifExportTaskProcessorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PostLdifExportTaskProcessorAPI.UpdatePostLdifExportTaskProcessorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Post Ldif Export Task Processor", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.UploadToS3PostLdifExportTaskProcessorResponse != nil {
			readUploadToS3PostLdifExportTaskProcessorResponse(ctx, updateResponse.UploadToS3PostLdifExportTaskProcessorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyPostLdifExportTaskProcessorResponse != nil {
			readThirdPartyPostLdifExportTaskProcessorResponse(ctx, updateResponse.ThirdPartyPostLdifExportTaskProcessorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *postLdifExportTaskProcessorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPostLdifExportTaskProcessor(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultPostLdifExportTaskProcessorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPostLdifExportTaskProcessor(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readPostLdifExportTaskProcessor(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state postLdifExportTaskProcessorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.PostLdifExportTaskProcessorAPI.GetPostLdifExportTaskProcessor(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Post Ldif Export Task Processor", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Post Ldif Export Task Processor", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.UploadToS3PostLdifExportTaskProcessorResponse != nil {
		readUploadToS3PostLdifExportTaskProcessorResponse(ctx, readResponse.UploadToS3PostLdifExportTaskProcessorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyPostLdifExportTaskProcessorResponse != nil {
		readThirdPartyPostLdifExportTaskProcessorResponse(ctx, readResponse.ThirdPartyPostLdifExportTaskProcessorResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *postLdifExportTaskProcessorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePostLdifExportTaskProcessor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPostLdifExportTaskProcessorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePostLdifExportTaskProcessor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updatePostLdifExportTaskProcessor(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan postLdifExportTaskProcessorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state postLdifExportTaskProcessorResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.PostLdifExportTaskProcessorAPI.UpdatePostLdifExportTaskProcessor(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createPostLdifExportTaskProcessorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.PostLdifExportTaskProcessorAPI.UpdatePostLdifExportTaskProcessorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Post Ldif Export Task Processor", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.UploadToS3PostLdifExportTaskProcessorResponse != nil {
			readUploadToS3PostLdifExportTaskProcessorResponse(ctx, updateResponse.UploadToS3PostLdifExportTaskProcessorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyPostLdifExportTaskProcessorResponse != nil {
			readThirdPartyPostLdifExportTaskProcessorResponse(ctx, updateResponse.ThirdPartyPostLdifExportTaskProcessorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultPostLdifExportTaskProcessorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *postLdifExportTaskProcessorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state postLdifExportTaskProcessorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PostLdifExportTaskProcessorAPI.DeletePostLdifExportTaskProcessorExecute(r.apiClient.PostLdifExportTaskProcessorAPI.DeletePostLdifExportTaskProcessor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Post Ldif Export Task Processor", err, httpResp)
		return
	}
}

func (r *postLdifExportTaskProcessorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPostLdifExportTaskProcessor(ctx, req, resp)
}

func (r *defaultPostLdifExportTaskProcessorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPostLdifExportTaskProcessor(ctx, req, resp)
}

func importPostLdifExportTaskProcessor(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
