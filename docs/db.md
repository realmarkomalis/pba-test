docker run --rm   --name pg-packback -e POSTGRES_PASSWORD=packback -d -p 5433:5432 -v $HOME/docker/volumes/postgres:/var/lib/postgresql/data  postgres

CREATE EXTENSION IF NOT EXISTS "uuid-ossp ";