APPLICATIONS_REDIS_CLIENT_OBJ ||= Redis.new(host: AppConfig.redis['host'],
                                            port: AppConfig.redis['port'],
                                            db: AppConfig.redis['applications_redis_cache']['db'])
CHATS_REDIS_CLIENT_OBJ ||= Redis.new(host: AppConfig.redis['host'],
                                     port: AppConfig.redis['port'],
                                     db: AppConfig.redis['chats_redis_cache']['db'])

APPLICATIONS_CACHE ||= Redis::Namespace.new(AppConfig.redis['namespace'].to_sym, redis: APPLICATIONS_REDIS_CLIENT_OBJ)
CHATS_CACHE ||= Redis::Namespace.new(AppConfig.redis['namespace'].to_sym, redis: CHATS_REDIS_CLIENT_OBJ)
