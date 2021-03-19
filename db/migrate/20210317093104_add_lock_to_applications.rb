class AddLockToApplications < ActiveRecord::Migration[5.2]
  def change
    add_column :applications, :lock_version, :integer
  end
end
