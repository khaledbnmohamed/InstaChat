# frozen_string_literal: true

Resque.redis = if Rails.env.production?
                 Redis.new(url: "redis://#{AppConfig.redis['master_name']}",
                           sentinels: AppConfig.redis['sentinels'],
                           role: :master,
                           db: AppConfig.redis['rescue']['db'])
               else
                 Redis.new(host: AppConfig.redis['host'],
                           port: AppConfig.redis['port'],
                           db: AppConfig.redis['rescue']['db'])
               end

Resque.logger.level = Logger::DEBUG
