import CopyToClipboard from "react-copy-to-clipboard";
import { Button } from "@/components/ui/button";
import { RiFileCopyLine } from "@remixicon/react";
import { toast } from "sonner";

const Copy = ({ text, onCopyText }: { text: string, onCopyText?: string }) => {

  const copiedTextNotifcation = onCopyText ?? "copied to clipboard"

  return (
    <CopyToClipboard text={text}
      onCopy={() => toast.success(copiedTextNotifcation)}>
      <Button variant="ghost" size="icon">
        <RiFileCopyLine className="w-4 h-4" />
      </Button>
    </CopyToClipboard>
  )
}

export default Copy;
