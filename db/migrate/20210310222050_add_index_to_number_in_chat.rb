class AddIndexToNumberInChat < ActiveRecord::Migration[6.0]
  def change
    add_index :chats, :number
  end
end
