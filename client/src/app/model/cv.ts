import { Diploma } from './diploma';

export interface CV {
  _id?: string;
  citizen_ucn: string;
  name: string;
  email: string;
  description: string;
  work_experience: string[];
  education: Diploma[];
}