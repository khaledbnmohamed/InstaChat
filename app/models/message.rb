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
class Message < ApplicationRecord
  include Elasticsearch::Model
  include Elasticsearch::Model::Callbacks

  belongs_to :chat, inverse_of: :messages

  before_save :increment_messages_counter

  index_name    'text_index'
  document_type 'text'

  settings do
    mappings dynamic: false do
      indexes :text, type: :text
    end
  end

  def increment_messages_counter
    chat.with_lock do
      chat.increment!(:messages_count)
      message.number = messages_count
    end
  end
end
