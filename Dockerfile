FROM golang as builder
WORKDIR /lalash
COPY . /lalash
RUN make

FROM ubuntu
COPY --from=builder /lalash/lalash /bin
RUN useradd -ms /bin/lalash alice
USER alice
WORKDIR /home/alice
CMD /bin/lalash
