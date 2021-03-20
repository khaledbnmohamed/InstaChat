# frozen_string_literal: true

module Api::V1
  class MessagesController < ::Api::BaseController
    before_action :set_application, :set_chat, except: :reindex

    before_action :set_message, only: %i[show edit update destroy]

    # GET /messages
    def index
      messages = if params[:keyword].present?
                   @chat.messages.search(params[:keyword])
                 else
                   @chat.messages
                 end

      render_option = params[:keyword].present? ? elastic_search_result(messages) : messages
      render json: { chat: ChatBlueprint.render_as_json(@chat),
                     messages: MessageBlueprint.render_as_json(render_option) }
    end

    # GET /messages/1
    def show; end

    # GET /messages/1/edit
    def edit; end

    # POST /messages
    def create
      begin
        creation_response = MessageCreationService.create_message(@chat, message_params[:text])
      rescue ActiveRecord::StaleObjectError => e
        Rails.logger.info '============= Message Version Error =============== '
        Rails.logger.info e
        Rails.logger.info '============= Message Version Error =============== '
        creation_response = MessageCreationService.create_message(@chat, message_params[:text])
      end

      render json: creation_response
    end

    # PATCH/PUT /messages/1
    def update
      raise Errors::CustomError.new(:bad_request, 400, @message.errors.messages) unless @message.update(message_params)

      render json: MessageBlueprint.render(@message)
    end

    # DELETE /messages/1
    def destroy
      @message.destroy
      render json: MessageBlueprint.render(@message)
    end

    def reindex
      last_inserted_message = Message.last
      last_inserted_message.__elasticsearch__.index_document
      render json: MessageBlueprint.render(last_inserted_message)
    end

    private

    def set_application
      @application = Application.find_by!(number: params[:application_token])
    end

    # Use callbacks to share common setup or constraints between actions.
    def set_chat
      @chat = @application.chats.find_by!(number: params[:chat_id])
    end

    # Use callbacks to share common setup or constraints between actions.
    def set_message
      @message = @chat.messages.find_by!(number: params[:id])
    end

    # Only allow a trusted parameter "white list" through.
    def message_params
      params.require(:message).permit(:text)
    end

    def elastic_search_result(messages)
      hits = []
      if messages.present?
        messages.each do |m|
          hits << m['_source']
        end
      end
      hits
    end
  end
end
