FROM golang:1.20 as build-env

WORKDIR /go/src/app
ADD . /go/src/app

RUN go test -mod=vendor -cover ./...
RUN go build -mod=vendor -o /go/bin/app


FROM gcr.io/distroless/static:966f4bd97f611354c4ad829f1ed298df9386c2ec
# latest-amd64 -> 966f4bd97f611354c4ad829f1ed298df9386c2ec
# https://github.com/GoogleContainerTools/distroless/tree/master/base

LABEL name="R2 Action"
LABEL repository="https://github.com/alsosee/r2action"
LABEL homepage="https://github.com/alsosee/r2action"

LABEL maintainer="Konstantin Chukhlomin <mail@chuhlomin.com>"
LABEL com.github.actions.name="R2 Action"
LABEL com.github.actions.description="Get an object from CloudFlare R2 bucket"
LABEL com.github.actions.icon="database"
LABEL com.github.actions.color="purple"

COPY --from=build-env /go/bin/app /app

CMD ["/app"]
