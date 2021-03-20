# frozen_string_literal: true

class AddIndexToNumberInChat < ActiveRecord::Migration[5.2]
  def change
    add_index :chats, :number
  end
end
