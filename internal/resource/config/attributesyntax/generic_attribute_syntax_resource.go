package attributesyntax

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &genericAttributeSyntaxResource{}
	_ resource.ResourceWithConfigure   = &genericAttributeSyntaxResource{}
	_ resource.ResourceWithImportState = &genericAttributeSyntaxResource{}
)

// Create a Generic Attribute Syntax resource
func NewGenericAttributeSyntaxResource() resource.Resource {
	return &genericAttributeSyntaxResource{}
}

// genericAttributeSyntaxResource is the resource implementation.
type genericAttributeSyntaxResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *genericAttributeSyntaxResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_generic_attribute_syntax"
}

// Configure adds the provider configured client to the resource.
func (r *genericAttributeSyntaxResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type genericAttributeSyntaxResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	LastUpdated           types.String `tfsdk:"last_updated"`
	Notifications         types.Set    `tfsdk:"notifications"`
	RequiredActions       types.Set    `tfsdk:"required_actions"`
	Enabled               types.Bool   `tfsdk:"enabled"`
	RequireBinaryTransfer types.Bool   `tfsdk:"require_binary_transfer"`
}

// GetSchema defines the schema for the resource.
func (r *genericAttributeSyntaxResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Generic Attribute Syntax.",
		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Attribute Syntax is enabled.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"require_binary_transfer": schema.BoolAttribute{
				Description: "Indicates whether values of this attribute are required to have a \"binary\" transfer option as described in RFC 4522. Attributes with this syntax will generally be referenced with names including \";binary\" (e.g., \"userCertificate;binary\").",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Read a GenericAttributeSyntaxResponse object into the model struct
func readGenericAttributeSyntaxResponse(ctx context.Context, r *client.GenericAttributeSyntaxResponse, state *genericAttributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createGenericAttributeSyntaxOperations(plan genericAttributeSyntaxResourceModel, state genericAttributeSyntaxResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequireBinaryTransfer, state.RequireBinaryTransfer, "require-binary-transfer")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *genericAttributeSyntaxResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan genericAttributeSyntaxResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AttributeSyntaxApi.GetAttributeSyntax(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Generic Attribute Syntax", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state genericAttributeSyntaxResourceModel
	readGenericAttributeSyntaxResponse(ctx, readResponse.GenericAttributeSyntaxResponse, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.AttributeSyntaxApi.UpdateAttributeSyntax(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createGenericAttributeSyntaxOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AttributeSyntaxApi.UpdateAttributeSyntaxExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Generic Attribute Syntax", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGenericAttributeSyntaxResponse(ctx, updateResponse.GenericAttributeSyntaxResponse, &state, &resp.Diagnostics)
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
func (r *genericAttributeSyntaxResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state genericAttributeSyntaxResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AttributeSyntaxApi.GetAttributeSyntax(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Generic Attribute Syntax", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readGenericAttributeSyntaxResponse(ctx, readResponse.GenericAttributeSyntaxResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *genericAttributeSyntaxResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan genericAttributeSyntaxResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state genericAttributeSyntaxResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.AttributeSyntaxApi.UpdateAttributeSyntax(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createGenericAttributeSyntaxOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AttributeSyntaxApi.UpdateAttributeSyntaxExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Generic Attribute Syntax", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGenericAttributeSyntaxResponse(ctx, updateResponse.GenericAttributeSyntaxResponse, &state, &resp.Diagnostics)
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
func (r *genericAttributeSyntaxResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *genericAttributeSyntaxResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
