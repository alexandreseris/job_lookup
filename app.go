package main

import (
	"context"
	"fmt"
	"job_lookup/db"
	"strings"
	"time"
)

type App struct {
	appctx  context.Context
	dbctx   context.Context
	queries *db.Queries
}

func NewApp() *App {
	c, q := MakeDb()
	return &App{
		dbctx:   c,
		queries: q,
	}
}

func (a *App) startup(ctx context.Context) {
	a.appctx = ctx
}

func (a *App) Log(m string) {
	fmt.Println(m)
}

func groupByMap[R interface{}, G interface{}, K comparable](
	data []R,
	groupKey func(*R) K,
	setGroupFromRaw func([]R, *G),
) []G {
	groupedMap := make(map[K][]R)
	for _, line := range data {
		key := groupKey(&line)
		lookup := groupedMap[key]
		groupedMap[key] = append(lookup, line)
	}
	grouped := make([]G, 0)
	for _, v := range groupedMap {
		var newItem G
		setGroupFromRaw(v, &newItem)
		grouped = append(grouped, newItem)
	}
	return grouped
}

func (a *App) ListCompanyTypes() ([]db.CompanyType, error) {
	return a.queries.ListCompanyType(a.dbctx)
}

func (a *App) DeleteCompanyType(item db.CompanyType) error {
	return a.queries.DeleteCompanyType(a.dbctx, item.Name)
}

func (a *App) UpdateCompanyType(item db.CompanyType) error {
	return a.queries.UpdateCompanyType(a.dbctx, db.UpdateCompanyTypeParams{
		ID:   item.ID,
		Name: item.Name,
	})
}

func (a *App) InsertCompanyType(item db.CompanyType) (db.CompanyType, error) {
	c, err := a.queries.InsertCompanyType(a.dbctx, item.Name)
	return c, err
}

type Company struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Notes          string    `json:"notes"`
	CompanyTypes   []string  `json:"company_types"`
	ApplicationCnt int64     `json:"application_cnt"`
	LastEvent      time.Time `json:"last_event"`
	NextEvent      time.Time `json:"next_event"`
}

func (a *App) ListCompanies() ([]Company, error) {
	c, err := a.queries.ListCompany(a.dbctx)
	if err != nil {
		return nil, err
	}
	return groupByMap(
		c,
		func(line *db.ListCompanyRow) int64 { return line.ID },
		func(r []db.ListCompanyRow, g *Company) {
			g.ID = r[0].ID
			g.Name = r[0].Name
			g.Notes = r[0].Notes
			g.ApplicationCnt = r[0].ApplicationCnt
			g.LastEvent = time.Unix(r[0].LastEvent, 0)
			g.NextEvent = time.Unix(r[0].NextEvent, 0)
			for _, sub := range r {
				g.CompanyTypes = append(g.CompanyTypes, sub.CompanyType.Name)
			}
		},
	), nil
}

func (a *App) DeleteCompany(item Company) error {
	return a.queries.DeleteCompany(a.dbctx, item.Name)
}

