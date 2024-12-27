package logfieldmapping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &logFieldMappingDataSource{}
	_ datasource.DataSourceWithConfigure = &logFieldMappingDataSource{}
)

// Create a Log Field Mapping data source
func NewLogFieldMappingDataSource() datasource.DataSource {
	return &logFieldMappingDataSource{}
}

// logFieldMappingDataSource is the datasource implementation.
type logFieldMappingDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *logFieldMappingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_field_mapping"
}

// Configure adds the provider configured client to the data source.
func (r *logFieldMappingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type logFieldMappingDataSourceModel struct {
	Id                                  types.String `tfsdk:"id"`
	Name                                types.String `tfsdk:"name"`
	Type                                types.String `tfsdk:"type"`
	LogFieldTimestamp                   types.String `tfsdk:"log_field_timestamp"`
	LogFieldConnectionID                types.String `tfsdk:"log_field_connection_id"`
	LogFieldStartupid                   types.String `tfsdk:"log_field_startupid"`
	LogFieldProductName                 types.String `tfsdk:"log_field_product_name"`
	LogFieldCategory                    types.String `tfsdk:"log_field_category"`
	LogFieldSeverity                    types.String `tfsdk:"log_field_severity"`
	LogFieldInstanceName                types.String `tfsdk:"log_field_instance_name"`
	LogFieldOperationID                 types.String `tfsdk:"log_field_operation_id"`
	LogFieldMessageType                 types.String `tfsdk:"log_field_message_type"`
	LogFieldOperationType               types.String `tfsdk:"log_field_operation_type"`
	LogFieldMessageID                   types.String `tfsdk:"log_field_message_id"`
	LogFieldResultCode                  types.String `tfsdk:"log_field_result_code"`
	LogFieldMessage                     types.String `tfsdk:"log_field_message"`
	LogFieldOrigin                      types.String `tfsdk:"log_field_origin"`
	LogFieldRequesterDN                 types.String `tfsdk:"log_field_requester_dn"`
	LogFieldDisconnectReason            types.String `tfsdk:"log_field_disconnect_reason"`
	LogFieldDeleteOldRDN                types.String `tfsdk:"log_field_delete_old_rdn"`
	LogFieldAuthenticatedUserDN         types.String `tfsdk:"log_field_authenticated_user_dn"`
	LogFieldProcessingTime              types.String `tfsdk:"log_field_processing_time"`
	LogFieldRequestedAttributes         types.String `tfsdk:"log_field_requested_attributes"`
	LogFieldSASLMechanismName           types.String `tfsdk:"log_field_sasl_mechanism_name"`
	LogFieldNewRDN                      types.String `tfsdk:"log_field_new_rdn"`
	LogFieldBaseDN                      types.String `tfsdk:"log_field_base_dn"`
	LogFieldBindDN                      types.String `tfsdk:"log_field_bind_dn"`
	LogFieldMatchedDN                   types.String `tfsdk:"log_field_matched_dn"`
	LogFieldRequesterIPAddress          types.String `tfsdk:"log_field_requester_ip_address"`
	LogFieldAuthenticationType          types.String `tfsdk:"log_field_authentication_type"`
	LogFieldNewSuperiorDN               types.String `tfsdk:"log_field_new_superior_dn"`
	LogFieldFilter                      types.String `tfsdk:"log_field_filter"`
	LogFieldAlternateAuthorizationDN    types.String `tfsdk:"log_field_alternate_authorization_dn"`
	LogFieldEntryDN                     types.String `tfsdk:"log_field_entry_dn"`
	LogFieldEntriesReturned             types.String `tfsdk:"log_field_entries_returned"`
	LogFieldAuthenticationFailureID     types.String `tfsdk:"log_field_authentication_failure_id"`
	LogFieldRequestOID                  types.String `tfsdk:"log_field_request_oid"`
	LogFieldResponseOID                 types.String `tfsdk:"log_field_response_oid"`
	LogFieldTargetProtocol              types.String `tfsdk:"log_field_target_protocol"`
	LogFieldTargetPort                  types.String `tfsdk:"log_field_target_port"`
	LogFieldTargetAddress               types.String `tfsdk:"log_field_target_address"`
	LogFieldTargetAttribute             types.String `tfsdk:"log_field_target_attribute"`
	LogFieldTargetHost                  types.String `tfsdk:"log_field_target_host"`
	LogFieldProtocolVersion             types.String `tfsdk:"log_field_protocol_version"`
	LogFieldProtocolName                types.String `tfsdk:"log_field_protocol_name"`
	LogFieldAuthenticationFailureReason types.String `tfsdk:"log_field_authentication_failure_reason"`
	LogFieldAdditionalInformation       types.String `tfsdk:"log_field_additional_information"`
	LogFieldUnindexed                   types.String `tfsdk:"log_field_unindexed"`
	LogFieldScope                       types.String `tfsdk:"log_field_scope"`
	LogFieldReferralUrls                types.String `tfsdk:"log_field_referral_urls"`
	LogFieldSourceAddress               types.String `tfsdk:"log_field_source_address"`
	LogFieldMessageIDToAbandon          types.String `tfsdk:"log_field_message_id_to_abandon"`
	LogFieldResponseControls            types.String `tfsdk:"log_field_response_controls"`
	LogFieldRequestControls             types.String `tfsdk:"log_field_request_controls"`
	LogFieldIntermediateClientResult    types.String `tfsdk:"log_field_intermediate_client_result"`
	LogFieldIntermediateClientRequest   types.String `tfsdk:"log_field_intermediate_client_request"`
	LogFieldReplicationChangeID         types.String `tfsdk:"log_field_replication_change_id"`
	Description                         types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the datasource.
func (r *logFieldMappingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Log Field Mapping.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Log Field Mapping resource. Options are ['access', 'error']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_timestamp": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `access`: The time that the operation was processed. When the `type` attribute is set to `error`: The time that the log message was generated.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `access`: The time that the operation was processed.\n  - `error`: The time that the log message was generated.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"log_field_connection_id": schema.StringAttribute{
				Description: "The connection ID assigned to the client connection.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_startupid": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `access`: The startup ID for the Directory Server. A different value will be generated each time the server is started, and it may be used to distinguish between operations with the same connection ID and operation ID across server restarts. When the `type` attribute is set to `error`: The startup ID for the Directory Server. A different value will be generated each time the server is started.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `access`: The startup ID for the Directory Server. A different value will be generated each time the server is started, and it may be used to distinguish between operations with the same connection ID and operation ID across server restarts.\n  - `error`: The startup ID for the Directory Server. A different value will be generated each time the server is started.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"log_field_product_name": schema.StringAttribute{
				Description: "The name for this Directory Server product, which may be used to identify which product was used to log the message if multiple products log to the same database table.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_category": schema.StringAttribute{
				Description: "The category for the log message.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_severity": schema.StringAttribute{
				Description: "The severity for the log message.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_instance_name": schema.StringAttribute{
				Description: "A name that uniquely identifies this Directory Server instance, which may be used to identify which instance was used to log the message if multiple server instances log to the same database table.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_operation_id": schema.StringAttribute{
				Description: "The operation ID for the operation processed by the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_message_type": schema.StringAttribute{
				Description: "The type of log message. Message types may include \"CONNECT\", \"DISCONNECT\", \"FORWARD\", \"RESULT\", \"ENTRY\", or \"REFERENCE\".",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_operation_type": schema.StringAttribute{
				Description: "The type of operation that was processed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_message_id": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `access`: The message ID included in the client request. When the `type` attribute is set to `error`: The numeric value which uniquely identifies the type of message.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `access`: The message ID included in the client request.\n  - `error`: The numeric value which uniquely identifies the type of message.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"log_field_result_code": schema.StringAttribute{
				Description: "The numeric result code for the operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_message": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `access`: The diagnostic message for the operation. When the `type` attribute is set to `error`: The text of the log message.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `access`: The diagnostic message for the operation.\n  - `error`: The text of the log message.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"log_field_origin": schema.StringAttribute{
				Description: "The origin for the operation. Values may include \"replication\" (if the operation was received via replication), \"internal\" (if it was an internal operation processed by a server component), or \"external\" (if it was a request from a client).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_requester_dn": schema.StringAttribute{
				Description: "The DN of the user that requested the operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_disconnect_reason": schema.StringAttribute{
				Description: "The reason that the client connection was closed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_delete_old_rdn": schema.StringAttribute{
				Description: "Indicates whether the old RDN values should be removed from an entry while processing a modify DN operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_authenticated_user_dn": schema.StringAttribute{
				Description: "The DN of the user that authenticated to the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_processing_time": schema.StringAttribute{
				Description: "The length of time (in milliseconds with microsecond accuracy) required to process the operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_requested_attributes": schema.StringAttribute{
				Description: "The set of requested attributes for the search operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_sasl_mechanism_name": schema.StringAttribute{
				Description: "The name of the SASL mechanism used to authenticate.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_new_rdn": schema.StringAttribute{
				Description: "The new RDN to use for a modify DN operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_base_dn": schema.StringAttribute{
				Description: "The base DN for a search operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_bind_dn": schema.StringAttribute{
				Description: "The bind DN for a bind operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_matched_dn": schema.StringAttribute{
				Description: "The DN of the superior entry closest to the DN specified by the client.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_requester_ip_address": schema.StringAttribute{
				Description: "The IP address of the client that requested the operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_authentication_type": schema.StringAttribute{
				Description: "The type of authentication requested by the client.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_new_superior_dn": schema.StringAttribute{
				Description: "The new superior DN from a modify DN operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_filter": schema.StringAttribute{
				Description: "The filter from a search operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_alternate_authorization_dn": schema.StringAttribute{
				Description: "The DN of the alternate authorization identity used when processing the operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_entry_dn": schema.StringAttribute{
				Description: "The DN of the entry targeted by the operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_entries_returned": schema.StringAttribute{
				Description: "The number of search result entries returned to the client.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_authentication_failure_id": schema.StringAttribute{
				Description: "The numeric identifier for the authentication failure reason.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_request_oid": schema.StringAttribute{
				Description: "The OID of an extended request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_response_oid": schema.StringAttribute{
				Description: "The OID of an extended response.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_target_protocol": schema.StringAttribute{
				Description: "The protocol used when forwarding the request to a backend server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_target_port": schema.StringAttribute{
				Description: "The network port of the Directory Server to which the client connection has been established, or of the backend server to which the request has been forwarded.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_target_address": schema.StringAttribute{
				Description: "The network address of the Directory Server to which the client connection has been established.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_target_attribute": schema.StringAttribute{
				Description: "The name of the attribute targeted by a compare operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_target_host": schema.StringAttribute{
				Description: "The address of the server to which the request has been forwarded.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_protocol_version": schema.StringAttribute{
				Description: "The protocol version used by the client when communicating with the Directory Server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_protocol_name": schema.StringAttribute{
				Description: "The name of the protocol the client is using to communicate with the Directory Server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_authentication_failure_reason": schema.StringAttribute{
				Description: "A message explaining the reason that the authentication attempt failed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_additional_information": schema.StringAttribute{
				Description: "Additional information about the operation that was processed which was not returned to the client.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_unindexed": schema.StringAttribute{
				Description: "Indicates whether the requested search operation was unindexed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_scope": schema.StringAttribute{
				Description: "The scope for the search operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_referral_urls": schema.StringAttribute{
				Description: "The referral URLs returned to the client.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_source_address": schema.StringAttribute{
				Description: "The address of the client from which the connection was established.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_message_id_to_abandon": schema.StringAttribute{
				Description: "The message ID of the operation to be abandoned.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_response_controls": schema.StringAttribute{
				Description: "The OIDs of the response controls returned to the client.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_request_controls": schema.StringAttribute{
				Description: "The OIDs of the request controls returned to the client.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_intermediate_client_result": schema.StringAttribute{
				Description: "The contents of the intermediate client response control returned to the client.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_intermediate_client_request": schema.StringAttribute{
				Description: "The contents of the intermediate client request control provided by the client.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_replication_change_id": schema.StringAttribute{
				Description: "The replication change ID.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log Field Mapping",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a AccessLogFieldMappingResponse object into the model struct
func readAccessLogFieldMappingResponseDataSource(ctx context.Context, r *client.AccessLogFieldMappingResponse, state *logFieldMappingDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFieldTimestamp = internaltypes.StringTypeOrNil(r.LogFieldTimestamp, false)
	state.LogFieldConnectionID = internaltypes.StringTypeOrNil(r.LogFieldConnectionID, false)
	state.LogFieldStartupid = internaltypes.StringTypeOrNil(r.LogFieldStartupid, false)
	state.LogFieldProductName = internaltypes.StringTypeOrNil(r.LogFieldProductName, false)
	state.LogFieldInstanceName = internaltypes.StringTypeOrNil(r.LogFieldInstanceName, false)
	state.LogFieldOperationID = internaltypes.StringTypeOrNil(r.LogFieldOperationID, false)
	state.LogFieldMessageType = internaltypes.StringTypeOrNil(r.LogFieldMessageType, false)
	state.LogFieldOperationType = internaltypes.StringTypeOrNil(r.LogFieldOperationType, false)
	state.LogFieldMessageID = internaltypes.StringTypeOrNil(r.LogFieldMessageID, false)
	state.LogFieldResultCode = internaltypes.StringTypeOrNil(r.LogFieldResultCode, false)
	state.LogFieldMessage = internaltypes.StringTypeOrNil(r.LogFieldMessage, false)
	state.LogFieldOrigin = internaltypes.StringTypeOrNil(r.LogFieldOrigin, false)
	state.LogFieldRequesterDN = internaltypes.StringTypeOrNil(r.LogFieldRequesterDN, false)
	state.LogFieldDisconnectReason = internaltypes.StringTypeOrNil(r.LogFieldDisconnectReason, false)
	state.LogFieldDeleteOldRDN = internaltypes.StringTypeOrNil(r.LogFieldDeleteOldRDN, false)
	state.LogFieldAuthenticatedUserDN = internaltypes.StringTypeOrNil(r.LogFieldAuthenticatedUserDN, false)
	state.LogFieldProcessingTime = internaltypes.StringTypeOrNil(r.LogFieldProcessingTime, false)
	state.LogFieldRequestedAttributes = internaltypes.StringTypeOrNil(r.LogFieldRequestedAttributes, false)
	state.LogFieldSASLMechanismName = internaltypes.StringTypeOrNil(r.LogFieldSASLMechanismName, false)
	state.LogFieldNewRDN = internaltypes.StringTypeOrNil(r.LogFieldNewRDN, false)
	state.LogFieldBaseDN = internaltypes.StringTypeOrNil(r.LogFieldBaseDN, false)
	state.LogFieldBindDN = internaltypes.StringTypeOrNil(r.LogFieldBindDN, false)
	state.LogFieldMatchedDN = internaltypes.StringTypeOrNil(r.LogFieldMatchedDN, false)
	state.LogFieldRequesterIPAddress = internaltypes.StringTypeOrNil(r.LogFieldRequesterIPAddress, false)
	state.LogFieldAuthenticationType = internaltypes.StringTypeOrNil(r.LogFieldAuthenticationType, false)
	state.LogFieldNewSuperiorDN = internaltypes.StringTypeOrNil(r.LogFieldNewSuperiorDN, false)
	state.LogFieldFilter = internaltypes.StringTypeOrNil(r.LogFieldFilter, false)
	state.LogFieldAlternateAuthorizationDN = internaltypes.StringTypeOrNil(r.LogFieldAlternateAuthorizationDN, false)
	state.LogFieldEntryDN = internaltypes.StringTypeOrNil(r.LogFieldEntryDN, false)
	state.LogFieldEntriesReturned = internaltypes.StringTypeOrNil(r.LogFieldEntriesReturned, false)
	state.LogFieldAuthenticationFailureID = internaltypes.StringTypeOrNil(r.LogFieldAuthenticationFailureID, false)
	state.LogFieldRequestOID = internaltypes.StringTypeOrNil(r.LogFieldRequestOID, false)
	state.LogFieldResponseOID = internaltypes.StringTypeOrNil(r.LogFieldResponseOID, false)
	state.LogFieldTargetProtocol = internaltypes.StringTypeOrNil(r.LogFieldTargetProtocol, false)
	state.LogFieldTargetPort = internaltypes.StringTypeOrNil(r.LogFieldTargetPort, false)
	state.LogFieldTargetAddress = internaltypes.StringTypeOrNil(r.LogFieldTargetAddress, false)
	state.LogFieldTargetAttribute = internaltypes.StringTypeOrNil(r.LogFieldTargetAttribute, false)
	state.LogFieldTargetHost = internaltypes.StringTypeOrNil(r.LogFieldTargetHost, false)
	state.LogFieldProtocolVersion = internaltypes.StringTypeOrNil(r.LogFieldProtocolVersion, false)
	state.LogFieldProtocolName = internaltypes.StringTypeOrNil(r.LogFieldProtocolName, false)
	state.LogFieldAuthenticationFailureReason = internaltypes.StringTypeOrNil(r.LogFieldAuthenticationFailureReason, false)
	state.LogFieldAdditionalInformation = internaltypes.StringTypeOrNil(r.LogFieldAdditionalInformation, false)
	state.LogFieldUnindexed = internaltypes.StringTypeOrNil(r.LogFieldUnindexed, false)
	state.LogFieldScope = internaltypes.StringTypeOrNil(r.LogFieldScope, false)
	state.LogFieldReferralUrls = internaltypes.StringTypeOrNil(r.LogFieldReferralUrls, false)
	state.LogFieldSourceAddress = internaltypes.StringTypeOrNil(r.LogFieldSourceAddress, false)
	state.LogFieldMessageIDToAbandon = internaltypes.StringTypeOrNil(r.LogFieldMessageIDToAbandon, false)
	state.LogFieldResponseControls = internaltypes.StringTypeOrNil(r.LogFieldResponseControls, false)
	state.LogFieldRequestControls = internaltypes.StringTypeOrNil(r.LogFieldRequestControls, false)
	state.LogFieldIntermediateClientResult = internaltypes.StringTypeOrNil(r.LogFieldIntermediateClientResult, false)
	state.LogFieldIntermediateClientRequest = internaltypes.StringTypeOrNil(r.LogFieldIntermediateClientRequest, false)
	state.LogFieldReplicationChangeID = internaltypes.StringTypeOrNil(r.LogFieldReplicationChangeID, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a ErrorLogFieldMappingResponse object into the model struct
func readErrorLogFieldMappingResponseDataSource(ctx context.Context, r *client.ErrorLogFieldMappingResponse, state *logFieldMappingDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("error")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFieldTimestamp = internaltypes.StringTypeOrNil(r.LogFieldTimestamp, false)
	state.LogFieldProductName = internaltypes.StringTypeOrNil(r.LogFieldProductName, false)
	state.LogFieldInstanceName = internaltypes.StringTypeOrNil(r.LogFieldInstanceName, false)
	state.LogFieldStartupid = internaltypes.StringTypeOrNil(r.LogFieldStartupid, false)
	state.LogFieldCategory = internaltypes.StringTypeOrNil(r.LogFieldCategory, false)
	state.LogFieldSeverity = internaltypes.StringTypeOrNil(r.LogFieldSeverity, false)
	state.LogFieldMessageID = internaltypes.StringTypeOrNil(r.LogFieldMessageID, false)
	state.LogFieldMessage = internaltypes.StringTypeOrNil(r.LogFieldMessage, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read resource information
func (r *logFieldMappingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state logFieldMappingDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogFieldMappingAPI.GetLogFieldMapping(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Log Field Mapping", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.AccessLogFieldMappingResponse != nil {
		readAccessLogFieldMappingResponseDataSource(ctx, readResponse.AccessLogFieldMappingResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ErrorLogFieldMappingResponse != nil {
		readErrorLogFieldMappingResponseDataSource(ctx, readResponse.ErrorLogFieldMappingResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
