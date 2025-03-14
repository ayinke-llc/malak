"use client"

import "@blocknote/core/fonts/inter.css";
import "@blocknote/mantine/style.css";
import { useTheme } from "next-themes";
import type {
  MalakBlock,
  MalakUpdate,
  ServerAPIStatus,
  ServerContentUpdateRequest,
} from "@/client/Api";
import { Badge } from "@/components/ui/badge";
import client from "@/lib/client";
import { UPDATE_CONTENT } from "@/lib/query-constants";
import {
  BlockNoteSchema,
  defaultBlockSpecs,
  filterSuggestionItems,
  PartialBlock,
  insertOrUpdateBlock,
  type BlockSchemaFromSpecs
} from "@blocknote/core";
import { BlockNoteView } from "@blocknote/mantine";
import {
  SuggestionMenuController,
  getDefaultReactSlashMenuItems,
  useCreateBlockNote
} from "@blocknote/react";
import { useMutation } from "@tanstack/react-query";
import type { AxiosError } from "axios";
import { useState } from "react";
import { toast } from "sonner";
import { useDebouncedCallback } from "use-debounce";
import { defaultEditorContent } from "./default-value";
import fileUploader from "./image-upload";
import { Alert } from "./blocks/alert";
import { RiAlertLine, RiBarChartLine } from "@remixicon/react";
import { Dashboard } from "./blocks/dashboard";
import { Chart } from "./blocks/chart";

const schema = BlockNoteSchema.create({
  blockSpecs: {
    ...defaultBlockSpecs,

    // custom blocks
    alert: Alert,
    chart: Chart,
    // dashboards disabled because 
    // 1. insanely hard to chart them. and that is because a dashboard can contain both pie and bar charts
    // 2. it is really really hard. there is no way to make the dashboard toggleable in email 
    // so it will just be so hard to read other contents below the dashboard. imagine you hve a dashbaord with 10 charts inside
    // . Endless scroll before you get to the content beneath and rest of the update 
    // 3. Might make more sense to just make the dashboard shareable via link instead. Add the link and call it a day. 
    // there are already plans to add dashboard sharing. 
    //
    // But i think 3 is the most sane path but i will sleep over it
    // dashboard: Dashboard,
  },
});

type EditorBlock = BlockSchemaFromSpecs<typeof schema.blockSpecs>;

const insertChart = (editor: typeof schema.BlockNoteEditor) => ({
  title: "Chart",
  onItemClick: () => {
    insertOrUpdateBlock(editor, {
      type: "chart",
    });
  },
  group: "Data",
  icon: <RiBarChartLine />,
});

const insertAlert = (editor: typeof schema.BlockNoteEditor) => ({
  title: "Alert",
  onItemClick: () => {
    insertOrUpdateBlock(editor, {
      type: "alert",
    });
  },
  aliases: [
    "alert",
    "notification",
    "emphasize",
    "warning",
    "error",
    "info",
    "success",
  ],
  group: "Other",
  icon: <RiAlertLine />,
});

// const insertDashboard = (editor: typeof schema.BlockNoteEditor) => ({
//   title: "Dashboard",
//   onItemClick: () => {
//     insertOrUpdateBlock(editor, {
//       type: "dashboard",
//     });
//   },
//   aliases: [
//   ],
//   group: "Data",
//   icon: <RiBarChartLine />,
// });

export type EditorProps = {
  reference: string;
  loading: boolean;
  update: MalakUpdate | undefined;
};

const BlockNoteJSEditor = ({ reference, update }: EditorProps) => {
  const { theme } = useTheme();
  const [saveStatus, setSaveStatus] = useState<"Saved" | "Unsaved" | "Storing" | "Sent">(
    update?.status == "sent" ? "Sent" : "Saved",
  );

  let initialContent = defaultEditorContent(reference);
  if (update) {
    initialContent = update?.content as PartialBlock[];
  }

  const editor = useCreateBlockNote({
    initialContent,
    schema,
    uploadFile: fileUploader,
  });

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

  const debouncedUpdates = useDebouncedCallback(async (blocks: EditorBlock[]) => {
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

    if (!titleContent) {
      toast.error("Your update must contain a title")
      return
    }

    if (titleContent.type !== "text" || !titleContent.text) {
      toast.error("Your update title can be only text");
      return;
    }

    mutation.mutate({
      title: titleContent.text,
      update: blocks as MalakBlock[],
    });
  }, 1000);

  if (reference === undefined || reference === "") {
    return null;
  }

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
        slashMenu={false}
        editor={editor}
        theme={theme as "light" | "dark"}
        editable={update?.status !== 'sent'}
        onChange={() => {
          if (update?.status === "sent") {
            return
          }
          setSaveStatus("Storing");
          debouncedUpdates(editor.document as EditorBlock[]);
          setSaveStatus("Unsaved");
        }}
      >
        <SuggestionMenuController
          triggerCharacter={"/"}
          getItems={async (query) =>
            filterSuggestionItems(
              [
                ...getDefaultReactSlashMenuItems(editor),
                insertAlert(editor),
                insertChart(editor),
                // insertDashboard(editor)
              ],
              query
            )
          }
        />
      </BlockNoteView>
    </div>
  );
};

export default BlockNoteJSEditor;
