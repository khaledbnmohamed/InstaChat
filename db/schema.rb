# frozen_string_literal: true

# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# This file is the source Rails uses to define your schema when running `rails
# db:schema:load`. When creating a new database, `rails db:schema:load` tends to
# be faster and is potentially less error prone than running all of your
# migrations from scratch. Old migrations may fail to apply correctly if those
# migrations use external dependencies or application code.
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema.define(version: 20_210_308_224_009) do
  create_table 'applications', options: 'ENGINE=InnoDB DEFAULT CHARSET=utf8mb4', force: :cascade do |t|
    t.string 'name', null: false
    t.string 'number', null: false
    t.integer 'chats_counter', default: 0
    t.datetime 'created_at', precision: 6, null: false
    t.datetime 'updated_at', precision: 6, null: false
  end

  create_table 'chats', options: 'ENGINE=InnoDB DEFAULT CHARSET=utf8mb4', force: :cascade do |t|
    t.bigint 'application_id', null: false
    t.string 'number'
    t.datetime 'created_at', precision: 6, null: false
    t.datetime 'updated_at', precision: 6, null: false
    t.index ['application_id'], name: 'index_chats_on_application_id'
  end

  create_table 'messages', options: 'ENGINE=InnoDB DEFAULT CHARSET=utf8mb4', force: :cascade do |t|
    t.text 'text', null: false
    t.datetime 'created_at', precision: 6, null: false
    t.datetime 'updated_at', precision: 6, null: false
    t.bigint 'chat_id'
    t.index ['chat_id'], name: 'index_messages_on_chat_id'
  end

  add_foreign_key 'chats', 'applications'
  add_foreign_key 'messages', 'chats'
end
