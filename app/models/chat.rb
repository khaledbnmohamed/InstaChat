# frozen_string_literal: true

# == Schema Information
#
# Table name: chats
#
#  id             :bigint           not null, primary key
#  messages_count :integer          default(0)
#  number         :string(255)
#  created_at     :datetime         not null
#  updated_at     :datetime         not null
#  application_id :bigint           not null
#
# Indexes
#
#  index_chats_on_application_id  (application_id)
#
# Foreign Keys
#
#  fk_rails_...  (application_id => applications.id)
#
class Chat < ApplicationRecord
  has_reference :number

  belongs_to :application, inverse_of: :chats

  has_many :messages, dependent: :restrict_with_exception, inverse_of: :chat

  before_save :increment_chats_counter

  def increment_chats_counter
    application.with_lock do
      application.increment!(:chats_count)
    end
  end
end
