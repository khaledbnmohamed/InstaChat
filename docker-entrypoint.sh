#!/bin/bash
set -e

if [ -f /home/instachat/instachat/tmp/pids/server.pid ]; then
  rm /home/instachat/instachat/tmp/pids/server.pid
fi

gem install bundler
bundle install

# overcommit --install
# overcommit --sign
bundle exec rake db:create
bundle exec rake db:migrate
bundle exec rake es:build_index
exec "$@"
