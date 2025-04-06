import { useState } from "react";
import Link from "next/link";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import {
  RiSearchLine,
  RiStarFill,
  RiStarLine,
  RiArrowRightLine,
  RiUserAddLine,
  RiContactsLine,
} from "@remixicon/react";

export interface SearchResult {
  id: string;
  name: string;
  company: string;
  email: string;
  image: string;
}

interface AddInvestorDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onAddInvestor: (investor: SearchResult & {
    checkSize: string;
    initialContactDate: string;
    isLeadInvestor: boolean;
    rating: number;
  }) => void;
}

// Mock search results
const mockSearchResults: SearchResult[] = [
  {
    id: "sr1",
    name: "David Marcus",
    company: "Breakthrough Capital",
    email: "david@breakthrough.vc",
    image: "/avatars/david.jpg"
  },
  {
    id: "sr2",
    name: "Lisa Wong",
    company: "Horizon Ventures",
    email: "lisa@horizon.vc",
    image: "/avatars/lisa.jpg"
  },
  {
    id: "sr3",
    name: "Michael Chen",
    company: "Dragon Capital",
    email: "michael@dragon.vc",
    image: "/avatars/michael.jpg"
  }
];

// Mock search function
const searchInvestors = async (query: string): Promise<SearchResult[]> => {
  // Simulate API delay
  await new Promise(resolve => setTimeout(resolve, 500));

  // Filter mock results based on query
  return mockSearchResults.filter(result =>
    result.name.toLowerCase().includes(query.toLowerCase()) ||
    result.company.toLowerCase().includes(query.toLowerCase()) ||
    result.email.toLowerCase().includes(query.toLowerCase())
  );
};

