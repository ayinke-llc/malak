import "@blocknote/core/fonts/inter.css";
import "@blocknote/mantine/style.css";
import "@blocknote/shadcn/style.css";
import {
  DefaultReactSuggestionItem,
  getDefaultReactSlashMenuItems,
  SuggestionMenuController,
  useCreateBlockNote,
} from "@blocknote/react";
import {
  BlockNoteEditor,
  filterSuggestionItems,
} from "@blocknote/core"
import { BlockNoteView } from "@blocknote/shadcn";
import { defaultEditorContent } from "./default-value";
import fileUploader from "./image-upload";

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
