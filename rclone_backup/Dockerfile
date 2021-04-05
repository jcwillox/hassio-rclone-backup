ARG BUILD_FROM
FROM $BUILD_FROM

ARG BUILD_ARCH

ENV LANG C.UTF-8
ENV RCLONE_INSTALLED_VERSION 1.55.0

# Copy data for add-on
COPY ./app /
COPY ./run.sh /
COPY ./stdout.sh /usr/sbin/stdout

RUN chmod +x /run.sh \
    && chmod +x /usr/sbin/stdout

# Install requirements for add-on
RUN if [ "${BUILD_ARCH}" = "armhf" ]; then \
        curl -LJo /tmp/rclone.zip https://github.com/rclone/rclone/releases/download/v"${RCLONE_INSTALLED_VERSION}"/rclone-v"${RCLONE_INSTALLED_VERSION}"-linux-arm.zip; \
    elif [ "${BUILD_ARCH}" = "armv7" ]; then \
        curl -LJo /tmp/rclone.zip https://github.com/rclone/rclone/releases/download/v"${RCLONE_INSTALLED_VERSION}"/rclone-v"${RCLONE_INSTALLED_VERSION}"-linux-arm-v7.zip; \
    elif [ "${BUILD_ARCH}" = "aarch64" ]; then \
        curl -LJo /tmp/rclone.zip https://github.com/rclone/rclone/releases/download/v"${RCLONE_INSTALLED_VERSION}"/rclone-v"${RCLONE_INSTALLED_VERSION}"-linux-arm64.zip; \
    elif [ "${BUILD_ARCH}" = "i386" ]; then \
        curl -LJo /tmp/rclone.zip https://github.com/rclone/rclone/releases/download/v"${RCLONE_INSTALLED_VERSION}"/rclone-v"${RCLONE_INSTALLED_VERSION}"-linux-386.zip; \
    elif [ "${BUILD_ARCH}" = "amd64" ]; then \
        curl -LJo /tmp/rclone.zip https://github.com/rclone/rclone/releases/download/v"${RCLONE_INSTALLED_VERSION}"/rclone-v"${RCLONE_INSTALLED_VERSION}"-linux-amd64.zip; \
    else \
        exit 1; \
    fi \
    && unzip /tmp/rclone.zip -d /tmp \
    && cp /tmp/rclone-*/rclone /usr/bin/rclone  \
    && chmod 755 /usr/bin/rclone \
    && rm -fr /tmp/rclone* \
    && mkdir -p /root/.config/rclone \
    && touch /root/.config/rclone/rclone.conf \
    && apk add --no-cache tini \
    && pip install --no-cache-dir python-slugify

CMD [ "/run.sh" ]
