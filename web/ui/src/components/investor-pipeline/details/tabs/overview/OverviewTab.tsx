import { Card } from "../../../types";
import { ContactInfo } from "./ContactInfo";
import { DealInfo } from "./DealInfo";
import { SuggestedActions } from "./SuggestedActions";

interface OverviewTabProps {
  investor: Card;
  onAddActivity: () => void;
  onUploadDocument: () => void;
}

export function OverviewTab({ investor, onAddActivity, onUploadDocument }: OverviewTabProps) {
  return (
    <div className="space-y-6">
      <ContactInfo contact={investor.contact} />
      <DealInfo investor={investor} />
      <SuggestedActions
        onAddActivity={onAddActivity}
        onUploadDocument={onUploadDocument}
      />
    </div>
  );
} 
