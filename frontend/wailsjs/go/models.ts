export namespace db {
	
	export class ListCompanyTypeRow {
	    id: number;
	    name: string;
	    companies: number;
	
	    static createFrom(source: any = {}) {
	        return new ListCompanyTypeRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.companies = source["companies"];
	    }
	}
	export class ListEventSourceRow {
	    id: number;
	    name: string;
	    events: number;
	
	    static createFrom(source: any = {}) {
	        return new ListEventSourceRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.events = source["events"];
	    }
	}
	export class ListJobApplicationStatusRow {
	    id: number;
	    name: string;
	    applications: number;
	
	    static createFrom(source: any = {}) {
	        return new ListJobApplicationStatusRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.applications = source["applications"];
	    }
	}

}

export namespace main {
	
	export class Company {
	    id: number;
	    name: string;
	    notes: string;
	    company_types: string[];
	    application_cnt: number;
	    // Go type: time
	    last_event?: any;
	    // Go type: time
	    next_event?: any;
	
	    static createFrom(source: any = {}) {
	        return new Company(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.notes = source["notes"];
	        this.company_types = source["company_types"];
	        this.application_cnt = source["application_cnt"];
	        this.last_event = this.convertValues(source["last_event"], null);
	        this.next_event = this.convertValues(source["next_event"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Contact {
	    id: number;
	    company_id: number;
	    job_position: string;
	    fist_name: string;
	    last_name: string;
	    email: string;
	    phone_number: string;
	    notes: string;
	    company_name: string;
	    // Go type: time
	    last_event?: any;
	    // Go type: time
	    next_event?: any;
	
	    static createFrom(source: any = {}) {
	        return new Contact(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.company_id = source["company_id"];
	        this.job_position = source["job_position"];
	        this.fist_name = source["fist_name"];
	        this.last_name = source["last_name"];
	        this.email = source["email"];
	        this.phone_number = source["phone_number"];
	        this.notes = source["notes"];
	        this.company_name = source["company_name"];
	        this.last_event = this.convertValues(source["last_event"], null);
	        this.next_event = this.convertValues(source["next_event"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Event {
	    id: number;
	    title: string;
	    // Go type: time
	    date: any;
	    notes: string;
	    source: string;
	    job_title: string;
	    company_name: string;
	    contacts: string[];
	
	    static createFrom(source: any = {}) {
	        return new Event(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.date = this.convertValues(source["date"], null);
	        this.notes = source["notes"];
	        this.source = source["source"];
	        this.job_title = source["job_title"];
	        this.company_name = source["company_name"];
	        this.contacts = source["contacts"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class JobApplication {
	    id: number;
	    company_id: number;
	    status_id: number;
	    job_title: string;
	    notes: string;
	    status_name: string;
	    company_name: string;
	    event_cnt: number;
	    // Go type: time
	    last_event?: any;
	    // Go type: time
	    next_event?: any;
	
	    static createFrom(source: any = {}) {
	        return new JobApplication(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.company_id = source["company_id"];
	        this.status_id = source["status_id"];
	        this.job_title = source["job_title"];
	        this.notes = source["notes"];
	        this.status_name = source["status_name"];
	        this.company_name = source["company_name"];
	        this.event_cnt = source["event_cnt"];
	        this.last_event = this.convertValues(source["last_event"], null);
	        this.next_event = this.convertValues(source["next_event"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

