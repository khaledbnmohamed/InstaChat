# frozen_string_literal: true

class AddNumberToMessages < ActiveRecord::Migration[5.2]
  def change
    add_column :messages, :number, :string, index: true
  end
end
