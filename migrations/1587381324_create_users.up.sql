CREATE TABLE "farmers"(
    "id" SERIAL NOT NULL,
    "fname" TEXT NOT NULL,
    "lname" TEXT NOT NULL,
    "email" TEXT NOT NULL UNIQUE,
    "phone" VARCHAR(10) NOT NULL UNIQUE,
    "address" TEXT NOT NULL,
    "password" TEXT NOT NULL
);
ALTER TABLE
    "farmers" ADD PRIMARY KEY("id");
CREATE TABLE "machines"(
    "id" SERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "base_hourly_charge" BIGINT NOT NULL,
    "owner_id" BIGINT NOT NULL
);
ALTER TABLE
    "machines" ADD PRIMARY KEY("id");
CREATE TABLE "bookings"(
    "id" SERIAL NOT NULL,
    "farmer_id" BIGINT NOT NULL,
    "machine_id" BIGINT NOT NULL
);
ALTER TABLE
    "bookings" ADD PRIMARY KEY("id");
CREATE TABLE "invoices"(
    "id" SERIAL NOT NULL,
    "booking_id" BIGINT NOT NULL,
    "date_generated" DATE NOT NULL,
    "total_amount" BIGINT NOT NULL
);
ALTER TABLE
    "invoices" ADD PRIMARY KEY("id");
CREATE TABLE "slots_booked"(
    "id" SERIAL NOT NULL,
    "booking_id" BIGINT NOT NULL,
    "slot_id" BIGINT NOT NULL,
    "date" DATE NOT NULL
);
ALTER TABLE
    "slots_booked" ADD PRIMARY KEY("id");

ALTER TABLE
    "invoices" ADD CONSTRAINT "invoices_booking_id_foreign" FOREIGN KEY("booking_id") REFERENCES "bookings"("id");
ALTER TABLE
    "slots_booked" ADD CONSTRAINT "slots_booked_booking_id_foreign" FOREIGN KEY("booking_id") REFERENCES "bookings"("id");
ALTER TABLE
    "bookings" ADD CONSTRAINT "bookings_farmer_id_foreign" FOREIGN KEY("farmer_id") REFERENCES "farmers"("id");
ALTER TABLE
    "machines" ADD CONSTRAINT "machines_owner_id_foreign" FOREIGN KEY("owner_id") REFERENCES "farmers"("id");
ALTER TABLE
    "bookings" ADD CONSTRAINT "bookings_machine_id_foreign" FOREIGN KEY("machine_id") REFERENCES "machines"("id");
