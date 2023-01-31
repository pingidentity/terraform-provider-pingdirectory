# PingDirectory Terraform Provider Plugin Notes

Terraform is a powerful tool for managing configuration, particularly infrastructure.  When it comes to applications, however, some of the Terraform assumptions and conventions do not apply directly. As a result, there are a few things specific to PingDirectory that have had to be addressed in a non-typical manner under Terraform.  Some of these situations result in Terraform behaving a bit different than other providers; others require the use of attributes in a different manner in the provider when compared to normal PingDirectory operations.

This document will provide guidance on things you should know when either using the provider or if you wish to contribute to the provider.

## Edit-only Resources

Certain configuration objects are edit-only. These objects do not fit in the typical Terraform model. Terraform expects to manage the entire lifecycle of an object, from creation to deletion. In the case of PingDirectory, some objects are inherent, often with default values.  However, leaving these objects out of the control of Terraform would severely restrict the ability to manage PingDirectory configuration using Terraform HCL. As an example, consider the global configuration, which is very useful and often modified in PingDirectory; this resource is edit-only.

Other providers have functionality to manage resources in this category. One particular example is in this AWS provider resource: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/default_security_group.  For items such as this and others in the registry, the strategy in Terraform is to “adopt” them on a typical create, and "forget" them on a delete action.  When Terraform forgets, the resource still exists and can be managed elsewhere. This pattern is what the PingDirectory provider has implemented for edit-only objects.

However, there is one major difference between how the PingDirectory provider handles edit-only resources versus the AWS provider example above. This difference is in how the resource is initialized.  In the ***create*** action of the AWS provider, all of the properties of the existing edit-only object are wiped and replaced completely with what is specified in the Terraform file. In the global configuration of PingDirectory, for example, this action was determined to be impractical.  Doing a wipe would require the person managing the global configuration to specify every single property of the global configuration (over 80 in total) in order to manage it with Terraform.  All properties would be managed, even if only one of them needed to be changed from the default. Instead, a strategy of marking every property as **_Optional_** and **_Computed_** in the Terraform schema was adopted.  In this way, defaults and existing values are allowed to flow through from PingDirectory, and from that point the user only needs to specify the values to be modified and managed in the Terraform file.

While the global configuration is used as an example here, all edit-only resources are handled in this manner.

### Attributes in edit-only resources

As part of the edit-only strategy, all attributes of these resources are defined as **_computed_** so that users do not need to specify every available attribute.

## id instead of name

When running **_dsconfig_** commands, the instance is referred to using the **name** variable.  The convention in a Terraform provider is to use **id** as the naming attribute in most cases.

## Boolean and Optional fields

In Terraform, you can mark attributes of the schema as **_Computed_**. This annotation means that the provider can set its own value for the attribute. In some cases with PingDirectory, the use of this convention becomes necessary. For example, for any optional boolean value in a configuration object, the attribute must be marked as Computed. This setting is because PingDirectory will return a default boolean value (*false* or *true* depending on the attribute) if no value is specified. The same is true for any optional Set values, as PingDirectory will return a default empty set if no value is specified.

Computed values lead to some complications when trying to “unset” a value for edit-only configuration objects. If you simply remove a value from your Terraform file for a resource, the former Computed value will just remain as-is; Terraform will see this as no changes having been made. This behavior means that unsetting a value for an edit-only resource requires a deliberate action.

- For strings, PingDirectory sees the empty string as equivalent to unsetting that string, so you can simply set that string property to empty. 
- For sets, an empty set can be used.
- For booleans, there is no need to unset because *true* or *false* can be used.
- For integer values, however, **there is no way to unset a value using Terraform after it has been set for PingDirectory edit-only objects** due to the way they are currently implemented.

## Empty strings

Empty strings are treated as the equivalent to null in the provider.

## Contributing

We appreciate your help! To contribute through logging issues or creating pull requests, please read the [contribution guidelines](CONTRIBUTING.md)
