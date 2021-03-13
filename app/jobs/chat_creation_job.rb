# frozen_string_literal: true

class ChatCreationJob < ApplicationJob
  queue_as :chats

  def perform(application_id:, chat_number:)
    chat = Chat.new(application_id: application_id, number: chat_number)

    chat_saved = chat.save
    Resque.logger.info "========= chat below #{chat_saved ? ' ' : 'not '} saved ========="
    Resque.logger.info chat
    Resque.logger.info '================================'

    # keeping data integrity if record not saved, as a rollback methdology
    # better approaches might come later
    chat.decrement_chats_counter unless chat_saved
  end
end
