# Build stage
FROM golang:latest

RUN mkdir /build

WORKDIR /build

RUN export GO111MODULE=on

# Clone the repository
RUN go install github.com/fegig/goLang_class@latest
RUN cd /build/ && git clone https://github.com/fegig/goLang_class.git

# Change directory to the cloned repository
RUN cd /build/goLang_class && go build

# Expose port 8080
EXPOSE 8080

# Run the application
ENTRYPOINT ["./build/goLang_class/main"] 