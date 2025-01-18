import type { ServerAPIStatus } from "@/client/Api";
import client from "@/lib/client";
import type { AxiosError } from "axios";
import { toast } from "sonner";

const fileUploader = async (file: File) => {
  if (!file.type.includes("image/")) {
    toast.error("File type not supported.");
    return "";
  }

  if (file.size / 1024 / 1024 > 10) {
    toast.error("File size too big (max 10MB).");
    return "";
  }

  return client.upload
    .uploadImage(
      {
        image_body: file,
      },
      {
        headers: {
          "content-type": file?.type || "application/octet-stream",
        },
      },
    )
    .then(async (res) => {
      return res.data.url;
    })
    .catch((err: AxiosError<ServerAPIStatus>) => {
      let msg = err.message;
      if (err.response !== undefined) {
        msg = err.response.data.message;
      }
      toast.error(msg);
      return "";
    });
};

export default fileUploader;
