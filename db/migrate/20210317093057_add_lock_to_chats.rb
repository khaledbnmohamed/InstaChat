class AddLockToChats < ActiveRecord::Migration[5.2]
  def change
    add_column :chats, :lock_version, :integer
  end
end
