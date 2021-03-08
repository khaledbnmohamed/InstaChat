# frozen_string_literal: true

class JsonWebToken
  EXPIRATION_TOKEN_TIME = 1.day.to_s

  class << self
    def encode(payload)
      exp = Time.zone.now + EXPIRATION_TOKEN_TIME
      payload = { data: payload, exp: exp.to_i }
      JWT.encode(payload, Rails.application.credentials.secret_key_base)
    end

    def decode(token)
      body = JWT.decode(token, Rails.application.credentials.secret_key_base)[0]
      HashWithIndifferentAccess.new body
    rescue StandardError
      nil
    end
  end
end
