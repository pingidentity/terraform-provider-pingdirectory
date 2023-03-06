#
# Configure pf-connected-identities for DA configuration
#
resource "pingdirectory_composed_attribute_plugin" "pfConnectedIdentitiesPlugin" {
  id                                                         = "pf-connected-identities"
  enabled                                                    = true
  attribute_type                                             = "objectClass"
  value_pattern                                              = ["pf-connected-identities"]
  target_attribute_exists_during_initial_population_behavior = "merge-existing-and-composed-values"
  include_base_dn                                            = ["${var.user_base_dn}"]
  include_filter                                             = ["(objectClass=inetOrgPerson)"]
}

resource "pingdirectory_composed_attribute_plugin" "pfConnectedIdentityPlugin" {
  id              = "pf-connected-identity"
  enabled         = true
  attribute_type  = "pf-connected-identity"
  value_pattern   = ["auth-source=pf-local-identity:user-id={uid}"]
  include_base_dn = ["${var.user_base_dn}"]
  include_filter  = ["(objectClass=inetOrgPerson)"]
}

#
# The search-base-dn value is the DN of a valid base entry where
# managed users are stored.
#
resource "pingdirectory_user_rest_resource_type" "usersRestResourceType" {
  id                             = "users"
  display_name                   = "Users"
  enabled                        = true
  search_base_dn                 = "ou=people,${var.user_base_dn}"
  primary_display_attribute_type = "cn"
  resource_endpoint              = "users"
  search_filter_pattern          = "(|(cn=*%%*)(mail=%%*)(uid=%%*)(sn=*%%*))"
  structural_ldap_objectclass    = "inetOrgPerson"
  parent_dn                      = "ou=people,${var.user_base_dn}"
  create_rdn_attribute_type      = "uid"
}

resource "pingdirectory_group_rest_resource_type" "groupsRestResourceType" {
  id                             = "groups"
  display_name                   = "Groups"
  enabled                        = true
  search_base_dn                 = "ou=groups,${var.user_base_dn}"
  primary_display_attribute_type = "cn"
  resource_endpoint              = "groups"
  search_filter_pattern          = "(cn=*%%*)"
  structural_ldap_objectclass    = "groupOfUniqueNames"
  parent_dn                      = "ou=groups,${var.user_base_dn}"
}

#
# Specify the attributes that will be made available through the Delegated Admin API
#
resource "pingdirectory_generic_delegated_admin_attribute" "cnAttribute" {
  rest_resource_type_name = pingdirectory_user_rest_resource_type.usersRestResourceType.id
  attribute_type          = "cn"
  display_name            = "Full Name"
  display_order_index     = 0
}

resource "pingdirectory_generic_delegated_admin_attribute" "givenNameAttribute" {
  rest_resource_type_name = pingdirectory_user_rest_resource_type.usersRestResourceType.id
  attribute_type          = "givenName"
  display_name            = "First Name"
  display_order_index     = 1
}

resource "pingdirectory_generic_delegated_admin_attribute" "snAttribute" {
  rest_resource_type_name = pingdirectory_user_rest_resource_type.usersRestResourceType.id
  attribute_type          = "sn"
  display_name            = "Last Name"
  display_order_index     = 2
}

resource "pingdirectory_generic_delegated_admin_attribute" "mailAttribute" {
  rest_resource_type_name = pingdirectory_user_rest_resource_type.usersRestResourceType.id
  attribute_type          = "mail"
  display_name            = "Email"
  display_order_index     = 3
}

resource "pingdirectory_generic_delegated_admin_attribute" "uidAttribute" {
  rest_resource_type_name = pingdirectory_user_rest_resource_type.usersRestResourceType.id
  attribute_type          = "uid"
  display_name            = "User ID"
  display_order_index     = 4
}

resource "pingdirectory_generic_delegated_admin_attribute" "accountDisabledAttribute" {
  rest_resource_type_name = pingdirectory_user_rest_resource_type.usersRestResourceType.id
  attribute_type          = "ds-pwp-account-disabled"
  display_name            = "Account Disabled"
}


resource "pingdirectory_generic_delegated_admin_attribute" "cnGroupAttribute" {
  rest_resource_type_name = pingdirectory_group_rest_resource_type.groupsRestResourceType.id
  attribute_type          = "cn"
  display_name            = "Group"
}

resource "pingdirectory_generic_delegated_admin_attribute" "descriptionGroupAttribute" {
  rest_resource_type_name = pingdirectory_group_rest_resource_type.groupsRestResourceType.id
  attribute_type          = "description"
  display_name            = "Description"
}

#
# Create Delegated Admin Rights
#
resource "pingdirectory_delegated_admin_rights" "deladminRights" {
  id            = "deladmin"
  enabled       = true
  admin_user_dn = "uid=administrator,ou=people,${var.user_base_dn}"
}

