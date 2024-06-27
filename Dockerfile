# Build
FROM golang:alpine as builder

LABEL stage=gobuilder
RUN apk update --no-cache && apk add --no-cache tzdata
 
WORKDIR /build

# ADD crontab /etc/cron.d/crontab
# ADD script.sh /script.sh
# COPY entry.sh /entry.sh
# RUN chmod 755 /script.sh /entry.sh

# Add crontab file
COPY crontab /crontab

# Give execution rights on the cron job
RUN chmod 0644 /crontab
RUN /usr/bin/crontab /crontab


ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .
RUN go build -o /app/openserp .


FROM zenika/alpine-chrome:with-chromedriver
USER root

COPY --from=builder /app/openserp /usr/local/bin/openserp
ADD config.yaml /usr/src/app
COPY crontab /crontab

# Give execution rights on the cron job
RUN chmod 0644 /crontab
RUN /usr/bin/crontab /crontab

ENTRYPOINT ["openserp", "serve"]

