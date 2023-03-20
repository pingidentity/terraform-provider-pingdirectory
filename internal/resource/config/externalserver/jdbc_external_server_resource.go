package externalserver

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &jdbcExternalServerResource{}
	_ resource.ResourceWithConfigure   = &jdbcExternalServerResource{}
	_ resource.ResourceWithImportState = &jdbcExternalServerResource{}
	_ resource.Resource                = &defaultJdbcExternalServerResource{}
	_ resource.ResourceWithConfigure   = &defaultJdbcExternalServerResource{}
	_ resource.ResourceWithImportState = &defaultJdbcExternalServerResource{}
)

// Create a Jdbc External Server resource
func NewJdbcExternalServerResource() resource.Resource {
	return &jdbcExternalServerResource{}
}

func NewDefaultJdbcExternalServerResource() resource.Resource {
	return &defaultJdbcExternalServerResource{}
}

// jdbcExternalServerResource is the resource implementation.
type jdbcExternalServerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultJdbcExternalServerResource is the resource implementation.
type defaultJdbcExternalServerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *jdbcExternalServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jdbc_external_server"
}

func (r *defaultJdbcExternalServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_jdbc_external_server"
}

// Configure adds the provider configured client to the resource.
func (r *jdbcExternalServerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultJdbcExternalServerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type jdbcExternalServerResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	LastUpdated               types.String `tfsdk:"last_updated"`
	Notifications             types.Set    `tfsdk:"notifications"`
	RequiredActions           types.Set    `tfsdk:"required_actions"`
	JdbcDriverType            types.String `tfsdk:"jdbc_driver_type"`
	JdbcDriverURL             types.String `tfsdk:"jdbc_driver_url"`
	DatabaseName              types.String `tfsdk:"database_name"`
	ServerHostName            types.String `tfsdk:"server_host_name"`
	ServerPort                types.Int64  `tfsdk:"server_port"`
	UserName                  types.String `tfsdk:"user_name"`
	Password                  types.String `tfsdk:"password"`
	PassphraseProvider        types.String `tfsdk:"passphrase_provider"`
	ValidationQuery           types.String `tfsdk:"validation_query"`
	ValidationQueryTimeout    types.String `tfsdk:"validation_query_timeout"`
	JdbcConnectionProperties  types.Set    `tfsdk:"jdbc_connection_properties"`
	TransactionIsolationLevel types.String `tfsdk:"transaction_isolation_level"`
	Description               types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *jdbcExternalServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	jdbcExternalServerSchema(ctx, req, resp, false)
}

func (r *defaultJdbcExternalServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	jdbcExternalServerSchema(ctx, req, resp, true)
}

func jdbcExternalServerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Jdbc External Server.",
		Attributes: map[string]schema.Attribute{
			"jdbc_driver_type": schema.StringAttribute{
				Description: "Specifies a supported database driver type. The driver class will be automatically selected based on this selection. We highly recommend using a JDBC 4 driver that is suitable for the current Java platform.",
				Required:    true,
			},
			"jdbc_driver_url": schema.StringAttribute{
				Description: "Specify the complete JDBC URL which will be used instead of the automatic URL format. You must select type 'other' for the jdbc-driver-type.",
				Optional:    true,
			},
			"database_name": schema.StringAttribute{
				Description: "Specifies which database to connect to. This is ignored if jdbc-driver-url is specified.",
				Optional:    true,
			},
			"server_host_name": schema.StringAttribute{
				Description: "The host name of the database server. This is ignored if jdbc-driver-url is specified.",
				Optional:    true,
			},
			"server_port": schema.Int64Attribute{
				Description: "The port number where the database server listens for requests. This is ignored if jdbc-driver-url is specified",
				Optional:    true,
			},
			"user_name": schema.StringAttribute{
				Description: "The name of the login account to use when connecting to the database server.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "The login password for the specified user name.",
				Optional:    true,
				Sensitive:   true,
			},
			"passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider to use to obtain the login password for the specified user.",
				Optional:    true,
			},
			"validation_query": schema.StringAttribute{
				Description: "The SQL query that will be used to validate connections to the database before making them available to the Directory Server.",
				Optional:    true,
			},
			"validation_query_timeout": schema.StringAttribute{
				Description: "Specifies the amount of time to wait for a response from the database when executing the validation query, if one is set. If the timeout is exceeded, the Directory Server will drop the connection and obtain a new one. A value of zero indicates no timeout.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"jdbc_connection_properties": schema.SetAttribute{
				Description: "Specifies the connection properties for the JDBC datasource.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"transaction_isolation_level": schema.StringAttribute{
				Description: "This property specifies the default transaction isolation level for connections to this JDBC External Server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this External Server",
				Optional:    true,
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
func addOptionalJdbcExternalServerFields(ctx context.Context, addRequest *client.AddJdbcExternalServerRequest, plan jdbcExternalServerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.JdbcDriverURL) {
		stringVal := plan.JdbcDriverURL.ValueString()
		addRequest.JdbcDriverURL = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DatabaseName) {
		stringVal := plan.DatabaseName.ValueString()
		addRequest.DatabaseName = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ServerHostName) {
		stringVal := plan.ServerHostName.ValueString()
		addRequest.ServerHostName = &stringVal
	}
	if internaltypes.IsDefined(plan.ServerPort) {
		intVal := int32(plan.ServerPort.ValueInt64())
		addRequest.ServerPort = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UserName) {
		stringVal := plan.UserName.ValueString()
		addRequest.UserName = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Password) {
		stringVal := plan.Password.ValueString()
		addRequest.Password = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PassphraseProvider) {
		stringVal := plan.PassphraseProvider.ValueString()
		addRequest.PassphraseProvider = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidationQuery) {
		stringVal := plan.ValidationQuery.ValueString()
		addRequest.ValidationQuery = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidationQueryTimeout) {
		stringVal := plan.ValidationQueryTimeout.ValueString()
		addRequest.ValidationQueryTimeout = &stringVal
	}
	if internaltypes.IsDefined(plan.JdbcConnectionProperties) {
		var slice []string
		plan.JdbcConnectionProperties.ElementsAs(ctx, &slice, false)
		addRequest.JdbcConnectionProperties = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TransactionIsolationLevel) {
		transactionIsolationLevel, err := client.NewEnumexternalServerTransactionIsolationLevelPropFromValue(plan.TransactionIsolationLevel.ValueString())
		if err != nil {
			return err
		}
		addRequest.TransactionIsolationLevel = transactionIsolationLevel
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	return nil
}

// Read a JdbcExternalServerResponse object into the model struct
func readJdbcExternalServerResponse(ctx context.Context, r *client.JdbcExternalServerResponse, state *jdbcExternalServerResourceModel, expectedValues *jdbcExternalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.JdbcDriverType = types.StringValue(r.JdbcDriverType.String())
	state.JdbcDriverURL = internaltypes.StringTypeOrNil(r.JdbcDriverURL, internaltypes.IsEmptyString(expectedValues.JdbcDriverURL))
	state.DatabaseName = internaltypes.StringTypeOrNil(r.DatabaseName, internaltypes.IsEmptyString(expectedValues.DatabaseName))
	state.ServerHostName = internaltypes.StringTypeOrNil(r.ServerHostName, internaltypes.IsEmptyString(expectedValues.ServerHostName))
	state.ServerPort = internaltypes.Int64TypeOrNil(r.ServerPort)
	state.UserName = internaltypes.StringTypeOrNil(r.UserName, internaltypes.IsEmptyString(expectedValues.UserName))
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.Password = expectedValues.Password
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, internaltypes.IsEmptyString(expectedValues.PassphraseProvider))
	state.ValidationQuery = internaltypes.StringTypeOrNil(r.ValidationQuery, internaltypes.IsEmptyString(expectedValues.ValidationQuery))
	state.ValidationQueryTimeout = internaltypes.StringTypeOrNil(r.ValidationQueryTimeout, internaltypes.IsEmptyString(expectedValues.ValidationQueryTimeout))
	config.CheckMismatchedPDFormattedAttributes("validation_query_timeout",
		expectedValues.ValidationQueryTimeout, state.ValidationQueryTimeout, diagnostics)
	state.JdbcConnectionProperties = internaltypes.GetStringSet(r.JdbcConnectionProperties)
	state.TransactionIsolationLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumexternalServerTransactionIsolationLevelProp(r.TransactionIsolationLevel), internaltypes.IsEmptyString(expectedValues.TransactionIsolationLevel))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createJdbcExternalServerOperations(plan jdbcExternalServerResourceModel, state jdbcExternalServerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.JdbcDriverType, state.JdbcDriverType, "jdbc-driver-type")
	operations.AddStringOperationIfNecessary(&ops, plan.JdbcDriverURL, state.JdbcDriverURL, "jdbc-driver-url")
	operations.AddStringOperationIfNecessary(&ops, plan.DatabaseName, state.DatabaseName, "database-name")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerHostName, state.ServerHostName, "server-host-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.ServerPort, state.ServerPort, "server-port")
	operations.AddStringOperationIfNecessary(&ops, plan.UserName, state.UserName, "user-name")
	operations.AddStringOperationIfNecessary(&ops, plan.Password, state.Password, "password")
	operations.AddStringOperationIfNecessary(&ops, plan.PassphraseProvider, state.PassphraseProvider, "passphrase-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.ValidationQuery, state.ValidationQuery, "validation-query")
	operations.AddStringOperationIfNecessary(&ops, plan.ValidationQueryTimeout, state.ValidationQueryTimeout, "validation-query-timeout")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.JdbcConnectionProperties, state.JdbcConnectionProperties, "jdbc-connection-properties")
	operations.AddStringOperationIfNecessary(&ops, plan.TransactionIsolationLevel, state.TransactionIsolationLevel, "transaction-isolation-level")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *jdbcExternalServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan jdbcExternalServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	jdbcDriverType, err := client.NewEnumexternalServerJdbcDriverTypePropFromValue(plan.JdbcDriverType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for JdbcDriverType", err.Error())
		return
	}
	addRequest := client.NewAddJdbcExternalServerRequest(plan.Id.ValueString(),
		[]client.EnumjdbcExternalServerSchemaUrn{client.ENUMJDBCEXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERJDBC},
		*jdbcDriverType)
	err = addOptionalJdbcExternalServerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Jdbc External Server", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExternalServerApi.AddExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExternalServerRequest(
		client.AddJdbcExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerApi.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Jdbc External Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state jdbcExternalServerResourceModel
	readJdbcExternalServerResponse(ctx, addResponse.JdbcExternalServerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultJdbcExternalServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan jdbcExternalServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ExternalServerApi.GetExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Jdbc External Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state jdbcExternalServerResourceModel
	readJdbcExternalServerResponse(ctx, readResponse.JdbcExternalServerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ExternalServerApi.UpdateExternalServer(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createJdbcExternalServerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ExternalServerApi.UpdateExternalServerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Jdbc External Server", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readJdbcExternalServerResponse(ctx, updateResponse.JdbcExternalServerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *jdbcExternalServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readJdbcExternalServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultJdbcExternalServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readJdbcExternalServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readJdbcExternalServer(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state jdbcExternalServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ExternalServerApi.GetExternalServer(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Jdbc External Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readJdbcExternalServerResponse(ctx, readResponse.JdbcExternalServerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *jdbcExternalServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateJdbcExternalServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultJdbcExternalServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateJdbcExternalServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateJdbcExternalServer(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan jdbcExternalServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state jdbcExternalServerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ExternalServerApi.UpdateExternalServer(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createJdbcExternalServerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ExternalServerApi.UpdateExternalServerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Jdbc External Server", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readJdbcExternalServerResponse(ctx, updateResponse.JdbcExternalServerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultJdbcExternalServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *jdbcExternalServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state jdbcExternalServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ExternalServerApi.DeleteExternalServerExecute(r.apiClient.ExternalServerApi.DeleteExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Jdbc External Server", err, httpResp)
		return
	}
}

func (r *jdbcExternalServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importJdbcExternalServer(ctx, req, resp)
}

func (r *defaultJdbcExternalServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importJdbcExternalServer(ctx, req, resp)
}

func importJdbcExternalServer(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
