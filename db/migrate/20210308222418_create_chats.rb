# frozen_string_literal: true

class CreateChats < ActiveRecord::Migration[6.0]
  def change
    create_table :chats do |t|
      t.references :application, null: false, foreign_key: true
      t.string :number
      t.integer :messages_count, default: 0

      t.timestamps
    end
  end
end
