# frozen_string_literal: true

# == Schema Information
#
# Table name: applications
#
#  id           :bigint           not null, primary key
#  chats_count  :integer          default(0)
#  lock_version :integer
#  name         :string(255)      not null
#  number       :string(255)      not null
#  created_at   :datetime         not null
#  updated_at   :datetime         not null
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

  #-- Callbacks
  after_update :update_applications_cache, if: :chats_count_changed?

  #-- Instance Methods
  def update_applications_cache
    APPLICATIONS_CACHE.set(number, chats_count)
  end

  def chats_cache_count(number)
    cached_chats_count = APPLICATIONS_CACHE.get(number.to_s)
    if cached_chats_count.blank?
      cached_chats_count = self.chats_count
      APPLICATIONS_CACHE.set(number.to_s, cached_chats_count)
    end
    cached_chats_count if cached_chats_count.present?
  end

  def increment_chats_counter
    cached_chats_count = chats_cache_count(self.number)
    incremented_chats_count = cached_chats_count.to_i + 1
    APPLICATIONS_CACHE.set(number.to_s, incremented_chats_count)
    incremented_chats_count
  end
end
