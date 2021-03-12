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
    def show; end

    # GET /chats/1/edit
    def edit; end

    # POST /chats
    def create
      ChatCreationJob.perform_later(application_id: @application.id)

      render json: 'Chat enqueued'
    end

    # PATCH/PUT /chats/1
    def update
      unless @chat.update(chat_params)
        raise Errors::CustomError.new(:bad_request, 400, @chat.errors.messages)
      end

      render json: ChatBlueprint.render(@chat)
    end

    # DELETE /chats/1
    def destroy
      @chat.destroy
      redirect_to chats_url, notice: 'Chat was successfully destroyed.'
    end

    private

    # Use callbacks to share common setup or constraints between actions.
    def set_application
      @application = Application.find_by(number: params[:application_token])
    end

    # Use callbacks to share common setup or constraints between actions.
    def set_chat
      @chat = Chat.find_by(number: params[:number])
    end
  end
end
