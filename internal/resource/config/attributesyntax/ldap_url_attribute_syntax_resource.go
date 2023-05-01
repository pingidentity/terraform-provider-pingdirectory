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
	   _ resource.Resource                = &ldapUrlAttributeSyntaxResource{}
	   _ resource.ResourceWithConfigure   = &ldapUrlAttributeSyntaxResource{}
	   _ resource.ResourceWithImportState = &ldapUrlAttributeSyntaxResource{}
)

// Create a Ldap Url Attribute Syntax resource
func NewLdapUrlAttributeSyntaxResource() resource.Resource {
	   return &ldapUrlAttributeSyntaxResource{}
}

// ldapUrlAttributeSyntaxResource is the resource implementation.
type ldapUrlAttributeSyntaxResource struct {
	   providerConfig internaltypes.ProviderConfiguration
	   apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *ldapUrlAttributeSyntaxResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	   resp.TypeName = req.ProviderTypeName + "_default_ldap_url_attribute_syntax"
}

// Configure adds the provider configured client to the resource.
func (r *ldapUrlAttributeSyntaxResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
    r.providerConfig = providerCfg.ProviderConfig
    r.apiClient = providerCfg.ApiClientV9200
}

type ldapUrlAttributeSyntaxResourceModel struct {
    Id              types.String `tfsdk:"id"`
    LastUpdated     types.String `tfsdk:"last_updated"`
    Notifications   types.Set    `tfsdk:"notifications"`
    RequiredActions types.Set    `tfsdk:"required_actions"`
    StrictFormat types.Bool `tfsdk:"strict_format"`
    Enabled types.Bool `tfsdk:"enabled"`
    RequireBinaryTransfer types.Bool `tfsdk:"require_binary_transfer"`
}

// GetSchema defines the schema for the resource.
func (r *ldapUrlAttributeSyntaxResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    schema := schema.Schema{
        Description: "Manages a Ldap Url Attribute Syntax.",
        Attributes: map[string]schema.Attribute{
            "strict_format": schema.BoolAttribute{
                Description: "Indicates whether values for attributes with this syntax will be required to be in the valid LDAP URL format. If this is set to false, then arbitrary strings will be allowed.",
                Optional:    true,
                Computed:    true,
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
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

// Read a LdapUrlAttributeSyntaxResponse object into the model struct
func readLdapUrlAttributeSyntaxResponse(ctx context.Context, r *client.LdapUrlAttributeSyntaxResponse, state *ldapUrlAttributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
    state.Id = types.StringValue(r.Id)
    state.StrictFormat = internaltypes.BoolTypeOrNil(r.StrictFormat)
    state.Enabled = types.BoolValue(r.Enabled)
    state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
    state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createLdapUrlAttributeSyntaxOperations(plan ldapUrlAttributeSyntaxResourceModel, state ldapUrlAttributeSyntaxResourceModel) []client.Operation {
    var ops []client.Operation
    operations.AddBoolOperationIfNecessary(&ops, plan.StrictFormat, state.StrictFormat, "strict-format")
    operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
    operations.AddBoolOperationIfNecessary(&ops, plan.RequireBinaryTransfer, state.RequireBinaryTransfer, "require-binary-transfer")
    return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *ldapUrlAttributeSyntaxResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Retrieve values from plan
    var plan ldapUrlAttributeSyntaxResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    readResponse, httpResp, err := r.apiClient.AttributeSyntaxApi.GetAttributeSyntax(
        config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
    if err != nil {
        config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ldap Url Attribute Syntax", err, httpResp)
        return
    }

    // Log response JSON
    responseJson, err := readResponse.MarshalJSON()
    if err == nil {
        tflog.Debug(ctx, "Read response: "+string(responseJson))
    }

    // Read the existing configuration
    var state ldapUrlAttributeSyntaxResourceModel
    readLdapUrlAttributeSyntaxResponse(ctx, readResponse.LdapUrlAttributeSyntaxResponse, &state, &resp.Diagnostics)

    // Determine what changes are needed to match the plan
    updateRequest := r.apiClient.AttributeSyntaxApi.UpdateAttributeSyntax(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	   ops := createLdapUrlAttributeSyntaxOperations(plan, state)
	   if len(ops) > 0 {
	   	   updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
	   	   // Log operations
	   	   operations.LogUpdateOperations(ctx, ops)

        updateResponse, httpResp, err := r.apiClient.AttributeSyntaxApi.UpdateAttributeSyntaxExecute(updateRequest)
        if err != nil {
            config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ldap Url Attribute Syntax", err, httpResp)
            return
        }

	       // Log response JSON
        responseJson, err := updateResponse.MarshalJSON()
        if err == nil {
        	   tflog.Debug(ctx, "Update response: "+string(responseJson))
        }

	       // Read the response
    readLdapUrlAttributeSyntaxResponse(ctx, updateResponse.LdapUrlAttributeSyntaxResponse, &state, &resp.Diagnostics)
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
func (r *ldapUrlAttributeSyntaxResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    // Get current state
  var state ldapUrlAttributeSyntaxResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    readResponse, httpResp, err := r.apiClient.AttributeSyntaxApi.GetAttributeSyntax(
        config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
    if err != nil {
        config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ldap Url Attribute Syntax", err, httpResp)
        return
    }

    // Log response JSON
    responseJson, err := readResponse.MarshalJSON()
    if err == nil {
        tflog.Debug(ctx, "Read response: "+string(responseJson))
    }

    // Read the response into the state
    readLdapUrlAttributeSyntaxResponse(ctx, readResponse.LdapUrlAttributeSyntaxResponse, &state, &resp.Diagnostics)

    // Set refreshed state
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Update a resource
func (r *ldapUrlAttributeSyntaxResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
  	 // Retrieve values from plan
	   var plan ldapUrlAttributeSyntaxResourceModel
	   diags := req.Plan.Get(ctx, &plan)
	   resp.Diagnostics.Append(diags...)
	   if resp.Diagnostics.HasError() {
	       return
	   }

	   // Get the current state to see how any attributes are changing
	   var state ldapUrlAttributeSyntaxResourceModel
	   req.State.Get(ctx, &state)
	   updateRequest := r.apiClient.AttributeSyntaxApi.UpdateAttributeSyntax(
        config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	   // Determine what update operations are necessary
	   ops := createLdapUrlAttributeSyntaxOperations(plan, state)
	   if len(ops) > 0 {
	   	   updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
	   	   // Log operations
	   	   operations.LogUpdateOperations(ctx, ops)

        updateResponse, httpResp, err := r.apiClient.AttributeSyntaxApi.UpdateAttributeSyntaxExecute(updateRequest)
        if err != nil {
            config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ldap Url Attribute Syntax", err, httpResp)
            return
        }

	       // Log response JSON
        responseJson, err := updateResponse.MarshalJSON()
        if err == nil {
        	   tflog.Debug(ctx, "Update response: "+string(responseJson))
        }

	       // Read the response
    readLdapUrlAttributeSyntaxResponse(ctx, updateResponse.LdapUrlAttributeSyntaxResponse, &state, &resp.Diagnostics)
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
func (r *ldapUrlAttributeSyntaxResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // No implementation necessary
}


func (r *ldapUrlAttributeSyntaxResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    // Retrieve import ID and save to id attribute
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

