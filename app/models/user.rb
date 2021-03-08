# frozen_string_literal: true

class User < ApplicationRecord
  validates :email, :mobile, presence: true
  validates :email, :mobile, uniqueness: true
  validates :email, format: { with: EMAIL_REGEX }
  validates :mobile, format: { with: MOBILE_REGEX }
  validates :password, confirmation: true
  validates :password, format: { with: PASSWORD_FORMAT }, on: :create
  validates :password, allow_nil: true, format: { with: PASSWORD_FORMAT }, on: :update

  # Instance Methods
  def auth_token(type)
    JsonWebToken.encode("#{id}+#{type}")
  end
end
