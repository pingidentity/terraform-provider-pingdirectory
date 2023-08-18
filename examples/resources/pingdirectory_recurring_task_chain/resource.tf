resource "pingdirectory_recurring_task_chain" "myRecurringTaskChain" {
  name                          = "MyRecurringTaskChain"
  recurring_task                = ["Export All Non-Administrative Backends"]
  scheduled_date_selection_type = "every-day"
  scheduled_time_of_day         = ["10:00", "11:00"]
}
