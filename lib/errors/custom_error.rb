# frozen_string_literal: true

module Errors
  class CustomError < StandardError
    attr_reader :status, :error, :message, :extra

    def initialize(error = nil, status = nil, message = nil, extra = nil)
      @error = error || :unprocessable_entity
      @status = status || 422
      @extra = extra
      @message = message || 'Something went wrong'
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
