import { MalakContact } from "@/client/Api"

export const fullName = (contact: MalakContact): string => {

  if (contact?.first_name === undefined && contact?.last_name === undefined) {
    return contact.email as string
  }

  if (contact?.first_name as string !== "" && contact?.last_name === undefined) {
    return contact?.first_name as string
  }

  return `${contact?.first_name as string} ${contact?.last_name as string}`
}
