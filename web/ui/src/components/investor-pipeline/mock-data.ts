import type { Columns, Board } from "./types";

export const initialColumns: Columns = {
  "backlog": {
    id: "backlog",
    title: "Backlog",
    cards: [
      {
        id: "1",
        title: "Sequoia Capital",
        amount: "5M",
        stage: "Initial Research",
        dueDate: "2024-04-20",
        contact: {
          name: "Sarah Chen",
          image: "/avatars/sarah.jpg",
        },
        roundDetails: {
          raising: "20M",
          type: "Series A",
          ownership: "15-20%"
        },
        checkSize: "TBD",
        initialContactDate: new Date().toISOString().split('T')[0],
        isLeadInvestor: false,
        rating: 0
      },
      {
        id: "2",
        title: "Andreessen Horowitz",
        amount: "10M",
        stage: "Initial Contact",
        dueDate: "2024-04-25",
        contact: {
          name: "Marc A.",
          image: "/avatars/marc.jpg",
        },
        roundDetails: {
          raising: "50M",
          type: "Series B",
          ownership: "10-15%"
        },
        checkSize: "TBD",
        initialContactDate: new Date().toISOString().split('T')[0],
        isLeadInvestor: false,
        rating: 0
      },
      {
        id: "4",
        title: "Greylock Partners",
        amount: "8M",
        stage: "Initial Research",
        dueDate: "2024-05-01",
        contact: {
          name: "Reid Hoffman",
          image: "/avatars/reid.jpg",
        },
        roundDetails: {
          raising: "40M",
          type: "Series C",
          ownership: "12-18%"
        },
        checkSize: "TBD",
        initialContactDate: new Date().toISOString().split('T')[0],
        isLeadInvestor: false,
        rating: 0
      },
      {
        id: "5",
        title: "Benchmark",
        amount: "6M",
        stage: "Initial Contact",
        dueDate: "2024-05-05",
        contact: {
          name: "Bill Gurley",
          image: "/avatars/bill.jpg",
        },
        roundDetails: {
          raising: "25M",
          type: "Series A",
          ownership: "15-20%"
        },
        checkSize: "TBD",
        initialContactDate: new Date().toISOString().split('T')[0],
        isLeadInvestor: false,
        rating: 0
      },
    ],
  },
  "research": {
    id: "research",
    title: "Research",
    cards: [
      {
        id: "3",
        title: "Lightspeed Ventures",
        amount: "7M",
        stage: "Due Diligence",
        dueDate: "2024-04-28",
        contact: {
          name: "Amy Wu",
          image: "/avatars/amy.jpg",
        },
        roundDetails: {
          raising: "30M",
          type: "Series A",
          ownership: "18-22%"
        },
        checkSize: "TBD",
        initialContactDate: new Date().toISOString().split('T')[0],
        isLeadInvestor: false,
        rating: 0
      },
      {
        id: "6",
        title: "Accel Partners",
        amount: "9M",
        stage: "Due Diligence",
        dueDate: "2024-05-10",
        contact: {
          name: "Jim Breyer",
          image: "/avatars/jim.jpg",
        },
        roundDetails: {
          raising: "35M",
          type: "Series B",
          ownership: "10-15%"
        },
        checkSize: "TBD",
        initialContactDate: new Date().toISOString().split('T')[0],
        isLeadInvestor: false,
        rating: 0
      },
    ],
  },
  "partner-call": {
    id: "partner-call",
    title: "Partner Call",
    cards: [
      {
        id: "7",
        title: "Kleiner Perkins",
        amount: "11M",
        stage: "Partner Call",
        dueDate: "2024-05-15",
        contact: {
          name: "John Doerr",
          image: "/avatars/john.jpg",
        },
        roundDetails: {
          raising: "45M",
          type: "Series C",
          ownership: "20-25%"
        },
        checkSize: "TBD",
        initialContactDate: new Date().toISOString().split('T')[0],
        isLeadInvestor: false,
        rating: 0
      },
    ],
  },
  "passed": {
    id: "passed",
    title: "Passed",
    cards: [
      {
        id: "8",
        title: "Union Square Ventures",
        amount: "4M",
        stage: "Passed",
        dueDate: "2024-05-20",
        contact: {
          name: "Fred Wilson",
          image: "/avatars/fred.jpg",
        },
        roundDetails: {
          raising: "20M",
          type: "Series A",
          ownership: "15-20%"
        },
        checkSize: "TBD",
        initialContactDate: new Date().toISOString().split('T')[0],
        isLeadInvestor: false,
        rating: 0
      },
    ],
  },
  "termsheet": {
    id: "termsheet",
    title: "Termsheet",
    cards: [
      {
        id: "9",
        title: "Bessemer Venture Partners",
        amount: "12M",
        stage: "Termsheet",
        dueDate: "2024-05-25",
        contact: {
          name: "David Cowan",
          image: "/avatars/david.jpg",
        },
        roundDetails: {
          raising: "50M",
          type: "Series D",
          ownership: "10-15%"
        },
        checkSize: "TBD",
        initialContactDate: new Date().toISOString().split('T')[0],
        isLeadInvestor: false,
        rating: 0
      },
    ],
  },
  "closed": {
    id: "closed",
    title: "Closed",
    cards: [
      {
        id: "10",
        title: "General Catalyst",
        amount: "15M",
        stage: "Closed",
        dueDate: "2024-05-30",
        contact: {
          name: "Joel Cutler",
          image: "/avatars/joel.jpg",
        },
        roundDetails: {
          raising: "60M",
          type: "Series E",
          ownership: "5-10%"
        },
        checkSize: "TBD",
        initialContactDate: new Date().toISOString().split('T')[0],
        isLeadInvestor: false,
        rating: 0
      },
    ],
  },
};

export const initialBoard: Board = {
  columns: initialColumns,
  isArchived: false
}; 