CREATE TABLE owner (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    name STRING NOT NULL,
    date_of_birth TIMESTAMP NOT NULL,
    CONSTRAINT "primary" PRIMARY KEY (id ASC)
);

CREATE TABLE pet (
    pid UUID NOT NULL,
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    name STRING NOT NULL,
    owner_name STRING NOT NULL,
    CONSTRAINT "primary" PRIMARY KEY (pid ASC, id ASC),
    CONSTRAINT fk_pid_ref_owner FOREIGN KEY (pid) REFERENCES owner (id)
) INTERLEAVE IN PARENT owner (pid);