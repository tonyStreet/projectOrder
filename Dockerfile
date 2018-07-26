FROM golang:latest 
RUN mkdir /app 
ADD project-order /app/
ADD config.yml /app/
WORKDIR /app 
CMD ["/app/project-order"]