func setCompanyType(a *App, companyId int64, companyTypes []string) error {
	err := a.queries.DeleteCompanyTypeRel(a.dbctx, companyId)
	if err != nil {
		return err
	}
	for _, ct := range companyTypes {
		dbCt, err := a.queries.GetCompanyType(a.dbctx, ct)
		if err != nil {
			return err
		}
		_, err = a.queries.InsertCompanyTypeRel(a.dbctx, db.InsertCompanyTypeRelParams{
			CompanyID:     companyId,
			CompanyTypeID: dbCt.ID,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) UpdateCompany(item Company) error {
	err := a.queries.UpdateCompany(a.dbctx, db.UpdateCompanyParams{
		ID:    item.ID,
		Name:  item.Name,
		Notes: item.Notes,
	})
	if err != nil {
		return nil
	}
	return setCompanyType(a, item.ID, item.CompanyTypes)
}

func (a *App) InsertCompany(item Company) (Company, error) {
	c, err := a.queries.InsertCompany(a.dbctx, db.InsertCompanyParams{
		Name:  item.Name,
		Notes: item.Notes,
	})
	if err != nil {
		return Company{}, err
	}
	err = setCompanyType(a, c.ID, item.CompanyTypes)
	if err != nil {
		return Company{}, err
	}
	item.ID = c.ID
	return item, nil
}

func (a *App) ListJobApplicationStatus() ([]db.JobApplicationStatus, error) {
	return a.queries.ListJobApplicationStatus(a.dbctx)
}
func (a *App) InsertJobApplicationStatus(item db.JobApplicationStatus) (db.JobApplicationStatus, error) {
	return a.queries.InsertJobApplicationStatus(a.dbctx, item.Name)
}
func (a *App) DeleteJobApplicationStatus(item db.JobApplicationStatus) error {
	return a.queries.DeleteJobApplicationStatus(a.dbctx, item.Name)
}
func (a *App) UpdateJobApplicationStatus(item db.JobApplicationStatus) error {
	return a.queries.UpdateJobApplicationStatus(a.dbctx, db.UpdateJobApplicationStatusParams{
		Name: item.Name,
		ID:   item.ID,
	})
}

func (a *App) ListEventSource() ([]db.EventSource, error) {
	return a.queries.ListEventSource(a.dbctx)
}
func (a *App) InsertEventSource(item db.EventSource) (db.EventSource, error) {
	return a.queries.InsertEventSource(a.dbctx, item.Name)
}
func (a *App) DeleteEventSource(item db.EventSource) error {
	return a.queries.DeleteEventSource(a.dbctx, item.Name)
}
func (a *App) UpdateEventSource(item db.EventSource) error {
	return a.queries.UpdateEventSource(a.dbctx, db.UpdateEventSourceParams{
		Name: item.Name,
		ID:   item.ID,
	})
}

type Contact struct {
	ID          int64     `json:"id"`
	CompanyID   int64     `json:"company_id"`
	JobPosition string    `json:"job_position"`
	FistName    string    `json:"fist_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Notes       string    `json:"notes"`
	CompanyName string    `json:"company_name"`
	LastEvent   time.Time `json:"last_event"`
	NextEvent   time.Time `json:"next_event"`
}

func (a *App) ListContact() ([]Contact, error) {
	res, err := a.queries.ListContact(a.dbctx)
	if err != nil {
		return nil, err
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
			LastEvent:   time.Unix(dbContact.LastEvent, 0),
			NextEvent:   time.Unix(dbContact.NextEvent, 0),
		}
	}
	return contacts, err
}
func (a *App) InsertContact(item Contact) (Contact, error) {
	c, err := a.queries.GetCompanyIdByName(a.dbctx, item.CompanyName)
	if err != nil {
		return Contact{}, err
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
		return Contact{}, err
	}
	item.ID = newitem.ID
	return item, nil
}
func (a *App) DeleteContact(item Contact) error {
	return a.queries.DeleteContact(a.dbctx, item.ID)
}
func (a *App) UpdateContact(item Contact) error {
	c, err := a.queries.GetCompanyIdByName(a.dbctx, item.CompanyName)
	if err != nil {
		return err
	}
	return a.queries.UpdateContact(a.dbctx, db.UpdateContactParams{
		CompanyID:   c,
		JobPosition: item.JobPosition,
		FistName:    item.FistName,
		LastName:    item.LastName,
		Email:       item.Email,
		PhoneNumber: item.PhoneNumber,
		Notes:       item.Notes,
		ID:          item.ID,
	})
}

type JobApplication struct {
	ID          int64     `json:"id"`
	CompanyID   int64     `json:"company_id"`
	StatusID    int64     `json:"status_id"`
	JobTitle    string    `json:"job_title"`
	Notes       string    `json:"notes"`
	StatusName  string    `json:"status_name"`
	CompanyName string    `json:"company_name"`
	EventCnt    int64     `json:"event_cnt"`
	LastEvent   time.Time `json:"last_event"`
	NextEvent   time.Time `json:"next_event"`
}

func (a *App) ListJobApplication() ([]JobApplication, error) {
	res, err := a.queries.ListJobApplication(a.dbctx)
	if err != nil {
		return nil, err
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
			LastEvent:   time.Unix(dbApp.LastEvent, 0),
			NextEvent:   time.Unix(dbApp.NextEvent, 0),
		}
	}
	return applications, err
}
func (a *App) InsertJobApplication(item JobApplication) (JobApplication, error) {
	c, err := a.queries.GetCompanyIdByName(a.dbctx, item.CompanyName)
	if err != nil {
		return JobApplication{}, err
	}
	s, err := a.queries.GetJobApplicationStatusIdByName(a.dbctx, item.StatusName)
	if err != nil {
		return JobApplication{}, err
	}
	newitem, err := a.queries.InsertJobApplication(a.dbctx, db.InsertJobApplicationParams{
		CompanyID: c,
		StatusID:  s,
		JobTitle:  item.JobTitle,
		Notes:     item.Notes,
	})
	if err != nil {
		return JobApplication{}, err
	}
	item.ID = newitem.ID
	return item, nil
}
func (a *App) DeleteJobApplication(item JobApplication) error {
	return a.queries.DeleteJobApplication(a.dbctx, item.ID)
}
func (a *App) UpdateJobApplication(item JobApplication) error {
	c, err := a.queries.GetCompanyIdByName(a.dbctx, item.CompanyName)
	if err != nil {
		return err
	}
	s, err := a.queries.GetJobApplicationStatusIdByName(a.dbctx, item.StatusName)
	if err != nil {
		return err
	}
	return a.queries.UpdateJobApplication(a.dbctx, db.UpdateJobApplicationParams{
		CompanyID: c,
		StatusID:  s,
		JobTitle:  item.JobTitle,
		Notes:     item.Notes,
		ID:        item.ID,
	})
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

func joinContactNames(contact db.Contact) string {
	return fmt.Sprintf("%s, %s", contact.FistName, contact.LastName)
}

func splitContactNames(contactNamesJoined string) (firstName string, lastName string) {
	spl := strings.Split(contactNamesJoined, ", ")
	return spl[0], spl[1]
}

func (a *App) ListEvents() ([]Event, error) {
	e, err := a.queries.ListEvent(a.dbctx)
	if err != nil {
		return nil, err
	}
	return groupByMap(
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
				g.Contacts = append(g.Contacts, joinContactNames(s.Contact))
			}
		},
	), nil
}

func setEventContacts(a *App, event Event) error {
	err := a.queries.DeleteEventContact(a.dbctx, event.ID)
	if err != nil {
		return err
	}
	for _, c := range event.Contacts {
		firstName, lastName := splitContactNames(c)
		contactId, err := a.queries.GetContactIdByNames(a.dbctx, db.GetContactIdByNamesParams{
			FistName: firstName,
			LastName: lastName,
			Name:     event.CompanyName,
		})
		if err != nil {
			return err
		}
		_, err = a.queries.InsertEventContact(a.dbctx, db.InsertEventContactParams{
			EventID:   event.ID,
			ContactID: contactId,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) InsertEvent(item Event) (Event, error) {
	s, err := a.queries.GetEventSourceIdByName(a.dbctx, item.Source)
	if err != nil {
		return Event{}, err
	}
	ja, err := a.queries.GetJobApplicationIdByName(a.dbctx, db.GetJobApplicationIdByNameParams{
		JobTitle: item.JobTitle,
		Name:     item.CompanyName,
	})
	if err != nil {
		return Event{}, err
	}
	newitem, err := a.queries.InsertEvent(a.dbctx, db.InsertEventParams{
		SourceID:         s,
		JobApplicationID: ja,
		Title:            item.Title,
		Date:             item.Date.Unix(),
		Notes:            item.Notes,
	})
	if err != nil {
		return Event{}, nil
	}
	item.ID = newitem.ID
	err = setEventContacts(a, item)
	if err != nil {
		return Event{}, err
	}
	return item, nil
}
func (a *App) DeleteEvent(item Event) error {
	return a.queries.DeleteEvent(a.dbctx, item.ID)
}
func (a *App) UpdateEvent(item Event) error {
	s, err := a.queries.GetEventSourceIdByName(a.dbctx, item.Source)
	if err != nil {
		return err
	}
	ja, err := a.queries.GetJobApplicationIdByName(a.dbctx, db.GetJobApplicationIdByNameParams{
		JobTitle: item.JobTitle,
		Name:     item.CompanyName,
	})
	if err != nil {
		return err
	}
	err = a.queries.UpdateEvent(a.dbctx, db.UpdateEventParams{
		SourceID:         s,
		JobApplicationID: ja,
		Title:            item.Title,
		Date:             item.Date.Unix(),
		Notes:            item.Notes,
		ID:               item.ID,
	})
	if err != nil {
		return err
	}
	return setEventContacts(a, item)
}
