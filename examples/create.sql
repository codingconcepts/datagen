CREATE DATABASE "sandbox";

USE "sandbox";

CREATE TABLE "owner" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "email" STRING NOT NULL,
    "date_of_birth" TIMESTAMP NOT NULL
);

CREATE TABLE "pet" (
    "id" UUID DEFAULT gen_random_uuid(),
    "pid" UUID NOT NULL,
    "name" STRING NOT NULL,
    "type" STRING NOT NULL,
    PRIMARY KEY ("pid", "id"),
    CONSTRAINT type_v1 CHECK ("type" IN ('cat', 'dog'))
) INTERLEAVE IN PARENT "owner" ("pid");