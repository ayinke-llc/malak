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
import { CurrencyInput } from "@/components/ui/currency-input";

export interface SearchResult {
  reference: string;
  name: string;
  company: string;
  email: string;
  image: string;
  isExisting?: boolean;
}

const mapContactToSearchResult = (contact: MalakContact): SearchResult => ({
  reference: contact.reference || "",
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
  isLoading?: boolean;
  existingContacts?: string[]; // Array of contact IDs that are already in the board
}

export function AddInvestorDialog({
  open,
  onOpenChange,
  onAddInvestor,
  isLoading = false,
  existingContacts = []
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

  const { data: results = [], isLoading: queryLoading, error } = useQuery<(SearchResult & { isExisting: boolean })[]>({
    queryKey: [SEARCH_CONTACTS, debouncedQuery],
    queryFn: async () => {
      if (!debouncedQuery || debouncedQuery.length < 4) return [];
      const response = await client.contacts.searchList({
        search: debouncedQuery
      });
      
      // Create a Set for O(1) lookup of existing contacts
      const existingContactsSet = new Set(existingContacts);
      
      // Map all contacts and include whether they're already in the board
      return (response.data.contacts || []).map(contact => ({
        ...mapContactToSearchResult(contact),
        isExisting: existingContactsSet.has(contact.id || "")
      }));
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

    // Don't reset state or close dialog here - let the parent component handle it after successful API call
  };

  const handleClose = () => {
    if (isLoading) return; // Prevent closing while loading
    
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
                {queryLoading && (
                  <div className="text-center py-4 text-sm text-muted-foreground">
                    Searching...
                  </div>
                )}

                {error && (
                  <div className="text-center py-4">
                    <RiErrorWarningLine className="w-6 h-6 text-destructive mx-auto mb-2" />
                    <p className="text-sm text-destructive">Failed to search contacts</p>
                  </div>
                )}

                {!queryLoading && !error && results.length === 0 && debouncedQuery && (
                  <div className="text-center py-4">
                    <RiContactsLine className="w-6 h-6 text-muted-foreground mx-auto mb-2" />
                    <p className="text-sm text-muted-foreground">No contacts found</p>
                  </div>
                )}

                {!queryLoading && !error && results.map((result) => (
                  <Card
                    key={result.reference}
                    className={`cursor-pointer transition-colors ${
                      result.isExisting ? 'opacity-50 pointer-events-none' : 'hover:bg-accent/5'
                    }`}
                    onClick={() => !result.isExisting && handleSelectInvestor(result)}
                  >
                    <CardContent className="p-3">
                      <div className="flex items-center gap-3">
                        <Avatar className="h-8 w-8">
                          <AvatarImage src={result.image} alt={result.name} />
                          <AvatarFallback>
                            {result.name.split(" ").map((n) => n[0]).join("")}
                          </AvatarFallback>
                        </Avatar>
                        <div className="min-w-0 flex-1">
                          <h4 className="truncate font-medium text-sm">
                            {result.name}
                            {result.isExisting && (
                              <span className="ml-2 text-xs text-muted-foreground">(Already in board)</span>
                            )}
                          </h4>
                          <p className="truncate text-xs text-muted-foreground">
                            {result.company}
                          </p>
                          <p className="truncate text-xs text-muted-foreground">
                            {result.email}
                          </p>
                        </div>
                        {!result.isExisting && (
                          <Button
                            variant="ghost"
                            size="icon"
                            className="ml-2"
                            onClick={() => handleSelectInvestor(result)}
                          >
                            <RiArrowRightLine className="h-4 w-4" />
                          </Button>
                        )}
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
                <CurrencyInput
                  value={investorDetails.checkSize}
                  onValueChange={(values) => {
                    setInvestorDetails({
                      ...investorDetails,
                      checkSize: values.value
                    });
                  }}
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
                  disabled={isLoading}
                >
                  Back to Search
                </Button>
                <Button type="submit" disabled={isLoading}>
                  {isLoading ? (
                    <>
                      <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-background mr-2" />
                      Adding...
                    </>
                  ) : (
                    "Add to Pipeline"
                  )}
                </Button>
              </div>
            </form>
          </>
        )}
      </DialogContent>
    </Dialog>
  );
} 
