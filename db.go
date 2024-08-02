package main

import (
	"job_lookup/db"
	"os"
	"time"

	"context"
	"database/sql"
	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var ddl string

const DB_FILE = "data.db"

func MakeDb() (context.Context, *db.Queries) {
	isTest := IsTest()
	if isTest {
		os.Remove(DB_FILE)
	}

	ctx := context.Background()

	dbConn, err := sql.Open("sqlite3", DB_FILE)

	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	if _, err := dbConn.ExecContext(ctx, "PRAGMA foreign_keys = ON"); err != nil {
		panic("failed to enable foreign keys: " + err.Error())
	}

	rows, err := dbConn.Query("SELECT 1 FROM sqlite_master WHERE type='table' AND name='company';")
	if err != nil {
		panic("failed to check if migration is needed: " + err.Error())
	}
	defer rows.Close()
	applyFixtures := false
	if !rows.Next() {
		applyFixtures = true
		if _, err := dbConn.ExecContext(ctx, ddl); err != nil {
			panic("failed to migrate: " + err.Error())
		}
	}

	queries := db.New(dbConn)

	if isTest && applyFixtures {
		setFixtures(ctx, queries)
	}

	return ctx, queries
}

func unwrapFixture[T interface{}](res T, err error) T {
	if err != nil {
		panic("fixture failled: " + err.Error())
	}
	return res
}

func setFixtures(ctx context.Context, queries *db.Queries) {
	ct1 := unwrapFixture(queries.InsertCompanyType(ctx,
		"type A"))
	ct2 := unwrapFixture(queries.InsertCompanyType(ctx,
		"type B"))

	c1 := unwrapFixture(queries.InsertCompany(ctx, db.InsertCompanyParams{
		Name: "company 1", Notes: "company 1 yada yada"}))
	c2 := unwrapFixture(queries.InsertCompany(ctx, db.InsertCompanyParams{
		Name: "company 2"}))

	unwrapFixture(queries.InsertCompanyTypeRel(ctx, db.InsertCompanyTypeRelParams{
		CompanyID: c1.ID, CompanyTypeID: ct1.ID}))
	unwrapFixture(queries.InsertCompanyTypeRel(ctx, db.InsertCompanyTypeRelParams{
		CompanyID: c1.ID, CompanyTypeID: ct2.ID}))
	unwrapFixture(queries.InsertCompanyTypeRel(ctx, db.InsertCompanyTypeRelParams{
		CompanyID: c2.ID, CompanyTypeID: ct1.ID}))

	con1 := unwrapFixture(queries.InsertContact(ctx, db.InsertContactParams{
		CompanyID: c1.ID, FistName: "John", LastName: "Doe", JobPosition: "manager",
	}))
	con2 := unwrapFixture(queries.InsertContact(ctx, db.InsertContactParams{
		CompanyID: c2.ID, FistName: "Jane", LastName: "Dae", JobPosition: "dev", Email: "jane.dae@something.com", PhoneNumber: "561-555-7689", Notes: "contact 2 yada yada",
	}))
	con3 := unwrapFixture(queries.InsertContact(ctx, db.InsertContactParams{
		CompanyID: c2.ID, FistName: "June", LastName: "Duh", JobPosition: "tech lead",
	}))

	jas1 := unwrapFixture(queries.InsertJobApplicationStatus(ctx,
		"refused"))
	jas2 := unwrapFixture(queries.InsertJobApplicationStatus(ctx,
		"pending"))
	jas3 := unwrapFixture(queries.InsertJobApplicationStatus(ctx,
		"accepted"))

	j1 := unwrapFixture(queries.InsertJobApplication(ctx, db.InsertJobApplicationParams{
		CompanyID: c1.ID, JobTitle: "job 1", StatusID: jas1.ID, Notes: "job 1 yada yada"}))
	j2 := unwrapFixture(queries.InsertJobApplication(ctx, db.InsertJobApplicationParams{
		CompanyID: c2.ID, JobTitle: "job 2", StatusID: jas2.ID}))
	j3 := unwrapFixture(queries.InsertJobApplication(ctx, db.InsertJobApplicationParams{
		CompanyID: c2.ID, JobTitle: "job 3", StatusID: jas3.ID}))

	es1 := unwrapFixture(queries.InsertEventSource(ctx,
		"email"))
	es2 := unwrapFixture(queries.InsertEventSource(ctx,
		"phone"))

	now := time.Now()
	loc := now.Local().Location()
	// event.date is just a date in the application but go does not seems to manage DST,
	// so setting an hour in the middle of the day avoid some confusion on the local tests
	baseDate := time.Date(now.Year(), now.Month(), 1, 15, 0, 0, 0, loc)
	baseDate1 := baseDate.AddDate(0, 0, 1)
	baseDate2 := baseDate1.AddDate(0, 0, 1)
	e1 := unwrapFixture(queries.InsertEvent(ctx, db.InsertEventParams{
		JobApplicationID: j1.ID, SourceID: es1.ID, Title: "first interview", Date: baseDate, Notes: "event 1 yada yada"}))
	e2 := unwrapFixture(queries.InsertEvent(ctx, db.InsertEventParams{
		JobApplicationID: j2.ID, SourceID: es2.ID, Title: "first interview", Date: baseDate, Notes: "event 2 yada yada"}))
	e3 := unwrapFixture(queries.InsertEvent(ctx, db.InsertEventParams{
		JobApplicationID: j3.ID, SourceID: es1.ID, Title: "first interview", Date: baseDate1, Notes: "event 3 yada yada"}))
	e4 := unwrapFixture(queries.InsertEvent(ctx, db.InsertEventParams{
		JobApplicationID: j3.ID, SourceID: es2.ID, Title: "second interview", Date: baseDate2, Notes: "event 4 yada yada"}))

	unwrapFixture(queries.InsertEventContact(ctx, db.InsertEventContactParams{
		EventID: e1.ID, ContactID: con1.ID}))
	unwrapFixture(queries.InsertEventContact(ctx, db.InsertEventContactParams{
		EventID: e2.ID, ContactID: con2.ID}))
	unwrapFixture(queries.InsertEventContact(ctx, db.InsertEventContactParams{
		EventID: e3.ID, ContactID: con2.ID}))
	unwrapFixture(queries.InsertEventContact(ctx, db.InsertEventContactParams{
		EventID: e4.ID, ContactID: con2.ID}))
	unwrapFixture(queries.InsertEventContact(ctx, db.InsertEventContactParams{
		EventID: e4.ID, ContactID: con3.ID}))
}
