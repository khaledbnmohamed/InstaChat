module ChatCreationService
  def self.create_chat(application)
    chat_number = application.increment_chats_counter

    ChatCreationWorker.perform_async(application.id, chat_number)

    { number: chat_number, message: "Chat with number: #{chat_number} will be created shorlty" }
  end
end
