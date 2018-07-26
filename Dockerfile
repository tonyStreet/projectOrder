FROM golang:latest 
RUN mkdir /app 
ADD project-order /app/ 
WORKDIR /app 
CMD ["/app/project-order"]