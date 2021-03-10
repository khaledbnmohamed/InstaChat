class MessageCreationJob < ApplicationJob
  queue_as :messages

  def perform(chat_id:,text:)
    message = Message.new(chat_id: chat_id, text: text)

    message_saved = message.save
    Resque.logger.info "========= Message below #{message_saved ? " " : "not "} saved ========="
    Resque.logger.info message
    Resque.logger.info "================================"
  end
end
