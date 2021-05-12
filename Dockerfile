FROM golang:1.8 as builder
WORKDIR /lalash
COPY . /lalash
RUN make

FROM ubuntu
COPY --from=builder /lalash /lalash
RUN useradd -ms /lalash/lalash alice
USER alice
WORKDIR /home/alice
CMD /lalash/lalash
