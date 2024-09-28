import "@blocknote/core/fonts/inter.css";
import "@blocknote/mantine/style.css";
import type { ServerAPIStatus, ServerContentUpdateRequest } from "@/client/Api";
import { Badge } from "@/components/Badge";
import client from "@/lib/client";
import { UPDATE_CONTENT } from "@/lib/query-constants";
import {
  type Block,
  type BlockNoteEditor,
  filterSuggestionItems,
} from "@blocknote/core";
import { BlockNoteView } from "@blocknote/mantine";
import {
  type DefaultReactSuggestionItem,
  SuggestionMenuController,
  getDefaultReactSlashMenuItems,
  useCreateBlockNote,
} from "@blocknote/react";
import { useMutation } from "@tanstack/react-query";
import type { AxiosError } from "axios";
import { useState } from "react";
import { toast } from "sonner";
import { useDebouncedCallback } from "use-debounce";
import { defaultEditorContent } from "./default-value";
import fileUploader from "./image-upload";

const getCustomSlashMenuItems = (
  editor: BlockNoteEditor,
): DefaultReactSuggestionItem[] => [...getDefaultReactSlashMenuItems(editor)];

export type EditorProps = {
  reference: string | undefined;
};

const BlockNoteJSEditor = ({ reference }: EditorProps) => {
  if (reference === undefined || reference === "") {
    return null;
  }

  const editor = useCreateBlockNote({
    initialContent: defaultEditorContent(reference),
    uploadFile: fileUploader,
  });

  const [saveStatus, setSaveStatus] = useState<"Saved" | "Unsaved" | "Storing">(
    "Saved",
  );

  const mutation = useMutation({
    mutationKey: [UPDATE_CONTENT],
    mutationFn: (data: ServerContentUpdateRequest) =>
      client.workspaces.updateContent(reference, data),
    onSuccess: () => {
      setSaveStatus("Saved");
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      let msg = err.message;
      if (err.response !== undefined) {
        msg = err.response.data.message;
      }
      toast.error(msg);
      setSaveStatus("Unsaved");
    },
    retry: false,
    gcTime: Number.POSITIVE_INFINITY,
  });

  const debouncedUpdates = useDebouncedCallback(async (blocks: Block[]) => {
    const title = blocks[0];

    if (!title) {
      toast.error(
        "updates must include a title. Please add a title using a heading ( level 2 )",
      );
      return;
    }

    if (title.type !== "heading" || title.props.level !== 2) {
      toast.error(
        "Your heading must be the first item in the editor. It serves as the title of your update.",
      );
      return;
    }

    const titleContent = title.content[0] as {
      type?: string;
      text: string;
    };

    if (titleContent.type !== "text" || !titleContent.text) {
      toast.error("Your update title can be only text");
      return;
    }

    mutation.mutate({
      title: titleContent.text,
      update: JSON.stringify(blocks),
    });
  }, 1000);

  const getVariant = (): "warning" | "error" | "success" | "neutral" => {
    switch (saveStatus) {
      case "Saved":
        return "success";
      case "Unsaved":
        return "warning";
      case "Storing":
        return "warning";
      default:
        return "neutral";
    }
  };

  return (
    <div className="relative w-full max-w-screen-lg">
      <div className="flex absolute right-5 top-5 z-10 mb-5 gap-2">
        <Badge className="uppercase" variant={getVariant()}>
          {saveStatus}
        </Badge>
      </div>
      <BlockNoteView
        editor={editor}
        theme={"light"}
        onChange={() => {
          setSaveStatus("Storing");
          debouncedUpdates(editor.document);
          setSaveStatus("Unsaved");
        }}
      >
        <SuggestionMenuController
          triggerCharacter={"/"}
          getItems={async (query) =>
            filterSuggestionItems(getCustomSlashMenuItems(editor), query)
          }
        />
      </BlockNoteView>
    </div>
  );
};

export default BlockNoteJSEditor;
