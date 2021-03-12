#!/bin/bash
set -e

if [ -f /home/instachat/instachat/tmp/pids/server.pid ]; then
  rm /home/instachat/instachat/tmp/pids/server.pid
fi

gem install bundler
bundle install

# overcommit --install
# overcommit --sign

exec "$@"
