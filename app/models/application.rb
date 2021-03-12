# frozen_string_literal: true

# == Schema Information
#
# Table name: applications
#
#  id          :bigint           not null, primary key
#  chats_count :integer          default(0)
#  name        :string(255)      not null
#  number      :string(255)      not null
#  created_at  :datetime         not null
#  updated_at  :datetime         not null
#
# Indexes
#
#  index_applications_on_number  (number)
#
class Application < ApplicationRecord
  has_reference :number

  has_many :chats, dependent: :restrict_with_exception, inverse_of: :application

  # validations
  validates :name, presence: true
end
