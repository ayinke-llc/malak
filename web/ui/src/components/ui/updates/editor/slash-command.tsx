import { createSuggestionItems } from "novel/extensions";
import { Command, renderItems } from "novel/extensions";
import { uploadFn } from "./image-upload";
import {
  RiCodeLine, RiH1, RiH2, RiH3, RiImageLine,
  RiListOrdered, RiListUnordered, RiQuoteText,
  RiSquareLine, RiText
} from "@remixicon/react";

export const suggestionItems = createSuggestionItems([
  {
    title: "Text",
    description: "Just start typing with plain text.",
    searchTerms: ["p", "paragraph"],
    icon: <RiText size={18} />,
    command: ({ editor, range }) => {
      editor.chain().
        focus().
        deleteRange(range).
        toggleNode("paragraph", "paragraph").
        run();
    },
  },
  {
    title: "To-do List",
    description: "Track tasks with a to-do list.",
    searchTerms: ["todo", "task", "list", "check", "checkbox"],
    icon: <RiSquareLine size={18} />,
    command: ({ editor, range }) => {
      editor.chain().
        focus().
        deleteRange(range).
        toggleTaskList().
        run();
    },
  },
  {
    title: "Heading 1",
    description: "Big section heading.",
    searchTerms: ["title", "big", "large"],
    icon: <RiH1 size={18} />,
    command: ({ editor, range }) => {
      editor.chain().
        focus().
        deleteRange(range).
        setNode("heading", { level: 1 }).
        run();
    },
  },
  {
    title: "Heading 2",
    description: "Medium section heading.",
    searchTerms: ["subtitle", "medium"],
    icon: <RiH2 size={18} />,
    command: ({ editor, range }) => {
      editor.chain().
        focus().
        deleteRange(range).
        setNode("heading", { level: 2 }).
        run();
    },
  },
  {
    title: "Heading 3",
    description: "Small section heading.",
    searchTerms: ["subtitle", "small"],
    icon: <RiH3 size={18} />,
    command: ({ editor, range }) => {
      editor.chain().
        focus().
        deleteRange(range).
        setNode("heading", { level: 3 }).
        run();
    },
  },
  {
    title: "Bullet List",
    description: "Create a simple bullet list.",
    searchTerms: ["unordered", "point"],
    icon: <RiListUnordered size={18} />,
    command: ({ editor, range }) => {
      editor.chain().
        focus().
        deleteRange(range).
        toggleBulletList().
        run();
    },
  },
  {
    title: "Numbered List",
    description: "Create a list with numbering.",
    searchTerms: ["ordered"],
    icon: <RiListOrdered size={18} />,
    command: ({ editor, range }) => {
      editor.chain().
        focus().
        deleteRange(range).
        toggleOrderedList().
        run();
    },
  },
  {
    title: "Quote",
    description: "Capture a quote.",
    searchTerms: ["blockquote"],
    icon: <RiQuoteText size={18} />,
    command: ({ editor, range }) =>
      editor.chain().
        focus().
        deleteRange(range).
        toggleNode("paragraph", "paragraph").
        toggleBlockquote().
        run(),
  },
  {
    title: "Code",
    description: "Capture a code snippet.",
    searchTerms: ["codeblock"],
    icon: <RiCodeLine size={18} />,
    command: ({ editor, range }) => editor.chain().
      focus().
      deleteRange(range).
      toggleCodeBlock().
      run(),
  },
  {
    title: "Image",
    description: "Upload an image from your computer.",
    searchTerms: ["photo", "picture", "media"],
    icon: <RiImageLine size={18} />,
    command: ({ editor, range }) => {
      editor.chain().
        focus().
        deleteRange(range).
        run();
      // upload image
      const input = document.createElement("input");
      input.type = "file";
      input.accept = "image/*";
      input.onchange = async () => {
        if (input.files?.length) {
          const file = input.files[0];
          const pos = editor.view.state.selection.from;
          uploadFn(file, editor.view, pos);
        }
      };
      input.click();
    },
  },
]);

export const slashCommand = Command.configure({
  suggestion: {
    items: () => suggestionItems,
    render: renderItems,
  },
});
