# frozen_string_literal: true

module Api::V1
  class MessageBlueprint < Blueprinter::Base
    fields :id

    association :chat, blueprint: ChatBlueprint
  end
end
