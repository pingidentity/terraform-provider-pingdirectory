package pingdirectory

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdata-config-api-go-client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &blindTrustManagerProviderResource{}
	_ resource.ResourceWithConfigure   = &blindTrustManagerProviderResource{}
	_ resource.ResourceWithImportState = &blindTrustManagerProviderResource{}
)

// NewBlindTrustManagerProviderResource is a helper function to simplify the provider implementation.
func NewBlindTrustManagerProviderResource() resource.Resource {
	return &blindTrustManagerProviderResource{}
}

// blindTrustManagerProviderResource is the resource implementation.
type blindTrustManagerProviderResource struct {
	providerConfig pingdirectoryProviderModel
	apiClient      *client.APIClient
}

// blindTrustManagerProviderResourceModel maps the resource schema data.
type blindTrustManagerProviderResourceModel struct {
	Name                     types.String `tfsdk:"name"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	IncludeJVMDefaultIssuers types.Bool   `tfsdk:"include_jvm_default_issuers"`
	LastUpdated              types.String `tfsdk:"last_updated"`
}

// Metadata returns the resource type name.
func (r *blindTrustManagerProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blind_trust_manager_provider"
}

// GetSchema defines the schema for the resource.
func (r *blindTrustManagerProviderResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Manages a Blind Trust Manager Provider.",
		Attributes: map[string]tfsdk.Attribute{
			"name": {
				Description: "Name of the Trust Manager Provider.",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"enabled": {
				Description: "Indicate whether the Trust Manager Provider is enabled for use.",
				Type:        types.BoolType,
				Required:    true,
			},
			"include_jvm_default_issuers": {
				Description: "Indicates whether certificates issued by an authority included in the JVM's set of default issuers should be automatically trusted, even if they would not otherwise be trusted by this provider.",
				Type:        types.BoolType,
				//TODO is this right? Do I need a default value for this via a plan modifier? If we have defaults then they'd have to be updated whenever directory defaults change.
				Optional: true,
				Computed: true,
			},
			"last_updated": {
				Description: "Timestamp of the last Terraform update of the Trust Manager Provider.",
				Type:        types.StringType,
				Computed:    true,
				Required:    false,
				Optional:    false,
			},
		},
	}, nil
}

// Configure adds the provider configured client to the resource.
func (r *blindTrustManagerProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(apiClientConfig)
	r.providerConfig = providerCfg.providerConfig
	r.apiClient = providerCfg.apiClient
}

// Create a new resource
func (r *blindTrustManagerProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan blindTrustManagerProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addProviderRequest := client.NewAddBlindTrustManagerProviderRequest(plan.Name.Value,
		[]client.EnumblindTrustManagerProviderSchemaUrn{client.ENUMBLINDTRUSTMANAGERPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0TRUST_MANAGER_PROVIDERBLIND},
		plan.Enabled.Value)
	if !plan.IncludeJVMDefaultIssuers.IsNull() && !plan.IncludeJVMDefaultIssuers.IsUnknown() {
		boolVal := plan.IncludeJVMDefaultIssuers.ValueBool()
		addProviderRequest.IncludeJVMDefaultIssuers = &boolVal
	}
	apiAddProviderRequest := r.apiClient.TrustManagerProviderApi.AddTrustManagerProvider(BasicAuthContext(ctx, r.providerConfig))
	apiAddProviderRequest = apiAddProviderRequest.AddTrustManagerProviderRequest(
		client.AddBlindTrustManagerProviderRequestAsAddTrustManagerProviderRequest(addProviderRequest))

	trustManagerResponse, httpResp, err := r.apiClient.TrustManagerProviderApi.AddTrustManagerProviderExecute(apiAddProviderRequest)
	if err != nil {
		ReportHttpError(&resp.Diagnostics, "An error occurred while creating the Trust Manager Provider", err, httpResp)
		return
	}

	// Read the response into the state
	ReadBlindTrustManagerProviderResponse(trustManagerResponse.BlindTrustManagerProviderResponse, &plan)

	// Populate Computed attribute values
	plan.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read a BlindTrustManagerProviderResponse object into the model struct
func ReadBlindTrustManagerProviderResponse(r *client.BlindTrustManagerProviderResponse, state *blindTrustManagerProviderResourceModel) {
	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeJVMDefaultIssuers = BoolTypeOrNil(r.IncludeJVMDefaultIssuers)
}

// Read resource information
func (r *blindTrustManagerProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state blindTrustManagerProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	trustManagerResponse, httpResp, err := r.apiClient.TrustManagerProviderApi.GetTrustManagerProvider(
		BasicAuthContext(ctx, r.providerConfig), state.Name.Value).Execute()
	if err != nil {
		ReportHttpError(&resp.Diagnostics, "An error occurred while getting the Trust Manager Provider", err, httpResp)
		return
	}

	// Read the response into the state
	ReadBlindTrustManagerProviderResponse(trustManagerResponse.BlindTrustManagerProviderResponse, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create any update operations necessary to make the state match the plan
func CreateBlindTrustManagerProviderOperations(plan blindTrustManagerProviderResourceModel, state blindTrustManagerProviderResourceModel) []client.Operation {
	var ops []client.Operation

	AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	AddBoolOperationIfNecessary(&ops, plan.IncludeJVMDefaultIssuers, state.IncludeJVMDefaultIssuers, "include-jvm-default-issuers")
	return ops
}

// Update a resource
func (r *blindTrustManagerProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan blindTrustManagerProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state blindTrustManagerProviderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.TrustManagerProviderApi.UpdateTrustManagerProvider(BasicAuthContext(ctx, r.providerConfig), plan.Name.Value)

	// Determine what update operations are necessary
	ops := CreateBlindTrustManagerProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))

		trustManagerResponse, httpResp, err := r.apiClient.TrustManagerProviderApi.UpdateTrustManagerProviderExecute(updateRequest)
		if err != nil {
			ReportHttpError(&resp.Diagnostics, "An error occurred while updating the Trust Manager Provider", err, httpResp)
			return
		}

		// Read the response
		ReadBlindTrustManagerProviderResponse(trustManagerResponse.BlindTrustManagerProviderResponse, &plan)
		// Update computed values
		plan.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *blindTrustManagerProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state blindTrustManagerProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.TrustManagerProviderApi.DeleteTrustManagerProviderExecute(
		r.apiClient.TrustManagerProviderApi.DeleteTrustManagerProvider(BasicAuthContext(ctx, r.providerConfig), state.Name.Value))
	if err != nil {
		ReportHttpError(&resp.Diagnostics, "An error occurred while deleting the Trust Manager Provider", err, httpResp)
		return
	}
}

func (r *blindTrustManagerProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to Name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
