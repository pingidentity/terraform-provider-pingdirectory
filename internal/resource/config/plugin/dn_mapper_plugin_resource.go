package plugin

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
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dnMapperPluginResource{}
	_ resource.ResourceWithConfigure   = &dnMapperPluginResource{}
	_ resource.ResourceWithImportState = &dnMapperPluginResource{}
	_ resource.Resource                = &defaultDnMapperPluginResource{}
	_ resource.ResourceWithConfigure   = &defaultDnMapperPluginResource{}
	_ resource.ResourceWithImportState = &defaultDnMapperPluginResource{}
)

// Create a Dn Mapper Plugin resource
func NewDnMapperPluginResource() resource.Resource {
	return &dnMapperPluginResource{}
}

func NewDefaultDnMapperPluginResource() resource.Resource {
	return &defaultDnMapperPluginResource{}
}

// dnMapperPluginResource is the resource implementation.
type dnMapperPluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultDnMapperPluginResource is the resource implementation.
type defaultDnMapperPluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *dnMapperPluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dn_mapper_plugin"
}

func (r *defaultDnMapperPluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_dn_mapper_plugin"
}

// Configure adds the provider configured client to the resource.
func (r *dnMapperPluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultDnMapperPluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type dnMapperPluginResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	LastUpdated                 types.String `tfsdk:"last_updated"`
	Notifications               types.Set    `tfsdk:"notifications"`
	RequiredActions             types.Set    `tfsdk:"required_actions"`
	PluginType                  types.Set    `tfsdk:"plugin_type"`
	SourceDN                    types.String `tfsdk:"source_dn"`
	TargetDN                    types.String `tfsdk:"target_dn"`
	EnableAttributeMapping      types.Bool   `tfsdk:"enable_attribute_mapping"`
	MapAttribute                types.Set    `tfsdk:"map_attribute"`
	EnableControlMapping        types.Bool   `tfsdk:"enable_control_mapping"`
	AlwaysMapResponses          types.Bool   `tfsdk:"always_map_responses"`
	Description                 types.String `tfsdk:"description"`
	Enabled                     types.Bool   `tfsdk:"enabled"`
	InvokeForInternalOperations types.Bool   `tfsdk:"invoke_for_internal_operations"`
}

// GetSchema defines the schema for the resource.
func (r *dnMapperPluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	dnMapperPluginSchema(ctx, req, resp, false)
}

func (r *defaultDnMapperPluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	dnMapperPluginSchema(ctx, req, resp, true)
}

func dnMapperPluginSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Dn Mapper Plugin.",
		Attributes: map[string]schema.Attribute{
			"plugin_type": schema.SetAttribute{
				Description: "Specifies the set of plug-in types for the plug-in, which specifies the times at which the plug-in is invoked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"source_dn": schema.StringAttribute{
				Description: "Specifies the source DN that may appear in client requests which should be remapped to the target DN. Note that the source DN must not be equal to the target DN.",
				Required:    true,
			},
			"target_dn": schema.StringAttribute{
				Description: "Specifies the DN to which the source DN should be mapped. Note that the target DN must not be equal to the source DN.",
				Required:    true,
			},
			"enable_attribute_mapping": schema.BoolAttribute{
				Description: "Indicates whether DN mapping should be applied to the values of attributes with appropriate syntaxes.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"map_attribute": schema.SetAttribute{
				Description: "Specifies a set of specific attributes for which DN mapping should be applied. This will only be applicable if the enable-attribute-mapping property has a value of \"true\". Any attributes listed must be defined in the server schema with either the distinguished name syntax or the name and optional UID syntax.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"enable_control_mapping": schema.BoolAttribute{
				Description: "Indicates whether DN mapping should be applied to DNs that may be present in specific controls. DN mapping will only be applied for control types which are specifically supported by the DN mapper plugin.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"always_map_responses": schema.BoolAttribute{
				Description: "Indicates whether DNs in response messages containing the target DN should always be remapped back to the source DN. If this is \"false\", then mapping will be performed for a response message only if one or more elements of the associated request are mapped. Otherwise, the mapping will be performed for all responses regardless of whether the mapping was applied to the request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Plugin",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the plug-in is enabled for use.",
				Required:    true,
			},
			"invoke_for_internal_operations": schema.BoolAttribute{
				Description: "Indicates whether the plug-in should be invoked for internal operations.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
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
func addOptionalDnMapperPluginFields(ctx context.Context, addRequest *client.AddDnMapperPluginRequest, plan dnMapperPluginResourceModel) error {
	if internaltypes.IsDefined(plan.PluginType) {
		var slice []string
		plan.PluginType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginPluginTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginPluginTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.PluginType = enumSlice
	}
	if internaltypes.IsDefined(plan.EnableAttributeMapping) {
		boolVal := plan.EnableAttributeMapping.ValueBool()
		addRequest.EnableAttributeMapping = &boolVal
	}
	if internaltypes.IsDefined(plan.MapAttribute) {
		var slice []string
		plan.MapAttribute.ElementsAs(ctx, &slice, false)
		addRequest.MapAttribute = slice
	}
	if internaltypes.IsDefined(plan.EnableControlMapping) {
		boolVal := plan.EnableControlMapping.ValueBool()
		addRequest.EnableControlMapping = &boolVal
	}
	if internaltypes.IsDefined(plan.AlwaysMapResponses) {
		boolVal := plan.AlwaysMapResponses.ValueBool()
		addRequest.AlwaysMapResponses = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	if internaltypes.IsDefined(plan.InvokeForInternalOperations) {
		boolVal := plan.InvokeForInternalOperations.ValueBool()
		addRequest.InvokeForInternalOperations = &boolVal
	}
	return nil
}

// Read a DnMapperPluginResponse object into the model struct
func readDnMapperPluginResponse(ctx context.Context, r *client.DnMapperPluginResponse, state *dnMapperPluginResourceModel, expectedValues *dnMapperPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.SourceDN = types.StringValue(r.SourceDN)
	state.TargetDN = types.StringValue(r.TargetDN)
	state.EnableAttributeMapping = types.BoolValue(r.EnableAttributeMapping)
	state.MapAttribute = internaltypes.GetStringSet(r.MapAttribute)
	state.EnableControlMapping = types.BoolValue(r.EnableControlMapping)
	state.AlwaysMapResponses = types.BoolValue(r.AlwaysMapResponses)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createDnMapperPluginOperations(plan dnMapperPluginResourceModel, state dnMapperPluginResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.PluginType, state.PluginType, "plugin-type")
	operations.AddStringOperationIfNecessary(&ops, plan.SourceDN, state.SourceDN, "source-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.TargetDN, state.TargetDN, "target-dn")
	operations.AddBoolOperationIfNecessary(&ops, plan.EnableAttributeMapping, state.EnableAttributeMapping, "enable-attribute-mapping")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.MapAttribute, state.MapAttribute, "map-attribute")
	operations.AddBoolOperationIfNecessary(&ops, plan.EnableControlMapping, state.EnableControlMapping, "enable-control-mapping")
	operations.AddBoolOperationIfNecessary(&ops, plan.AlwaysMapResponses, state.AlwaysMapResponses, "always-map-responses")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.InvokeForInternalOperations, state.InvokeForInternalOperations, "invoke-for-internal-operations")
	return ops
}

// Create a new resource
func (r *dnMapperPluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan dnMapperPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddDnMapperPluginRequest(plan.Id.ValueString(),
		[]client.EnumdnMapperPluginSchemaUrn{client.ENUMDNMAPPERPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINDN_MAPPER},
		plan.SourceDN.ValueString(),
		plan.TargetDN.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalDnMapperPluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Dn Mapper Plugin", err.Error())
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
		client.AddDnMapperPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Dn Mapper Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state dnMapperPluginResourceModel
	readDnMapperPluginResponse(ctx, addResponse.DnMapperPluginResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultDnMapperPluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan dnMapperPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Dn Mapper Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state dnMapperPluginResourceModel
	readDnMapperPluginResponse(ctx, readResponse.DnMapperPluginResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.PluginApi.UpdatePlugin(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createDnMapperPluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Dn Mapper Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDnMapperPluginResponse(ctx, updateResponse.DnMapperPluginResponse, &state, &plan, &resp.Diagnostics)
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
func (r *dnMapperPluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDnMapperPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultDnMapperPluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDnMapperPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readDnMapperPlugin(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state dnMapperPluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Dn Mapper Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDnMapperPluginResponse(ctx, readResponse.DnMapperPluginResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *dnMapperPluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDnMapperPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultDnMapperPluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDnMapperPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateDnMapperPlugin(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan dnMapperPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state dnMapperPluginResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.PluginApi.UpdatePlugin(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createDnMapperPluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Dn Mapper Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDnMapperPluginResponse(ctx, updateResponse.DnMapperPluginResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultDnMapperPluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *dnMapperPluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state dnMapperPluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PluginApi.DeletePluginExecute(r.apiClient.PluginApi.DeletePlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Dn Mapper Plugin", err, httpResp)
		return
	}
}

func (r *dnMapperPluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDnMapperPlugin(ctx, req, resp)
}

func (r *defaultDnMapperPluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDnMapperPlugin(ctx, req, resp)
}

func importDnMapperPlugin(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
