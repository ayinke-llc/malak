import { useState, useEffect, useCallback } from "react";
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
  RiErrorWarningLine,
} from "@remixicon/react";
import { useQuery } from "@tanstack/react-query";
import { SEARCH_CONTACTS } from "@/lib/query-constants";
import client from "@/lib/client";
import debounce from "lodash/debounce";
import type { ServerListContactsResponse, MalakContact } from "@/client/Api";

export interface SearchResult {
  id: string;
  name: string;
  company: string;
  email: string;
  image: string;
}

const mapContactToSearchResult = (contact: MalakContact): SearchResult => ({
  id: contact.id || "",
  name: `${contact.first_name || ""} ${contact.last_name || ""}`.trim(),
  company: contact.company || "",
  email: contact.email || "",
  image: "", // MalakContact doesn't have an image field
});

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

export function AddInvestorDialog({
  open,
  onOpenChange,
  onAddInvestor
}: AddInvestorDialogProps) {
  const [step, setStep] = useState<'search' | 'details'>('search');
  const [searchQuery, setSearchQuery] = useState("");
  const [debouncedQuery, setDebouncedQuery] = useState("");
  const [selectedInvestor, setSelectedInvestor] = useState<SearchResult | null>(null);
  const [investorDetails, setInvestorDetails] = useState({
    checkSize: "",
    initialContactDate: new Date().toISOString().split('T')[0],
    isLeadInvestor: false,
    rating: 0
  });
  const [hoveredRating, setHoveredRating] = useState<number | null>(null);

  // Create a debounced function to update the search query
  const debouncedSetQuery = useCallback(
    debounce((query: string) => {
      setDebouncedQuery(query);
    }, 300),
    []
  );

  // Update the debounced query when the search input changes
  useEffect(() => {
    if (searchQuery.length >= 4) {
      debouncedSetQuery(searchQuery);
    } else {
      debouncedSetQuery("");
    }
  }, [searchQuery, debouncedSetQuery]);

  const { data: results = [], isLoading, error } = useQuery({
    queryKey: [SEARCH_CONTACTS, debouncedQuery],
    queryFn: async () => {
      if (!debouncedQuery || debouncedQuery.length < 4) return [];
      const response = await client.contacts.searchList({
        search: debouncedQuery
      });
      return (response.data.contacts || []).map(mapContactToSearchResult);
    },
    enabled: debouncedQuery.length >= 4,
  });

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
    setDebouncedQuery("");
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
    setDebouncedQuery("");
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
                    onChange={(e) => setSearchQuery(e.target.value)}
                  />
                </div>
                {searchQuery && searchQuery.length < 4 && (
                  <p className="text-sm text-muted-foreground mt-2">
                    Please enter at least 4 characters to search
                  </p>
                )}
              </div>

              <div className="space-y-2">
                {isLoading && (
                  <div className="text-center py-4 text-sm text-muted-foreground">
                    Searching...
                  </div>
                )}

                {error && (
                  <div className="flex flex-col items-center justify-center py-4">
                    <RiErrorWarningLine className="w-8 h-8 text-destructive mb-2" />
                    <p className="text-sm text-destructive text-center">
                      An error occurred while searching. Please try again.
                    </p>
                  </div>
                )}

                {!isLoading && !error && searchQuery.length >= 4 && results.length === 0 && (
                  <Card className="py-8">
                    <div className="flex flex-col items-center justify-center text-center px-4">
                      <RiContactsLine className="h-8 w-8 text-muted-foreground/50 mb-4" />
                      <h3 className="text-lg font-medium mb-2">No contacts found</h3>
                      <p className="text-sm text-muted-foreground mb-4">
                        We couldn't find any contacts matching your search. Would you like to add a new contact?
                      </p>
                      <Link 
                        href="/contacts"
                        className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2"
                      >
                        <RiUserAddLine className="mr-2 h-4 w-4" />
                        Add New Contact
                      </Link>
                    </div>
                  </Card>
                )}

                {!isLoading && results.map((result: SearchResult) => (
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
                            {result.name.split(" ").map((n: string) => n[0]).join("")}
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
