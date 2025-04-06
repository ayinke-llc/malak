import { useState, useEffect } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";
import { RiStarFill, RiStarLine } from "@remixicon/react";
import { NumericFormat } from "react-number-format";
import { Card } from "../../../types";

interface EditInvestorDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  investor: Card | null;
  onSave: (updatedInvestor: Card) => void;
}

export function EditInvestorDialog({ 
  open, 
  onOpenChange,
  investor,
  onSave
}: EditInvestorDialogProps) {
  const [editedInvestor, setEditedInvestor] = useState<Card | null>(null);
  const [hoveredRating, setHoveredRating] = useState<number | null>(null);

  useEffect(() => {
    if (investor) {
      setEditedInvestor({ ...investor });
    }
  }, [investor]);

  if (!editedInvestor) return null;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (editedInvestor) {
      onSave(editedInvestor);
      onOpenChange(false);
    }
  };

  const displayRating = hoveredRating !== null ? hoveredRating : editedInvestor.rating;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Investor Details</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label>Amount</Label>
            <NumericFormat
              value={editedInvestor.amount.replace(/[^0-9.]/g, '')}
              onValueChange={(values) => {
                const { value } = values;
                const formattedValue = `${value}M`;
                setEditedInvestor({ ...editedInvestor, amount: formattedValue });
              }}
              thousandSeparator
              prefix="$"
              customInput={Input}
              placeholder="Enter amount (e.g. 5M)"
              className="[appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
            />
          </div>

          <div className="space-y-2">
            <Label>Check Size</Label>
            <NumericFormat
              value={editedInvestor.checkSize.replace(/[^0-9.]/g, '')}
              onValueChange={(values) => {
                const { value } = values;
                const formattedValue = `${value}M`;
                setEditedInvestor({ ...editedInvestor, checkSize: formattedValue });
              }}
              thousandSeparator
              prefix="$"
              customInput={Input}
              placeholder="Enter check size"
              className="[appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
            />
          </div>
          
          <div className="space-y-2">
            <Label>Initial Contact Date</Label>
            <Input
              type="date"
              value={editedInvestor.initialContactDate}
              onChange={(e) => setEditedInvestor({ ...editedInvestor, initialContactDate: e.target.value })}
            />
          </div>
          
          <div className="space-y-2">
            <Label>Lead Investor</Label>
            <div className="flex items-center space-x-2">
              <Switch
                checked={editedInvestor.isLeadInvestor}
                onCheckedChange={(checked) => setEditedInvestor({ ...editedInvestor, isLeadInvestor: checked })}
              />
              <span className="text-sm text-muted-foreground">
                {editedInvestor.isLeadInvestor ? 'Yes' : 'No'}
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
                  onClick={() => setEditedInvestor({ ...editedInvestor, rating: star })}
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
            <Button variant="outline" type="button" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit">Save Changes</Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
} 