export interface Employee {
    _id?: string;
    citizen_ucn: string;
    job_id: string;
    start_date: string; // ISO string for time.Time
    end_date: string;
    employer_review: string;
  }