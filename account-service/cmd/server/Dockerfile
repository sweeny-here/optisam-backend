FROM ubuntu

WORKDIR /usr/local/bin

ADD server .
ADD cert.pem .
ADD config-local.toml .
#RUN chmod +x ./server
CMD [ "./server" ]
#CMD [ "go", "run" ]
