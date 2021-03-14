# frozen_string_literal: true

# == Schema Information
#
# Table name: applications
#
#  id                :bigint           not null, primary key
#  area              :string
#  building_number   :string
#  email             :string           not null
#  latitude          :string
#  location          :string
#  longitude         :string
#  mobile            :string           not null
#  name              :string           not null
#  neighborhood      :string
#  password_digest   :string           not null
#  primary_address   :string           not null
#  secondary_address :string
#  street            :string
#  created_at        :datetime         not null
#  updated_at        :datetime         not null
#
require 'rails_helper'

RSpec.describe Chat, type: :model do
  it { is_expected.to belong_to(:application) }
end
