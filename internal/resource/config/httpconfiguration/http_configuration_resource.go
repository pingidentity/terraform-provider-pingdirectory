// Copyright © 2025 Ping Identity Corporation

package httpconfiguration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &httpConfigurationResource{}
	_ resource.ResourceWithConfigure   = &httpConfigurationResource{}
	_ resource.ResourceWithImportState = &httpConfigurationResource{}
)

// Create a Http Configuration resource
func NewHttpConfigurationResource() resource.Resource {
	return &httpConfigurationResource{}
}

// httpConfigurationResource is the resource implementation.
type httpConfigurationResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *httpConfigurationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_http_configuration"
}

// Configure adds the provider configured client to the resource.
func (r *httpConfigurationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type httpConfigurationResourceModel struct {
	Id                                    types.String `tfsdk:"id"`
	Notifications                         types.Set    `tfsdk:"notifications"`
	RequiredActions                       types.Set    `tfsdk:"required_actions"`
	Type                                  types.String `tfsdk:"type"`
	IncludeStackTracesInErrorPages        types.Bool   `tfsdk:"include_stack_traces_in_error_pages"`
	IncludeServletInformationInErrorPages types.Bool   `tfsdk:"include_servlet_information_in_error_pages"`
}

// GetSchema defines the schema for the resource.
func (r *httpConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a Http Configuration.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of HTTP Configuration resource. Options are ['http-configuration']",
				Optional:    false,
				Required:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"http-configuration"}...),
				},
			},
			"include_stack_traces_in_error_pages": schema.BoolAttribute{
				Description: "Indicates whether exceptions thrown by servlet or web application extensions will be included in the resulting error page response. Stack traces can be helpful in diagnosing application errors, but in production they may reveal information that might be useful to a malicious attacker.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_servlet_information_in_error_pages": schema.BoolAttribute{
				Description: "Supported in PingDirectory product version 9.3.0.0+. Indicates whether to expose servlet information in the error page response.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonResourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *httpConfigurationResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingDirectory9300)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model httpConfigurationResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.IncludeServletInformationInErrorPages) {
		resp.Diagnostics.AddError("Attribute 'include_servlet_information_in_error_pages' not supported by PingDirectory version "+r.providerConfig.ProductVersion, "")
	}
}

// Read a HttpConfigurationResponse object into the model struct
func readHttpConfigurationResponse(ctx context.Context, r *client.HttpConfigurationResponse, state *httpConfigurationResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("http-configuration")
	// Placeholder id value required by test framework
	state.Id = types.StringValue("id")
	state.IncludeStackTracesInErrorPages = internaltypes.BoolTypeOrNil(r.IncludeStackTracesInErrorPages)
	state.IncludeServletInformationInErrorPages = internaltypes.BoolTypeOrNil(r.IncludeServletInformationInErrorPages)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createHttpConfigurationOperations(plan httpConfigurationResourceModel, state httpConfigurationResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeStackTracesInErrorPages, state.IncludeStackTracesInErrorPages, "include-stack-traces-in-error-pages")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeServletInformationInErrorPages, state.IncludeServletInformationInErrorPages, "include-servlet-information-in-error-pages")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *httpConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan httpConfigurationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.HttpConfigurationAPI.GetHttpConfiguration(
		config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Http Configuration", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state httpConfigurationResourceModel
	readHttpConfigurationResponse(ctx, readResponse, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.HttpConfigurationAPI.UpdateHttpConfiguration(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	ops := createHttpConfigurationOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.HttpConfigurationAPI.UpdateHttpConfigurationExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Http Configuration", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readHttpConfigurationResponse(ctx, updateResponse, &state, &resp.Diagnostics)
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *httpConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state httpConfigurationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.HttpConfigurationAPI.GetHttpConfiguration(
		config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Http Configuration", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readHttpConfigurationResponse(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *httpConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan httpConfigurationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state httpConfigurationResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.HttpConfigurationAPI.UpdateHttpConfiguration(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))

	// Determine what update operations are necessary
	ops := createHttpConfigurationOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.HttpConfigurationAPI.UpdateHttpConfigurationExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Http Configuration", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readHttpConfigurationResponse(ctx, updateResponse, &state, &resp.Diagnostics)
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
func (r *httpConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *httpConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Set a placeholder id value to appease terraform.
	// The real attributes will be imported when terraform performs a read after the import.
	// If no value is set here, Terraform will error out when importing.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), "id")...)
}
