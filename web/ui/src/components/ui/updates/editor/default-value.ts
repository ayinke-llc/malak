import { format } from 'date-fns';

export const defaultEditorContent = {
  type: "doc",
  content: [
    {
      type: "heading",
      attrs: { level: 2 },
      content: [{ type: "text", text: `${format(new Date(), "EEEE, MMMM do, yyyy")} Update` }],
    }
  ]
};
