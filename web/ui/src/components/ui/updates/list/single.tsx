import { MalakUpdate } from "@/client/Api";
import { Badge } from "@/components/Badge";
import { Button } from "@/components/Button";
import { Divider } from "@/components/Divider";
import { RiMoreLine, RiPushpinLine } from "@remixicon/react";

const SingleUpdate = (update: MalakUpdate) => {
  return (
    <>
      <div key={update.id} className="flex items-center justify-between p-2 hover:bg-accent rounded-lg transition-colors">
        <div className="flex flex-col space-y-1">
          <div className="flex items-center space-x-2">
            <h3 className="font-semibold">{update.title}</h3>
            <Badge variant="error" className="text-xs">
              {update.status}
            </Badge>
          </div>
          <p className="text-sm text-muted-foreground">{update.created_at}</p>
        </div>
        <div className="flex space-x-2">
          <Button variant="ghost" size="icon" aria-label="Pin update">
            <RiPushpinLine className="h-4 w-4" />
          </Button>
          <Button variant="ghost" size="icon" aria-label="More options">
            <RiMoreLine className="h-4 w-4" />
          </Button>
        </div>
      </div>
      <Divider />
    </>
  )
}

export default SingleUpdate;
