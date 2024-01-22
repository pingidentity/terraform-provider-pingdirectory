resource "pingdirectory_post_ldif_export_task_processor" "myPostLdifExportTaskProcessor" {
  id                  = "MyPostLdifExportTaskProcessor"
  type                = "upload-to-s3"
  aws_external_server = "myExternalServer"
  s3_bucket_name      = "myS3Bucket"
  enabled             = false
}
