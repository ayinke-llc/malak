import { Badge } from "@/components/Badge";

const UpdateBadge = ({ status }: { status: string }) => {
  return (
    <Badge variant={"warning"} className="text-xs">
      {status}
    </Badge>
  );
};

export default UpdateBadge;
