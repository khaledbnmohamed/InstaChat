class ChatCreationWorker
  include Sidekiq::Worker

  def perform(application_id, chat_number)
    # Do something
  end
end
