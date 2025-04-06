import { Badge } from "@/components/ui/badge";
import { RiStarLine } from "@remixicon/react";
import { Card } from "../../../types";

interface DealInfoProps {
  investor: Card;
}

export function DealInfo({ investor }: DealInfoProps) {
  return (
    <div className="bg-card rounded-lg p-4 border">
      <h3 className="font-medium mb-3">Deal Information</h3>
      <div className="space-y-3 text-sm">
        <div className="flex items-center justify-between">
          <span className="text-muted-foreground">Deal Size</span>
          <span className="font-medium">${investor.amount}</span>
        </div>
        <div className="flex items-center justify-between">
          <span className="text-muted-foreground">Check Size</span>
          <span className="font-medium">{investor.checkSize}</span>
        </div>
        <div className="flex items-center justify-between">
          <span className="text-muted-foreground">Initial Contact</span>
          <span className="font-medium">{investor.initialContactDate}</span>
        </div>
        <div className="flex items-center justify-between">
          <span className="text-muted-foreground">Lead Investor</span>
          <Badge variant={investor.isLeadInvestor ? "default" : "outline"}>
            {investor.isLeadInvestor ? "Yes" : "No"}
          </Badge>
        </div>
        <div className="flex items-center justify-between">
          <span className="text-muted-foreground">Rating</span>
          <div className="flex items-center">
            {[1, 2, 3, 4, 5].map((star) => (
              <RiStarLine
                key={star}
                className={`w-4 h-4 ${star <= (investor.rating || 0) ? 'text-yellow-400' : 'text-muted-foreground'}`}
              />
            ))}
          </div>
        </div>
        <div className="flex items-center justify-between">
          <span className="text-muted-foreground">Stage</span>
          <Badge variant="outline">{investor.stage}</Badge>
        </div>
        <div className="flex items-center justify-between">
          <span className="text-muted-foreground">Due Date</span>
          <span className="font-medium">{investor.dueDate}</span>
        </div>
      </div>
    </div>
  );
} 
