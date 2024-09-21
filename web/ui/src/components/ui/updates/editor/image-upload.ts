import client from "@/lib/client";
import { createImageUpload } from "novel/plugins";
import { toast } from "sonner";

const onUpload = (file: File) => {

  const promise = client.images.uploadImage(
    {
      image_body: file
    },
    {
      headers: {
        "content-type": file?.type || "application/octet-stream",
      },
    })

  return new Promise((resolve, reject) => {
    toast.promise(
      promise.then(async (res) => {
        if (res.status === 200) {
          const { url } = res.data

          const image = new Image();
          image.src = url;
          image.onload = () => {
            resolve(url);
          };

          return
        }

        throw new Error("Error uploading image. Please try again.");
      }),
      {
        loading: "Uploading image...",
        success: "Image uploaded successfully.",
        error: (e) => {
          reject(e);
          return e.message;
        },
      },
    );
  });
};

export const uploadFn = createImageUpload({
  onUpload,
  validateFn: (file) => {
    if (!file.type.includes("image/")) {
      toast.error("File type not supported.");
      return false;
    }
    if (file.size / 1024 / 1024 > 10) {
      toast.error("File size too big (max 10MB).");
      return false;
    }
    return true;
  },
});
