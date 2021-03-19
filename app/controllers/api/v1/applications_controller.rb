# frozen_string_literal: true

module Api::V1
  class ApplicationsController < ::Api::BaseController
    before_action :set_application, only: %i[show edit update destroy]

    # GET /applications
    def index
      @applications = Application.all
    end

    # GET /applications/1
    def show
      render json: ApplicationBlueprint.render(@application)
    end

    # GET /applications/1/edit
    def edit; end

    # POST /applications
    def create
      application = Application.new(application_params)

      raise Errors::CustomError.new(:bad_request, 400, application.errors.messages) unless application.save

      render json: ApplicationBlueprint.render(application)
    end

    # PATCH/PUT /applications/1
    def update
      unless @application.update(application_params)
        raise Errors::CustomError.new(:bad_request, 400, @application.errors.messages)
      end

      render json: ApplicationBlueprint.render(@application)
    end

    # DELETE /applications/1
    def destroy
      @application.destroy
      render json: ApplicationBlueprint.render(@application)
    end

    private

    # Use callbacks to share common setup or constraints between actions.
    def set_application
      @application = Application.find_by!(number: params[:token])
    end

    # Only allow a trusted parameter "white list" through.
    def application_params
      params.require(:application).permit(:name)
    end
  end
end
