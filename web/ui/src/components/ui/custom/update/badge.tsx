import { Badge } from "@/components/ui/badge";

const UpdateBadge = ({ status }: { status: string }) => {
  return (
    <Badge
      variant={status === 'sent' ? 'default' : 'secondary'}
      className="text-xs">
      {status}
    </Badge>
  );
};

export default UpdateBadge;
