# frozen_string_literal: true

module Api::V1
  class ChatBlueprint < Blueprinter::Base
    fields :number, :messages_count

    association :application, blueprint: ApplicationBlueprint
  end
end
