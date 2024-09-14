import { Button } from "@/components/Button";
import { RiArrowDownSLine, RiArrowDownWideLine, RiArrowDropDownLine, RiMailSendLine } from "@remixicon/react";
import { ButtonProps } from "./props";

const SendUpdateButton = ({ }: ButtonProps) => {

  return (
    <Button type="submit"
      size="lg"
      variant="primary"
      className="gap-1">
      <RiMailSendLine size={18} />
      Send
    </Button>
  )
}

export default SendUpdateButton;
