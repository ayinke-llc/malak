import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { RiMailLine, RiPhoneLine } from "@remixicon/react";
import { Contact } from "../../../types";

interface ContactInfoProps {
  contact: Contact;
}

export function ContactInfo({ contact }: ContactInfoProps) {
  return (
    <div className="bg-card rounded-lg p-4 border">
      <h3 className="font-medium mb-3">Contact Information</h3>
      <div className="space-y-3">
        <div className="flex items-center">
          <Avatar className="w-10 h-10">
            <AvatarImage
              src={contact.image}
              alt={contact.name}
            />
            <AvatarFallback>
              {contact.name
                .split(" ")
                .map((n) => n[0])
                .join("")}
            </AvatarFallback>
          </Avatar>
          <div className="ml-3">
            <p className="font-medium">{contact.name}</p>
            <p className="text-sm text-muted-foreground">
              Primary Contact
            </p>
          </div>
        </div>
        <div className="flex items-center text-sm">
          <RiMailLine className="w-4 h-4 mr-2" />
          <span>contact@example.com</span>
        </div>
        <div className="flex items-center text-sm">
          <RiPhoneLine className="w-4 h-4 mr-2" />
          <span>+1 (555) 123-4567</span>
        </div>
      </div>
    </div>
  );
} 
