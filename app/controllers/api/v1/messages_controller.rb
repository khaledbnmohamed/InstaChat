# frozen_string_literal: true

module Api::V1
  class MessagesController < ::Api::BaseController
    before_action :set_chat

    before_action :set_message, only: %i[show edit update destroy]

    # GET /messages
    def index
      messages = if params[:keyword].present?
                   @chat.messages.search(params[:keyword])
                 else
                   @chat.messages
      end
      render json: params[:keyword].present? ? messages : MessageBlueprint.render(messages)
    end

    # GET /messages/1
    def show; end

    # GET /messages/1/edit
    def edit; end

    # POST /messages
    def create
      MessageCreationJob.perform_later(chat_id: @chat.id, text: message_params[:text] )

      render json: "enqueued"
    end

    # PATCH/PUT /messages/1
    def update
      unless @message.update(message_params)
        raise Errors::CustomError.new(:bad_request, 400, @message.errors.messages)
      end

      render json: MessageBlueprint.render(@message)
    end

    # DELETE /messages/1
    def destroy
      @message.destroy
      redirect_to messages_url, notice: 'Message was successfully destroyed.'
    end

    private

    # Use callbacks to share common setup or constraints between actions.
    def set_chat
      @chat = Chat.find_by(number: params[:chat_id])
    end

    # Use callbacks to share common setup or constraints between actions.
    def set_message
      @message = Message.find_by(number: params[:number])
    end

    # Only allow a trusted parameter "white list" through.
    def message_params
      params.require(:message).permit(:text)
    end
  end
end
