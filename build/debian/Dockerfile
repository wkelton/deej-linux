FROM debian:12-slim

WORKDIR /var/temp/deej
RUN chmod 777 .

# Install gnome desktop
RUN apt update && apt install -y gnome-shell

# Install the build dependencies
RUN apt update && apt install -y git golang libayatana-appindicator3-dev && go version

VOLUME install

RUN git config --system --add safe.directory '*'

# Copy git repo for build scripts
COPY .git ./.git
# Copy in the source code
COPY go.mod go.sum build/share/build-all.sh ./
COPY pkg ./pkg

ENTRYPOINT ["./build-all.sh"]
