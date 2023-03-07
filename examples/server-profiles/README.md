# Server profile examples
The examples in this directory correspond to the [Ping Identity server profiles](https://github.com/pingidentity/pingidentity-server-profiles) used with the PingDirectory Docker image. The examples replace the typical `pd.profile/dsconfig` folder found in PingDirectory server profiles. The dsconfig is instead applied with the Terraform provider. Other aspects of the server profile (such as schema and ldif files) are not managed by the Terraform provider. For the baseline example in particular, non-default schema from the baseline server profile is necessary for some of the provided config objects to be created.

The `getting-started` directory corresponds to the PingDirectory [getting-started server profile](https://github.com/pingidentity/pingidentity-server-profiles/tree/master/getting-started/pingdirectory/pd.profile/dsconfig).

The `baseline` directory corresponds to the PingDirectory [baseline server profile](https://github.com/pingidentity/pingidentity-server-profiles/tree/master/baseline/pingdirectory/pd.profile/dsconfig).
