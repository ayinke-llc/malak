import { MalakContact } from "@/client/Api"

export const fullName = (contact: MalakContact): string => {
  // If both names are undefined or empty, return email
  if ((!contact?.first_name || contact.first_name === "") && (!contact?.last_name || contact.last_name === "")) {
    return contact?.email as string
  }

  // If only first name is present
  if (contact?.first_name && contact.first_name !== "" && (!contact?.last_name || contact.last_name === "")) {
    return contact.first_name
  }

  // If only last name is present
  if ((!contact?.first_name || contact.first_name === "") && contact?.last_name && contact.last_name !== "") {
    return contact.last_name
  }

  // Both names are present
  return `${contact?.first_name} ${contact?.last_name}`.trim()
}
