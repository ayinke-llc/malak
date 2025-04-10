import { Button } from "@/components/ui/button";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { RiFileCopyLine } from "@remixicon/react";
import CopyToClipboard from "react-copy-to-clipboard";
import { toast } from "sonner";

const Copy = ({
  text,
  onCopyText = "Copied to clipboard",
  tooltipText = "Copy this"
}: { text: string, onCopyText?: string, tooltipText?: string }) => {

  return (
    <Tooltip>
      <TooltipTrigger>
        <CopyToClipboard text={text}
          onCopy={() => toast.success(onCopyText)}>
          <Button variant="ghost" size="icon">
            <RiFileCopyLine className="w-4 h-4" />
          </Button>
        </CopyToClipboard>
      </TooltipTrigger>
      <TooltipContent>
        {tooltipText}
      </TooltipContent>
    </Tooltip>
  )
}

export default Copy;
