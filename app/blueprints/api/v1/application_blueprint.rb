# frozen_string_literal: true

module Api::V1
  class ApplicationBlueprint < Blueprinter::Base
    fields :name, :number, :chats_count
  end
end
