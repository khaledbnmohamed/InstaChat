require 'rails_helper'

RSpec.describe ChatCreationJob, type: :job do
  describe '#perform_later' do
    ActiveJob::Base.queue_adapter = :test
    let(:application) { FactoryBot.create(:application) }

    describe 'chat service enqueue chat creation job' do
      it 'enqueued on calling chat service' do
        expect { ChatCreationService.create_chat(application) }.to have_enqueued_job.exactly(:once)
      end
    end
  end
end
