FROM golang:1.8 as builder
WORKDIR /lalash
COPY . /lalash
RUN make

FROM ubuntu
WORKDIR /lalash
COPY --from=builder /lalash /lalash
CMD ["./lalash"]
