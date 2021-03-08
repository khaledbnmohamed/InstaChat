# frozen_string_literal: true

# == Schema Information
#
# Table name: applications
#
#  id            :bigint           not null, primary key
#  chats_counter :integer          default(0)
#  name          :string(255)      not null
#  number        :string(255)      not null
#  created_at    :datetime         not null
#  updated_at    :datetime         not null
#
class Application < ApplicationRecord
  has_reference :number

  has_many :chats, dependent: :restrict_with_exception, inverse_of: :application
end
