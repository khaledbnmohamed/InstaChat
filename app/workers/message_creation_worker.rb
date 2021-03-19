class MessageCreationWorker
  include Sidekiq::Worker

  def perform(chat_id, text, number)
    # Do something
  end
end
