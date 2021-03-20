namespace :lazy_sync_cache do
  desc ' sync redis cached values with the database records'
  task sync_applications: :environment do
    Application.in_batches do |batch|
      batch.each do |app|
        cached_value = APPLICATIONS_CACHE.get(app.number.to_s)
        # if value not found, means that no interaction happened
        # on the object so nothing to be saved
        if cached_value.present?
          app.chats_count = cached_value
          app.save
          puts "Updated app number: #{app.number} succefully with chat_counts= #{cached_value} "
        end
      end
    end
  end

  task sync_chats: :environment do
    Chat.in_batches do |batch|
      batch.each do |chat|
        cached_value = CHATS_CACHE.get(chat.number.to_s)
        # if value not found, means that no interaction happened
        # on the object so nothing to be saved
        if cached_value.present?
          chat.messages_count = cached_value
          chat.save
          puts "Updated chat number: #{chat.number} succefully with messages_count= #{cached_value} "
        end
      end
    end
  end
end
