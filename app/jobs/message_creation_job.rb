# frozen_string_literal: true

class MessageCreationJob < ApplicationJob
  queue_as :default

  def perform(chat_id:, text:, number:)
    message = Message.new(chat_id: chat_id, text: text, number: number)

    message_saved = message.save
    Resque.logger.info "========= Message below #{message_saved ? ' ' : 'not '} saved ========="
    Resque.logger.info message
    Resque.logger.info '================================'

    # keeping data integrity if record not saved, as a rollback methdology
    # better approaches might come later
    message.decrement_messages_counter unless message_saved
  end
end
