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
require 'rails_helper'

RSpec.describe Chat, type: :model do
  it { is_expected.to belong_to(:application) }
end
