# The builder go-image
#FROM golang:1.14-alpine as builderimage
FROM golang:1.17 as builderimage

ENV http_proxy "http://172.18.104.20:1707"
ENV https_proxy "http://172.18.104.20:1707"

RUN mkdir /app && chmod -R 777 /app
RUN touch /app/ejlog-server.log
WORKDIR /app

# Copy and Download all necessary module
COPY go.mod go.sum ./
RUN go mod download 

# Copy All Local Files to Image
COPY . .

# Build Docker Image with CGO Enabled
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
#RUN GOOS=linux GOARCH=amd64 go build -o main

#FROM golang:1.14-alpine
FROM golang:1.17
#FROM alpine:latest

# Set Working Directory in new Image
RUN mkdir /app && chmod -R 777 /app
WORKDIR /app

# Get Executable Binary file to new Image
#COPY --from=builderimage /app/bigetron /app/.env /app/bigetron.log ./
COPY --from=builderimage /app/ejlog /app/.env /app/ejlog-server.log ./

# Expose port 3000 to the outside world
EXPOSE 3000/tcp

ENV http_proxy ""
ENV https_proxy ""

# Run the server executable
CMD [ "./main serve" ]