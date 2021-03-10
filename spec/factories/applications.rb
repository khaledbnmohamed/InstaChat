# frozen_string_literal: true

# == Schema Information
#
# Table name: applications
#
#  id          :bigint           not null, primary key
#  chats_count :integer          default(0)
#  name        :string(255)      not null
#  number      :string(255)      not null
#  created_at  :datetime         not null
#  updated_at  :datetime         not null
#
FactoryBot.define do
  factory :application do
    sequence(:name) { FFaker::LoremAR.word }
  end
end
