module ChatCreationService
  def self.create_chat(application)
    chat_number = application.increment_chats_counter

    ChatCreationJob.perform_later(application_id: application.id, chat_number: chat_number)

    { number: chat_number, message: "Chat with number: #{chat_number} will be created shorlty" }
  end
end
