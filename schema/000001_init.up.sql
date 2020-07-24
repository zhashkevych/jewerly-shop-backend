CREATE TABLE "users" (
    "id" serial NOT NULL,
    "email" varchar(255) NOT NULL UNIQUE,
    "password_hash" varchar(255) NOT NULL,
    "first_name" varchar(255) NOT NULL,
    "last_name" varchar(255) NOT NULL,
    "registered_at" TIMESTAMP NOT NULL DEFAULT NOW()
)