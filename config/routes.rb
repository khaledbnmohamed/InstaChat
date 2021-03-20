# frozen_string_literal: true

# require 'resque/server'
require 'sidekiq/web'

Rails.application.routes.draw do
  # For details on the DSL available within this file, see https://guides.rubyonrails.org/routing.html

  # mount(Resque::Server.new, at: '/resque') if Rails.env.development?
  mount Sidekiq::Web => '/sidekiq' # mount Sidekiq::Web in your Rails app

  swagger_documentation_constraint = lambda do |_|
    !Rails.env.production?
  end

  constraints swagger_documentation_constraint do
    mount Rswag::Ui::Engine => '/api-docs'
    mount Rswag::Api::Engine => '/api-docs'
  end
  namespace :api, constraints: { format: 'json' } do
    namespace :v1 do
      resources :applications, param: :token do
        resources :chats, params: :number do
          resources :messages, params: :number do
          end
        end
      end
      get 'messages/reindex', to: 'messages#reindex'
    end
  end
end
