# frozen_string_literal: true

class AddReferenceToChat < ActiveRecord::Migration[5.2]
  def change
    add_reference :messages, :chat, index: true, foreign_key: true
  end
end
