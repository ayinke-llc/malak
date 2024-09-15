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
import { Label } from "@/components/Label";
import { Switch } from "@/components/Switch";
import { yupResolver } from "@hookform/resolvers/yup";
import { RiMailSendLine } from "@remixicon/react";
import * as EmailValidator from 'email-validator';
import { Option } from "lucide-react";
import { useEditor } from "novel";
import { useState } from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import CreatableSelect from 'react-select/creatable';
import { toast } from "sonner";
import * as yup from "yup";
import { ButtonProps } from "./props";
import { Select } from "../../custom/select/select";

interface Option {
  readonly label: string;
  readonly value: string;
}

type SendUpdateInput = {
  email: string
  link?: boolean
  recipients?: Option[]
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

    if (!EmailValidator.validate(inputValue)) {
      toast.error("you can only add an email address as a new recipient")
      setLoading(false)
      return
    }

    setTimeout(() => {
      const newOption = createOption(inputValue);
      setLoading(false);
      setOptions((prev) => [...prev, newOption]);
      setValue((prev) => [...prev, newOption]);
    }, 9000);
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
                  <Label htmlFor="select-author" className="font-medium">
                    Select Author
                  </Label>
                  <Select data={[
                    {
                      label: "Lanre Adelowo",
                      value: "lanre"
                    },
                    {
                      label: "Ayinke",
                      value: "ayinke"
                    }
                  ]} />
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
                    isLoading={loading}
                  >
                    Cancel
                  </Button>
                </DialogClose>
                <Button type="submit"
                  className="w-full sm:w-fit"
                  isLoading={loading}>
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
