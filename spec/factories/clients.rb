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
FactoryBot.define do
  factory :application do
    sequence(:name) { |i| "عميل #{i}" }
    sequence(:email) { |i| "test-#{i}@instachat.com" }
    sequence(:mobile) { "0#{rand(505_000_000..505_999_999)}" }
    password { 'QWer12%!' }
    password_confirmation { 'QWer12%!' }
    latitude { rand(1..10) }
    location { |n| "#{FFaker::LoremAR.word}_#{n}" }
    longitude { rand(1..10) }
    neighborhood       { |n| "#{FFaker::LoremAR.word}_#{n}" }
    primary_address    { |n| "#{FFaker::LoremAR.word}_#{n}" }
    secondary_address  { |n| "#{FFaker::LoremAR.word}_#{n}" }
    street { rand(1..10) }
  end
end
