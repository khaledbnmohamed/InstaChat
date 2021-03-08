# frozen_string_literal: true

# Here we can include our constant variables that can be shared across multiple models and controllers

OTP_CODE_EXPIRATION_TIME = 180
INVALID_DAILY_TRIES_QUOTA = ENV['INVALID_DAILY_TRIES_QUOTA'].try(:to_i) || 5
EXPIRATION_TOKEN_TIME = 1.day.to_s
