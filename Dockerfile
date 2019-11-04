FROM debian:buster-slim
WORKDIR /app
RUN apt update
RUN apt install -y ca-certificates
RUN apt install openssl
RUN update-ca-certificates --fresh
ARG listen=0.0.0.0
ARG port=8088
ARG tlskey=""
ARG tlscert=""
ENV listen=$listen
ENV port=$port
ENV tlskey=$tlskey
ENV tlscert=$tlscert
COPY builds/linux/amd64/getsum ./ 
CMD ls -laZ && /app/getsum -s -l $listen -p $port -dir /tmp -tk ""$tlskey -tc ""$tlscert
EXPOSE $port
