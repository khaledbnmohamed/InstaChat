module MessageCreationService
  def self.create_message(chat, text)
    message_number = chat.increment_messages_counter

    MessageCreationWorker.perform_async(chat.id, text, message_number)

    { number: message_number, message: "Message with number: #{message_number} will be created shorlty" }
  end
end
