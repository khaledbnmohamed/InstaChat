# Use this file to easily define all of your cron jobs.
#
# It's helpful, but not entirely necessary to understand cron before proceeding.
# http://en.wikipedia.org/wiki/Cron

# Example:
#
# set :output, "/path/to/my/cron_log.log"
#
every 1.minute do
  puts "#################################################"
  puts "#################################################"
  puts "#################################################"
  puts "#################################################"
  puts "Running the syncing tasks"
  puts "#################################################"
  puts "#################################################"
  puts "#################################################"
  puts "#################################################"

  rake "lazy_sync_cache:sync_applications"
  rake "lazy_sync_cache:sync_chats"

end
#
# every 4.days do
#   runner "AnotherModel.prune_old_records"
# end

# Learn more: http://github.com/javan/whenever
