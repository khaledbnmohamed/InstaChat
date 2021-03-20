# frozen_string_literal: true

class AddIndexToNumberInApplication < ActiveRecord::Migration[5.2]
  def change
    add_index :applications, :number
  end
end
