# frozen_string_literal: true

module Api::V1
  class ChatsController < ::Api::BaseController
    before_action :set_application

    before_action :set_chat, only: %i[show edit update destroy]

    # GET /chats
    def index
      chats = @application.chats
      render json: ChatBlueprint.render(chats)
    end

    # GET /chats/1
    def show
      render json: ChatBlueprint.render(@chat)
    end

    # GET /chats/1/edit
    def edit; end

    # POST /chats
    def create
      begin
        creation_response = ChatCreationService.create_chat(@application)
      rescue ActiveRecord::StaleObjectError => e
        Rails.logger.info '============= Chat Version Error =============== '
        Rails.logger.info e
        Rails.logger.info '============= Chat Version Error =============== '
        creation_response = ChatCreationService.create_chat(@application)
      end

      render json: creation_response
    end

    # PATCH/PUT /chats/1
    def update
      raise Errors::CustomError.new(:bad_request, 400, @chat.errors.messages) unless @chat.update(chat_params)

      render json: ChatBlueprint.render(@chat)
    end

    # DELETE /chats/1
    def destroy
      raise Errors::CustomError.new(:bad_request, 400, @chat.errors.messages) unless @chat.destroy

      render json: ChatBlueprint.render(@chat)
    end

    private

    # Use callbacks to share common setup or constraints between actions.
    def set_application
      @application = Application.find_by!(number: params[:application_token])
    end

    # Use callbacks to share common setup or constraints between actions.
    def set_chat
      @chat = @application.chats.find_by!(number: params[:id])
    end
  end
end
