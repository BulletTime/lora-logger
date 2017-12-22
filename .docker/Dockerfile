FROM karalabe/xgo-latest

MAINTAINER Sven Agneessens <sven.agneessens@gmail.com>

# Get libpcap binaries for linux
RUN \
	dpkg --add-architecture armel && \
	apt-get update && \
	apt-get install -y libpcap0.8-dev

# Get libpcap binaries for armel
ENV ARMLIBPCAP_URL http://ftp.de.debian.org/debian/pool/main/libp/libpcap/libpcap0.8-dev_1.8.1-5_armel.deb
RUN \
  mkdir -p /libpcap && \
	./fetch.sh $ARMLIBPCAP_URL 4c8cd05c693fcc6e5aeccd189a2222b66ce7b1a8 && \
  dpkg -x libpcap0.8-dev_*_armel.deb /libpcap/armel && \
	rm libpcap0.8-dev*.deb
