import "@blocknote/core/fonts/inter.css";
import "@blocknote/mantine/style.css";
import {
  DefaultReactSuggestionItem,
  getDefaultReactSlashMenuItems,
  SuggestionMenuController,
  useCreateBlockNote,
} from "@blocknote/react";
import {
  BlockNoteEditor,
  filterSuggestionItems,
} from "@blocknote/core";
import { BlockNoteView } from "@blocknote/mantine";
import { defaultEditorContent } from "./default-value";
import fileUploader from "./image-upload";
import { ServerAPIStatus, ServerContentUpdateRequest } from "@/client/Api";
import client from "@/lib/client";
import { UPDATE_CONTENT } from "@/lib/query-constants";
import { useMutation } from "@tanstack/react-query";
import { AxiosError } from "axios";
import { useState } from "react";
import { toast } from "sonner";
import { useDebouncedCallback } from "use-debounce";
import { Converter } from "showdown";

const getCustomSlashMenuItems = (
  editor: BlockNoteEditor
): DefaultReactSuggestionItem[] => [
    ...getDefaultReactSlashMenuItems(editor),
  ];

export type EditorProps = {
  reference: string | undefined
}

const BlockNoteJSEditor = ({ reference }: EditorProps) => {

  if (reference === undefined || reference === "") {
    return null
  }

  const editor = useCreateBlockNote({
    initialContent: defaultEditorContent(reference),
    uploadFile: fileUploader,
  });

  const [saveStatus, setSaveStatus] = useState<"Saved" | "Unsaved" | "Storing">("Saved");
  const [charsCount, setCharsCount] = useState();
  const [openNode, setOpenNode] = useState(false);
  const [openColor, setOpenColor] = useState(false);
  const [openLink, setOpenLink] = useState(false);
  const [openAI, setOpenAI] = useState(false);


  const showdownConverter = new Converter()
  showdownConverter.setFlavor("github")

  const mutation = useMutation({
    mutationKey: [UPDATE_CONTENT],
    mutationFn: (data: ServerContentUpdateRequest) => client.workspaces.updateContent(reference, data),
    onSuccess: ({ data }) => {
      setSaveStatus("Saved")
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      let msg = err.message
      if (err.response !== undefined) {
        msg = err.response.data.message
      }
      toast.error(msg)
      setSaveStatus("Unsaved")
    },
    retry: false,
    gcTime: Infinity,
  })

  const debouncedUpdates = useDebouncedCallback(async () => {

    // const updatedJSON = editor.getJSON();
    //
    // const title = updatedJSON.content?.at(0)
    //
    // if (!title) {
    //   toast.error("updates must include a title. Please add a title using a heading")
    //   return
    // }
    //
    // if (title.type !== "heading") {
    //   toast.error("Your heading must be the first item in the editor. It serves as the title of your update.")
    //   return
    // }
    //
    // const content = title.content?.at(0)
    //
    // if (content?.type != "text" && content?.text?.trim().length === 0) {
    //   toast.error("Title can only include text and must not be empty")
    //   return
    // }
    //
    // mutation.mutate({
    //   update: showdownConverter.makeMarkdown(editor.getHTML()),
    // })
    //
    // setCharsCount(editor.storage.characterCount.words())
  }, 1000);

  const getVariant = (): "warning" | "error" | "success" | "neutral" => {
    switch (saveStatus) {
      case "Saved":
        return "success"
      case "Unsaved":
        return "warning"
      case "Storing":
        return "warning"
      default:
        return "neutral"
    }
  }

  return (

    <BlockNoteView
      editor={editor}
      theme={"light"}
    >
      <SuggestionMenuController
        triggerCharacter={"/"}
        getItems={async (query) =>
          filterSuggestionItems(getCustomSlashMenuItems(editor), query)
        }
      />
    </BlockNoteView>
  )
}

export default BlockNoteJSEditor;
