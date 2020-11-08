FROM debian:stretch

ENV BAZEL_VERSION=3.7.0
ENV SRVC_NAME=chat_srvc_grpc

# install pre-requisites
RUN apt update && apt install -y                \
    build-essential autoconf libtool pkg-config \
    g++ unzip zip                               \
    git curl gnupg                              \
    apt-transport-https ca-certificates

# install bazel-${BAZEL_VERSION}
RUN curl -fsSL https://bazel.build/bazel-release.pub.gpg | gpg --dearmor > bazel.gpg && \
    mv bazel.gpg /etc/apt/trusted.gpg.d/                                             && \
    echo "deb [arch=amd64] https://storage.googleapis.com/bazel-apt stable jdk1.8" |    \
      tee /etc/apt/sources.list.d/bazel.list                                         && \
    apt update && apt install -y bazel                                               && \
    apt update && apt full-upgrade -y                                                && \
    apt install -y bazel-${BAZEL_VERSION}                                            && \
    ln -fs /usr/bin/bazel-${BAZEL_VERSION} /usr/bin/bazel                            && \
    bazel --version

# initialize grpc service
COPY ./ /${SRVC_NAME}
CMD cd /${SRVC_NAME} && bazel run ${SRVC_NAME}
