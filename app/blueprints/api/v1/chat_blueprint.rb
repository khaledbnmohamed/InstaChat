# frozen_string_literal: true

module Api::V1
  class ChatBlueprint < Blueprinter::Base
    fields :number

    association :application, blueprint: ApplicationBlueprint
  end
end
