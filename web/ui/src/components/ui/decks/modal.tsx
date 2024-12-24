"use client";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { RiAddLine, RiUploadCloud2Line } from "@remixicon/react";
import { useState, useRef } from "react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";
import { toast } from "sonner";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import client from "@/lib/client";
import type { ServerAPIStatus } from "@/client/Api";
import type { AxiosError } from "axios";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { LIST_DECKS } from "@/lib/query-constants";

type FormData = {
  title: string;
  pdfUrl: string;
};

const schema = yup.object({
  title: yup.string().required("Title is required").min(3, "Title must be at least 3 characters"),
  pdfUrl: yup.string().required("Please upload a PDF file"),
}).required();

const truncateText = (text: string, maxLength: number) => {
  if (text.length <= maxLength) return text;
  return text.slice(0, maxLength) + "...";
};

export default function UploadDeckModal() {
  const [open, setOpen] = useState(false);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const queryClient = useQueryClient();

  const {
    register,
    handleSubmit,
    setValue,
    formState: { errors },
    reset
  } = useForm<FormData>({
    resolver: yupResolver(schema)
  });

  const uploadMutation = useMutation({
    mutationFn: (file: File) => client.images.uploadImage({ image_body: file }),
    onSuccess: ({ data }) => {
      setValue("pdfUrl", data.url);
      toast.success("File uploaded successfully");
    },
    onError: (err: AxiosError<ServerAPIStatus>) => {
      let msg = err.message;
      if (err.response !== undefined) {
        msg = err.response.data.message;
      }
      toast.error(msg);
    }
  });

  const createDeckMutation = useMutation({
    mutationFn: (data: FormData) => client.decks.decksCreate({
      title: data.title,
      deck_url: data.pdfUrl,
    }),
    onSuccess: ({ data }) => {
      toast.success(data.message);
      setOpen(false);
      reset();
      setSelectedFile(null);
      queryClient.invalidateQueries({ queryKey: [LIST_DECKS] });
    },
    onError: (err: AxiosError<ServerAPIStatus>) => {
      let msg = err.message;
      if (err.response !== undefined) {
        msg = err.response.data.message;
      }
      toast.error(msg);
    }
  });

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      if (file.type !== "application/pdf") {
        toast.error("Please upload a PDF file");
        return;
      }
      setSelectedFile(file);
      // Set the title to the file name without the extension
      const fileName = file.name.replace(/\.[^/.]+$/, "");
      setValue("title", fileName);

      uploadMutation.mutate(file);
    }
  };

  const handleDrop = async (e: React.DragEvent<HTMLLabelElement>) => {
    e.preventDefault();
    const file = e.dataTransfer.files?.[0];
    if (file && file.type === "application/pdf") {
      setSelectedFile(file);
      const fileName = file.name.replace(/\.[^/.]+$/, "");
      setValue("title", fileName);

      uploadMutation.mutate(file);
    } else {
      toast.error("Please upload a PDF file");
    }
  };

  const handleDragOver = (e: React.DragEvent<HTMLLabelElement>) => {
    e.preventDefault();
  };

  const onSubmit = async (data: FormData) => {
    createDeckMutation.mutate(data);
  };

  const isSubmitting = uploadMutation.isPending || createDeckMutation.isPending;

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button
          type="button"
          variant="default"
          className="whitespace-nowrap bg-zinc-100 text-zinc-900 hover:bg-zinc-200"
        >
          <RiAddLine className="mr-2 h-4 w-4" />
          Upload Deck
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="text-zinc-100">Upload Company Deck</DialogTitle>
          <DialogDescription className="text-zinc-400">
            Upload a PDF file for your company.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)} className="mt-6 space-y-6">
          <div className="space-y-2">
            <Label htmlFor="title" className="text-zinc-100">Deck Title</Label>
            <Input
              id="title"
              type="text"
              {...register("title")}
              className="bg-zinc-800 border-zinc-700 text-zinc-100 placeholder:text-zinc-500"
              placeholder="Enter deck title"
            />
            {errors.title && (
              <p className="text-sm text-red-500">{errors.title.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="deck" className="text-zinc-100">Deck File</Label>
            <div className="flex items-center justify-center w-full">
              <label
                htmlFor="deck"
                className="flex flex-col items-center justify-center w-full h-32 border-2 border-dashed rounded-lg cursor-pointer border-zinc-700 bg-zinc-800/50 hover:bg-zinc-800"
                onDrop={handleDrop}
                onDragOver={handleDragOver}
              >
                <div className="flex flex-col items-center justify-center pt-5 pb-6">
                  <RiUploadCloud2Line className="w-8 h-8 mb-3 text-zinc-400" />
                  {selectedFile ? (
                    <div className="text-center">
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <p className="mb-2 text-sm text-zinc-400">
                              Selected: <span className="font-medium text-zinc-300">{truncateText(selectedFile.name, 40)}</span>
                            </p>
                          </TooltipTrigger>
                          <TooltipContent>
                            <p>{selectedFile.name}</p>
                          </TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                      <p className="text-xs text-zinc-500">Click or drag to change file</p>
                    </div>
                  ) : (
                    <>
                      <p className="mb-2 text-sm text-zinc-400">
                        <span className="font-semibold">Click to upload</span> or drag and drop
                      </p>
                      <p className="text-xs text-zinc-500">PDF files only (MAX. 100MB)</p>
                    </>
                  )}
                </div>
                <Input
                  {...register("pdfUrl")}
                  ref={fileInputRef}
                  id="deck"
                  type="file"
                  className="hidden"
                  accept=".pdf,application/pdf"
                  onChange={handleFileChange}
                />
              </label>
            </div>
            {errors.pdfUrl && (
              <p className="text-sm text-red-500">{errors.pdfUrl.message}</p>
            )}
          </div>

          <div className="flex justify-end gap-3">
            <Button
              type="button"
              variant="ghost"
              onClick={() => {
                setOpen(false);
                reset();
                setSelectedFile(null);
              }}
              className="text-zinc-100 hover:text-zinc-200 hover:bg-zinc-800"
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={isSubmitting}
              className="bg-zinc-100 text-zinc-900 hover:bg-zinc-200 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isSubmitting ? "Uploading..." : "Upload"}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
} 
