# frozen_string_literal: true

if Rails.env.production?
  Resque.redis = Redis.new(url: "redis://#{AppConfig.redis['master_name']}",
                           sentinels: AppConfig.redis['sentinels'],
                           role: :master,
                           db: AppConfig.redis['rescue']['db'])
else
  Resque.redis = Redis.new(host: AppConfig.redis['host'],
                           port: AppConfig.redis['port'],
                           db: AppConfig.redis['rescue']['db'])
end

Resque.logger.level = Logger::DEBUG
