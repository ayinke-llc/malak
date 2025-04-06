export interface Contact {
  name: string;
  image: string;
}

export interface Card {
  id: string;
  title: string;
  amount: string;
  stage: string;
  dueDate: string;
  contact: Contact;
  checkSize: string;
  initialContactDate: string;
  isLeadInvestor: boolean;
  rating: number;
}

export interface Note {
  id: string;
  title: string;
  content: string;
  createdAt: string;
  updatedAt: string;
}

export type ActivityType = 'email' | 'meeting' | 'note';

export interface Activity {
  id: string;
  type: ActivityType;
  title: string;
  description: string;
  timestamp: string;
  content?: string;
}

export interface Document {
  id: string;
  name: string;
  type: 'pdf' | 'excel' | 'image' | 'other';
  size: number;
  uploadedAt: Date;
  uploadedBy: string;
}

export const TOTAL_ACTIVITIES_LIMIT = 250;
export const ACTIVITIES_PER_PAGE = 25; 