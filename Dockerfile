FROM golang:1.18 AS builder

# Build go binary
COPY . /go/src/hoppscotch-backend
WORKDIR /go/src/hoppscotch-backend
RUN go build -v

FROM alpine:3.15

# Update
RUN apk update upgrade

# Add libc6 compat, needed to run Go binaries in Alpine
RUN apk add --no-cache libc6-compat

# Set timezone to Europe/Amsterdam
RUN apk add tzdata
RUN ln -s /usr/share/zoneinfo/Europe/Amsterdam /etc/localtime

# Copy Go binary from builder stage to Alpine stage.
COPY --from=builder /go/src/hoppscotch-backend/hoppscotch-backend /usr/bin/

CMD [ "/usr/bin/hoppscotch-backend"]
