
# Stage 1
FROM golang:1.17-alpine3.16 AS builder
LABEL stage=builder 
ENV GO111MODULE=on
WORKDIR /app 
COPY . .
RUN apk add build-base && go build -o main ./cmd/main.go

# stage 2
FROM alpine:3.16 AS runner 
LABEL stage=runner 
LABEL maintainer="Made by AmayevArtyom && Mr.RobotDumbazz"
LABEL org.label-schema.description="Docker image for Forum"
WORKDIR /app
COPY --from=builder /app/main ./
COPY /internal /app/internal
COPY /static /app/static 
COPY /templates /app/templates
COPY Forum.db /app/
EXPOSE 8181 
CMD ["./main"]  
