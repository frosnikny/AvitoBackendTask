FROM golang:1.23-alpine AS builder

RUN apk add --update make git curl

ARG MODULE_NAME=backend

COPY ../zadanie-6105 /home/${MODULE_NAME}/

WORKDIR /home/${MODULE_NAME}/

RUN go build -o main cmd/app/main.go

RUN go build -o migrate cmd/migrate/new_models.go

FROM alpine:latest as production

ARG BUILDER_MODULE_NAME=backend

WORKDIR /root/

COPY --from=builder /home/${BUILDER_MODULE_NAME}/main .
COPY --from=builder /home/${BUILDER_MODULE_NAME}/migrate .

RUN chown root:root main migrate

CMD ["sh", "-c", "./migrate && ./main"]
