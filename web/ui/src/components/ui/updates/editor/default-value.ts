import { type PartialBlock } from "@blocknote/core";
import { format } from "date-fns";

export const defaultEditorContent = (reference: string): PartialBlock[] => {
  return [
    {
      id: reference,
      type: "heading",
      props: {
        level: 2,
      },
      content: `${format(new Date(), "EEEE, MMMM do, yyyy")} Update`,
    },
    {
      type: "paragraph",
    },
  ];
};
