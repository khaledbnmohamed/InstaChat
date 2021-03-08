# frozen_string_literal: true

module Api
  class BaseController < ApplicationController
    include ::Errors::ErrorHandler
    respond_to :json
    protect_from_forgery with: :null_session, if: -> { request.format.json? }
  end
end
