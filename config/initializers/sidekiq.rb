Sidekiq.configure_server do |config|
  config.redis = { namespace: 'instachat', url: AppConfig.redis['url'] }
end

Sidekiq.configure_client do |config|
  config.redis = { namespace: 'instachat', url: AppConfig.redis['url'] }
end
