package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx"
)

func Connection(host, uname, pass string) *pgx.Conn {

	config := pgx.ConnConfig{
		Host:     host,
		User:     uname,
		Password: pass,
		Database: "postgres",
	}
	db, err := pgx.Connect(config)
	if err != nil {
		panic(err)
	}

	// defer db.Close()
	err = db.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	return db
}

func CreateTables(db *pgx.Conn) error {

	// create table if not exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS public.users (
		id UUID PRIMARY KEY,
		name VARCHAR (50) UNIQUE NOT NULL
	  );`)
	if err != nil {
		log.Println("Error creating table: ", err)
		return err
	}

	// create table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS public.time_slots(
		id uuid NOT NULL,
		user_id uuid NOT NULL,
		time_slot character varying NOT NULL,
		PRIMARY KEY (id),
		CONSTRAINT user_id_foreign_key FOREIGN KEY (user_id)
			REFERENCES public.users (id) MATCH SIMPLE
			ON UPDATE NO ACTION
			ON DELETE NO ACTION
			NOT VALID
	);`)
	if err != nil {
		log.Println("Error creating table: ", err)
		return err
	}

	// create table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS public.events
	(
		id uuid NOT NULL,
		event_owner uuid NOT NULL,
		title character varying NOT NULL,
		event_start_time timestamp with time zone NOT NULL,
		event_end_time timestamp with time zone NOT NULL,
		participants character varying[] NOT NULL,
		PRIMARY KEY (id),
		CONSTRAINT event_owner_user_id__foreign_key FOREIGN KEY (event_owner)
			REFERENCES public.users (id) MATCH SIMPLE
			ON UPDATE NO ACTION
			ON DELETE NO ACTION
			NOT VALID
	);`)
	if err != nil {
		log.Println("Error creating table: ", err)
		return err
	}
	return nil
}
