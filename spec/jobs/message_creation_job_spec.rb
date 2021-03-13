require 'rails_helper'

RSpec.describe MessageCreationJob, type: :job do
  describe '#perform_later' do
    ActiveJob::Base.queue_adapter = :test
    let(:chat) { FactoryBot.create(:chat) }

    describe 'message service enqueue message creation job' do
      it 'enqueued on calling message service' do
        expect { MessageCreationService.create_message(chat, 'test') }.to have_enqueued_job.exactly(:once)
      end
    end
  end
end
