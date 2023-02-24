package backend

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &configFileHandlerBackendResource{}
	_ resource.ResourceWithConfigure   = &configFileHandlerBackendResource{}
	_ resource.ResourceWithImportState = &configFileHandlerBackendResource{}
)

// Create a Config File Handler Backend resource
func NewConfigFileHandlerBackendResource() resource.Resource {
	return &configFileHandlerBackendResource{}
}

// configFileHandlerBackendResource is the resource implementation.
type configFileHandlerBackendResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *configFileHandlerBackendResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config_file_handler_backend"
}

// Configure adds the provider configured client to the resource.
func (r *configFileHandlerBackendResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type configFileHandlerBackendResourceModel struct {
	Id                                  types.String `tfsdk:"id"`
	LastUpdated                         types.String `tfsdk:"last_updated"`
	Notifications                       types.Set    `tfsdk:"notifications"`
	RequiredActions                     types.Set    `tfsdk:"required_actions"`
	BackendID                           types.String `tfsdk:"backend_id"`
	BaseDN                              types.Set    `tfsdk:"base_dn"`
	WritabilityMode                     types.String `tfsdk:"writability_mode"`
	InsignificantConfigArchiveAttribute types.Set    `tfsdk:"insignificant_config_archive_attribute"`
	MirroredSubtreePeerPollingInterval  types.String `tfsdk:"mirrored_subtree_peer_polling_interval"`
	MirroredSubtreeEntryUpdateTimeout   types.String `tfsdk:"mirrored_subtree_entry_update_timeout"`
	MirroredSubtreeSearchTimeout        types.String `tfsdk:"mirrored_subtree_search_timeout"`
	Description                         types.String `tfsdk:"description"`
	Enabled                             types.Bool   `tfsdk:"enabled"`
	SetDegradedAlertWhenDisabled        types.Bool   `tfsdk:"set_degraded_alert_when_disabled"`
	ReturnUnavailableWhenDisabled       types.Bool   `tfsdk:"return_unavailable_when_disabled"`
	BackupFilePermissions               types.String `tfsdk:"backup_file_permissions"`
	NotificationManager                 types.String `tfsdk:"notification_manager"`
}

// GetSchema defines the schema for the resource.
func (r *configFileHandlerBackendResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Config File Handler Backend.",
		Attributes: map[string]schema.Attribute{
			"backend_id": schema.StringAttribute{
				Description: "Specifies a name to identify the associated backend.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"base_dn": schema.SetAttribute{
				Description: "Specifies the base DN(s) for the data that the backend handles.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"writability_mode": schema.StringAttribute{
				Description: "Specifies the behavior that the backend should use when processing write operations.",
				Optional:    true,
				Computed:    true,
			},
			"insignificant_config_archive_attribute": schema.SetAttribute{
				Description: "The name or OID of an attribute type that is considered insignificant for the purpose of maintaining the configuration archive.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"mirrored_subtree_peer_polling_interval": schema.StringAttribute{
				Description: "Tells the server component that is responsible for mirroring configuration data across a topology of servers the maximum amount of time to wait before polling the peer servers in the topology to determine if there are any changes in the topology. Mirrored data includes meta-data about the servers in the topology as well as cluster-wide configuration data.",
				Optional:    true,
				Computed:    true,
			},
			"mirrored_subtree_entry_update_timeout": schema.StringAttribute{
				Description: "Tells the server component that is responsible for mirroring configuration data across a topology of servers the maximum amount of time to wait for an update operation (add, delete, modify and modify-dn) on an entry to be applied on all servers in the topology. Mirrored data includes meta-data about the servers in the topology as well as cluster-wide configuration data.",
				Optional:    true,
				Computed:    true,
			},
			"mirrored_subtree_search_timeout": schema.StringAttribute{
				Description: "Tells the server component that is responsible for mirroring configuration data across a topology of servers the maximum amount of time to wait for a search operation to complete. Mirrored data includes meta-data about the servers in the topology as well as cluster-wide configuration data. Search requests that take longer than this timeout will be canceled and considered failures.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Backend",
				Optional:    true,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the backend is enabled in the server.",
				Optional:    true,
				Computed:    true,
			},
			"set_degraded_alert_when_disabled": schema.BoolAttribute{
				Description: "Determines whether the Directory Server enters a DEGRADED state (and sends a corresponding alert) when this Backend is disabled.",
				Optional:    true,
				Computed:    true,
			},
			"return_unavailable_when_disabled": schema.BoolAttribute{
				Description: "Determines whether any LDAP operation that would use this Backend is to return UNAVAILABLE when this Backend is disabled.",
				Optional:    true,
				Computed:    true,
			},
			"backup_file_permissions": schema.StringAttribute{
				Description: "Specifies the permissions that should be applied to files and directories created by a backup of the backend.",
				Optional:    true,
				Computed:    true,
			},
			"notification_manager": schema.StringAttribute{
				Description: "Specifies a notification manager for changes resulting from operations processed through this Backend",
				Optional:    true,
				Computed:    true,
			},
		},
	}
	config.AddCommonSchema(&schema, false)
	resp.Schema = schema
}

