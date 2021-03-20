FROM ruby:2.7.1

ARG APP_USER_ID

ENV PROJECT_DIRECTORY    instachat
ENV APP_USER             instachat
ENV USER_HOME_DIRECTORY  /home/instachat
ENV APP_PATH             ${USER_HOME_DIRECTORY}/${PROJECT_DIRECTORY}
ENV BUNDLED_GEMS_PATH    /bundle
ENV BUNDLE_PATH          ${BUNDLED_GEMS_PATH}
ENV BUNDLE_APP_CONFIG    ${BUNDLED_GEMS_PATH}
ENV BUNDLE_BIN           ${BUNDLED_GEMS_PATH}/bin
ENV GEM_HOME             ${BUNDLED_GEMS_PATH}
ENV PATH ${BUNDLE_BIN}:${PATH}
ENV RUBYGEMS_VERSION     ${RUBYGEMS_VERSION}
ENV NODE_MODULES_DIR          ${APP_PATH}/node_modules
ENV CACHE_DIR                 ${APP_PATH}/tmp/cache
ENV PACKS_DIR                 ${APP_PATH}/public/packs

WORKDIR /tmp



RUN useradd -md ${USER_HOME_DIRECTORY} -u ${APP_USER_ID} -s /bin/bash ${APP_USER} && \
    mkdir -p ${APP_PATH} ${BUNDLED_GEMS_PATH} ${NODE_MODULES_DIR} ${CACHE_DIR} ${PACKS_DIR} && \
    chown -R ${APP_USER}:${APP_USER} ${APP_PATH} ${BUNDLED_GEMS_PATH} ${NODE_MODULES_DIR} ${CACHE_DIR} ${PACKS_DIR}


RUN apt-get update && apt-get -y install cron

# Copy entrypoint to docker home directory
WORKDIR ${USER_HOME_DIRECTORY}


#Enable Git completion (pressing tab key)
RUN echo "source /usr/share/bash-completion/completions/git" >> ${USER_HOME_DIRECTORY}/.bashrc


COPY docker-entrypoint.sh docker-entrypoint.sh
RUN chmod +x docker-entrypoint.sh

USER ${APP_USER}

WORKDIR ${APP_PATH}


# Left for debugging
# CMD while true ; do sleep 1 ; done
