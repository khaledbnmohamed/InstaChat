# frozen_string_literal: true

# == Schema Information
#
# Table name: messages
#
#  id         :bigint           not null, primary key
#  text       :text(65535)      not null
#  created_at :datetime         not null
#  updated_at :datetime         not null
#  chat_id    :bigint
#
# Indexes
#
#  index_messages_on_chat_id  (chat_id)
#
# Foreign Keys
#
#  fk_rails_...  (chat_id => chats.id)
#
FactoryBot.define do
  factory :message do
    text { FFaker::Lorem.phrase }
    association :chat
  end
end
