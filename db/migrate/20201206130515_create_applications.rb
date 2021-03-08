# frozen_string_literal: true

class CreateApplications < ActiveRecord::Migration[6.0]
  def change
    create_table :applications do |t|
      t.string :name, null: false
      t.string :number, null: false
      t.integer :chats_counter, default: 0

      t.timestamps
    end
  end
end
