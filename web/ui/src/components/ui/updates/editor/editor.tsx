"use client";

import { ServerAPIStatus, ServerContentUpdateRequest } from "@/client/Api";
import { Badge } from "@/components/Badge";
import client from "@/lib/client";
import { UPDATE_CONTENT } from "@/lib/query-constants";
import { useMutation } from "@tanstack/react-query";
import { AxiosError } from "axios";
import {
  EditorCommand,
  EditorCommandEmpty,
  EditorCommandItem,
  EditorCommandList,
  EditorContent,
  type EditorInstance,
  EditorRoot
} from "novel";
import { handleCommandNavigation } from "novel/extensions";
import { handleImageDrop, handleImagePaste } from "novel/plugins";
import { useState } from "react";
import { toast } from "sonner";
import { useDebouncedCallback } from "use-debounce";
import { defaultEditorContent } from "./default-value";
import { defaultExtensions } from "./extensions";
import GenerativeMenuSwitch from "./generative/generative-menu-switch";
import { uploadFn } from "./image-upload";
import { ColorSelector } from "./selectors/color-selector";
import { LinkSelector } from "./selectors/link-selector";
import { MathSelector } from "./selectors/math-selector";
import { NodeSelector } from "./selectors/node-selector";
import { TextButtons } from "./selectors/text-buttons";
import { slashCommand, suggestionItems } from "./slash-command";
import { Separator } from "./ui/separator";
import { Converter } from "showdown";

export type EditorProps = {
  reference: string | undefined
}

const NovelEditor = ({ reference }: EditorProps) => {

  if (reference === undefined || reference === "") {
    return null
  }

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

  const debouncedUpdates = useDebouncedCallback(async (editor: EditorInstance) => {

    const updatedJSON = editor.getJSON();

    const title = updatedJSON.content?.at(0)

    if (!title) {
      toast.error("updates must include a title. Please add a title using a heading")
      return
    }

    if (title.type !== "heading") {
      toast.error("Your heading must be the first item in the editor. It serves as the title of your update.")
      return
    }

    const content = title.content?.at(0)

    if (content?.type != "text" && content?.text?.trim().length === 0) {
      toast.error("Title can only include text and must not be empty")
      return
    }

    mutation.mutate({
      update: showdownConverter.makeMarkdown(editor.getHTML()),
    })

    setCharsCount(editor.storage.characterCount.words())
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
    <div className="relative w-full max-w-screen-lg">
      <div className="flex absolute right-5 top-5 z-10 mb-5 gap-2">
        <Badge className="uppercase" variant={getVariant()}>{saveStatus}</Badge>
        <Badge className={charsCount ? "uppercase px-2 py-1" : "hidden"} variant="neutral">{charsCount} Words</Badge>
      </div>
      <EditorRoot>
        <EditorContent
          initialContent={defaultEditorContent}
          extensions={[...defaultExtensions, slashCommand]}
          className="relative min-h-[500px] w-full max-w-screen-lg border-muted bg-background sm:mb-[calc(20vh)] sm:rounded-lg sm:border sm:shadow-lg"
          editorProps={{
            handleDOMEvents: {
              keydown: (_view, event) => handleCommandNavigation(event),
            },
            handlePaste: (view, event) => handleImagePaste(view, event, uploadFn),
            handleDrop: (view, event, _slice, moved) => handleImageDrop(view, event, moved, uploadFn),
            attributes: {
              class:
                "prose prose-lg dark:prose-invert prose-headings:font-title font-default focus:outline-none max-w-full",
            },
          }}
          onUpdate={({ editor }) => {
            setSaveStatus("Storing")
            debouncedUpdates(editor);
            setSaveStatus("Unsaved");
          }}
        >
          <EditorCommand className="bg-white dark:bg-indigo-500 z-50 h-auto max-h-[330px] overflow-y-auto rounded-md border border-muted bg-background px-1 py-2 shadow-md transition-all">
            <EditorCommandEmpty className="px-2 text-muted-foreground">No results</EditorCommandEmpty>
            <EditorCommandList>
              {suggestionItems.map((item) => (
                <EditorCommandItem
                  value={item.title}
                  onCommand={(val) => {
                    if (item.command !== undefined) {
                      item.command(val)
                    }
                  }}
                  className="flex w-full items-center space-x-2 rounded-md px-2 py-1 text-left text-sm hover:bg-accent aria-selected:bg-accent"
                  key={item.title}
                >
                  <div className="flex h-10 w-10 items-center justify-center rounded-md border border-muted bg-background">
                    {item.icon}
                  </div>
                  <div>
                    <p className="font-medium">{item.title}</p>
                    <p className="text-xs text-muted-foreground">{item.description}</p>
                  </div>
                </EditorCommandItem>
              ))}
            </EditorCommandList>
          </EditorCommand>

          <GenerativeMenuSwitch open={openAI} onOpenChange={setOpenAI} withAI={false}>
            <Separator orientation="vertical" />
            <NodeSelector open={openNode} onOpenChange={setOpenNode} />
            <Separator orientation="vertical" />

            <LinkSelector open={openLink} onOpenChange={setOpenLink} />
            <Separator orientation="vertical" />
            <MathSelector />
            <Separator orientation="vertical" />
            <TextButtons />
            <Separator orientation="vertical" />
            <ColorSelector open={openColor} onOpenChange={setOpenColor} />
          </GenerativeMenuSwitch>
        </EditorContent>
      </EditorRoot>
    </div >
  );
};

export default NovelEditor;
