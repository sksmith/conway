# GoReleaser builds the binaries, so we just need to copy them
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binaries (GoReleaser provides these in build context)
COPY conway-basic conway-advanced ./

# Make binaries executable
RUN chmod +x conway-basic conway-advanced

CMD ["./conway-basic"]