// Read a ConfigFileHandlerBackendResponse object into the model struct
func readConfigFileHandlerBackendResponse(ctx context.Context, r *client.ConfigFileHandlerBackendResponse, state *configFileHandlerBackendResourceModel, expectedValues *configFileHandlerBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.BackendID = types.StringValue(r.BackendID)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.WritabilityMode = types.StringValue(r.WritabilityMode.String())
	state.InsignificantConfigArchiveAttribute = internaltypes.GetStringSet(r.InsignificantConfigArchiveAttribute)
	state.MirroredSubtreePeerPollingInterval = internaltypes.StringTypeOrNil(r.MirroredSubtreePeerPollingInterval, true)
	config.CheckMismatchedPDFormattedAttributes("mirrored_subtree_peer_polling_interval",
		expectedValues.MirroredSubtreePeerPollingInterval, state.MirroredSubtreePeerPollingInterval, diagnostics)
	state.MirroredSubtreeEntryUpdateTimeout = internaltypes.StringTypeOrNil(r.MirroredSubtreeEntryUpdateTimeout, true)
	config.CheckMismatchedPDFormattedAttributes("mirrored_subtree_entry_update_timeout",
		expectedValues.MirroredSubtreeEntryUpdateTimeout, state.MirroredSubtreeEntryUpdateTimeout, diagnostics)
	state.MirroredSubtreeSearchTimeout = internaltypes.StringTypeOrNil(r.MirroredSubtreeSearchTimeout, true)
	config.CheckMismatchedPDFormattedAttributes("mirrored_subtree_search_timeout",
		expectedValues.MirroredSubtreeSearchTimeout, state.MirroredSubtreeSearchTimeout, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.BackupFilePermissions = internaltypes.StringTypeOrNil(r.BackupFilePermissions, true)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createConfigFileHandlerBackendOperations(plan configFileHandlerBackendResourceModel, state configFileHandlerBackendResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.BackendID, state.BackendID, "backend-id")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.WritabilityMode, state.WritabilityMode, "writability-mode")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.InsignificantConfigArchiveAttribute, state.InsignificantConfigArchiveAttribute, "insignificant-config-archive-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.MirroredSubtreePeerPollingInterval, state.MirroredSubtreePeerPollingInterval, "mirrored-subtree-peer-polling-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.MirroredSubtreeEntryUpdateTimeout, state.MirroredSubtreeEntryUpdateTimeout, "mirrored-subtree-entry-update-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.MirroredSubtreeSearchTimeout, state.MirroredSubtreeSearchTimeout, "mirrored-subtree-search-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.SetDegradedAlertWhenDisabled, state.SetDegradedAlertWhenDisabled, "set-degraded-alert-when-disabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.ReturnUnavailableWhenDisabled, state.ReturnUnavailableWhenDisabled, "return-unavailable-when-disabled")
	operations.AddStringOperationIfNecessary(&ops, plan.BackupFilePermissions, state.BackupFilePermissions, "backup-file-permissions")
	operations.AddStringOperationIfNecessary(&ops, plan.NotificationManager, state.NotificationManager, "notification-manager")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *configFileHandlerBackendResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan configFileHandlerBackendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.BackendApi.GetBackend(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.BackendID.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Config File Handler Backend", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state configFileHandlerBackendResourceModel
	readConfigFileHandlerBackendResponse(ctx, readResponse.ConfigFileHandlerBackendResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.BackendApi.UpdateBackend(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.BackendID.ValueString())
	ops := createConfigFileHandlerBackendOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.BackendApi.UpdateBackendExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Config File Handler Backend", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readConfigFileHandlerBackendResponse(ctx, updateResponse.ConfigFileHandlerBackendResponse, &state, &plan, &resp.Diagnostics)
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
func (r *configFileHandlerBackendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state configFileHandlerBackendResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.BackendApi.GetBackend(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.BackendID.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Config File Handler Backend", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readConfigFileHandlerBackendResponse(ctx, readResponse.ConfigFileHandlerBackendResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *configFileHandlerBackendResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan configFileHandlerBackendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state configFileHandlerBackendResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.BackendApi.UpdateBackend(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.BackendID.ValueString())

	// Determine what update operations are necessary
	ops := createConfigFileHandlerBackendOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.BackendApi.UpdateBackendExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Config File Handler Backend", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readConfigFileHandlerBackendResponse(ctx, updateResponse.ConfigFileHandlerBackendResponse, &state, &plan, &resp.Diagnostics)
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
func (r *configFileHandlerBackendResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *configFileHandlerBackendResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to backend_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("backend_id"), req, resp)
}
