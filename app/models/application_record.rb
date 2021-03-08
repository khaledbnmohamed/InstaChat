# frozen_string_literal: true

class ApplicationRecord < ActiveRecord::Base
  include HasReference

  self.abstract_class = true
end
