package pingdirectory

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingdata-config-api-go-client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &locationsResource{}
	_ resource.ResourceWithConfigure   = &locationsResource{}
	_ resource.ResourceWithImportState = &locationsResource{}
)

// NewLocationsResource is a helper function to simplify the provider implementation.
func NewLocationsResource() resource.Resource {
	return &locationsResource{}
}

// locationsResource is the resource implementation.
type locationsResource struct {
	providerConfig pingdirectoryProviderModel
	apiClient      *client.APIClient
}

// locationsResourceModel maps the resource schema data.
type locationsResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

// Metadata returns the resource type name.
func (r *locationsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_location"
}

// GetSchema defines the schema for the resource.
func (r *locationsResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Manages a location.",
		Attributes: map[string]tfsdk.Attribute{
			"name": {
				Description: "Name of the location.",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"description": {
				Description: "Description of the location.",
				Type:        types.StringType,
				Optional:    true,
			},
			"last_updated": {
				Description: "Timestamp of the last Terraform update of the location.",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}, nil
}

// Configure adds the provider configured client to the resource.
func (r *locationsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(locationsResource)
	r.providerConfig = providerCfg.providerConfig
	r.apiClient = providerCfg.apiClient
}

//TODO does it make sense to do this for each call?
func (r *locationsResource) BasicAuthContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{
		UserName: r.providerConfig.Username.Value,
		Password: r.providerConfig.Password.Value,
	})
}

// Create a new resource
func (r *locationsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan locationsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addLocRequest := client.NewAddLocationRequest(plan.Name.Value)
	addLocRequest.Description = &plan.Description.Value
	apiAddLocationRequest := r.apiClient.LocationApi.AddLocation(r.BasicAuthContext(ctx))
	apiAddLocationRequest = apiAddLocationRequest.AddLocationRequest(*addLocRequest)

	//TODO any reason to look at the http response here rather than just checking the error? Maybe for a more descriptive error?
	_, err := r.apiClient.LocationApi.AddLocationExecute(apiAddLocationRequest)
	if err != nil {
		resp.Diagnostics.AddError("An error occurred while creating the Location", err.Error())
		return
	}

	// Populate Computed attribute values
	plan.LastUpdated = types.String{Value: string(time.Now().Format(time.RFC850))}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *locationsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state locationsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO again any reason to use the HTTP response?
	getResp, _, err := r.apiClient.LocationApi.GetLocation(r.BasicAuthContext(ctx), state.Name.Value).Execute()
	if err != nil {
		resp.Diagnostics.AddError("An error occurred while getting the Location", err.Error())
		return
	}

	// Read the updated description
	state.Description = types.String{Value: *getResp.Description}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *locationsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan locationsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//TODO update the location with the client here
	updateOperation := client.NewOperation()
	operation := "replace"
	path := "description"
	updateOperation.Op = &operation
	updateOperation.Path = &path
	//TODO this breaks when removing the description
	updateOperation.Value = &plan.Description.Value
	updateLocRequest := r.apiClient.LocationApi.UpdateLocation(r.BasicAuthContext(ctx), plan.Name.Value)
	updateLocRequest = updateLocRequest.UpdateLocationRequest(*client.NewUpdateLocationRequest([]client.Operation{*updateOperation}))
	// TODO again any reason to use the HTTP response?
	_, err := r.apiClient.LocationApi.UpdateLocationExecute(updateLocRequest)
	if err != nil {
		resp.Diagnostics.AddError("An error occurred while updating the Location", err.Error())
		return
	}

	// Update resource state with updated items and timestamp
	plan.LastUpdated = types.String{Value: string(time.Now().Format(time.RFC850))}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *locationsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state locationsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//TODO use for HTTP response?
	_, err := r.apiClient.LocationApi.DeleteLocationExecute(r.apiClient.LocationApi.DeleteLocation(r.BasicAuthContext(ctx), state.Name.Value))
	if err != nil {
		resp.Diagnostics.AddError("An error occurred while deleting the Location", err.Error())
		return
	}
}

func (r *locationsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to Name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