export function AddInvestorDialog({
  open,
  onOpenChange,
  onAddInvestor
}: AddInvestorDialogProps) {
  const [step, setStep] = useState<'search' | 'details'>('search');
  const [searchQuery, setSearchQuery] = useState("");
  const [results, setResults] = useState<SearchResult[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedInvestor, setSelectedInvestor] = useState<SearchResult | null>(null);
  const [investorDetails, setInvestorDetails] = useState({
    checkSize: "",
    initialContactDate: new Date().toISOString().split('T')[0],
    isLeadInvestor: false,
    rating: 0
  });
  const [hoveredRating, setHoveredRating] = useState<number | null>(null);

  const handleSearch = async (query: string) => {
    setSearchQuery(query);
    if (query.length < 2) {
      setResults([]);
      return;
    }

    setLoading(true);
    try {
      const searchResults = await searchInvestors(query);
      setResults(searchResults);
    } catch (error) {
      console.error("Error searching investors:", error);
      setResults([]);
    } finally {
      setLoading(false);
    }
  };

  const handleSelectInvestor = (investor: SearchResult) => {
    setSelectedInvestor(investor);
    setStep('details');
  };

  const handleSubmitDetails = (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedInvestor) return;

    onAddInvestor({
      ...selectedInvestor,
      ...investorDetails
    });

    // Reset state
    setStep('search');
    setSearchQuery("");
    setResults([]);
    setSelectedInvestor(null);
    setInvestorDetails({
      checkSize: "",
      initialContactDate: new Date().toISOString().split('T')[0],
      isLeadInvestor: false,
      rating: 0
    });
    setHoveredRating(null);
    onOpenChange(false);
  };

  const handleClose = () => {
    // Reset state
    setStep('search');
    setSearchQuery("");
    setResults([]);
    setSelectedInvestor(null);
    setInvestorDetails({
      checkSize: "",
      initialContactDate: new Date().toISOString().split('T')[0],
      isLeadInvestor: false,
      rating: 0
    });
    setHoveredRating(null);
    onOpenChange(false);
  };

  const displayRating = hoveredRating !== null ? hoveredRating : investorDetails.rating;

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-lg">
        {step === 'search' ? (
          <>
            <DialogHeader>
              <DialogTitle>Add Investor</DialogTitle>
              <DialogDescription>
                Search for an investor by name, company, or email to add them to your pipeline.
              </DialogDescription>
            </DialogHeader>

            <div className="mt-4 space-y-4">
              <div>
                <Label htmlFor="search">Search</Label>
                <div className="relative mt-2">
                  <RiSearchLine className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground h-4 w-4" />
                  <Input
                    id="search"
                    placeholder="Search investors..."
                    className="pl-9"
                    value={searchQuery}
                    onChange={(e) => handleSearch(e.target.value)}
                  />
                </div>
              </div>

              <div className="space-y-2">
                {loading && (
                  <div className="text-center py-4 text-sm text-muted-foreground">
                    Searching...
                  </div>
                )}

                {!loading && results.length === 0 && searchQuery.length >= 2 && (
                  <div className="text-center py-6 space-y-4">
                    <div className="flex justify-center">
                      <div className="bg-muted p-4 rounded-full">
                        <RiUserAddLine className="h-8 w-8 text-muted-foreground" />
                      </div>
                    </div>
                    <div className="space-y-2">
                      <h3 className="font-medium">Contact Not Found</h3>
                      <p className="text-sm text-muted-foreground max-w-sm mx-auto">
                        The contact you're looking for doesn't exist yet. You'll need to create it first in your contacts.
                      </p>
                    </div>
                    <Link
                      href="/contacts"
                      className="inline-flex items-center gap-2 px-4 py-2 rounded-md bg-primary text-primary-foreground hover:bg-primary/90 transition-colors"
                    >
                      <RiContactsLine className="h-4 w-4" />
                      Create New Contact
                    </Link>
                  </div>
                )}

                {!loading && results.map((result) => (
                  <Card
                    key={result.id}
                    className="cursor-pointer transition-colors hover:bg-muted/50"
                    onClick={() => handleSelectInvestor(result)}
                  >
                    <CardContent className="p-3">
                      <div className="flex items-center gap-3">
                        <Avatar className="h-8 w-8">
                          <AvatarImage src={result.image} alt={result.name} />
                          <AvatarFallback>
                            {result.name.split(" ").map(n => n[0]).join("")}
                          </AvatarFallback>
                        </Avatar>
                        <div className="min-w-0 flex-1">
                          <h4 className="truncate font-medium text-sm">{result.name}</h4>
                          <p className="truncate text-xs text-muted-foreground">{result.company}</p>
                          <p className="truncate text-xs text-muted-foreground/75">{result.email}</p>
                        </div>
                        <RiArrowRightLine className="h-4 w-4 text-muted-foreground" />
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            </div>
          </>
        ) : (
          <>
            <DialogHeader>
              <DialogTitle>Investor Details</DialogTitle>
              <DialogDescription>
                Add additional details about {selectedInvestor?.name} from {selectedInvestor?.company}.
              </DialogDescription>
            </DialogHeader>

            <form onSubmit={handleSubmitDetails} className="mt-4 space-y-4">
              <div className="space-y-2">
                <Label>Check Size</Label>
                <Input
                  value={investorDetails.checkSize}
                  onChange={(e) => setInvestorDetails({ ...investorDetails, checkSize: e.target.value })}
                  placeholder="Enter check size"
                />
              </div>

              <div className="space-y-2">
                <Label>Initial Contact Date</Label>
                <Input
                  type="date"
                  value={investorDetails.initialContactDate}
                  onChange={(e) => setInvestorDetails({ ...investorDetails, initialContactDate: e.target.value })}
                />
              </div>

              <div className="space-y-2">
                <Label>Lead Investor</Label>
                <div className="flex items-center space-x-2">
                  <Switch
                    checked={investorDetails.isLeadInvestor}
                    onCheckedChange={(checked) => setInvestorDetails({ ...investorDetails, isLeadInvestor: checked })}
                  />
                  <span className="text-sm text-muted-foreground">
                    {investorDetails.isLeadInvestor ? 'Yes' : 'No'}
                  </span>
                </div>
              </div>

              <div className="space-y-2">
                <Label>Rating</Label>
                <div className="flex items-center gap-1">
                  {[1, 2, 3, 4, 5].map((star) => (
                    <Button
                      key={star}
                      type="button"
                      variant="ghost"
                      size="icon"
                      className="hover:bg-transparent"
                      onClick={() => setInvestorDetails({ ...investorDetails, rating: star })}
                      onMouseEnter={() => setHoveredRating(star)}
                      onMouseLeave={() => setHoveredRating(null)}
                    >
                      {star <= displayRating ? (
                        <RiStarFill className="w-6 h-6 text-yellow-400" />
                      ) : (
                        <RiStarLine className="w-6 h-6 text-muted-foreground" />
                      )}
                    </Button>
                  ))}
                  <span className="ml-2 text-sm text-muted-foreground">
                    {displayRating} of 5
                  </span>
                </div>
              </div>

              <div className="flex justify-end gap-2">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setStep('search')}
                >
                  Back to Search
                </Button>
                <Button type="submit">Add to Pipeline</Button>
              </div>
            </form>
          </>
        )}
      </DialogContent>
    </Dialog>
  );
} 
