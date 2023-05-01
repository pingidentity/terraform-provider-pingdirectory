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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &generalizedTimeAttributeSyntaxResource{}
	_ resource.ResourceWithConfigure   = &generalizedTimeAttributeSyntaxResource{}
	_ resource.ResourceWithImportState = &generalizedTimeAttributeSyntaxResource{}
)

// Create a Generalized Time Attribute Syntax resource
func NewGeneralizedTimeAttributeSyntaxResource() resource.Resource {
	return &generalizedTimeAttributeSyntaxResource{}
}

// generalizedTimeAttributeSyntaxResource is the resource implementation.
type generalizedTimeAttributeSyntaxResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *generalizedTimeAttributeSyntaxResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_generalized_time_attribute_syntax"
}

// Configure adds the provider configured client to the resource.
func (r *generalizedTimeAttributeSyntaxResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type generalizedTimeAttributeSyntaxResourceModel struct {
	Id                             types.String `tfsdk:"id"`
	LastUpdated                    types.String `tfsdk:"last_updated"`
	Notifications                  types.Set    `tfsdk:"notifications"`
	RequiredActions                types.Set    `tfsdk:"required_actions"`
	EnableCompaction               types.Bool   `tfsdk:"enable_compaction"`
	IncludeAttributeInCompaction   types.Set    `tfsdk:"include_attribute_in_compaction"`
	ExcludeAttributeFromCompaction types.Set    `tfsdk:"exclude_attribute_from_compaction"`
	Enabled                        types.Bool   `tfsdk:"enabled"`
	RequireBinaryTransfer          types.Bool   `tfsdk:"require_binary_transfer"`
}

// GetSchema defines the schema for the resource.
func (r *generalizedTimeAttributeSyntaxResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Generalized Time Attribute Syntax.",
		Attributes: map[string]schema.Attribute{
			"enable_compaction": schema.BoolAttribute{
				Description: "Indicates whether values of attributes with this syntax should be compacted when stored in a local DB database.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_attribute_in_compaction": schema.SetAttribute{
				Description: "Specifies the specific attributes (which should be associated with this syntax) whose values should be compacted. If one or more include attributes are specified, then only those attributes will have their values compacted. If not set then all attributes will have their values compacted. The exclude-attribute-from-compaction property takes precedence over this property.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"exclude_attribute_from_compaction": schema.SetAttribute{
				Description: "Specifies the specific attributes (which should be associated with this syntax) whose values should not be compacted. If one or more exclude attributes are specified, then values of those attributes will not have their values compacted. This property takes precedence over the include-attribute-in-compaction property.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
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

// Read a GeneralizedTimeAttributeSyntaxResponse object into the model struct
func readGeneralizedTimeAttributeSyntaxResponse(ctx context.Context, r *client.GeneralizedTimeAttributeSyntaxResponse, state *generalizedTimeAttributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createGeneralizedTimeAttributeSyntaxOperations(plan generalizedTimeAttributeSyntaxResourceModel, state generalizedTimeAttributeSyntaxResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddBoolOperationIfNecessary(&ops, plan.EnableCompaction, state.EnableCompaction, "enable-compaction")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeAttributeInCompaction, state.IncludeAttributeInCompaction, "include-attribute-in-compaction")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeAttributeFromCompaction, state.ExcludeAttributeFromCompaction, "exclude-attribute-from-compaction")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequireBinaryTransfer, state.RequireBinaryTransfer, "require-binary-transfer")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *generalizedTimeAttributeSyntaxResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan generalizedTimeAttributeSyntaxResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AttributeSyntaxApi.GetAttributeSyntax(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Generalized Time Attribute Syntax", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state generalizedTimeAttributeSyntaxResourceModel
	readGeneralizedTimeAttributeSyntaxResponse(ctx, readResponse.GeneralizedTimeAttributeSyntaxResponse, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.AttributeSyntaxApi.UpdateAttributeSyntax(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createGeneralizedTimeAttributeSyntaxOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AttributeSyntaxApi.UpdateAttributeSyntaxExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Generalized Time Attribute Syntax", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGeneralizedTimeAttributeSyntaxResponse(ctx, updateResponse.GeneralizedTimeAttributeSyntaxResponse, &state, &resp.Diagnostics)
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
func (r *generalizedTimeAttributeSyntaxResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state generalizedTimeAttributeSyntaxResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AttributeSyntaxApi.GetAttributeSyntax(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Generalized Time Attribute Syntax", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readGeneralizedTimeAttributeSyntaxResponse(ctx, readResponse.GeneralizedTimeAttributeSyntaxResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *generalizedTimeAttributeSyntaxResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan generalizedTimeAttributeSyntaxResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state generalizedTimeAttributeSyntaxResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.AttributeSyntaxApi.UpdateAttributeSyntax(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createGeneralizedTimeAttributeSyntaxOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AttributeSyntaxApi.UpdateAttributeSyntaxExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Generalized Time Attribute Syntax", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGeneralizedTimeAttributeSyntaxResponse(ctx, updateResponse.GeneralizedTimeAttributeSyntaxResponse, &state, &resp.Diagnostics)
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
func (r *generalizedTimeAttributeSyntaxResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *generalizedTimeAttributeSyntaxResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
