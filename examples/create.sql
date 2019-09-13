CREATE DATABASE sandbox;

CREATE TABLE sandbox.owner (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    name STRING NOT NULL,
    date_of_birth TIMESTAMP NOT NULL,
    CONSTRAINT "primary" PRIMARY KEY (id ASC)
);

CREATE TABLE sandbox.pet (
    pid UUID NOT NULL,
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    name STRING NOT NULL,
    owner_name STRING NOT NULL,
    owner_date_of_birth TIMESTAMP NOT NULL,
    CONSTRAINT "primary" PRIMARY KEY (pid ASC, id ASC),
    CONSTRAINT fk_pid_ref_owner FOREIGN KEY (pid) REFERENCES sandbox.owner (id)
) INTERLEAVE IN PARENT sandbox.owner (pid);

CREATE TABLE sandbox.one (
    id int primary key,
    name STRING NOT NULL
);

CREATE TABLE sandbox.two (
    one_id int
);