import "@blocknote/core/fonts/inter.css";
import "@blocknote/mantine/style.css";
import type {
  MalakBlockContent,
  MalakUpdate,
  ServerAPIStatus,
  ServerContentUpdateRequest,
} from "@/client/Api";
import { Badge, type badgeVariants } from "@/components/ui/badge";
import client from "@/lib/client";
import { UPDATE_CONTENT } from "@/lib/query-constants";
import {
  type Block,
  BlockNoteEditor,
  filterSuggestionItems,
  PartialBlock,
} from "@blocknote/core";
import { BlockNoteView } from "@blocknote/mantine";
import {
  type DefaultReactSuggestionItem,
  SuggestionMenuController,
  getDefaultReactSlashMenuItems,
} from "@blocknote/react";
import { useMutation } from "@tanstack/react-query";
import type { AxiosError } from "axios";
import { useMemo, useState } from "react";
import { toast } from "sonner";
import { useDebouncedCallback } from "use-debounce";
import { defaultEditorContent } from "./default-value";
import fileUploader from "./image-upload";

const getCustomSlashMenuItems = (
  editor: BlockNoteEditor,
): DefaultReactSuggestionItem[] => {
  return [
    ...getDefaultReactSlashMenuItems(editor).filter((item) => {
      const exclude = ["Video", "Audio", "File"];
      return !exclude.includes(item.title);
    }),
  ];
};

export type EditorProps = {
  reference: string;
  loading: boolean;
  update: MalakUpdate | undefined;
};

const BlockNoteJSEditor = ({ reference, update }: EditorProps) => {
  if (reference === undefined || reference === "") {
    return null;
  }

  const [saveStatus, setSaveStatus] = useState<"Saved" | "Unsaved" | "Storing">(
    "Saved",
  );

  let initialContent = defaultEditorContent(reference);

  if (update) {
    initialContent = update?.content as PartialBlock[];
  }

  const editor = useMemo(() => {
    return BlockNoteEditor.create({
      initialContent,
      uploadFile: fileUploader,
    });
  }, [initialContent]);

  const mutation = useMutation({
    mutationKey: [UPDATE_CONTENT],
    mutationFn: async (data: ServerContentUpdateRequest) => {
      return client.workspaces.updateContent(reference, data);
    },
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
      update: blocks as MalakBlockContent[],
    });
  }, 1000);

  const getVariant = (): "default" | "secondary" | "destructive" | "outline" => {
    switch (saveStatus) {
      case "Saved":
        return "default";
      case "Unsaved":
        return "destructive";
      case "Storing":
        return "secondary";
      default:
        return "default";
    }
  };

  return (
    <div className="relative w-full max-w-screen pt-10">
      <div className="flex absolute right-5 top-5 z-10 mb-15 gap-2">
        <Badge className="uppercase" variant={getVariant()}>
          {saveStatus}
        </Badge>
      </div>
      <BlockNoteView
        editor={editor}
        theme={"dark"}
        editable={update?.status !== 'sent'}
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
