# frozen_string_literal: true

# == Schema Information
#
# Table name: chats
#
#  id             :bigint           not null, primary key
#  lock_version   :integer
#  messages_count :integer          default(0)
#  number         :string(255)
#  created_at     :datetime         not null
#  updated_at     :datetime         not null
#  application_id :bigint           not null
#
# Indexes
#
#  index_chats_on_application_id  (application_id)
#  index_chats_on_number          (number)
#
# Foreign Keys
#
#  fk_rails_...  (application_id => applications.id)
#
class Chat < ApplicationRecord
  # relations
  belongs_to :application, inverse_of: :chats

  has_many :messages, dependent: :restrict_with_exception, inverse_of: :chat

  # callbacks
  before_destroy :decrement_chats_counter

  #-- Callbacks
  after_update :update_chats_cache, if: :messages_count_changed?

  #-- Instance Methods
  def update_chats_cache
    CHATS_CACHE.set(number, messages_count)
  end

  def messages_cache_count(number)
    cached_messages_count = CHATS_CACHE.get(number.to_s)
    if cached_messages_count.blank?
      cached_messages_count = self.messages_count
      CHATS_CACHE.set(number.to_s, cached_messages_count.to_s)
    end
    cached_messages_count if cached_messages_count.present?
  end

  def increment_messages_counter
    cached_messages_count = messages_cache_count(self.number)
    incremented_chats_count = cached_messages_count.to_i + 1
    CHATS_CACHE.set(number.to_s, incremented_chats_count)
    incremented_chats_count
  end

  # Not using locks as it's less likely to conflict
  def decrement_chats_counter
    application.decrement!(:chats_count)
  end
end
