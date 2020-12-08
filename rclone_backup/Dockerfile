ARG BUILD_FROM
FROM $BUILD_FROM

ENV LANG C.UTF-8

# Copy data for add-on
COPY ./app /
COPY ./run.sh /
COPY ./tostdout.sh /usr/sbin/tostdout

RUN chmod +x /run.sh \
  && chmod +x /usr/sbin/tostdout

# Install requirements for add-on
RUN apk add --no-cache tini rclone \
    && pip install --no-cache-dir python-slugify

CMD [ "/run.sh" ]