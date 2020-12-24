FROM registry.gitlab.pztrn.name/containers/mirror/golang:1.15.5-alpine AS build

WORKDIR /go/src/gitlab.pztrn.name/pztrn/metricator
COPY . .

ARG CI_PROJECT_NAME
ARG CI_COMMIT_SHA
ARG CI_COMMIT_REF_NAME
ARG CI_COMMIT_TAG

ENV CGO_ENABLED=0
RUN apk add bash git make
RUN make metricatord-build
RUN make metricator-client-build

FROM registry.gitlab.pztrn.name/containers/mirror/golang:1.15.5-alpine
LABEL maintainer="Stanislav N. <pztrn@pztrn.name>"

COPY --from=build /go/src/gitlab.pztrn.name/pztrn/metricator/._bin/metricatord /usr/local/bin/metricatord
COPY --from=build /go/src/gitlab.pztrn.name/pztrn/metricator/._bin/metricator-client /usr/local/bin/metricator-client

RUN apk add tzdata

ENTRYPOINT [ "/usr/local/bin/metricatord", "-config", "/config.yaml" ]
