"use client";

import type { ServerAPIStatus } from "@/client/Api";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import client from "@/lib/client";
import { LIST_DECKS, UPLOAD_DECK } from "@/lib/query-constants";
import { yupResolver } from "@hookform/resolvers/yup";
import { RiAddLine, RiUploadCloud2Line } from "@remixicon/react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { AxiosError } from "axios";
import { useRef, useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import * as yup from "yup";

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
    mutationKey: [UPLOAD_DECK],
    mutationFn: (file: File) => client.uploads.uploadDeck({ image_body: file }),
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
          className="whitespace-nowrap"
        >
          <RiAddLine className="mr-2 h-4 w-4" />
          Upload Deck
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="text-foreground">Upload Company Deck</DialogTitle>
          <DialogDescription className="text-muted-foreground">
            Upload a PDF file for your company.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)} className="mt-6 space-y-6">
          <div className="space-y-2">
            <Label htmlFor="title" className="text-foreground">Deck Title</Label>
            <Input
              id="title"
              type="text"
              {...register("title")}
              className="bg-background border-input text-foreground placeholder:text-muted-foreground"
              placeholder="Enter deck title"
            />
            {errors.title && (
              <p className="text-sm text-red-500">{errors.title.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="deck" className="text-foreground">Deck File</Label>
            <div className="flex items-center justify-center w-full">
              <label
                htmlFor="deck"
                className="flex flex-col items-center justify-center w-full h-32 border-2 border-dashed rounded-lg cursor-pointer border-input bg-background/50 hover:bg-accent"
                onDrop={handleDrop}
                onDragOver={handleDragOver}
              >
                <div className="flex flex-col items-center justify-center pt-5 pb-6">
                  <RiUploadCloud2Line className="w-8 h-8 mb-3 text-muted-foreground" />
                  {selectedFile ? (
                    <div className="text-center">
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <p className="mb-2 text-sm text-muted-foreground">
                              Selected: <span className="font-medium text-foreground">{truncateText(selectedFile.name, 40)}</span>
                            </p>
                          </TooltipTrigger>
                          <TooltipContent>
                            <p>{selectedFile.name}</p>
                          </TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                      <p className="text-xs text-muted-foreground">Click or drag to change file</p>
                    </div>
                  ) : (
                    <>
                      <p className="mb-2 text-sm text-muted-foreground">
                        <span className="font-semibold">Click to upload</span> or drag and drop
                      </p>
                      <p className="text-xs text-muted-foreground">PDF files only (MAX. 100MB)</p>
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
              className="text-foreground hover:text-foreground hover:bg-accent"
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={isSubmitting}
              className="bg-primary text-primary-foreground hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isSubmitting ? "Uploading..." : "Upload"}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
} 
