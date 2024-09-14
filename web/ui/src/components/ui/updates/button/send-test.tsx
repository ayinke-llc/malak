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
import { RiEyeLine } from "@remixicon/react";
import { useState } from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import * as yup from "yup";
import { ButtonProps } from "./props";

type PreviewUpdateInput = {
  email: string
  link?: boolean
}

const schema = yup
  .object({
    email: yup.string().min(5).max(50).required(),
    link: yup.boolean().optional(),
  })
  .required()

const SendTestButton = ({ }: ButtonProps) => {

  const [loading, setLoading] = useState<boolean>(false)

  const {
    register,
    formState: { errors },
    handleSubmit,
  } = useForm({
    resolver: yupResolver(schema),
  })

  const onSubmit: SubmitHandler<PreviewUpdateInput> = (data) => {
    setLoading(true)
  }

  return (
    <>
      <div className="flex justify-center">
        <Dialog>
          <DialogTrigger asChild>
            <Button type="button"
              variant="secondary" size="lg" className="gap-1">
              <RiEyeLine size={18} />
              Preview
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-lg">
            <form onSubmit={handleSubmit(onSubmit)}>
              <DialogHeader>
                <DialogTitle>
                  Send a test email
                </DialogTitle>
                <DialogDescription className="mt-1 text-sm leading-6">
                  Send a test email to preview your updates before
                  sending to your investors
                </DialogDescription>

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
                  <Switch disabled id="r3" />
                  <Label disabled htmlFor="r3">
                    Coming soon. Generate a test preview link you can share
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

export default SendTestButton;
