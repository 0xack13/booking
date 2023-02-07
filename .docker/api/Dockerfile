#STEP 1/2
FROM golang:1.19 AS build

ENV GO111MODULE=on
ENV SERVICE=api
ENV APP=booking
ENV GOMODCACHE=/root/.cache/go

WORKDIR /go/src/github.com/carlos/${APP}

COPY go.mod ./

COPY . .
RUN go build -o booking-service-${SERVICE} /go/src/github.com/carlos/${APP}/cmd/${SERVICE}

CMD ["./booking-service-api"]

#STEP 2/2
FROM gcr.io/distroless/base
ENV APP=booking-service-api

WORKDIR /home/nonroot
USER nonroot

COPY --chown=nonroot --from=build /go/src/github.com/carlos/${APP}/${APP} /home/nonroot/
CMD ["./booking-service-api"]
