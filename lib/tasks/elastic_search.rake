namespace :elasticsearch do
  desc ' create and refresh elastic search indices'
  task build_index: :environment do
    Message.__elasticsearch__.create_index!
    Message.__elasticsearch__.refresh_index!
  end
end
