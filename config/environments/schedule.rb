set :output, "#{path}/log/cron_log.log"

every 1.minutes do
  rake "lazy_sync_cache:sync_applications"
  rake "lazy_sync_cache:sync_chats"
end
