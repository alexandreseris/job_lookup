// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {main} from '../models';
import {db} from '../models';

export function DeleteCompany(arg1:main.Company):Promise<void>;

export function DeleteCompanyType(arg1:db.ListCompanyTypeRow):Promise<void>;

export function DeleteContact(arg1:main.Contact):Promise<void>;

export function DeleteEvent(arg1:main.Event):Promise<void>;

export function DeleteEventSource(arg1:db.ListEventSourceRow):Promise<void>;

export function DeleteJobApplication(arg1:main.JobApplication):Promise<void>;

export function DeleteJobApplicationStatus(arg1:db.ListJobApplicationStatusRow):Promise<void>;

export function InsertCompany(arg1:main.Company):Promise<main.Company>;

export function InsertCompanyType(arg1:db.ListCompanyTypeRow):Promise<db.ListCompanyTypeRow>;

export function InsertContact(arg1:main.Contact):Promise<main.Contact>;

export function InsertEvent(arg1:main.Event):Promise<main.Event>;

export function InsertEventSource(arg1:db.ListEventSourceRow):Promise<db.ListEventSourceRow>;

export function InsertJobApplication(arg1:main.JobApplication):Promise<main.JobApplication>;

export function InsertJobApplicationStatus(arg1:db.ListJobApplicationStatusRow):Promise<db.ListJobApplicationStatusRow>;

export function ListCompanies():Promise<Array<main.Company>>;

export function ListCompanyTypes():Promise<Array<db.ListCompanyTypeRow>>;

export function ListContact():Promise<Array<main.Contact>>;

export function ListEventSource():Promise<Array<db.ListEventSourceRow>>;

export function ListEvents():Promise<Array<main.Event>>;

export function ListJobApplication():Promise<Array<main.JobApplication>>;

export function ListJobApplicationStatus():Promise<Array<db.ListJobApplicationStatusRow>>;

export function Log(arg1:string):Promise<void>;

export function UpdateCompany(arg1:main.Company):Promise<void>;

export function UpdateCompanyType(arg1:db.ListCompanyTypeRow):Promise<void>;

export function UpdateContact(arg1:main.Contact):Promise<void>;

export function UpdateEvent(arg1:main.Event):Promise<void>;

export function UpdateEventSource(arg1:db.ListEventSourceRow):Promise<void>;

export function UpdateJobApplication(arg1:main.JobApplication):Promise<void>;

export function UpdateJobApplicationStatus(arg1:db.ListJobApplicationStatusRow):Promise<void>;
