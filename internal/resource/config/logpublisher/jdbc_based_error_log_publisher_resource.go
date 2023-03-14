package logpublisher

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
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
	_ resource.Resource                = &jdbcBasedErrorLogPublisherResource{}
	_ resource.ResourceWithConfigure   = &jdbcBasedErrorLogPublisherResource{}
	_ resource.ResourceWithImportState = &jdbcBasedErrorLogPublisherResource{}
	_ resource.Resource                = &defaultJdbcBasedErrorLogPublisherResource{}
	_ resource.ResourceWithConfigure   = &defaultJdbcBasedErrorLogPublisherResource{}
	_ resource.ResourceWithImportState = &defaultJdbcBasedErrorLogPublisherResource{}
)

// Create a Jdbc Based Error Log Publisher resource
func NewJdbcBasedErrorLogPublisherResource() resource.Resource {
	return &jdbcBasedErrorLogPublisherResource{}
}

func NewDefaultJdbcBasedErrorLogPublisherResource() resource.Resource {
	return &defaultJdbcBasedErrorLogPublisherResource{}
}

// jdbcBasedErrorLogPublisherResource is the resource implementation.
type jdbcBasedErrorLogPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultJdbcBasedErrorLogPublisherResource is the resource implementation.
type defaultJdbcBasedErrorLogPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *jdbcBasedErrorLogPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jdbc_based_error_log_publisher"
}

func (r *defaultJdbcBasedErrorLogPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_jdbc_based_error_log_publisher"
}

// Configure adds the provider configured client to the resource.
func (r *jdbcBasedErrorLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultJdbcBasedErrorLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type jdbcBasedErrorLogPublisherResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	LastUpdated          types.String `tfsdk:"last_updated"`
	Notifications        types.Set    `tfsdk:"notifications"`
	RequiredActions      types.Set    `tfsdk:"required_actions"`
	Server               types.String `tfsdk:"server"`
	LogFieldMapping      types.String `tfsdk:"log_field_mapping"`
	LogTableName         types.String `tfsdk:"log_table_name"`
	QueueSize            types.Int64  `tfsdk:"queue_size"`
	DefaultSeverity      types.Set    `tfsdk:"default_severity"`
	OverrideSeverity     types.Set    `tfsdk:"override_severity"`
	Description          types.String `tfsdk:"description"`
	Enabled              types.Bool   `tfsdk:"enabled"`
	LoggingErrorBehavior types.String `tfsdk:"logging_error_behavior"`
}

// GetSchema defines the schema for the resource.
func (r *jdbcBasedErrorLogPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	jdbcBasedErrorLogPublisherSchema(ctx, req, resp, false)
}

func (r *defaultJdbcBasedErrorLogPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	jdbcBasedErrorLogPublisherSchema(ctx, req, resp, true)
}

func jdbcBasedErrorLogPublisherSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Jdbc Based Error Log Publisher.",
		Attributes: map[string]schema.Attribute{
			"server": schema.StringAttribute{
				Description: "The JDBC-based Database Server to use for a connection.",
				Required:    true,
			},
			"log_field_mapping": schema.StringAttribute{
				Description: "The log field mapping associates loggable fields to database column names. The table name is not part of this mapping.",
				Required:    true,
			},
			"log_table_name": schema.StringAttribute{
				Description: "The table name to log entries to the database server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"queue_size": schema.Int64Attribute{
				Description: "The maximum number of log records that can be stored in the asynchronous queue.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"default_severity": schema.SetAttribute{
				Description: "Specifies the default severity levels for the logger.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"override_severity": schema.SetAttribute{
				Description: "Specifies the override severity levels for the logger based on the category of the messages.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log Publisher",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Log Publisher is enabled for use.",
				Required:    true,
			},
			"logging_error_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the server should exhibit if an error occurs during logging processing.",
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
func addOptionalJdbcBasedErrorLogPublisherFields(ctx context.Context, addRequest *client.AddJdbcBasedErrorLogPublisherRequest, plan jdbcBasedErrorLogPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogTableName) {
		stringVal := plan.LogTableName.ValueString()
		addRequest.LogTableName = &stringVal
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		intVal := int32(plan.QueueSize.ValueInt64())
		addRequest.QueueSize = &intVal
	}
	if internaltypes.IsDefined(plan.DefaultSeverity) {
		var slice []string
		plan.DefaultSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherDefaultSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherDefaultSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefaultSeverity = enumSlice
	}
	if internaltypes.IsDefined(plan.OverrideSeverity) {
		var slice []string
		plan.OverrideSeverity.ElementsAs(ctx, &slice, false)
		addRequest.OverrideSeverity = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Read a JdbcBasedErrorLogPublisherResponse object into the model struct
func readJdbcBasedErrorLogPublisherResponse(ctx context.Context, r *client.JdbcBasedErrorLogPublisherResponse, state *jdbcBasedErrorLogPublisherResourceModel, expectedValues *jdbcBasedErrorLogPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Server = types.StringValue(r.Server)
	state.LogFieldMapping = types.StringValue(r.LogFieldMapping)
	state.LogTableName = types.StringValue(r.LogTableName)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createJdbcBasedErrorLogPublisherOperations(plan jdbcBasedErrorLogPublisherResourceModel, state jdbcBasedErrorLogPublisherResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Server, state.Server, "server")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldMapping, state.LogFieldMapping, "log-field-mapping")
	operations.AddStringOperationIfNecessary(&ops, plan.LogTableName, state.LogTableName, "log-table-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.QueueSize, state.QueueSize, "queue-size")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DefaultSeverity, state.DefaultSeverity, "default-severity")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.OverrideSeverity, state.OverrideSeverity, "override-severity")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.LoggingErrorBehavior, state.LoggingErrorBehavior, "logging-error-behavior")
	return ops
}

// Create a new resource
func (r *jdbcBasedErrorLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan jdbcBasedErrorLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddJdbcBasedErrorLogPublisherRequest(plan.Id.ValueString(),
		[]client.EnumjdbcBasedErrorLogPublisherSchemaUrn{client.ENUMJDBCBASEDERRORLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERJDBC_BASED_ERROR},
		plan.Server.ValueString(),
		plan.LogFieldMapping.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalJdbcBasedErrorLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Jdbc Based Error Log Publisher", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddJdbcBasedErrorLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Jdbc Based Error Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state jdbcBasedErrorLogPublisherResourceModel
	readJdbcBasedErrorLogPublisherResponse(ctx, addResponse.JdbcBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultJdbcBasedErrorLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan jdbcBasedErrorLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Jdbc Based Error Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state jdbcBasedErrorLogPublisherResourceModel
	readJdbcBasedErrorLogPublisherResponse(ctx, readResponse.JdbcBasedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogPublisherApi.UpdateLogPublisher(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createJdbcBasedErrorLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Jdbc Based Error Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readJdbcBasedErrorLogPublisherResponse(ctx, updateResponse.JdbcBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
func (r *jdbcBasedErrorLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readJdbcBasedErrorLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultJdbcBasedErrorLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readJdbcBasedErrorLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readJdbcBasedErrorLogPublisher(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state jdbcBasedErrorLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Jdbc Based Error Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readJdbcBasedErrorLogPublisherResponse(ctx, readResponse.JdbcBasedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *jdbcBasedErrorLogPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateJdbcBasedErrorLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultJdbcBasedErrorLogPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateJdbcBasedErrorLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateJdbcBasedErrorLogPublisher(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan jdbcBasedErrorLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state jdbcBasedErrorLogPublisherResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogPublisherApi.UpdateLogPublisher(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createJdbcBasedErrorLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Jdbc Based Error Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readJdbcBasedErrorLogPublisherResponse(ctx, updateResponse.JdbcBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultJdbcBasedErrorLogPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *jdbcBasedErrorLogPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state jdbcBasedErrorLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogPublisherApi.DeleteLogPublisherExecute(r.apiClient.LogPublisherApi.DeleteLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Jdbc Based Error Log Publisher", err, httpResp)
		return
	}
}

func (r *jdbcBasedErrorLogPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importJdbcBasedErrorLogPublisher(ctx, req, resp)
}

func (r *defaultJdbcBasedErrorLogPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importJdbcBasedErrorLogPublisher(ctx, req, resp)
}

func importJdbcBasedErrorLogPublisher(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
