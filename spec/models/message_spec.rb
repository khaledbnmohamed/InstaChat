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
require 'rails_helper'

RSpec.describe Message, type: :model do
  it { is_expected.to belong_to(:chat) }

  it { is_expected.to validate_presence_of(:text) }
end
