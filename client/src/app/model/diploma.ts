export interface Diploma {
    _id?: string;
    institution_id: string; // UUID as string
    institution_name: string;
    institution_type: string;
    average_grade: number;
    ucn: string;
  }