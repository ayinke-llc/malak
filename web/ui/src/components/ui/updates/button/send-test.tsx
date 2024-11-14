import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { yupResolver } from "@hookform/resolvers/yup";
import { RiEyeLine } from "@remixicon/react";
import { useState } from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
import * as yup from "yup";
import type { ButtonProps } from "./props";
import { useMutation } from "@tanstack/react-query";
import { SEND_PREVIEW_UPDATE } from "@/lib/query-constants";
import client from "@/lib/client";
import { toast } from "sonner";
import { AxiosError } from "axios";
import { ServerAPIStatus } from "@/client/Api";

type PreviewUpdateInput = {
  email: string;
};

const schema = yup
  .object({
    email: yup.string().min(5).max(50).required(),
  })
  .required();

const SendTestButton = ({ reference }: ButtonProps) => {
  const [loading, setLoading] = useState<boolean>(false);

  const {
    register,
    formState: { errors },
    handleSubmit,
    reset,
  } = useForm({
    resolver: yupResolver(schema),
  });

  const mutation = useMutation({
    mutationKey: [SEND_PREVIEW_UPDATE],
    mutationFn: (data: PreviewUpdateInput) =>
      client.workspaces.previewUpdate(reference, data),
    onSuccess: ({ data }) => {
      toast.success(data.message);
      reset();
    },
    onError(err: AxiosError<ServerAPIStatus>) {
      let msg = err.message;
      if (err.response !== undefined) {
        msg = err.response.data.message;
      }
      toast.error(msg);
    },
    retry: false,
    gcTime: Number.POSITIVE_INFINITY,
    onSettled: () => setLoading(false),
    onMutate: () => setLoading(true),
  });

  const onSubmit: SubmitHandler<PreviewUpdateInput> = (data) => {
    mutation.mutate(data);
  };

  return (
    <>
      <div className="flex justify-center">
        <Dialog>
          <DialogTrigger asChild>
            <Button
              type="button"
              variant="secondary"
              size="lg"
              className="gap-1"
            >
              <RiEyeLine size={18} />
              Preview
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-lg">
            <form onSubmit={handleSubmit(onSubmit)}>
              <DialogHeader>
                <DialogTitle>Send a test email</DialogTitle>
                <DialogDescription className="mt-1 text-sm leading-6">
                  Send a test email to preview your updates before sending to
                  your investors
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
                    <p className="mt-4 text-xs text-red-600 dark:text-red-500">
                      <span className="font-medium">
                        {errors.email.message}
                      </span>
                    </p>
                  )}
                </div>
              </DialogHeader>
              <DialogFooter className="mt-6">
                <DialogClose asChild>
                  <Button
                    type={"button"}
                    className="mt-2 w-full sm:mt-0 sm:w-fit"
                    variant="secondary"
                    loading={loading}
                  >
                    Cancel
                  </Button>
                </DialogClose>
                <Button
                  type="submit"
                  className="w-full sm:w-fit"
                  loading={loading}
                >
                  Preview
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </div>
    </>
  );
};

export default SendTestButton;
