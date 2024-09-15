import { Button } from "@/components/Button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/Dialog";
import { Input } from "@/components/Input";
import { Label } from "@/components/Label";
import { Switch } from "@/components/Switch";
import { yupResolver } from "@hookform/resolvers/yup";
import { RiMailSendLine } from "@remixicon/react";
import { useState } from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import * as yup from "yup";
import { ButtonProps } from "./props";
import { Option } from "lucide-react";
import CreatableSelect from 'react-select/creatable';
import { EditorBubbleItem, useEditor } from "novel";

interface Option {
  readonly label: string;
  readonly value: string;
}

type SendUpdateInput = {
  email: string
  link?: boolean
  recipients: Option[]
}

const schema = yup
  .object({
    email: yup.string().min(5).max(50).required(),
    link: yup.boolean().optional(),
  })
  .required()

const SendUpdateButton = ({ }: ButtonProps) => {

  const [loading, setLoading] = useState<boolean>(false)

  const { editor } = useEditor()

  const [options, setOptions] = useState<Option[]>([
    { value: "oops", label: "oops" },
    { value: "test", label: "test" }
  ]);

  const [value, setValue] = useState<Option[]>([]);

  const createOption = (input: string): Option => {
    return { value: input, label: input }
  }

  const createNewContact = (inputValue: string) => {
    setLoading(true);
    setTimeout(() => {
      const newOption = createOption(inputValue);
      setLoading(false);
      setOptions((prev) => [...prev, newOption]);
      setValue((prev) => [...prev, newOption]);
    }, 1000);
  };

  const {
    register,
    formState: { errors },
    handleSubmit,
  } = useForm({
    resolver: yupResolver(schema),
  })

  const onSubmit: SubmitHandler<SendUpdateInput> = (data) => {
    setLoading(true)
    console.log(data)
  }

  return (
    <>
      <div className="flex justify-center">
        <Dialog>
          <DialogTrigger asChild>
            <Button type="submit"
              size="lg"
              variant="primary"
              className="gap-1">
              <RiMailSendLine size={18} />
              Send
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-lg">
            <form onSubmit={handleSubmit(onSubmit)}>
              <DialogHeader>
                <DialogTitle>
                  Send this update
                </DialogTitle>
                <DialogDescription className="mt-1 text-sm leading-6">
                  An email will be sent to all selected contacts immediately.
                  Please re-verify your content is ready and good to go
                </DialogDescription>

                <div className="mt-4">
                  <CreatableSelect
                    isMulti
                    isClearable
                    isDisabled={loading}
                    isLoading={loading}
                    onChange={(newValue) => {
                      setValue(newValue)
                    }}
                    onCreateOption={createNewContact}
                    options={options}
                    value={value}
                  />
                </div>
                <div className="mt-4">
                  <Label htmlFor="workspace-name" className="font-medium">
                    Email address
                  </Label>
                  <Input
                    id="email"
                    placeholder="yo@lanre.wtf"
                    className="mt-2"
                    type="email"
                    {...register("email")}
                  />
                  {errors.email && (
                    < p className="mt-4 text-xs text-red-600 dark:text-red-500">
                      <span className="font-medium">
                        {errors.email.message}
                      </span>
                    </p>
                  )}
                </div>

                <div className="mt-4">
                  <Switch disabled id="r3" {...register("link")} />
                  <Label disabled htmlFor="r3">
                    Coming soon. Generate a public viewable link for this update
                  </Label>
                </div>
              </DialogHeader>
              <DialogFooter className="mt-6">
                <DialogClose asChild>
                  <Button
                    type={"button"}
                    className="mt-2 w-full sm:mt-0 sm:w-fit"
                    variant="secondary"
                  >
                    Cancel
                  </Button>
                </DialogClose>
                <Button type="submit"
                  className="w-full sm:w-fit">
                  Send
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </div>
    </>
  )
}

export default SendUpdateButton;
