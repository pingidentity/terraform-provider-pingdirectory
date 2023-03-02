package logpublisher

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
	_ resource.Resource                = &groovyScriptedErrorLogPublisherResource{}
	_ resource.ResourceWithConfigure   = &groovyScriptedErrorLogPublisherResource{}
	_ resource.ResourceWithImportState = &groovyScriptedErrorLogPublisherResource{}
	_ resource.Resource                = &defaultGroovyScriptedErrorLogPublisherResource{}
	_ resource.ResourceWithConfigure   = &defaultGroovyScriptedErrorLogPublisherResource{}
	_ resource.ResourceWithImportState = &defaultGroovyScriptedErrorLogPublisherResource{}
)

// Create a Groovy Scripted Error Log Publisher resource
func NewGroovyScriptedErrorLogPublisherResource() resource.Resource {
	return &groovyScriptedErrorLogPublisherResource{}
}

func NewDefaultGroovyScriptedErrorLogPublisherResource() resource.Resource {
	return &defaultGroovyScriptedErrorLogPublisherResource{}
}

// groovyScriptedErrorLogPublisherResource is the resource implementation.
type groovyScriptedErrorLogPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultGroovyScriptedErrorLogPublisherResource is the resource implementation.
type defaultGroovyScriptedErrorLogPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *groovyScriptedErrorLogPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groovy_scripted_error_log_publisher"
}

func (r *defaultGroovyScriptedErrorLogPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_groovy_scripted_error_log_publisher"
}

// Configure adds the provider configured client to the resource.
func (r *groovyScriptedErrorLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultGroovyScriptedErrorLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type groovyScriptedErrorLogPublisherResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	LastUpdated          types.String `tfsdk:"last_updated"`
	Notifications        types.Set    `tfsdk:"notifications"`
	RequiredActions      types.Set    `tfsdk:"required_actions"`
	ScriptClass          types.String `tfsdk:"script_class"`
	ScriptArgument       types.Set    `tfsdk:"script_argument"`
	DefaultSeverity      types.Set    `tfsdk:"default_severity"`
	OverrideSeverity     types.Set    `tfsdk:"override_severity"`
	Description          types.String `tfsdk:"description"`
	Enabled              types.Bool   `tfsdk:"enabled"`
	LoggingErrorBehavior types.String `tfsdk:"logging_error_behavior"`
}

// GetSchema defines the schema for the resource.
func (r *groovyScriptedErrorLogPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	groovyScriptedErrorLogPublisherSchema(ctx, req, resp, false)
}

func (r *defaultGroovyScriptedErrorLogPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	groovyScriptedErrorLogPublisherSchema(ctx, req, resp, true)
}

func groovyScriptedErrorLogPublisherSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Groovy Scripted Error Log Publisher.",
		Attributes: map[string]schema.Attribute{
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Error Log Publisher.",
				Required:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Error Log Publisher. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"default_severity": schema.SetAttribute{
				Description: "Specifies the default severity levels for the logger.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"override_severity": schema.SetAttribute{
				Description: "Specifies the override severity levels for the logger based on the category of the messages.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log Publisher",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Log Publisher is enabled for use.",
				Required:    true,
			},
			"logging_error_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the server should exhibit if an error occurs during logging processing.",
				Optional:    true,
				Computed:    true,
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
func addOptionalGroovyScriptedErrorLogPublisherFields(ctx context.Context, addRequest *client.AddGroovyScriptedErrorLogPublisherRequest, plan groovyScriptedErrorLogPublisherResourceModel) error {
	if internaltypes.IsDefined(plan.ScriptArgument) {
		var slice []string
		plan.ScriptArgument.ElementsAs(ctx, &slice, false)
		addRequest.ScriptArgument = slice
	}
	if internaltypes.IsDefined(plan.DefaultSeverity) {
		var slice []string
		plan.DefaultSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherDefaultSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherDefaultSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefaultSeverity = enumSlice
	}
	if internaltypes.IsDefined(plan.OverrideSeverity) {
		var slice []string
		plan.OverrideSeverity.ElementsAs(ctx, &slice, false)
		addRequest.OverrideSeverity = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Read a GroovyScriptedErrorLogPublisherResponse object into the model struct
func readGroovyScriptedErrorLogPublisherResponse(ctx context.Context, r *client.GroovyScriptedErrorLogPublisherResponse, state *groovyScriptedErrorLogPublisherResourceModel, expectedValues *groovyScriptedErrorLogPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createGroovyScriptedErrorLogPublisherOperations(plan groovyScriptedErrorLogPublisherResourceModel, state groovyScriptedErrorLogPublisherResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ScriptClass, state.ScriptClass, "script-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScriptArgument, state.ScriptArgument, "script-argument")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DefaultSeverity, state.DefaultSeverity, "default-severity")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.OverrideSeverity, state.OverrideSeverity, "override-severity")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.LoggingErrorBehavior, state.LoggingErrorBehavior, "logging-error-behavior")
	return ops
}

// Create a new resource
func (r *groovyScriptedErrorLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan groovyScriptedErrorLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddGroovyScriptedErrorLogPublisherRequest(plan.Id.ValueString(),
		[]client.EnumgroovyScriptedErrorLogPublisherSchemaUrn{client.ENUMGROOVYSCRIPTEDERRORLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERGROOVY_SCRIPTED_ERROR},
		plan.ScriptClass.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalGroovyScriptedErrorLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Groovy Scripted Error Log Publisher", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddGroovyScriptedErrorLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Groovy Scripted Error Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state groovyScriptedErrorLogPublisherResourceModel
	readGroovyScriptedErrorLogPublisherResponse(ctx, addResponse.GroovyScriptedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultGroovyScriptedErrorLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan groovyScriptedErrorLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Groovy Scripted Error Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state groovyScriptedErrorLogPublisherResourceModel
	readGroovyScriptedErrorLogPublisherResponse(ctx, readResponse.GroovyScriptedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogPublisherApi.UpdateLogPublisher(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createGroovyScriptedErrorLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Groovy Scripted Error Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGroovyScriptedErrorLogPublisherResponse(ctx, updateResponse.GroovyScriptedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
func (r *groovyScriptedErrorLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readGroovyScriptedErrorLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultGroovyScriptedErrorLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readGroovyScriptedErrorLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readGroovyScriptedErrorLogPublisher(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state groovyScriptedErrorLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Groovy Scripted Error Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readGroovyScriptedErrorLogPublisherResponse(ctx, readResponse.GroovyScriptedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *groovyScriptedErrorLogPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateGroovyScriptedErrorLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultGroovyScriptedErrorLogPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateGroovyScriptedErrorLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateGroovyScriptedErrorLogPublisher(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan groovyScriptedErrorLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state groovyScriptedErrorLogPublisherResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogPublisherApi.UpdateLogPublisher(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createGroovyScriptedErrorLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Groovy Scripted Error Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGroovyScriptedErrorLogPublisherResponse(ctx, updateResponse.GroovyScriptedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultGroovyScriptedErrorLogPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *groovyScriptedErrorLogPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state groovyScriptedErrorLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogPublisherApi.DeleteLogPublisherExecute(r.apiClient.LogPublisherApi.DeleteLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Groovy Scripted Error Log Publisher", err, httpResp)
		return
	}
}

func (r *groovyScriptedErrorLogPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importGroovyScriptedErrorLogPublisher(ctx, req, resp)
}

func (r *defaultGroovyScriptedErrorLogPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importGroovyScriptedErrorLogPublisher(ctx, req, resp)
}

func importGroovyScriptedErrorLogPublisher(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
