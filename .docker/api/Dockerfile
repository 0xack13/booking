#STEP 1/2
FROM golang:1.19 AS build

ENV GO111MODULE=on
ENV SERVICE=api
ENV APP=s4l
ENV GOMODCACHE=/root/.cache/go

WORKDIR /go/src/github.com/carlos/${APP}

COPY go.mod ./

COPY . .
RUN go build -o s4l-service-${SERVICE} /go/src/github.com/carlos/${APP}/cmd/${SERVICE}

CMD ["./s4l-service-api"]

#STEP 2/2
FROM gcr.io/distroless/base
ENV APP=s4l-service-api

WORKDIR /home/nonroot
USER nonroot

COPY --chown=nonroot --from=build /go/src/github.com/carlos/${APP}/${APP} /home/nonroot/
CMD ["./s4l-service-api"]
