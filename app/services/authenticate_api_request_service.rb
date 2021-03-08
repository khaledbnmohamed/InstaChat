# frozen_string_literal: true

class AuthenticateApiRequestService
  def initialize(token)
    @token = token
  end

  def call
    return unless decoded_auth_token
    if expired_token?
      raise Errors::CustomError.new(:unauthorized, 401, 'Authentication Error')
    end

    employee || company
  end

  private

  attr_reader :token

  def employee
    if decoded_auth_token[:employee_id] && decoded_auth_token[:authenticated]
      @employee ||= Employee.find(decoded_auth_token[:employee_id])
    end
  end

  def application
    if decoded_auth_token[:company_id]
      @application ||= Application.find(decoded_auth_token[:company_id])
    end
  end

  def decoded_auth_token
    @decoded_auth_token ||= JsonWebToken.decode(token)
  end

  def expired_token?
    Time.zone.now > Time.zone.at(decoded_auth_token[:exp]) || AUTH_TOKENS_BLACKLIST.get(token)
  end
end
