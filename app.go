package main

import (
	"context"
	"database/sql"
	"fmt"
	"job_lookup/db"
	"strings"
	"time"
)

type App struct {
	appctx  context.Context
	dbctx   context.Context
	queries *db.Queries
	db      *sql.DB
}

func NewApp() *App {
	c, q, db := MakeDb()
	return &App{
		dbctx:   c,
		queries: q,
		db:      db,
	}
}

func (a *App) startup(ctx context.Context) {
	a.appctx = ctx
}

func groupBy[R interface{}, G interface{}, K comparable](
	data []R,
	groupKey func(*R) K,
	setGroupFromRaw func([]R, *G),
) []G {
	// does not scale well but that should be ok since there's not that much data
	// and the function keep db order
	treadtedSet := make(map[K]bool)
	grouped := make([]G, 0)
	l := len(data)
	for i := range data {
		item := data[i]
		id := groupKey(&item)
		if treadtedSet[id] {
			continue
		}
		lines := []R{item}
		for j := i + 1; j < l; j++ {
			nextLine := data[j]
			if id == groupKey(&nextLine) {
				lines = append(lines, nextLine)
			}
		}
		var groupedItem G
		setGroupFromRaw(lines, &groupedItem)
		grouped = append(grouped, groupedItem)
		treadtedSet[id] = true
	}
	return grouped
}

func wrapError(err error, operation string) error {
	if err == nil {
		return nil
	}
	newErr := fmt.Errorf("%s: %w", operation, err)
	fmt.Printf("ERROR: %v\n", newErr)
	return newErr
}

func (a *App) ListCompanyTypes() ([]db.ListCompanyTypeRow, error) {
	data, err := a.queries.ListCompanyType(a.dbctx)
	return data, wrapError(err, "ListCompanyTypes")
}

func (a *App) DeleteCompanyType(item db.ListCompanyTypeRow) error {
	return wrapError(a.queries.DeleteCompanyType(a.dbctx, item.Name), "DeleteCompanyType")
}

func (a *App) UpdateCompanyType(item db.ListCompanyTypeRow) error {
	return wrapError(a.queries.UpdateCompanyType(a.dbctx, db.UpdateCompanyTypeParams{
		ID:   item.ID,
		Name: item.Name,
	}), "UpdateCompanyType")
}

func (a *App) InsertCompanyType(item db.ListCompanyTypeRow) (db.ListCompanyTypeRow, error) {
	c, err := a.queries.InsertCompanyType(a.dbctx, item.Name)
	item.ID = c.ID
	return item, wrapError(err, "InsertCompanyType")
}

func timpestampToDate(timestamp int64) *time.Time {
	if timestamp == 0 {
		return nil
	}
	date := time.Unix(timestamp, 0)
	return &date
}

type Company struct {
	ID             int64      `json:"id"`
	Name           string     `json:"name"`
	Notes          string     `json:"notes"`
	CompanyTypes   []string   `json:"company_types"`
	ApplicationCnt int64      `json:"application_cnt"`
	LastEvent      *time.Time `json:"last_event"`
	NextEvent      *time.Time `json:"next_event"`
}

func (a *App) ListCompanies() ([]Company, error) {
	c, err := a.queries.ListCompany(a.dbctx)
	if err != nil {
		return nil, wrapError(err, "ListCompanies")
	}
	return groupBy(
		c,
		func(line *db.ListCompanyRow) int64 { return line.ID },
		func(r []db.ListCompanyRow, g *Company) {
			g.ID = r[0].ID
			g.Name = r[0].Name
			g.Notes = r[0].Notes
			g.ApplicationCnt = r[0].ApplicationCnt
			g.LastEvent = timpestampToDate(r[0].LastEvent)
			g.NextEvent = timpestampToDate(r[0].NextEvent)
			for _, sub := range r {
				g.CompanyTypes = append(g.CompanyTypes, sub.CompanyType.Name)
			}
		},
	), nil
}

func (a *App) DeleteCompany(item Company) error {
	return wrapError(a.queries.DeleteCompany(a.dbctx, item.Name), "DeleteCompany")
}

