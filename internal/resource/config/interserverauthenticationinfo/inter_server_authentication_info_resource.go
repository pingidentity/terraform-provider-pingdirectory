package interserverauthenticationinfo

import (
	"context"
	"strings"
	"time"

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
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &interServerAuthenticationInfoResource{}
	_ resource.ResourceWithConfigure   = &interServerAuthenticationInfoResource{}
	_ resource.ResourceWithImportState = &interServerAuthenticationInfoResource{}
)

// Create a Inter Server Authentication Info resource
func NewInterServerAuthenticationInfoResource() resource.Resource {
	return &interServerAuthenticationInfoResource{}
}

// interServerAuthenticationInfoResource is the resource implementation.
type interServerAuthenticationInfoResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *interServerAuthenticationInfoResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_inter_server_authentication_info"
}

// Configure adds the provider configured client to the resource.
func (r *interServerAuthenticationInfoResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type interServerAuthenticationInfoResourceModel struct {
	Id                         types.String `tfsdk:"id"`
	LastUpdated                types.String `tfsdk:"last_updated"`
	Notifications              types.Set    `tfsdk:"notifications"`
	RequiredActions            types.Set    `tfsdk:"required_actions"`
	Type                       types.String `tfsdk:"type"`
	ServerInstanceListenerName types.String `tfsdk:"server_instance_listener_name"`
	ServerInstanceName         types.String `tfsdk:"server_instance_name"`
	AuthenticationType         types.String `tfsdk:"authentication_type"`
	BindDN                     types.String `tfsdk:"bind_dn"`
	Username                   types.String `tfsdk:"username"`
	Password                   types.String `tfsdk:"password"`
	Purpose                    types.Set    `tfsdk:"purpose"`
}

// GetSchema defines the schema for the resource.
func (r *interServerAuthenticationInfoResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a Inter Server Authentication Info.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Inter Server Authentication Info resource. Options are ['password', 'certificate']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"password", "certificate"}...),
				},
			},
			"server_instance_listener_name": schema.StringAttribute{
				Description: "Name of the parent Server Instance Listener",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"server_instance_name": schema.StringAttribute{
				Description: "Name of the parent Server Instance",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"authentication_type": schema.StringAttribute{
				Description: "Identifies the type of password authentication that will be used.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bind_dn": schema.StringAttribute{
				Description: "A DN of the username that should be used for the bind request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"username": schema.StringAttribute{
				Description: "The username that should be used for the bind request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password": schema.StringAttribute{
				Description: "The password for the username or bind-dn.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"purpose": schema.SetAttribute{
				Description: "Identifies the purpose of this Inter Server Authentication Info.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan
func (r *interServerAuthenticationInfoResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var model interServerAuthenticationInfoResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.BindDN) && model.Type.ValueString() != "password" {
		resp.Diagnostics.AddError("Attribute 'bind_dn' not supported by pingdirectory_inter_server_authentication_info resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'bind_dn', the 'type' attribute must be one of ['password']")
	}
	if internaltypes.IsDefined(model.Password) && model.Type.ValueString() != "password" {
		resp.Diagnostics.AddError("Attribute 'password' not supported by pingdirectory_inter_server_authentication_info resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'password', the 'type' attribute must be one of ['password']")
	}
	if internaltypes.IsDefined(model.AuthenticationType) && model.Type.ValueString() != "password" {
		resp.Diagnostics.AddError("Attribute 'authentication_type' not supported by pingdirectory_inter_server_authentication_info resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'authentication_type', the 'type' attribute must be one of ['password']")
	}
	if internaltypes.IsDefined(model.Username) && model.Type.ValueString() != "password" {
		resp.Diagnostics.AddError("Attribute 'username' not supported by pingdirectory_inter_server_authentication_info resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'username', the 'type' attribute must be one of ['password']")
	}
}

// Read a PasswordInterServerAuthenticationInfoResponse object into the model struct
func readPasswordInterServerAuthenticationInfoResponse(ctx context.Context, r *client.PasswordInterServerAuthenticationInfoResponse, state *interServerAuthenticationInfoResourceModel, expectedValues *interServerAuthenticationInfoResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("password")
	state.Id = types.StringValue(r.Id)
	state.ServerInstanceListenerName = expectedValues.ServerInstanceListenerName
	state.ServerInstanceName = expectedValues.ServerInstanceName
	state.AuthenticationType = internaltypes.StringTypeOrNil(
		client.StringPointerEnuminterServerAuthenticationInfoAuthenticationTypeProp(r.AuthenticationType), true)
	state.BindDN = internaltypes.StringTypeOrNil(r.BindDN, true)
	state.Username = internaltypes.StringTypeOrNil(r.Username, true)
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.Password = expectedValues.Password
	state.Purpose = internaltypes.GetStringSet(
		client.StringSliceEnuminterServerAuthenticationInfoPurposeProp(r.Purpose))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Read a CertificateInterServerAuthenticationInfoResponse object into the model struct
func readCertificateInterServerAuthenticationInfoResponse(ctx context.Context, r *client.CertificateInterServerAuthenticationInfoResponse, state *interServerAuthenticationInfoResourceModel, expectedValues *interServerAuthenticationInfoResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("certificate")
	state.Id = types.StringValue(r.Id)
	state.ServerInstanceListenerName = expectedValues.ServerInstanceListenerName
	state.ServerInstanceName = expectedValues.ServerInstanceName
	state.Purpose = internaltypes.GetStringSet(
		client.StringSliceEnuminterServerAuthenticationInfoPurposeProp(r.Purpose))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createInterServerAuthenticationInfoOperations(plan interServerAuthenticationInfoResourceModel, state interServerAuthenticationInfoResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.AuthenticationType, state.AuthenticationType, "authentication-type")
	operations.AddStringOperationIfNecessary(&ops, plan.BindDN, state.BindDN, "bind-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.Username, state.Username, "username")
	operations.AddStringOperationIfNecessary(&ops, plan.Password, state.Password, "password")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.Purpose, state.Purpose, "purpose")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *interServerAuthenticationInfoResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan interServerAuthenticationInfoResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.InterServerAuthenticationInfoApi.GetInterServerAuthenticationInfo(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString(), plan.ServerInstanceListenerName.ValueString(), plan.ServerInstanceName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Inter Server Authentication Info", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state interServerAuthenticationInfoResourceModel
	if plan.Type.ValueString() == "password" {
		readPasswordInterServerAuthenticationInfoResponse(ctx, readResponse.PasswordInterServerAuthenticationInfoResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "certificate" {
		readCertificateInterServerAuthenticationInfoResponse(ctx, readResponse.CertificateInterServerAuthenticationInfoResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.InterServerAuthenticationInfoApi.UpdateInterServerAuthenticationInfo(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString(), plan.ServerInstanceListenerName.ValueString(), plan.ServerInstanceName.ValueString())
	ops := createInterServerAuthenticationInfoOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.InterServerAuthenticationInfoApi.UpdateInterServerAuthenticationInfoExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Inter Server Authentication Info", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "password" {
			readPasswordInterServerAuthenticationInfoResponse(ctx, updateResponse.PasswordInterServerAuthenticationInfoResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "certificate" {
			readCertificateInterServerAuthenticationInfoResponse(ctx, updateResponse.CertificateInterServerAuthenticationInfoResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *interServerAuthenticationInfoResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state interServerAuthenticationInfoResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.InterServerAuthenticationInfoApi.GetInterServerAuthenticationInfo(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString(), state.ServerInstanceListenerName.ValueString(), state.ServerInstanceName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Inter Server Authentication Info", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.PasswordInterServerAuthenticationInfoResponse != nil {
		readPasswordInterServerAuthenticationInfoResponse(ctx, readResponse.PasswordInterServerAuthenticationInfoResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CertificateInterServerAuthenticationInfoResponse != nil {
		readCertificateInterServerAuthenticationInfoResponse(ctx, readResponse.CertificateInterServerAuthenticationInfoResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *interServerAuthenticationInfoResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan interServerAuthenticationInfoResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state interServerAuthenticationInfoResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.InterServerAuthenticationInfoApi.UpdateInterServerAuthenticationInfo(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString(), plan.ServerInstanceListenerName.ValueString(), plan.ServerInstanceName.ValueString())

	// Determine what update operations are necessary
	ops := createInterServerAuthenticationInfoOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.InterServerAuthenticationInfoApi.UpdateInterServerAuthenticationInfoExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Inter Server Authentication Info", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "password" {
			readPasswordInterServerAuthenticationInfoResponse(ctx, updateResponse.PasswordInterServerAuthenticationInfoResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "certificate" {
			readCertificateInterServerAuthenticationInfoResponse(ctx, updateResponse.CertificateInterServerAuthenticationInfoResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *interServerAuthenticationInfoResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *interServerAuthenticationInfoResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 3 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [server-instance-name]/[server-instance-listener-name]/[inter-server-authentication-info-name]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("server_instance_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("server_instance_listener_name"), split[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), split[2])...)
}
