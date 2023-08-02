package certificatemapper

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &certificateMappersDataSource{}
	_ datasource.DataSourceWithConfigure = &certificateMappersDataSource{}
)

// Create a Certificate Mappers data source
func NewCertificateMappersDataSource() datasource.DataSource {
	return &certificateMappersDataSource{}
}

// certificateMappersDataSource is the datasource implementation.
type certificateMappersDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *certificateMappersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate_mappers"
}

// Configure adds the provider configured client to the data source.
func (r *certificateMappersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type certificateMappersDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *certificateMappersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Certificate Mapper objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Certificate Mapper objects found in the configuration",
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
func (r *certificateMappersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state certificateMappersDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.CertificateMapperApi.ListCertificateMappers(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.CertificateMapperApi.ListCertificateMappersExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Certificate Mapper objects", err, httpResp)
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
		if response.SubjectEqualsDnCertificateMapperResponse != nil {
			attributes["id"] = types.StringValue(response.SubjectEqualsDnCertificateMapperResponse.Id)
			attributes["type"] = types.StringValue("subject-equals-dn")
		}
		if response.SubjectDnToUserAttributeCertificateMapperResponse != nil {
			attributes["id"] = types.StringValue(response.SubjectDnToUserAttributeCertificateMapperResponse.Id)
			attributes["type"] = types.StringValue("subject-dn-to-user-attribute")
		}
		if response.GroovyScriptedCertificateMapperResponse != nil {
			attributes["id"] = types.StringValue(response.GroovyScriptedCertificateMapperResponse.Id)
			attributes["type"] = types.StringValue("groovy-scripted")
		}
		if response.SubjectAttributeToUserAttributeCertificateMapperResponse != nil {
			attributes["id"] = types.StringValue(response.SubjectAttributeToUserAttributeCertificateMapperResponse.Id)
			attributes["type"] = types.StringValue("subject-attribute-to-user-attribute")
		}
		if response.FingerprintCertificateMapperResponse != nil {
			attributes["id"] = types.StringValue(response.FingerprintCertificateMapperResponse.Id)
			attributes["type"] = types.StringValue("fingerprint")
		}
		if response.ThirdPartyCertificateMapperResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyCertificateMapperResponse.Id)
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
