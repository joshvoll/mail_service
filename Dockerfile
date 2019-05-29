#BUILD STAGE
FROM golang:alpine AS builder

#add label to be able to filter temporary images
LABEL stage=builder

#set working directory outside gopath
WORKDIR /mailer
COPY . .
RUN apk add --no-cache git && \
    go get -d -v ./...

RUN go build -o ./mailer

# --------------------------------------------------

#FINAL STAGE
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /mailer

#copy the whole project folder from temporary image and set entrypoint
COPY --from=builder /mailer .
ENTRYPOINT ./mailer
LABEL Name=mailer Version=1.0
EXPOSE 5010


## build image and remove temporary images:
#docker build -t mailer . && docker image prune --filter label=stage=builder -f