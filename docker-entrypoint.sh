#!/bin/bash
set -e

if [ -f /home/instachat/instachat/tmp/pids/server.pid ]; then
  rm /home/instachat/instachat/tmp/pids/server.pid
fi

gem install bundler
bundle install

bundle exec rake db:create
bundle exec rake db:migrate
bundle exec rake elasticsearch:build_index

# overcommit --install
# overcommit --sign
exec "$@"
