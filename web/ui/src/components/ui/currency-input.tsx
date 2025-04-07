import { NumericFormat, NumericFormatProps } from "react-number-format";
import { Input } from "./input";
import { cn } from "@/lib/utils";

interface CurrencyInputProps extends Omit<NumericFormatProps, "customInput"> {
  className?: string;
}

export function CurrencyInput({ className, ...props }: CurrencyInputProps) {
  return (
    <NumericFormat
      customInput={Input}
      thousandSeparator=","
      prefix="$"
      className={cn(className)}
      decimalScale={2}
      allowNegative={false}
      {...props}
    />
  );
} 
