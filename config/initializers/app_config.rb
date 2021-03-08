# frozen_string_literal: true

class AppConfig
  def self.load
    config = Rails.application.config_for 'app_config'

    config.each_pair do |key, value|
      cattr_accessor key

      if value.is_a?(Hash)
        send("#{key}=", hash_to_struct(value))
      else
        send("#{key}=", value)
      end
    end
  end

  def self.hash_to_struct(hash)
    keys = hash.keys
    values = hash.values_at(*keys)
    result = Struct.new(*keys).new(*values)

    hash.each_pair do |key, value|
      result.send("#{key}=", hash_to_struct(value)) if value.is_a?(Hash)
    end

    result
  end
end

AppConfig.load
