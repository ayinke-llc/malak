import { Badge } from "@/components/ui/badge";

const UpdateBadge = ({ status }: { status: string }) => {
  return (
    <Badge variant={"secondary"} className="text-xs">
      {status}
    </Badge>
  );
};

export default UpdateBadge;
