FROM debian

WORKDIR /app
COPY ms-go-example ./
COPY profiles/default.env ./profiles/
COPY migrations/*.sql ./migrations/

EXPOSE 80

CMD [ "./ms-go-example" ]
