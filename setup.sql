CREATE TABLE "person" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "s" STRING NOT NULL,
    "i" INT NOT NULL,
	"d" TIMESTAMP NOT NULL,
	"f32" FLOAT NOT NULL,
	"f64" DECIMAL NOT NULL,
    CONSTRAINT "primary" PRIMARY KEY ("id" ASC)
);

CREATE TABLE "pet" (
    "pid" UUID NOT NULL,
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "name" STRING NOT NULL,
    CONSTRAINT "primary" PRIMARY KEY ("pid" ASC, "id" ASC),
    CONSTRAINT fk_pid_ref_person FOREIGN KEY ("pid") REFERENCES person ("id")
) INTERLEAVE IN PARENT person ("pid");