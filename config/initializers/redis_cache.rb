# CONFIGURATIONS_CACHE_ORGINAL_OBJ ||= Redis.new( host: AppConfig.redis[:host],
#                                               port: AppConfig.redis[:port],
#                                               db: AppConfig.redis[:system_configuration_redis_cache][:db]
#                                             )

APPLICATIONS_REDIS_CLIENT ||= Redis.new( host: AppConfig.redis['host'],
                                    port: AppConfig.redis['port'],
                                    db: AppConfig.redis['applications_redis_cache']['db']
                                  )
CHATS_REDIS_CLIENT ||= Redis.new( host: AppConfig.redis['host'],
                                    port: AppConfig.redis['port'],
                                    db: AppConfig.redis['chats_redis_cache']['db']
                                  )

# CONFIGURATIONS_CACHE ||= Redis::Namespace.new(AppConfig.redis[:namespace].to_sym, redis: CONFIGURATIONS_CACHE_ORGINAL_OBJ)
