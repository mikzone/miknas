ARG GOLANG_VERSION=1.22
FROM golang:${GOLANG_VERSION}-bullseye AS builder

ARG VIPS_VERSION=8.14.2
ARG CGIF_VERSION=0.3.0
ARG LIBSPNG_VERSION=0.7.3
ARG TARGETARCH

ENV PKG_CONFIG_PATH=/usr/local/lib/pkgconfig

# Installs libvips + required libraries
RUN DEBIAN_FRONTEND=noninteractive \
  apt-get update && \
  apt-get install --no-install-recommends -y \
  ca-certificates \
  automake build-essential curl \
  python3-pip ninja-build pkg-config \
  gobject-introspection gtk-doc-tools libglib2.0-dev libjpeg62-turbo-dev libpng-dev \
  libwebp-dev libtiff5-dev libexif-dev libxml2-dev libpoppler-glib-dev \
  swig libpango1.0-dev libmatio-dev libopenslide-dev libcfitsio-dev libopenjp2-7-dev liblcms2-dev \
  libgsf-1-dev fftw3-dev liborc-0.4-dev librsvg2-dev libimagequant-dev libheif-dev && \
  pip3 install meson && \
  cd /tmp && \
    curl -fsSLO https://github.com/dloebl/cgif/archive/refs/tags/V${CGIF_VERSION}.tar.gz && \
    tar xf V${CGIF_VERSION}.tar.gz && \
    cd cgif-${CGIF_VERSION} && \
    meson build --prefix=/usr/local --libdir=/usr/local/lib --buildtype=release && \
    cd build && \
    ninja && \
    ninja install && \
  cd /tmp && \
    curl -fsSLO https://github.com/randy408/libspng/archive/refs/tags/v${LIBSPNG_VERSION}.tar.gz && \
    tar xf v${LIBSPNG_VERSION}.tar.gz && \
    cd libspng-${LIBSPNG_VERSION} && \
    meson setup _build \
      --buildtype=release \
      --strip \
      --prefix=/usr/local \
      --libdir=lib && \
    ninja -C _build && \
    ninja -C _build install && \
  cd /tmp && \
    curl -fsSLO https://github.com/libvips/libvips/releases/download/v${VIPS_VERSION}/vips-${VIPS_VERSION}.tar.xz && \
    tar xf vips-${VIPS_VERSION}.tar.xz && \
    cd vips-${VIPS_VERSION} && \
    meson setup _build \
    --buildtype=release \
    --strip \
    --prefix=/usr/local \
    --libdir=lib \
    -Dgtk_doc=false \
    -Dmagick=disabled \
    -Dintrospection=false && \
    ninja -C _build && \
    ninja -C _build install && \
  ldconfig && \
  rm -rf /usr/local/lib/python* && \
  rm -rf /usr/local/lib/libvips-cpp.* && \
  rm -rf /usr/local/lib/*.a && \
  rm -rf /usr/local/lib/*.la

# Cache go modules
ENV GO111MODULE=on

WORKDIR ${GOPATH}/src/github.com/mikzone/miknas

COPY go.work .
COPY go.work.sum .

WORKDIR ${GOPATH}/src/github.com/mikzone/miknas

COPY server/go.mod server/go.mod
COPY server/go.sum server/go.sum
COPY app/server/go.mod app/server/go.mod

# RUN go work sync
RUN go mod download

# Copy imaginary sources
COPY server server
COPY app/server app/server

WORKDIR ${GOPATH}/src/github.com/mikzone/miknas/app/server

RUN go build \
    -o ${GOPATH}/bin/miknas_server \
    -ldflags="-s -w" \
    main.go

FROM debian:bullseye-slim

COPY --from=builder /usr/local/lib /usr/local/lib
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

# Install runtime dependencies
RUN DEBIAN_FRONTEND=noninteractive \
  apt-get update && \
  apt-get install --no-install-recommends -y \
  procps libglib2.0-0 libjpeg62-turbo libpng16-16 libopenexr25 \
  libwebp6 libwebpmux3 libwebpdemux2 libtiff5 libexif12 libxml2 libpoppler-glib8 \
  libpango1.0-0 libmatio11 libopenslide0 libopenjp2-7 libjemalloc2 \
  libgsf-1-114 fftw3 liborc-0.4-0 librsvg2-2 libcfitsio9 libimagequant0 libheif1 && \
  ln -s /usr/lib/$(uname -m)-linux-gnu/libjemalloc.so.2 /usr/local/lib/libjemalloc.so && \
  apt-get autoremove -y && \
  apt-get autoclean && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=builder /go/bin/miknas_server /usr/local/bin/miknas_server

ENV VIPS_WARNING=0
ENV GIN_MODE=release
ENV LD_PRELOAD=/usr/local/lib/libjemalloc.so

RUN groupadd -r --gid 1000 miknas && useradd --no-log-init -r -g 1000 -u 1000 miknas && \
  mkdir -p /web/workspace && mkdir -p /web/config && chmod -R 777 /web
USER miknas

ENV PORT=2020
WORKDIR /web
ENTRYPOINT ["/usr/local/bin/miknas_server"]
EXPOSE ${PORT}
