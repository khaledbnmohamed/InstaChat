# frozen_string_literal: true

module Errors
  class ServiceUnavailableError < StandardError
    attr_reader :status, :error, :message, :extra

    def initialize(error = nil, status = nil, message = nil, extra = nil)
      @error = error || :service_unavailable
      @status = status || 503
      @extra = extra
      @message = message || 'Service is currently unavailable'
    end

    def fetch_json
      {
        status: @status,
        error: @error,
        message: @message,
        extra: @extra
      }
    end
  end
end
