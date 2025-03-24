"use client"

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { RiAddLine, RiUploadCloud2Line } from "@remixicon/react";
import { useState } from "react";
import { toast } from "sonner";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";

const truncateText = (text: string, maxLength: number) => {
  if (text.length <= maxLength) return text;
  return text.slice(0, maxLength) + "...";
};

export default function CSVUploadModal() {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [uploadSuccess, setUploadSuccess] = useState(false);

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      if (file.type !== "text/csv") {
        toast.error("Please upload a CSV file");
        setUploadSuccess(false);
        return;
      }
      setSelectedFile(file);
      setUploadSuccess(false);
      // Here you would typically handle the CSV file
      // For now, we'll just show a success message
      setUploadSuccess(true);
      toast.success("CSV file uploaded successfully");
    }
  };

  const handleDrop = async (e: React.DragEvent<HTMLLabelElement>) => {
    e.preventDefault();
    const file = e.dataTransfer.files?.[0];
    if (file && file.type === "text/csv") {
      setSelectedFile(file);
      setUploadSuccess(false);
      // Here you would typically handle the CSV file
      // For now, we'll just show a success message
      setUploadSuccess(true);
      toast.success("CSV file uploaded successfully");
    } else {
      toast.error("Please upload a CSV file");
      setUploadSuccess(false);
    }
  };

  const handleDragOver = (e: React.DragEvent<HTMLLabelElement>) => {
    e.preventDefault();
  };

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button
          type="button"
          variant="default"
          className="whitespace-nowrap gap-1"
        >
          <RiAddLine />
          Import CSV
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="text-foreground">Import Contacts from CSV</DialogTitle>
          <DialogDescription className="text-muted-foreground">
            Upload a CSV file containing your contacts. The CSV should have the following columns:
            email, first_name, last_name, phone, company
          </DialogDescription>
        </DialogHeader>

        <div className="mt-6">
          <Label htmlFor="csv" className="text-foreground">CSV File</Label>
          <div className="flex items-center justify-center w-full">
            <label
              htmlFor="csv"
              className="flex flex-col items-center justify-center w-full h-32 border-2 border-dashed rounded-lg cursor-pointer border-input bg-background/50 hover:bg-hover transition-colors duration-200"
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
                    <p className="text-xs text-muted-foreground">CSV files only</p>
                  </>
                )}
              </div>
              <input
                id="csv"
                type="file"
                className="hidden"
                accept=".csv"
                onChange={handleFileChange}
              />
            </label>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
} 