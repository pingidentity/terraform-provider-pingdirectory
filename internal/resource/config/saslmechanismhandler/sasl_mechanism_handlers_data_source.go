package saslmechanismhandler

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &saslMechanismHandlersDataSource{}
	_ datasource.DataSourceWithConfigure = &saslMechanismHandlersDataSource{}
)

// Create a Sasl Mechanism Handlers data source
func NewSaslMechanismHandlersDataSource() datasource.DataSource {
	return &saslMechanismHandlersDataSource{}
}

// saslMechanismHandlersDataSource is the datasource implementation.
type saslMechanismHandlersDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *saslMechanismHandlersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sasl_mechanism_handlers"
}

// Configure adds the provider configured client to the data source.
func (r *saslMechanismHandlersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type saslMechanismHandlersDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *saslMechanismHandlersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Sasl Mechanism Handler objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Sasl Mechanism Handler objects found in the configuration",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: internaltypes.ObjectsObjectType(),
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read resource information
func (r *saslMechanismHandlersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state saslMechanismHandlersDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.SaslMechanismHandlerAPI.ListSaslMechanismHandlers(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.SaslMechanismHandlerAPI.ListSaslMechanismHandlersExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Sasl Mechanism Handler objects", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	objects := []attr.Value{}
	for _, response := range readResponse.Resources {
		attributes := map[string]attr.Value{}
		if response.UnboundidMsChapV2SaslMechanismHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.UnboundidMsChapV2SaslMechanismHandlerResponse.Id)
			attributes["type"] = types.StringValue("unboundid-ms-chap-v2")
		}
		if response.UnboundidTotpSaslMechanismHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.UnboundidTotpSaslMechanismHandlerResponse.Id)
			attributes["type"] = types.StringValue("unboundid-totp")
		}
		if response.UnboundidInterServerSaslMechanismHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.UnboundidInterServerSaslMechanismHandlerResponse.Id)
			attributes["type"] = types.StringValue("unboundid-inter-server")
		}
		if response.PingIdentityInterServerSaslMechanismHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.PingIdentityInterServerSaslMechanismHandlerResponse.Id)
			attributes["type"] = types.StringValue("ping-identity-inter-server")
		}
		if response.UnboundidYubikeyOtpSaslMechanismHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.UnboundidYubikeyOtpSaslMechanismHandlerResponse.Id)
			attributes["type"] = types.StringValue("unboundid-yubikey-otp")
		}
		if response.ExternalSaslMechanismHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.ExternalSaslMechanismHandlerResponse.Id)
			attributes["type"] = types.StringValue("external")
		}
		if response.DigestMd5SaslMechanismHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.DigestMd5SaslMechanismHandlerResponse.Id)
			attributes["type"] = types.StringValue("digest-md5")
		}
		if response.PlainSaslMechanismHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.PlainSaslMechanismHandlerResponse.Id)
			attributes["type"] = types.StringValue("plain")
		}
		if response.UnboundidDeliveredOtpSaslMechanismHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.UnboundidDeliveredOtpSaslMechanismHandlerResponse.Id)
			attributes["type"] = types.StringValue("unboundid-delivered-otp")
		}
		if response.UnboundidExternalAuthSaslMechanismHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.UnboundidExternalAuthSaslMechanismHandlerResponse.Id)
			attributes["type"] = types.StringValue("unboundid-external-auth")
		}
		if response.AnonymousSaslMechanismHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.AnonymousSaslMechanismHandlerResponse.Id)
			attributes["type"] = types.StringValue("anonymous")
		}
		if response.CramMd5SaslMechanismHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.CramMd5SaslMechanismHandlerResponse.Id)
			attributes["type"] = types.StringValue("cram-md5")
		}
		if response.OauthBearerSaslMechanismHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.OauthBearerSaslMechanismHandlerResponse.Id)
			attributes["type"] = types.StringValue("oauth-bearer")
		}
		if response.UnboundidCertificatePlusPasswordSaslMechanismHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.UnboundidCertificatePlusPasswordSaslMechanismHandlerResponse.Id)
			attributes["type"] = types.StringValue("unboundid-certificate-plus-password")
		}
		if response.GssapiSaslMechanismHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.GssapiSaslMechanismHandlerResponse.Id)
			attributes["type"] = types.StringValue("gssapi")
		}
		if response.ThirdPartySaslMechanismHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartySaslMechanismHandlerResponse.Id)
			attributes["type"] = types.StringValue("third-party")
		}
		obj, diags := types.ObjectValue(internaltypes.ObjectsAttrTypes(), attributes)
		resp.Diagnostics.Append(diags...)
		objects = append(objects, obj)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	state.Objects, diags = types.SetValue(internaltypes.ObjectsObjectType(), objects)
	resp.Diagnostics.Append(diags...)
	state.Id = types.StringValue("id")

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
