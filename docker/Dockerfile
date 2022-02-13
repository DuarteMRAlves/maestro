ARG GO=1.17.6

FROM golang:${GO}-bullseye

ARG WORKSPACE=/opt/maestro
ARG USER=maestro
ARG GROUP=maestro-group

# Create non-privileged user to run
RUN addgroup --system ${GROUP} && \
    adduser --system --no-create-home --ingroup ${GROUP} ${USER} && \
    mkdir ${WORKSPACE} && \
    chown -R ${USER}:${GROUP} ${WORKSPACE}

WORKDIR ${WORKSPACE}

COPY ./target/maestro maestro

RUN chmod +x maestro

CMD ./maestro