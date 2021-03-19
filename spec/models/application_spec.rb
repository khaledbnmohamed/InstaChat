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
require 'rails_helper'

RSpec.describe Application, type: :model do
  it { is_expected.to validate_presence_of(:name) }
end
