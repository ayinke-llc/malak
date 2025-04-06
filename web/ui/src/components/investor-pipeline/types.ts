import type { MalakContact, MalakFundraiseContactDealDetails } from "@/client/Api";

export interface Contact {
  name: string;
  image?: string;
  company?: string;
}

export interface Card {
  id: string;
  title: string;
  amount: string;
  stage: string;
  dueDate: string;
  contact: {
    name: string;
    company?: string;
    email?: string;
    phone?: string;
    title?: string;
  };
  roundDetails: {
    raising: string;
    type: string;
    ownership: string;
  };
  checkSize: string;
  initialContactDate: string;
  isLeadInvestor: boolean;
  rating: number;
  originalContact?: MalakContact;
  originalDeal?: MalakFundraiseContactDealDetails;
}

export interface Column {
  id: string;
  title: string;
  description: string;
  cards: Card[];
}

export interface Columns {
  [key: string]: Column;
}

export interface Board {
  columns: Columns;
  isArchived: boolean;
}

export interface ShareSettings {
  isEnabled: boolean;
  shareLink: string;
  requireEmail: boolean;
  requirePassword: boolean;
  password?: string;
}

export interface Note {
  id: string;
  title: string;
  content: string;
  createdAt: string;
  updatedAt: string;
}

export interface Activity {
  id: string;
  type: 'email' | 'meeting' | 'document' | 'team' | 'stage_change';
  title: string;
  description: string;
  timestamp: string;
  content?: string;
  metadata?: {
    fromStage?: string;
    toStage?: string;
  };
}

export interface Document {
  id: string;
  name: string;
  type: 'pdf' | 'excel' | 'image' | 'other';
  size: number;
  uploadedAt: Date;
  uploadedBy: string;
} 