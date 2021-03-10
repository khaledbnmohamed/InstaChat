class ChatCreationJob < ApplicationJob
  queue_as :chats

  def perform(application_id:)
    chat = Chat.new(application_id: application_id)

    chat_saved = chat.save
    Resque.logger.info "========= chat below #{chat_saved ? " " : "not "} saved ========="
    Resque.logger.info chat
    Resque.logger.info "================================"
  end
end
