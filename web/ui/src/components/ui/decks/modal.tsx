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
import { useState, useRef, ChangeEvent } from "react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

export default function UploadDeckModal() {
  const [open, setOpen] = useState(false);
  const [title, setTitle] = useState("");
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setSelectedFile(file);
      // Set the title to the file name without the extension
      const fileName = file.name.replace(/\.[^/.]+$/, "");
      setTitle(fileName);
    }
  };

  const handleDrop = (e: React.DragEvent<HTMLLabelElement>) => {
    e.preventDefault();
    const file = e.dataTransfer.files?.[0];
    if (file && file.type === "application/pdf") {
      setSelectedFile(file);
      const fileName = file.name.replace(/\.[^/.]+$/, "");
      setTitle(fileName);
    }
  };

  const handleDragOver = (e: React.DragEvent<HTMLLabelElement>) => {
    e.preventDefault();
  };

  const handleSubmit = () => {
    if (!selectedFile || !title.trim()) return;
    // TODO: Handle file upload
    setOpen(false);
    setTitle("");
    setSelectedFile(null);
  };

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
        <div className="mt-6 space-y-6">
          <div className="space-y-2">
            <Label htmlFor="title" className="text-zinc-100">Deck Title</Label>
            <Input
              id="title"
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Enter deck title"
              className="bg-zinc-800 border-zinc-700 text-zinc-100 placeholder:text-zinc-500"
            />
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
                      <p className="mb-2 text-sm text-zinc-400">
                        Selected: <span className="font-medium text-zinc-300">{selectedFile.name}</span>
                      </p>
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
                  ref={fileInputRef}
                  id="deck"
                  type="file"
                  className="hidden"
                  accept=".pdf,application/pdf"
                  onChange={handleFileChange}
                />
              </label>
            </div>
          </div>

          <div className="flex justify-end gap-3">
            <Button
              type="button"
              variant="ghost"
              onClick={() => {
                setOpen(false);
                setTitle("");
                setSelectedFile(null);
              }}
              className="text-zinc-100 hover:text-zinc-200 hover:bg-zinc-800"
            >
              Cancel
            </Button>
            <Button
              type="button"
              onClick={handleSubmit}
              disabled={!selectedFile || !title.trim()}
              className="bg-zinc-100 text-zinc-900 hover:bg-zinc-200 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Upload
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
} 