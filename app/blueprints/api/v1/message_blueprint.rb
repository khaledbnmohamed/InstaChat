# frozen_string_literal: true

module Api::V1
  class MessageBlueprint < Blueprinter::Base
    fields :text, :number

    association :chat, blueprint: ChatBlueprint
  end
end
