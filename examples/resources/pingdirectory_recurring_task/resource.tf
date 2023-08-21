resource "pingdirectory_recurring_task" "myRecurringTask" {
  name                          = "MyRecurringTask"
  type                          = "generate-server-profile"
  profile_directory             = "/opt/out/instance"
  retain_previous_profile_count = 10
}