func setCompanyType(a *App, queries *db.Queries, companyId int64, companyTypes []string) error {
	err := queries.DeleteCompanyTypeRel(a.dbctx, companyId)
	if err != nil {
		return wrapError(err, "setCompanyType/DeleteCompanyTypeRel")
	}
	for _, ct := range companyTypes {
		dbCt, err := queries.GetCompanyType(a.dbctx, ct)
		if err != nil {
			return wrapError(err, "setCompanyType/GetCompanyType")
		}
		_, err = queries.InsertCompanyTypeRel(a.dbctx, db.InsertCompanyTypeRelParams{
			CompanyID:     companyId,
			CompanyTypeID: dbCt.ID,
		})
		if err != nil {
			return wrapError(err, "setCompanyType/InsertCompanyTypeRel")
		}
	}
	return nil
}

func (a *App) UpdateCompany(item Company) error {
	tx, err := a.db.Begin()
	if err != nil {
		return wrapError(err, "UpdateCompany/Begin")
	}
	defer tx.Rollback()
	queries := a.queries.WithTx(tx)
	err = queries.UpdateCompany(a.dbctx, db.UpdateCompanyParams{
		ID:    item.ID,
		Name:  item.Name,
		Notes: item.Notes,
	})
	if err != nil {
		return wrapError(err, "UpdateCompany")
	}
	err = setCompanyType(a, queries, item.ID, item.CompanyTypes)
	if err != nil {
		return wrapError(err, "UpdateCompany/setCompanyType")
	}
	return wrapError(tx.Commit(), "UpdateCompany/Commit")
}

func (a *App) InsertCompany(item Company) (Company, error) {
	tx, err := a.db.Begin()
	if err != nil {
		return Company{}, wrapError(err, "InsertCompany/Begin")
	}
	defer tx.Rollback()
	queries := a.queries.WithTx(tx)

	c, err := queries.InsertCompany(a.dbctx, db.InsertCompanyParams{
		Name:  item.Name,
		Notes: item.Notes,
	})
	if err != nil {
		return Company{}, wrapError(err, "InsertCompany")
	}
	err = setCompanyType(a, queries, c.ID, item.CompanyTypes)
	if err != nil {
		return Company{}, wrapError(err, "InsertCompany/setCompanyType")
	}
	item.ID = c.ID
	err = tx.Commit()
	if err != nil {
		return Company{}, wrapError(tx.Commit(), "InsertCompany/Commit")
	}
	return item, nil
}

func (a *App) ListJobApplicationStatus() ([]db.ListJobApplicationStatusRow, error) {
	data, err := a.queries.ListJobApplicationStatus(a.dbctx)
	return data, wrapError(err, "ListJobApplicationStatus")
}
func (a *App) InsertJobApplicationStatus(item db.ListJobApplicationStatusRow) (db.ListJobApplicationStatusRow, error) {
	data, err := a.queries.InsertJobApplicationStatus(a.dbctx, item.Name)
	item.ID = data.ID
	return item, wrapError(err, "InsertJobApplicationStatus")
}
func (a *App) DeleteJobApplicationStatus(item db.ListJobApplicationStatusRow) error {
	return wrapError(a.queries.DeleteJobApplicationStatus(a.dbctx, item.Name), "DeleteJobApplicationStatus")
}
func (a *App) UpdateJobApplicationStatus(item db.ListJobApplicationStatusRow) error {
	return wrapError(a.queries.UpdateJobApplicationStatus(a.dbctx, db.UpdateJobApplicationStatusParams{
		Name: item.Name,
		ID:   item.ID,
	}), "UpdateJobApplicationStatus")
}

