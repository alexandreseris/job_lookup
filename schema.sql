CREATE TABLE company (
    id integer PRIMARY KEY AUTOINCREMENT,
    name text NOT NULL,
    notes text NOT NULL
);

CREATE UNIQUE INDEX idx_company_name ON company(name);

CREATE TABLE company_type (
    id integer PRIMARY KEY AUTOINCREMENT,
    name text NOT NULL
);

CREATE UNIQUE INDEX idx_company_type_name ON company_type(name);

CREATE TABLE company_type_rel (
    id integer PRIMARY KEY AUTOINCREMENT,
    company_id integer NOT NULL,
    company_type_id integer NOT NULL,
    CONSTRAINT fk_company_type_rel_company FOREIGN KEY (company_id) REFERENCES company(id) ON DELETE CASCADE,
    CONSTRAINT fk_company_type_rel_company_type FOREIGN KEY (company_type_id) REFERENCES company_type(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_company_type_rel_company_company_type ON company_type_rel(company_id, company_type_id);

CREATE INDEX idx_company_type_rel_company_type_id ON company_type_rel(company_type_id);

CREATE INDEX idx_company_type_rel_company_id ON company_type_rel(company_id);

CREATE TABLE contact (
    id integer PRIMARY KEY AUTOINCREMENT,
    company_id integer NOT NULL,
    job_position text NOT NULL,
    fist_name text NOT NULL,
    last_name text NOT NULL,
    email text NOT NULL,
    phone_number text NOT NULL,
    notes text NOT NULL,
    CONSTRAINT fk_contact_company FOREIGN KEY (company_id) REFERENCES company(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_contact_company_contact_name ON contact(company_id, fist_name, last_name);

CREATE INDEX idx_contact_company_id ON contact(company_id);

CREATE TABLE job_application_status (
    id integer PRIMARY KEY AUTOINCREMENT,
    name text NOT NULL
);

CREATE UNIQUE INDEX idx_job_application_status_name ON job_application_status(name);

CREATE TABLE job_application (
    id integer PRIMARY KEY AUTOINCREMENT,
    company_id integer NOT NULL,
    status_id integer NOT NULL,
    job_title text NOT NULL,
    notes text NOT NULL,
    CONSTRAINT fk_job_application_company FOREIGN KEY (company_id) REFERENCES company(id) ON DELETE CASCADE,
    CONSTRAINT fk_job_application_status FOREIGN KEY (status_id) REFERENCES job_application_status(id) ON DELETE CASCADE
);

CREATE INDEX idx_job_application_status_id ON job_application(status_id);

CREATE UNIQUE INDEX job_application_uniq ON job_application(company_id, job_title);

CREATE TABLE event_source (
    id integer PRIMARY KEY AUTOINCREMENT,
    name text NOT NULL
);

CREATE UNIQUE INDEX idx_event_source_name ON event_source(name);

CREATE TABLE event (
    id integer PRIMARY KEY AUTOINCREMENT,
    source_id integer NOT NULL,
    job_application_id integer NOT NULL,
    title text NOT NULL,
    date integer NOT NULL,
    notes text NOT NULL,
    CONSTRAINT fk_event_job_application FOREIGN KEY (job_application_id) REFERENCES job_application(id) ON DELETE CASCADE,
    CONSTRAINT fk_event_source FOREIGN KEY (source_id) REFERENCES event_source(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_event_application_title ON event(job_application_id, title);

CREATE INDEX idx_event_job_application_id ON event(job_application_id);

CREATE INDEX idx_event_source_id ON event(source_id);

CREATE TABLE event_contacts (
    id integer PRIMARY KEY AUTOINCREMENT,
    event_id integer NOT NULL,
    contact_id integer NOT NULL,
    CONSTRAINT fk_event_contacts_contact FOREIGN KEY (contact_id) REFERENCES contact(id) ON DELETE CASCADE,
    CONSTRAINT fk_event_contacts_event FOREIGN KEY (event_id) REFERENCES event(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_event_contacts_event_contact ON event_contacts(event_id, contact_id);

CREATE INDEX idx_event_contacts_event_id ON event_contacts(event_id);

CREATE INDEX idx_event_contacts_contact_id ON event_contacts(contact_id);