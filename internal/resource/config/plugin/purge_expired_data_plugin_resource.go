package plugin

import (
	"context"
	"time"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &purgeExpiredDataPluginResource{}
	_ resource.ResourceWithConfigure   = &purgeExpiredDataPluginResource{}
	_ resource.ResourceWithImportState = &purgeExpiredDataPluginResource{}
)

// Create a Purge Expired Data Plugin resource
func NewPurgeExpiredDataPluginResource() resource.Resource {
	return &purgeExpiredDataPluginResource{}
}

// purgeExpiredDataPluginResource is the resource implementation.
type purgeExpiredDataPluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *purgeExpiredDataPluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_purge_expired_data_plugin"
}

// Configure adds the provider configured client to the resource.
func (r *purgeExpiredDataPluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type purgeExpiredDataPluginResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	LastUpdated             types.String `tfsdk:"last_updated"`
	Notifications           types.Set    `tfsdk:"notifications"`
	RequiredActions         types.Set    `tfsdk:"required_actions"`
	DatetimeAttribute       types.String `tfsdk:"datetime_attribute"`
	DatetimeJSONField       types.String `tfsdk:"datetime_json_field"`
	DatetimeFormat          types.String `tfsdk:"datetime_format"`
	CustomDatetimeFormat    types.String `tfsdk:"custom_datetime_format"`
	CustomTimezone          types.String `tfsdk:"custom_timezone"`
	ExpirationOffset        types.String `tfsdk:"expiration_offset"`
	PurgeBehavior           types.String `tfsdk:"purge_behavior"`
	BaseDN                  types.String `tfsdk:"base_dn"`
	Filter                  types.String `tfsdk:"filter"`
	PollingInterval         types.String `tfsdk:"polling_interval"`
	MaxUpdatesPerSecond     types.Int64  `tfsdk:"max_updates_per_second"`
	PeerServerPriorityIndex types.Int64  `tfsdk:"peer_server_priority_index"`
	NumDeleteThreads        types.Int64  `tfsdk:"num_delete_threads"`
	Description             types.String `tfsdk:"description"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *purgeExpiredDataPluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Purge Expired Data Plugin.",
		Attributes: map[string]schema.Attribute{
			"datetime_attribute": schema.StringAttribute{
				Description: "The LDAP attribute that determines when data should be deleted. This could store the expiration time, or it could store the creation time and the expiration-offset property specifies the duration before data is deleted.",
				Required:    true,
			},
			"datetime_json_field": schema.StringAttribute{
				Description: "The top-level JSON field within the configured datetime-attribute that determines when data should be deleted. This could store the expiration time, or it could store the creation time and the expiration-offset property specifies the duration before data is deleted.",
				Optional:    true,
			},
			"datetime_format": schema.StringAttribute{
				Description: "Specifies the format of the datetime stored within the entry that determines when data should be purged.",
				Required:    true,
			},
			"custom_datetime_format": schema.StringAttribute{
				Description: "When the datetime-format property is configured with a value of \"custom\", this specifies the format (using a string compatible with the java.text.SimpleDateFormat class) that will be used to search for expired data.",
				Optional:    true,
			},
			"custom_timezone": schema.StringAttribute{
				Description: "Specifies the time zone to use when generating a date string using the configured custom-datetime-format value. The provided value must be accepted by java.util.TimeZone.getTimeZone.",
				Optional:    true,
			},
			"expiration_offset": schema.StringAttribute{
				Description: "The duration to wait after the value specified in datetime-attribute (and optionally datetime-json-field) before purging the data.",
				Required:    true,
			},
			"purge_behavior": schema.StringAttribute{
				Description: "Specifies whether to delete expired entries or attribute values. By default entries are deleted.",
				Optional:    true,
			},
			"base_dn": schema.StringAttribute{
				Description: "Only entries located within the subtree specified by this base DN are eligible for purging.",
				Optional:    true,
			},
			"filter": schema.StringAttribute{
				Description: "Only entries that match this LDAP filter will be eligible for having data purged.",
				Optional:    true,
			},
			"polling_interval": schema.StringAttribute{
				Description: "This specifies how often the plugin should check for expired data. It also controls the offset of peer servers (see the peer-server-priority-index for more information).",
				Required:    true,
			},
			"max_updates_per_second": schema.Int64Attribute{
				Description: "This setting smooths out the performance impact on the server by throttling the purging to the specified maximum number of updates per second. To avoid a large backlog, this value should be set comfortably above the average rate that expired data is generated. When purge-behavior is set to subtree-delete-entries, then deletion of the entire subtree is considered a single update for the purposes of throttling.",
				Required:    true,
			},
			"peer_server_priority_index": schema.Int64Attribute{
				Description: "In a replicated environment, this determines the order in which peer servers should attempt to purge data.",
				Optional:    true,
			},
			"num_delete_threads": schema.Int64Attribute{
				Description: "The number of threads used to delete expired entries.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Plugin",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the plug-in is enabled for use.",
				Required:    true,
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalPurgeExpiredDataPluginFields(ctx context.Context, addRequest *client.AddPurgeExpiredDataPluginRequest, plan purgeExpiredDataPluginResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DatetimeJSONField) {
		stringVal := plan.DatetimeJSONField.ValueString()
		addRequest.DatetimeJSONField = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CustomDatetimeFormat) {
		stringVal := plan.CustomDatetimeFormat.ValueString()
		addRequest.CustomDatetimeFormat = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CustomTimezone) {
		stringVal := plan.CustomTimezone.ValueString()
		addRequest.CustomTimezone = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PurgeBehavior) {
		purgeBehavior, err := client.NewEnumpluginPurgeBehaviorPropFromValue(plan.PurgeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.PurgeBehavior = purgeBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BaseDN) {
		stringVal := plan.BaseDN.ValueString()
		addRequest.BaseDN = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Filter) {
		stringVal := plan.Filter.ValueString()
		addRequest.Filter = &stringVal
	}
	if internaltypes.IsDefined(plan.PeerServerPriorityIndex) {
		intVal := int32(plan.PeerServerPriorityIndex.ValueInt64())
		addRequest.PeerServerPriorityIndex = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	return nil
}

// Read a PurgeExpiredDataPluginResponse object into the model struct
func readPurgeExpiredDataPluginResponse(ctx context.Context, r *client.PurgeExpiredDataPluginResponse, state *purgeExpiredDataPluginResourceModel, expectedValues *purgeExpiredDataPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.DatetimeAttribute = types.StringValue(r.DatetimeAttribute)
	state.DatetimeJSONField = internaltypes.StringTypeOrNil(r.DatetimeJSONField, internaltypes.IsEmptyString(expectedValues.DatetimeJSONField))
	state.DatetimeFormat = types.StringValue(r.DatetimeFormat.String())
	state.CustomDatetimeFormat = internaltypes.StringTypeOrNil(r.CustomDatetimeFormat, internaltypes.IsEmptyString(expectedValues.CustomDatetimeFormat))
	state.CustomTimezone = internaltypes.StringTypeOrNil(r.CustomTimezone, internaltypes.IsEmptyString(expectedValues.CustomTimezone))
	state.ExpirationOffset = types.StringValue(r.ExpirationOffset)
	config.CheckMismatchedPDFormattedAttributes("expiration_offset",
		expectedValues.ExpirationOffset, state.ExpirationOffset, diagnostics)
	state.PurgeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginPurgeBehaviorProp(r.PurgeBehavior), internaltypes.IsEmptyString(expectedValues.PurgeBehavior))
	state.BaseDN = internaltypes.StringTypeOrNil(r.BaseDN, internaltypes.IsEmptyString(expectedValues.BaseDN))
	state.Filter = internaltypes.StringTypeOrNil(r.Filter, internaltypes.IsEmptyString(expectedValues.Filter))
	state.PollingInterval = types.StringValue(r.PollingInterval)
	config.CheckMismatchedPDFormattedAttributes("polling_interval",
		expectedValues.PollingInterval, state.PollingInterval, diagnostics)
	state.MaxUpdatesPerSecond = types.Int64Value(int64(r.MaxUpdatesPerSecond))
	state.PeerServerPriorityIndex = internaltypes.Int64TypeOrNil(r.PeerServerPriorityIndex)
	state.NumDeleteThreads = types.Int64Value(int64(r.NumDeleteThreads))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createPurgeExpiredDataPluginOperations(plan purgeExpiredDataPluginResourceModel, state purgeExpiredDataPluginResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.DatetimeAttribute, state.DatetimeAttribute, "datetime-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.DatetimeJSONField, state.DatetimeJSONField, "datetime-json-field")
	operations.AddStringOperationIfNecessary(&ops, plan.DatetimeFormat, state.DatetimeFormat, "datetime-format")
	operations.AddStringOperationIfNecessary(&ops, plan.CustomDatetimeFormat, state.CustomDatetimeFormat, "custom-datetime-format")
	operations.AddStringOperationIfNecessary(&ops, plan.CustomTimezone, state.CustomTimezone, "custom-timezone")
	operations.AddStringOperationIfNecessary(&ops, plan.ExpirationOffset, state.ExpirationOffset, "expiration-offset")
	operations.AddStringOperationIfNecessary(&ops, plan.PurgeBehavior, state.PurgeBehavior, "purge-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.Filter, state.Filter, "filter")
	operations.AddStringOperationIfNecessary(&ops, plan.PollingInterval, state.PollingInterval, "polling-interval")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxUpdatesPerSecond, state.MaxUpdatesPerSecond, "max-updates-per-second")
	operations.AddInt64OperationIfNecessary(&ops, plan.PeerServerPriorityIndex, state.PeerServerPriorityIndex, "peer-server-priority-index")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumDeleteThreads, state.NumDeleteThreads, "num-delete-threads")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *purgeExpiredDataPluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan purgeExpiredDataPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	datetimeFormat, err := client.NewEnumpluginDatetimeFormatPropFromValue(plan.DatetimeFormat.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for DatetimeFormat", err.Error())
		return
	}
	addRequest := client.NewAddPurgeExpiredDataPluginRequest(plan.Id.ValueString(),
		[]client.EnumpurgeExpiredDataPluginSchemaUrn{client.ENUMPURGEEXPIREDDATAPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINPURGE_EXPIRED_DATA},
		plan.DatetimeAttribute.ValueString(),
		*datetimeFormat,
		plan.ExpirationOffset.ValueString(),
		plan.PollingInterval.ValueString(),
		int32(plan.MaxUpdatesPerSecond.ValueInt64()),
		int32(plan.NumDeleteThreads.ValueInt64()),
		plan.Enabled.ValueBool())
	err = addOptionalPurgeExpiredDataPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Purge Expired Data Plugin", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddPurgeExpiredDataPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Purge Expired Data Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state purgeExpiredDataPluginResourceModel
	readPurgeExpiredDataPluginResponse(ctx, addResponse.PurgeExpiredDataPluginResponse, &state, &plan, &resp.Diagnostics)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *purgeExpiredDataPluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state purgeExpiredDataPluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Purge Expired Data Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readPurgeExpiredDataPluginResponse(ctx, readResponse.PurgeExpiredDataPluginResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *purgeExpiredDataPluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan purgeExpiredDataPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state purgeExpiredDataPluginResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.PluginApi.UpdatePlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createPurgeExpiredDataPluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Purge Expired Data Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readPurgeExpiredDataPluginResponse(ctx, updateResponse.PurgeExpiredDataPluginResponse, &state, &plan, &resp.Diagnostics)
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
func (r *purgeExpiredDataPluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state purgeExpiredDataPluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PluginApi.DeletePluginExecute(r.apiClient.PluginApi.DeletePlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Purge Expired Data Plugin", err, httpResp)
		return
	}
}

func (r *purgeExpiredDataPluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
