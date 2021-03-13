# frozen_string_literal: true

module HasReference
  extend ActiveSupport::Concern

  MAX_REF_NUMBER = 10_000_000_000

  def generate_unique_reference(attribute = :reference_number)
    ref_num = SecureRandom.random_number(MAX_REF_NUMBER)
    ref_num = SecureRandom.random_number(MAX_REF_NUMBER) while self.class.where(attribute => ref_num).exists?

    ref_num
  end

  module ClassMethods
    # Inspired by ActiveRecord's has_secure_token method
    # rubocop:disable Namin/PredicateName
    def has_reference(attribute = :reference_number)
      validates attribute, uniqueness: true

      before_create do
        send("#{attribute}=", generate_unique_reference(attribute)) unless send("#{attribute}?")
      end
    end
    # rubocop:enable Namin/PredicateName
  end
end