func (a *App) ListEventSource() ([]db.ListEventSourceRow, error) {
	data, err := a.queries.ListEventSource(a.dbctx)
	return data, wrapError(err, "ListEventSource")
}
func (a *App) InsertEventSource(item db.ListEventSourceRow) (db.ListEventSourceRow, error) {
	data, err := a.queries.InsertEventSource(a.dbctx, item.Name)
	item.ID = data.ID
	return item, wrapError(err, "InsertEventSource")
}
func (a *App) DeleteEventSource(item db.ListEventSourceRow) error {
	return wrapError(a.queries.DeleteEventSource(a.dbctx, item.Name), "DeleteEventSource")
}
func (a *App) UpdateEventSource(item db.ListEventSourceRow) error {
	return wrapError(a.queries.UpdateEventSource(a.dbctx, db.UpdateEventSourceParams{
		Name: item.Name,
		ID:   item.ID,
	}), "UpdateEventSource")
}

type Contact struct {
	ID          int64      `json:"id"`
	CompanyID   int64      `json:"company_id"`
	JobPosition string     `json:"job_position"`
	FistName    string     `json:"fist_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	PhoneNumber string     `json:"phone_number"`
	Notes       string     `json:"notes"`
	CompanyName string     `json:"company_name"`
	LastEvent   *time.Time `json:"last_event"`
	NextEvent   *time.Time `json:"next_event"`
}

func (a *App) ListContact() ([]Contact, error) {
	res, err := a.queries.ListContact(a.dbctx)
	if err != nil {
		return nil, wrapError(err, "ListContact")
	}
	contacts := make([]Contact, len(res))
	for i := range res {
		dbContact := res[i]
		contacts[i] = Contact{
			ID:          dbContact.ID,
			CompanyID:   dbContact.CompanyID,
			JobPosition: dbContact.JobPosition,
			FistName:    dbContact.FistName,
			LastName:    dbContact.LastName,
			Email:       dbContact.Email,
			PhoneNumber: dbContact.PhoneNumber,
			Notes:       dbContact.Notes,
			CompanyName: dbContact.CompanyName,
			LastEvent:   timpestampToDate(dbContact.LastEvent),
			NextEvent:   timpestampToDate(dbContact.NextEvent),
		}
	}
	return contacts, nil
}
func (a *App) InsertContact(item Contact) (Contact, error) {
	c, err := a.queries.GetCompanyIdByName(a.dbctx, item.CompanyName)
	if err != nil {
		return Contact{}, wrapError(err, "InsertContact/GetCompanyIdByName")
	}
	newitem, err := a.queries.InsertContact(a.dbctx, db.InsertContactParams{
		CompanyID:   c,
		JobPosition: item.JobPosition,
		FistName:    item.FistName,
		LastName:    item.LastName,
		Email:       item.Email,
		PhoneNumber: item.PhoneNumber,
		Notes:       item.Notes,
	})
	if err != nil {
		return Contact{}, wrapError(err, "InsertContact")
	}
	item.ID = newitem.ID
	return item, nil
}
func (a *App) DeleteContact(item Contact) error {
	return wrapError(a.queries.DeleteContact(a.dbctx, item.ID), "DeleteContact")
}
func (a *App) UpdateContact(item Contact) error {
	c, err := a.queries.GetCompanyIdByName(a.dbctx, item.CompanyName)
	if err != nil {
		return wrapError(err, "UpdateContact/GetCompanyIdByName")
	}
	return wrapError(a.queries.UpdateContact(a.dbctx, db.UpdateContactParams{
		CompanyID:   c,
		JobPosition: item.JobPosition,
		FistName:    item.FistName,
		LastName:    item.LastName,
		Email:       item.Email,
		PhoneNumber: item.PhoneNumber,
		Notes:       item.Notes,
		ID:          item.ID,
	}), "UpdateContact")
}

type JobApplication struct {
	ID          int64      `json:"id"`
	CompanyID   int64      `json:"company_id"`
	StatusID    int64      `json:"status_id"`
	JobTitle    string     `json:"job_title"`
	Notes       string     `json:"notes"`
	StatusName  string     `json:"status_name"`
	CompanyName string     `json:"company_name"`
	EventCnt    int64      `json:"event_cnt"`
	LastEvent   *time.Time `json:"last_event"`
	NextEvent   *time.Time `json:"next_event"`
}

func (a *App) ListJobApplication() ([]JobApplication, error) {
	res, err := a.queries.ListJobApplication(a.dbctx)
	if err != nil {
		return nil, wrapError(err, "ListJobApplication")
	}
	applications := make([]JobApplication, len(res))
	for i := range res {
		dbApp := res[i]
		applications[i] = JobApplication{
			ID:          dbApp.ID,
			CompanyID:   dbApp.CompanyID,
			StatusID:    dbApp.StatusID,
			JobTitle:    dbApp.JobTitle,
			Notes:       dbApp.Notes,
			StatusName:  dbApp.StatusName,
			CompanyName: dbApp.CompanyName,
			EventCnt:    dbApp.EventCnt,
			LastEvent:   timpestampToDate(dbApp.LastEvent),
			NextEvent:   timpestampToDate(dbApp.NextEvent),
		}
	}
	return applications, nil
}
func (a *App) InsertJobApplication(item JobApplication) (JobApplication, error) {
	c, err := a.queries.GetCompanyIdByName(a.dbctx, item.CompanyName)
	if err != nil {
		return JobApplication{}, wrapError(err, "InsertJobApplication/GetCompanyIdByName")
	}
	s, err := a.queries.GetJobApplicationStatusIdByName(a.dbctx, item.StatusName)
	if err != nil {
		return JobApplication{}, wrapError(err, "InsertJobApplication/GetJobApplicationStatusIdByName")
	}
	newitem, err := a.queries.InsertJobApplication(a.dbctx, db.InsertJobApplicationParams{
		CompanyID: c,
		StatusID:  s,
		JobTitle:  item.JobTitle,
		Notes:     item.Notes,
	})
	if err != nil {
		return JobApplication{}, wrapError(err, "InsertJobApplication")
	}
	item.ID = newitem.ID
	return item, nil
}
func (a *App) DeleteJobApplication(item JobApplication) error {
	return wrapError(a.queries.DeleteJobApplication(a.dbctx, item.ID), "DeleteJobApplication")
}
func (a *App) UpdateJobApplication(item JobApplication) error {
	c, err := a.queries.GetCompanyIdByName(a.dbctx, item.CompanyName)
	if err != nil {
		return wrapError(err, "UpdateJobApplication/GetCompanyIdByName")
	}
	s, err := a.queries.GetJobApplicationStatusIdByName(a.dbctx, item.StatusName)
	if err != nil {
		return wrapError(err, "UpdateJobApplication/GetJobApplicationStatusIdByName")
	}
	return wrapError(a.queries.UpdateJobApplication(a.dbctx, db.UpdateJobApplicationParams{
		CompanyID: c,
		StatusID:  s,
		JobTitle:  item.JobTitle,
		Notes:     item.Notes,
		ID:        item.ID,
	}), "UpdateJobApplication")
}

type Event struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Date        time.Time `json:"date"`
	Notes       string    `json:"notes"`
	Source      string    `json:"source"`
	JobTitle    string    `json:"job_title"`
	CompanyName string    `json:"company_name"`
	Contacts    []string  `json:"contacts"`
}

func joinContactNames(fistName *string, lastName *string) *string {
	if fistName == nil || lastName == nil {
		return nil
	}
	fmtName := fmt.Sprintf("%s, %s", *fistName, *lastName)
	return &fmtName
}

func splitContactNames(contactNamesJoined string) (firstName string, lastName string) {
	spl := strings.Split(contactNamesJoined, ", ")
	return spl[0], spl[1]
}

func (a *App) ListEvents() ([]Event, error) {
	e, err := a.queries.ListEvent(a.dbctx)
	if err != nil {
		return nil, wrapError(err, "ListEvents")
	}
	return groupBy(
		e,
		func(r *db.ListEventRow) int64 { return r.ID },
		func(r []db.ListEventRow, g *Event) {
			g.ID = r[0].ID
			g.Title = r[0].Title
			g.Date = time.Unix(r[0].Date, 0)
			g.Notes = r[0].Notes
			g.Source = r[0].Source
			g.JobTitle = r[0].JobTitle
			g.CompanyName = r[0].CompanyName
			for _, s := range r {
				fmtName := joinContactNames(s.ContactFistName, s.ContactLastName)
				if fmtName != nil {
					g.Contacts = append(g.Contacts, *fmtName)
				}
			}
		},
	), nil
}

func setEventContacts(a *App, queries *db.Queries, event Event) error {
	err := queries.DeleteEventContact(a.dbctx, event.ID)
	if err != nil {
		return wrapError(err, "setEventContacts/DeleteEventContact")
	}
	for _, c := range event.Contacts {
		firstName, lastName := splitContactNames(c)
		contactId, err := queries.GetContactIdByNames(a.dbctx, db.GetContactIdByNamesParams{
			FistName: firstName,
			LastName: lastName,
			Name:     event.CompanyName,
		})
		if err != nil {
			return wrapError(err, "setEventContacts/GetContactIdByNames")
		}
		_, err = queries.InsertEventContact(a.dbctx, db.InsertEventContactParams{
			EventID:   event.ID,
			ContactID: contactId,
		})
		if err != nil {
			return wrapError(err, "setEventContacts/InsertEventContact")
		}
	}
	return nil
}

func (a *App) InsertEvent(item Event) (Event, error) {
	tx, err := a.db.Begin()
	if err != nil {
		return Event{}, wrapError(err, "InsertEvent/Begin")
	}
	defer tx.Rollback()
	queries := a.queries.WithTx(tx)

	s, err := queries.GetEventSourceIdByName(a.dbctx, item.Source)
	if err != nil {
		return Event{}, wrapError(err, "InsertEvent/GetEventSourceIdByName")
	}
	ja, err := queries.GetJobApplicationIdByName(a.dbctx, db.GetJobApplicationIdByNameParams{
		JobTitle: item.JobTitle,
		Name:     item.CompanyName,
	})
	if err != nil {
		return Event{}, wrapError(err, "InsertEvent/GetJobApplicationIdByName")
	}
	newitem, err := queries.InsertEvent(a.dbctx, db.InsertEventParams{
		SourceID:         s,
		JobApplicationID: ja,
		Title:            item.Title,
		Date:             item.Date.Unix(),
		Notes:            item.Notes,
	})
	if err != nil {
		return Event{}, wrapError(err, "InsertEvent")
	}
	item.ID = newitem.ID
	err = setEventContacts(a, queries, item)
	if err != nil {
		return Event{}, wrapError(err, "InsertEvent/setEventContacts")
	}
	err = tx.Commit()
	if err != nil {
		return Event{}, wrapError(err, "InsertEvent/Commit")
	}
	return item, nil
}
func (a *App) DeleteEvent(item Event) error {
	return wrapError(a.queries.DeleteEvent(a.dbctx, item.ID), "DeleteEvent")
}
func (a *App) UpdateEvent(item Event) error {
	tx, err := a.db.Begin()
	if err != nil {
		return wrapError(err, "UpdateEvent/Begin")
	}
	defer tx.Rollback()
	queries := a.queries.WithTx(tx)

	s, err := queries.GetEventSourceIdByName(a.dbctx, item.Source)
	if err != nil {
		return wrapError(err, "UpdateEvent/GetEventSourceIdByName")
	}
	ja, err := queries.GetJobApplicationIdByName(a.dbctx, db.GetJobApplicationIdByNameParams{
		JobTitle: item.JobTitle,
		Name:     item.CompanyName,
	})
	if err != nil {
		return wrapError(err, "UpdateEvent/GetJobApplicationIdByName")
	}
	err = queries.UpdateEvent(a.dbctx, db.UpdateEventParams{
		SourceID:         s,
		JobApplicationID: ja,
		Title:            item.Title,
		Date:             item.Date.Unix(),
		Notes:            item.Notes,
		ID:               item.ID,
	})
	if err != nil {
		return wrapError(err, "UpdateEvent")
	}
	err = setEventContacts(a, queries, item)
	if err != nil {
		return wrapError(err, "UpdateEvent/setEventContacts")
	}
	return wrapError(tx.Commit(), "UpdateEvent/Commit")
}
