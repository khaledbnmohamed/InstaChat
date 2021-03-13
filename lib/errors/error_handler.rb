# frozen_string_literal: true

# Error module to Handle errors globally
module Errors
  module ErrorHandler
    def self.included(clazz)
      clazz.class_eval do
        rescue_from StandardError do |e|
          unless handled_errors(e)
            log_error(e, true)
            # Standard Error
            respond_backtrace(:standard_error, 500, e.to_s, e.backtrace)
          end
        end
      end
    end

    private

    def handled_errors(err)
      case err.class.name
      when 'ActiveRecord::RecordNotFound'
        respond(:record_not_found, 404, 're-check the params, it might be a typo!', nil)
      when 'Errors::CustomError'
        respond(err.error, err.status, err.message, err.extra)
      when 'ArgumentError'
        log_error(err, true)
        respond(:bad_request, 400, err.message, nil)
        false
      end
    end

    def log_error(error, raven = false, _extra_context_for_raven = {})
      Rails.logger.error error.message
      Rails.logger.error error.backtrace.join("\n")
      return unless raven
    end

    def respond(error, status, message, extra)
      json = {
        status: status,
        error: error,
        message: message,
        extra: extra
      }
      render json: json, status: status
    end

    def respond_backtrace(error, status, message, backtrace)
      json = {
        status: status,
        error: error,
        message: message,
        back_trace: backtrace
      }
      render json: json, status: status
    end
  end
end
