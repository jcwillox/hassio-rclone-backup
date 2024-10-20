ARG BUILD_FROM
FROM golang:1.23-alpine as build

ARG BUILD_ARCH
ARG BUILD_FROM

ENV CGO_ENABLED=0
WORKDIR /app

COPY scheduler .

RUN if [ "${BUILD_ARCH}" = "armhf" ]; then \
        export GOARCH=arm GOARM=6; \
    elif [ "${BUILD_ARCH}" = "armv7" ]; then \
        export GOARCH=arm GOARM=7; \
    elif [ "${BUILD_ARCH}" = "aarch64" ]; then \
        export GOARCH=arm64; \
    elif [ "${BUILD_ARCH}" = "i386" ]; then \
        export GOARCH=386; \
    elif [ "${BUILD_ARCH}" = "amd64" ]; then \
        export GOARCH=amd64; \
    else \
        exit 1; \
    fi \
    && go build -o scheduler.bin


FROM $BUILD_FROM

# Build arguments
ARG BUILD_ARCH
ARG BUILD_DATE
ARG BUILD_DESCRIPTION
ARG BUILD_NAME
ARG BUILD_REF
ARG BUILD_REPOSITORY
ARG BUILD_VERSION

ENV LANG C.UTF-8
# renovate: datasource=github-releases depName=rclone packageName=rclone/rclone
ENV RCLONE_INSTALLED_VERSION=1.68.1
# renovate: datasource=github-releases depName=rclone-webui packageName=rclone/rclone-webui-react
ENV RCLONE_WEBUI_INSTALLED_VERSION=2.0.5

# Install fuse
RUN apk add fuse \
    && sed -i 's/#user_allow_other/user_allow_other/' /etc/fuse.conf
    && ln -s /bin/fusermount /bin/fusermount3

# Copy root filesystem
COPY rootfs /
COPY --from=build /app/scheduler.bin /usr/bin/scheduler

# Download Web GUI
RUN curl -LJo /tmp/currentbuild.zip https://github.com/rclone/rclone-webui-react/releases/download/v"${RCLONE_WEBUI_INSTALLED_VERSION}"/currentbuild.zip \
    && unzip /tmp/currentbuild.zip -d /tmp \
    && mkdir -p /usr/share/rclone-webui-react \
    && mv /tmp/build/* /usr/share/rclone-webui-react \
    && rm -rf /tmp/*

# Fix permissions
RUN chmod a+x /etc/cont-init.d/*
RUN chmod a+x /etc/services.d/*/*

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
    && rm -rf /tmp/rclone* \
    && mkdir -p /root/.config/rclone

# Labels
LABEL \
    io.hass.name="${BUILD_NAME}" \
    io.hass.description="${BUILD_DESCRIPTION}" \
    io.hass.arch="${BUILD_ARCH}" \
    io.hass.type="addon" \
    io.hass.version=${BUILD_VERSION} \
    maintainer="Josh Willox" \
    org.opencontainers.image.title="${BUILD_NAME}" \
    org.opencontainers.image.description="${BUILD_DESCRIPTION}" \
    org.opencontainers.image.vendor="" \
    org.opencontainers.image.authors="Josh Willox" \
    org.opencontainers.image.licenses="MIT" \
    org.opencontainers.image.url="https://github.com/${BUILD_REPOSITORY}" \
    org.opencontainers.image.source="https://github.com/${BUILD_REPOSITORY}" \
    org.opencontainers.image.documentation="https://github.com/${BUILD_REPOSITORY}/blob/main/README.md" \
    org.opencontainers.image.created=${BUILD_DATE} \
    org.opencontainers.image.revision=${BUILD_REF} \
    org.opencontainers.image.version=${BUILD_VERSION}