#
# Create Delegated Admin Resource User and Group Rights
#
# This will add/update aci's found on the User and Group resource trees, defined in rest resource
#
resource "pingdirectory_delegated_admin_resource_rights" "usersRights" {
  delegated_admin_rights_name = pingdirectory_delegated_admin_rights.deladminRights.id
  rest_resource_type          = pingdirectory_user_rest_resource_type.usersRestResourceType.id
  admin_scope                 = "all-resources-in-base"
  admin_permission            = ["create", "read", "update", "delete", "manage-group-membership"]
  enabled                     = true
}

resource "pingdirectory_delegated_admin_resource_rights" "groupsRights" {
  delegated_admin_rights_name = pingdirectory_delegated_admin_rights.deladminRights.id
  rest_resource_type          = pingdirectory_group_rest_resource_type.groupsRestResourceType.id
  admin_scope                 = "all-resources-in-base"
  admin_permission            = ["create", "read", "update", "delete", "manage-group-membership"]
  enabled                     = true
}

#
# Create an access token validator for PingFederate tokens.
#
# WARNING: Use of the Blind Trust Trust Manager Provider is not recommended for production.  Instead, obtain PingFederate's
#          server certificate and add it to the JKS trust store using the 'manage-certificates trust-server-certificate'
#          command.  Then, update the PingFederateInstance external server to use the JKS Trust Manager Provider.
#          Consult the PingDirectory and PingData Security Guide for more information about configuring Trust Manager Providers.
#
resource "pingdirectory_default_blind_trust_manager_provider" "blindTrustManagerProvider" {
  id      = "Blind Trust"
  enabled = true
}

resource "pingdirectory_http_external_server" "pfExternalServer" {
  id                           = "pingfederate"
  base_url                     = "https://${var.pingfederate_hostname}:${var.pingfederate_https_port}"
  hostname_verification_method = "allow-all"
  trust_manager_provider       = pingdirectory_default_blind_trust_manager_provider.blindTrustManagerProvider.id
}

resource "pingdirectory_exact_match_identity_mapper" "entryUUIDMatchMapper" {
  id              = "entryUUIDMatch"
  enabled         = true
  match_attribute = ["entryUUID"]
  match_base_dn   = ["${var.user_base_dn}"]
}

resource "pingdirectory_ping_federate_access_token_validator" "pfAccessTokenValidator" {
  id                   = "pingfederate-validator"
  enabled              = true
  identity_mapper      = pingdirectory_exact_match_identity_mapper.entryUUIDMatchMapper.id
  subject_claim_name   = "Username"
  authorization_server = pingdirectory_http_external_server.pfExternalServer.id
  client_id            = "pingdirectory"
  client_secret        = "2FederateM0re"
}

#
# Complete the configuration of the Delegated Admin API.
#
resource "pingdirectory_custom_virtual_attribute" "delegatedAdminPrivilegeVirtualAttribute" {
  id      = "Delegated Admin Privilege"
  enabled = true
}

#
# A CORS policy is not needed when the app is running in the Ping Directory Server or Ping Proxy Server.
# To prevent a potential security vulnerability in the CORS policy, cors-allowed-origins should instead be set to the
# public name of the host, proxy, or load balancer that is going to be presenting the delegated admin web application.
#
resource "pingdirectory_http_servlet_cross_origin_policy" "daCrossOriginPolicy" {
  id                   = "Delegated Admin Cross-Origin Policy"
  cors_allowed_methods = ["GET", "OPTIONS", "POST", "DELETE", "PATCH"]
  cors_allowed_origins = ["*"]
}

resource "pingdirectory_delegated_admin_http_servlet_extension" "daServletExtension" {
  id                  = "Delegated Admin"
  access_token_scope  = "urn:pingidentity:directory-delegated-admin"
  response_header     = ["Cache-Control: no-cache, no-store, must-revalidate", "Expires: 0", "Pragma: no-cache"]
  cross_origin_policy = pingdirectory_http_servlet_cross_origin_policy.daCrossOriginPolicy.id
  # The above attribute must be removed (by setting to "") to allow destroying the Delegated Admin Cross-Origin Policy object
  # cross_origin_policy = ""
}


#
# Create an email account status notification handler for user creation.
# This handler cannot be enabled until an SMTP server is available in the global configuration.
#
resource "pingdirectory_simple_request_criteria" "daUserCreationRequestCriteria" {
  id                               = "Delegated Admin User Creation Request Criteria"
  operation_type                   = ["add"]
  included_target_entry_dn         = ["ou=people,${var.user_base_dn}"]
  any_included_target_entry_filter = ["(objectClass=inetOrgPerson)"]
  included_application_name        = ["PingDirectory Delegated Admin"]
}

resource "pingdirectory_multi_part_email_account_status_notification_handler" "daEmailAccountStatusNotificationHandler" {
  id                                             = "Delegated Admin Email Account Status Notification Handler"
  enabled                                        = false
  account_creation_notification_request_criteria = pingdirectory_simple_request_criteria.daUserCreationRequestCriteria.id
  account_created_message_template               = "config/account-status-notification-email-templates/delegated-admin-account-created.template"
}
