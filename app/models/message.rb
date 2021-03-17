# frozen_string_literal: true

# == Schema Information
#
# Table name: messages
#
#  id         :bigint           not null, primary key
#  number     :string(255)
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

  # relations
  belongs_to :chat, inverse_of: :messages

  # validations
  validates :text, presence: true

  # elastic search configurations
  index_name    'text_index'
  document_type 'text'

  settings do
    mappings dynamic: false do
      indexes :text, type: :text
    end
  end

  def self.search(query)
    __elasticsearch__.search(
      {
        query: {
          wildcard: {
            text: "*#{query}*"
          }
        }
      }
    )
  end

  # Not using locks as it's less likely to conflict
  def decrement_messages_counter
    chat.decrement!(:messages_count)
  end
end
