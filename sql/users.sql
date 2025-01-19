CREATE TABLE users_next (
  id UUID PRIMARY KEY,
  name VARCHAR (50) UNIQUE NOT NULL,
  timeslots VARCHAR[]
);