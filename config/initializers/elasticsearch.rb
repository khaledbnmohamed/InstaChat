# frozen_string_literal: true

Elasticsearch::Model.client = Elasticsearch::Client.new log: true, host: ENV['ES_HOST'] || 'localhost:9200',
                                                        retry_on_failure: true
