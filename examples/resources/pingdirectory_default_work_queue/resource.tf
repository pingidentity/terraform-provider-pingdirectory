resource "pingdirectory_default_work_queue" "myWorkQueue" {
  num_worker_threads      = 2
  max_work_queue_capacity = 800
}
