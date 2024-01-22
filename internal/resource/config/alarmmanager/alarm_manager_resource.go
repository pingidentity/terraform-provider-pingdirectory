package alarmmanager

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10000/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &alarmManagerResource{}
	_ resource.ResourceWithConfigure   = &alarmManagerResource{}
	_ resource.ResourceWithImportState = &alarmManagerResource{}
)

// Create a Alarm Manager resource
func NewAlarmManagerResource() resource.Resource {
	return &alarmManagerResource{}
}

// alarmManagerResource is the resource implementation.
type alarmManagerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *alarmManagerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_alarm_manager"
}

// Configure adds the provider configured client to the resource.
func (r *alarmManagerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type alarmManagerResourceModel struct {
	Id                     types.String `tfsdk:"id"`
	Notifications          types.Set    `tfsdk:"notifications"`
	RequiredActions        types.Set    `tfsdk:"required_actions"`
	Type                   types.String `tfsdk:"type"`
	DefaultGaugeAlertLevel types.String `tfsdk:"default_gauge_alert_level"`
	GeneratedAlertTypes    types.Set    `tfsdk:"generated_alert_types"`
	SuppressedAlarm        types.Set    `tfsdk:"suppressed_alarm"`
}

// GetSchema defines the schema for the resource.
func (r *alarmManagerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a Alarm Manager.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Alarm Manager resource. Options are ['alarm-manager']",
				Optional:    false,
				Required:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"alarm-manager"}...),
				},
			},
			"default_gauge_alert_level": schema.StringAttribute{
				Description: "Specifies the level at which alerts are sent for alarms raised by the Alarm Manager.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"generated_alert_types": schema.SetAttribute{
				Description: "Indicates what kind of alert types should be generated.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"suppressed_alarm": schema.SetAttribute{
				Description: "Specifies the names of the alarm alert types that should be suppressed. If the condition that triggers an alarm in this list occurs, then the alarm will not be raised and no alerts will be generated. Only a subset of alarms can be suppressed in this way. Alarms triggered by a gauge can be disabled by disabling the gauge.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonResourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a AlarmManagerResponse object into the model struct
func readAlarmManagerResponse(ctx context.Context, r *client.AlarmManagerResponse, state *alarmManagerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("alarm-manager")
	// Placeholder id value required by test framework
	state.Id = types.StringValue("id")
	state.DefaultGaugeAlertLevel = types.StringValue(r.DefaultGaugeAlertLevel.String())
	state.GeneratedAlertTypes = internaltypes.GetStringSet(
		client.StringSliceEnumalarmManagerGeneratedAlertTypesProp(r.GeneratedAlertTypes))
	state.SuppressedAlarm = internaltypes.GetStringSet(
		client.StringSliceEnumalarmManagerSuppressedAlarmProp(r.SuppressedAlarm))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createAlarmManagerOperations(plan alarmManagerResourceModel, state alarmManagerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.DefaultGaugeAlertLevel, state.DefaultGaugeAlertLevel, "default-gauge-alert-level")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.GeneratedAlertTypes, state.GeneratedAlertTypes, "generated-alert-types")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SuppressedAlarm, state.SuppressedAlarm, "suppressed-alarm")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *alarmManagerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan alarmManagerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AlarmManagerAPI.GetAlarmManager(
		config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Alarm Manager", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state alarmManagerResourceModel
	readAlarmManagerResponse(ctx, readResponse, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.AlarmManagerAPI.UpdateAlarmManager(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	ops := createAlarmManagerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AlarmManagerAPI.UpdateAlarmManagerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Alarm Manager", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAlarmManagerResponse(ctx, updateResponse, &state, &resp.Diagnostics)
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *alarmManagerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state alarmManagerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AlarmManagerAPI.GetAlarmManager(
		config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Alarm Manager", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAlarmManagerResponse(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *alarmManagerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan alarmManagerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state alarmManagerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.AlarmManagerAPI.UpdateAlarmManager(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))

	// Determine what update operations are necessary
	ops := createAlarmManagerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AlarmManagerAPI.UpdateAlarmManagerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Alarm Manager", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAlarmManagerResponse(ctx, updateResponse, &state, &resp.Diagnostics)
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
func (r *alarmManagerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *alarmManagerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Set a placeholder id value to appease terraform.
	// The real attributes will be imported when terraform performs a read after the import.
	// If no value is set here, Terraform will error out when importing.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), "id")...)
}
