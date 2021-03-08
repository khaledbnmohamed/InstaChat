# frozen_string_literal: true

source 'https://rubygems.org'
git_source(:github) { |repo| "https://github.com/#{repo}.git" }

ruby '2.7.1'

# Bundle edge Rails instead: gem 'rails', github: 'rails/rails'
gem 'rails', '~> 6.0.2', '>= 6.0.2.2'
# Use Puma as the app server
gem 'puma', '~> 4.3'
# Turbolinks makes navigating your web application faster. Read more: https://github.com/turbolinks/turbolinks
gem 'turbolinks', '~> 5'
# Simple, Fast, and Declarative Serialization Library for Ruby
gem 'blueprinter', '~> 0.23.4'
# Use Redis adapter to run Action Cable in production
gem 'redis', '~> 4.0'
# Use Active Model has_secure_password
gem 'bcrypt', '~> 3.1.7'
# Use wisper to broadcast events
gem 'wisper', '2.0.0'
# Use wisper-activejob for async broadcasting
gem 'wisper-activejob', '~> 1.0.0'
# Use Active Storage variant
# gem 'image_processing', '~> 1.2'
# Use mysql as the database for Active Record
gem 'mysql2', '>= 0.4.4'

# Reduces boot times through caching; required in config/boot.rb
gem 'bootsnap', '>= 1.4.2', require: false

# Localization Helpers
gem 'devise-i18n', '=1.9.1'
gem 'rails-i18n', '=6.0.0'

# Authentication
gem 'devise', '4.7.1'
gem 'jwt', '2.2.1'

# Authorization
gem 'pundit', '=2.1.0'

# Search and filtering for AR
gem 'ransack', '=2.3.2'

gem 'resque', '=2.0.0'
gem 'rswag', '~> 2.3.1'

gem 'rack-cors', '~> 1.1.0'

gem 'kaminari', '=1.2.1'

gem 'factory_bot_rails', '=5.2.0'
gem 'ffaker', '=2.14.0'

gem 'timecop', '~> 0.9.1'

gem 'rspec-collection_matchers', '=1.2.0'
gem 'rspec-rails', '=4.0.0'

group :development, :test do
  # Call 'byebug' anywhere in the code to stop execution and get a debugger console
  gem 'byebug', platforms: %i[mri mingw x64_mingw]
  gem 'dotenv-rails', '=2.7.5'

  gem 'pry', '=0.13.1'
  gem 'pry-doc', '=1.1.0'
  gem 'pry-rails', '=0.3.9'
  gem 'webmock', '=3.8.3'
end

group :development do
  # Access an interactive console on exception pages or by calling 'console' anywhere in the code.
  gem 'listen', '>= 3.0.5', '< 3.2'
  gem 'web-console', '>= 3.3.0'
  # Spring speeds up development by keeping your application running in the background. Read more: https://github.com/rails/spring
  gem 'overcommit', '~> 0.49.0'
  gem 'rubocop', '=0.82.0', require: false
  gem 'rubocop-performance', '=1.5.2', require: false
  gem 'rubocop-rails', '=2.5.2', require: false
  gem 'rubocop-rspec', '=1.38.1', require: false
  gem 'spring'
  gem 'spring-watcher-listen', '~> 2.0.0'

  gem 'annotate', '=3.1.1'
  gem 'brakeman', '=4.8.1', require: false
  gem 'bullet', '=6.1.0'
end

group :test do
  gem 'shoulda-matchers', '~> 4.3.0'
  gem 'simplecov', '=0.18.5', require: false
end

# Windows does not include zoneinfo files, so bundle the tzinfo-data gem
gem 'tzinfo-data', platforms: %i[mingw mswin x64_mingw jruby]
