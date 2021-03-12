# frozen_string_literal: true

class AddIndexToNumberInApplication < ActiveRecord::Migration[6.0]
  def change
    add_index :applications, :number
  end
end
